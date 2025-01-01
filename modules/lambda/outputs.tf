output "analysis_lambda_arn" {
  value       = module.analysis-lambda.analysis_lambda_arn
  description = "ARN of the analysis Lambda function"
}

output "api_req_lambda_arn" {
  value       = module.api-req-lambda.api_req_lambda_arn
  description = "ARN of the API request Lambda function"
}

output "scraper_lambda_arn" {
  value       = module.scraper-lambda.scraper_lambda_arn
  description = "ARN of the scraper Lambda function"
}

output "analysis_lambda_name" {
  value = module.analysis-lambda.analysis_lambda_name
  description = "Name of the analysis Lambda function"
}

output "api_req_lambda_name" {
    value = module.api-req-lambda.api_req_lambda_name
    description = "Name of the api-req Lambda function"
}

output "scraper_lambda_name" {
    value = module.scraper-lambda.scraper_lambda_name
    description = "Name of the scraper Lambda function"
}