output "sqs_id" {
  value = aws_sqs_queue.main.id
}

output "sqs_arn" {
  value = aws_sqs_queue.main.arn
}

output "sqs_url" {
  value = aws_sqs_queue.main.url
}