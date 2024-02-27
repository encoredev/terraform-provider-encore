package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewPubSubSubscription() datasource.DataSource {
	return NewEncoreDataSource(
		"need.Subscription",
		"pubsub_subscription",
		"Encore provisioned Pub/Sub Subscription information",
		"AWSSNSSubscription",
		"GCPPubSubSubscription",
	)
}

type AWSSNSSubscription struct {
	Arn                string
	WrappedAWSSNSTopic `graphql:"topic"`
	Queue              AWSSQSQueue
}

func (a *AWSSNSSubscription) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"arn":   "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for this resource",
		"queue": "The sqs queue which this subscription forwards messages to",
	}
}

type AWSSQSQueue struct {
	Arn        string
	DeadLetter AWSDeadLetterQueue `graphql:"dlq"`
}

func (a *AWSSQSQueue) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"arn":         "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for this resource",
		"dead_letter": "The dead letter queue for this subscription",
	}
}

type WrappedAWSSNSTopic struct {
	Topic AWSSNSTopic `graphql:"... on AWSSNSTopic"`
}

func (a *WrappedAWSSNSTopic) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"topic": "The topic which this subscription is subscribed to",
	}
}

type AWSDeadLetterQueue struct {
	Arn string
}

func (a *AWSDeadLetterQueue) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"arn": "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for this resource",
	}
}

type GCPPubSubSubscription struct {
	SelfLink              string `tf:"id"`
	WrappedGCPPubSubTopic `graphql:"topic"`
	DeadLetter            GCPDeadLetterQueue `graphql:"dlq"`
}

func (a *GCPPubSubSubscription) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"id":          "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) in the form of `projects/{project}/subscriptions/{subscription}`",
		"dead_letter": "The dead letter queue for this subscription",
	}
}

type GCPDeadLetterQueue struct {
	SelfLink              string `tf:"id"`
	WrappedGCPPubSubTopic `graphql:"topic"`
}

func (a *GCPDeadLetterQueue) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"id": "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) in the form of `projects/{project}/subscriptions/{subscription}`",
	}
}

type WrappedGCPPubSubTopic struct {
	Topic GCPPubSubTopic `graphql:"... on GCPPubSubTopic"`
}
