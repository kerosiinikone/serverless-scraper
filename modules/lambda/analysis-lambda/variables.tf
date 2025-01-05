variable "project_name" {
    type = string
}

variable "s3_bucket_name" {
  type        = string
}

variable "sqs_queue_url" {
  type        = string
}

variable "openai_api_key" {
  type        = string
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