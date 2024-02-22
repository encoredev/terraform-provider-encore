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
		"SNS Topic information",
		map[string]string{
			"arn": "The ARN for this  sns topic",
		}
}

type GCPPubSubTopic struct {
	SelfLink string
}

func (a *GCPPubSubTopic) GetDocs() (subkey string, mdDesc string, attrDesc map[string]string) {
	return "gcp_pubsub",
		"GCP Pubsub Topic information",
		map[string]string{
			"self_link": "The GCP self link to the topic",
		}
}
