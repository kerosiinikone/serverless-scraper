output "sqs_id" {
  value = aws_sqs_queue.sqs_queue.id
}

output "sqs_arn" {
  value = aws_sqs_queue.sqs_queue.arn
}

output "sqs_url" {
  value = aws_sqs_queue.sqs_queue.url
}

output "sqs_name" {
  value = aws_sqs_queue.sqs_queue.name
}