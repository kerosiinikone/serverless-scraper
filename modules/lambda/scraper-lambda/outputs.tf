output "scraper_lambda_arn" {
  value = aws_lambda_function.scraper-lambda.arn
}

output "scraper_lambda_name" {
    value = aws_lambda_function.scraper-lambda.function_name
}