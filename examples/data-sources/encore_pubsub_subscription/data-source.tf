data "encore_pubsub_subscription" "aws" {
  name = "test"
}

data "encore_pubsub_subscription" "gcp" {
  name = "test"
  env  = "gcp"
}

output "aws_sub" {
  value = data.encore_pubsub_subscription.aws.name
}

output "gcp_sub" {
  value = data.encore_pubsub_subscription.gcp.gcp_pubsub.self_link
}