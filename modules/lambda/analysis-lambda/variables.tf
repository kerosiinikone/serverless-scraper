variable "project_name" {
    type = string
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
  default     = ""
}

variable "lambda_image_uris" {
    type = map(string)
}

variable "lambda_memory_sizes" {
    type = map(number)
}

variable "lambda_timeouts" {
    type = map(number)
}

variable "lambda_iam_role_arns" {
    type = map(string)
}