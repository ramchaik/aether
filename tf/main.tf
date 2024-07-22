data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "main" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name = "main-vpc"
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_subnet" "public" {
  count                   = 2
  vpc_id                  = aws_vpc.main.id
  cidr_block              = element(["10.1.1.0/24", "10.1.2.0/24"], count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  tags = {
    Name = "public-subnet-${count.index + 1}"
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = element(["10.1.3.0/24", "10.1.4.0/24"], count.index)
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags = {
    Name = "private-subnet-${count.index + 1}"
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
}

resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_eip" "nat" {
  domain = "vpc"
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_nat_gateway" "main" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public[0].id
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main.id
  }
}

resource "aws_route_table_association" "private" {
  count          = 2
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}

resource "aws_eks_cluster" "main" {
  name     = "main-cluster"
  role_arn = var.eks_role_arn
  vpc_config {
    subnet_ids = aws_subnet.public[*].id
  }
  version = "1.30"
}

resource "aws_eks_node_group" "general" {
  cluster_name    = aws_eks_cluster.main.name
  node_group_name = "general-node-group"
  node_role_arn   = var.node_role_arn
  subnet_ids      = concat(aws_subnet.public[*].id, aws_subnet.private[*].id)
  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }
  instance_types = ["t3.small", "t3a.small"]
  capacity_type  = "SPOT"
  disk_size      = 20

  labels = {
    "node-group" = "general"
  }

  depends_on = [aws_eks_cluster.main]
}

resource "aws_eks_node_group" "forge" {
  cluster_name    = aws_eks_cluster.main.name
  node_group_name = "forge-node-group"
  node_role_arn   = var.node_role_arn
  subnet_ids      = concat(aws_subnet.public[*].id, aws_subnet.private[*].id)
  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }
  instance_types = ["t3.medium", "t3a.medium"]
  capacity_type  = "SPOT"
  disk_size      = 20

  labels = {
    "node-group" = "forge"
  }

  depends_on = [aws_eks_cluster.main]
}

resource "null_resource" "create_namespace" {
  provisioner "local-exec" {
    command = <<EOT
      aws eks update-kubeconfig --name ${aws_eks_cluster.main.name} --region us-east-1
      kubectl create ns aether
    EOT
  }

  depends_on = [
    aws_eks_node_group.general,
    aws_eks_node_group.forge
  ]
}

resource "kubernetes_secret" "aws_credentials" {
  metadata {
    name      = "aws-credentials"
    namespace = "aether"
  }

  data = {
    AWS_ACCESS_KEY_ID     = var.aws_access_key_id
    AWS_SECRET_ACCESS_KEY = var.aws_secret_access_key
    AWS_SESSION_TOKEN     = var.aws_session_token
  }
}

resource "null_resource" "delete_eks_resources" {
  triggers = {
    cluster_name = aws_eks_cluster.main.name
    region       = var.region
  }

  provisioner "local-exec" {
    when    = destroy
    command = <<-EOT
      aws eks update-kubeconfig --name ${self.triggers.cluster_name} --region ${self.triggers.region}
      kubectl delete services --all
      kubectl delete deployments --all
      kubectl delete pods --all
      kubectl delete daemonsets --all
      kubectl get svc --all-namespaces -o json | jq -r '.items[] | select(.spec.type=="LoadBalancer") | .metadata.name' | xargs -I {} kubectl delete svc {}
      sleep 300  # Increased wait time for EKS resources to be deleted
    EOT
  }

  depends_on = [aws_eks_node_group.general, aws_eks_node_group.forge]
}

resource "aws_db_parameter_group" "custom_pg" {
  family = "postgres15"
  name   = "custom-pg-hba-conf"

  parameter {
    name  = "rds.force_ssl"
    value = "0"
  }

  parameter {
    name  = "log_connections"
    value = "1"
  }

  parameter {
    name  = "log_hostname"
    value = "1"
  }

  parameter {
    name  = "pgaudit.log"
    value = "all"
  }
}

resource "aws_db_subnet_group" "main" {
  name       = "main"
  subnet_ids = aws_subnet.private[*].id
}

resource "aws_security_group" "lambda" {
  name        = "lambda-sg"
  description = "Security group for Lambda functions"
  vpc_id      = aws_vpc.main.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "rds" {
  name        = "rds-sg"
  description = "Allow inbound traffic from EKS and Lambda"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Allow inbound from EKS"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_eks_cluster.main.vpc_config[0].cluster_security_group_id]
  }

  ingress {
    description     = "Allow inbound from Lambda"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.lambda.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_instance" "main" {
  identifier           = "main-db"
  engine               = "postgres"
  engine_version       = "15.7"
  instance_class       = "db.t3.micro"
  allocated_storage    = 10
  storage_type         = "gp2"
  db_name              = "aether"
  username             = var.db_username
  password             = var.db_password
  parameter_group_name = aws_db_parameter_group.custom_pg.name
  skip_final_snapshot  = true

  vpc_security_group_ids = [aws_security_group.rds.id, aws_security_group.lambda.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name

  tags = {
    Name = "main-db"
  }
}

resource "aws_s3_bucket" "aether" {
  bucket        = var.s3_bucket_name
  force_destroy = true
  tags = {
    Name = var.s3_bucket_name
  }
}

resource "aws_s3_bucket_policy" "aether_bucket_policy" {
  bucket = aws_s3_bucket.aether.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.aether.arn}/*"
      }
    ]
  })
}

resource "aws_s3_bucket_public_access_block" "aether" {
  bucket = aws_s3_bucket.aether.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_sqs_queue" "aether" {
  name = var.sqs_queue_name
}

resource "aws_ecr_repository" "forge" {
  name                 = var.ecr_forge_repository_name
  image_tag_mutability = "MUTABLE"
  tags = {
    Name = var.ecr_forge_repository_name
  }
}

resource "aws_ecr_repository" "frontstage" {
  name                 = var.ecr_frontstage_repository_name
  image_tag_mutability = "MUTABLE"
  tags = {
    Name = var.ecr_frontstage_repository_name
  }
}

resource "aws_ecr_repository" "launchpad" {
  name                 = var.ecr_launchpad_repository_name
  image_tag_mutability = "MUTABLE"
  tags = {
    Name = var.ecr_launchpad_repository_name
  }
}

resource "aws_ecr_repository" "logify" {
  name                 = var.ecr_logify_repository_name
  image_tag_mutability = "MUTABLE"
  tags = {
    Name = var.ecr_logify_repository_name
  }
}

resource "aws_ecr_repository" "proxy" {
  name                 = var.ecr_proxy_repository_name
  image_tag_mutability = "MUTABLE"
  tags = {
    Name = var.ecr_proxy_repository_name
  }
}

resource "aws_kinesis_stream" "aether" {
  name             = var.kinesis_stream_name
  shard_count      = 1
  retention_period = 24

  stream_mode_details {
    stream_mode = "PROVISIONED"
  }

  tags = {
    Name = var.kinesis_stream_name
  }
}

resource "helm_release" "argocd" {
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  namespace        = "argocd"
  create_namespace = true
  version          = "5.36.1"

  values = [
    <<-EOT
    server:
      extraArgs:
        - --insecure
    configs:
      secret:
        argocdServerAdminPassword: "${bcrypt(var.argocd_admin_password)}"
    EOT
  ]

  depends_on = [aws_eks_node_group.general, aws_eks_node_group.forge]
}

resource "kubernetes_secret" "argocd_ssh_key" {
  metadata {
    name      = "argocd-ssh-key"
    namespace = "argocd"
    labels = {
      "argocd.argoproj.io/secret-type" = "repository"
    }
  }

  data = {
    "sshPrivateKey" = file("${var.ssh_private_key_path}")
  }

  type = "Opaque"

  depends_on = [helm_release.argocd]
}

resource "kubectl_manifest" "argocd_repository" {
  yaml_body = <<-YAML
apiVersion: v1
kind: Secret
metadata:
  name: ather-repo
  namespace: argocd
  labels:
    argocd.argoproj.io/secret-type: repository
stringData:
  type: git
  url: git@github.com:ramchaik/aether.git
  sshPrivateKey: |
    ${indent(4, file("${var.ssh_private_key_path}"))}
YAML

  depends_on = [helm_release.argocd, kubernetes_secret.argocd_ssh_key]
}

resource "kubectl_manifest" "argocd_application" {
  yaml_body = <<-YAML
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: aether
  namespace: argocd
spec:
  project: default
  source:
    repoURL: git@github.com:ramchaik/aether.git
    targetRevision: "${var.argocd_repo_branch}"
    path: "${var.argocd_repo_path}"
  destination:
    server: https://kubernetes.default.svc
    namespace: aether
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
YAML

  depends_on = [helm_release.argocd, kubectl_manifest.argocd_repository]
}
