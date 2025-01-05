variable "region" {
    type = string
    default = "eu-north-1"
}

variable "s3_bucket_name" {
    type = string
}

variable "aws_account_id" {
    type = string
}

variable "project_name" {
    type = string
    default = "scraperdev"
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
      analysis = 120
      api_req  = 60
      scraper  = 120
    }
}

variable "openai_api_key" {
  type        = string
}