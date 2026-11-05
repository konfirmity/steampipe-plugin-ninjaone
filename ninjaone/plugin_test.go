package ninjaone

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func TestTablesFollowTheHubStandards(t *testing.T) {
	p := Plugin(context.Background())
	for name, table := range p.TableMap {
		if table.Name != name {
			t.Errorf("%s registered under %s", table.Name, name)
		}
		if !strings.HasPrefix(name, "ninjaone_") {
			t.Errorf("%s lacks the plugin prefix", name)
		}
		if table.List == nil && table.Get == nil {
			t.Errorf("%s has no list or get", name)
		}
		if !strings.HasSuffix(table.Description, ".") {
			t.Errorf("%s description must end with a period", name)
		}
		seen := map[string]bool{}
		for _, column := range table.Columns {
			if seen[column.Name] {
				t.Errorf("%s repeats column %s", name, column.Name)
			}
			seen[column.Name] = true
			if !strings.HasSuffix(column.Description, ".") {
				t.Errorf("%s.%s description must end with a period", name, column.Name)
			}
		}
		if !seen["device_id"] && !seen["id"] && !seen["uid"] {
			t.Errorf("%s has no identifying column", name)
		}
	}
}

func TestNullIfEmptyKeepsZeroFacts(t *testing.T) {
	zero := int64(0)
	no := false
	cases := []struct {
		in   interface{}
		want interface{}
	}{
		{"", nil},
		{flexString(""), nil},
		{flexString("x"), "x"},
		{json.RawMessage("null"), nil},
		{json.RawMessage(nil), nil},
		{&zero, &zero},
		{&no, &no},
	}
	for _, tc := range cases {
		got, err := nullIfEmpty(context.Background(), &transform.TransformData{Value: tc.in})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("nullIfEmpty(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestEpochToTimeKeepsTheInstant(t *testing.T) {
	seconds := 1793872707.5
	got, err := epochToTime(context.Background(), &transform.TransformData{Value: &seconds})
	if err != nil {
		t.Fatal(err)
	}
	if text := got.(interface{ String() string }).String(); !strings.HasPrefix(text, "2026-11-05") {
		t.Fatalf("got %v", got)
	}
	milliseconds := 1793295720001.0
	got, err = epochToTime(context.Background(), &transform.TransformData{Value: &milliseconds})
	if err != nil {
		t.Fatal(err)
	}
	if text := got.(interface{ String() string }).String(); !strings.HasPrefix(text, "2026-10-") {
		t.Fatalf("milliseconds read as %v", got)
	}
	var missing *float64
	if got, _ := epochToTime(context.Background(), &transform.TransformData{Value: missing}); got != nil {
		t.Fatalf("got %v, want null", got)
	}
}
