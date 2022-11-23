package cloudafrica

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// serversModel maps servers schema data.
type ServerModel struct {
	ID      types.Int64       `tfsdk:"id"`
	Name    types.String      `tfsdk:"name"`
	State   types.String      `tfsdk:"state"`
	CPUs    types.Int64       `tfsdk:"cpus"`
	RamMiB  types.Int64       `tfsdk:"ram_mib"`
	SSHKeys []SSHKeyModel     `tfsdk:"ssh_keys"`
	Disks   []ServerDiskModel `tfsdk:"disks"`
}
