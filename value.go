package flagga

import (
	"fmt"
	"strconv"
	"time"
)

// Value is a flag value that can be set.
type Value interface {
	// Set the given value as the new value. For slice types, this performs an
	// append if the value is not a slice.
	Set(val interface{}) error
}

type value struct {
	value interface{}
}

// NewValue wraps the pointer into a Value type.
func NewValue(val interface{}) Value {
	return &value{val}
}

func (v *value) Set(val interface{}) error {
	switch v := v.value.(type) {
	case *string:
		assignString(v, val)
		return nil
	case *float64:
		return assignFloat64(v, val)
	case *bool:
		return assignBool(v, val)
	case *uint:
		return assignUint(v, val)
	case *int:
		return assignInt(v, val)
	case *uint64:
		return assignUint64(v, val)
	case *int64:
		return assignInt64(v, val)
	case *time.Duration:
		return assignDuration(v, val)
	case *[]string:
		assignStringList(v, val)
		return nil
	case *[]float64:
		return assignFloat64List(v, val)
	case *[]int:
		return assignIntList(v, val)
	case *[]uint:
		return assignUintList(v, val)
	case *[]int64:
		return assignInt64List(v, val)
	case *[]uint64:
		return assignUint64List(v, val)
	case *[]time.Duration:
		return assignDurationList(v, val)
	}

	panic(fmt.Errorf("invalid value of type: %T", v.value))
}

func assignString(dst *string, val interface{}) {
	switch val := val.(type) {
	case string:
		*dst = val
	case []byte:
		*dst = string(val)
	default:
		*dst = fmt.Sprint(val)
	}
}

func assignFloat64(dst *float64, val interface{}) error {
	switch val := val.(type) {
	case float64:
		*dst = val
	case float32:
		*dst = float64(val)
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		*dst = f
	case int:
		*dst = float64(val)
	case int64:
		*dst = float64(val)
	case uint:
		*dst = float64(val)
	case uint64:
		*dst = float64(val)
	case []byte:
		return assignFloat64(dst, string(val))
	default:
		return fmt.Errorf("cannot assign type %T to float64", val)
	}

	return nil
}

func assignInt(dst *int, val interface{}) error {
	switch val := val.(type) {
	case int:
		*dst = val
	case int64:
		*dst = int(val)
	case uint:
		*dst = int(val)
	case uint64:
		*dst = int(val)
	case float64:
		*dst = int(val)
	case string:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		*dst = int(n)
	case []byte:
		return assignInt(dst, string(val))
	default:
		return fmt.Errorf("cannot assign type %T to int", val)
	}

	return nil
}

func assignUint(dst *uint, val interface{}) error {
	switch val := val.(type) {
	case uint:
		*dst = val
	case int64:
		*dst = uint(val)
	case int:
		*dst = uint(val)
	case uint64:
		*dst = uint(val)
	case float64:
		*dst = uint(val)
	case string:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		*dst = uint(n)
	case []byte:
		return assignUint(dst, string(val))
	default:
		return fmt.Errorf("cannot assign type %T to uint", val)
	}

	return nil
}

func assignInt64(dst *int64, val interface{}) error {
	switch val := val.(type) {
	case int:
		*dst = int64(val)
	case int64:
		*dst = val
	case uint:
		*dst = int64(val)
	case uint64:
		*dst = int64(val)
	case float64:
		*dst = int64(val)
	case string:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		*dst = n
	case []byte:
		return assignInt64(dst, string(val))
	default:
		return fmt.Errorf("cannot assign type %T to int", val)
	}

	return nil
}

func assignUint64(dst *uint64, val interface{}) error {
	switch val := val.(type) {
	case uint:
		*dst = uint64(val)
	case int64:
		*dst = uint64(val)
	case int:
		*dst = uint64(val)
	case uint64:
		*dst = val
	case float64:
		*dst = uint64(val)
	case string:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		*dst = n
	case []byte:
		return assignUint64(dst, string(val))
	default:
		return fmt.Errorf("cannot assign type %T to uint64", val)
	}

	return nil
}

func assignBool(dst *bool, val interface{}) error {
	switch val := val.(type) {
	case bool:
		*dst = val
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		*dst = b
	case []byte:
		return assignBool(dst, string(val))
	default:
		return fmt.Errorf("cannot assign type %T to uint64", val)
	}

	return nil
}

func assignDuration(dst *time.Duration, val interface{}) error {
	switch val := val.(type) {
	case time.Duration:
		*dst = val
	case int:
		*dst = time.Duration(val)
	case int64:
		*dst = time.Duration(val)
	case uint:
		*dst = time.Duration(val)
	case uint64:
		*dst = time.Duration(val)
	case float64:
		*dst = time.Duration(val)
	case string:
		d, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		*dst = d
	case []byte:
		return assignDuration(dst, string(val))
	default:
		return fmt.Errorf("cannot assign type %T to time.Duration", val)
	}

	return nil
}

func assignStringList(dst *[]string, val interface{}) {
	switch val := val.(type) {
	case []interface{}:
		*dst = make([]string, len(val))
		for i, v := range val {
			assignString(&(*dst)[i], v)
		}
	case []string:
		*dst = val
	default:
		var s string
		assignString(&s, val)
		*dst = append(*dst, s)
	}
}

func assignIntList(dst *[]int, val interface{}) error {
	switch val := val.(type) {
	case []interface{}:
		*dst = make([]int, len(val))
		for i, v := range val {
			if err := assignInt(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []string:
		*dst = make([]int, len(val))
		for i, v := range val {
			if err := assignInt(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []int:
		*dst = val
	case []uint:
		*dst = make([]int, len(val))
		for i, v := range val {
			(*dst)[i] = int(v)
		}
	case []int64:
		*dst = make([]int, len(val))
		for i, v := range val {
			(*dst)[i] = int(v)
		}
	case []uint64:
		*dst = make([]int, len(val))
		for i, v := range val {
			(*dst)[i] = int(v)
		}
	case []float64:
		*dst = make([]int, len(val))
		for i, v := range val {
			(*dst)[i] = int(v)
		}
	case []byte, string, int, uint, uint64, int64, float64:
		var v int
		if err := assignInt(&v, val); err != nil {
			return err
		}
		*dst = append(*dst, v)
	default:
		return fmt.Errorf("cannot assign type %T to []int", val)
	}

	return nil
}

func assignUintList(dst *[]uint, val interface{}) error {
	switch val := val.(type) {
	case []interface{}:
		*dst = make([]uint, len(val))
		for i, v := range val {
			if err := assignUint(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []string:
		*dst = make([]uint, len(val))
		for i, v := range val {
			if err := assignUint(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []uint:
		*dst = val
	case []int:
		*dst = make([]uint, len(val))
		for i, v := range val {
			(*dst)[i] = uint(v)
		}
	case []int64:
		*dst = make([]uint, len(val))
		for i, v := range val {
			(*dst)[i] = uint(v)
		}
	case []uint64:
		*dst = make([]uint, len(val))
		for i, v := range val {
			(*dst)[i] = uint(v)
		}
	case []float64:
		*dst = make([]uint, len(val))
		for i, v := range val {
			(*dst)[i] = uint(v)
		}
	case []byte, string, int, uint, uint64, int64, float64:
		var v uint
		if err := assignUint(&v, val); err != nil {
			return err
		}
		*dst = append(*dst, v)
	default:
		return fmt.Errorf("cannot assign type %T to []uint", val)
	}

	return nil
}

func assignInt64List(dst *[]int64, val interface{}) error {
	switch val := val.(type) {
	case []interface{}:
		*dst = make([]int64, len(val))
		for i, v := range val {
			if err := assignInt64(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []string:
		*dst = make([]int64, len(val))
		for i, v := range val {
			if err := assignInt64(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []int64:
		*dst = val
	case []uint:
		*dst = make([]int64, len(val))
		for i, v := range val {
			(*dst)[i] = int64(v)
		}
	case []int:
		*dst = make([]int64, len(val))
		for i, v := range val {
			(*dst)[i] = int64(v)
		}
	case []uint64:
		*dst = make([]int64, len(val))
		for i, v := range val {
			(*dst)[i] = int64(v)
		}
	case []float64:
		*dst = make([]int64, len(val))
		for i, v := range val {
			(*dst)[i] = int64(v)
		}
	case []byte, string, int, uint, uint64, int64, float64:
		var v int64
		if err := assignInt64(&v, val); err != nil {
			return err
		}
		*dst = append(*dst, v)
	default:
		return fmt.Errorf("cannot assign type %T to []int64", val)
	}

	return nil
}

func assignUint64List(dst *[]uint64, val interface{}) error {
	switch val := val.(type) {
	case []interface{}:
		*dst = make([]uint64, len(val))
		for i, v := range val {
			if err := assignUint64(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []string:
		*dst = make([]uint64, len(val))
		for i, v := range val {
			if err := assignUint64(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []uint64:
		*dst = val
	case []int:
		*dst = make([]uint64, len(val))
		for i, v := range val {
			(*dst)[i] = uint64(v)
		}
	case []int64:
		*dst = make([]uint64, len(val))
		for i, v := range val {
			(*dst)[i] = uint64(v)
		}
	case []uint:
		*dst = make([]uint64, len(val))
		for i, v := range val {
			(*dst)[i] = uint64(v)
		}
	case []float64:
		*dst = make([]uint64, len(val))
		for i, v := range val {
			(*dst)[i] = uint64(v)
		}
	case []byte, string, int, uint, uint64, int64, float64:
		var v uint64
		if err := assignUint64(&v, val); err != nil {
			return err
		}
		*dst = append(*dst, v)
	default:
		return fmt.Errorf("cannot assign type %T to []uint64", val)
	}

	return nil
}

func assignDurationList(dst *[]time.Duration, val interface{}) error {
	switch val := val.(type) {
	case []interface{}:
		*dst = make([]time.Duration, len(val))
		for i, v := range val {
			if err := assignDuration(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []string:
		*dst = make([]time.Duration, len(val))
		for i, v := range val {
			if err := assignDuration(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []time.Duration:
		*dst = val
	case []int:
		*dst = make([]time.Duration, len(val))
		for i, v := range val {
			(*dst)[i] = time.Duration(v)
		}
	case []uint:
		*dst = make([]time.Duration, len(val))
		for i, v := range val {
			(*dst)[i] = time.Duration(v)
		}
	case []int64:
		*dst = make([]time.Duration, len(val))
		for i, v := range val {
			(*dst)[i] = time.Duration(v)
		}
	case []uint64:
		*dst = make([]time.Duration, len(val))
		for i, v := range val {
			(*dst)[i] = time.Duration(v)
		}
	case []float64:
		*dst = make([]time.Duration, len(val))
		for i, v := range val {
			(*dst)[i] = time.Duration(v)
		}
	case []byte, string, int, uint, uint64, int64, time.Duration, float64:
		var v time.Duration
		if err := assignDuration(&v, val); err != nil {
			return err
		}
		*dst = append(*dst, v)
	default:
		return fmt.Errorf("cannot assign type %T to []time.Duration", val)
	}

	return nil
}

func assignFloat64List(dst *[]float64, val interface{}) error {
	switch val := val.(type) {
	case []interface{}:
		*dst = make([]float64, len(val))
		for i, v := range val {
			if err := assignFloat64(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []string:
		*dst = make([]float64, len(val))
		for i, v := range val {
			if err := assignFloat64(&(*dst)[i], v); err != nil {
				return err
			}
		}
	case []float64:
		*dst = val
	case []float32:
		*dst = make([]float64, len(val))
		for i, v := range val {
			(*dst)[i] = float64(v)
		}
	case []int:
		*dst = make([]float64, len(val))
		for i, v := range val {
			(*dst)[i] = float64(v)
		}
	case []uint:
		*dst = make([]float64, len(val))
		for i, v := range val {
			(*dst)[i] = float64(v)
		}
	case []int64:
		*dst = make([]float64, len(val))
		for i, v := range val {
			(*dst)[i] = float64(v)
		}
	case []uint64:
		*dst = make([]float64, len(val))
		for i, v := range val {
			(*dst)[i] = float64(v)
		}
	case []byte, string, float32, float64, int, uint, int64, uint64:
		var v float64
		if err := assignFloat64(&v, val); err != nil {
			return err
		}
		*dst = append(*dst, v)
	default:
		return fmt.Errorf("cannot assign type %T to []float64", val)
	}

	return nil
}

func isBool(v Value) bool {
	vb, ok := v.(*value)
	if !ok {
		return false
	}

	_, ok = vb.value.(*bool)
	return ok
}

func isSlice(v Value) bool {
	vb, ok := v.(*value)
	if !ok {
		return false
	}

	switch vb.value.(type) {
	case *[]string,
		*[]float64,
		*[]int,
		*[]uint,
		*[]int64,
		*[]uint64,
		*[]time.Duration:
		return true
	default:
		return false
	}
}
