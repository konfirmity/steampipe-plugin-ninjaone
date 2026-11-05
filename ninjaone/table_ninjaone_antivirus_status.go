package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaAntivirusStatus struct {
	DeviceId         *int64     `json:"deviceId"`
	DefinitionStatus flexString `json:"definitionStatus"`
	ProductName      flexString `json:"productName"`
	ProductState     flexString `json:"productState"`
	Timestamp        *float64   `json:"timestamp"`
	Version          flexString `json:"version"`
}

//// TABLE DEFINITION

func tableNinjaoneAntivirusStatus(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_antivirus_status",
		Description: "Antivirus product and state on each device. A product name of NONE means no antivirus was detected.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaAntivirusStatus]("ninjaone_antivirus_status", "/v2/queries/antivirus-status"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "definition_status", Type: proto.ColumnType_STRING, Description: "Whether definitions are up to date."},
			{Name: "product_name", Type: proto.ColumnType_STRING, Description: "Antivirus product name, or NONE."},
			{Name: "product_state", Type: proto.ColumnType_STRING, Description: "Antivirus product state."},
			timeColumn("timestamp", "Time the device reported this data."),
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Antivirus version."},
			titleColumn("ProductName"),
		},
	}
}
