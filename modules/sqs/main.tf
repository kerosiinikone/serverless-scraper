resource "aws_sqs_queue" "sqs_queue" {
  content_based_deduplication       = "false"
  delay_seconds                     = "0"
  fifo_queue                        = "false"
  kms_data_key_reuse_period_seconds = "300"
  max_message_size                  = "262144"
  message_retention_seconds         = "345600"
  name                              = var.queue_name

  policy = <<POLICY
{
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Action": "SQS:*",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::941377157937:root"
      },
      "Resource": "arn:aws:sqs:eu-north-1:941377157937:poc-analysis-queue",
      "Sid": "__owner_statement"
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  receive_wait_time_seconds  = "0"
  sqs_managed_sse_enabled    = "true"
  visibility_timeout_seconds = "60"
}
