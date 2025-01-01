resource "aws_lambda_function" "api-req-lambda" {
  architectures = ["x86_64"]

  environment {
    variables = {
      RUNTIME = "lambda"
    }
  }

  ephemeral_storage {
    size = 512
  }

  function_name = "${var.project_name}-api-req-lambda"
  image_uri     = var.lambda_image_uris["api_req"]

  logging_config {
    log_format = "Text"
    log_group  = "/aws/lambda/${var.project_name}-api-req-lambda"
  }

  memory_size                    = var.lambda_memory_sizes["api_req"]
  package_type                   = "Image"
  reserved_concurrent_executions = -1
  role                           = var.lambda_iam_role_arns["api_req"]
  skip_destroy                   = false
  timeout                        = var.lambda_timeouts["api_req"]

  tracing_config {
    mode = "PassThrough"
  }
}