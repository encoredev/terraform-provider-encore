terraform {
    required_providers {
        encore = {
            source = "registry.terraform.io/encore/encore"
        }
    }
}

provider "encore" {
  app_id = "my-app"
  env_name = "production"
}

data "encore_pubsub_topic" "my-topic" {}

output "arn" {
  value = data.encore_pubsub_topic.my-topic.aws.arn
}
