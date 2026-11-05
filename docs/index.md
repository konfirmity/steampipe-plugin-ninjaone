---
organization: konfirmity
category: ["security"]
display_name: NinjaOne
name: ninjaone
description: "Steampipe plugin to query devices, patches, software, antivirus, health and users from NinjaOne."
og_description: "Query NinjaOne with SQL! Open source CLI. No DB required."
---

# NinjaOne + Steampipe

[NinjaOne](https://www.ninjaone.com) is a remote monitoring and management platform for endpoints.

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

For example:

```sql
select
  id,
  title,
  os_name,
  last_contact
from
  ninjaone_device;
```

```
+-----+----------------+------------+----------------------+
| id  | title          | os_name    | last_contact         |
+-----+----------------+------------+----------------------+
| 101 | EXAMPLE-LAP-01 | Windows 11 | 2026-11-05T09:12:00Z |
| 102 | EXAMPLE-MAC-02 | macOS 15   | 2026-11-05T08:47:00Z |
+-----+----------------+------------+----------------------+
```

## Documentation

- **[Table definitions and examples](/plugins/konfirmity/ninjaone/tables)**

## Quick start

### Install

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

### Credentials

The plugin signs in with a client app using the client credentials grant. It asks for the `monitoring` scope and only reads.

Create the client app in NinjaOne:

1. Go to **Administration > Apps > API > Client app IDs** and select **Add client app**.
2. Set **Application platform** to **API Services (machine-to-machine)**.
3. Enter a name.
4. Enter a **Redirect URI**. The form requires one, and `https://localhost` works.
5. Under **Scopes**, select **Monitoring** only.
6. Under **Allowed grant types**, select **Client credentials** only. Untick **Authorization code**.
7. Save the app and copy the client ID and client secret. The client secret is shown once.

### Configuration

The example config (`config/ninjaone.spc`) has a single connection named `ninjaone`:

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

Environment variables must be set where the Steampipe service starts. They are used only when the argument is not set in the connection.

## Limiting a connection to one organization

A client app cannot be limited to one organization in NinjaOne, so a managed service provider's credential sees every tenant. `organization_id` is what limits the plugin. With it set, every table returns only that organization:

| Tables | How the limit works |
|---|---|
| Device reports, `ninjaone_device`, `ninjaone_alert`, `ninjaone_activity` | NinjaOne's device filter (`df`) |
| `ninjaone_organization`, `ninjaone_location` | The plugin filters the rows it fetches |
| `ninjaone_user` | The organization's end-users endpoint. The provider's technicians are never returned. |
| `ninjaone_policy`, `ninjaone_role` | Only those the organization's devices use |

An empty `organization_id` counts as unset. The plugin does not check that the value is a real organization, so a mistyped identifier can look like an organization with nothing to report. To check a connection, run `select id, name from ninjaone_organization`: it should return exactly one row.

To find the identifier, run the same query on a connection without `organization_id`.

## What the data does not show

- **An empty patch table, or `HEALTHY` status, does not prove a device is patched.** Patching may be off in its policy. Software patching is off by default in the Mac policy we inspected, and Mac OS update settings sit under MDM. This plugin does not expose policy contents, so confirm the setting in the NinjaOne console.
- **NinjaOne exposes vulnerability counts per device, not CVE lists.** The counts are in `ninjaone_device_health`.
- **`chassis_type` does not tell laptops from desktops.** It reported `DESKTOP` for a MacBook Pro.
- **`bios_serial_number` was empty on a Mac.** Use `serial_number`.
- **`last_logged_in_user` and `ninjaone_logged_on_user` hold local login names, not email addresses.**

Backup usage and custom fields are deliberately not tables yet. `/v2/queries/backup/usage` returned 500 on the test tenant, and custom fields can hold sensitive values. For the same reason, the `userData` attributes on devices and locations are not exposed.

## Multiple Connections

To query several organizations, define one connection per organization and combine them with an [aggregator connection](https://steampipe.io/docs/managing/connections#using-aggregators). Each connection keeps its own `organization_id` limit:

```hcl
connection "ninjaone_example_one" {
  plugin          = "ninjaone"
  client_id       = "example-client-id-not-real"
  client_secret   = "example-client-secret-not-real"
  base_url        = "https://oc.ninjarmm.com"
  organization_id = "12345"
}

connection "ninjaone_example_two" {
  plugin          = "ninjaone"
  client_id       = "example-client-id-not-real"
  client_secret   = "example-client-secret-not-real"
  base_url        = "https://oc.ninjarmm.com"
  organization_id = "67890"
}

connection "ninjaone_all" {
  plugin      = "ninjaone"
  type        = "aggregator"
  connections = ["ninjaone_example_one", "ninjaone_example_two"]
}
```
