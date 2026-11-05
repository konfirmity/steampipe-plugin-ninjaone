# Table: ninjaone_software_patch

Third-party software patches that are pending, failed or rejected. The API's `title` field is the `patch_title` column here, because `title` is the standard Steampipe column. Install attempts are in `ninjaone_software_patch_install`.

An empty result here, or a `HEALTHY` status in `ninjaone_device_health`, does not prove a device is patched: patching may be off in its policy. Software patching is off by default in the Mac policy we inspected, and Mac OS update settings sit under MDM. This plugin does not expose policy contents, so confirm the setting in the NinjaOne console.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List outstanding third-party software patches.

```sql+postgres
select
  device_id,
  id,
  patch_title,
  impact,
  status,
  product_identifier
from
  ninjaone_software_patch;
```

```sql+sqlite
select
  device_id,
  id,
  patch_title,
  impact,
  status,
  product_identifier
from
  ninjaone_software_patch;
```

### Find outstanding critical patches
Rank what to fix first.

```sql+postgres
select
  d.title,
  p.patch_title,
  p.status,
  p.timestamp
from
  ninjaone_software_patch as p
  join ninjaone_device as d on d.id = p.device_id
where
  p.impact = 'CRITICAL'
order by
  p.timestamp;
```

```sql+sqlite
select
  d.title,
  p.patch_title,
  p.status,
  p.timestamp
from
  ninjaone_software_patch as p
  join ninjaone_device as d on d.id = p.device_id
where
  p.impact = 'CRITICAL'
order by
  p.timestamp;
```

### Count patches by status and impact
Summarize the patch backlog.

```sql+postgres
select
  status,
  impact,
  count(*) as patches
from
  ninjaone_software_patch
group by
  status,
  impact
order by
  patches desc;
```

```sql+sqlite
select
  status,
  impact,
  count(*) as patches
from
  ninjaone_software_patch
group by
  status,
  impact
order by
  patches desc;
```

### Count devices needing each patch
Find the single patches that would clear the most devices.

```sql+postgres
select
  patch_title,
  count(distinct device_id) as devices
from
  ninjaone_software_patch
group by
  patch_title
order by
  devices desc;
```

```sql+sqlite
select
  patch_title,
  count(distinct device_id) as devices
from
  ninjaone_software_patch
group by
  patch_title
order by
  devices desc;
```

### List Mac devices with no rows in this report
Treat these as unknown, not as patched: software patching is off by default in the Mac policy we inspected.

```sql+postgres
select
  d.id,
  d.title,
  coalesce(d.policy_id, d.role_policy_id) as policy_id
from
  ninjaone_device as d
  left join ninjaone_software_patch as p on p.device_id = d.id
where
  d.node_class = 'MAC'
  and p.device_id is null;
```

```sql+sqlite
select
  d.id,
  d.title,
  coalesce(d.policy_id, d.role_policy_id) as policy_id
from
  ninjaone_device as d
  left join ninjaone_software_patch as p on p.device_id = d.id
where
  d.node_class = 'MAC'
  and p.device_id is null;
```
