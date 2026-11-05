package ninjaone

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaAlert struct {
	Uid             flexString      `json:"uid"`
	DeviceId        *int64          `json:"deviceId"`
	Message         flexString      `json:"message"`
	Subject         flexString      `json:"subject"`
	SourceType      flexString      `json:"sourceType"`
	SourceName      flexString      `json:"sourceName"`
	SourceConfigUid flexString      `json:"sourceConfigUid"`
	CreateTime      *float64        `json:"createTime"`
	UpdateTime      *float64        `json:"updateTime"`
	Data            json.RawMessage `json:"data"`
}

//// TABLE DEFINITION

func tableNinjaoneAlert(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_alert",
		Description: "Active NinjaOne alerts on devices.",
		List:        &plugin.ListConfig{Hydrate: listAlerts, KeyColumns: plugin.OptionalColumns([]string{"device_id"})},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device the alert is on."},
			timeColumn("create_time", "Time the alert was raised."),
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Alert detail."},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "Alert message."},
			{Name: "source_config_uid", Type: proto.ColumnType_STRING, Description: "Identifier of the condition that raised the alert."},
			{Name: "source_name", Type: proto.ColumnType_STRING, Description: "Name of the condition that raised the alert."},
			{Name: "source_type", Type: proto.ColumnType_STRING, Description: "Kind of condition that raised the alert."},
			{Name: "subject", Type: proto.ColumnType_STRING, Description: "Alert subject."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "Unique alert identifier."},
			timeColumn("update_time", "Time the alert last changed."),
			titleColumn("Subject"),
		},
	}
}

//// LIST FUNCTION

func listAlerts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_alert.listAlerts", "connection_error", err)
		return nil, err
	}
	allowed, err := c.scopedDevices(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_alert.listAlerts", "api_error", err)
		return nil, err
	}
	var rows []json.RawMessage
	if err := c.get(ctx, "/v2/alerts", filterQuery(c.deviceFilter(deviceIDQual(d))), &rows); err != nil {
		plugin.Logger(ctx).Error("ninjaone_alert.listAlerts", "api_error", err)
		return nil, err
	}
	stream := onlyDevices(allowed, streamRows[ninjaAlert](ctx, d))
	for _, row := range rows {
		more, err := stream(row)
		if err != nil {
			return nil, err
		}
		if !more {
			break
		}
	}
	return nil, nil
}
