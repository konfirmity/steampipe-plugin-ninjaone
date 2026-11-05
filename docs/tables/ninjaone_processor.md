# Table: ninjaone_processor

Processors on each device.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List the processors on each device.

```sql+postgres
select
  device_id,
  name,
  architecture,
  num_cores,
  num_logical_cores
from
  ninjaone_processor;
```

```sql+sqlite
select
  device_id,
  name,
  architecture,
  num_cores,
  num_logical_cores
from
  ninjaone_processor;
```

### Count devices by processor
Plan hardware refresh, since old processors can fall out of vendor support.

```sql+postgres
select
  name,
  count(distinct device_id) as devices
from
  ninjaone_processor
group by
  name
order by
  devices desc;
```

```sql+sqlite
select
  name,
  count(distinct device_id) as devices
from
  ninjaone_processor
group by
  name
order by
  devices desc;
```

### Count devices by architecture
Check which builds of security tools you need to deploy.

```sql+postgres
select
  architecture,
  count(distinct device_id) as devices
from
  ninjaone_processor
group by
  architecture
order by
  devices desc;
```

```sql+sqlite
select
  architecture,
  count(distinct device_id) as devices
from
  ninjaone_processor
group by
  architecture
order by
  devices desc;
```

### Find devices with fewer than 4 cores
These devices may struggle to run endpoint protection.

```sql+postgres
select
  d.title,
  p.name,
  p.num_cores
from
  ninjaone_processor as p
  join ninjaone_device as d on d.id = p.device_id
where
  p.num_cores < 4;
```

```sql+sqlite
select
  d.title,
  p.name,
  p.num_cores
from
  ninjaone_processor as p
  join ninjaone_device as d on d.id = p.device_id
where
  p.num_cores < 4;
```
