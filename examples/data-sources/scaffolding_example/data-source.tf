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

// Use cases:
//
// 1. Get information about the infrastructure Encore has provisioned,
//    to be able to provision additional infrastructure alongside it.
//
// 2. Propagate the provisioned infrastructure resources back to Encore,
//    for use inside the running services.
//
// 3. Set up IAM permissions for services to use the additional infrastructure.

// Solutions:
// 1. Add a TF data source for querying infra information.
// 2. Provide a way to export information about the additional infra resources,
//    for example by generating a config file that can be statically used by the service.
// 3. Support defining policies to add to services
