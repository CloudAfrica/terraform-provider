package cloudafrica

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServerDiskModel struct {
	ID types.Int64 `tfsdk:"id"`
}
