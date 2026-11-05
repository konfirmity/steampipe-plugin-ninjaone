package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaWindowsService struct {
	DeviceId    *int64     `json:"deviceId"`
	Description flexString `json:"description"`
	DisplayName flexString `json:"displayName"`
	Name        flexString `json:"name"`
	StartType   flexString `json:"startType"`
	State       flexString `json:"state"`
	Timestamp   *float64   `json:"timestamp"`
	UserName    flexString `json:"userName"`
}

//// TABLE DEFINITION

func tableNinjaoneWindowsService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_windows_service",
		Description: "Windows services on each device.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaWindowsService]("ninjaone_windows_service", "/v2/queries/windows-services"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Service description."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Service display name."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Service name."},
			{Name: "start_type", Type: proto.ColumnType_STRING, Description: "Start type."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "Service state."},
			timeColumn("timestamp", "Time the device reported this data."),
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "Account the service runs as."},
			titleColumn("DisplayName"),
		},
	}
}
