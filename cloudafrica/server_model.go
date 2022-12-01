package cloudafrica

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CloudAfrica/client"
)

// ServerModel maps servers schema data.
type ServerModel struct {
	ID      types.Int64       `tfsdk:"id"`
	Name    types.String      `tfsdk:"name"`
	State   types.String      `tfsdk:"state"`
	CPUs    types.Int64       `tfsdk:"cpus"`
	RamMiB  types.Int64       `tfsdk:"ram_mib"`
	SSHKeys []SSHKeyModel     `tfsdk:"ssh_keys"`
	Disks   []ServerDiskModel `tfsdk:"disks"`
}

func ServerModelFromApi(server cloudafrica.Server) ServerModel {
	s := ServerModel{
		ID:     types.Int64Value(int64(*server.Id)),
		Name:   types.StringValue(*server.Name),
		State:  types.StringValue(*server.State),
		CPUs:   types.Int64Value(int64(*server.Cpus)),
		RamMiB: types.Int64Value(int64(*server.RamMib)),
	}

	for _, disk := range server.Disks {
		s.Disks = append(s.Disks, ServerDiskModel{
			ID: types.Int64Value(*disk.Id),
		})
	}

	for _, key := range server.SshKeys {
		s.SSHKeys = append(s.SSHKeys, SSHKeyModel{
			ID:   types.Int64Value(*key.Id),
			Name: types.StringValue(*key.Name),
			Body: types.StringValue(*key.Body),
		})
	}

	return s
}
