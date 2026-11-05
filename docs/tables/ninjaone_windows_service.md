# Table: ninjaone_windows_service

Windows services on each device. Devices that are not running Windows have no rows.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List Windows services on each device.

```sql+postgres
select
  device_id,
  name,
  display_name,
  state,
  start_type,
  user_name
from
  ninjaone_windows_service;
```

```sql+sqlite
select
  device_id,
  name,
  display_name,
  state,
  start_type,
  user_name
from
  ninjaone_windows_service;
```

### Check core security services
Show the state of Windows Update, Microsoft Defender and the firewall on each device.

```sql+postgres
select
  device_id,
  name,
  state,
  start_type
from
  ninjaone_windows_service
where
  name in ('wuauserv', 'WinDefend', 'mpssvc')
order by
  device_id,
  name;
```

```sql+sqlite
select
  device_id,
  name,
  state,
  start_type
from
  ninjaone_windows_service
where
  name in ('wuauserv', 'WinDefend', 'mpssvc')
order by
  device_id,
  name;
```

### Find services running under a named account
A service account's password lives on the device, so these accounts are worth reviewing.

```sql+postgres
select
  device_id,
  name,
  display_name,
  user_name
from
  ninjaone_windows_service
where
  lower(user_name) not in ('localsystem', 'nt authority\localservice', 'nt authority\networkservice');
```

```sql+sqlite
select
  device_id,
  name,
  display_name,
  user_name
from
  ninjaone_windows_service
where
  lower(user_name) not in ('localsystem', 'nt authority\localservice', 'nt authority\networkservice');
```

### Find remote access services
Compare against your list of approved tools.

```sql+postgres
select
  device_id,
  name,
  display_name,
  state
from
  ninjaone_windows_service
where
  lower(display_name) like '%teamviewer%'
  or lower(display_name) like '%anydesk%'
  or lower(display_name) like '%vnc%';
```

```sql+sqlite
select
  device_id,
  name,
  display_name,
  state
from
  ninjaone_windows_service
where
  lower(display_name) like '%teamviewer%'
  or lower(display_name) like '%anydesk%'
  or lower(display_name) like '%vnc%';
```
