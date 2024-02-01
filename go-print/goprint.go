package goprint

import (
	"fmt"
	"reflect"
	"strings"
)

type Formatter struct {
	// indent for array&list
	ListIndent int
	// indent for map
	MapIndent int
	// indent for struct
	StructIndent int
	// Whether show array&slice as a tag.
	//
	// If true, example: <[]int :len=2 items=[1, 2]>
	// Otherwise, example: [1, 2, "4"]
	ListShowAsTag bool
	// Whether show map as a tag.
	//
	// If true, example: <map[string, any] :len=3 "a"=1>
	// Otherwise, example: {"a": 1}
	MapShowAsTag bool
	// Max number for display items in array&list.
	// If `ListDisplayNum <= 0`, means infinity.
	// If `len(data) > ListDisplayNum`, extra parts are shown as ellipsis.
	ListDisplayNum int
	// Max number for display items in map.
	// If `MapDisplayNum <= 0`, means infinity.
	// If `len(data) > MapDisplayNum`, extra parts are shown as ellipsis.
	MapDisplayNum int
	// Only when `StructIndent > 0`, `StructMethodShow` is valid.
	StructMethodShow bool
}

var Default = Formatter{
	StructIndent: 2, ListDisplayNum: 100, MapDisplayNum: 100,
}

func ToString(data any) string {
	return Default.ToString(data)
}

func (f *Formatter) ToString(data any) string {
	return f.toString(data, 0)
}

func (f *Formatter) toString(data any, layer int) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf("%q", data)
	case reflect.Pointer:
		if value.IsNil() {
			return fmt.Sprintf("<%s nil>", f.typeToString(value.Type()))
		}
		if !value.CanInterface() {
			return fmt.Sprintf("%#v", data)
		}
		return "&" + f.toString(value.Elem().Interface(), layer)
	case reflect.Array, reflect.Slice:
		return f.listToString(value, layer)
	case reflect.Map:
		return f.mapToString(value, layer)
	case reflect.Chan:
		return fmt.Sprintf("<%s len=%d cap=%d ptr=%#x>", f.typeToString(value.Type()), value.Len(), value.Cap(), value.Pointer())
	case reflect.Struct:
		return f.structToString(value, layer)
	case reflect.Func:
		return f.funcToString(value)
	default:
		return fmt.Sprintf("%#v", data)
	}
}

func (f *Formatter) listToString(value reflect.Value, layer int) string {
	var sb strings.Builder
	length := value.Len()
	displayLength := length
	if f.ListDisplayNum > 0 && length > f.ListDisplayNum {
		displayLength = f.ListDisplayNum
	}

	if f.ListShowAsTag {
		sb.WriteString(fmt.Sprintf("<%s :len=%d items=", f.typeToString(value.Type()), value.Len()))
	}
	sb.WriteString("[")
	if displayLength > 0 {
		if f.ListIndent > 0 {
			layer++
		}
		appendIndent(&sb, f.ListIndent, layer, false)
		for i := 0; i < displayLength; i++ {
			sb.WriteString(f.toString(value.Index(i).Interface(), layer))
			if i < displayLength-1 { // before the last one
				sb.WriteString(",")
				appendIndent(&sb, f.ListIndent, layer, true)
			} else if displayLength < length { // last one, but need ellipsis
				sb.WriteString(",")
				appendIndent(&sb, f.ListIndent, layer, true)
				sb.WriteString("...")
			}
		}
		if f.ListIndent > 0 {
			layer--
		}
		appendIndent(&sb, f.ListIndent, layer, false)
	}
	sb.WriteString("]")
	if f.ListShowAsTag {
		sb.WriteString(">")
	}

	return sb.String()
}

func (f *Formatter) mapToString(value reflect.Value, layer int) string {
	var sb strings.Builder
	length := value.Len()
	displayLength := length
	if f.MapDisplayNum > 0 && length > f.MapDisplayNum {
		displayLength = f.MapDisplayNum
	}

	if f.MapShowAsTag {
		sb.WriteString(fmt.Sprintf("<%s :len=%d", f.typeToString(value.Type()), value.Len()))
	} else {
		sb.WriteString("{")
	}
	if displayLength > 0 {
		if f.MapIndent > 0 {
			layer++
		}
		appendIndent(&sb, f.MapIndent, layer, false)
		for i, key := range value.MapKeys() {
			sb.WriteString(f.toString(key.Interface(), layer))
			if f.MapShowAsTag {
				sb.WriteString("=")
			} else {
				sb.WriteString(": ")
			}
			sb.WriteString(f.toString(value.MapIndex(key).Interface(), layer))
			if i < displayLength-1 { // before the last one
				if !f.MapShowAsTag {
					sb.WriteString(",")
				}
				appendIndent(&sb, f.MapIndent, layer, true)
			} else { // last one
				if displayLength < length { // need ellipsis
					if !f.MapShowAsTag {
						sb.WriteString(",")
					}
					appendIndent(&sb, f.MapIndent, layer, true)
					sb.WriteString("...")
				}
				break
			}
		}
		if f.MapIndent > 0 {
			layer--
		}
		appendIndent(&sb, f.MapIndent, layer, false)
	}
	if f.MapShowAsTag {
		sb.WriteString(">")
	} else {
		sb.WriteString("}")
	}

	return sb.String()
}

func (f *Formatter) typeToString(type_ reflect.Type) string {
	switch type_.Kind() {
	case reflect.Pointer:
		return "*" + f.typeToString(type_.Elem())
	case reflect.Array, reflect.Slice:
		return fmt.Sprintf("[]%s", f.typeToString(type_.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s, %s]", f.typeToString(type_.Key()), f.typeToString(type_.Elem()))
	case reflect.Chan:
		return fmt.Sprintf("chan[%s]", f.typeToString(type_.Elem()))
	case reflect.Func:
		return fmt.Sprintf("func[%s -> %s]", f.funcInTypeString(type_), f.funcOutTypeString(type_))
	default:
		if type_.Name() == "" {
			return "?"
		} else {
			return type_.String()
		}
	}
}

func (f *Formatter) funcInTypeString(type_ reflect.Type) string {
	numIn := type_.NumIn()
	if numIn == 0 {
		return "()"
	}
	if numIn == 1 {
		return f.typeToString(type_.In(0))
	}
	var sb strings.Builder
	sb.WriteString("(")
	for i := 0; i < numIn; i++ {
		sb.WriteString(f.typeToString(type_.In(i)))
		if i < numIn-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (f *Formatter) funcOutTypeString(type_ reflect.Type) string {
	numOut := type_.NumOut()
	if numOut == 0 {
		return "()"
	}
	if numOut == 1 {
		return f.typeToString(type_.Out(0))
	}
	var sb strings.Builder
	sb.WriteString("(")
	for i := 0; i < numOut; i++ {
		sb.WriteString(f.typeToString(type_.Out(i)))
		if i < numOut-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (f *Formatter) structToString(value reflect.Value, layer int) string {
	type_ := value.Type()
	var fields []reflect.StructField
	for i := 0; i < value.NumField(); i++ {
		field := type_.Field(i)
		if field.IsExported() {
			fields = append(fields, field)
		}
	}
	length := len(fields)

	var sb strings.Builder

	sb.WriteString("<")
	sb.WriteString(f.typeToString(type_))
	if length > 0 {
		if f.StructIndent > 0 {
			layer++
		}
		appendIndent(&sb, f.StructIndent, layer, true)
		for i, field := range fields {
			sb.WriteString(field.Name)
			sb.WriteString("=")
			sb.WriteString(f.toString(value.Field(i).Interface(), layer))
			if i < length-1 {
				appendIndent(&sb, f.StructIndent, layer, true)
			}
		}
		if f.StructIndent > 0 {
			layer--
		}
		appendIndent(&sb, f.StructIndent, layer, false)
	}
	sb.WriteString(">")

	return sb.String()
}

func (f *Formatter) funcToString(value reflect.Value) string {
	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(f.typeToString(value.Type()))
	if value.IsNil() {
		sb.WriteString(" nil")
	} else {
		sb.WriteString(fmt.Sprintf(" ptr=%#x", value.Pointer()))
	}
	sb.WriteString(">")
	return sb.String()
}
