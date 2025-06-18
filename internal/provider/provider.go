package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/360-build/terraform-provider-vaultassist/internal/provider/vaultclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &vaultassistProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &vaultassistProvider{
			version: version,
		}
	}
}

// vaultassistProvider is the provider implementation.
type vaultassistProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *vaultassistProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vaultassist"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *vaultassistProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Required: true,
			},
			"role": schema.StringAttribute{
				Required: true,
			},
			"mountpoint": schema.StringAttribute{
				Required: true,
			},
			"headers": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a vaultassist API client for data sources and resources.
func (p *vaultassistProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve the configuration values.
	var config struct {
		Address    string            `tfsdk:"address"`
		Role       string            `tfsdk:"role"`
		Mountpoint string            `tfsdk:"mountpoint"`
		Headers    map[string]string `tfsdk:"headers"`
	}
	req.Config.Get(ctx, &config)

	//Initialize the Vault client.
	client, err := vaultclient.NewVaultClient(config.Address, config.Mountpoint, config.Role, config.Headers)
	if err != nil {
		resp.Diagnostics.AddError("failed to initialize Vault client", err.Error())
		return
	}

	//Store the client in the context.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *vaultassistProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *vaultassistProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBootstrapSecretResource,
		NewPatchSecretResource,
	}
}
