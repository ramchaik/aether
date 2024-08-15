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
}

resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = element(["10.1.3.0/24", "10.1.4.0/24"], count.index)
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags = {
    Name = "private-subnet-${count.index + 1}"
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
}

resource "aws_nat_gateway" "main" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public[0].id

  depends_on = [aws_internet_gateway.main]

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

resource "kubernetes_namespace" "aether" {
  metadata {
    name = "aether"
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
    AWS_ACCESS_KEY_ID     = sensitive(coalesce(var.aws_access_key_id, data.external.env_vars.result["AWS_ACCESS_KEY_ID"]))
    AWS_SECRET_ACCESS_KEY = sensitive(coalesce(var.aws_secret_access_key, data.external.env_vars.result["AWS_SECRET_ACCESS_KEY"]))
    AWS_SESSION_TOKEN     = sensitive(coalesce(var.aws_session_token, data.external.env_vars.result["AWS_SESSION_TOKEN"]))
  }

  depends_on = [kubernetes_namespace.aether]
}

data "external" "env_vars" {
  program = ["sh", "-c", "echo '{\"AWS_ACCESS_KEY_ID\":\"'$AWS_ACCESS_KEY_ID'\",\"AWS_SECRET_ACCESS_KEY\":\"'$AWS_SECRET_ACCESS_KEY'\",\"AWS_SESSION_TOKEN\":\"'$AWS_SESSION_TOKEN'\"}'"]
}

resource "kubernetes_secret" "clerk_keys" {
  metadata {
    name      = "clerk-keys"
    namespace = "aether"
  }

  data = {
    PUBLIC_CLERK_PUBLISHABLE_KEY = var.public_clerk_publishable_key
    CLERK_SECRET_KEY             = var.clerk_secret_key
  }

  depends_on = [kubernetes_namespace.aether]
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

resource "aws_kms_key" "rds_encryption_key" {
  description             = "KMS key for RDS encryption"
  deletion_window_in_days = 7
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

  storage_encrypted = true
  kms_key_id        = aws_kms_key.rds_encryption_key.arn

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

resource "helm_release" "reloader" {
  name       = "reloader"
  repository = "https://stakater.github.io/stakater-charts"
  chart      = "reloader"
  namespace  = "kube-system"
  version    = "v1.0.48"

  set {
    name  = "reloader.watchGlobally"
    value = "false"
  }

  depends_on = [aws_eks_node_group.general, aws_eks_node_group.forge]
}

# kubectl port-forward svc/argocd-server -n argocd 8080:443
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

  depends_on = [
    helm_release.argocd,
    kubectl_manifest.argocd_repository,
    kubernetes_namespace.aether
  ]
}

resource "kubernetes_secret" "db_url" {
  metadata {
    name      = "db-url"
    namespace = "aether"
  }

  data = {
    DATABASE_URL = "postgresql://${var.db_username}:${var.db_password}@${aws_db_instance.main.endpoint}/aether?sslmode=disable"
  }

  depends_on = [kubernetes_namespace.aether]
}

resource "null_resource" "lambda_zip" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = <<EOT
      cd ../src/logify/serverless
      make clean build
      zip -j lambda_function.zip bootstrap
      mv lambda_function.zip ../../../tf/
      cd ../../../tf
      echo "$(pwd)/lambda_function.zip" > ${path.module}/zip_file_path.txt
      echo "Created zip file at: $(pwd)/lambda_function.zip"
    EOT
  }
}

data "local_file" "zip_file_path" {
  filename   = "${path.module}/zip_file_path.txt"
  depends_on = [null_resource.lambda_zip]
}

resource "aws_lambda_function" "kinesis_consumer" {
  filename         = trimspace(data.local_file.zip_file_path.content)
  function_name    = "kinesis-consumer-lambda"
  role             = "arn:aws:iam::502413910473:role/LabRole"
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  source_code_hash = filebase64sha256(trimspace(data.local_file.zip_file_path.content))

  environment {
    variables = {
      DATABASE_URL   = "postgresql://${var.db_username}:${var.db_password}@${aws_db_instance.main.endpoint}/aether-logs?sslmode=disable"
      KINESIS_STREAM = var.kinesis_stream_name
    }
  }

  vpc_config {
    subnet_ids         = aws_subnet.private[*].id
    security_group_ids = [aws_security_group.lambda.id]
  }

  depends_on = [
    null_resource.lambda_zip,
    aws_security_group.lambda,
    aws_subnet.private
  ]
}

resource "aws_lambda_event_source_mapping" "kinesis_trigger" {
  event_source_arn  = aws_kinesis_stream.aether.arn
  function_name     = aws_lambda_function.kinesis_consumer.function_name
  starting_position = "LATEST"
}

resource "kubernetes_secret" "db_credentials" {
  metadata {
    name      = "db-credentials"
    namespace = "aether"
  }

  data = {
    DB_USERNAME = var.db_username
    DB_PASSWORD = var.db_password
  }

  depends_on = [kubernetes_namespace.aether]
}

resource "helm_release" "keda" {
  name             = "keda"
  repository       = "https://kedacore.github.io/charts"
  chart            = "keda"
  namespace        = "keda"
  create_namespace = true
  version          = "2.9.0"

  set {
    name  = "rbac.create"
    value = "true"
  }

  set {
    name  = "serviceAccount.create"
    value = "true"
  }

  depends_on = [aws_eks_node_group.general, aws_eks_node_group.forge]
}

resource "helm_release" "prometheus" {
  name             = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  namespace        = "monitoring"
  create_namespace = true

  set {
    name  = "prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues"
    value = "false"
  }
}

# kubectl port-forward svc/grafana 3000:80 -n monitoring
resource "helm_release" "grafana" {
  name             = "grafana"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "grafana"
  namespace        = "monitoring"
  create_namespace = true

  set {
    name  = "adminPassword"
    value = var.grafana_admin_password
  }
}

resource "helm_release" "karpenter" {
  name       = "karpenter"
  repository = "https://charts.karpenter.sh"
  chart      = "karpenter"
  namespace  = "karpenter"
  version    = "v0.29.0"

  create_namespace = true

  set {
    name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = "arn:aws:iam::502413910473:role/LabRole"
  }

  set {
    name  = "settings.aws.clusterName"
    value = aws_eks_cluster.main.name
  }

  set {
    name  = "settings.aws.clusterEndpoint"
    value = aws_eks_cluster.main.endpoint
  }

  set {
    name  = "settings.aws.defaultInstanceProfile"
    value = "KarpenterNodeInstanceProfile-${aws_eks_cluster.main.name}"
  }

  depends_on = [
    aws_eks_node_group.general,
    aws_eks_node_group.forge,
  ]
}

resource "kubectl_manifest" "karpenter_provisioner" {
  yaml_body = <<-YAML
apiVersion: karpenter.sh/v1alpha5
kind: Provisioner
metadata:
  name: default
spec:
  requirements:
    - key: karpenter.sh/capacity-type
      operator: In
      values: ["spot"]
  limits:
    resources:
      cpu: 1000
  providerRef:
    name: default
  consolidation:
    enabled: true
YAML

  depends_on = [helm_release.karpenter]
}

resource "kubectl_manifest" "karpenter_node_pool" {
  yaml_body = <<-YAML
apiVersion: karpenter.k8s.aws/v1alpha1
kind: AWSNodeTemplate
metadata:
  name: default
spec:
  subnetSelector:
    karpenter.sh/discovery: ${aws_eks_cluster.main.name}
  securityGroupSelector:
    karpenter.sh/discovery: ${aws_eks_cluster.main.name}
  tags:
    karpenter.sh/discovery: ${aws_eks_cluster.main.name}
YAML

  depends_on = [helm_release.karpenter]
}