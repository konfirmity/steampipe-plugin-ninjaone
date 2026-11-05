# Table: ninjaone_operating_system

Operating system on each device.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List the operating system on each device.

```sql+postgres
select
  device_id,
  name,
  build_number,
  release_id,
  architecture,
  last_boot_time
from
  ninjaone_operating_system;
```

```sql+sqlite
select
  device_id,
  name,
  build_number,
  release_id,
  architecture,
  last_boot_time
from
  ninjaone_operating_system;
```

### Count devices by operating system and build
Show version spread, and how far behind the oldest builds are.

```sql+postgres
select
  name,
  build_number,
  count(*) as devices
from
  ninjaone_operating_system
group by
  name,
  build_number
order by
  devices desc;
```

```sql+sqlite
select
  name,
  build_number,
  count(*) as devices
from
  ninjaone_operating_system
group by
  name,
  build_number
order by
  devices desc;
```

### Find Windows 10 devices
Find devices that need an upgrade plan.

```sql+postgres
select
  device_id,
  name,
  build_number
from
  ninjaone_operating_system
where
  name like '%Windows 10%';
```

```sql+sqlite
select
  device_id,
  name,
  build_number
from
  ninjaone_operating_system
where
  name like '%Windows 10%';
```

### Find devices not rebooted for 30 days
Patches often need a reboot to take effect.

```sql+postgres
select
  device_id,
  name,
  last_boot_time
from
  ninjaone_operating_system
where
  last_boot_time < now() - interval '30 days'
order by
  last_boot_time;
```

```sql+sqlite
select
  device_id,
  name,
  last_boot_time
from
  ninjaone_operating_system
where
  last_boot_time < datetime('now', '-30 days')
order by
  last_boot_time;
```

### List devices waiting for a reboot
A pending reboot can mean an installed patch is not yet active.

```sql+postgres
select
  device_id,
  name,
  last_boot_time
from
  ninjaone_operating_system
where
  needs_reboot = true;
```

```sql+sqlite
select
  device_id,
  name,
  last_boot_time
from
  ninjaone_operating_system
where
  needs_reboot = true;
```
