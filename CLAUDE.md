# CLAUDE.md

A read-only Steampipe plugin (Go) that exposes the NinjaOne public API as Postgres tables. The Konfirmity proxy builds it from a git tag as a Postgres FDW extension and writes one connection per client.

## Tenancy is the main invariant

A client app cannot be limited to one organization in NinjaOne, so a managed service provider's credential sees every tenant. `organization_id` is the only thing that limits the plugin. With it set, every table returns only that organization:

- Device reports, `ninjaone_device`, `ninjaone_alert`, `ninjaone_activity`: the `df` filter from `client.deviceFilter`.
- `ninjaone_device` lookup by `id`: fetched, then checked with `client.inScope`.
- `ninjaone_organization`, `ninjaone_location`: rows filtered with `client.inScope`.
- `ninjaone_user`: `/v2/organization/{id}/end-users`, so the provider's technicians are never returned.
- `ninjaone_policy`, `ninjaone_role`: `client.deviceReferences` lists what the organization's devices use. The organization's own policy map cannot be used: it lists every role in the account. `/v2/policies` and `/v2/roles` span every tenant and have no server-side filter.

Every new table states which of these it uses, and gets a row in `test/fake_ninjaone.py` and a check in `test/verify-scope.sh`. The fake returns every tenant's rows unless `df` narrows them, so a table that skips the filter fails the check.

## Verifying a change

- `go build ./...`, `go test ./...`, `go vet ./...` and `gofmt -l .`: the four CI steps in `.github/workflows/build.yml`.
- `test/verify-scope.sh` runs real SQL through Steampipe against the fake and fails on any cross-tenant row. CI does not run it, so run it by hand after touching a table or the client. Leave `STEAMPIPE_INSTALL_DIR` unset: the script writes a fake `ninjaone.spc` into that directory and force-stops its service.
- Fake values are not API values. `os-patches` rows in the fake use status `PENDING`; a third-party API wrapper lists `MANUAL`, `APPROVED`, `FAILED` and `REJECTED`. Check a value against a real tenant before a doc or query depends on it.

## SDK pin

Stay on `steampipe-plugin-sdk` v6.0.0, not v6.1.0: the proxy's FDW fork pins it. Check `go.mod` after any `go get` or `go mod tidy`.

## Code conventions

- Retries live in `client.get`, not the SDK's hydrate retry, because a hydrate retry restarts a list and streams its rows twice.
- Report rows use `flexString` so one mistyped field cannot fail a whole table.
- The default transform is `nullIfEmpty`, not the SDK's `NullIfZero`: a zero threat count is a fact and stays `0`.
- JSON columns backed by `json.RawMessage` hold JSON `null`, not SQL NULL, when the API omits them. Queries test them with `->>`, not `is null`.
- `ninjaone_software_patch.patch_title` is the API's `title`, because `title` is the standard Steampipe column.
- `make install` writes to `plugins/local/ninjaone`, so a local connection needs `plugin = "local/ninjaone"`.

## Facts found by testing

These are easy to get wrong in a table, a query or a doc:

- An empty patch table, or `HEALTHY` status, does not prove a device is patched: patching may be off in its policy. Software patching is off by default in the Mac policy we inspected, and Mac OS update settings sit under MDM.
- NinjaOne exposes vulnerability counts per device (`ninjaone_device_health`), not CVE lists.
- `chassis_type` reported `DESKTOP` for a MacBook Pro, so it cannot tell laptops from desktops. `bios_serial_number` was empty on a Mac; use `serial_number`.
- `last_logged_in_user` and `ninjaone_logged_on_user` are local login names, not email addresses.
- `ninjaone_activity` supports `id > n`, sent as `newerThan`, to read only new rows.

## Deliberately not tables yet

- Backup usage: `/v2/queries/backup/usage` returned 500 on the test tenant.
- Custom fields: they can hold sensitive values.

Raise either with a person before adding it.

## Docs

`docs/tables/<table>.md` has `# Table: <name>`, a description, then `## Examples` with `### Basic info` first (select and from only, named columns) and 3-4 security or compliance queries with imperative titles. Each query appears as a `sql+postgres` and a `sql+sqlite` block. Every column must exist in that table's Go file. The config block in `docs/index.md` and `README.md` is a copy of `config/ninjaone.spc`, so change all three together. Use fake example values, and no real tenant names or secrets.
