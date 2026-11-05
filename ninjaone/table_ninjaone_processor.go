package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaProcessor struct {
	DeviceId        *int64     `json:"deviceId"`
	Architecture    flexString `json:"architecture"`
	ClockSpeed      *int64     `json:"clockSpeed"`
	MaxClockSpeed   *int64     `json:"maxClockSpeed"`
	Name            flexString `json:"name"`
	NumCores        *int64     `json:"numCores"`
	NumLogicalCores *int64     `json:"numLogicalCores"`
	Timestamp       *float64   `json:"timestamp"`
}

//// TABLE DEFINITION

func tableNinjaoneProcessor(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_processor",
		Description: "Processors on each device.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaProcessor]("ninjaone_processor", "/v2/queries/processors"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "architecture", Type: proto.ColumnType_STRING, Description: "Processor architecture."},
			{Name: "clock_speed", Type: proto.ColumnType_INT, Description: "Current clock speed."},
			{Name: "max_clock_speed", Type: proto.ColumnType_INT, Description: "Maximum clock speed."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Processor name."},
			{Name: "num_cores", Type: proto.ColumnType_INT, Description: "Physical cores."},
			{Name: "num_logical_cores", Type: proto.ColumnType_INT, Description: "Logical cores."},
			timeColumn("timestamp", "Time the device reported this data."),
			titleColumn("Name"),
		},
	}
}
