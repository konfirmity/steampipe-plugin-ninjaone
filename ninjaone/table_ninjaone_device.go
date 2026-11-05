package ninjaone

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

type ninjaDeviceOS struct {
	Manufacturer flexString `json:"manufacturer"`
	Name         flexString `json:"name"`
	Architecture flexString `json:"architecture"`
	BuildNumber  flexString `json:"buildNumber"`
	ReleaseId    flexString `json:"releaseId"`
	LastBootTime *float64   `json:"lastBootTime"`
	NeedsReboot  *bool      `json:"needsReboot"`
}

type ninjaDeviceSystem struct {
	Manufacturer        flexString `json:"manufacturer"`
	Model               flexString `json:"model"`
	BiosSerialNumber    flexString `json:"biosSerialNumber"`
	SerialNumber        flexString `json:"serialNumber"`
	AssetSerialNumber   flexString `json:"assetSerialNumber"`
	Domain              flexString `json:"domain"`
	DomainRole          flexString `json:"domainRole"`
	TotalPhysicalMemory *int64     `json:"totalPhysicalMemory"`
	VirtualMachine      *bool      `json:"virtualMachine"`
	ChassisType         flexString `json:"chassisType"`
}

type ninjaDevice struct {
	Id               int64             `json:"id"`
	Uid              flexString        `json:"uid"`
	OrganizationId   int64             `json:"organizationId"`
	LocationId       *int64            `json:"locationId"`
	ParentDeviceId   *int64            `json:"parentDeviceId"`
	AssignedOwnerUid flexString        `json:"assignedOwnerUid"`
	NodeClass        flexString        `json:"nodeClass"`
	NodeRoleId       *int64            `json:"nodeRoleId"`
	RolePolicyId     *int64            `json:"rolePolicyId"`
	PolicyId         *int64            `json:"policyId"`
	ApprovalStatus   flexString        `json:"approvalStatus"`
	Offline          *bool             `json:"offline"`
	DisplayName      flexString        `json:"displayName"`
	SystemName       flexString        `json:"systemName"`
	DnsName          flexString        `json:"dnsName"`
	NetbiosName      flexString        `json:"netbiosName"`
	Created          *float64          `json:"created"`
	LastContact      *float64          `json:"lastContact"`
	LastUpdate       *float64          `json:"lastUpdate"`
	Tags             []string          `json:"tags"`
	IpAddresses      []string          `json:"ipAddresses"`
	MacAddresses     []string          `json:"macAddresses"`
	PublicIp         flexString        `json:"publicIP"`
	LastLoggedInUser flexString        `json:"lastLoggedInUser"`
	DeviceType       flexString        `json:"deviceType"`
	Os               ninjaDeviceOS     `json:"os"`
	System           ninjaDeviceSystem `json:"system"`
	Memory           json.RawMessage   `json:"memory"`
	Processors       json.RawMessage   `json:"processors"`
	Volumes          json.RawMessage   `json:"volumes"`
	Maintenance      json.RawMessage   `json:"maintenance"`
}

// Akas returns the device's globally unique identifiers.
func (n ninjaDevice) Akas() []string {
	return []string{"ninjaone:device:" + string(n.Uid)}
}

// Title returns the name a technician sees for the device.
func (n ninjaDevice) Title() string {
	for _, name := range []flexString{n.DisplayName, n.SystemName, n.DnsName} {
		if name != "" {
			return string(name)
		}
	}
	return fmt.Sprintf("%d", n.Id)
}

//// TABLE DEFINITION

func tableNinjaoneDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_device",
		Description: "Devices managed by NinjaOne, with their operating system, hardware, network and last check-in.",
		List: &plugin.ListConfig{
			Hydrate: listDevices,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDevice,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Device identifier."},
			{Name: "approval_status", Type: proto.ColumnType_STRING, Description: "Approval state: PENDING, STAGED, APPROVED or DECOMMISSIONED."},
			{Name: "asset_serial_number", Type: proto.ColumnType_STRING, Description: "Asset serial number.", Transform: field("System.AssetSerialNumber")},
			{Name: "assigned_owner_uid", Type: proto.ColumnType_STRING, Description: "Unique identifier of the user assigned as the device owner."},
			{Name: "bios_serial_number", Type: proto.ColumnType_STRING, Description: "BIOS serial number. Empty on Macs.", Transform: field("System.BiosSerialNumber")},
			{Name: "chassis_type", Type: proto.ColumnType_STRING, Description: "Reported chassis type. Macs have been seen reporting DESKTOP for a laptop.", Transform: field("System.ChassisType")},
			timeColumn("created", "Time the device was enrolled."),
			{Name: "device_type", Type: proto.ColumnType_STRING, Description: "Kind of device record, such as AgentDevice."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Display name set on the device."},
			{Name: "dns_name", Type: proto.ColumnType_STRING, Description: "DNS name."},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Domain the device has joined.", Transform: field("System.Domain")},
			{Name: "domain_role", Type: proto.ColumnType_STRING, Description: "Role of the device in its domain.", Transform: field("System.DomainRole")},
			{Name: "ip_addresses", Type: proto.ColumnType_JSON, Description: "IP addresses reported by the device."},
			timeColumn("last_contact", "Last time the agent contacted NinjaOne."),
			{Name: "last_logged_in_user", Type: proto.ColumnType_STRING, Description: "Last user who logged in. A local login name, not an email address."},
			timeColumn("last_update", "Last time the device record was updated."),
			{Name: "location_id", Type: proto.ColumnType_INT, Description: "Identifier of the device's location."},
			{Name: "mac_addresses", Type: proto.ColumnType_JSON, Description: "MAC addresses reported by the device."},
			{Name: "maintenance", Type: proto.ColumnType_JSON, Description: "Maintenance window state."},
			{Name: "manufacturer", Type: proto.ColumnType_STRING, Description: "Hardware manufacturer.", Transform: field("System.Manufacturer")},
			{Name: "memory", Type: proto.ColumnType_JSON, Description: "Installed memory."},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "Hardware model.", Transform: field("System.Model")},
			{Name: "netbios_name", Type: proto.ColumnType_STRING, Description: "NetBIOS name."},
			{Name: "node_class", Type: proto.ColumnType_STRING, Description: "Device class, such as WINDOWS_WORKSTATION, WINDOWS_SERVER or MAC."},
			{Name: "node_role_id", Type: proto.ColumnType_INT, Description: "Identifier of the device role."},
			{Name: "offline", Type: proto.ColumnType_BOOL, Description: "True if the agent is offline."},
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: "Identifier of the organization the device belongs to."},
			{Name: "os_architecture", Type: proto.ColumnType_STRING, Description: "Operating system architecture.", Transform: field("Os.Architecture")},
			{Name: "os_build_number", Type: proto.ColumnType_STRING, Description: "Operating system build number.", Transform: field("Os.BuildNumber")},
			{Name: "os_last_boot_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time of the last boot.", Transform: transform.FromField("Os.LastBootTime").Transform(epochToTime)},
			{Name: "os_manufacturer", Type: proto.ColumnType_STRING, Description: "Operating system vendor.", Transform: field("Os.Manufacturer")},
			{Name: "os_name", Type: proto.ColumnType_STRING, Description: "Operating system name and version.", Transform: field("Os.Name")},
			{Name: "os_needs_reboot", Type: proto.ColumnType_BOOL, Description: "True if the operating system is waiting for a reboot.", Transform: transform.FromField("Os.NeedsReboot")},
			{Name: "os_release_id", Type: proto.ColumnType_STRING, Description: "Operating system release identifier.", Transform: field("Os.ReleaseId")},
			{Name: "parent_device_id", Type: proto.ColumnType_INT, Description: "Identifier of the parent device, if any."},
			{Name: "policy_id", Type: proto.ColumnType_INT, Description: "Identifier of a policy set directly on the device."},
			{Name: "processors", Type: proto.ColumnType_JSON, Description: "Installed processors."},
			{Name: "public_ip", Type: proto.ColumnType_IPADDR, Description: "Public IP address last seen."},
			{Name: "role_policy_id", Type: proto.ColumnType_INT, Description: "Identifier of the policy that comes from the device role."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Hardware serial number.", Transform: field("System.SerialNumber")},
			{Name: "system_name", Type: proto.ColumnType_STRING, Description: "Host name reported by the device."},
			{Name: "total_physical_memory", Type: proto.ColumnType_INT, Description: "Installed memory in bytes.", Transform: transform.FromField("System.TotalPhysicalMemory")},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "Unique device identifier."},
			{Name: "virtual_machine", Type: proto.ColumnType_BOOL, Description: "True if the device is a virtual machine.", Transform: transform.FromField("System.VirtualMachine")},
			{Name: "volumes", Type: proto.ColumnType_JSON, Description: "Storage volumes."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromMethod("Title")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "Labels on the device as a map of label to true.", Transform: transform.FromField("Tags").Transform(labelsToTags)},
			{Name: "tags_raw", Type: proto.ColumnType_JSON, Description: "Labels on the device as NinjaOne returns them.", Transform: transform.FromField("Tags")},
			akasColumn(),
		},
	}
}

//// LIST FUNCTION

func listDevices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_device.listDevices", "connection_error", err)
		return nil, err
	}
	if err := c.requireOrganization(ctx); err != nil {
		plugin.Logger(ctx).Error("ninjaone_device.listDevices", "api_error", err)
		return nil, err
	}
	stream := streamRows[ninjaDevice](ctx, d)
	err = c.listAfter(ctx, "/v2/devices-detailed", c.deviceFilter(nil), d.QueryContext.GetLimit(), func(raw json.RawMessage) (bool, error) {
		var device struct {
			OrganizationID int64 `json:"organizationId"`
		}
		if err := json.Unmarshal(raw, &device); err != nil {
			return false, fmt.Errorf("decode device: %w", err)
		}
		if !c.inScope(device.OrganizationID) {
			return true, nil
		}
		return stream(raw)
	})
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_device.listDevices", "api_error", err)
		return nil, err
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDevice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ninjaone_device.getDevice", "connection_error", err)
		return nil, err
	}
	id := d.EqualsQuals["id"].GetInt64Value()
	var device ninjaDevice
	if err := c.get(ctx, fmt.Sprintf("/v2/device/%d", id), nil, &device); err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("ninjaone_device.getDevice", "api_error", err)
		return nil, err
	}
	if !c.inScope(device.OrganizationId) {
		return nil, nil
	}
	return device, nil
}
