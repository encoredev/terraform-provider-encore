package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAWSSNS() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.encore_pubsub_topic.topic", "aws_sns.arn", "arn:aws:sns:region:account:app-env-events"),
	)
}

func testGCPTopic() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.encore_pubsub_topic.topic", "gcp_pubsub.id", "projects/app-env/topics/events"),
	)
}

func TestPubsubTopicDataSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testV6ProviderFactories,
		Steps: []resource.TestStep{
			testStepForEnv(
				"eks",
				testPubsubTopicDataSourceConfig,
				testAWSSNS(),
			),
			testStepForEnv(
				"fargate",
				testPubsubTopicDataSourceConfig,
				testAWSSNS(),
			),
			testStepForEnv(
				"cloudrun",
				testPubsubTopicDataSourceConfig,
				testGCPTopic(),
			),
			testStepForEnv(
				"gke",
				testPubsubTopicDataSourceConfig,
				testGCPTopic(),
			),
		},
	})
}

const testPubsubTopicDataSourceConfig = `
provider "encore" {
	auth_key = "test"
	env = "%s"
}

data "encore_pubsub_topic" "topic" {
    name = "events"
}

`
