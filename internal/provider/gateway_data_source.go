package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewGateway() datasource.DataSource {
	return NewEncoreDataSource(
		"need.Gateway",
		"gateway",
		"Encore provisioned gateway.",
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

func (g *Ingress) GetDocs() map[string]string {
	return map[string]string{
		"k8s_ingress": "Kubernetes Ingress. Set if the gateway is provisioned on a Kubernetes cluster.",
		"aws_alb":     "AWS Application Load Balancer. Set if the gateway is provisioned on AWS.",
	}
}

type K8sIngress struct {
	K8sData `graphql:"data"`
}

type AWSAppLoadBalancer struct {
	Arn       string
	Listeners []AWSAppLoadBalancerListener
}

func (a *AWSAppLoadBalancer) GetDocs() map[string]string {
	return map[string]string{
		"arn":       "[ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the AWS Application Load Balancer.",
		"listeners": "Listeners of the AWS Application Load Balancer.",
	}
}

type AWSAppLoadBalancerListener struct {
	Arn      string
	Port     int
	Protocol string
}

func (a *AWSAppLoadBalancerListener) GetDocs() map[string]string {
	return map[string]string{
		"arn":      "[ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the listener.",
		"port":     "Port of the listener.",
		"protocol": "Protocol of the listener.",
	}
}
