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

type ComputeInstance struct {
	GcpCloudRun              GCPCloudRun              `graphql:"... on GCPCloudRun"`
	AwsFargateTaskDefinition AWSFargateTaskDefinition `graphql:"... on AWSFargateTaskDefinition"`
	K8sContainer             `graphql:"... on K8sContainer"`
}

type GCPCloudRun struct {
	SelfLink               string                    `tf:"id"`
	ServerlessVpcConnector GCPServerlessVpcConnector `graphql:"serverlessVPCConnector"`
	ServiceAccount         GCPServiceAccount
}

type GCPServerlessVpcConnector struct {
	SelfLink string `tf:"id"`
}

type GCPServiceAccount struct {
	SelfLink string
}

type AWSFargateTaskDefinition struct {
	Arn           string
	Service       AWSFargateService
	TaskRole      AWSRole
	ExecutionRole AWSRole
}

type AWSRole struct {
	Arn string
}

type AWSFargateService struct {
	Arn            string
	Cluster        AWSFargateCluster
	Subnets        []AWSSubnet
	SecurityGroups []AWSSecurityGroup
}

type AWSFargateCluster struct {
	Arn string
}

type K8sData struct {
	Name string
}

type K8sContainer struct {
	K8sDeployment K8sDeployment `graphql:"deployment"`
}

type K8sDeployment struct {
	K8sData        `graphql:"data"`
	Namespace      K8sNamespace
	ServiceAccount K8sServiceAccount
}

type K8sServiceAccount struct {
	K8sData             `graphql:"data"`
	K8sWorkloadIdentity `graphql:"workloadIdentity"`
}

type K8sWorkloadIdentity struct {
	GcpServiceAccount GCPServiceAccount `graphql:"... on GCPServiceAccount"`
	AwsRole           AWSRole           `graphql:"... on AWSRole"`
}

type K8sNamespace struct {
	K8sData    `graphql:"data"`
	K8sCluster `graphql:"cluster"`
}

type K8sCluster struct {
	GcpGke GCPK8sCluster `graphql:"... on GCPK8sCluster"`
	AwsEks AWSK8sCluster `graphql:"... on AWSK8sCluster"`
}

type GCPK8sCluster struct {
	SelfLink string `tf:"id"`
}

type AWSK8sCluster struct {
	Arn string
}

type K8sClusterIP struct {
	K8sData `graphql:"data"`
}
