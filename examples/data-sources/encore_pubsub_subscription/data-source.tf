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

output "aws_sqs" {
  value = {
    "arn" : data.encore_pubsub_subscription.subscription.aws_sns.arn,
    "topic" : data.encore_pubsub_subscription.subscription.aws_sns.topic.arn,
    "queue" : data.encore_pubsub_subscription.subscription.aws_sns.queue.arn,
    "dead_letter" : data.encore_pubsub_subscription.subscription.aws_sns.queue.dead_letter.arn
  }
}

output "gcp_pubsub" {
  value = {
    "id" : data.encore_pubsub_subscription.subscription.gcp_pubsub.id,
    "topic" : data.encore_pubsub_subscription.subscription.gcp_pubsub.topic.id,
    "dead_letter" : data.encore_pubsub_subscription.subscription.gcp_pubsub.dead_letter.id,
    "dead_letter_topic" : data.encore_pubsub_subscription.subscription.gcp_pubsub.dead_letter.topic.id
  }
}
