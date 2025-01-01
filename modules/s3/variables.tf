variable "force_destroy" {
  type        = bool
  description = "Whether to allow force destroy of the bucket"
  default     = false
}

variable "bucket_name" {
    type = string
    description = "The name of the S3 bucket"
    default = "raw-scraped-data"
  
}