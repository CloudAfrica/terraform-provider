package cloudafrica

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// serversModel maps servers schema data.
type BillingAccountModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
