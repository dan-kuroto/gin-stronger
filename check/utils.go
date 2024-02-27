package check

import (
	"reflect"

	gp "github.com/dan-kuroto/gin-stronger/go-print"
)

var formatter = gp.Formatter{}

// get length of string/array/chan/map/slice or pointer to them
func getLength(data any) (length int, ok bool) {
	if str, ok := toString(data); ok {
		return len(str), true
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				length, ok = 0, false
			}
		}()

		value := reflect.ValueOf(data)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		length, ok = value.Len(), true
	}()

	return length, ok
}

// convert to float64 from int/int8/.../uint/uint8/.../float32/float64 or their pointer
func toFloat64(data any) (float64, bool) {
	switch data := data.(type) {
	case int:
		return float64(data), true
	case *int:
		return float64(*data), true
	case int8:
		return float64(data), true
	case *int8:
		return float64(*data), true
	case int16:
		return float64(data), true
	case *int16:
		return float64(*data), true
	case int32:
		return float64(data), true
	case *int32:
		return float64(*data), true
	case int64:
		return float64(data), true
	case *int64:
		return float64(*data), true
	case uint:
		return float64(data), true
	case *uint:
		return float64(*data), true
	case uint8:
		return float64(data), true
	case *uint8:
		return float64(*data), true
	case uint16:
		return float64(data), true
	case *uint16:
		return float64(*data), true
	case uint32:
		return float64(data), true
	case *uint32:
		return float64(*data), true
	case uint64:
		return float64(data), true
	case *uint64:
		return float64(*data), true
	case float32:
		return float64(data), true
	case *float32:
		return float64(*data), true
	case float64:
		return data, true
	case *float64:
		return *data, true
	default:
		return 0, false
	}
}

// convert string/*string to string
func toString(data any) (string, bool) {
	switch data := data.(type) {
	case string:
		return data, true
	case *string:
		return *data, true
	default:
		return "", false
	}
}

// convert bool/*bool to string
func toBool(data any) (value bool, ok bool) {
	switch data := data.(type) {
	case bool:
		return data, true
	case *bool:
		return *data, true
	default:
		return false, false
	}
}

// Only valid for int/int8/.../uint/uint8/.../float32/float64/string and their pointer.
func basicEqual(a, b any) (equal bool, ok bool) {
	// int/int8/.../uint/uint8/.../float32/float64
	if a, ok := toFloat64(a); ok {
		if b, ok := toFloat64(b); ok {
			return a == b, true
		} else {
			return false, false
		}
	}

	// string
	if a, ok := toString(a); ok {
		if b, ok := toString(b); ok {
			return a == b, true
		} else {
			return false, false
		}
	}

	return false, false
}

// The type handling mechanism is the same as `basicEqual`.
func basicGreater(a, b any) (greater bool, ok bool) {
	// int/int8/.../uint/uint8/.../float32/float64
	if a, ok := toFloat64(a); ok {
		if b, ok := toFloat64(b); ok {
			return a > b, true
		} else {
			return false, false
		}
	}

	// string
	if a, ok := toString(a); ok {
		if b, ok := toString(b); ok {
			return a > b, true
		} else {
			return false, false
		}
	}

	return false, false
}

// The type handling mechanism is the same as `basicEqual`.
func basicLess(a, b any) (less bool, ok bool) {
	// int/int8/.../uint/uint8/.../float32/float64
	if a, ok := toFloat64(a); ok {
		if b, ok := toFloat64(b); ok {
			return a < b, true
		} else {
			return false, false
		}
	}

	// string
	if a, ok := toString(a); ok {
		if b, ok := toString(b); ok {
			return a < b, true
		} else {
			return false, false
		}
	}

	return false, false
}

// The type handling mechanism is the same as `basicEqual`.
func basicIn(a any, b []any) bool {
	for _, item := range b {
		if equal, ok := basicEqual(a, item); ok && equal {
			return true
		}
	}
	return false
}
