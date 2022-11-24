package cloudafrica

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"strconv"
)

var (
	_ resource.Resource              = &ServerResource{}
	_ resource.ResourceWithConfigure = &ServerResource{}
)

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

type ServerResource struct {
	client *TFClient
}

func (d *ServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *ServerResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: ServerSchema(),
	}, nil
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to import  CloudAfrica Server. Could not convert import ID to int64",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServerModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serversResp, _, err := d.client.Client.ServerApi.GetServer(d.client.Auth, state.ID.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read CloudAfrica Server",
			err.Error(),
		)
		return
	}

	// Map response body to model
	apiServer := serversResp.Server
	state = ServerModelFromApi(apiServer)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *ServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*TFClient)
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	//	var plan serversModel
	//	diags := req.Plan.Get(ctx, &plan)
	//	resp.Diagnostics.Append(diags...)
	//	if resp.Diagnostics.HasError() {
	//		return
	//	}
	//
	//	var networkID int64 = 17592186045489
	//	var sshKeyID int64 = 17592186045481
	//
	//	// Generate API request body from plan
	//	var order = cloudafrica.ServerOrder{
	//		17592186045479,
	//		17592186045485,
	//		"andre-test",
	//		1,
	//		1024,
	//		[]cloudafrica.ServerOrderDisk{cloudafrica.ServerOrderDisk{0, 10000}},
	//		[]cloudafrica.ServerOrderNetworkInterface{cloudafrica.ServerOrderNetworkInterface{0, true, networkID}},
	//		[]int64{sshKeyID}}
	//	//for _, item := range plan.Items {
	//	//	items = append(items, hashicups.OrderItem{
	//	//		Coffee: hashicups.Coffee{
	//	//			ID: int(item.Coffee.ID.ValueInt64()),
	//	//		},
	//	//		Quantity: int(item.Quantity.ValueInt64()),
	//	//	})
	//	//}
	//
	//	// Create new order
	//	orderResp, err := r.client.Servers.Create(order)
	//	if err != nil {
	//		resp.Diagnostics.AddError(
	//			"Error creating order",
	//			"Could not create order, unexpected error: "+err.Error(),
	//		)
	//		return
	//	}
	//
	//	// Map response body to schema and populate Computed attribute values
	//	plan.ID = types.Int64Value(orderResp.ServerID)
	//	//for orderItemIndex, orderItem := range order.Items {
	//	//	plan.Items[orderItemIndex] = orderItemModel{
	//	//		Coffee: orderItemCoffeeModel{
	//	//			ID:          types.Int64Value(int64(orderItem.Coffee.ID)),
	//	//			Name:        types.StringValue(orderItem.Coffee.Name),
	//	//			Teaser:      types.StringValue(orderItem.Coffee.Teaser),
	//	//			Description: types.StringValue(orderItem.Coffee.Description),
	//	//			Price:       types.Float64Value(orderItem.Coffee.Price),
	//	//			Image:       types.StringValue(orderItem.Coffee.Image),
	//	//		},
	//	//		Quantity: types.Int64Value(int64(orderItem.Quantity)),
	//	//	}
	//	//}
	//	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	//
	//	// Set state to fully populated data
	//	diags = resp.State.Set(ctx, plan)
	//	resp.Diagnostics.Append(diags...)
	//	if resp.Diagnostics.HasError() {
	//		return
	//	}
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
