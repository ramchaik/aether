variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "eks_role_arn" {
  description = "ARN of the IAM role for EKS"
  type        = string
}

variable "node_role_arn" {
  description = "ARN of the IAM role for EKS nodes"
  type        = string
}

variable "db_username" {
  description = "Username for the RDS instance"
  type        = string
}

variable "db_password" {
  description = "Password for the RDS instance"
  type        = string
}

variable "s3_bucket_name" {
  description = "Name of the S3 bucket"
  type        = string
}

variable "sqs_queue_name" {
  description = "Name of the SQS queue"
  type        = string
}

variable "kinesis_stream_name" {
  description = "Name of the Kinesis stream"
  type        = string
}

variable "argocd_admin_password" {
  description = "Admin password for ArgoCD"
  type        = string
  sensitive   = true
}

variable "argocd_repo_url" {
  description = "Git repository URL for ArgoCD applications"
  type        = string
}

variable "argocd_repo_path" {
  description = "Path in the Git repository containing Kubernetes manifests"
  type        = string
}

variable "argocd_repo_branch" {
  description = "Branch of the Git repository to track"
  type        = string
  default     = "main"
}

variable "ssh_private_key_path" {
  description = "Path to the SSH private key file"
  type        = string
  sensitive   = true
}

variable "aws_access_key_id" {
  type      = string
  sensitive = true
  default   = ""
}

variable "aws_secret_access_key" {
  type      = string
  sensitive = true
  default   = ""
}

variable "aws_session_token" {
  type      = string
  sensitive = true
  default   = ""
}

variable "public_clerk_publishable_key" {
  description = "Publishable key for Clerk"
  type        = string
}

variable "clerk_secret_key" {
  description = "Secret key for Clerk"
  type        = string
}

variable "grafana_admin_password" {
  description = "Admin password for Grafana"
  type        = string
  sensitive   = true
}