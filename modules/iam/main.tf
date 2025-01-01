resource "aws_iam_role" "ApiRequestLambdaRole" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description           = "Allows Lambda functions to call AWS services on your behalf."
  max_session_duration = "3600"
  name                 = "${var.project_name}-ApiRequestLambdaRole"
  path                 = "/"

  inline_policy {
    name = "dynamodb_create_empty_and_table"

    policy = jsonencode({
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "dynamodb:CreateTable",
            "dynamodb:PutItem",
            "dynamodb:DescribeTable",
            "dynamodb:ListTables",
            "dynamodb:GetItem",
            "dynamodb:UpdateItem"
          ],
          "Effect" : "Allow",
          "Resource" : [
            "arn:aws:dynamodb:${var.region}:${var.aws_account_id}:table/${var.project_name}-*"
          ],
          "Sid" : "DynamoDBAccess"
        }
      ]
    })
  }
  inline_policy {
    name = "lambda_invoke_function"

    policy = jsonencode({
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "lambda:InvokeFunction",
            "lambda:InvokeAsync"
          ],
          "Effect" : "Allow",
          "Resource" : [
            "arn:aws:lambda:${var.region}:${var.aws_account_id}:function:${var.project_name}-*"
          ],
          "Sid" : "LambdaInvoke"
        }
      ]
    })
  }
}

resource "aws_iam_role" "LambdaScraperAnalysisRole" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description           = "Allows Lambda functions to call AWS services on your behalf."
  max_session_duration = "3600"
  name                 = "${var.project_name}-LambdaScraperAnalysisRole"
  path                 = "/"

  inline_policy {
    name = "s3_list_and_get"

    policy = jsonencode({
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "s3:GetObject",
            "s3:ListBucket",
            "s3:GetBucketLocation"
          ],
          "Effect" : "Allow",
          "Resource" : [
            "arn:aws:s3:::${var.s3_bucket_name}",
            "arn:aws:s3:::${var.s3_bucket_name}/*"
          ],
          "Sid" : "S3ListAndGet"
        }
      ]
    })
  }
  inline_policy {
      name = "dynamodb_update_item"
      policy = jsonencode({
          "Version": "2012-10-17",
          "Statement": [
              {
                  "Action": [
                      "dynamodb:UpdateItem"
                  ],
                  "Effect": "Allow",
                  "Resource": [
                      "arn:aws:dynamodb:${var.region}:${var.aws_account_id}:table/${var.project_name}-*"
                  ],
                  "Sid": "DynamoDBUpdateItem"
              }
          ]
      })
  }
  inline_policy {
    name = "sqs_receive_and_delete"

    policy = jsonencode({
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "sqs:DeleteMessage",
            "sqs:ChangeMessageVisibility",
            "sqs:ReceiveMessage",
            "sqs:GetQueueAttributes"
          ],
          "Effect" : "Allow",
          "Resource" : [
            "arn:aws:sqs:${var.region}:${var.aws_account_id}:${var.sqs_queue_name}"
          ],
          "Sid" : "SQSReceiveAndDelete"
        }
      ]
    })
  }
}

resource "aws_iam_role" "LambdaScraperRole" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description           = "Allows Lambda functions to call AWS services on your behalf."
  max_session_duration = "3600"
  name                 = "${var.project_name}-LambdaScraperRole"
  path                 = "/"

  inline_policy {
    name = "s3_put"

    policy = jsonencode({
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "s3:PutObject"
          ],
          "Effect" : "Allow",
          "Resource" : [
            "arn:aws:s3:::${var.s3_bucket_name}/*"
          ],
          "Sid" : "S3Put"
        }
      ]
    })
  }
  inline_policy {
    name = "dynamodb_update_item"

    policy = jsonencode({
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "dynamodb:UpdateItem"
          ],
          "Effect" : "Allow",
          "Resource" : [
            "arn:aws:dynamodb:${var.region}:${var.aws_account_id}:table/${var.project_name}-*"
          ],
          "Sid" : "DynamoDBUpdateItem"
        }
      ]
    })
  }
  inline_policy {
    name = "sqs_send_message"

    policy = jsonencode({
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "sqs:SendMessage"
          ],
          "Effect" : "Allow",
          "Resource" : [
            "arn:aws:sqs:${var.region}:${var.aws_account_id}:${var.sqs_queue_name}"
          ],
          "Sid" : "SQSSendMessage"
        }
      ]
    })
  }
}