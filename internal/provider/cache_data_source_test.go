package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAWSRedis() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.encore_cache.cache", "aws_redis.arn", "arn:aws:elasticache:region:account:replicationgroup:app-env-cache-cluster"),
		resource.TestCheckResourceAttr("data.encore_cache.cache", "aws_redis.vpc.id", "vpc"),
		resource.TestCheckResourceAttr("data.encore_cache.cache", "aws_redis.security_group.id", "sg"),
		resource.TestCheckResourceAttr("data.encore_cache.cache", "aws_redis.parameter_group.arn", "arn:aws:elasticache:region:account:parametergroup:app-env-cache-cluster"),
		resource.TestCheckResourceAttr("data.encore_cache.cache", "aws_redis.subnet_group.arn", "arn:aws:elasticache:region:account:subnetgroup:redis"),
		testAWSSubnets("data.encore_cache.cache", "aws_redis.subnet_group"),
	)
}

func testGCPRedis() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.encore_cache.cache", "gcp_redis.id", "projects/app-env/locations/northamerica-northeast1/instances/app-env"),
		resource.TestCheckResourceAttr("data.encore_cache.cache", "gcp_redis.network.id", "projects/app-env/global/networks/default"),
	)
}

func TestCacheDataSource(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testV6ProviderFactories,
		Steps: []resource.TestStep{
			testStepForEnv(
				"eks",
				testCacheDataSourceConfig,
				testAWSRedis(),
			),
			testStepForEnv(
				"fargate",
				testCacheDataSourceConfig,
				testAWSRedis(),
			),
			testStepForEnv(
				"cloudrun",
				testCacheDataSourceConfig,
				testGCPRedis(),
			),
			testStepForEnv(
				"gke",
				testCacheDataSourceConfig,
				testGCPRedis(),
			),
		},
	})
}

const testCacheDataSourceConfig = `
provider "encore" {
	auth_key = "test"
	env = "%s"
}

data "encore_cache" "cache" {
    name = "cache-cluster"
}
`
