data "encore_pubsub_subscription" "aws" {
  name = "test"
  env  = "aws"
}

data "aws_iam_policy_document" "mypolicy" {
  statement {
    effect = "Allow"
    actions = ["sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes",
    "sqs:ChangeMessageVisibility"]
    resources = [data.encore_pubsub_topic.aws.aws_sqs.arn]
  }
}
