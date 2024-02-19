// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PubSubTopic{}

func NewPubSubTopic() datasource.DataSource {
	return &PubSubTopic{}
}

// PubSubTopic defines the data source implementation for Encore Pub/Sub topics.
type PubSubTopic struct {
	client *http.Client
}

// PubSubTopicModel describes the data source data model.
type PubSubTopicModel struct {
	// Cloud provider-related information for this resource.
	Cloud types.Object `tfsdk:"cloud"`

	// AWS-specific information for this resource. Only available if the resource is provisioned in AWS.
	AWS types.Object `tfsdk:"aws"`
}

func (d *PubSubTopic) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pubsub_topic"
}

func (d *PubSubTopic) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source that provides information about an Encore-managed Pub/Sub topic.",

		Attributes: map[string]schema.Attribute{
			"aws": schema.SingleNestedAttribute{
				MarkdownDescription: "AWS-specific information for this resource. Only available if the resource is provisioned in AWS.",
				Computed:            true,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"arn": schema.StringAttribute{
						MarkdownDescription: "The ARN for this resource.",
						Computed:            true,
					},
				},
			},
			"cloud": schema.SingleNestedAttribute{
				MarkdownDescription: "Cloud provider-related information for this resource.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"guid": schema.StringAttribute{
						MarkdownDescription: "A globally unique id for this resource. In AWS this is an ARN, and in GCP it generally a self-link.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *PubSubTopic) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *PubSubTopic) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PubSubTopicModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	var diags diag.Diagnostics
	data.Cloud = types.ObjectNull(map[string]attr.Type{
		"guid": types.StringType,
	})

	data.AWS, diags = types.ObjectValue(map[string]attr.Type{
		"arn": types.StringType,
	}, map[string]attr.Value{
		"arn": types.StringValue("arn:aws:sns:us-west-2:123456789012:MyTopic"),
	})
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
