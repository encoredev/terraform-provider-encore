package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewPubSubSubscription() datasource.DataSource {
	return NewEncoreDataSource(
		"need.Subscription",
		"pubsub_subscription",
		&GCPPubSubSubscription{},
		&AWSSQSQueue{})
}

type AWSSQSQueue struct {
	Arn string
}

func (a *AWSSQSQueue) GetDocs() (subkey string, mdDesc string, attrDesc map[string]string) {
	return "aws_sqs",
		"Encore provisioned SQS Queue information",
		map[string]string{
			"arn": "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for this sqs queue",
		}
}

type GCPPubSubSubscription struct {
	SelfLink string `tfsdk:"id"`
}

func (a *GCPPubSubSubscription) GetDocs() (subkey string, mdDesc string, attrDesc map[string]string) {
	return "gcp_pubsub",
		"Encore provisioned PubSub Subscription information",
		map[string]string{
			"id": "The [relative resource name](https://cloud.google.com/apis/design/resource_names#id) for this PubSub subscription in the form of `projects/{project}/subscriptions/{subscription}`",
		}
}
