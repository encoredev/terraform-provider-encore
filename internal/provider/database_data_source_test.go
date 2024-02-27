package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAWSRDS() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "database_name", "todo"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "aws_rds.arn", "arn:aws:rds:region:account:db"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "aws_rds.vpc.id", "vpc"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "aws_rds.subnet_group.arn", "arn:aws:rds:region:account:subgrp:app-env"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "aws_rds.security_group.id", "sg"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "aws_rds.parameter_group.arn", "arn:aws:rds:region:account:pg:rds-instance"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "aws_rds.subnet_group.arn", "arn:aws:rds:region:account:subgrp:app-env"),
		testAWSSubnets("data.encore_sql_database.database", "aws_rds.subnet_group"),
	)
}

func testGCPCloudSQL() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "database_name", "todo"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "gcp_cloud_sql.id", "projects/app-env/regions/northamerica-northeast1/instances/app-env"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "gcp_cloud_sql.network.id", "projects/app-env/global/networks/default"),
		resource.TestCheckResourceAttr("data.encore_sql_database.database", "gcp_cloud_sql.ssl_cert.fingerprint", "fingerprint"),
	)
}

func TestDatabaseDataSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testV6ProviderFactories,
		Steps: []resource.TestStep{
			testStepForEnv(
				"eks",
				testDatabaseDataSourceConfig,
				testAWSRDS(),
			),
			testStepForEnv(
				"fargate",
				testDatabaseDataSourceConfig,
				testAWSRDS(),
			),
			testStepForEnv(
				"cloudrun",
				testDatabaseDataSourceConfig,
				testGCPCloudSQL(),
			),
			testStepForEnv(
				"gke",
				testDatabaseDataSourceConfig,
				testGCPCloudSQL(),
			),
		},
	})
}

const testDatabaseDataSourceConfig = `
provider "encore" {
	auth_key = "test"
	env = "%s"
}

data "encore_sql_database" "database" {
    name = "todo"
}
`
