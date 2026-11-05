package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaAntivirusThreat struct {
	DeviceId        *int64     `json:"deviceId"`
	Category        flexString `json:"category"`
	DetectionSource flexString `json:"detectionSource"`
	Level           flexString `json:"level"`
	Name            flexString `json:"name"`
	ProductCode     flexString `json:"productCode"`
	QuarantineId    flexString `json:"quarantineId"`
	Status          flexString `json:"status"`
	ThreatId        flexString `json:"threatId"`
	Timestamp       *float64   `json:"timestamp"`
	TraceList       flexString `json:"traceList"`
	Type            flexString `json:"type"`
}

//// TABLE DEFINITION

func tableNinjaoneAntivirusThreat(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_antivirus_threat",
		Description: "Threats detected by antivirus on each device.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaAntivirusThreat]("ninjaone_antivirus_threat", "/v2/queries/antivirus-threats"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Threat category."},
			{Name: "detection_source", Type: proto.ColumnType_STRING, Description: "What detected the threat."},
			{Name: "level", Type: proto.ColumnType_STRING, Description: "Threat level."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Threat name."},
			{Name: "product_code", Type: proto.ColumnType_STRING, Description: "Antivirus product code."},
			{Name: "quarantine_id", Type: proto.ColumnType_STRING, Description: "Quarantine identifier."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Threat status."},
			{Name: "threat_id", Type: proto.ColumnType_STRING, Description: "Threat identifier."},
			timeColumn("timestamp", "Time the device reported this data."),
			{Name: "trace_list", Type: proto.ColumnType_STRING, Description: "Files and locations involved."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Threat type."},
			titleColumn("Name"),
		},
	}
}
