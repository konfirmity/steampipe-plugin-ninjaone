package ninjaone

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaOrganization struct {
	Id               int64           `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	NodeApprovalMode string          `json:"nodeApprovalMode"`
	Locations        json.RawMessage `json:"locations"`
	Policies         json.RawMessage `json:"policies"`
	Settings         json.RawMessage `json:"settings"`
}

// Akas returns the organization's globally unique identifiers.
func (o ninjaOrganization) Akas() []string {
	return []string{fmt.Sprintf("ninjaone:organization:%d", o.Id)}
}

//// TABLE DEFINITION

func tableNinjaoneOrganization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_organization",
		Description: "NinjaOne organizations, one per client of a managed service provider. A connection with organization_id set returns only that organization.",
		List:        &plugin.ListConfig{Hydrate: listOrganizations},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Organization identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Organization name."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Organization description."},
			{Name: "locations", Type: proto.ColumnType_JSON, Description: "Locations in the organization."},
			{Name: "node_approval_mode", Type: proto.ColumnType_STRING, Description: "How new devices are approved: AUTOMATIC, MANUAL or REJECT."},
			{Name: "policies", Type: proto.ColumnType_JSON, Description: "Device role to policy assignments."},
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "Organization settings."},
			titleColumn("Name"),
			akasColumn(),
		},
	}
}

//// LIST FUNCTION

func listOrganizations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_organization.listOrganizations", "connection_error", err)
		return nil, err
	}
	if err := c.requireOrganization(ctx); err != nil {
		plugin.Logger(ctx).Error("ninjaone_organization.listOrganizations", "api_error", err)
		return nil, err
	}
	err = c.listAfter(ctx, "/v2/organizations-detailed", "", -1, func(raw json.RawMessage) (bool, error) {
		var organization ninjaOrganization
		if err := json.Unmarshal(raw, &organization); err != nil {
			return false, fmt.Errorf("decode organization: %w", err)
		}
		if c.inScope(organization.Id) {
			d.StreamListItem(ctx, organization)
		}
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_organization.listOrganizations", "api_error", err)
		return nil, err
	}
	return nil, nil
}
