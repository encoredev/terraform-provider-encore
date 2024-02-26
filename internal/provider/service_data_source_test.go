package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAWSFargateService(res string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(res, "aws_fargate_task_definition.arn", "arn:aws:ecs:region:account:task-definition/app-env-encore"),
		resource.TestCheckResourceAttr(res, "aws_fargate_task_definition.service.arn", "arn:aws:ecs:region:account:service/app-env/encore"),
		resource.TestCheckResourceAttr(res, "aws_fargate_task_definition.service.cluster.arn", "arn:aws:ecs:region:account:cluster/app-env"),
		testAWSSubnets(res, "aws_fargate_task_definition.service"),
		resource.TestCheckResourceAttr(res, "aws_fargate_task_definition.service.security_groups.0.id", "sg"),
		resource.TestCheckResourceAttr(res, "aws_fargate_task_definition.task_role.arn", "arn:aws:iam::account:role/encore/app/env/app-env-encore-task-role"),
		resource.TestCheckResourceAttr(res, "aws_fargate_task_definition.execution_role.arn", "arn:aws:iam::account:role/encore/app/env/app-env-encore-execution-role"),
	)
}

func testKubernetesService(res, svcName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(res, "k8s_deployment.name", svcName),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.name", "app-env"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.service_account.name", svcName),
		resource.TestCheckResourceAttr(res, "k8s_cluster_ip.name", svcName),
	)
}

func testEKSService(res, svcName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testKubernetesService(res, svcName),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.aws_eks.arn", "arn:aws:eks:region:account:cluster/app-env"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.aws_eks.role.arn", "arn:aws:iam::account:role/encore/app/env/app-env-cache-task-role"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.aws_eks.security_group.id", "sg"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.aws_eks.vpc.id", "vpc"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.service_account.aws_role.arn", "arn:aws:iam::account:role/encore/app/env/app-env-"+svcName+"-task-role"),
		testAWSSubnets(res, "k8s_deployment.namespace.aws_eks"),
	)
}

func testGCPCloudRunService(res, svcName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testGCPCloudRun(res, svcName),
		resource.TestCheckResourceAttr(res, "gcp_cloud_run.serverless_vpc_connector.id", "projects/app-env/locations/northamerica-northeast1/connectors/appenv"),
		resource.TestCheckResourceAttr(res, "gcp_cloud_run.serverless_vpc_connector.network.id", "projects/app-env/global/networks/default"),
	)
}

func testGCPCloudRun(res, svcName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(res, "gcp_cloud_run.id", "projects/app-env/locations/northamerica-northeast1/services/"+svcName),
		resource.TestCheckResourceAttr(res, "gcp_cloud_run.service_account.id", "projects/app-env/serviceAccounts/"+svcName+"@app-env.iam.gserviceaccount.com"),
	)
}

func testGKEService(res, svcName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testKubernetesService(res, svcName),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.gcp_gke.id", "projects/app-env/locations/northamerica-northeast1/clusters/app-env"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.gcp_gke.network.id", "projects/app-env/global/networks/default"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.gcp_gke.service_account.id", "test-service-account"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.gcp_gke.node_pools.0.id", "test-node-pool"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.namespace.gcp_gke.node_pools.1.id", "test-node-pool"),
		resource.TestCheckResourceAttr(res, "k8s_deployment.service_account.gcp_service_account.id", "projects/app-env/serviceAccounts/"+svcName+"@app-env.iam.gserviceaccount.com"),
	)
}

func TestServiceDataSource(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testV6ProviderFactories,
		Steps: []resource.TestStep{
			testStepForEnv(
				"eks",
				testServiceDataSourceConfig,
				testEKSService("data.encore_service.service", "cache"),
			),
			testStepForEnv(
				"fargate",
				testServiceDataSourceConfig,
				testAWSFargateService("data.encore_service.service"),
			),
			testStepForEnv(
				"cloudrun",
				testServiceDataSourceConfig,
				testGCPCloudRunService("data.encore_service.service", "cache"),
			),
			testStepForEnv(
				"gke",
				testServiceDataSourceConfig,
				testGKEService("data.encore_service.service", "cache"),
			),
		},
	})
}

const testServiceDataSourceConfig = `
provider "encore" {
	auth_key = "test"
	env = "%s"
}

data "encore_service" "service" {
    name = "cache"
}
`
