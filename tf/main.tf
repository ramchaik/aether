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
  subnet_ids      = aws_subnet.public[*].id
  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }
  instance_types = ["t3.small", "t3a.small", "t3.medium", "t3a.medium"]
  capacity_type  = "SPOT"
  depends_on     = [aws_eks_cluster.main]
}


resource "aws_db_subnet_group" "main" {
  name       = "main"
  subnet_ids = aws_subnet.public[*].id
}

resource "aws_security_group" "rds" {
  name        = "rds-sg"
  description = "Allow inbound traffic from EKS"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "Allow inbound from EKS"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
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
  parameter_group_name = "default.postgres15"
  skip_final_snapshot  = true

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name

  tags = {
    Name = "main-db"
  }
}


resource "aws_s3_bucket" "aether" {
  bucket = var.s3_bucket_name
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

resource "aws_ecr_repository" "aether" {
  name                 = var.ecr_repository_name
  image_tag_mutability = "MUTABLE"
  tags = {
    Name = var.ecr_repository_name
  }
}

output "endpoint" {
  value = aws_eks_cluster.main.endpoint
}

output "kubeconfig-certificate-authority-data" {
  value = aws_eks_cluster.main.certificate_authority[0].data
}

output "rds_endpoint" {
  value = aws_db_instance.main.endpoint
}

output "s3_bucket_name" {
  value = aws_s3_bucket.aether.bucket
}

output "sqs_queue_url" {
  value = aws_sqs_queue.aether.url
}

output "ecr_repository_url" {
  value = aws_ecr_repository.aether.repository_url
}
