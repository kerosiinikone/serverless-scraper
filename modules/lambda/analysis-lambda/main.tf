resource "aws_lambda_function" "analysis-lambda" {
  architectures = ["x86_64"]

  environment {
    variables = {
      AWS_BUCKET = var.s3_bucket_name
      OPENAI_API = var.openai_api_key
      QUEUE_URL  = var.sqs_queue_url
      RUNTIME    = "lambda"
    }
  }

  ephemeral_storage {
    size = 512
  }

  function_name = "${var.project_name}-analysis-lambda"
  image_uri     = var.lambda_image_uris["analysis"]

  logging_config {
    log_format = "Text"
    log_group  = "/aws/lambda/${var.project_name}-analysis-lambda"
  }

  memory_size                    = var.lambda_memory_sizes["analysis"]
  package_type                   = "Image"
  reserved_concurrent_executions = -1
  role                           = var.lambda_iam_role_arns["analysis"]
  skip_destroy                   = false
  timeout                        = var.lambda_timeouts["analysis"]

  tracing_config {
    mode = "PassThrough"
  }
}