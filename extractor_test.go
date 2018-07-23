package flagga

import (
	"os"
	"reflect"
	"testing"
)

func TestEnv(t *testing.T) {
	testCases := []struct {
		key      string
		ok       bool
		expected string
	}{
		{"foo", false, ""},
		{"bar", false, ""},
		{"FOO", true, "bar"},
	}

	os.Setenv("FOO", "baz")
	os.Setenv("TEST_FOO", "bar")

	sources := []Source{EnvPrefix("TEST_")}
	for _, tt := range testCases {
		t.Run(tt.key, func(t *testing.T) {
			var s string
			ok, err := Env(tt.key).Get(sources, NewValue(&s))
			if err != nil {
				t.Error(err)
			}

			if ok != tt.ok {
				t.Errorf("expecting ok to be: %v, got: %v", tt.ok, ok)
			}

			if tt.ok && tt.expected != s {
				t.Errorf("expecting result to be %q, got: %q", tt.expected, s)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	testCases := []struct {
		key      string
		err      bool
		ok       bool
		value    interface{}
		expected interface{}
	}{
		{"foo", false, false, nil, nil},
		{"bar", false, true, new(int64), int64(42)},
		{"baz", true, false, new(bool), nil},
	}

	sources := []Source{
		&jsonSource{&FileSource{Value: map[string]interface{}{
			"bar": int64(42),
			"baz": float64(3.14),
		}}},
	}

	for _, tt := range testCases {
		t.Run(tt.key, func(t *testing.T) {
			ok, err := JSON(tt.key).Get(sources, NewValue(tt.value))
			if tt.err && err == nil {
				t.Errorf("expecting error, got nil instead")
			} else if !tt.err && err != nil {
				t.Errorf("got unexpected error: %s", err)
			}

			if tt.ok != ok {
				t.Errorf("expected ok to be: %v, got: %v", tt.ok, ok)
			}

			if tt.ok {
				val := reflect.ValueOf(tt.value).Elem().Interface()
				if !reflect.DeepEqual(val, tt.expected) {
					t.Errorf("expecting value to be: %v, got: %v", tt.expected, val)
				}
			}
		})
	}
}
