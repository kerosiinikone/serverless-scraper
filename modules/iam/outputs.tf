output "api_request_lambda_role_arn" {
  value = aws_iam_role.ApiRequestLambdaRole.arn
}

output "lambda_scraper_analysis_role_arn" {
  value = aws_iam_role.LambdaScraperAnalysisRole.arn
}

output "lambda_scraper_role_arn" {
  value = aws_iam_role.LambdaScraperRole.arn
}