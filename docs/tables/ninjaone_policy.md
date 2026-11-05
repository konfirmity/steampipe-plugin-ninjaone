# Table: ninjaone_policy

NinjaOne device policies. This table holds each policy's name and place in the hierarchy, not its settings, so it cannot show whether patching is on. Check that in the NinjaOne console.

NinjaOne has no filter for this list. With `organization_id` set, the plugin returns only policies the organization's devices use: those that come from a device's role, and those set directly on a device. A parent policy is returned only if it is itself in use, so `parent_policy_id` can point at a row that is not in the table.

## Examples

### Basic info
List device policies.

```sql+postgres
select
  id,
  name,
  node_class,
  node_class_default,
  updated
from
  ninjaone_policy;
```

```sql+sqlite
select
  id,
  name,
  node_class,
  node_class_default,
  updated
from
  ninjaone_policy;
```

### Find policies changed in the last 30 days
Review recent changes to device configuration.

```sql+postgres
select
  id,
  name,
  node_class,
  updated
from
  ninjaone_policy
where
  updated > now() - interval '30 days'
order by
  updated desc;
```

```sql+sqlite
select
  id,
  name,
  node_class,
  updated
from
  ninjaone_policy
where
  updated > datetime('now', '-30 days')
order by
  updated desc;
```

### Count devices on each policy
`policy_id` is a policy set on the device directly and `role_policy_id` comes from its role, so a direct policy is taken first.

```sql+postgres
select
  p.name,
  p.node_class,
  count(d.id) as devices
from
  ninjaone_policy as p
  left join ninjaone_device as d on coalesce(d.policy_id, d.role_policy_id) = p.id
group by
  p.name,
  p.node_class
order by
  devices desc;
```

```sql+sqlite
select
  p.name,
  p.node_class,
  count(d.id) as devices
from
  ninjaone_policy as p
  left join ninjaone_device as d on coalesce(d.policy_id, d.role_policy_id) = p.id
group by
  p.name,
  p.node_class
order by
  devices desc;
```

### List devices with a policy set directly
These devices do not follow their role's policy.

```sql+postgres
select
  d.id,
  d.title,
  p.name as policy
from
  ninjaone_device as d
  join ninjaone_policy as p on p.id = d.policy_id
where
  d.policy_id is distinct from d.role_policy_id;
```

```sql+sqlite
select
  d.id,
  d.title,
  p.name as policy
from
  ninjaone_device as d
  join ninjaone_policy as p on p.id = d.policy_id
where
  d.policy_id is distinct from d.role_policy_id;
```

### Show each policy's parent
Trace inheritance. On a connection with `organization_id`, a null parent name can mean the parent is not in use by that organization.

```sql+postgres
select
  p.name,
  p.parent_policy_id,
  parent.name as parent_name
from
  ninjaone_policy as p
  left join ninjaone_policy as parent on parent.id = p.parent_policy_id
where
  p.parent_policy_id is not null;
```

```sql+sqlite
select
  p.name,
  p.parent_policy_id,
  parent.name as parent_name
from
  ninjaone_policy as p
  left join ninjaone_policy as parent on parent.id = p.parent_policy_id
where
  p.parent_policy_id is not null;
```
