output "api_gateway_id" {
  description = "The ID of the API Gateway"
  value       = aws_api_gateway_rest_api.go_scraper_api.id
}

output "api_gateway_root_resource_id" {
  description = "The root resource ID of the API Gateway"
  value       = aws_api_gateway_rest_api.go_scraper_api.root_resource_id
}

output "api_gateway_deployment_id" {
  description = "The deployment ID of the API Gateway"
  value       = aws_api_gateway_deployment.deploy_api.id
}

output "api_gateway_stage_name" {
  description = "The stage name of the API Gateway"
  value       = aws_api_gateway_stage.dev.stage_name
}

output "api_gateway_method" {
  description = "The HTTP method of the API Gateway"
  value       = aws_api_gateway_method.get_resource.http_method
}