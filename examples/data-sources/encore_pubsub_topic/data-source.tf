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
