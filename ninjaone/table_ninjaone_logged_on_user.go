package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaLoggedOnUser struct {
	DeviceId  *int64     `json:"deviceId"`
	LogonTime *float64   `json:"logonTime"`
	UserName  flexString `json:"userName"`
}

//// TABLE DEFINITION

func tableNinjaoneLoggedOnUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_logged_on_user",
		Description: "Users logged on to each device. These are local login names, not email addresses.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaLoggedOnUser]("ninjaone_logged_on_user", "/v2/queries/logged-on-users"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			timeColumn("logon_time", "Time of the logon."),
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "Local login name."},
			titleColumn("UserName"),
		},
	}
}
