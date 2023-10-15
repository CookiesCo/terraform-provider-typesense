package typesense

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	keyEnvName = "TYPESENSE_MANAGEMENT_KEY"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &typesenseProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &typesenseProvider{}
}

// typesenseProvider is the provider implementation.
type typesenseProvider struct{}

// Metadata returns the provider type name.
func (p *typesenseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "typesense"
}

// Schema defines the provider-level schema for configuration data.
func (p *typesenseProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage your Typesense clusters",
		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				Description: "Cloud Management API Key",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a typesense client for data sources and resources.
func (p *typesenseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Typesense client")
	// Retrieve provider data from configuration
	var config typesenseProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Key.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Unknown Typesense API Key",
			"The provider cannot create the Typesense API client as there is an unknown configuration value for the Typesense API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the "+keyEnvName+" environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	key := os.Getenv(keyEnvName)
	if !config.Key.IsNull() {
		key = config.Key.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Missing Typesense API Key",
			"The provider cannot create the Typesense API client as there is a missing or empty value for the Typesense API key. "+
				"Set the key value in the configuration or use the "+keyEnvName+" environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "foo", "bar")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "foo")

	// Create a new Typesense client using the configuration values
	client, err := NewClient(key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Typesense API Client",
			"An unexpected error occurred when creating the Typesense API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Typesense Client Error: "+err.Error(),
		)
		return
	}

	// Make the Typesense client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Configured Typesense client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *typesenseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClusterDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *typesenseProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
		NewClusterApiKeysResource,
	}
}

// typesenseProviderModel maps provider schema data to a Go type.
type typesenseProviderModel struct {
	Key types.String `tfsdk:"key"`
}

// typesenseClusterModel maps Typesense cluster schema data.
type typesenseClusterModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Memory                 types.String `tfsdk:"memory"`
	VCPU                   types.String `tfsdk:"vcpu"`
	HighPerformanceDisk    types.String `tfsdk:"high_performance_disk"`
	TypesenseServerVersion types.String `tfsdk:"typesense_server_version"`
	HighAvailability       types.String `tfsdk:"high_availability"`
	SearchDeliveryNetwork  types.String `tfsdk:"search_delivery_network"`
	LoadBalancing          types.String `tfsdk:"load_balancing"`
	Region                 types.String `tfsdk:"region"`
	AutoUpgradeCapacity    types.Bool   `tfsdk:"auto_upgrade_capacity"`
	Status                 types.String `tfsdk:"status"`
}

type typesenseClusterApiKeysModel struct {
	ID            types.String `tfsdk:"id"`
	ClusterId     types.String `tfsdk:"cluster_id"`
	AdminKey      types.String `tfsdk:"admin_key"`
	SearchOnlyKey types.String `tfsdk:"search_only_key"`
}
