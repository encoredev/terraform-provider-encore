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
	DatabaseNumber int
	RedisCluster   `graphql:"cluster"`
}

type RedisCluster struct {
	AwsRedis AWSRedisCluster `graphql:"... on AWSRedisCluster"`
	GcpRedis GCPRedisCluster `graphql:"... on GCPRedisCluster"`
}

type GCPRedisCluster struct {
	SelfLink string `tf:"id"`
	Network  GCPNetwork
}

type AWSSubnetGroup struct {
	Arn     string
	Subnets []AWSSubnet
}

type AWSRedisParameterGroup struct {
	Arn string
}

type AWSRedisCluster struct {
	Arn            string
	VPC            AWSVPC
	SubnetGroup    AWSSubnetGroup
	SecurityGroup  AWSSecurityGroup
	ParameterGroup AWSParameterGroup
}
