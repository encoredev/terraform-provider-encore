package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var queryType = reflect.TypeOf((*SatisfierQuery)(nil)).Elem()

type SatisfierQuery struct {
	Type string `graphql:"__typename"`

	AWSSNSSubscription    AWSSNSSubscription    `graphql:"... on AWSSNSSubscription" tf:"aws_sns"`
	GCPPubSubSubscription GCPPubSubSubscription `graphql:"... on GCPPubSubSubscription" tf:"gcp_pubsub"`

	AWSSNSTopic    AWSSNSTopic    `graphql:"... on AWSSNSTopic" tf:"aws_sns"`
	GCPPubSubTopic GCPPubSubTopic `graphql:"... on GCPPubSubTopic" tf:"gcp_pubsub"`

	SQLDatabase `graphql:"... on SQLDatabase"`

	RedisKeyspace `graphql:"... on RedisKeyspace"`

	Service `graphql:"... on Service"`

	Gateway `graphql:"... on Gateway"`
}

func (a *SatisfierQuery) GetDocs() (attrDesc map[string]string) {
	return map[string]string{
		"gcp_pubsub": "Set if the resource is provisioned by GCP Pub/Sub",
		"aws_sns":    "Set if the resource is provisioned AWS SNS",
	}
}

var _ datasource.DataSource = &EncoreDataSource{}

func NewEncoreDataSource(typeRef TypeRef, name, desc string, fragments ...string) datasource.DataSource {
	return &EncoreDataSource{
		typeRef: typeRef,
		name:    name,
		schema:  createSchema(desc, fragments...),
	}
}

type EncoreDataSource struct {
	needs   *NeedsData
	typeRef TypeRef
	name    string
	schema  schema.Schema
}

func (d *EncoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.name
}

func (d *EncoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.schema
}

func (d *EncoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	needs, ok := req.ProviderData.(*NeedsData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *NeedsData, received %T", req.ProviderData),
		)

		return
	}

	d.needs = needs
}

func (d *EncoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	resp.Diagnostics.Append(d.needs.SetValue(ctx, d.typeRef, req.Config, &resp.State)...)
}
