package flagga

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestEnvPrefix(t *testing.T) {
	source := EnvPrefix("FOO_")

	os.Setenv("FOO_BAR", "foo")
	os.Setenv("BAR", "bar")

	testCases := []struct {
		key      string
		ok       bool
		expected string
	}{
		{"bar", false, ""},
		{"BAR", true, "foo"},
		{"QUX", false, ""},
	}

	for _, tt := range testCases {
		t.Run(tt.key, func(t *testing.T) {
			var s string
			ok, err := source.Get(tt.key, NewValue(&s))
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if ok != tt.ok {
				t.Errorf("expected ok to be: %v, got: %v", tt.ok, ok)
			}

			if tt.ok {
				if s != tt.expected {
					t.Errorf("expected value to be: %s, got: %s", tt.expected, s)
				}
			}
		})
	}
}

func TestJSONVia(t *testing.T) {
	data, err := json.Marshal(map[string]interface{}{
		"foo": "bar",
		"bar": 1,
		"baz": []interface{}{3, 1, "5"},
	})
	if err != nil {
		t.Fatalf("unexpected error encoding json: %s", err)
	}

	f, err := ioutil.TempFile(os.TempDir(), "json-test-flagga")
	if err != nil {
		t.Fatalf("unexpected error saving json file: %s", err)
	}

	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Errorf("error removing file: %s", err)
		}
	}()

	if _, err := io.Copy(f, bytes.NewBuffer(data)); err != nil {
		t.Fatalf("unexpected error copying json: %s", err)
	}

	if err := f.Close(); err != nil {
		t.Errorf("error closing file: %s", err)
	}

	source := JSONVia(f.Name())
	if err := source.Open(); err != nil {
		t.Fatalf("unable to open json file: %s", err)
	}

	testCases := []struct {
		dst      interface{}
		key      string
		expected interface{}
		err      bool
		ok       bool
	}{
		{new(string), "qux", nil, false, false},
		{new(string), "foo", "bar", false, true},
		{new(int), "foo", nil, true, false},
		{new([]int), "baz", []int{3, 1, 5}, false, true},
	}

	for _, tt := range testCases {
		t.Run(tt.key, func(t *testing.T) {
			ok, err := source.Get(tt.key, NewValue(tt.dst))
			if tt.err && err == nil {
				t.Errorf("expecting error, got nil instead")
			} else if !tt.err && err != nil {
				t.Errorf("got unexpected error: %s", err)
			}

			if tt.ok != ok {
				t.Errorf("expected ok to be: %v, got: %v", tt.ok, ok)
			}

			if tt.ok {
				val := reflect.ValueOf(tt.dst).Elem().Interface()
				if !reflect.DeepEqual(val, tt.expected) {
					t.Errorf("expecting value to be: %v, got: %v", tt.expected, val)
				}
			}
		})
	}
}
