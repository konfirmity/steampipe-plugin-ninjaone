# Table: ninjaone_antivirus_status

Antivirus product and state on each device. A `product_name` of `NONE` means no antivirus was detected.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List the antivirus product on each device.

```sql+postgres
select
  device_id,
  product_name,
  product_state,
  definition_status,
  version,
  timestamp
from
  ninjaone_antivirus_status;
```

```sql+sqlite
select
  device_id,
  product_name,
  product_state,
  definition_status,
  version,
  timestamp
from
  ninjaone_antivirus_status;
```

### Find devices with no antivirus detected
List devices that need antivirus deployed.

```sql+postgres
select
  d.id,
  d.title,
  d.node_class,
  d.last_contact
from
  ninjaone_antivirus_status as a
  join ninjaone_device as d on d.id = a.device_id
where
  a.product_name = 'NONE';
```

```sql+sqlite
select
  d.id,
  d.title,
  d.node_class,
  d.last_contact
from
  ninjaone_antivirus_status as a
  join ninjaone_device as d on d.id = a.device_id
where
  a.product_name = 'NONE';
```

### Find devices missing from the antivirus report
A device with no row has reported no antivirus state at all, which is not the same as `NONE`.

```sql+postgres
select
  d.id,
  d.title,
  d.node_class
from
  ninjaone_device as d
  left join ninjaone_antivirus_status as a on a.device_id = d.id
where
  a.device_id is null;
```

```sql+sqlite
select
  d.id,
  d.title,
  d.node_class
from
  ninjaone_device as d
  left join ninjaone_antivirus_status as a on a.device_id = d.id
where
  a.device_id is null;
```

### Count devices by product and state
Show antivirus coverage for an audit.

```sql+postgres
select
  product_name,
  product_state,
  count(*) as devices
from
  ninjaone_antivirus_status
group by
  product_name,
  product_state
order by
  devices desc;
```

```sql+sqlite
select
  product_name,
  product_state,
  count(*) as devices
from
  ninjaone_antivirus_status
group by
  product_name,
  product_state
order by
  devices desc;
```

### Count devices by definition status
Check whether definitions are being kept current.

```sql+postgres
select
  definition_status,
  count(*) as devices
from
  ninjaone_antivirus_status
group by
  definition_status
order by
  devices desc;
```

```sql+sqlite
select
  definition_status,
  count(*) as devices
from
  ninjaone_antivirus_status
group by
  definition_status
order by
  devices desc;
```
