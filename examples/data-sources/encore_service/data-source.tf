data "encore_service" "service" {
  name = "my-service"
  env  = "my-env"
}

output "aws_fargate" {
  value = {
    "taskdef_arn" : data.encore_service.service.aws_fargate_task_definition.arn,
    "service" : data.encore_service.service.aws_fargate_task_definition.service.arn,
    "cluster" : data.encore_service.service.aws_fargate_task_definition.service.cluster.arn,
    "security_group" : data.encore_service.service.aws_fargate_task_definition.service.security_groups.0.id,
    "task_role" : data.encore_service.service.aws_fargate_task_definition.task_role.arn,
    "execution_role" : data.encore_service.service.aws_fargate_task_definition.execution_role.arn,
    "subnet" : data.encore_service.service.aws_fargate_task_definition.service.subnets.0.arn,
    "subnet_az" : data.encore_service.service.aws_fargate_task_definition.service.subnets.0.az
  }
}

output "gcp_cloudrun" {
  value = {
    "id" : data.encore_service.service.gcp_cloud_run.id,
    "service_account" : data.encore_service.service.gcp_cloud_run.service_account.id,
    "serverless_vpc_connector" : data.encore_service.service.gcp_cloud_run.serverless_vpc_connector.id,
    "network" : data.encore_service.service.gcp_cloud_run.serverless_vpc_connector.network.id
  }
}

output "k8s_deployment" {
  value = {
    "deployment" : data.encore_service.service.k8s_deployment.name,
    "namespace" : data.encore_service.service.k8s_deployment.namespace.name,
    "service_account" : data.encore_service.service.k8s_deployment.service_account.name,
    "cluster_ip" : data.encore_service.service.k8s_cluster_ip.name
  }
}

output "gcp_gke_cluster" {
  value = {
    "cluster" : data.encore_service.service.k8s_deployment.namespace.gcp_gke.id,
    "network" : data.encore_service.service.k8s_deployment.namespace.gcp_gke.network.id,
    "cluster_service_account" : data.encore_service.service.k8s_deployment.namespace.gcp_gke.service_account.id,
    "node_pool" : data.encore_service.service.k8s_deployment.namespace.gcp_gke.node_pools.0.id,
    "deployment_service_account" : data.encore_service.service.k8s_deployment.service_account.gcp_service_account.id
  }
}

output "aws_eks_cluster" {
  value = {
    "cluster" : data.encore_service.service.k8s_deployment.namespace.aws_eks.arn,
    "cluster_role" : data.encore_service.service.k8s_deployment.namespace.aws_eks.role.arn,
    "security_group" : data.encore_service.service.k8s_deployment.namespace.aws_eks.security_group.id,
    "vpc" : data.encore_service.service.k8s_deployment.namespace.aws_eks.vpc.id,
    "deployment_role" : data.encore_service.service.k8s_deployment.service_account.aws_role.arn,
    "subnet" : data.encore_service.service.k8s_deployment.namespace.aws_eks.subnets.0.arn,
    "subnet_az" : data.encore_service.service.k8s_deployment.namespace.aws_eks.subnets.0.az
  }
}