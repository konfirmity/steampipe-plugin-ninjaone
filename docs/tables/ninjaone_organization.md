# Table: ninjaone_organization

NinjaOne organizations, one per client of a managed service provider.

A client app cannot be limited to one organization in NinjaOne, so a managed service provider's credential sees every tenant. NinjaOne has no filter for this list, so with `organization_id` set the plugin fetches the rows and keeps only that organization.

On a connection without `organization_id`, this table is how you find the identifier to put in the setting.

## Examples

### Basic info
List organizations. Use `id` as the connection's `organization_id`.

```sql+postgres
select
  id,
  name,
  description,
  node_approval_mode
from
  ninjaone_organization;
```

```sql+sqlite
select
  id,
  name,
  description,
  node_approval_mode
from
  ninjaone_organization;
```

### Find organizations that approve new devices automatically
With `AUTOMATIC`, a device that enrolls is managed without anyone reviewing it.

```sql+postgres
select
  id,
  name,
  node_approval_mode
from
  ninjaone_organization
where
  node_approval_mode = 'AUTOMATIC';
```

```sql+sqlite
select
  id,
  name,
  node_approval_mode
from
  ninjaone_organization
where
  node_approval_mode = 'AUTOMATIC';
```

### Count devices in each organization
Check that the device count matches what the client expects.

```sql+postgres
select
  o.name,
  count(d.id) as devices
from
  ninjaone_organization as o
  left join ninjaone_device as d on d.organization_id = o.id
group by
  o.name
order by
  devices desc;
```

```sql+sqlite
select
  o.name,
  count(d.id) as devices
from
  ninjaone_organization as o
  left join ninjaone_device as d on d.organization_id = o.id
group by
  o.name
order by
  devices desc;
```

### Show which policy each device role gets
Review the organization's role to policy assignments.

```sql+postgres
select
  o.name as organization,
  r.name as role,
  p.name as policy
from
  ninjaone_organization as o,
  jsonb_array_elements(o.policies) as a
  left join ninjaone_role as r on r.id = (a ->> 'nodeRoleId')::bigint
  left join ninjaone_policy as p on p.id = (a ->> 'policyId')::bigint
order by
  o.name,
  r.name;
```

```sql+sqlite
select
  o.name as organization,
  r.name as role,
  p.name as policy
from
  ninjaone_organization as o
  join json_each(o.policies) as a
  left join ninjaone_role as r on r.id = json_extract(a.value, '$.nodeRoleId')
  left join ninjaone_policy as p on p.id = json_extract(a.value, '$.policyId')
order by
  o.name,
  r.name;
```
