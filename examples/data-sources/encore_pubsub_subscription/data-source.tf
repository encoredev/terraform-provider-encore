data "encore_pubsub_subscription" "test_sub" {
  name = "test"
}

data "encore_pubsub_subscription" "other_sub" {
  name = "test"
  env  = "gcp"
}

output "aws_sub" {
  value = data.encore_pubsub_subscription.test_sub.aws_sqs.arn
}

output "gcp_sub" {
  value = data.encore_pubsub_subscription.other_sub.gcp_pubsub.relative_resource_name
}