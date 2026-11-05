# Table: ninjaone_role

NinjaOne device roles.

NinjaOne has no filter for this list. With `organization_id` set, the plugin returns only roles the organization's devices use. A parent role is returned only if it is itself in use, so `node_role_parent_id` can point at a row that is not in the table.

`chassis_type` here is the chassis a role applies to. It is not a reliable way to tell laptops from desktops: see `ninjaone_computer_system`.

## Examples

### Basic info
List device roles.

```sql+postgres
select
  id,
  name,
  node_class,
  custom,
  created
from
  ninjaone_role;
```

```sql+sqlite
select
  id,
  name,
  node_class,
  custom,
  created
from
  ninjaone_role;
```

### List custom roles
Review roles created by the tenant.

```sql+postgres
select
  id,
  name,
  node_class,
  created
from
  ninjaone_role
where
  custom = true
order by
  created desc;
```

```sql+sqlite
select
  id,
  name,
  node_class,
  created
from
  ninjaone_role
where
  custom = true
order by
  created desc;
```

### Count devices in each role
Check that devices are classified as you expect.

```sql+postgres
select
  r.name,
  r.node_class,
  count(d.id) as devices
from
  ninjaone_role as r
  left join ninjaone_device as d on d.node_role_id = r.id
group by
  r.name,
  r.node_class
order by
  devices desc;
```

```sql+sqlite
select
  r.name,
  r.node_class,
  count(d.id) as devices
from
  ninjaone_role as r
  left join ninjaone_device as d on d.node_role_id = r.id
group by
  r.name,
  r.node_class
order by
  devices desc;
```

### Find roles with no devices
Find unused roles. With `organization_id` set, this shows roles assigned to the organization that no device uses.

```sql+postgres
select
  r.id,
  r.name,
  r.node_class
from
  ninjaone_role as r
  left join ninjaone_device as d on d.node_role_id = r.id
where
  d.id is null;
```

```sql+sqlite
select
  r.id,
  r.name,
  r.node_class
from
  ninjaone_role as r
  left join ninjaone_device as d on d.node_role_id = r.id
where
  d.id is null;
```
