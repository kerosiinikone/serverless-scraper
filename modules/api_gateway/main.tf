resource "aws_api_gateway_rest_api" "go_scraper_api" {
  name        = "Scraper API Request"
}

resource "aws_api_gateway_resource" "api_resource" {
  rest_api_id = aws_api_gateway_rest_api.go_scraper_api.id
  parent_id   = aws_api_gateway_rest_api.go_scraper_api.root_resource_id
  path_part   = "api-req-lambda"
}

resource "aws_api_gateway_method" "get_resource" {
  rest_api_id   = aws_api_gateway_rest_api.go_scraper_api.id
  resource_id   = aws_api_gateway_resource.api_resource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_deployment" "deploy_api" {
  depends_on = [
    aws_api_gateway_integration.api_req_lambda_integration,
  ]

  rest_api_id = aws_api_gateway_rest_api.go_scraper_api.id
}

resource "aws_api_gateway_stage" "dev" {
  deployment_id = aws_api_gateway_deployment.deploy_api.id
  rest_api_id   = aws_api_gateway_rest_api.go_scraper_api.id
  stage_name    = "dev"
}
