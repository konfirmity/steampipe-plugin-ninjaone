package ninjaone

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

// flexString decodes a JSON string or number. Report rows are decoded with it
// because one mistyped field would otherwise fail the whole table.
type flexString string

// UnmarshalJSON accepts a string, a number or null.
func (s *flexString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*s = ""
		return nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*s = flexString(text)
		return nil
	}
	*s = flexString(strings.TrimSpace(string(data)))
	return nil
}

// nullIfEmpty turns an empty string into null and leaves false and 0 alone:
// a zero threat count is a fact, while SDK NullIfZero would erase it.
func nullIfEmpty(_ context.Context, d *transform.TransformData) (interface{}, error) {
	switch value := d.Value.(type) {
	case flexString:
		if value == "" {
			return nil, nil
		}
		return string(value), nil
	case string:
		if value == "" {
			return nil, nil
		}
	case json.RawMessage:
		if len(value) == 0 || string(value) == "null" {
			return nil, nil
		}
	}
	return d.Value, nil
}

const millisecondsThreshold = 1e11

func epochToTime(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var seconds float64
	switch value := d.Value.(type) {
	case nil:
		return nil, nil
	case float64:
		seconds = value
	case *float64:
		if value == nil {
			return nil, nil
		}
		seconds = *value
	default:
		return nil, fmt.Errorf("unexpected timestamp type %T", d.Value)
	}
	if seconds == 0 {
		return nil, nil
	}
	// logged-on-users sends milliseconds where every other report sends seconds.
	if seconds > millisecondsThreshold {
		seconds /= 1000
	}
	whole, fraction := math.Modf(seconds)
	return time.Unix(int64(whole), int64(fraction*float64(time.Second))).UTC(), nil
}

func labelsToTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	labels, ok := d.Value.([]string)
	if !ok || len(labels) == 0 {
		return nil, nil
	}
	tags := make(map[string]bool, len(labels))
	for _, label := range labels {
		tags[label] = true
	}
	return tags, nil
}

func field(name string) *transform.ColumnTransforms {
	return transform.FromField(name).Transform(nullIfEmpty)
}

func timeColumn(name, description string) *plugin.Column {
	return &plugin.Column{Name: name, Type: proto.ColumnType_TIMESTAMP, Description: description, Transform: transform.FromCamel().Transform(epochToTime)}
}

func titleColumn(fieldName string) *plugin.Column {
	return &plugin.Column{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: field(fieldName)}
}

func akasColumn() *plugin.Column {
	return &plugin.Column{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifiers for the resource.", Transform: transform.FromMethod("Akas")}
}

func deviceIDQual(d *plugin.QueryData) *int64 {
	qual, ok := d.EqualsQuals["device_id"]
	if !ok || qual == nil {
		return nil
	}
	id := qual.GetInt64Value()
	return &id
}

// listReport streams a /v2/queries report as rows of T.
func listReport[T any](table, path string) plugin.HydrateFunc {
	return func(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
		c, err := connect(ctx, d)
		if err != nil {
			plugin.Logger(ctx).Error(table+".list", "connection_error", err)
			return nil, err
		}
		allowed, err := c.scopedDevices(ctx)
		if err != nil {
			plugin.Logger(ctx).Error(table+".list", "api_error", err)
			return nil, err
		}
		err = c.listCursor(ctx, path, c.deviceFilter(deviceIDQual(d)), d.QueryContext.GetLimit(), onlyDevices(allowed, streamRows[T](ctx, d)))
		if err != nil {
			plugin.Logger(ctx).Error(table+".list", "api_error", err)
			return nil, err
		}
		return nil, nil
	}
}

func streamRows[T any](ctx context.Context, d *plugin.QueryData) func(json.RawMessage) (bool, error) {
	return func(raw json.RawMessage) (bool, error) {
		var row T
		if err := json.Unmarshal(raw, &row); err != nil {
			return false, fmt.Errorf("decode row: %w", err)
		}
		d.StreamListItem(ctx, row)
		return d.RowsRemaining(ctx) != 0, nil
	}
}

// deviceReferences names the policies and roles the organization's devices use.
// /v2/policies and /v2/roles span every tenant, and an organization's own policy
// map lists every role in the account, so only the devices narrow them.
func (c *client) deviceReferences(ctx context.Context) (policies, roles map[int64]bool, err error) {
	if err := c.requireOrganization(ctx); err != nil {
		return nil, nil, err
	}
	policies, roles = map[int64]bool{}, map[int64]bool{}
	err = c.listAfter(ctx, "/v2/devices-detailed", c.deviceFilter(nil), -1, func(raw json.RawMessage) (bool, error) {
		var device ninjaDevice
		if err := json.Unmarshal(raw, &device); err != nil {
			return false, fmt.Errorf("decode device: %w", err)
		}
		if !c.inScope(device.OrganizationId) {
			return true, nil
		}
		for _, id := range []*int64{device.PolicyId, device.RolePolicyId} {
			if id != nil {
				policies[*id] = true
			}
		}
		if device.NodeRoleId != nil {
			roles[*device.NodeRoleId] = true
		}
		return true, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return policies, roles, nil
}
