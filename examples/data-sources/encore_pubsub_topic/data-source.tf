data "encore_pubsub_topic" "aws" {
  name = "test"
}

data "encore_pubsub_topic" "gcp" {
  name = "test"
  env  = "gcp"
}

output "aws_topic" {
  value = data.encore_pubsub_topic.aws.name
}

output "gcp_topic" {
  value = data.encore_pubsub_topic.gcp.gcp_pubsub.self_link
}