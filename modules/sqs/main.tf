resource "aws_sqs_queue" "sqs_queue" {
  content_based_deduplication       = false
  delay_seconds                     = 0
  fifo_queue                        = false
  kms_data_key_reuse_period_seconds = 300
  max_message_size                  = 262144
  visibility_timeout_seconds        = 120
  message_retention_seconds         = 345600
  name                              = var.queue_name

  policy = <<POLICY
  {
    "Id": "__default_policy_ID",
    "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "__owner_statement",
          "Action": "sqs:*",
          "Effect": "Allow",
          "Resource": "arn:aws:sqs:${var.region}:${var.aws_account_id}:${var.queue_name}",
          "Principal": {
            "AWS": [
              "${var.aws_account_id}"
            ]
          }
        }
      ]
    }
  POLICY
}

resource "aws_lambda_event_source_mapping" "queue-lambda-trigger" {
  batch_size                         = 1
  bisect_batch_on_function_error     = false
  enabled                            = true
  event_source_arn                   = aws_sqs_queue.sqs_queue.arn
  function_name                      = var.analysis_lambda_name
}

