package ninjaone

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaNetworkInterface struct {
	DeviceId       *int64          `json:"deviceId"`
	AdapterName    flexString      `json:"adapterName"`
	DefaultGateway flexString      `json:"defaultGateway"`
	DnsHostName    flexString      `json:"dnsHostName"`
	DnsServers     json.RawMessage `json:"dnsServers"`
	InterfaceIndex flexString      `json:"interfaceIndex"`
	InterfaceName  flexString      `json:"interfaceName"`
	InterfaceType  flexString      `json:"interfaceType"`
	IpAddress      json.RawMessage `json:"ipAddress"`
	LinkSpeed      flexString      `json:"linkSpeed"`
	MacAddress     json.RawMessage `json:"macAddress"`
	Mtu            flexString      `json:"mtu"`
	Status         flexString      `json:"status"`
	SubnetMask     flexString      `json:"subnetMask"`
}

//// TABLE DEFINITION

func tableNinjaoneNetworkInterface(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_network_interface",
		Description: "Network interfaces on each device.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaNetworkInterface]("ninjaone_network_interface", "/v2/queries/network-interfaces"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "adapter_name", Type: proto.ColumnType_STRING, Description: "Adapter name."},
			{Name: "default_gateway", Type: proto.ColumnType_STRING, Description: "Default gateway."},
			{Name: "dns_host_name", Type: proto.ColumnType_STRING, Description: "DNS host name."},
			{Name: "dns_servers", Type: proto.ColumnType_JSON, Description: "DNS servers."},
			{Name: "interface_index", Type: proto.ColumnType_STRING, Description: "Interface index."},
			{Name: "interface_name", Type: proto.ColumnType_STRING, Description: "Interface name."},
			{Name: "interface_type", Type: proto.ColumnType_STRING, Description: "Interface type."},
			{Name: "ip_address", Type: proto.ColumnType_JSON, Description: "IP addresses."},
			{Name: "link_speed", Type: proto.ColumnType_STRING, Description: "Link speed."},
			{Name: "mac_address", Type: proto.ColumnType_JSON, Description: "MAC addresses."},
			{Name: "mtu", Type: proto.ColumnType_STRING, Description: "Maximum transmission unit."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Interface status."},
			{Name: "subnet_mask", Type: proto.ColumnType_STRING, Description: "Subnet mask."},
			titleColumn("InterfaceName"),
		},
	}
}
