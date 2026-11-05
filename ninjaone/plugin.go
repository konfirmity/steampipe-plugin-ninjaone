// Package ninjaone implements a Steampipe plugin for reading devices, patches,
// software, antivirus, health and users from the NinjaOne public API.
package ninjaone

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

// Plugin creates this (ninjaone) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	return &plugin.Plugin{
		Name:             "steampipe-plugin-ninjaone",
		DefaultTransform: transform.FromCamel().Transform(nullIfEmpty),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"ninjaone_activity":               tableNinjaoneActivity(ctx),
			"ninjaone_alert":                  tableNinjaoneAlert(ctx),
			"ninjaone_antivirus_status":       tableNinjaoneAntivirusStatus(ctx),
			"ninjaone_antivirus_threat":       tableNinjaoneAntivirusThreat(ctx),
			"ninjaone_computer_system":        tableNinjaoneComputerSystem(ctx),
			"ninjaone_device":                 tableNinjaoneDevice(ctx),
			"ninjaone_device_health":          tableNinjaoneDeviceHealth(ctx),
			"ninjaone_disk":                   tableNinjaoneDisk(ctx),
			"ninjaone_location":               tableNinjaoneLocation(ctx),
			"ninjaone_logged_on_user":         tableNinjaoneLoggedOnUser(ctx),
			"ninjaone_network_interface":      tableNinjaoneNetworkInterface(ctx),
			"ninjaone_operating_system":       tableNinjaoneOperatingSystem(ctx),
			"ninjaone_organization":           tableNinjaoneOrganization(ctx),
			"ninjaone_os_patch":               tableNinjaoneOsPatch(ctx),
			"ninjaone_os_patch_install":       tableNinjaoneOsPatchInstall(ctx),
			"ninjaone_policy":                 tableNinjaonePolicy(ctx),
			"ninjaone_policy_override":        tableNinjaonePolicyOverride(ctx),
			"ninjaone_processor":              tableNinjaoneProcessor(ctx),
			"ninjaone_role":                   tableNinjaoneRole(ctx),
			"ninjaone_software":               tableNinjaoneSoftware(ctx),
			"ninjaone_software_patch":         tableNinjaoneSoftwarePatch(ctx),
			"ninjaone_software_patch_install": tableNinjaoneSoftwarePatchInstall(ctx),
			"ninjaone_user":                   tableNinjaoneUser(ctx),
			"ninjaone_volume":                 tableNinjaoneVolume(ctx),
			"ninjaone_windows_service":        tableNinjaoneWindowsService(ctx),
		},
	}
}
