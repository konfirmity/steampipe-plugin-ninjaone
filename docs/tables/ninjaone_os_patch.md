# Table: ninjaone_os_patch

Operating system patches that are pending, failed or rejected. Installed patches and install attempts are in `ninjaone_os_patch_install`.

An empty result here, or a `HEALTHY` status in `ninjaone_device_health`, does not prove a device is patched: patching may be off in its policy. Software patching is off by default in the Mac policy we inspected, and Mac OS update settings sit under MDM. This plugin does not expose policy contents, so confirm the setting in the NinjaOne console.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List outstanding operating system patches.

```sql+postgres
select
  device_id,
  id,
  name,
  severity,
  status,
  kb_number
from
  ninjaone_os_patch;
```

```sql+sqlite
select
  device_id,
  id,
  name,
  severity,
  status,
  kb_number
from
  ninjaone_os_patch;
```

### Find outstanding critical patches
Rank what to fix first.

```sql+postgres
select
  d.title,
  p.name,
  p.kb_number,
  p.status,
  p.timestamp
from
  ninjaone_os_patch as p
  join ninjaone_device as d on d.id = p.device_id
where
  p.severity = 'CRITICAL'
order by
  p.timestamp;
```

```sql+sqlite
select
  d.title,
  p.name,
  p.kb_number,
  p.status,
  p.timestamp
from
  ninjaone_os_patch as p
  join ninjaone_device as d on d.id = p.device_id
where
  p.severity = 'CRITICAL'
order by
  p.timestamp;
```

### Count patches by status and severity
Summarize the patch backlog.

```sql+postgres
select
  status,
  severity,
  count(*) as patches
from
  ninjaone_os_patch
group by
  status,
  severity
order by
  patches desc;
```

```sql+sqlite
select
  status,
  severity,
  count(*) as patches
from
  ninjaone_os_patch
group by
  status,
  severity
order by
  patches desc;
```

### Find failed and rejected patches
A rejected patch was refused on purpose, so it should have a recorded reason.

```sql+postgres
select
  device_id,
  name,
  kb_number,
  status
from
  ninjaone_os_patch
where
  status in ('FAILED', 'REJECTED');
```

```sql+sqlite
select
  device_id,
  name,
  kb_number,
  status
from
  ninjaone_os_patch
where
  status in ('FAILED', 'REJECTED');
```

### List devices with no rows in this report
Treat these as unknown, not as patched. Check patching is on in each device's policy.

```sql+postgres
select
  d.id,
  d.title,
  d.node_class,
  coalesce(d.policy_id, d.role_policy_id) as policy_id
from
  ninjaone_device as d
  left join ninjaone_os_patch as p on p.device_id = d.id
where
  p.device_id is null;
```

```sql+sqlite
select
  d.id,
  d.title,
  d.node_class,
  coalesce(d.policy_id, d.role_policy_id) as policy_id
from
  ninjaone_device as d
  left join ninjaone_os_patch as p on p.device_id = d.id
where
  p.device_id is null;
```
