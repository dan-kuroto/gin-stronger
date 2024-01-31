package goprint

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
	case reflect.Chan:
		return fmt.Sprintf("<%s len=%d cap=%d ptr=%#x>", typeToString(value.Type()), value.Len(), value.Cap(), value.Pointer())
	case reflect.Struct:
		return structToString(value)
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
		if type_.Name() == "" {
			return "?"
		} else {
			return type_.String()
		}
	}
}

func structToString(value reflect.Value) string {
	type_ := value.Type()

	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(typeToString(type_))
	for i := 0; i < value.NumField(); i++ {
		field := type_.Field(i)
		if field.IsExported() {
			sb.WriteString(" ")
			sb.WriteString(field.Name)
			sb.WriteString("=")
			sb.WriteString(ToString(value.Field(i).Interface()))
		}
	}
	sb.WriteString(">")
	// TODO: 显示method?(允许换行的时候才显示,现在还不支持indent)

	return sb.String()
}

// TODO: config: 是否换行indent/pointer是否&/map等容器是否带type/...(可参考objprint)
