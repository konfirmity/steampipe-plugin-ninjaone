# Table: ninjaone_activity

The NinjaOne activity log: patching, scripts, logins and other events, newest first.

Filter on `id > n` to read only rows newer than the last one you saw. Only the `>` operator on `id` is sent to the API. Any other comparison on `id`, or a filter on `activity_time`, reads the log and filters afterwards.

With `organization_id` set, the device filter (`df`) returns only activity on that organization's devices. Activity that is not tied to a device is left out.

## Examples

### Basic info
List activity with its device and status.

```sql+postgres
select
  id,
  activity_time,
  device_id,
  activity_type,
  status,
  message
from
  ninjaone_activity;
```

```sql+sqlite
select
  id,
  activity_time,
  device_id,
  activity_type,
  status,
  message
from
  ninjaone_activity;
```

### Read only new activity
Store the highest `id` you have processed and pass it on the next run. The plugin sends this to NinjaOne, so older rows are not fetched.

```sql+postgres
select
  id,
  activity_time,
  device_id,
  activity_type,
  status,
  message
from
  ninjaone_activity
where
  id > 12345
order by
  id;
```

```sql+sqlite
select
  id,
  activity_time,
  device_id,
  activity_type,
  status,
  message
from
  ninjaone_activity
where
  id > 12345
order by
  id;
```

### List recent activity on one device
Build a timeline for a device during an incident review.

```sql+postgres
select
  id,
  activity_time,
  activity_type,
  status,
  message
from
  ninjaone_activity
where
  device_id = 101
order by
  activity_time desc
limit 50;
```

```sql+sqlite
select
  id,
  activity_time,
  activity_type,
  status,
  message
from
  ninjaone_activity
where
  device_id = 101
order by
  activity_time desc
limit 50;
```

### Show patching activity from the last 30 days
Collect evidence that patch management ran.

```sql+postgres
select
  activity_time,
  device_id,
  status,
  message
from
  ninjaone_activity
where
  activity_type = 'PATCH_MANAGEMENT'
  and activity_time > now() - interval '30 days'
order by
  activity_time desc;
```

```sql+sqlite
select
  activity_time,
  device_id,
  status,
  message
from
  ninjaone_activity
where
  activity_type = 'PATCH_MANAGEMENT'
  and activity_time > datetime('now', '-30 days')
order by
  activity_time desc;
```

### Attribute activity to the user who triggered it
Review who acted on devices. With `organization_id` set, the provider's technicians are not in `ninjaone_user`, so their rows show a `user_id` and no email.

```sql+postgres
select
  a.activity_time,
  a.device_id,
  a.user_id,
  u.email,
  a.message
from
  ninjaone_activity as a
  left join ninjaone_user as u on u.id = a.user_id
where
  a.user_id is not null
order by
  a.activity_time desc
limit 100;
```

```sql+sqlite
select
  a.activity_time,
  a.device_id,
  a.user_id,
  u.email,
  a.message
from
  ninjaone_activity as a
  left join ninjaone_user as u on u.id = a.user_id
where
  a.user_id is not null
order by
  a.activity_time desc
limit 100;
```
