package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaDisk struct {
	DeviceId       *int64     `json:"deviceId"`
	BytesPerSector *int64     `json:"bytesPerSector"`
	Description    flexString `json:"description"`
	InterfaceType  flexString `json:"interfaceType"`
	Manufacturer   flexString `json:"manufacturer"`
	MediaType      flexString `json:"mediaType"`
	Model          flexString `json:"model"`
	Name           flexString `json:"name"`
	PartitionCount *int64     `json:"partitionCount"`
	SerialNumber   flexString `json:"serialNumber"`
	Size           *int64     `json:"size"`
	SmartCapable   *bool      `json:"smartCapable"`
	Status         flexString `json:"status"`
	Timestamp      *float64   `json:"timestamp"`
}

//// TABLE DEFINITION

func tableNinjaoneDisk(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_disk",
		Description: "Physical disks on each device.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaDisk]("ninjaone_disk", "/v2/queries/disks"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "bytes_per_sector", Type: proto.ColumnType_INT, Description: "Bytes per sector."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Disk description."},
			{Name: "interface_type", Type: proto.ColumnType_STRING, Description: "Disk interface."},
			{Name: "manufacturer", Type: proto.ColumnType_STRING, Description: "Disk manufacturer."},
			{Name: "media_type", Type: proto.ColumnType_STRING, Description: "Media type."},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "Disk model."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Disk name."},
			{Name: "partition_count", Type: proto.ColumnType_INT, Description: "Number of partitions."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Disk serial number."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "Size in bytes."},
			{Name: "smart_capable", Type: proto.ColumnType_BOOL, Description: "True if the disk supports SMART."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Disk status."},
			timeColumn("timestamp", "Time the device reported this data."),
			titleColumn("Name"),
		},
	}
}
