# Table: ninjaone_software_patch_install

Third-party software patch installation attempts, successful or failed. The API's `title` field is the `patch_title` column here.

No rows for a device does not prove it needs nothing: software patching may be off in its policy. See the note on `ninjaone_software_patch`.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List third-party software patch installation attempts.

```sql+postgres
select
  device_id,
  id,
  patch_title,
  impact,
  status,
  installed_at
from
  ninjaone_software_patch_install;
```

```sql+sqlite
select
  device_id,
  id,
  patch_title,
  impact,
  status,
  installed_at
from
  ninjaone_software_patch_install;
```

### Find failed installs
Show patches that need another attempt.

```sql+postgres
select
  device_id,
  patch_title,
  impact,
  timestamp
from
  ninjaone_software_patch_install
where
  status = 'FAILED'
order by
  timestamp desc;
```

```sql+sqlite
select
  device_id,
  patch_title,
  impact,
  timestamp
from
  ninjaone_software_patch_install
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
  patch_title,
  impact,
  installed_at
from
  ninjaone_software_patch_install
where
  installed_at > now() - interval '30 days'
order by
  installed_at desc;
```

```sql+sqlite
select
  device_id,
  patch_title,
  impact,
  installed_at
from
  ninjaone_software_patch_install
where
  installed_at > datetime('now', '-30 days')
order by
  installed_at desc;
```

### Count installs by status
See how often third-party patching succeeds.

```sql+postgres
select
  status,
  count(*) as installs
from
  ninjaone_software_patch_install
group by
  status
order by
  installs desc;
```

```sql+sqlite
select
  status,
  count(*) as installs
from
  ninjaone_software_patch_install
group by
  status
order by
  installs desc;
```
