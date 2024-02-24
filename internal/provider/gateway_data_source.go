package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewGateway() datasource.DataSource {
	return NewEncoreDataSource(
		"need.Gateway",
		"gateway",
		"Encore provisioned gateway information",
		"Gateway")
}

type Gateway struct {
	ComputeInstance `graphql:"compute"`
	Route           `graphql:"route"`
	Ingress         `graphql:"ingress"`
}

type Ingress struct {
	K8sIngress K8sIngress         `graphql:"... on K8sIngress"`
	AwsAlb     AWSAppLoadBalancer `graphql:"... on AWSAppLoadBalancer"`
}

type K8sIngress struct {
	K8sData `graphql:"data"`
}

type AWSAppLoadBalancer struct {
	Arn       string
	Listeners []AWSAppLoadBalancerListener
}

type AWSAppLoadBalancerListener struct {
	Arn      string
	Port     int
	Protocol string
}
