locals {
  lambda_image_uris = {
    analysis = "${var.aws_account_id}.dkr.ecr.${var.region}.amazonaws.com/${var.project_name}:analysis-lambda"
    api_req  = "${var.aws_account_id}.dkr.ecr.${var.region}.amazonaws.com/${var.project_name}:api-req-lambda"
    scraper  = "${var.aws_account_id}.dkr.ecr.${var.region}.amazonaws.com/${var.project_name}:scraper-lambda"
  }
}

module "iam" {
    source = "./modules/iam"

    region = var.region
    aws_account_id = var.aws_account_id
    project_name = var.project_name
    s3_bucket_name = var.s3_bucket_name
    sqs_queue_name = module.sqs_queue.sqs_name
}

module "analysis-lambda" {
    source = "./modules/lambda/analysis-lambda"

    project_name = var.project_name
    s3_bucket_name = var.s3_bucket_name
    sqs_queue_url = module.sqs_queue.sqs_url
    openai_api_key = var.openai_api_key
    lambda_image_uris = local.lambda_image_uris
    lambda_memory_sizes = var.lambda_memory_sizes
    lambda_timeouts = var.lambda_timeouts
    lambda_iam_role_arns = {
        analysis = module.iam.lambda_scraper_analysis_role_arn
    }
}

module "api-req-lambda" {
    source = "./modules/lambda/api-req-lambda"

    project_name = var.project_name
    lambda_image_uris = local.lambda_image_uris
    lambda_memory_sizes = var.lambda_memory_sizes
    lambda_timeouts = var.lambda_timeouts
    lambda_iam_role_arns = {
        api_req = module.iam.api_request_lambda_role_arn
    }
}

module "scraper-lambda" {
    source = "./modules/lambda/scraper-lambda"

    project_name = var.project_name
    s3_bucket_name = var.s3_bucket_name
    sqs_queue_url = module.sqs_queue.sqs_url
    lambda_image_uris = local.lambda_image_uris
    lambda_memory_sizes = var.lambda_memory_sizes
    lambda_timeouts = var.lambda_timeouts
    lambda_iam_role_arns = {
        scraper = module.iam.lambda_scraper_role_arn
    }
}

module "api_gateway" {
    source = "./modules/api_gateway"

    lambda_invoke_arn = module.api-req-lambda.invoke_arn
    lambda_function_name = module.api-req-lambda.api_req_lambda_name
}

module "sqs_queue" {
    source = "./modules/sqs"

    aws_account_id = var.aws_account_id
    region = var.region
    analysis_lambda_name = module.analysis-lambda.analysis_lambda_name
}

module "s3_bucket" {
    source = "./modules/s3"

    bucket_name = var.s3_bucket_name
}
