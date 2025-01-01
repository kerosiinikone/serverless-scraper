output "analysis_lambda_arn" {
  value       = aws_lambda_function.analysis-lambda.arn
  description = "ARN of the analysis Lambda function"
}

output "analysis_lambda_name" {
    value = aws_lambda_function.analysis-lambda.function_name
    description = "Name of the analysis Lambda function"
}