package ninjaone

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaVolume struct {
	DeviceId        *int64          `json:"deviceId"`
	AutoMount       *bool           `json:"autoMount"`
	BitLockerStatus json.RawMessage `json:"bitLockerStatus"`
	Capacity        *int64          `json:"capacity"`
	Compressed      *bool           `json:"compressed"`
	DeviceType      flexString      `json:"deviceType"`
	DriveLetter     flexString      `json:"driveLetter"`
	FileSystem      flexString      `json:"fileSystem"`
	FreeSpace       *int64          `json:"freeSpace"`
	Label           flexString      `json:"label"`
	Name            flexString      `json:"name"`
	SerialNumber    flexString      `json:"serialNumber"`
	Timestamp       *float64        `json:"timestamp"`
}

//// TABLE DEFINITION

func tableNinjaoneVolume(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_volume",
		Description: "Storage volumes on each device, with BitLocker encryption state where Windows reports it.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaVolume]("ninjaone_volume", "/v2/queries/volumes"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "auto_mount", Type: proto.ColumnType_BOOL, Description: "True if the volume mounts automatically."},
			{Name: "bit_locker_status", Type: proto.ColumnType_JSON, Description: "BitLocker conversion, encryption, protection and lock state. Null outside Windows."},
			{Name: "capacity", Type: proto.ColumnType_INT, Description: "Capacity in bytes."},
			{Name: "compressed", Type: proto.ColumnType_BOOL, Description: "True if the volume is compressed."},
			{Name: "device_type", Type: proto.ColumnType_STRING, Description: "Kind of volume."},
			{Name: "drive_letter", Type: proto.ColumnType_STRING, Description: "Drive letter."},
			{Name: "file_system", Type: proto.ColumnType_STRING, Description: "File system."},
			{Name: "free_space", Type: proto.ColumnType_INT, Description: "Free space in bytes."},
			{Name: "label", Type: proto.ColumnType_STRING, Description: "Volume label."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Volume name."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Volume serial number."},
			timeColumn("timestamp", "Time the device reported this data."),
			titleColumn("Name"),
		},
	}
}
