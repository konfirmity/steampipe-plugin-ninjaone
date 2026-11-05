package ninjaone

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaActivity struct {
	Id             int64           `json:"id"`
	ActivityTime   *float64        `json:"activityTime"`
	DeviceId       *int64          `json:"deviceId"`
	ActivityType   flexString      `json:"activityType"`
	StatusCode     flexString      `json:"statusCode"`
	Status         flexString      `json:"status"`
	ActivityResult flexString      `json:"activityResult"`
	UserId         *int64          `json:"userId"`
	Message        flexString      `json:"message"`
	Type           flexString      `json:"type"`
	Data           json.RawMessage `json:"data"`
}

//// TABLE DEFINITION

func tableNinjaoneActivity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_activity",
		Description: "The NinjaOne activity log: patching, scripts, logins and other events, newest first. Filter on id > n to read only what is new.",
		List: &plugin.ListConfig{
			Hydrate: listActivities,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "device_id", Require: plugin.Optional},
				{Name: "id", Require: plugin.Optional, Operators: []string{">"}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Activity identifier."},
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device the activity is on."},
			{Name: "activity_result", Type: proto.ColumnType_STRING, Description: "Outcome of the activity."},
			timeColumn("activity_time", "Time of the activity."),
			{Name: "activity_type", Type: proto.ColumnType_STRING, Description: "Category of activity, such as PATCH_MANAGEMENT or SYSTEM."},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Activity detail."},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "Human-readable description."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Activity status."},
			{Name: "status_code", Type: proto.ColumnType_STRING, Description: "Activity status code."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Specific activity type."},
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "Identifier of the user who triggered the activity."},
			titleColumn("Message"),
		},
	}
}

//// LIST FUNCTION

func listActivities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_activity.listActivities", "connection_error", err)
		return nil, err
	}
	allowed, err := c.scopedDevices(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_activity.listActivities", "api_error", err)
		return nil, err
	}
	var newerThan *int64
	if quals := d.Quals["id"]; quals != nil {
		for _, qual := range quals.Quals {
			if qual.Operator == ">" {
				id := qual.Value.GetInt64Value()
				newerThan = &id
			}
		}
	}
	err = c.listActivities(ctx, c.deviceFilter(deviceIDQual(d)), newerThan, d.QueryContext.GetLimit(), onlyDevices(allowed, streamRows[ninjaActivity](ctx, d)))
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_activity.listActivities", "api_error", err)
		return nil, err
	}
	return nil, nil
}
