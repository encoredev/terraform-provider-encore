package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAWSSQS() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.encore_pubsub_subscription.subscription", "aws_sns.arn", "arn:aws:sns:region:account:app-env-events"),
		resource.TestCheckResourceAttr("data.encore_pubsub_subscription.subscription", "aws_sns.topic.arn", "arn:aws:sns:region:account:app-env-events"),
		resource.TestCheckResourceAttr("data.encore_pubsub_subscription.subscription", "aws_sns.queue.arn", "arn:aws:sqs:region:account:app-env-events-log-event"),
		resource.TestCheckResourceAttr("data.encore_pubsub_subscription.subscription", "aws_sns.queue.dead_letter.arn", "arn:aws:sqs:region:account:app-env-events-log-event-dlq"),
	)
}

func testGCPSub() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.encore_pubsub_subscription.subscription", "gcp_pubsub.id", "projects/app-env/subscriptions/events.log-event"),
		resource.TestCheckResourceAttr("data.encore_pubsub_subscription.subscription", "gcp_pubsub.topic.id", "projects/app-env/topics/events"),
		resource.TestCheckResourceAttr("data.encore_pubsub_subscription.subscription", "gcp_pubsub.dead_letter.id", "projects/app-env/subscriptions/events.log-event.deadletter.encore"),
		resource.TestCheckResourceAttr("data.encore_pubsub_subscription.subscription", "gcp_pubsub.dead_letter.topic.id", "projects/app-env/topics/events.log-event.deadletter"),
	)
}

func TestPubsubSubDataSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testV6ProviderFactories,
		Steps: []resource.TestStep{
			testStepForEnv(
				"eks",
				testPubsubSubDataSourceConfig,
				testAWSSQS(),
			),
			testStepForEnv(
				"fargate",
				testPubsubSubDataSourceConfig,
				testAWSSQS(),
			),
			testStepForEnv(
				"cloudrun",
				testPubsubSubDataSourceConfig,
				testGCPSub(),
			),
			testStepForEnv(
				"gke",
				testPubsubSubDataSourceConfig,
				testGCPSub(),
			),
		},
	})
}

const testPubsubSubDataSourceConfig = `
provider "encore" {
	auth_key = "test"
	env = "%s"
}

data "encore_pubsub_subscription" "subscription" {
    name = "log-event"
}

`
