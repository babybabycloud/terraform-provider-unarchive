package unarchive

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &unarchiveProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &unarchiveProvider{}
}

// unarchiveProvider is the provider implementation.
type unarchiveProvider struct{}

// Metadata returns the provider type name.
func (p *unarchiveProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "unarchive"
}

// Schema defines the provider-level schema for configuration data.
func (p *unarchiveProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *unarchiveProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *unarchiveProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFileDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *unarchiveProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
