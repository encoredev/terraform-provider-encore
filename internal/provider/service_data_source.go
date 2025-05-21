package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewService() datasource.DataSource {
	return NewEncoreDataSource(
		"need.Service",
		"service",
		"Encore provisioned service information",
		"Service")
}

type Service struct {
	ComputeInstance `graphql:"compute"`
	Route           `graphql:"route"`
}

type Route struct {
	K8sClusterIP K8sClusterIP `graphql:"... on K8sClusterIP"`
}

func (a *Route) GetDocs() map[string]string {
	return map[string]string{
		"k8s_cluster_ip": "The cluster IP of the service. Set if the service is a Kubernetes service",
	}
}

type ComputeInstance struct {
	GcpCloudRun              GCPCloudRun              `graphql:"... on GCPCloudRun"`
	AwsFargateTaskDefinition AWSFargateTaskDefinition `graphql:"... on AWSFargateTaskDefinition"`
	K8sContainer             `graphql:"... on K8sContainer"`
}

func (a *ComputeInstance) GetDocs() map[string]string {
	return map[string]string{
		"gcp_cloud_run":               "The Cloud Run service. Set if the service is a Google Cloud Run service",
		"aws_fargate_task_definition": "The Fargate task definition. Set if the service is an AWS Fargate service",
	}
}

type GCPCloudRun struct {
	SelfLink               string                    `tf:"id"`
	ServerlessVpcConnector GCPServerlessVpcConnector `graphql:"serverlessVPCConnector"`
	ServiceAccount         GCPServiceAccount
	Subnet                 GCPSubnet
}

func (a *GCPCloudRun) GetDocs() map[string]string {
	return map[string]string{
		"id":                       "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) of the Cloud Run service in the form of `projects/{project}/locations/{location}/services/{service}`",
		"serverless_vpc_connector": "The serverless VPC connector. Set if the service is a Google Cloud Run service with a serverless VPC connector",
		"subnet":                   "The subnet the Cloud Run service is associated with. Set if the service is a Google Cloud Run service with Direct VPC Access",
		"service_account":          "The GCP service account of the Cloud Run service",
	}
}

type GCPServerlessVpcConnector struct {
	SelfLink string `tf:"id"`
	Network  GCPNetwork
}

func (a *GCPServerlessVpcConnector) GetDocs() map[string]string {
	return map[string]string{
		"id": "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) of the serverless VPC connector in the form of `projects/{project}/locations/{location}/connectors/{connector}`",
	}
}

type GCPServiceAccount struct {
	SelfLink string `tf:"id"`
}

func (a *GCPServiceAccount) GetDocs() map[string]string {
	return map[string]string{
		"id": "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) of the service account in the form of `projects/{project}/serviceAccounts/{service_account}`",
	}

}

type AWSFargateTaskDefinition struct {
	Arn           string
	Service       AWSFargateService
	TaskRole      AWSRole
	ExecutionRole AWSRole
	VPC           AWSVPC
}

func (a *AWSFargateTaskDefinition) GetDocs() map[string]string {
	return map[string]string{
		"arn":            "The [arn](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the Fargate task definition",
		"service":        "The Fargate service the task definition is associated with",
		"task_role":      "The task role of the Fargate task definition",
		"execution_role": "The execution role of the Fargate task definition",
		"vpc":            "The VPC the Fargate Service is associated with",
	}
}

type AWSRole struct {
	Arn string
}

func (a *AWSRole) GetDocs() map[string]string {
	return map[string]string{
		"arn": "The [arn](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the role",
	}

}

type AWSFargateService struct {
	Arn            string
	Cluster        AWSFargateCluster
	Subnets        []AWSSubnet
	SecurityGroups []AWSSecurityGroup
}

func (a *AWSFargateService) GetDocs() map[string]string {
	return map[string]string{
		"arn":             "The [arn](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the Fargate service",
		"cluster":         "The Fargate cluster the service is associated with",
		"subnets":         "The subnets the Fargate service is associated with",
		"security_groups": "The security groups the Fargate service is associated with",
	}
}

type AWSFargateCluster struct {
	Arn string
}

func (a *AWSFargateCluster) GetDocs() map[string]string {
	return map[string]string{
		"arn": "The [arn](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the Fargate cluster",
	}
}

type K8sData struct {
	Name string
}

func (a *K8sData) GetDocs() map[string]string {
	return map[string]string{
		"name": "The name of the Kubernetes resource",
	}
}

type K8sContainer struct {
	K8sDeployment K8sDeployment `graphql:"deployment"`
}

func (a *K8sContainer) GetDocs() map[string]string {
	return map[string]string{
		"k8s_deployment": "The deployment the service is part of",
	}
}

type K8sDeployment struct {
	K8sData        `graphql:"data"`
	Namespace      K8sNamespace
	ServiceAccount K8sServiceAccount
}

func (a *K8sDeployment) GetDocs() map[string]string {
	return map[string]string{
		"namespace":       "The namespace the deployment is part of",
		"service_account": "The service account of the deployment",
	}

}

type K8sServiceAccount struct {
	K8sData             `graphql:"data"`
	K8sWorkloadIdentity `graphql:"workloadIdentity"`
}

type K8sWorkloadIdentity struct {
	GcpServiceAccount GCPServiceAccount `graphql:"... on GCPServiceAccount"`
	AwsRole           AWSRole           `graphql:"... on AWSRole"`
}

func (a *K8sWorkloadIdentity) GetDocs() map[string]string {
	return map[string]string{
		"gcp_service_account": "The GCP service account the K8s service account is mapped to. Set if the workload identity is a GCP service account",
		"aws_role":            "The AWS role the K8s service account is mapped to. Set if the workload identity is an AWS role",
	}
}

type K8sNamespace struct {
	K8sData    `graphql:"data"`
	K8sCluster `graphql:"cluster"`
}

type K8sCluster struct {
	GcpGke GCPK8sCluster `graphql:"... on GCPK8sCluster"`
	AwsEks AWSK8sCluster `graphql:"... on AWSK8sCluster"`
}

func (a *K8sCluster) GetDocs() map[string]string {
	return map[string]string{
		"gcp_gke": "The GCP GKE cluster the namespace is part of. Set if the cluster is a GCP GKE cluster",
		"aws_eks": "The AWS EKS cluster the namespace is part of. Set if the cluster is an AWS EKS cluster",
	}
}

type GCPK8sCluster struct {
	SelfLink       string `tf:"id"`
	Network        GCPNetwork
	ServiceAccount GCPServiceAccount
	NodePools      []GCPK8sNodePool
}

func (a *GCPK8sCluster) GetDocs() map[string]string {
	return map[string]string{
		"id":              "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) of the GKE cluster in the form of `projects/{project}/locations/{location}/clusters/{cluster}`",
		"network":         "The network the GKE cluster is part of",
		"service_account": "The GCP service account of the GKE cluster",
		"node_pools":      "The node pools of the GKE cluster",
	}
}

type GCPK8sNodePool struct {
	SelfLink string `tf:"id"`
}

func (a *GCPK8sNodePool) GetDocs() map[string]string {
	return map[string]string{
		"id": "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) of the node pool in the form of `projects/{project}/locations/{location}/clusters/{cluster}/nodePools/{node_pool}`",
	}
}

type AWSK8sCluster struct {
	Arn           string
	Subnets       []AWSSubnet
	SecurityGroup AWSSecurityGroup
	Role          AWSRole
	VPC           AWSVPC
}

func (a *AWSK8sCluster) GetDocs() map[string]string {
	return map[string]string{
		"arn":            "The [arn](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the EKS cluster",
		"subnet_group":   "The subnet group the EKS cluster is part of",
		"security_group": "The security group the EKS cluster is part of",
		"role":           "The role of the EKS cluster",
		"vpc":            "The VPC the EKS cluster is part of",
	}
}

type K8sClusterIP struct {
	K8sData `graphql:"data"`
}
