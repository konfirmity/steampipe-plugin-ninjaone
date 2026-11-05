package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaSoftware struct {
	DeviceId          *int64     `json:"deviceId"`
	InstallDate       flexString `json:"installDate"`
	IsSystemComponent *bool      `json:"isSystemComponent"`
	Location          flexString `json:"location"`
	Name              flexString `json:"name"`
	ProductCode       flexString `json:"productCode"`
	Publisher         flexString `json:"publisher"`
	Size              *int64     `json:"size"`
	Timestamp         *float64   `json:"timestamp"`
	Version           flexString `json:"version"`
}

//// TABLE DEFINITION

func tableNinjaoneSoftware(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_software",
		Description: "Software installed on each device.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaSoftware]("ninjaone_software", "/v2/queries/software"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "install_date", Type: proto.ColumnType_STRING, Description: "Install date as reported by the device."},
			{Name: "is_system_component", Type: proto.ColumnType_BOOL, Description: "True if the item is an operating system component."},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "Install location."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Software name."},
			{Name: "product_code", Type: proto.ColumnType_STRING, Description: "Vendor product code."},
			{Name: "publisher", Type: proto.ColumnType_STRING, Description: "Software publisher."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "Install size in bytes."},
			timeColumn("timestamp", "Time the device reported this data."),
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Software version."},
			titleColumn("Name"),
		},
	}
}
