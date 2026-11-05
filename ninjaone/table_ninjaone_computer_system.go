package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaComputerSystem struct {
	DeviceId            *int64     `json:"deviceId"`
	AssetSerialNumber   flexString `json:"assetSerialNumber"`
	BiosSerialNumber    flexString `json:"biosSerialNumber"`
	ChassisType         flexString `json:"chassisType"`
	Domain              flexString `json:"domain"`
	DomainRole          flexString `json:"domainRole"`
	Manufacturer        flexString `json:"manufacturer"`
	Model               flexString `json:"model"`
	Name                flexString `json:"name"`
	NumberOfProcessors  *int64     `json:"numberOfProcessors"`
	SerialNumber        flexString `json:"serialNumber"`
	Timestamp           *float64   `json:"timestamp"`
	TotalPhysicalMemory *int64     `json:"totalPhysicalMemory"`
	VirtualMachine      *bool      `json:"virtualMachine"`
}

//// TABLE DEFINITION

func tableNinjaoneComputerSystem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_computer_system",
		Description: "Hardware identity of each device.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaComputerSystem]("ninjaone_computer_system", "/v2/queries/computer-systems"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "asset_serial_number", Type: proto.ColumnType_STRING, Description: "Asset serial number."},
			{Name: "bios_serial_number", Type: proto.ColumnType_STRING, Description: "BIOS serial number. Empty on Macs."},
			{Name: "chassis_type", Type: proto.ColumnType_STRING, Description: "Reported chassis type."},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Domain the device has joined."},
			{Name: "domain_role", Type: proto.ColumnType_STRING, Description: "Role of the device in its domain."},
			{Name: "manufacturer", Type: proto.ColumnType_STRING, Description: "Hardware manufacturer."},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "Hardware model."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Host name."},
			{Name: "number_of_processors", Type: proto.ColumnType_INT, Description: "Number of processors."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Hardware serial number."},
			timeColumn("timestamp", "Time the device reported this data."),
			{Name: "total_physical_memory", Type: proto.ColumnType_INT, Description: "Installed memory in bytes."},
			{Name: "virtual_machine", Type: proto.ColumnType_BOOL, Description: "True if the device is a virtual machine."},
			titleColumn("Name"),
		},
	}
}
