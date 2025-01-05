variable "queue_name" {
    type = string
    default = "analysis-queue"
}

variable "region" {
    type = string
    default = "eu-north-1"
}

variable "aws_account_id" {
    type = string
}

variable "analysis_lambda_name" {
    type = string
}