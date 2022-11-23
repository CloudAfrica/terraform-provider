package cloudafrica

import (
	"context"
	"os"

	"github.com/CloudAfrica/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TFClient struct {
	Client *cloudafrica.APIClient
	Auth   context.Context
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &cloudafricaProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &cloudafricaProvider{}
}

// cloudafricaProvider is the provider implementation.
type cloudafricaProvider struct{}

type cloudafricaProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *cloudafricaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudafrica"
}

// GetSchema defines the provider-level schema for configuration data.
func (p *cloudafricaProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:     types.StringType,
				Optional: true,
			},
			"token": {
				Description: "API Key to access the CloudAfrica API. May also be provided via CLOUDAFRICA_TOKEN environment variable.",
				Type:        types.StringType,
				Optional:    true,
			},
		},
	}, nil
}

// Configure prepares a Cloudafrica API client for data sources and resources.
func (p *cloudafricaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config cloudafricaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown CloudAfrica API Host",
			"The provider cannot create the CloudAfrica API client as there is an unknown configuration value for the CloudAfrica API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CLOUDAFRICA_HOST environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown CloudAfrica API Token",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the CloudAfrica API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CloudAfrica_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	//	host := os.Getenv("CLOUDAFRICA_HOST")

	//if !config.Host.IsNull() {
	//	host = config.Host.ValueString()
	//}

	token := os.Getenv("CLOUDAFRICA_TOKEN")

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	//
	tflog.Info(ctx, "Creating the client", map[string]any{"token": token})

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing CloudAfrica API Token",
			"The provider cannot create the CloudAfrica API client as there is a missing or empty value for the CloudAfrica API token. "+
				"Set the username value in the configuration or use the CLOUDAFRICA_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new CloudAfrica client using the configuration values

	auth := context.WithValue(context.Background(), cloudafrica.ContextAccessToken, token)
	cfg := cloudafrica.NewConfiguration()
	client := cloudafrica.NewAPIClient(cfg)

	tf_client := TFClient{Client: client, Auth: auth}

	//client, err := cloudafrica.NewClient(&host, &token)
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Unable to Create CloudAfrica API Client",
	//		"An unexpected error occurred when creating the CloudAfrica API client. "+
	//			"If the error is not clear, please contact support@cloudafrica.net.\n\n"+
	//			"CloudAfrica Client Error: "+err.Error(),
	//	)
	//	return
	//}

	// Make the CloudAfrica client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &tf_client
	resp.ResourceData = &tf_client
}

// DataSources defines the data sources implemented in the provider.
func (p *cloudafricaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServersDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *cloudafricaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServersResource,
	}
}
