# Table: ninjaone_computer_system

Hardware identity of each device: manufacturer, model, serial number and domain.

Do not use `chassis_type` to tell laptops from desktops: it reported `DESKTOP` for a MacBook Pro. `bios_serial_number` was empty on a Mac, so use `serial_number` to identify hardware.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List the hardware identity of each device.

```sql+postgres
select
  device_id,
  name,
  manufacturer,
  model,
  serial_number,
  virtual_machine
from
  ninjaone_computer_system;
```

```sql+sqlite
select
  device_id,
  name,
  manufacturer,
  model,
  serial_number,
  virtual_machine
from
  ninjaone_computer_system;
```

### Find devices with no serial number
Asset records need a serial number to tie a device to a purchase or disposal record.

```sql+postgres
select
  device_id,
  name,
  manufacturer,
  model
from
  ninjaone_computer_system
where
  serial_number is null;
```

```sql+sqlite
select
  device_id,
  name,
  manufacturer,
  model
from
  ninjaone_computer_system
where
  serial_number is null;
```

### Count devices by manufacturer and model
Summarize the hardware estate.

```sql+postgres
select
  manufacturer,
  model,
  count(*) as devices
from
  ninjaone_computer_system
group by
  manufacturer,
  model
order by
  devices desc;
```

```sql+sqlite
select
  manufacturer,
  model,
  count(*) as devices
from
  ninjaone_computer_system
group by
  manufacturer,
  model
order by
  devices desc;
```

### List virtual machines
Separate virtual machines from physical assets.

```sql+postgres
select
  device_id,
  name,
  manufacturer,
  model
from
  ninjaone_computer_system
where
  virtual_machine = true;
```

```sql+sqlite
select
  device_id,
  name,
  manufacturer,
  model
from
  ninjaone_computer_system
where
  virtual_machine = true;
```

### Find devices not joined to a domain
Identify devices outside central account management.

```sql+postgres
select
  device_id,
  name,
  model,
  serial_number
from
  ninjaone_computer_system
where
  domain is null;
```

```sql+sqlite
select
  device_id,
  name,
  model,
  serial_number
from
  ninjaone_computer_system
where
  domain is null;
```
