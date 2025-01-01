module "analysis-lambda" {
    source = "./modules/lambda/analysis-lambda"

    project_name = var.project_name
    s3_bucket_name = module.s3_bucket.bucket_id
    sqs_queue_url = module.sqs_queue.queue_url
    openai_api_key = var.openai_api_key
    lambda_image_uris = var.lambda_image_uris
    lambda_memory_sizes = var.lambda_memory_sizes
    lambda_timeouts = var.lambda_timeouts
    lambda_iam_role_arns = var.lambda_iam_role_arns
}

module "api-req-lambda" {
    source = "./modules/lambda/api-req-lambda"

    project_name = var.project_name
    lambda_image_uris = var.lambda_image_uris
    lambda_memory_sizes = var.lambda_memory_sizes
    lambda_timeouts = var.lambda_timeouts
    lambda_iam_role_arns = var.lambda_iam_role_arns
}

module "scraper-lambda" {
    source = "./modules/lambda/scraper-lambda"

    project_name = var.project_name
    s3_bucket_name = module.s3_bucket.bucket_id
    sqs_queue_url = module.sqs_queue.sqs_url
    lambda_image_uris = var.lambda_image_uris
    lambda_memory_sizes = var.lambda_memory_sizes
    lambda_timeouts = var.lambda_timeouts
    lambda_iam_role_arns = var.lambda_iam_role_arns
}

resource "aws_lambda_event_source_mapping" "queue-lambda-trigger" {
  batch_size                         = "1"
  bisect_batch_on_function_error     = "false"
  enabled                            = "true"
  event_source_arn                   = module.sqs_queue.sqs_arn
  function_name                      = module.scraper-lambda.function_name
  maximum_batching_window_in_seconds = "0"
  maximum_record_age_in_seconds      = "0"
  maximum_retry_attempts             = "0"
  parallelization_factor             = "0"
  tumbling_window_in_seconds         = "0"
}
