package ninjaone

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaPolicy struct {
	Id               int64    `json:"id"`
	ParentPolicyId   *int64   `json:"parentPolicyId"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	NodeClass        string   `json:"nodeClass"`
	Updated          *float64 `json:"updated"`
	NodeClassDefault *bool    `json:"nodeClassDefault"`
}

// Akas returns the policy's globally unique identifiers.
func (p ninjaPolicy) Akas() []string {
	return []string{fmt.Sprintf("ninjaone:policy:%d", p.Id)}
}

//// TABLE DEFINITION

func tableNinjaonePolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_policy",
		Description: "NinjaOne device policies. A connection with organization_id set returns only policies that organization's devices use.",
		List:        &plugin.ListConfig{Hydrate: listPolicies},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Policy identifier."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Policy description."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Policy name."},
			{Name: "node_class", Type: proto.ColumnType_STRING, Description: "Device class the policy applies to."},
			{Name: "node_class_default", Type: proto.ColumnType_BOOL, Description: "True if this is the default policy for its device class."},
			{Name: "parent_policy_id", Type: proto.ColumnType_INT, Description: "Identifier of the policy this one inherits from."},
			timeColumn("updated", "Time the policy last changed."),
			titleColumn("Name"),
			akasColumn(),
		},
	}
}

//// LIST FUNCTION

func listPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_policy.listPolicies", "connection_error", err)
		return nil, err
	}
	var allowed map[int64]bool
	if c.organizationID != nil {
		if allowed, _, err = c.deviceReferences(ctx); err != nil {
			plugin.Logger(ctx).Error("ninjaone_policy.listPolicies", "api_error", err)
			return nil, err
		}
	}
	var rows []json.RawMessage
	if err := c.get(ctx, "/v2/policies", nil, &rows); err != nil {
		plugin.Logger(ctx).Error("ninjaone_policy.listPolicies", "api_error", err)
		return nil, err
	}
	for _, raw := range rows {
		var policy ninjaPolicy
		if err := json.Unmarshal(raw, &policy); err != nil {
			return nil, fmt.Errorf("decode policy: %w", err)
		}
		if allowed == nil || allowed[policy.Id] {
			d.StreamListItem(ctx, policy)
		}
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}
	return nil, nil
}
