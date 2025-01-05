resource "aws_lambda_function" "scraper-lambda" {
  architectures = ["x86_64"]

  environment {
    variables = {
      AWS_BUCKET = var.s3_bucket_name
      QUEUE_URL  = var.sqs_queue_url
      RUNTIME    = "lambda"
    }
  }

  ephemeral_storage {
    size = 512
  }

  function_name = "scraper-lambda"
  image_uri     = var.lambda_image_uris["scraper"]

  logging_config {
    log_format = "Text"
    log_group  = "/aws/lambda/scraper-lambda"
  }

  memory_size                    = var.lambda_memory_sizes["scraper"]
  package_type                   = "Image"
  reserved_concurrent_executions = -1
  role                           = var.lambda_iam_role_arns["scraper"]
  skip_destroy                   = false
  timeout                        = var.lambda_timeouts["scraper"]

  tracing_config {
    mode = "PassThrough"
  }
}