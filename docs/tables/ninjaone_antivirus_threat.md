# Table: ninjaone_antivirus_threat

Threats detected by antivirus on each device.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List detected threats.

```sql+postgres
select
  device_id,
  name,
  category,
  level,
  status,
  timestamp
from
  ninjaone_antivirus_threat;
```

```sql+sqlite
select
  device_id,
  name,
  category,
  level,
  status,
  timestamp
from
  ninjaone_antivirus_threat;
```

### Find threats reported in the last 30 days
Review recent detections.

```sql+postgres
select
  device_id,
  name,
  level,
  status,
  detection_source,
  timestamp
from
  ninjaone_antivirus_threat
where
  timestamp > now() - interval '30 days'
order by
  timestamp desc;
```

```sql+sqlite
select
  device_id,
  name,
  level,
  status,
  detection_source,
  timestamp
from
  ninjaone_antivirus_threat
where
  timestamp > datetime('now', '-30 days')
order by
  timestamp desc;
```

### Count threats by level and status
Summarize detections and how they were handled.

```sql+postgres
select
  level,
  status,
  count(*) as threats
from
  ninjaone_antivirus_threat
group by
  level,
  status
order by
  threats desc;
```

```sql+sqlite
select
  level,
  status,
  count(*) as threats
from
  ninjaone_antivirus_threat
group by
  level,
  status
order by
  threats desc;
```

### Find threats seen on more than one device
Spot something spreading across the estate.

```sql+postgres
select
  name,
  count(distinct device_id) as devices
from
  ninjaone_antivirus_threat
group by
  name
having
  count(distinct device_id) > 1
order by
  devices desc;
```

```sql+sqlite
select
  name,
  count(distinct device_id) as devices
from
  ninjaone_antivirus_threat
group by
  name
having
  count(distinct device_id) > 1
order by
  devices desc;
```

### Show threats with the device name
Give a responder the device name and the files involved.

```sql+postgres
select
  d.title as device,
  t.name,
  t.status,
  t.trace_list,
  t.timestamp
from
  ninjaone_antivirus_threat as t
  join ninjaone_device as d on d.id = t.device_id
order by
  t.timestamp desc;
```

```sql+sqlite
select
  d.title as device,
  t.name,
  t.status,
  t.trace_list,
  t.timestamp
from
  ninjaone_antivirus_threat as t
  join ninjaone_device as d on d.id = t.device_id
order by
  t.timestamp desc;
```
