output "api_request_lambda_role_arn" {
  value       = aws_iam_role.ApiRequestLambdaRole.arn
  description = "ARN of the API Request Lambda Role"
}

output "lambda_scraper_analysis_role_arn" {
  value       = aws_iam_role.LambdaScraperAnalysisRole.arn
  description = "ARN of the Lambda Scraper Analysis Role"
}

output "lambda_scraper_role_arn" {
  value       = aws_iam_role.LambdaScraperRole.arn
  description = "ARN of the Lambda Scraper Role"
}