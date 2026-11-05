package ninjaone

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaPolicyOverride struct {
	DeviceId  *int64          `json:"deviceId"`
	Overrides json.RawMessage `json:"overrides"`
}

//// TABLE DEFINITION

func tableNinjaonePolicyOverride(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_policy_override",
		Description: "Policy settings overridden on individual devices.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaPolicyOverride]("ninjaone_policy_override", "/v2/queries/policy-overrides"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "overrides", Type: proto.ColumnType_JSON, Description: "Names of the overridden sections."},
		},
	}
}
