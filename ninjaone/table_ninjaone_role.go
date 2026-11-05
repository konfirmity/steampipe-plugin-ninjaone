package ninjaone

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaRole struct {
	Id               int64    `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	NodeClass        string   `json:"nodeClass"`
	Custom           *bool    `json:"custom"`
	ChassisType      string   `json:"chassisType"`
	Created          *float64 `json:"created"`
	NodeRoleParentId *int64   `json:"nodeRoleParentId"`
}

// Akas returns the role's globally unique identifiers.
func (r ninjaRole) Akas() []string {
	return []string{fmt.Sprintf("ninjaone:role:%d", r.Id)}
}

//// TABLE DEFINITION

func tableNinjaoneRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_role",
		Description: "NinjaOne device roles. A connection with organization_id set returns only roles that organization's devices use.",
		List:        &plugin.ListConfig{Hydrate: listRoles},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Role identifier."},
			{Name: "chassis_type", Type: proto.ColumnType_STRING, Description: "Chassis type the role applies to."},
			timeColumn("created", "Time the role was created."),
			{Name: "custom", Type: proto.ColumnType_BOOL, Description: "True if the role was created by the tenant."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Role description."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Role name."},
			{Name: "node_class", Type: proto.ColumnType_STRING, Description: "Device class the role applies to."},
			{Name: "node_role_parent_id", Type: proto.ColumnType_INT, Description: "Identifier of the parent role."},
			titleColumn("Name"),
			akasColumn(),
		},
	}
}

//// LIST FUNCTION

func listRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_role.listRoles", "connection_error", err)
		return nil, err
	}
	var allowed map[int64]bool
	if c.organizationID != nil {
		if _, allowed, err = c.deviceReferences(ctx); err != nil {
			plugin.Logger(ctx).Error("ninjaone_role.listRoles", "api_error", err)
			return nil, err
		}
	}
	var rows []json.RawMessage
	if err := c.get(ctx, "/v2/roles", nil, &rows); err != nil {
		plugin.Logger(ctx).Error("ninjaone_role.listRoles", "api_error", err)
		return nil, err
	}
	for _, raw := range rows {
		var role ninjaRole
		if err := json.Unmarshal(raw, &role); err != nil {
			return nil, fmt.Errorf("decode role: %w", err)
		}
		if allowed == nil || allowed[role.Id] {
			d.StreamListItem(ctx, role)
		}
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}
	return nil, nil
}
