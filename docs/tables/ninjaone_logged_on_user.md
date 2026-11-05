# Table: ninjaone_logged_on_user

Users logged on to each device. These are local login names, not email addresses, so they do not join to `ninjaone_user.email`.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List who is logged on to each device.

```sql+postgres
select
  device_id,
  user_name,
  logon_time
from
  ninjaone_logged_on_user;
```

```sql+sqlite
select
  device_id,
  user_name,
  logon_time
from
  ninjaone_logged_on_user;
```

### Find logons with a built-in admin account
Shared and built-in accounts make actions hard to attribute.

```sql+postgres
select
  d.title,
  u.user_name,
  u.logon_time
from
  ninjaone_logged_on_user as u
  join ninjaone_device as d on d.id = u.device_id
where
  lower(u.user_name) in ('administrator', 'admin', 'root');
```

```sql+sqlite
select
  d.title,
  u.user_name,
  u.logon_time
from
  ninjaone_logged_on_user as u
  join ninjaone_device as d on d.id = u.device_id
where
  lower(u.user_name) in ('administrator', 'admin', 'root');
```

### Find devices with more than one user logged on
Spot shared devices.

```sql+postgres
select
  device_id,
  count(distinct user_name) as users
from
  ninjaone_logged_on_user
group by
  device_id
having
  count(distinct user_name) > 1
order by
  users desc;
```

```sql+sqlite
select
  device_id,
  count(distinct user_name) as users
from
  ninjaone_logged_on_user
group by
  device_id
having
  count(distinct user_name) > 1
order by
  users desc;
```

### Find sessions older than 7 days
Long sessions can outlive a screen lock or an access review.

```sql+postgres
select
  device_id,
  user_name,
  logon_time
from
  ninjaone_logged_on_user
where
  logon_time < now() - interval '7 days'
order by
  logon_time;
```

```sql+sqlite
select
  device_id,
  user_name,
  logon_time
from
  ninjaone_logged_on_user
where
  logon_time < datetime('now', '-7 days')
order by
  logon_time;
```
