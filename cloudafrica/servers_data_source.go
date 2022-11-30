package cloudafrica

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &ServersDataSource{}
	_ datasource.DataSourceWithConfigure = &ServersDataSource{}
)

// ServersDataSourceModel maps the data source schema data.
type ServersDataSourceModel struct {
	Servers []ServerModel `tfsdk:"servers"`
}

func NewServersDataSource() datasource.DataSource {
	return &ServersDataSource{}
}

type ServersDataSource struct {
	client *TFClient
}

func (d *ServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servers"
}

func ServerSchema() map[string]tfsdk.Attribute {
	return map[string]tfsdk.Attribute(map[string]tfsdk.Attribute{
		"id": {
			Type:     types.Int64Type,
			Computed: true,
		},
		"name": {
			Type:     types.StringType,
			Required: true,
		},
		"state": {
			Type:     types.StringType,
			Required: true,
		},
		"cpus": {
			Type:     types.Int64Type,
			Required: true,
		},
		"ram_mib": {
			Type:     types.Int64Type,
			Required: true,
		},
		"ssh_keys": {
			Required: true,
			Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				"id": {
					Type:     types.Int64Type,
					Computed: true,
				},
				"name": {
					Type:     types.StringType,
					Required: true,
				},
				"body": {
					Type:     types.StringType,
					Required: true,
				},
			}),
		},
		"disks": {
			Required: true,
			Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				"id": {
					Type:     types.Int64Type,
					Computed: true,
				},
			}),
		},
	})
}

func (d *ServersDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"servers": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(
					ServerSchema(),
				),
			},
		},
	}, nil
}

func (d *ServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ServersDataSourceModel

	serversResp, _, err := d.client.Client.ServerApi.ListServers(d.client.Auth).Execute()
	if err != nil {
		tflog.Error(ctx, "Error Reading servers resource", map[string]any{"err": err.Error()})
		resp.Diagnostics.AddError(
			"Unable to Read CloudAfrica Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	servers := serversResp.Servers
	tflog.Trace(ctx, "Got servers", map[string]any{"servers": servers})
	for _, server := range servers {
		serverState := ServerModelFromApi(server)

		state.Servers = append(state.Servers, serverState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *ServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*TFClient)
}
