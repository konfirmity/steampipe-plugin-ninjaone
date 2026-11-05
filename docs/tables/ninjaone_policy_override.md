# Table: ninjaone_policy_override

Policy settings overridden on individual devices. `overrides` is a JSON array of the overridden section names, not the values they were set to.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List devices with policy overrides.

```sql+postgres
select
  device_id,
  overrides
from
  ninjaone_policy_override;
```

```sql+sqlite
select
  device_id,
  overrides
from
  ninjaone_policy_override;
```

### Count devices by overridden sections
See which sections are overridden most often.

```sql+postgres
select
  overrides,
  count(*) as devices
from
  ninjaone_policy_override
group by
  overrides
order by
  devices desc;
```

```sql+sqlite
select
  overrides,
  count(*) as devices
from
  ninjaone_policy_override
group by
  overrides
order by
  devices desc;
```

### Find overrides that touch patching
An override may be why a device differs from its policy. Check the device in the NinjaOne console for the value.

```sql+postgres
select
  device_id,
  overrides
from
  ninjaone_policy_override
where
  lower(overrides::text) like '%patch%';
```

```sql+sqlite
select
  device_id,
  overrides
from
  ninjaone_policy_override
where
  lower(overrides) like '%patch%';
```

### Show overrides with the device and its policy
Review exceptions to a baseline.

```sql+postgres
select
  d.title,
  p.name as policy,
  o.overrides
from
  ninjaone_policy_override as o
  join ninjaone_device as d on d.id = o.device_id
  left join ninjaone_policy as p on p.id = coalesce(d.policy_id, d.role_policy_id);
```

```sql+sqlite
select
  d.title,
  p.name as policy,
  o.overrides
from
  ninjaone_policy_override as o
  join ninjaone_device as d on d.id = o.device_id
  left join ninjaone_policy as p on p.id = coalesce(d.policy_id, d.role_policy_id);
```
