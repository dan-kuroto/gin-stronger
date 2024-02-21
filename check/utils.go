package check

import (
	"reflect"

	gp "github.com/dan-kuroto/gin-stronger/go-print"
)

var formatter = gp.Formatter{}

func getLength(data any) (length int, ok bool) {
	if str, ok := data.(string); ok {
		return len(str), true
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				length, ok = 0, false
			}
		}()

		value := reflect.ValueOf(data)
		length, ok = value.Len(), true
	}()

	return length, ok
}

func toFloat64(data any) (float64, bool) {
	switch data := data.(type) {
	case int:
		return float64(data), true
	case int8:
		return float64(data), true
	case int16:
		return float64(data), true
	case int32:
		return float64(data), true
	case int64:
		return float64(data), true
	case uint:
		return float64(data), true
	case uint8:
		return float64(data), true
	case uint16:
		return float64(data), true
	case uint32:
		return float64(data), true
	case uint64:
		return float64(data), true
	case float32:
		return float64(data), true
	case float64:
		return data, true
	default:
		return 0, false
	}
}
