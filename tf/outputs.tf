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

output "forge_ecr_repository_url" {
  value = aws_ecr_repository.forge.repository_url
}

output "launchpad_ecr_repository_url" {
  value = aws_ecr_repository.launchpad.repository_url
}

output "frontstage_ecr_repository_url" {
  value = aws_ecr_repository.frontstage.repository_url
}

output "logify_ecr_repository_url" {
  value = aws_ecr_repository.logify.repository_url
}

output "proxy_ecr_repository_url" {
  value = aws_ecr_repository.proxy.repository_url
}

output "kinesis_stream_arn" {
  value = aws_kinesis_stream.aether.arn
}
