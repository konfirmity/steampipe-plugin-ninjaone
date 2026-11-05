# Table: ninjaone_location

Locations inside NinjaOne organizations.

NinjaOne has no filter for this list, so with `organization_id` set the plugin fetches the rows and keeps only that organization's locations.

## Examples

### Basic info
List locations.

```sql+postgres
select
  id,
  name,
  organization_id,
  address,
  description
from
  ninjaone_location;
```

```sql+sqlite
select
  id,
  name,
  organization_id,
  address,
  description
from
  ninjaone_location;
```

### Count devices at each location
Show where managed devices are.

```sql+postgres
select
  l.name,
  count(d.id) as devices
from
  ninjaone_location as l
  left join ninjaone_device as d on d.location_id = l.id
group by
  l.name
order by
  devices desc;
```

```sql+sqlite
select
  l.name,
  count(d.id) as devices
from
  ninjaone_location as l
  left join ninjaone_device as d on d.location_id = l.id
group by
  l.name
order by
  devices desc;
```

### Find locations with no devices
Find locations that can be retired, or sites that were never enrolled.

```sql+postgres
select
  l.id,
  l.name,
  l.address
from
  ninjaone_location as l
  left join ninjaone_device as d on d.location_id = l.id
where
  d.id is null;
```

```sql+sqlite
select
  l.id,
  l.name,
  l.address
from
  ninjaone_location as l
  left join ninjaone_device as d on d.location_id = l.id
where
  d.id is null;
```

### Show locations with their organization
On a connection without `organization_id`, list the sites of each client.

```sql+postgres
select
  o.name as organization,
  l.name as location,
  l.address
from
  ninjaone_location as l
  join ninjaone_organization as o on o.id = l.organization_id
order by
  o.name,
  l.name;
```

```sql+sqlite
select
  o.name as organization,
  l.name as location,
  l.address
from
  ninjaone_location as l
  join ninjaone_organization as o on o.id = l.organization_id
order by
  o.name,
  l.name;
```
