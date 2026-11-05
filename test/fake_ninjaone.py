"""A two-tenant fake of the NinjaOne API for test/verify-scope.sh.

Like the real API, a report returns every tenant's rows unless df narrows them,
so a table that forgets the organization filter shows up as a leak. Each page
holds two rows, so paging is exercised on every table. With a second argument of
ignore-df the fake disregards df entirely, which proves the plugin's own check.
"""
import json
import re
import sys
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.parse import parse_qs, urlparse

PAGE = 2
IGNORE_DF = len(sys.argv) > 2 and sys.argv[2] == "ignore-df"
NOW = 1793872707.589
OWN, OTHER = 1, 2


def device(id, org):
    return {
        "id": id, "uid": f"uid-{id}", "organizationId": org, "locationId": org, "nodeClass": "WINDOWS_WORKSTATION",
        "nodeRoleId": 100 + org, "rolePolicyId": 200 + org, "approvalStatus": "APPROVED", "offline": False,
        "systemName": f"org{org}-pc{id}", "created": NOW, "lastContact": NOW, "lastUpdate": NOW,
        "tags": ["finance"], "ipAddresses": ["10.0.0.1"], "macAddresses": ["00:11"], "publicIP": "203.0.113.5",
        "os": {"name": "Windows 11", "needsReboot": True, "lastBootTime": NOW},
        "system": {"serialNumber": f"SER{id}", "virtualMachine": False, "totalPhysicalMemory": 1024},
    }


DEVICES = [device(1, OWN), device(2, OWN), device(3, OWN), device(4, OTHER), device(5, OTHER)]


def per_device(fields):
    return [{**fields(d), "deviceId": d["id"], "timestamp": NOW} for d in DEVICES]


REPORTS = {
    "device-health": per_device(lambda d: {"healthStatus": "NEEDS_ATTENTION", "pendingOSPatchesCount": 0, "activeThreatsCount": 0, "offline": False, "productsInstallationStatuses": {}}),
    "software": per_device(lambda d: {"name": f"App{d['id']}", "version": "1.0", "publisher": "Acme", "isSystemComponent": False, "size": 10}),
    "os-patches": per_device(lambda d: {"id": f"KB{d['id']}", "name": "Update", "severity": "CRITICAL", "status": "PENDING", "type": "SECURITY_UPDATES", "kbNumber": 500123}),
    "os-patch-installs": per_device(lambda d: {"id": f"KB{d['id']}", "name": "Update", "status": "INSTALLED", "installedAt": NOW}),
    "software-patches": per_device(lambda d: {"id": f"S{d['id']}", "productIdentifier": "chrome", "title": "Chrome", "impact": "CRITICAL", "status": "PENDING"}),
    "software-patch-installs": per_device(lambda d: {"id": f"S{d['id']}", "title": "Chrome", "status": "FAILED", "installedAt": NOW}),
    "antivirus-status": per_device(lambda d: {"productName": "Bitdefender", "productState": "ON", "definitionStatus": "UP_TO_DATE"}),
    "antivirus-threats": per_device(lambda d: {"name": "Eicar", "status": "QUARANTINED", "threatId": 12345, "level": "HIGH"}),
    "volumes": per_device(lambda d: {"name": "C:", "capacity": 100, "freeSpace": 50, "bitLockerStatus": {"protectionStatus": "PROTECTED", "conversionStatus": "FULLY_ENCRYPTED"}}),
    "disks": per_device(lambda d: {"name": "disk0", "size": 100, "smartCapable": True}),
    "network-interfaces": per_device(lambda d: {"interfaceName": "en0", "ipAddress": ["10.0.0.1"], "macAddress": ["00:11"], "dnsServers": ["1.1.1.1"]}),
    "processors": per_device(lambda d: {"name": "Core", "numCores": 8}),
    "computer-systems": per_device(lambda d: {"name": d["systemName"], "serialNumber": d["system"]["serialNumber"]}),
    "operating-systems": per_device(lambda d: {"name": "Windows 11", "needsReboot": True, "lastBootTime": NOW}),
    "logged-on-users": per_device(lambda d: {"userName": f"user{d['id']}", "logonTime": NOW * 1000}),
    "windows-services": per_device(lambda d: {"name": "wuauserv", "state": "RUNNING", "startType": "AUTO_START"}),
    "policy-overrides": per_device(lambda d: {"overrides": ["patching"]}),
}
ORGS = [
    {"id": OWN, "name": "Own", "nodeApprovalMode": "AUTOMATIC", "policies": [{"nodeRoleId": 101, "policyId": 201}, {"nodeRoleId": 102, "policyId": 202}], "locations": [], "settings": {}},
    {"id": OTHER, "name": "OtherTenant", "nodeApprovalMode": "MANUAL", "policies": [{"nodeRoleId": 102, "policyId": 202}], "locations": [], "settings": {}},
]
LOCATIONS = [{"id": o["id"], "name": f"Loc{o['id']}", "organizationId": o["id"]} for o in ORGS]
POLICIES = [{"id": 201, "name": "OwnPolicy", "nodeClass": "MAC", "nodeClassDefault": False, "updated": NOW}, {"id": 202, "name": "OtherPolicy", "nodeClass": "MAC", "nodeClassDefault": False, "updated": NOW}]
ROLES = [{"id": 101, "name": "OwnRole", "nodeClass": "MAC", "custom": True, "created": NOW}, {"id": 102, "name": "OtherRole", "nodeClass": "MAC", "custom": True, "created": NOW}]
USERS = [
    {"id": 1, "uid": "u1", "email": "tech@msp.example", "userType": "TECHNICIAN", "organizationId": None, "mfaConfigured": True, "enabled": True},
    {"id": 2, "uid": "u2", "firstName": "Own", "lastName": "User", "email": "own@client.example", "userType": "END_USER", "organizationId": OWN, "mfaConfigured": False, "enabled": True, "assignedDeviceIds": [1]},
    {"id": 3, "uid": "u3", "email": "other@tenant.example", "userType": "END_USER", "organizationId": OTHER, "mfaConfigured": False, "enabled": True},
]
ALERTS = [{"uid": f"a{d['id']}", "deviceId": d["id"], "subject": "Disk low", "message": "m", "sourceType": "CONDITION", "createTime": NOW, "updateTime": NOW, "data": {}} for d in DEVICES]
ACTIVITIES = [{"id": 10 + d["id"], "activityTime": NOW, "deviceId": d["id"], "activityType": "PATCH_MANAGEMENT", "status": "OK", "statusCode": "OK", "message": "scan", "type": "t", "data": {}} for d in DEVICES]
ACTIVITIES.append({"id": 99, "activityTime": NOW, "activityType": "SYSTEM", "status": "OK", "statusCode": "OK", "message": "tenant-wide", "type": "t", "data": {}})


def matching_device_ids(df):
    ids = {d["id"] for d in DEVICES}
    if not df:
        return ids
    for term in re.split(r"\s+AND\s+", df):
        key, value = re.fullmatch(r"\s*(org|id)\s*=\s*(\d+)\s*", term).groups()
        ids &= {d["id"] for d in DEVICES if (d["organizationId"] if key == "org" else d["id"]) == int(value)}
    return ids


class Handler(BaseHTTPRequestHandler):
    def send(self, body, status=200):
        data = json.dumps(body).encode()
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def do_POST(self):
        self.rfile.read(int(self.headers.get("content-length", 0)))
        self.send({"access_token": "fake", "token_type": "Bearer", "expires_in": 3600, "scope": "monitoring"})

    def do_GET(self):
        url = urlparse(self.path)
        query = {k: v[0] for k, v in parse_qs(url.query).items()}
        path = url.path
        size = min(int(query.get("pageSize", PAGE)), PAGE)
        df = None if IGNORE_DF else query.get("df")
        ids = matching_device_ids(df)

        report = re.fullmatch(r"/v2/queries/(.+)", path)
        if report:
            rows = [r for r in REPORTS[report[1]] if r["deviceId"] in ids]
            offset = int(query.get("cursor", "0"))
            page = rows[offset:offset + size]
            more = offset + size < len(rows)
            return self.send({"cursor": {"name": str(offset + size), "offset": offset, "count": len(rows)} if more else {}, "results": page})
        if path in ("/v2/devices-detailed", "/v2/organizations-detailed", "/v2/locations"):
            source = {"/v2/devices-detailed": [d for d in DEVICES if d["id"] in ids], "/v2/organizations-detailed": ORGS, "/v2/locations": LOCATIONS}[path]
            after = int(query.get("after", 0))
            return self.send([r for r in source if r["id"] > after][:size])
        organization = re.fullmatch(r"/v2/organization/(\d+)", path)
        if organization:
            found = [o for o in ORGS if o["id"] == int(organization[1])]
            return self.send(found[0]) if found else self.send({"errorMessage": "not found"}, 404)
        one = re.fullmatch(r"/v2/device/(\d+)", path)
        if one:
            found = [d for d in DEVICES if d["id"] == int(one[1])]
            return self.send(found[0]) if found else self.send({"errorMessage": "not found"}, 404)
        end_users = re.fullmatch(r"/v2/organization/(\d+)/end-users", path)
        if end_users:
            return self.send([u for u in USERS if u["organizationId"] == int(end_users[1])])
        if path == "/v2/users":
            return self.send(USERS)
        if path == "/v2/policies":
            return self.send(POLICIES)
        if path == "/v2/roles":
            return self.send(ROLES)
        if path == "/v2/alerts":
            return self.send([a for a in ALERTS if a["deviceId"] in ids])
        if path == "/v2/activities":
            rows = [a for a in ACTIVITIES if (a.get("deviceId") in ids) or (not df)]
            rows.sort(key=lambda a: -a["id"])
            if "olderThan" in query:
                rows = [a for a in rows if a["id"] < int(query["olderThan"])]
            if "newerThan" in query:
                rows = [a for a in rows if a["id"] > int(query["newerThan"])]
            return self.send({"activities": rows[:size]})
        self.send({"errorMessage": path}, 404)

    def log_message(self, *args):
        pass


ThreadingHTTPServer(("127.0.0.1", int(sys.argv[1])), Handler).serve_forever()
