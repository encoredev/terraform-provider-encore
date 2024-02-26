package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewCache() datasource.DataSource {
	return NewEncoreDataSource(
		"need.CacheKeyspace",
		"cache",
		"Encore provisioned cache information",
		"RedisKeyspace")
}

type RedisKeyspace struct {
	RedisCluster `graphql:"cluster"`
}

type RedisCluster struct {
	AwsRedis AWSRedisCluster `graphql:"... on AWSRedisCluster"`
	GcpRedis GCPRedisCluster `graphql:"... on GCPRedisCluster"`
}

func (r *RedisCluster) GetDocs() map[string]string {
	return map[string]string{
		"aws_redis": "Set if the Redis cluster is provisioned on AWS",
		"gcp_redis": "Set if the Redis cluster is provisioned on GCP",
	}
}

type GCPRedisCluster struct {
	SelfLink string `tf:"id"`
	Network  GCPNetwork
}

func (r *GCPRedisCluster) GetDocs() map[string]string {
	return map[string]string{
		"id":      "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) in the form of `projects/{project}/locations/{location}/instances/{instance}`",
		"network": "The network the Redis cluster is provisioned in",
	}
}

type AWSSubnetGroup struct {
	Arn     string
	Subnets []AWSSubnet
}

func (r *AWSSubnetGroup) GetDocs() map[string]string {
	return map[string]string{
		"arn":     "The [Amazon Resource Name (ARN)](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the subnet group",
		"subnets": "The subnets the resource is provisioned in",
	}
}

type AWSRedisParameterGroup struct {
	Arn string
}

func (r *AWSRedisParameterGroup) GetDocs() map[string]string {
	return map[string]string{
		"arn": "The [Amazon Resource Name (ARN)](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the parameter group",
	}
}

type AWSRedisCluster struct {
	Arn            string
	VPC            AWSVPC
	SubnetGroup    AWSSubnetGroup
	SecurityGroup  AWSSecurityGroup
	ParameterGroup AWSParameterGroup
}

func (r *AWSRedisCluster) GetDocs() map[string]string {
	return map[string]string{
		"arn":             "The [Amazon Resource Name (ARN)](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the Redis cluster",
		"vpc":             "The VPC the Redis cluster is provisioned in",
		"subnet_group":    "The subnet group the Redis cluster is provisioned in",
		"security_group":  "The security group of the Redis cluster",
		"parameter_group": "The parameter group of the Redis cluster",
	}
}
