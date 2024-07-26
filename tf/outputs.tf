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

output "kinesis_stream_arn" {
  value = aws_kinesis_stream.aether.arn
}
