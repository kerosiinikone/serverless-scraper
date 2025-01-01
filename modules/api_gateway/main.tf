resource "aws_api_gateway_rest_api" "api" {
  name        = "Scraper API Request"
}

resource "aws_api_gateway_resource" "my_resource" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "myresource"
}

resource "aws_api_gateway_method" "get" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.my_resource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "api-req-lambda-integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.my_resource.id
  http_method = aws_api_gateway_method.get.http_method

  integration_http_method = "GET"
  type                    = "AWS_PROXY"
  uri                     = module.lambda.api_req_lambda_arn
}

resource "aws_api_gateway_deployment" "my_deployment" {
  depends_on = [
    aws_api_gateway_integration.api-req-lambda-integration,
  ]

  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name  = "dev"
}