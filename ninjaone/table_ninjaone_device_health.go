package ninjaone

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type ninjaDeviceHealth struct {
	DeviceId                     *int64          `json:"deviceId"`
	ActiveJobCount               *int64          `json:"activeJobCount"`
	ActiveThreatsCount           *int64          `json:"activeThreatsCount"`
	AlertCount                   *int64          `json:"alertCount"`
	AvInstallStatus              flexString      `json:"avInstallStatus"`
	BlockedThreatsCount          *int64          `json:"blockedThreatsCount"`
	CriticalVulnerabilityCount   *int64          `json:"criticalVulnerabilityCount"`
	FailedOsPatchesCount         *int64          `json:"failedOSPatchesCount"`
	FailedSoftwarePatchesCount   *int64          `json:"failedSoftwarePatchesCount"`
	HealthStatus                 flexString      `json:"healthStatus"`
	HighVulnerabilityCount       *int64          `json:"highVulnerabilityCount"`
	InstallationIssuesCount      *int64          `json:"installationIssuesCount"`
	LowVulnerabilityCount        *int64          `json:"lowVulnerabilityCount"`
	MediumVulnerabilityCount     *int64          `json:"mediumVulnerabilityCount"`
	Offline                      *bool           `json:"offline"`
	ParentDeviceId               *int64          `json:"parentDeviceId"`
	ParentOffline                *bool           `json:"parentOffline"`
	PendingOsPatchesCount        *int64          `json:"pendingOSPatchesCount"`
	PendingRebootReason          flexString      `json:"pendingRebootReason"`
	PendingSoftwarePatchesCount  *int64          `json:"pendingSoftwarePatchesCount"`
	ProductsInstallationStatuses json.RawMessage `json:"productsInstallationStatuses"`
	QuarantinedThreatsCount      *int64          `json:"quarantinedThreatsCount"`
}

//// TABLE DEFINITION

func tableNinjaoneDeviceHealth(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ninjaone_device_health",
		Description: "Health summary per device: threat, patch, alert and vulnerability counts. A zero patch count can also mean patching is switched off in the device's policy.",
		List: &plugin.ListConfig{
			Hydrate:    listReport[ninjaDeviceHealth]("ninjaone_device_health", "/v2/queries/device-health"),
			KeyColumns: plugin.OptionalColumns([]string{"device_id"}),
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "Identifier of the device."},
			{Name: "active_job_count", Type: proto.ColumnType_INT, Description: "Jobs running on the device."},
			{Name: "active_threats_count", Type: proto.ColumnType_INT, Description: "Active antivirus threats."},
			{Name: "alert_count", Type: proto.ColumnType_INT, Description: "Open alerts."},
			{Name: "av_install_status", Type: proto.ColumnType_STRING, Description: "Antivirus installation state."},
			{Name: "blocked_threats_count", Type: proto.ColumnType_INT, Description: "Blocked antivirus threats."},
			{Name: "critical_vulnerability_count", Type: proto.ColumnType_INT, Description: "Critical vulnerabilities."},
			{Name: "failed_os_patches_count", Type: proto.ColumnType_INT, Description: "Operating system patches that failed."},
			{Name: "failed_software_patches_count", Type: proto.ColumnType_INT, Description: "Software patches that failed."},
			{Name: "health_status", Type: proto.ColumnType_STRING, Description: "Overall state: UNKNOWN, HEALTHY, NEEDS_ATTENTION or UNHEALTHY."},
			{Name: "high_vulnerability_count", Type: proto.ColumnType_INT, Description: "High vulnerabilities."},
			{Name: "installation_issues_count", Type: proto.ColumnType_INT, Description: "Product installation problems."},
			{Name: "low_vulnerability_count", Type: proto.ColumnType_INT, Description: "Low vulnerabilities."},
			{Name: "medium_vulnerability_count", Type: proto.ColumnType_INT, Description: "Medium vulnerabilities."},
			{Name: "offline", Type: proto.ColumnType_BOOL, Description: "True if the agent is offline."},
			{Name: "parent_device_id", Type: proto.ColumnType_INT, Description: "Identifier of the parent device."},
			{Name: "parent_offline", Type: proto.ColumnType_BOOL, Description: "True if the parent device is offline."},
			{Name: "pending_os_patches_count", Type: proto.ColumnType_INT, Description: "Operating system patches waiting to install."},
			{Name: "pending_reboot_reason", Type: proto.ColumnType_STRING, Description: "Why the device is waiting for a reboot."},
			{Name: "pending_software_patches_count", Type: proto.ColumnType_INT, Description: "Software patches waiting to install."},
			{Name: "products_installation_statuses", Type: proto.ColumnType_JSON, Description: "Installation state of each NinjaOne product."},
			{Name: "quarantined_threats_count", Type: proto.ColumnType_INT, Description: "Quarantined antivirus threats."},
		},
	}
}
