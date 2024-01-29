package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func ToString(data any) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf("%q", data)
	case reflect.Pointer:
		return "&" + ToString(value.Elem().Interface())
	case reflect.Array, reflect.Slice:
		var sb strings.Builder
		sb.WriteString("[")
		for i := 0; i < value.Len(); i++ {
			sb.WriteString(ToString(value.Index(i).Interface()))
			if i < value.Len()-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("]")
		return sb.String()
	case reflect.Map:
		var sb strings.Builder
		sb.WriteString("{")
		for i, key := range value.MapKeys() {
			sb.WriteString(ToString(key.Interface()))
			sb.WriteString(": ")
			sb.WriteString(ToString(value.MapIndex(key).Interface()))
			if i < value.Len()-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("}")
		return sb.String()
	// TODO: 1. interface&struct 2. config: 是否换行indent/pointer是否&/...(可参考objprint)
	default:
		return fmt.Sprint(data)
	}
}
