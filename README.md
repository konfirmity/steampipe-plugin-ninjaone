# NinjaOne Plugin for Steampipe

Use SQL to query devices, patches, software, antivirus, health and users from the NinjaOne public API. The plugin only reads.

- [Plugin overview](docs/index.md)
- [Table docs](docs/tables)
- [Changelog](CHANGELOG.md)

## Quick start

Build from source and copy the example config:

```sh
make install
cp config/ninjaone.spc ~/.steampipe/config/ninjaone.spc
```

Then edit `~/.steampipe/config/ninjaone.spc`:

1. Replace the example `client_id` and `client_secret` with your own.
2. Change `plugin` to `"local/ninjaone"`, because `make install` puts the binary under `plugins/local/ninjaone`.
3. Set `base_url` if you are not on `https://app.ninjarmm.com`.
4. Set `organization_id` if your credential can see more than one organization.

Run a query:

```sh
steampipe query "select id, title, os_name from ninjaone_device limit 10"
```

## Credentials

The plugin signs in with a client app using the client credentials grant. It asks for the `monitoring` scope.

Create the client app in NinjaOne:

1. Go to **Administration > Apps > API > Client app IDs** and select **Add client app**.
2. Set **Application platform** to **API Services (machine-to-machine)**.
3. Enter a name.
4. Enter a **Redirect URI**. The form requires one, and `https://localhost` works.
5. Under **Scopes**, select **Monitoring** only.
6. Under **Allowed grant types**, select **Client credentials** only. Untick **Authorization code**.
7. Save the app and copy the client ID and client secret. The client secret is shown once.

## Configuration

```hcl
connection "ninjaone" {
  plugin = "ninjaone"

  # `make install` puts the binary under plugins/local/ninjaone. For a build
  # installed that way, change the line above to: plugin = "local/ninjaone"

  # Client ID of an API Services (machine-to-machine) client app.
  # Can also be set with the NINJAONE_CLIENT_ID environment variable.
  # A value here wins over the environment variable, so delete this line to use the variable.
  client_id = "example-client-id-not-real"

  # Client secret of the same client app. NinjaOne shows it once, when the app is created.
  # Can also be set with the NINJAONE_CLIENT_SECRET environment variable.
  # A value here wins over the environment variable, so delete this line to use the variable.
  client_secret = "example-client-secret-not-real"

  # Regional instance, including https://. Defaults to https://app.ninjarmm.com.
  # https://oc.ninjarmm.com is verified. NinjaOne also has eu, ca and us2 regions:
  # for those, copy the host from your browser's address bar while signed in to NinjaOne.
  # Can also be set with the NINJAONE_BASE_URL environment variable.
  # base_url = "https://oc.ninjarmm.com"

  # Limits every table to one NinjaOne organization. Write it as a quoted whole number.
  # A client app cannot be limited to one organization in NinjaOne, so a managed service
  # provider's credential sees every tenant. This setting is what limits the plugin.
  # Leave it unset and every table returns every organization the credential can see.
  # Can also be set with the NINJAONE_ORGANIZATION_ID environment variable.
  # organization_id = "12345"
}
```

## Limiting a connection to one organization

A client app cannot be limited to one organization in NinjaOne, so a managed service provider's credential sees every tenant. `organization_id` is what limits the plugin. With it set, every table returns only that organization:

| Tables | How the limit works |
|---|---|
| Device reports, `ninjaone_device`, `ninjaone_alert`, `ninjaone_activity` | NinjaOne's device filter (`df`) |
| `ninjaone_organization`, `ninjaone_location` | The plugin filters the rows it fetches |
| `ninjaone_user` | The organization's end-users endpoint. The provider's technicians are never returned. |
| `ninjaone_policy`, `ninjaone_role` | Only those the organization's devices use |

An empty `organization_id` counts as unset. The plugin does not check that the value is a real organization. To check a connection, run `select id, name from ninjaone_organization`: it should return exactly one row.

## What the data does not show

- An empty patch table, or `HEALTHY` status, does not prove a device is patched. Patching may be off in its policy. Software patching is off by default in the Mac policy we inspected, and Mac OS update settings sit under MDM. This plugin does not expose policy contents, so confirm the setting in the NinjaOne console.
- NinjaOne exposes vulnerability counts per device (`ninjaone_device_health`), not CVE lists.
- `chassis_type` reported `DESKTOP` for a MacBook Pro, so do not use it to tell laptops from desktops.
- `bios_serial_number` was empty on a Mac. Use `serial_number`.
- `last_logged_in_user` and `ninjaone_logged_on_user` hold local login names, not email addresses.

Backup usage and custom fields are deliberately not tables yet. `/v2/queries/backup/usage` returned 500 on the test tenant, and custom fields can hold sensitive values. For the same reason, the `userData` attributes on devices and locations are not exposed.

## Example queries

```sql
-- Devices not seen for 30 days
select
  id,
  title,
  last_contact
from
  ninjaone_device
where
  last_contact < now() - interval '30 days'
order by
  last_contact;

-- Devices with critical vulnerabilities (counts, not CVE lists)
select
  device_id,
  critical_vulnerability_count
from
  ninjaone_device_health
where
  critical_vulnerability_count > 0
order by
  critical_vulnerability_count desc;

-- Devices with no antivirus detected
select
  device_id
from
  ninjaone_antivirus_status
where
  product_name = 'NONE';

-- Enabled users without multi-factor authentication
select
  email,
  user_type
from
  ninjaone_user
where
  enabled = true
  and mfa_configured is not true;

-- Activity newer than the last id you processed
select
  id,
  activity_time,
  device_id,
  message
from
  ninjaone_activity
where
  id > 12345
order by
  id;
```

## Development

```sh
go build ./...          # compile
go test ./...           # unit tests
go vet ./...
make install            # build into ~/.steampipe/plugins/local/ninjaone/
make sqlite             # SQLite extension
make postgres           # Postgres FDW extension, needs pg_config
```

`test/verify-scope.sh` runs real SQL through Steampipe against `test/fake_ninjaone.py`, a two-tenant fake of the API, and fails on any cross-tenant row. It needs Steampipe and Python 3. Set `STEAMPIPE` if the binary is not on your `PATH`:

```sh
STEAMPIPE=/path/to/steampipe test/verify-scope.sh
```

Leave `STEAMPIPE_INSTALL_DIR` unset when you run it. The script then uses a temporary directory. If the variable is set, the script overwrites `ninjaone.spc` there and stops that install's service.

The plugin must stay on `steampipe-plugin-sdk` v6.0.0, not v6.1.0, because the proxy's FDW fork pins it.

---

NinjaOne is a trademark of its owner, used here for identification only. This is a community plugin and is not affiliated with or endorsed by NinjaOne.
