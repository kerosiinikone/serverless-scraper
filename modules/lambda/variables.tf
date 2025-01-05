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
  default = {
    analysis = ""
    api_req  = ""
    scraper  = ""
  }
}

variable "lambda_memory_sizes" {
  type = map(number)
  default = {
    analysis = 256
    api_req  = 128
    scraper  = 256
  }
}

variable "lambda_timeouts" {
  type = map(number)
  default = {
    analysis = 300
    api_req  = 60
    scraper  = 300
  }
}

variable "lambda_iam_role_arns" {
  type = map(string)
  default = {
    analysis = ""
    api_req  = ""
    scraper  = ""
  }
}

variable "s3_bucket_name" {
  type        = string
}

variable "sqs_queue_url" {
  type        = string
}

variable "openai_api_key" {
  type        = string
  default = ""
}