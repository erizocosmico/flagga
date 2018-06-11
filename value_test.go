package flagga

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestValue(t *testing.T) {
	testCases := []struct {
		dst      interface{}
		value    interface{}
		expected interface{}
		err      bool
	}{
		{new(string), "foo", "foo", false},
		{new(string), []byte("foo"), "foo", false},
		{new(string), int64(1), "1", false},

		{new(bool), true, true, false},
		{new(bool), false, false, false},
		{new(bool), "true", true, false},
		{new(bool), "false", false, false},
		{new(bool), []byte("true"), true, false},
		{new(bool), 1, false, true},
		{new(bool), "asldkjsa", false, true},

		{new(float64), 3.14, 3.14, false},
		{new(float64), float32(3.), 3., false},
		{new(float64), "3.14", 3.14, false},
		{new(float64), []byte("3.14"), 3.14, false},
		{new(float64), "asldkjsa", 0, true},
		{new(float64), true, 0, true},
		{new(float64), int(1), float64(1), false},
		{new(float64), uint(1), float64(1), false},
		{new(float64), int64(1), float64(1), false},
		{new(float64), uint64(1), float64(1), false},

		{new(int), int(1), int(1), false},
		{new(int), uint(1), int(1), false},
		{new(int), int64(1), int(1), false},
		{new(int), uint64(1), int(1), false},
		{new(int), "1", int(1), false},
		{new(int), []byte("1"), int(1), false},
		{new(int), 3.14, int(3), false},
		{new(int), "asldkjsa", 0, true},

		{new(uint), int(1), uint(1), false},
		{new(uint), uint(1), uint(1), false},
		{new(uint), int64(1), uint(1), false},
		{new(uint), uint64(1), uint(1), false},
		{new(uint), "1", uint(1), false},
		{new(uint), []byte("1"), uint(1), false},
		{new(uint), 3.14, uint(3), false},
		{new(uint), "asldkjsa", 0, true},

		{new(int64), int(1), int64(1), false},
		{new(int64), uint(1), int64(1), false},
		{new(int64), int64(1), int64(1), false},
		{new(int64), uint64(1), int64(1), false},
		{new(int64), "1", int64(1), false},
		{new(int64), []byte("1"), int64(1), false},
		{new(int64), 3.14, int64(3), false},
		{new(int64), "asldkjsa", 0, true},

		{new(uint64), int(1), uint64(1), false},
		{new(uint64), uint(1), uint64(1), false},
		{new(uint64), int64(1), uint64(1), false},
		{new(uint64), uint64(1), uint64(1), false},
		{new(uint64), "1", uint64(1), false},
		{new(uint64), []byte("1"), uint64(1), false},
		{new(uint64), 3.14, uint64(3), false},
		{new(uint64), "asldkjsa", 0, true},

		{new(time.Duration), int(1), time.Duration(1), false},
		{new(time.Duration), uint(1), time.Duration(1), false},
		{new(time.Duration), int64(1), time.Duration(1), false},
		{new(time.Duration), uint64(1), time.Duration(1), false},
		{new(time.Duration), 1 * time.Second, 1 * time.Second, false},
		{new(time.Duration), "1s", 1 * time.Second, false},
		{new(time.Duration), []byte("1s"), 1 * time.Second, false},
		{new(time.Duration), 3.14, time.Duration(3), false},
		{new(time.Duration), "1", 0, true},

		{new([]string), "foo", []string{"foo"}, false},
		{new([]string), []byte("foo"), []string{"foo"}, false},
		{new([]string), []string{"f", "o", "o"}, []string{"f", "o", "o"}, false},
		{new([]string), []interface{}{"f", 1, 3.14}, []string{"f", "1", "3.14"}, false},

		{new([]int), []int{1, 2}, []int{1, 2}, false},
		{new([]int), []uint{1, 2}, []int{1, 2}, false},
		{new([]int), []int64{1, 2}, []int{1, 2}, false},
		{new([]int), []uint64{1, 2}, []int{1, 2}, false},
		{new([]int), []float64{1, 2}, []int{1, 2}, false},
		{new([]int), []string{"1", "2"}, []int{1, 2}, false},
		{new([]int), []string{"a", "2"}, nil, true},
		{new([]int), []interface{}{0, 1}, []int{0, 1}, false},
		{new([]int), "1", []int{1}, false},
		{new([]int), []byte("1"), []int{1}, false},
		{new([]int), int(1), []int{1}, false},
		{new([]int), int64(1), []int{1}, false},
		{new([]int), uint(1), []int{1}, false},
		{new([]int), uint64(1), []int{1}, false},
		{new([]int), float64(1), []int{1}, false},
		{new([]int), []interface{}{"a", 1}, nil, true},
		{new([]int), "a", nil, true},

		{new([]uint), []int{1, 2}, []uint{1, 2}, false},
		{new([]uint), []uint{1, 2}, []uint{1, 2}, false},
		{new([]uint), []int64{1, 2}, []uint{1, 2}, false},
		{new([]uint), []uint64{1, 2}, []uint{1, 2}, false},
		{new([]uint), []float64{1, 2}, []uint{1, 2}, false},
		{new([]uint), []string{"1", "2"}, []uint{1, 2}, false},
		{new([]uint), []string{"a", "2"}, nil, true},
		{new([]uint), []interface{}{0, 1}, []uint{0, 1}, false},
		{new([]uint), "1", []uint{1}, false},
		{new([]uint), []byte("1"), []uint{1}, false},
		{new([]uint), uint(1), []uint{1}, false},
		{new([]uint), int64(1), []uint{1}, false},
		{new([]uint), uint(1), []uint{1}, false},
		{new([]uint), uint64(1), []uint{1}, false},
		{new([]uint), float64(1), []uint{1}, false},
		{new([]uint), []interface{}{"a", 1}, nil, true},
		{new([]uint), "a", nil, true},

		{new([]int64), []int{1, 2}, []int64{1, 2}, false},
		{new([]int64), []uint{1, 2}, []int64{1, 2}, false},
		{new([]int64), []int64{1, 2}, []int64{1, 2}, false},
		{new([]int64), []uint64{1, 2}, []int64{1, 2}, false},
		{new([]int64), []float64{1, 2}, []int64{1, 2}, false},
		{new([]int64), []string{"1", "2"}, []int64{1, 2}, false},
		{new([]int64), []string{"a", "2"}, nil, true},
		{new([]int64), []interface{}{0, 1}, []int64{0, 1}, false},
		{new([]int64), "1", []int64{1}, false},
		{new([]int64), []byte("1"), []int64{1}, false},
		{new([]int64), int64(1), []int64{1}, false},
		{new([]int64), int64(1), []int64{1}, false},
		{new([]int64), uint(1), []int64{1}, false},
		{new([]int64), uint64(1), []int64{1}, false},
		{new([]int64), float64(1), []int64{1}, false},
		{new([]int64), []interface{}{"a", 1}, nil, true},
		{new([]int64), "a", nil, true},

		{new([]uint64), []int{1, 2}, []uint64{1, 2}, false},
		{new([]uint64), []uint{1, 2}, []uint64{1, 2}, false},
		{new([]uint64), []int64{1, 2}, []uint64{1, 2}, false},
		{new([]uint64), []uint64{1, 2}, []uint64{1, 2}, false},
		{new([]uint64), []float64{1, 2}, []uint64{1, 2}, false},
		{new([]uint64), []string{"1", "2"}, []uint64{1, 2}, false},
		{new([]uint64), []string{"a", "2"}, nil, true},
		{new([]uint64), []interface{}{0, 1}, []uint64{0, 1}, false},
		{new([]uint64), "1", []uint64{1}, false},
		{new([]uint64), []byte("1"), []uint64{1}, false},
		{new([]uint64), uint64(1), []uint64{1}, false},
		{new([]uint64), int64(1), []uint64{1}, false},
		{new([]uint64), uint(1), []uint64{1}, false},
		{new([]uint64), uint64(1), []uint64{1}, false},
		{new([]uint64), float64(1), []uint64{1}, false},
		{new([]uint64), []interface{}{"a", 1}, nil, true},
		{new([]uint64), "a", nil, true},

		{new([]time.Duration), []time.Duration{1, 2}, []time.Duration{1, 2}, false},
		{new([]time.Duration), []int{1, 2}, []time.Duration{1, 2}, false},
		{new([]time.Duration), []uint{1, 2}, []time.Duration{1, 2}, false},
		{new([]time.Duration), []int64{1, 2}, []time.Duration{1, 2}, false},
		{new([]time.Duration), []uint64{1, 2}, []time.Duration{1, 2}, false},
		{new([]time.Duration), []float64{1, 2}, []time.Duration{1, 2}, false},
		{new([]time.Duration), []string{"1s", "2s"}, []time.Duration{1 * time.Second, 2 * time.Second}, false},
		{new([]time.Duration), []string{"a", "2"}, nil, true},
		{new([]time.Duration), []interface{}{0, 1}, []time.Duration{0, 1}, false},
		{new([]time.Duration), "1s", []time.Duration{1 * time.Second}, false},
		{new([]time.Duration), []byte("1s"), []time.Duration{1 * time.Second}, false},
		{new([]time.Duration), time.Duration(1), []time.Duration{1}, false},
		{new([]time.Duration), int(1), []time.Duration{1}, false},
		{new([]time.Duration), int64(1), []time.Duration{1}, false},
		{new([]time.Duration), uint(1), []time.Duration{1}, false},
		{new([]time.Duration), uint64(1), []time.Duration{1}, false},
		{new([]time.Duration), float64(1), []time.Duration{1}, false},
		{new([]time.Duration), []interface{}{"a", 1}, nil, true},
		{new([]time.Duration), "a", nil, true},

		{new([]float64), []float64{1., 2.}, []float64{1., 2.}, false},
		{new([]float64), []float32{1., 2.}, []float64{1., 2.}, false},
		{new([]float64), []string{"1.0", "2.0"}, []float64{1., 2.}, false},
		{new([]float64), []string{"a", "2."}, nil, true},
		{new([]float64), []interface{}{"a", "2."}, nil, true},
		{new([]float64), []interface{}{2.0, 3.0}, []float64{2., 3.}, false},
		{new([]float64), "1.0", []float64{1}, false},
		{new([]float64), "a", nil, true},
		{new([]float64), []byte("1.0"), []float64{1}, false},
		{new([]float64), float64(1), []float64{1.}, false},
		{new([]float64), float32(1), []float64{1.}, false},
		{new([]float64), []int{1, 2}, []float64{1, 2}, false},
		{new([]float64), []uint{1, 2}, []float64{1, 2}, false},
		{new([]float64), []int64{1, 2}, []float64{1, 2}, false},
		{new([]float64), []uint64{1, 2}, []float64{1, 2}, false},
		{new([]float64), uint64(1), []float64{1}, false},
		{new([]float64), int64(1), []float64{1}, false},
		{new([]float64), uint(1), []float64{1}, false},
		{new([]float64), uint64(1), []float64{1}, false},
	}

	for _, tt := range testCases {
		name := fmt.Sprintf("set (%T)(%v) to %T", tt.value, tt.value, tt.dst)
		t.Run(name, func(t *testing.T) {
			err := NewValue(tt.dst).Set(tt.value)
			if tt.err && err == nil {
				t.Errorf("expecting error, got nil instead")
			} else if !tt.err && err != nil {
				t.Errorf("got unexpected error: %s", err)
			}

			if !tt.err {
				val := reflect.ValueOf(tt.dst).Elem().Interface()
				if !reflect.DeepEqual(val, tt.expected) {
					t.Errorf(
						"expecting value to be: (%T)(%v), got: (%T)(%v)",
						tt.expected,
						tt.expected,
						val,
						val,
					)
				}
			}
		})
	}
}
