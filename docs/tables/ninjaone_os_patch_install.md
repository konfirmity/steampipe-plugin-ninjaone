# Table: ninjaone_os_patch_install

Operating system patch installation attempts, successful or failed.

No rows for a device does not prove it needs nothing: patching may be off in its policy. See the note on `ninjaone_os_patch`.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List operating system patch installation attempts.

```sql+postgres
select
  device_id,
  id,
  name,
  severity,
  status,
  installed_at
from
  ninjaone_os_patch_install;
```

```sql+sqlite
select
  device_id,
  id,
  name,
  severity,
  status,
  installed_at
from
  ninjaone_os_patch_install;
```

### Find failed installs
Show patches that need another attempt.

```sql+postgres
select
  device_id,
  name,
  kb_number,
  severity,
  timestamp
from
  ninjaone_os_patch_install
where
  status = 'FAILED'
order by
  timestamp desc;
```

```sql+sqlite
select
  device_id,
  name,
  kb_number,
  severity,
  timestamp
from
  ninjaone_os_patch_install
where
  status = 'FAILED'
order by
  timestamp desc;
```

### List patches installed in the last 30 days
Collect evidence for a patching control.

```sql+postgres
select
  device_id,
  name,
  kb_number,
  severity,
  installed_at
from
  ninjaone_os_patch_install
where
  installed_at > now() - interval '30 days'
order by
  installed_at desc;
```

```sql+sqlite
select
  device_id,
  name,
  kb_number,
  severity,
  installed_at
from
  ninjaone_os_patch_install
where
  installed_at > datetime('now', '-30 days')
order by
  installed_at desc;
```

### Find devices with no successful install in 90 days
These devices may have patching off, or may be stuck.

```sql+postgres
select
  d.id,
  d.title,
  max(i.installed_at) as last_install
from
  ninjaone_device as d
  left join ninjaone_os_patch_install as i on i.device_id = d.id and i.status = 'INSTALLED'
group by
  d.id,
  d.title
having
  max(i.installed_at) is null
  or max(i.installed_at) < now() - interval '90 days';
```

```sql+sqlite
select
  d.id,
  d.title,
  max(i.installed_at) as last_install
from
  ninjaone_device as d
  left join ninjaone_os_patch_install as i on i.device_id = d.id and i.status = 'INSTALLED'
group by
  d.id,
  d.title
having
  max(i.installed_at) is null
  or max(i.installed_at) < datetime('now', '-90 days');
```
