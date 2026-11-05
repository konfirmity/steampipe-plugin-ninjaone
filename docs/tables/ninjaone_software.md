# Table: ninjaone_software

Software installed on each device. `install_date` is text, as reported by the device.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List installed software.

```sql+postgres
select
  device_id,
  name,
  version,
  publisher,
  install_date
from
  ninjaone_software;
```

```sql+sqlite
select
  device_id,
  name,
  version,
  publisher,
  install_date
from
  ninjaone_software;
```

### Find remote access tools
Remote access software outside your approved list is a common way in.

```sql+postgres
select
  d.title,
  s.name,
  s.version
from
  ninjaone_software as s
  join ninjaone_device as d on d.id = s.device_id
where
  lower(s.name) like '%teamviewer%'
  or lower(s.name) like '%anydesk%'
  or lower(s.name) like '%vnc%';
```

```sql+sqlite
select
  d.title,
  s.name,
  s.version
from
  ninjaone_software as s
  join ninjaone_device as d on d.id = s.device_id
where
  lower(s.name) like '%teamviewer%'
  or lower(s.name) like '%anydesk%'
  or lower(s.name) like '%vnc%';
```

### Count installs of each product
Build a software inventory, leaving out operating system components.

```sql+postgres
select
  name,
  publisher,
  count(distinct device_id) as devices
from
  ninjaone_software
where
  is_system_component is not true
group by
  name,
  publisher
order by
  devices desc;
```

```sql+sqlite
select
  name,
  publisher,
  count(distinct device_id) as devices
from
  ninjaone_software
where
  is_system_component is not true
group by
  name,
  publisher
order by
  devices desc;
```

### List installed versions of one product
Find devices on a version with a known problem.

```sql+postgres
select
  version,
  count(*) as devices
from
  ninjaone_software
where
  lower(name) like '%chrome%'
group by
  version
order by
  devices desc;
```

```sql+sqlite
select
  version,
  count(*) as devices
from
  ninjaone_software
where
  lower(name) like '%chrome%'
group by
  version
order by
  devices desc;
```

### Find software with no publisher
Unattributed software deserves a closer look.

```sql+postgres
select
  name,
  version,
  location,
  count(*) as devices
from
  ninjaone_software
where
  publisher is null
  and is_system_component is not true
group by
  name,
  version,
  location
order by
  devices desc;
```

```sql+sqlite
select
  name,
  version,
  location,
  count(*) as devices
from
  ninjaone_software
where
  publisher is null
  and is_system_component is not true
group by
  name,
  version,
  location
order by
  devices desc;
```
