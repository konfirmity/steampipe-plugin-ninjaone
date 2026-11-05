package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaOperatingSystem struct {
	DeviceId                *int64     `json:"deviceId"`
	Architecture            flexString `json:"architecture"`
	BuildNumber             flexString `json:"buildNumber"`
	Language                flexString `json:"language"`
	LastBootTime            *float64   `json:"lastBootTime"`
	Locale                  flexString `json:"locale"`
	Manufacturer            flexString `json:"manufacturer"`
	Name                    flexString `json:"name"`
	NeedsReboot             *bool      `json:"needsReboot"`
	ReleaseId               flexString `json:"releaseId"`
	ServicePackMajorVersion *int64     `json:"servicePackMajorVersion"`
	ServicePackMinorVersion *int64     `json:"servicePackMinorVersion"`
	Timestamp               *float64   `json:"timestamp"`
}

//// TABLE DEFINITION

func tableNinjaoneOperatingSystem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_operating_system",
		Description: "Operating system on each device.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaOperatingSystem]("ninjaone_operating_system", "/v2/queries/operating-systems"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "architecture", Type: proto.ColumnType_STRING, Description: "Operating system architecture."},
			{Name: "build_number", Type: proto.ColumnType_STRING, Description: "Build number."},
			{Name: "language", Type: proto.ColumnType_STRING, Description: "Language."},
			timeColumn("last_boot_time", "Time of the last boot."),
			{Name: "locale", Type: proto.ColumnType_STRING, Description: "Locale."},
			{Name: "manufacturer", Type: proto.ColumnType_STRING, Description: "Operating system vendor."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Operating system name and version."},
			{Name: "needs_reboot", Type: proto.ColumnType_BOOL, Description: "True if the operating system is waiting for a reboot."},
			{Name: "release_id", Type: proto.ColumnType_STRING, Description: "Release identifier."},
			{Name: "service_pack_major_version", Type: proto.ColumnType_INT, Description: "Service pack major version."},
			{Name: "service_pack_minor_version", Type: proto.ColumnType_INT, Description: "Service pack minor version."},
			timeColumn("timestamp", "Time the device reported this data."),
			titleColumn("Name"),
		},
	}
}
