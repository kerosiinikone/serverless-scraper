output "api_req_lambda_arn" {
  value = aws_lambda_function.api-req-lambda.arn
}

output "api_req_lambda_name" {
    value = aws_lambda_function.api-req-lambda.function_name
}

output "invoke_arn" {
    value = aws_lambda_function.api-req-lambda.invoke_arn
}