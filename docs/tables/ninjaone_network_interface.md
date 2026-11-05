# Table: ninjaone_network_interface

Network interfaces on each device. `ip_address`, `mac_address` and `dns_servers` are JSON arrays.

With `organization_id` set, the plugin passes it to NinjaOne as a device filter (`df`), so only that organization's devices are returned. A `device_id` filter is sent to the API and cannot widen that limit.

## Examples

### Basic info
List the network interfaces on each device.

```sql+postgres
select
  device_id,
  interface_name,
  interface_type,
  ip_address,
  mac_address,
  status
from
  ninjaone_network_interface;
```

```sql+sqlite
select
  device_id,
  interface_name,
  interface_type,
  ip_address,
  mac_address,
  status
from
  ninjaone_network_interface;
```

### Show the interfaces on one device
Match a device to its addresses during an investigation.

```sql+postgres
select
  interface_name,
  adapter_name,
  ip_address,
  mac_address,
  default_gateway
from
  ninjaone_network_interface
where
  device_id = 101;
```

```sql+sqlite
select
  interface_name,
  adapter_name,
  ip_address,
  mac_address,
  default_gateway
from
  ninjaone_network_interface
where
  device_id = 101;
```

### Count interfaces by DNS servers
Find devices that do not use the DNS servers your policy expects.

```sql+postgres
select
  dns_servers,
  count(*) as interfaces
from
  ninjaone_network_interface
group by
  dns_servers
order by
  interfaces desc;
```

```sql+sqlite
select
  dns_servers,
  count(*) as interfaces
from
  ninjaone_network_interface
group by
  dns_servers
order by
  interfaces desc;
```

### Count interfaces by default gateway
Show which networks devices are attached to.

```sql+postgres
select
  default_gateway,
  count(distinct device_id) as devices
from
  ninjaone_network_interface
where
  default_gateway is not null
group by
  default_gateway
order by
  devices desc;
```

```sql+sqlite
select
  default_gateway,
  count(distinct device_id) as devices
from
  ninjaone_network_interface
where
  default_gateway is not null
group by
  default_gateway
order by
  devices desc;
```
