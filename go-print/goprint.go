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
	// Pointer preffix: such as '&' in `&time.Time{}` and '*' in `*time.Time`.
	PointerPreffixHide bool
	// Whether show container(map/array/slice) as a tag.
	//
	// If true, example: <map[string, any] :len=3 "a"=1 "b"=<[]int :len=2 items=[1, 2]>>
	//
	// Otherwise, example: {"a": 1, "b": {"c": [1, 2, "4"]}}
	ContainerShowAsTag bool
	// If `ContainerDisplayNum <= 0`, means infinity.
	// If `len(data) > ContainerDisplayNum`, extra parts are shown as ellipsis.
	ContainerDisplayNum int
	// Only when `StructIndent > 0`, `StructMethodShow` is valid.
	StructMethodShow bool
	Color            bool
}

var Default = Formatter{
	StructIndent: 2,
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

	sb.WriteString("[")
	if length > 0 {
		if f.ListIndent > 0 {
			layer++
		}
		appendIndent(&sb, f.ListIndent, layer, false)
		for i := 0; i < length; i++ {
			sb.WriteString(f.toString(value.Index(i).Interface(), layer))
			if i < length-1 {
				sb.WriteString(",")
				appendIndent(&sb, f.ListIndent, layer, true)
			}
		}
		if f.ListIndent > 0 {
			layer--
		}
		appendIndent(&sb, f.ListIndent, layer, false)
	}
	sb.WriteString("]")

	return sb.String()
}

func (f *Formatter) mapToString(value reflect.Value, layer int) string {
	var sb strings.Builder
	length := value.Len()

	sb.WriteString("{")
	if length > 0 {
		if f.MapIndent > 0 {
			layer++
		}
		appendIndent(&sb, f.MapIndent, layer, false)
		for i, key := range value.MapKeys() {
			sb.WriteString(f.toString(key.Interface(), layer))
			sb.WriteString(": ")
			sb.WriteString(f.toString(value.MapIndex(key).Interface(), layer))
			if i < length-1 {
				sb.WriteString(", ")
				appendIndent(&sb, f.MapIndent, layer, true)
			}
		}
		if f.MapIndent > 0 {
			layer--
		}
		appendIndent(&sb, f.MapIndent, layer, false)
	}
	sb.WriteString("}")

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
	var sb strings.Builder
	type_ := value.Type()
	length := value.NumField()

	sb.WriteString("<")
	sb.WriteString(f.typeToString(type_))
	if length > 0 {
		if f.StructIndent > 0 {
			layer++
		}
		appendIndent(&sb, f.StructIndent, layer, true)
		for i := 0; i < length; i++ {
			field := type_.Field(i)
			if field.IsExported() {
				sb.WriteString(field.Name)
				sb.WriteString("=")
				sb.WriteString(f.toString(value.Field(i).Interface(), layer))
			}
			if i < length-1 { // TODO: 这里有问题，应该直接遍历IsExported的,否则就会多出一些空格或换行出来
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
