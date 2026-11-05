package ninjaone

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func testClient(t *testing.T, organizationID *int64, handler http.HandlerFunc) *client {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("grant_type") != "client_credentials" || r.Form.Get("scope") != "monitoring" {
			t.Errorf("token request form: %v", r.Form)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("%s sent without the bearer token", r.URL.Path)
		}
		handler(w, r)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	c := newClient("id", "secret", server.URL+"/", organizationID)
	c.wait = func(context.Context, time.Duration) error { return nil }
	return c
}

func collect(rows *[]string) func(json.RawMessage) (bool, error) {
	return func(raw json.RawMessage) (bool, error) {
		*rows = append(*rows, string(raw))
		return true, nil
	}
}

func TestDeviceFilter(t *testing.T) {
	organization, device := int64(42), int64(7)
	cases := []struct {
		name         string
		organization *int64
		device       *int64
		want         string
	}{
		{"unscoped", nil, nil, ""},
		{"organization", &organization, nil, "org=42"},
		{"device cannot widen the organization", &organization, &device, "org=42 AND id=7"},
		{"device only", nil, &device, "id=7"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &client{organizationID: tc.organization}
			if got := c.deviceFilter(tc.device); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestListCursorFollowsEveryPageAndKeepsTheFilter(t *testing.T) {
	organization := int64(42)
	var calls atomic.Int32
	c := testClient(t, &organization, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("df"); got != "org=42" {
			t.Errorf("page %d sent df %q", calls.Load()+1, got)
		}
		switch calls.Add(1) {
		case 1:
			fmt.Fprint(w, `{"cursor":{"name":"next-page","offset":0,"count":3},"results":[{"n":1},{"n":2}]}`)
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next-page" {
				t.Errorf("second page cursor %q", got)
			}
			fmt.Fprint(w, `{"cursor":{},"results":[{"n":3}]}`)
		default:
			t.Error("requested past the last page")
		}
	})
	var rows []string
	if err := c.listCursor(context.Background(), "/v2/queries/software", c.deviceFilter(nil), -1, collect(&rows)); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("got %d rows, want 3", len(rows))
	}
}

func TestListFailsWhenALaterPageFails(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			fmt.Fprint(w, `{"cursor":{"name":"next-page"},"results":[{"n":1}]}`)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errorMessage":"boom"}`)
	})
	var rows []string
	err := c.listCursor(context.Background(), "/v2/queries/software", "", -1, collect(&rows))
	if err == nil {
		t.Fatal("a partial listing must be an error, or a sync could read it as devices gone")
	}
}

func TestListAfterUsesTheHighestIDSeen(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		switch calls.Add(1) {
		case 1:
			fmt.Fprint(w, `[{"id":5},{"id":9},{"id":7}]`)
		case 2:
			if got := r.URL.Query().Get("after"); got != "9" {
				t.Errorf("after = %q, want 9", got)
			}
			fmt.Fprint(w, `[]`)
		}
	})
	var rows []string
	if err := c.listAfter(context.Background(), "/v2/devices-detailed", "", -1, collect(&rows)); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("got %d rows, want 3", len(rows))
	}
}

func TestListActivitiesWalksBackward(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		switch calls.Add(1) {
		case 1:
			if got := r.URL.Query().Get("newerThan"); got != "2" {
				t.Errorf("newerThan = %q", got)
			}
			fmt.Fprint(w, `{"lastActivityId":4,"activities":[{"id":5},{"id":4}]}`)
		case 2:
			if got := r.URL.Query().Get("olderThan"); got != "4" {
				t.Errorf("olderThan = %q, want 4", got)
			}
			fmt.Fprint(w, `{"activities":[]}`)
		}
	})
	newerThan := int64(2)
	var rows []string
	if err := c.listActivities(context.Background(), "", &newerThan, -1, collect(&rows)); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}
}

func TestStopsEarlyWhenTheCallerHasEnoughRows(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if got := r.URL.Query().Get("pageSize"); got != "1" {
			t.Errorf("pageSize = %q, want the query limit", got)
		}
		fmt.Fprint(w, `{"cursor":{"name":"more"},"results":[{"n":1}]}`)
	})
	err := c.listCursor(context.Background(), "/v2/queries/software", "", 1, func(json.RawMessage) (bool, error) { return false, nil })
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatalf("made %d requests, want 1", calls.Load())
	}
}

func TestRetriesThrottlingThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) < 3 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		fmt.Fprint(w, `[]`)
	})
	var rows []json.RawMessage
	if err := c.get(context.Background(), "/v2/policies", nil, &rows); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 3 {
		t.Fatalf("made %d requests, want 3", calls.Load())
	}
}

func TestGivesUpAfterMaxAttempts(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	var rows []json.RawMessage
	if err := c.get(context.Background(), "/v2/policies", nil, &rows); err == nil {
		t.Fatal("expected an error")
	}
	if calls.Load() != maxAttempts {
		t.Fatalf("made %d requests, want %d", calls.Load(), maxAttempts)
	}
}

func TestDoesNotRetryOrdinaryErrors(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusNotFound)
	})
	var device ninjaDevice
	err := c.get(context.Background(), "/v2/device/1", nil, &device)
	if !isNotFound(err) {
		t.Fatalf("got %v", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("made %d requests, want 1", calls.Load())
	}
}

func TestErrorBodyIsTruncated(t *testing.T) {
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, strings.Repeat("x", errorBodyLimit*2))
	})
	var rows []json.RawMessage
	err := c.get(context.Background(), "/v2/policies", nil, &rows)
	if err == nil {
		t.Fatal("expected an error")
	}
	if len(err.Error()) > errorBodyLimit+100 {
		t.Fatalf("error body not truncated: %d bytes", len(err.Error()))
	}
}

func TestOnlyDevicesDropsRowsFromOtherDevices(t *testing.T) {
	var rows []string
	keep := onlyDevices(map[int64]bool{1: true}, collect(&rows))
	for _, row := range []string{`{"deviceId":1}`, `{"deviceId":2}`, `{"message":"tenant-wide"}`} {
		more, err := keep(json.RawMessage(row))
		if err != nil || !more {
			t.Fatalf("row %s: more=%v err=%v", row, more, err)
		}
	}
	if len(rows) != 1 || rows[0] != `{"deviceId":1}` {
		t.Fatalf("got %v", rows)
	}
}

func TestOnlyDevicesPassesEverythingOnAnUnscopedConnection(t *testing.T) {
	var rows []string
	keep := onlyDevices(nil, collect(&rows))
	if _, err := keep(json.RawMessage(`{"message":"tenant-wide"}`)); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %v", rows)
	}
}

func TestAnEmptyAllowListKeepsNothing(t *testing.T) {
	var rows []string
	keep := onlyDevices(map[int64]bool{}, collect(&rows))
	if _, err := keep(json.RawMessage(`{"deviceId":1}`)); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("got %v", rows)
	}
}

func TestScopedDevicesFailsForAnOrganizationThatDoesNotExist(t *testing.T) {
	organization := int64(42)
	c := testClient(t, &organization, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/organization/42" {
			t.Errorf("requested %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := c.scopedDevices(context.Background()); !isNotFound(err) {
		t.Fatalf("a wrong organization must not read as an empty tenant, got %v", err)
	}
}

func TestScopedDevicesIgnoresDevicesFromOtherOrganizations(t *testing.T) {
	organization := int64(42)
	var devicePages atomic.Int32
	c := testClient(t, &organization, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/organization/42" {
			fmt.Fprint(w, `{"id":42}`)
			return
		}
		if devicePages.Add(1) == 1 {
			fmt.Fprint(w, `[{"id":1,"organizationId":42},{"id":2,"organizationId":43}]`)
			return
		}
		fmt.Fprint(w, `[]`)
	})
	devices, err := c.scopedDevices(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !devices[1] || devices[2] {
		t.Fatalf("got %v", devices)
	}
}

func TestDeviceReferencesComeOnlyFromThisOrganizationsDevices(t *testing.T) {
	organization := int64(42)
	var devicePages atomic.Int32
	c := testClient(t, &organization, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/organization/42":
			fmt.Fprint(w, `{"id":42,"policies":[{"nodeRoleId":9,"policyId":90}]}`)
		case "/v2/devices-detailed":
			if devicePages.Add(1) == 1 {
				fmt.Fprint(w, `[{"id":1,"organizationId":42,"nodeRoleId":2,"rolePolicyId":20,"policyId":21},{"id":2,"organizationId":43,"nodeRoleId":8,"rolePolicyId":80}]`)
				return
			}
			fmt.Fprint(w, `[]`)
		default:
			t.Errorf("requested %s", r.URL.Path)
		}
	})
	policies, roles, err := c.deviceReferences(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !policies[20] || !policies[21] || policies[80] || policies[90] {
		t.Fatalf("policies: %v", policies)
	}
	if !roles[2] || roles[8] || roles[9] {
		t.Fatalf("roles: %v", roles)
	}
}

func TestListAfterStopsIfAnEndpointIgnoresPaging(t *testing.T) {
	c := testClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":5}]`)
	})
	var rows []string
	if err := c.listAfter(context.Background(), "/v2/locations", "", -1, collect(&rows)); err == nil {
		t.Fatal("expected an error instead of an endless loop")
	}
}

func TestInScope(t *testing.T) {
	organization := int64(42)
	scoped, unscoped := &client{organizationID: &organization}, &client{}
	if !scoped.inScope(42) || scoped.inScope(43) || !unscoped.inScope(43) {
		t.Fatal("scope check is wrong")
	}
}

func TestParseOrganizationID(t *testing.T) {
	for _, bad := range []string{"abc", "0", "-1", "1 OR 1=1", "1.5"} {
		if _, err := parseOrganizationID(bad); err == nil {
			t.Errorf("%q should be rejected", bad)
		}
	}
	id, err := parseOrganizationID("42")
	if err != nil || id == nil || *id != 42 {
		t.Fatalf("got %v, %v", id, err)
	}
	if id, err := parseOrganizationID(""); err != nil || id != nil {
		t.Fatalf("got %v, %v", id, err)
	}
}

func TestFlexStringAcceptsWhatTheSpecDoesNotPromise(t *testing.T) {
	var row struct {
		A flexString `json:"a"`
		B flexString `json:"b"`
		C flexString `json:"c"`
	}
	if err := json.Unmarshal([]byte(`{"a":"text","b":123,"c":null}`), &row); err != nil {
		t.Fatal(err)
	}
	if row.A != "text" || row.B != "123" || row.C != "" {
		t.Fatalf("got %+v", row)
	}
}
