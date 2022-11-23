package cloudafrica

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SSHKeyModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Body types.String `tfsdk:"body"`
}
