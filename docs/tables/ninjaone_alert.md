# Table: ninjaone_alert

Active alerts on NinjaOne devices. A `device_id` filter is sent to the API.

With `organization_id` set, the device filter (`df`) returns only alerts on that organization's devices.

## Examples

### Basic info
List active alerts.

```sql+postgres
select
  uid,
  device_id,
  subject,
  source_type,
  create_time
from
  ninjaone_alert;
```

```sql+sqlite
select
  uid,
  device_id,
  subject,
  source_type,
  create_time
from
  ninjaone_alert;
```

### Find alerts open for more than 7 days
Show alerts nobody has dealt with.

```sql+postgres
select
  uid,
  device_id,
  subject,
  source_name,
  create_time
from
  ninjaone_alert
where
  create_time < now() - interval '7 days'
order by
  create_time;
```

```sql+sqlite
select
  uid,
  device_id,
  subject,
  source_name,
  create_time
from
  ninjaone_alert
where
  create_time < datetime('now', '-7 days')
order by
  create_time;
```

### Count alerts by source
See which conditions raise the most alerts.

```sql+postgres
select
  source_type,
  source_name,
  count(*) as alerts
from
  ninjaone_alert
group by
  source_type,
  source_name
order by
  alerts desc;
```

```sql+sqlite
select
  source_type,
  source_name,
  count(*) as alerts
from
  ninjaone_alert
group by
  source_type,
  source_name
order by
  alerts desc;
```

### Show alerts with the device name
Give a responder the device name alongside the alert.

```sql+postgres
select
  d.title as device,
  a.subject,
  a.message,
  a.create_time
from
  ninjaone_alert as a
  join ninjaone_device as d on d.id = a.device_id
order by
  a.create_time desc;
```

```sql+sqlite
select
  d.title as device,
  a.subject,
  a.message,
  a.create_time
from
  ninjaone_alert as a
  join ninjaone_device as d on d.id = a.device_id
order by
  a.create_time desc;
```
