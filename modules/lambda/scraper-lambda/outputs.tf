output "scraper_lambda_arn" {
  value       = aws_lambda_function.scraper-lambda.arn
  description = "ARN of the scraper Lambda function"
}

output "scraper_lambda_name" {
    value = aws_lambda_function.scraper-lambda.function_name
    description = "Name of the scraper Lambda function"
}