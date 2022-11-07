package cloudafrica

import (
	"context"

	"github.com/CloudAfrica/goclient/cloudafrica"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &serversDataSource{}
	_ datasource.DataSourceWithConfigure = &serversDataSource{}
)

// serversDataSourceModel maps the data source schema data.
type serversDataSourceModel struct {
	Servers []serversModel `tfsdk:"servers"`
}

// serversModel maps servers schema data.
type serversModel struct {
	ID types.Int64 `tfsdk:"id"`
}

func NewServersDataSource() datasource.DataSource {
	return &serversDataSource{}
}

type serversDataSource struct {
	client *cloudafrica.CloudAfrica
}

func (d *serversDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servers"
}

func (d *serversDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"servers": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.Int64Type,
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

func (d *serversDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serversDataSourceModel

	servers, err := d.client.Servers.List()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read CloudAfrica Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, server := range servers {
		serverstate := serversModel{
			ID: types.Int64Value(int64(server.ID)),
		}

		// TODO: disks and network
		//for _, ingredient := range coffee.Ingredient {
		//	serverstate.Ingredients = append(serverstate.Ingredients, serversIngredientsModel{
		//		ID: types.Int64Value(int64(ingredient.ID)),
		//	})
		//}

		state.Servers = append(state.Servers, serverstate)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serversDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*cloudafrica.CloudAfrica)
}
