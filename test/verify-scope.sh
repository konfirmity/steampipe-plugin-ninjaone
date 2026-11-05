#!/usr/bin/env bash
# Runs real SQL through Steampipe against test/fake_ninjaone.py and fails if a
# connection limited to one organization ever returns another tenant's data.
# Usage: STEAMPIPE=/path/to/steampipe test/verify-scope.sh
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
steampipe="${STEAMPIPE:-steampipe}"
port="${FAKE_PORT:-48762}"
export STEAMPIPE_INSTALL_DIR="${STEAMPIPE_INSTALL_DIR:-$(mktemp -d)}"
export STEAMPIPE_UPDATE_CHECK=false STEAMPIPE_TELEMETRY=none STEAMPIPE_CACHE=false

mkdir -p "$STEAMPIPE_INSTALL_DIR/plugins/local/ninjaone" "$STEAMPIPE_INSTALL_DIR/config"
(cd "$root" && go build -o "$STEAMPIPE_INSTALL_DIR/plugins/local/ninjaone/ninjaone.plugin" -tags netgo .)

connect_to() {
  cat > "$STEAMPIPE_INSTALL_DIR/config/ninjaone.spc" <<EOF
connection "ninjaone" {
  plugin          = "local/ninjaone"
  client_id       = "fake"
  client_secret   = "fake"
  base_url        = "http://127.0.0.1:$port"
  organization_id = "$1"
}
EOF
}

fake=""
start_fake() {
  stop_fake
  python3 -I "$root/test/fake_ninjaone.py" "$port" "$@" &
  fake=$!
  sleep 1
}
stop_fake() {
  if [[ -n "$fake" ]]; then
    kill "$fake" 2>/dev/null || true
    wait "$fake" 2>/dev/null || true
    fake=""
  fi
}
cleanup() {
  stop_fake
  "$steampipe" service stop --force >/dev/null 2>&1 || true
}
trap cleanup EXIT

failures=0
query() { "$steampipe" query --output csv --header=false "$1"; }

check() {
  local label="$1" sql="$2" want="$3" got
  got="$(query "$sql" 2>&1 | tr -d '\r' | paste -sd, - || true)"
  if [[ "$got" == "$want" ]]; then
    printf 'ok    %s\n' "$label"
  else
    printf 'FAIL  %s: got [%s], want [%s]\n' "$label" "$got" "$want"
    failures=$((failures + 1))
  fi
}

refuses() {
  local label="$1" sql="$2"
  if query "$sql" >/dev/null 2>&1; then
    printf 'FAIL  %s: the query succeeded\n' "$label"
    failures=$((failures + 1))
  else
    printf 'ok    %s\n' "$label"
  fi
}

tenant_checks() {
  local table
  for table in device_health software os_patch os_patch_install software_patch software_patch_install \
    antivirus_status antivirus_threat volume disk network_interface processor computer_system \
    operating_system logged_on_user windows_service policy_override alert; do
    check "ninjaone_$table returns only this tenant's devices" \
      "select string_agg(distinct device_id::text, '|' order by device_id::text) from ninjaone_$table" "1|2|3"
  done
  check "ninjaone_device pages and stays in tenant" "select string_agg(id::text, '|' order by id) from ninjaone_device" "1|2|3"
  check "ninjaone_device get cannot reach another tenant" "select count(*) from ninjaone_device where id = 4" "0"
  check "ninjaone_device get reaches this tenant" "select count(*) from ninjaone_device where id = 2" "1"
  check "device_id qual cannot widen the tenant" "select count(*) from ninjaone_software where device_id = 4" "0"
  check "ninjaone_activity excludes other tenants and tenant-wide rows" "select string_agg(device_id::text, '|' order by device_id) from ninjaone_activity" "1|2|3"
  check "ninjaone_organization" "select string_agg(name, '|') from ninjaone_organization" "Own"
  check "ninjaone_location" "select string_agg(name, '|') from ninjaone_location" "Loc1"
  check "ninjaone_user hides technicians and other tenants" "select string_agg(email, '|') from ninjaone_user" "own@client.example"
  check "ninjaone_policy" "select string_agg(name, '|') from ninjaone_policy" "OwnPolicy"
  check "ninjaone_role" "select string_agg(name, '|') from ninjaone_role" "OwnRole"
}

echo "--- the API honours the organization filter"
connect_to 1
start_fake
tenant_checks
check "zero stays zero, not null" "select pending_os_patches_count, active_threats_count, offline from ninjaone_device_health where device_id = 1" "0,0,false"
check "numeric kb number survives" "select kb_number from ninjaone_os_patch where device_id = 1" "500123"
check "bitlocker detail is readable" "select bit_locker_status ->> 'protectionStatus' from ninjaone_volume where device_id = 1" "PROTECTED"
check "a missing json value is a real null" "select count(*) from ninjaone_device where id = 1 and maintenance is null" "1"
check "serial and nested fields" "select serial_number, os_name, os_needs_reboot from ninjaone_device where id = 1" "SER1,Windows 11,true"
check "tags map" "select tags ->> 'finance' from ninjaone_device where id = 1" "true"
check "timestamps are real times" "select extract(year from last_contact)::int from ninjaone_device where id = 1" "2026"
check "logon time in milliseconds is a real time" "select extract(year from logon_time)::int from ninjaone_logged_on_user where device_id = 1" "2026"

echo "--- the API ignores the organization filter, so only the plugin's own check holds"
start_fake ignore-df
"$steampipe" service stop --force >/dev/null 2>&1 || true
tenant_checks

echo "--- an organization that does not exist is an error, not an empty tenant"
start_fake
connect_to 999
"$steampipe" service stop --force >/dev/null 2>&1 || true
for table in device software os_patch alert activity organization location policy role; do
  refuses "ninjaone_$table refuses a missing organization" "select count(*) from ninjaone_$table"
done

if ((failures > 0)); then
  echo "$failures check(s) failed"
  exit 1
fi
echo "all checks passed"
