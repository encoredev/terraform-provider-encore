package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAWSFargateGateway() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testAWSFargateService("data.encore_gateway.gateway"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.arn", "arn:aws:elasticloadbalancing:region:account:loadbalancer/app/app-env"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.0.arn", "arn:aws:elasticloadbalancing:region:account:listener/app/app-env/l1"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.0.port", "80"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.0.protocol", "HTTP"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.1.arn", "arn:aws:elasticloadbalancing:region:account:listener/app/app-env/l2"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.1.port", "443"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.1.protocol", "HTTPS"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.2.arn", "arn:aws:elasticloadbalancing:region:account:listener/app/app-env/l3"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.2.port", "54355"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "aws_alb.listeners.2.protocol", "HTTPS"),
	)
}

func testEKSGateway() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testEKSService("data.encore_gateway.gateway", "api-gateway"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "k8s_ingress.name", "res"),
	)
}

func testGKEGateway() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testGKEService("data.encore_gateway.gateway", "api-gateway"),
		resource.TestCheckResourceAttr("data.encore_gateway.gateway", "k8s_ingress.name", "res-16or00pus0nak4albtkg-encore-aws-gateway-com"),
	)
}

func TestGatewayDataSource(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testV6ProviderFactories,
		Steps: []resource.TestStep{
			testStepForEnv(
				"eks",
				testGatewayDataSourceConfig,
				testEKSGateway(),
			),
			testStepForEnv(
				"fargate",
				testGatewayDataSourceConfig,
				testAWSFargateGateway(),
			),
			testStepForEnv(
				"cloudrun",
				testGatewayDataSourceConfig,
				testGCPCloudRun("data.encore_gateway.gateway", "api-gateway"),
			),
			testStepForEnv(
				"gke",
				testGatewayDataSourceConfig,
				testGKEGateway(),
			),
		},
	})
}

const testGatewayDataSourceConfig = `
provider "encore" {
	auth_key = "test"
	env = "%s"
}

data "encore_gateway" "gateway" {
    name = "api-gateway"
}
`
