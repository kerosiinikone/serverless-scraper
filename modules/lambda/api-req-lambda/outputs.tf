output "api_req_lambda_arn" {
  value       = aws_lambda_function.api-req-lambda.arn
  description = "ARN of the api-req Lambda function"
}

output "api_req_lambda_name" {
    value = aws_lambda_function.api-req-lambda.function_name
    description = "Name of the api-req Lambda function"
}