package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaSoftwarePatch struct {
	DeviceId          *int64     `json:"deviceId"`
	Id                flexString `json:"id"`
	Impact            flexString `json:"impact"`
	InstalledAt       *float64   `json:"installedAt"`
	PatchTitle        flexString `json:"title"`
	ProductIdentifier flexString `json:"productIdentifier"`
	Status            flexString `json:"status"`
	Timestamp         *float64   `json:"timestamp"`
	Type              flexString `json:"type"`
}

//// TABLE DEFINITION

func tableNinjaoneSoftwarePatch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_software_patch",
		Description: "Third-party software patches that are pending, failed or rejected. An empty result does not prove a device is patched: check that its policy enables software patching.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaSoftwarePatch]("ninjaone_software_patch", "/v2/queries/software-patches"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Patch identifier."},
			{Name: "impact", Type: proto.ColumnType_STRING, Description: "Patch impact."},
			timeColumn("installed_at", "Time the patch was installed."),
			{Name: "patch_title", Type: proto.ColumnType_STRING, Description: "Patch title."},
			{Name: "product_identifier", Type: proto.ColumnType_STRING, Description: "Identifier of the patched product."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Patch status."},
			timeColumn("timestamp", "Time the device reported this data."),
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Patch type."},
			titleColumn("PatchTitle"),
		},
	}
}
