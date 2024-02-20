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
