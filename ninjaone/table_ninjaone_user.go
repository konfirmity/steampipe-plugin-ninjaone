package ninjaone

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

type ninjaUser struct {
	Id                int64    `json:"id"`
	Uid               string   `json:"uid"`
	FirstName         string   `json:"firstName"`
	LastName          string   `json:"lastName"`
	Email             string   `json:"email"`
	Phone             string   `json:"phone"`
	Enabled           *bool    `json:"enabled"`
	Administrator     *bool    `json:"administrator"`
	PermitAllClients  *bool    `json:"permitAllClients"`
	NotifyAllClients  *bool    `json:"notifyAllClients"`
	MustChangePw      *bool    `json:"mustChangePw"`
	MfaConfigured     *bool    `json:"mfaConfigured"`
	UserType          string   `json:"userType"`
	InvitationStatus  string   `json:"invitationStatus"`
	OrganizationId    *int64   `json:"organizationId"`
	DeviceIds         []int64  `json:"deviceIds"`
	AssignedDeviceIds []int64  `json:"assignedDeviceIds"`
	Roles             []string `json:"roles"`
}

// Akas returns the user's globally unique identifiers.
func (u ninjaUser) Akas() []string {
	return []string{"ninjaone:user:" + u.Uid}
}

// Title returns the user's full name, or their email when no name is set.
func (u ninjaUser) Title() string {
	if name := strings.TrimSpace(u.FirstName + " " + u.LastName); name != "" {
		return name
	}
	return u.Email
}

//// TABLE DEFINITION

func tableNinjaoneUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_user",
		Description: "NinjaOne users. A connection with organization_id set returns that organization's end users only, never the provider's technicians.",
		List:        &plugin.ListConfig{Hydrate: listUsers},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "User identifier."},
			{Name: "administrator", Type: proto.ColumnType_BOOL, Description: "True if the user is a system administrator."},
			{Name: "assigned_device_ids", Type: proto.ColumnType_JSON, Description: "Identifiers of devices the user is assigned to as owner."},
			{Name: "device_ids", Type: proto.ColumnType_JSON, Description: "Identifiers of devices the user can access."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "True if the account is enabled."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "First name."},
			{Name: "invitation_status", Type: proto.ColumnType_STRING, Description: "Invitation state: REGISTERED, PENDING or EXPIRED."},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "Last name."},
			{Name: "mfa_configured", Type: proto.ColumnType_BOOL, Description: "True if multi-factor authentication is set up."},
			{Name: "must_change_pw", Type: proto.ColumnType_BOOL, Description: "True if the user must change their password."},
			{Name: "notify_all_clients", Type: proto.ColumnType_BOOL, Description: "True if the user receives notifications for all organizations."},
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: "Identifier of the end user's organization. Null for technicians."},
			{Name: "permit_all_clients", Type: proto.ColumnType_BOOL, Description: "True if the user can reach all organizations."},
			{Name: "phone", Type: proto.ColumnType_STRING, Description: "Phone number."},
			{Name: "roles", Type: proto.ColumnType_JSON, Description: "Names of the user's roles."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "Unique user identifier."},
			{Name: "user_type", Type: proto.ColumnType_STRING, Description: "TECHNICIAN or END_USER."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromMethod("Title")},
			akasColumn(),
		},
	}
}

//// LIST FUNCTION

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_user.listUsers", "connection_error", err)
		return nil, err
	}
	path := "/v2/users"
	if c.organizationID != nil {
		path = fmt.Sprintf("/v2/organization/%d/end-users", *c.organizationID)
	}
	var rows []json.RawMessage
	if err := c.get(ctx, path, nil, &rows); err != nil {
		plugin.Logger(ctx).Error("ninjaone_user.listUsers", "api_error", err)
		return nil, err
	}
	for _, raw := range rows {
		var user ninjaUser
		if err := json.Unmarshal(raw, &user); err != nil {
			return nil, fmt.Errorf("decode user: %w", err)
		}
		if c.organizationID != nil && (user.UserType == "TECHNICIAN" || (user.OrganizationId != nil && *user.OrganizationId != *c.organizationID)) {
			continue
		}
		d.StreamListItem(ctx, user)
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}
	return nil, nil
}
