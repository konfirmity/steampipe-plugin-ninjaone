package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaOsPatchInstall struct {
	DeviceId    *int64     `json:"deviceId"`
	Id          flexString `json:"id"`
	InstalledAt *float64   `json:"installedAt"`
	KbNumber    flexString `json:"kbNumber"`
	Name        flexString `json:"name"`
	Severity    flexString `json:"severity"`
	Status      flexString `json:"status"`
	Timestamp   *float64   `json:"timestamp"`
	Type        flexString `json:"type"`
}

//// TABLE DEFINITION

func tableNinjaoneOsPatchInstall(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_os_patch_install",
		Description: "Operating system patch installation attempts, successful or failed.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaOsPatchInstall]("ninjaone_os_patch_install", "/v2/queries/os-patch-installs"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Patch identifier."},
			timeColumn("installed_at", "Time the patch was installed."),
			{Name: "kb_number", Type: proto.ColumnType_STRING, Description: "Microsoft knowledge base number."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Patch name."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "Patch severity."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Install outcome."},
			timeColumn("timestamp", "Time the device reported this data."),
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Patch type."},
			titleColumn("Name"),
		},
	}
}
