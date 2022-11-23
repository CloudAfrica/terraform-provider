package cloudafrica

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ImageModel struct {
	ID      types.Int64  `tfsdk:"id"`
	OS      types.String `tfsdk:"os"`
	Version types.String `tfsdk:"version"`
}
