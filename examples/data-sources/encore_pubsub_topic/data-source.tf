data "encore_pubsub_topic" "aws" {
  name = "test"
  env  = "aws"
}

data "aws_iam_policy_document" "mypolicy" {
  statement {
    effect    = "Allow"
    actions   = ["sns:Publish"]
    resources = [data.encore_pubsub_topic.aws.aws_sns.arn]
  }
}



output "aws_sns" {
  value = {
    "arn" : data.encore_pubsub_topic.topic.aws_sns.arn,
  }
}

output "gcp_pubsub" {
  value = {
    "id" : data.encore_pubsub_topic.topic.gcp_pubsub.id,
  }
}
