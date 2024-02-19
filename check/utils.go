package check

import (
	"reflect"

	gp "github.com/dan-kuroto/gin-stronger/go-print"
)

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type orderable interface {
	number | rune | string
}

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
