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
	_ datasource.DataSource              = &ImageDataSource{}
	_ datasource.DataSourceWithConfigure = &ImageDataSource{}
)

// ImageDataSourceModel maps the data source schema data.
type ImageDataSourceModel struct {
	Images []ImageModel `tfsdk:"images"`
}

func NewImageDataSource() datasource.DataSource {
	return &ImageDataSource{}
}

type ImageDataSource struct {
	client *TFClient
}

func (d *ImageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

func (d *ImageDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"images": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.Int64Type,
						Computed: true,
					},
					"os": {
						Type:     types.StringType,
						Computed: true,
					},
					"version": {
						Type:     types.StringType,
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

func (d *ImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ImageDataSourceModel

	images_resp, _, err := d.client.Client.DefaultApi.ListImages(d.client.Auth).Execute()
	if err != nil {
		tflog.Error(ctx, "Error Reading images resource", map[string]any{"err": err.Error()})
		resp.Diagnostics.AddError(
			"Unable to Read CloudAfrica Images",
			err.Error(),
		)
		return
	}

	// Map response body to model
	images := images_resp.Images
	tflog.Trace(ctx, "Got images", map[string]any{"images": images})
	for _, image := range images {
		imagestate := ImageModel{
			ID:      types.Int64Value(int64(image.Id)),
			OS:      types.StringValue(image.Os),
			Version: types.StringValue(image.Version),
		}

		state.Images = append(state.Images, imagestate)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *ImageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*TFClient)
}
