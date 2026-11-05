# Table: ninjaone_device_health

Health summary per device: threat, patch, alert and vulnerability counts.

NinjaOne exposes vulnerability counts per device here, not CVE lists, so this plugin cannot tell you which CVEs affect a device.

A `HEALTHY` status or a zero patch count does not prove a device is patched: patching may be off in its policy. Software patching is off by default in the Mac policy we inspected, and Mac OS update settings sit under MDM. This plugin does not expose policy contents, so confirm the setting in the NinjaOne console.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List the health summary of each device.

```sql+postgres
select
  device_id,
  health_status,
  pending_os_patches_count,
  failed_os_patches_count,
  active_threats_count,
  critical_vulnerability_count
from
  ninjaone_device_health;
```

```sql+sqlite
select
  device_id,
  health_status,
  pending_os_patches_count,
  failed_os_patches_count,
  active_threats_count,
  critical_vulnerability_count
from
  ninjaone_device_health;
```

### Find devices with critical or high vulnerabilities
Rank devices by vulnerability count. Look in the NinjaOne console for the detail behind a count.

```sql+postgres
select
  d.title,
  h.critical_vulnerability_count,
  h.high_vulnerability_count
from
  ninjaone_device_health as h
  join ninjaone_device as d on d.id = h.device_id
where
  h.critical_vulnerability_count > 0
  or h.high_vulnerability_count > 0
order by
  h.critical_vulnerability_count desc,
  h.high_vulnerability_count desc;
```

```sql+sqlite
select
  d.title,
  h.critical_vulnerability_count,
  h.high_vulnerability_count
from
  ninjaone_device_health as h
  join ninjaone_device as d on d.id = h.device_id
where
  h.critical_vulnerability_count > 0
  or h.high_vulnerability_count > 0
order by
  h.critical_vulnerability_count desc,
  h.high_vulnerability_count desc;
```

### Find devices with active threats
List devices that need a responder now.

```sql+postgres
select
  d.title,
  h.active_threats_count,
  h.quarantined_threats_count,
  h.blocked_threats_count
from
  ninjaone_device_health as h
  join ninjaone_device as d on d.id = h.device_id
where
  h.active_threats_count > 0
order by
  h.active_threats_count desc;
```

```sql+sqlite
select
  d.title,
  h.active_threats_count,
  h.quarantined_threats_count,
  h.blocked_threats_count
from
  ninjaone_device_health as h
  join ninjaone_device as d on d.id = h.device_id
where
  h.active_threats_count > 0
order by
  h.active_threats_count desc;
```

### Find devices with failed patches
Show patching that was attempted and did not work.

```sql+postgres
select
  device_id,
  failed_os_patches_count,
  failed_software_patches_count,
  pending_reboot_reason
from
  ninjaone_device_health
where
  failed_os_patches_count > 0
  or failed_software_patches_count > 0;
```

```sql+sqlite
select
  device_id,
  failed_os_patches_count,
  failed_software_patches_count,
  pending_reboot_reason
from
  ninjaone_device_health
where
  failed_os_patches_count > 0
  or failed_software_patches_count > 0;
```

### Count devices by health status
Summarize the estate.

```sql+postgres
select
  health_status,
  count(*) as devices
from
  ninjaone_device_health
group by
  health_status
order by
  devices desc;
```

```sql+sqlite
select
  health_status,
  count(*) as devices
from
  ninjaone_device_health
group by
  health_status
order by
  devices desc;
```
