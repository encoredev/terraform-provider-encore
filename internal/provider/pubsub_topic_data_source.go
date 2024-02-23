package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewPubSubTopic() datasource.DataSource {
	return NewEncoreDataSource(
		"need.Topic",
		"pubsub_topic",
		&GCPPubSubTopic{},
		&AWSSNSTopic{})
}

type AWSSNSTopic struct {
	Arn string
}

func (a *AWSSNSTopic) GetDocs() (subkey string, mdDesc string, attrDesc map[string]string) {
	return "aws_sns",
		"Encore provisioned SNS Topic information",
		map[string]string{
			"arn": "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for this  sns topic",
		}
}

type GCPPubSubTopic struct {
	SelfLink string `tfsdk:"id"`
}

func (a *GCPPubSubTopic) GetDocs() (subkey string, mdDesc string, attrDesc map[string]string) {
	return "gcp_pubsub",
		"Encore provisioned GCP Pubsub Topic information",
		map[string]string{
			"id": "The [relative resource name](https://cloud.google.com/apis/design/resource_names#id) for this Pubsub topic in the form of `projects/{project}/topics/{topic}`",
		}
}
