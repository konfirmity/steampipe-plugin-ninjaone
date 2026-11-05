package ninjaone

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaLocation struct {
	Id             int64  `json:"id"`
	OrganizationId int64  `json:"organizationId"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	Description    string `json:"description"`
}

// Akas returns the location's globally unique identifiers.
func (l ninjaLocation) Akas() []string {
	return []string{fmt.Sprintf("ninjaone:location:%d", l.Id)}
}

//// TABLE DEFINITION

func tableNinjaoneLocation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_location",
		Description: "Locations inside NinjaOne organizations.",
		List:        &plugin.ListConfig{Hydrate: listLocations},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Location identifier."},
			{Name: "address", Type: proto.ColumnType_STRING, Description: "Location address."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Location description."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Location name."},
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: "Identifier of the organization the location belongs to."},
			titleColumn("Name"),
			akasColumn(),
		},
	}
}

//// LIST FUNCTION

func listLocations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_location.listLocations", "connection_error", err)
		return nil, err
	}
	if err := c.requireOrganization(ctx); err != nil {
		plugin.Logger(ctx).Error("ninjaone_location.listLocations", "api_error", err)
		return nil, err
	}
	err = c.listAfter(ctx, "/v2/locations", "", -1, func(raw json.RawMessage) (bool, error) {
		var location ninjaLocation
		if err := json.Unmarshal(raw, &location); err != nil {
			return false, fmt.Errorf("decode location: %w", err)
		}
		if c.inScope(location.OrganizationId) {
			d.StreamListItem(ctx, location)
		}
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_location.listLocations", "api_error", err)
		return nil, err
	}
	return nil, nil
}
