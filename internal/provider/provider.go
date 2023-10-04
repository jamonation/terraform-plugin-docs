package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &docsProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &docsProvider{
			version: version,
		}
	}
}

// docsProvider is the provider implementation.
type docsProvider struct {
	version string
}

type docsProviderModel struct{}

type ProviderOpts struct{}

// Metadata returns the provider type name.
func (p *docsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "docs"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *docsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *docsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var config docsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := &ProviderOpts{}

	resp.DataSourceData = opts
}

// Resources defines the resources implemented in the provider.
func (p *docsProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

// DataSources defines the data sources implemented in the provider.
func (p *docsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewReadmeDataSource,
	}
}
