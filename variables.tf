variable "region" {
    type = string
    default = "eu-north-1"
}

variable "aws_account_id" {
    type = string
}

variable "project_name" {
    type = string
    default = "poc"
}

variable "lambda_image_uris" {
  type        = map(string)
  description = "ECR Image URIs for each Lambda function"
  default = {
    analysis = ""
    api_req  = ""
    scraper  = ""
  }
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

variable "openai_api_key" {
  type        = string
  description = "OpenAI API key for the analysis-lambda"
  default     = ""
}