data "encore_pubsub_topic" "test_topic" {
  name = "test"
}

data "encore_pubsub_topic" "other_topic" {
  name = "test"
  env  = "gcp"
}

output "aws_topic" {
  value = data.encore_pubsub_topic.test_topic.aws_sns.arn
}

output "gcp_topic" {
  value = data.encore_pubsub_topic.other_topic.gcp_pubsub.relative_resource_name
}