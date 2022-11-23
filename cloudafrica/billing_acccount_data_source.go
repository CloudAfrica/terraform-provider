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
	_ datasource.DataSource              = &BillingAccountsDataSource{}
	_ datasource.DataSourceWithConfigure = &BillingAccountsDataSource{}
)

// BillingAccountsDataSourceModel maps the data source schema data.
type BillingAccountsDataSourceModel struct {
	BillingAccounts []BillingAccountModel `tfsdk:"billing_accounts"`
}

func NewBillingAccountsDataSource() datasource.DataSource {
	return &BillingAccountsDataSource{}
}

type BillingAccountsDataSource struct {
	client *TFClient
}

func (d *BillingAccountsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_accounts"
}

func (d *BillingAccountsDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"billing_accounts": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.Int64Type,
						Computed: true,
					},
					"name": {
						Type:     types.StringType,
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

func (d *BillingAccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BillingAccountsDataSourceModel

	items_resp, _, err := d.client.Client.BillingAccountApi.ListBillingAccounts(d.client.Auth).Execute()
	if err != nil {
		tflog.Error(ctx, "Error Reading billing-accounts resource", map[string]any{"err": err.Error()})
		resp.Diagnostics.AddError(
			"Unable to Read CloudAfrica BillingAccounts",
			err.Error(),
		)
		return
	}

	// Map response body to model
	items := items_resp.BillingAccounts
	tflog.Trace(ctx, "Got billing-accounts", map[string]any{"items": items})
	for _, item := range items {
		itemstate := BillingAccountModel{
			ID:   types.Int64Value(int64(item.Id)),
			Name: types.StringValue(item.Name),
		}

		state.BillingAccounts = append(state.BillingAccounts, itemstate)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *BillingAccountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*TFClient)
}
