// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure EncoreProvider satisfies various provider interfaces.
var _ provider.Provider = &EncoreProvider{}

// EncoreProvider defines the provider implementation.
type EncoreProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// EncoreProviderModel describes the provider data model.
type EncoreProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	AppID   types.String `tfsdk:"app_id"`
	EnvName types.String `tfsdk:"env_name"`
}

func (p *EncoreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "encore"
	resp.Version = p.version
}

func (p *EncoreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key to use to authenticate with the Encore Platform. If empty, the provider attempts to use the local Encore CLI for authentication.",
				Optional:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The default Encore application id to operate on, if not overridden on a resource.",
				Optional:            true,
			},
			"env_name": schema.StringAttribute{
				MarkdownDescription: "The default Encore environment to operate on, if not overridden on a resource.",
				Optional:            true,
			},
		},
	}
}

func (p *EncoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data EncoreProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *EncoreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *EncoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPubSubTopic,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &EncoreProvider{
			version: version,
		}
	}
}
