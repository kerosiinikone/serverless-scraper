resource "aws_apigatewayv2_api" "go_scraper_api" {
  name          = "Scraper API"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_route" "api_route" {
  api_id    = aws_apigatewayv2_api.go_scraper_api.id
  route_key = "POST /api-req-lambda"
  target    = "integrations/${aws_apigatewayv2_integration.api_integration.id}"
}

resource "aws_apigatewayv2_integration" "api_integration" {
  api_id           = aws_apigatewayv2_api.go_scraper_api.id
  integration_type = "AWS_PROXY"
  description      = "Integration for the api-req-lambda"
  integration_method = "POST"
  integration_uri  = var.lambda_invoke_arn
}

resource "aws_apigatewayv2_deployment" "deploy_api" {
  api_id      = aws_apigatewayv2_api.go_scraper_api.id
  description = "dev"
  depends_on = [aws_apigatewayv2_route.api_route, aws_apigatewayv2_integration.api_integration]
}

resource "aws_apigatewayv2_stage" "default" {
  api_id = aws_apigatewayv2_api.go_scraper_api.id
  name = "$default"
  auto_deploy = true
  deployment_id = aws_apigatewayv2_deployment.deploy_api.id
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.go_scraper_api.execution_arn}/*/*"
}
