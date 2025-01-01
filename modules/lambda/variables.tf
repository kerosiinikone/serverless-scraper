variable "region" {
    type = string
}

variable "aws_account_id" {
    type = string
}

variable "project_name" {
    type = string
    default = "poc"
}

variable "lambda_image_uris" {
  type = map(string)
  description = "ECR Image URIs for each Lambda function"
  default = {
    analysis = ""
    api_req  = ""
    scraper  = ""
  }
}

variable "lambda_memory_sizes" {
  type = map(number)
  description = "Memory sizes for each Lambda function"
  default = {
    analysis = 256
    api_req  = 128
    scraper  = 256
  }
}

variable "lambda_timeouts" {
  type = map(number)
  description = "Timeouts for each Lambda function"
  default = {
    analysis = 300
    api_req  = 60
    scraper  = 300
  }
}

variable "lambda_iam_role_arns" {
  type = map(string)
  description = "IAM role ARNs for each Lambda function"
  default = {
    analysis = ""
    api_req  = ""
    scraper  = ""
  }
}

variable "s3_bucket_name" {
  type        = string
  description = "Name of the S3 bucket for Lambda environment variables"
}

variable "sqs_queue_url" {
  type        = string
  description = "URL of the SQS queue for Lambda environment variables"
}

variable "openai_api_key" {
  type        = string
  description = "OpenAI API key for the analysis-lambda"
  default = ""
}