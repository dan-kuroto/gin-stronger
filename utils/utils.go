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
		if value.IsNil() {
			return fmt.Sprintf("<%s nil>", typeToString(value.Type()))
		}
		if !value.CanInterface() {
			return fmt.Sprintf("%#v", data)
		}
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
	case reflect.Interface, reflect.Struct:
		// TODO: 类似前面map的处理，但是要加上类名，且key不加引号了
		return fmt.Sprintf("%#v", data)
	case reflect.Chan:
		return fmt.Sprintf("<%s len=%d cap=%d ptr=%#x>", typeToString(value.Type()), value.Len(), value.Cap(), value.Pointer())
	case reflect.Func:
		// TODO
		return fmt.Sprintf("%#v", data)
	default:
		return fmt.Sprintf("%#v", data)
	}
}

func typeToString(type_ reflect.Type) string {
	switch type_.Kind() {
	case reflect.Pointer:
		return "*" + typeToString(type_.Elem())
	case reflect.Array, reflect.Slice:
		return fmt.Sprintf("[]%s", typeToString(type_.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s, %s]", typeToString(type_.Key()), typeToString(type_.Elem()))
	case reflect.Chan:
		return fmt.Sprintf("chan[%s]", typeToString(type_.Elem()))
	default:
		return type_.String()
	}
}

// TODO: config: 是否换行indent/pointer是否&/...(可参考objprint)
