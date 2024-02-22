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
		"SQS Queue information",
		map[string]string{
			"arn": "The ARN for this sqs queue",
		}
}

type GCPPubSubSubscription struct {
	SelfLink string
}

func (a *GCPPubSubSubscription) GetDocs() (subkey string, mdDesc string, attrDesc map[string]string) {
	return "gcp_pubsub",
		"SQS Queue information",
		map[string]string{
			"self_link": "The GCP self link to the topic",
		}
}
