# Table: ninjaone_volume

Storage volumes on each device, with BitLocker encryption state where Windows reports it.

`bit_locker_status` is a JSON object holding the conversion, encryption, protection and lock state. Read it with `->>`. A volume without BitLocker data holds a JSON null, so `bit_locker_status is null` does not find it.

BitLocker is a Windows feature. This table does not report FileVault, so no BitLocker data on a Mac says nothing about whether its disk is encrypted.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List the storage volumes on each device.

```sql+postgres
select
  device_id,
  name,
  drive_letter,
  file_system,
  capacity,
  free_space
from
  ninjaone_volume;
```

```sql+sqlite
select
  device_id,
  name,
  drive_letter,
  file_system,
  capacity,
  free_space
from
  ninjaone_volume;
```

### Show the BitLocker state of each volume
Collect evidence for a disk encryption control.

```sql+postgres
select
  device_id,
  name,
  drive_letter,
  bit_locker_status ->> 'protectionStatus' as protection_status,
  bit_locker_status ->> 'conversionStatus' as conversion_status
from
  ninjaone_volume
where
  bit_locker_status ->> 'protectionStatus' is not null;
```

```sql+sqlite
select
  device_id,
  name,
  drive_letter,
  json_extract(bit_locker_status, '$.protectionStatus') as protection_status,
  json_extract(bit_locker_status, '$.conversionStatus') as conversion_status
from
  ninjaone_volume
where
  json_extract(bit_locker_status, '$.protectionStatus') is not null;
```

### Count volumes by BitLocker protection status
Summarize encryption coverage.

```sql+postgres
select
  bit_locker_status ->> 'protectionStatus' as protection_status,
  count(*) as volumes
from
  ninjaone_volume
group by
  bit_locker_status ->> 'protectionStatus'
order by
  volumes desc;
```

```sql+sqlite
select
  json_extract(bit_locker_status, '$.protectionStatus') as protection_status,
  count(*) as volumes
from
  ninjaone_volume
group by
  json_extract(bit_locker_status, '$.protectionStatus')
order by
  volumes desc;
```

### Find Windows volumes with no BitLocker data
Check these by hand: the volume may be unencrypted, or may not be a fixed disk.

```sql+postgres
select
  d.title,
  v.name,
  v.drive_letter,
  v.device_type
from
  ninjaone_volume as v
  join ninjaone_device as d on d.id = v.device_id
where
  d.node_class like 'WINDOWS%'
  and v.bit_locker_status ->> 'protectionStatus' is null;
```

```sql+sqlite
select
  d.title,
  v.name,
  v.drive_letter,
  v.device_type
from
  ninjaone_volume as v
  join ninjaone_device as d on d.id = v.device_id
where
  d.node_class like 'WINDOWS%'
  and json_extract(v.bit_locker_status, '$.protectionStatus') is null;
```

### Find volumes with less than 10 percent free
A full disk can stop patches and logs from being written.

```sql+postgres
select
  device_id,
  name,
  capacity,
  free_space
from
  ninjaone_volume
where
  free_space * 1.0 / nullif(capacity, 0) < 0.1
order by
  free_space;
```

```sql+sqlite
select
  device_id,
  name,
  capacity,
  free_space
from
  ninjaone_volume
where
  free_space * 1.0 / nullif(capacity, 0) < 0.1
order by
  free_space;
```
