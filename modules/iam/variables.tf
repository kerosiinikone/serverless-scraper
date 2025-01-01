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

variable "s3_bucket_name" {
    type = string
}

variable "sqs_queue_name" {
    type = string
}