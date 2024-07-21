provider "aws" {
  region = var.region
}

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

resource "aws_eks_node_group" "main" {
  cluster_name    = aws_eks_cluster.main.name
  node_group_name = "main-node-group"
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
  depends_on     = [aws_eks_cluster.main]
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
      sleep 60  # Wait for resources to be deleted
    EOT
  }

  depends_on = [aws_eks_node_group.main]
}

resource "null_resource" "cleanup_resources" {
  triggers = {
    vpc_id         = aws_vpc.main.id
    nat_eip_id     = aws_eip.nat.id
    s3_bucket_name = aws_s3_bucket.aether.id
    region         = var.region
  }

  provisioner "local-exec" {
    when    = destroy
    command = <<-EOT
      # Release Elastic IP
      aws ec2 release-address --allocation-id ${self.triggers.nat_eip_id} --region ${self.triggers.region}

      # Delete all resources in the VPC
      aws ec2 describe-instances --filters "Name=vpc-id,Values=${self.triggers.vpc_id}" --query 'Reservations[].Instances[].InstanceId' --output text --region ${self.triggers.region} | xargs -r aws ec2 terminate-instances --instance-ids --region ${self.triggers.region}
      aws ec2 describe-nat-gateways --filter "Name=vpc-id,Values=${self.triggers.vpc_id}" --query 'NatGateways[].NatGatewayId' --output text --region ${self.triggers.region} | xargs -r -n1 aws ec2 delete-nat-gateway --nat-gateway-id --region ${self.triggers.region}
      aws ec2 describe-network-interfaces --filters "Name=vpc-id,Values=${self.triggers.vpc_id}" --query 'NetworkInterfaces[].NetworkInterfaceId' --output text --region ${self.triggers.region} | xargs -r -n1 aws ec2 delete-network-interface --network-interface-id --region ${self.triggers.region}

      # Wait for resources to be deleted
      sleep 300
    EOT
  }
  depends_on = [null_resource.delete_eks_resources]
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
  description = "Allow outbound traffic from Lambda to RDS"
  vpc_id      = aws_vpc.main.id

  egress {
    description     = "Allow outbound to RDS"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.rds.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "rds" {
  name        = "rds-sg"
  description = "Allow inbound traffic from EKS"
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

resource "aws_s3_bucket" "private_bucket" {
  bucket        = var.private_s3_bucket_name
  force_destroy = true

  tags = {
    Name = var.private_s3_bucket_name
  }
}

resource "aws_s3_bucket_ownership_controls" "private_bucket_ownership" {
  bucket = aws_s3_bucket.private_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "private_bucket" {
  bucket = aws_s3_bucket.private_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_sqs_queue" "aether" {
  name = var.sqs_queue_name
}

resource "aws_ecr_repository" "aether" {
  name                 = var.ecr_repository_name
  image_tag_mutability = "MUTABLE"
  tags = {
    Name = var.ecr_repository_name
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
