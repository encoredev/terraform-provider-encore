package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type Satisfier struct {
	Type string `graphql:"__typename"`

	AWSSQSQueue           `graphql:"... on AWSSQSQueue"`
	GCPPubSubSubscription `graphql:"... on GCPPubSubSubscription"`

	GCPPubSubTopic `graphql:"... on GCPPubSubTopic"`
	AWSSNSTopic    `graphql:"... on AWSSNSTopic"`
}

var _ datasource.DataSource = &EncoreDataSource{}

func NewEncoreDataSource(typeRef TypeRef, name string, tfTypes ...TFType) datasource.DataSource {
	return &EncoreDataSource{
		typeRef: typeRef,
		name:    name,
		schema:  createSchema(tfTypes...),
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
