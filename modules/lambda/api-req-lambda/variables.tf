variable "project_name" {
    type = string
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