# Table: ninjaone_device

Devices managed by NinjaOne, with their operating system, hardware, network and last check-in.

With `organization_id` set, the list is limited by the device filter (`df`). A lookup by `id` fetches that one device and returns nothing if it belongs to another organization.

Do not use `chassis_type` to tell laptops from desktops: it reported `DESKTOP` for a MacBook Pro. `bios_serial_number` was empty on a Mac, so use `serial_number` to identify hardware.

`last_logged_in_user` is a local login name, not an email address. `assigned_owner_uid` joins to `ninjaone_user.uid` when you need the person.

## Examples

### Basic info
List devices with their operating system and last check-in.

```sql+postgres
select
  id,
  title,
  node_class,
  organization_id,
  os_name,
  serial_number,
  last_contact
from
  ninjaone_device;
```

```sql+sqlite
select
  id,
  title,
  node_class,
  organization_id,
  os_name,
  serial_number,
  last_contact
from
  ninjaone_device;
```

### Find devices not seen for 30 days
Find agents that stopped reporting, whose patch and antivirus data has gone stale.

```sql+postgres
select
  id,
  title,
  node_class,
  last_logged_in_user,
  last_contact
from
  ninjaone_device
where
  approval_status = 'APPROVED'
  and last_contact < now() - interval '30 days'
order by
  last_contact;
```

```sql+sqlite
select
  id,
  title,
  node_class,
  last_logged_in_user,
  last_contact
from
  ninjaone_device
where
  approval_status = 'APPROVED'
  and last_contact < datetime('now', '-30 days')
order by
  last_contact;
```

### List devices waiting for a reboot
A pending reboot can mean an installed patch is not yet active.

```sql+postgres
select
  id,
  title,
  os_name,
  os_last_boot_time
from
  ninjaone_device
where
  os_needs_reboot = true
order by
  os_last_boot_time;
```

```sql+sqlite
select
  id,
  title,
  os_name,
  os_last_boot_time
from
  ninjaone_device
where
  os_needs_reboot = true
order by
  os_last_boot_time;
```

### List devices waiting for approval
Review devices that enrolled but have not been approved.

```sql+postgres
select
  id,
  title,
  node_class,
  public_ip,
  created
from
  ninjaone_device
where
  approval_status = 'PENDING'
order by
  created;
```

```sql+sqlite
select
  id,
  title,
  node_class,
  public_ip,
  created
from
  ninjaone_device
where
  approval_status = 'PENDING'
order by
  created;
```

### Find devices with no assigned owner
Asset registers need an accountable owner for each device.

```sql+postgres
select
  id,
  title,
  node_class,
  last_logged_in_user
from
  ninjaone_device
where
  assigned_owner_uid is null;
```

```sql+sqlite
select
  id,
  title,
  node_class,
  last_logged_in_user
from
  ninjaone_device
where
  assigned_owner_uid is null;
```
