# Table: ninjaone_user

NinjaOne users: technicians and end users.

With `organization_id` set, the plugin reads the organization's end-users endpoint instead of the full user list, so the provider's technicians are never returned. The columns about technician access, such as `administrator` and `permit_all_clients`, are only useful on a connection without `organization_id`.

`email` does not match the local login names in `ninjaone_logged_on_user` or `ninjaone_device.last_logged_in_user`. Use `ninjaone_device.assigned_owner_uid` to link a device to a user.

## Examples

### Basic info
List users.

```sql+postgres
select
  id,
  title,
  email,
  user_type,
  enabled,
  mfa_configured
from
  ninjaone_user;
```

```sql+sqlite
select
  id,
  title,
  email,
  user_type,
  enabled,
  mfa_configured
from
  ninjaone_user;
```

### Find enabled users without multi-factor authentication
Show accounts that need MFA set up.

```sql+postgres
select
  title,
  email,
  user_type
from
  ninjaone_user
where
  enabled = true
  and mfa_configured is not true;
```

```sql+sqlite
select
  title,
  email,
  user_type
from
  ninjaone_user
where
  enabled = true
  and mfa_configured is not true;
```

### List administrators
Review who holds the highest access. Run this on a connection without `organization_id`.

```sql+postgres
select
  title,
  email,
  enabled,
  mfa_configured
from
  ninjaone_user
where
  administrator = true;
```

```sql+sqlite
select
  title,
  email,
  enabled,
  mfa_configured
from
  ninjaone_user
where
  administrator = true;
```

### Find technicians who can reach every organization
Review access across tenants. Run this on a connection without `organization_id`.

```sql+postgres
select
  title,
  email,
  mfa_configured
from
  ninjaone_user
where
  user_type = 'TECHNICIAN'
  and permit_all_clients = true;
```

```sql+sqlite
select
  title,
  email,
  mfa_configured
from
  ninjaone_user
where
  user_type = 'TECHNICIAN'
  and permit_all_clients = true;
```

### Find disabled users who still own devices
Offboarding should reassign devices.

```sql+postgres
select
  u.email,
  d.id as device_id,
  d.title as device
from
  ninjaone_user as u
  join ninjaone_device as d on d.assigned_owner_uid = u.uid
where
  u.enabled = false;
```

```sql+sqlite
select
  u.email,
  d.id as device_id,
  d.title as device
from
  ninjaone_user as u
  join ninjaone_device as d on d.assigned_owner_uid = u.uid
where
  u.enabled = false;
```
