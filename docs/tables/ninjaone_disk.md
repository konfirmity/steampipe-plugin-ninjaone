# Table: ninjaone_disk

Physical disks on each device.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List the physical disks on each device.

```sql+postgres
select
  device_id,
  name,
  model,
  media_type,
  size,
  status
from
  ninjaone_disk;
```

```sql+sqlite
select
  device_id,
  name,
  model,
  media_type,
  size,
  status
from
  ninjaone_disk;
```

### List disk serial numbers for one device
Record which disks left with a device for a disposal or sanitization record.

```sql+postgres
select
  name,
  manufacturer,
  model,
  serial_number,
  size
from
  ninjaone_disk
where
  device_id = 101;
```

```sql+sqlite
select
  name,
  manufacturer,
  model,
  serial_number,
  size
from
  ninjaone_disk
where
  device_id = 101;
```

### Count disks by status
Spot disks reporting something other than the usual status.

```sql+postgres
select
  status,
  count(*) as disks
from
  ninjaone_disk
group by
  status
order by
  disks desc;
```

```sql+sqlite
select
  status,
  count(*) as disks
from
  ninjaone_disk
group by
  status
order by
  disks desc;
```

### Find disks that do not support SMART
These disks give no early warning of failure.

```sql+postgres
select
  device_id,
  name,
  model,
  interface_type
from
  ninjaone_disk
where
  smart_capable = false;
```

```sql+sqlite
select
  device_id,
  name,
  model,
  interface_type
from
  ninjaone_disk
where
  smart_capable = false;
```
