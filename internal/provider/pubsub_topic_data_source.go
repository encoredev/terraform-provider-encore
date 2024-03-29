package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewPubSubTopic() datasource.DataSource {
	return NewEncoreDataSource(
		"need.Topic",
		"pubsub_topic",
		"Encore provisioned Pub/Sub topic information",
		"GCPPubSubTopic",
		"AWSSNSTopic",
	)
}

type AWSSNSTopic struct {
	Arn string
}

func (a *AWSSNSTopic) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"arn": "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for this resource",
	}
}

type GCPPubSubTopic struct {
	SelfLink string `tf:"id"`
}

func (a *GCPPubSubTopic) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"id": "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) in the form of `projects/{project}/topics/{topic}`",
	}
}
