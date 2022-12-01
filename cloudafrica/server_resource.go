package cloudafrica

import (
	"context"
	cloudafrica "github.com/CloudAfrica/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"time"
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
	var plan ServerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var networkID int64 = 17592186045491
	var sshKeyID int64 = 17592186045481
	var billingAccountId int64 = 17592186045487
	var ImageID int64 = 17592186045479

	// Generate API request body from plan
	var order = cloudafrica.ServerOrder{
		ImageId:           ImageID,
		Name:              plan.Name.ValueString(),
		Cpus:              1,
		RamMib:            1024,
		BillingAccountId:  billingAccountId,
		Disks:             []cloudafrica.ServerOrderDisk{cloudafrica.ServerOrderDisk{SizeMb: 10000}},
		NetworkInterfaces: []cloudafrica.ServerOrderNetworkInterface{cloudafrica.ServerOrderNetworkInterface{Primary: true, NetworkId: networkID}},
		SshKeyIds:         []int64{sshKeyID}}
	//for _, item := range plan.Items {
	//	items = append(items, hashicups.OrderItem{
	//		Coffee: hashicups.Coffee{
	//			ID: int(item.Coffee.ID.ValueInt64()),
	//		},
	//		Quantity: int(item.Quantity.ValueInt64()),
	//	})
	//}

	// Create new order
	orderResp, _, err := r.client.Client.ServerApi.CreateServer(r.client.Auth).ServerOrder(order).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(orderResp.ServerId)

	err = orderResp.Wait(
		r.client.Client,
		r.client.Auth,
		20,
		1*time.Second,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for order task to finish",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	server, _, err := r.client.Client.ServerApi.GetServer(r.client.Auth, orderResp.ServerId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching server after order created",
			"Could not import server in TF, unexpected error: "+err.Error(),
		)
		return
	}

	plan = ServerModelFromApi(server.Server)

	tflog.Info(ctx, "Server:", map[string]any{"server": plan})

	//for orderItemIndex, orderItem := range order.Items {
	//	plan.Items[orderItemIndex] = orderItemModel{
	//		Coffee: orderItemCoffeeModel{
	//			ID:          types.Int64Value(int64(orderItem.Coffee.ID)),
	//			Name:        types.StringValue(orderItem.Coffee.Name),
	//			Teaser:      types.StringValue(orderItem.Coffee.Teaser),
	//			Description: types.StringValue(orderItem.Coffee.Description),
	//			Price:       types.Float64Value(orderItem.Coffee.Price),
	//			Image:       types.StringValue(orderItem.Coffee.Image),
	//		},
	//		Quantity: types.Int64Value(int64(orderItem.Quantity)),
	//	}
	//}
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
