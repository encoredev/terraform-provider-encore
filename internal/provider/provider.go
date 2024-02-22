package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	EnvName types.String `tfsdk:"env"`
}

func (p *EncoreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "encore"
	resp.Version = p.version
}

func (p *EncoreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"env": schema.StringAttribute{
				MarkdownDescription: "The default Encore environment to operate on, if not overridden on a resource.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key to use to authenticate with the Encore Platform. If empty, the provider attempts to use ENCORE_API_KEY env.",
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

	client := NewPlatformClient(p.version)

	apiKey := data.APIKey.ValueString()
	if apiKey == "" {
		apiKey = os.Getenv("ENCORE_API_KEY")
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(path.Root("api_key"), "missing key", "missing encore api key")
		return
	}

	err := client.Auth(ctx, apiKey)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("api_key"), "invalid key", "encore platform auth failed")
		return
	}
	needs := NewNeedsData(client, data.EnvName.ValueString(), p.DataSources(ctx))
	resp.DataSourceData = needs
	resp.ResourceData = needs
}

func (p *EncoreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *EncoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPubSubTopic,
		NewPubSubSubscription,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &EncoreProvider{
			version: version,
		}
	}
}
