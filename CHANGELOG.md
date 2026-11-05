# Changelog

## v0.0.1 [unreleased]

_What's new?_

- New tables:
  - Devices and hardware: [ninjaone_device](docs/tables/ninjaone_device.md), [ninjaone_device_health](docs/tables/ninjaone_device_health.md), [ninjaone_computer_system](docs/tables/ninjaone_computer_system.md), [ninjaone_operating_system](docs/tables/ninjaone_operating_system.md), [ninjaone_processor](docs/tables/ninjaone_processor.md), [ninjaone_disk](docs/tables/ninjaone_disk.md), [ninjaone_volume](docs/tables/ninjaone_volume.md), [ninjaone_network_interface](docs/tables/ninjaone_network_interface.md), [ninjaone_logged_on_user](docs/tables/ninjaone_logged_on_user.md), [ninjaone_windows_service](docs/tables/ninjaone_windows_service.md)
  - Patching: [ninjaone_os_patch](docs/tables/ninjaone_os_patch.md), [ninjaone_os_patch_install](docs/tables/ninjaone_os_patch_install.md), [ninjaone_software_patch](docs/tables/ninjaone_software_patch.md), [ninjaone_software_patch_install](docs/tables/ninjaone_software_patch_install.md)
  - Software and antivirus: [ninjaone_software](docs/tables/ninjaone_software.md), [ninjaone_antivirus_status](docs/tables/ninjaone_antivirus_status.md), [ninjaone_antivirus_threat](docs/tables/ninjaone_antivirus_threat.md)
  - Alerts and activity: [ninjaone_alert](docs/tables/ninjaone_alert.md), [ninjaone_activity](docs/tables/ninjaone_activity.md)
  - Organizations, policies and users: [ninjaone_organization](docs/tables/ninjaone_organization.md), [ninjaone_location](docs/tables/ninjaone_location.md), [ninjaone_policy](docs/tables/ninjaone_policy.md), [ninjaone_policy_override](docs/tables/ninjaone_policy_override.md), [ninjaone_role](docs/tables/ninjaone_role.md), [ninjaone_user](docs/tables/ninjaone_user.md)
- Client credentials sign-in (`client_id`, `client_secret`) with the `monitoring` scope, and regional `base_url`, each with an environment-variable fallback.
- `organization_id` limits every table to one organization, because a client app cannot be limited in NinjaOne.
- `ninjaone_activity` supports `id > n` to read only new rows.
- Automatic backoff and retry on rate-limit (429) and gateway (502, 503, 504) errors, honoring `Retry-After`.
- Builds as a Steampipe plugin (`make install`), a SQLite extension (`make sqlite`) and a Postgres FDW extension (`make postgres`).
- Backup usage and custom fields are deliberately not tables yet.
