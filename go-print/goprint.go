package goprint

import (
	"fmt"
	"reflect"
	"strings"
)

type Processor struct {
	// No line breaks when `Indent <= 0`.
	Indent int
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
	// Only when `Indent > 0`, `StructMethodShow` is valid.
	StructMethodShow bool
}

var Default = Processor{
	Indent: 2,
}

func ToString(data any) string {
	return Default.ToString(data)
}

func (p *Processor) ToString(data any) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf("%q", data)
	case reflect.Pointer:
		if value.IsNil() {
			return fmt.Sprintf("<%s nil>", p.typeToString(value.Type()))
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
		return fmt.Sprintf("<%s len=%d cap=%d ptr=%#x>", p.typeToString(value.Type()), value.Len(), value.Cap(), value.Pointer())
	case reflect.Struct:
		return p.structToString(value)
	case reflect.Func:
		return p.funcToString(value)
	default:
		return fmt.Sprintf("%#v", data)
	}
}

func (p *Processor) typeToString(type_ reflect.Type) string {
	switch type_.Kind() {
	case reflect.Pointer:
		return "*" + p.typeToString(type_.Elem())
	case reflect.Array, reflect.Slice:
		return fmt.Sprintf("[]%s", p.typeToString(type_.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s, %s]", p.typeToString(type_.Key()), p.typeToString(type_.Elem()))
	case reflect.Chan:
		return fmt.Sprintf("chan[%s]", p.typeToString(type_.Elem()))
	case reflect.Func:
		return fmt.Sprintf("func[%s -> %s]", p.funcInTypeString(type_), p.funcOutTypeString(type_))
	default:
		if type_.Name() == "" {
			return "?"
		} else {
			return type_.String()
		}
	}
}

func (p *Processor) funcInTypeString(type_ reflect.Type) string {
	numIn := type_.NumIn()
	if numIn == 0 {
		return "()"
	}
	if numIn == 1 {
		return p.typeToString(type_.In(0))
	}
	var sb strings.Builder
	sb.WriteString("(")
	for i := 0; i < numIn; i++ {
		sb.WriteString(p.typeToString(type_.In(i)))
		if i < numIn-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (p *Processor) funcOutTypeString(type_ reflect.Type) string {
	numOut := type_.NumOut()
	if numOut == 0 {
		return "()"
	}
	if numOut == 1 {
		return p.typeToString(type_.Out(0))
	}
	var sb strings.Builder
	sb.WriteString("(")
	for i := 0; i < numOut; i++ {
		sb.WriteString(p.typeToString(type_.Out(i)))
		if i < numOut-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (p *Processor) structToString(value reflect.Value) string {
	type_ := value.Type()

	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(p.typeToString(type_))
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

	return sb.String()
}

func (p *Processor) funcToString(value reflect.Value) string {
	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(p.typeToString(value.Type()))
	if value.IsNil() {
		sb.WriteString(" nil")
	} else {
		sb.WriteString(fmt.Sprintf(" ptr=%#x", value.Pointer()))
	}
	sb.WriteString(">")
	return sb.String()
}
