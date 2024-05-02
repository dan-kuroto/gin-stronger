package gp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

type Formatter struct {
	// MaxIndentLayer specifies the maximum indentation level for formatting.
	// A value <= 0 allows unlimited indentation, while non-zero values
	// will hide additional indentation beyond the specified limit.
	MaxIndentLayer int
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
	// Max number for display items in struct.
	// If `StructDisplayNum <= 0`, means infinity.
	// If `len(data) > StructDisplayNum`, extra parts are shown as ellipsis.
	StructDisplayNum int
	// Max number for display items in array&list.
	// If `ListDisplayNum <= 0`, means infinity.
	// If `len(data) > ListDisplayNum`, extra parts are shown as ellipsis.
	ListDisplayNum int
	// Max number for display items in map.
	// If `MapDisplayNum <= 0`, means infinity.
	// If `len(data) > MapDisplayNum`, extra parts are shown as ellipsis.
	MapDisplayNum int
	BracketColor  bool
	// Determines whether to display strings with quotes.
	//
	// When the string is nested (e.g., within a slice, struct, or similar composite),
	// it will always appear enclosed in quotes, regardless of the value of StrQuote.
	StrQuote bool
}

// context for go-print Formatter
type gpfContext struct {
	Indents []int
	IsArray bool
	// Primarily to resolve the problem to fetch the address of struct pointer.
	CurrPtr uintptr
	// A map of visited pointers to prevent infinite recursion.
	// Therefore, actually only pointer of list/map/struct are stored.
	VistedPointers map[uintptr]bool
}

var DefaultFormatter = Formatter{
	MaxIndentLayer:   10,
	ListDisplayNum:   100,
	MapDisplayNum:    100,
	StructDisplayNum: 100,
	BracketColor:     true,
	StrQuote:         true,
}

var (
	normalColors = [8]*color.Color{
		color.New(color.FgBlack),
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
		color.New(color.FgWhite),
	}
	boldColors = [8]*color.Color{
		color.New(color.FgBlack, color.Bold),
		color.New(color.FgRed, color.Bold),
		color.New(color.FgGreen, color.Bold),
		color.New(color.FgYellow, color.Bold),
		color.New(color.FgBlue, color.Bold),
		color.New(color.FgMagenta, color.Bold),
		color.New(color.FgCyan, color.Bold),
		color.New(color.FgWhite, color.Bold),
	}
)

func ToString(data any) string {
	return DefaultFormatter.ToString(data)
}

func (f *Formatter) ToString(data any) string {
	return f.toString(data, gpfContext{VistedPointers: map[uintptr]bool{}})
}

func (f *Formatter) toString(data any, ctx gpfContext) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.String:
		if f.StrQuote && len(ctx.Indents) != 0 {
			return fmt.Sprintf("%q", data)
		} else {
			return fmt.Sprint(data)
		}
	case reflect.Pointer:
		if value.IsNil() {
			return fmt.Sprintf("<%s nil>", f.typeToString(value.Type()))
		}
		if !value.CanInterface() {
			return fmt.Sprintf("%#v", data)
		}
		ctx.CurrPtr = value.Pointer()
		return "&" + f.toString(value.Elem().Interface(), ctx)
	case reflect.Array:
		ctx.IsArray = true
		return f.listToString(value, ctx)
	case reflect.Slice:
		ctx.IsArray = false
		ctx.CurrPtr = value.Pointer()
		return f.listToString(value, ctx)
	case reflect.Map:
		ctx.CurrPtr = value.Pointer()
		return f.mapToString(value, ctx)
	case reflect.Chan:
		return fmt.Sprintf("<%s :len=%d :cap=%d :ptr=%#x>", f.typeToString(value.Type()), value.Len(), value.Cap(), value.Pointer())
	case reflect.Struct:
		return f.structToString(value, ctx)
	case reflect.Func:
		return f.funcToString(value)
	default:
		return fmt.Sprintf("%#v", data)
	}
}

func (f *Formatter) listToString(value reflect.Value, ctx gpfContext) string {
	var sb strings.Builder
	length := value.Len()
	displayLength := length
	if f.ListDisplayNum > 0 && length > f.ListDisplayNum {
		displayLength = f.ListDisplayNum
	}
	if f.MaxIndentLayer > 0 && len(ctx.Indents) >= f.MaxIndentLayer {
		displayLength = 0
	}

	if f.ListShowAsTag {
		appendColoredString(&sb, fmt.Sprint("<", f.typeToString(value.Type())), len(ctx.Indents), f.BracketColor, true)
		if !ctx.IsArray {
			sb.WriteString(fmt.Sprintf(" :len=%d :cap=%d", value.Len(), value.Cap()))
		}
		if ctx.CurrPtr != 0 { // pointer
			sb.WriteString(fmt.Sprintf(" :ptr=%#x", ctx.CurrPtr))
		}
		sb.WriteString(" :items=")
	}
	if ctx.CurrPtr != 0 {
		if ctx.VistedPointers[ctx.CurrPtr] { // skip visited pointer
			displayLength = 0
		} else { // store unvisited pointer
			ctx.VistedPointers[ctx.CurrPtr] = true
		}
		// reset after use
		ctx.CurrPtr = 0
	}
	appendColoredString(&sb, "[", len(ctx.Indents), f.BracketColor, true)
	if displayLength > 0 { // can show items
		ctx.Indents = append(ctx.Indents, f.ListIndent)
		appendIndent(&sb, f.ListIndent, ctx.Indents, false, f.BracketColor)
		for i := 0; i < displayLength; i++ {
			sb.WriteString(f.toString(value.Index(i).Interface(), ctx))
			if i < displayLength-1 { // before the last one
				sb.WriteString(",")
				appendIndent(&sb, f.ListIndent, ctx.Indents, true, f.BracketColor)
			} else if displayLength < length { // last one, but need ellipsis
				sb.WriteString(",")
				appendIndent(&sb, f.ListIndent, ctx.Indents, true, f.BracketColor)
				sb.WriteString("...")
			}
		}
		ctx.Indents = ctx.Indents[:len(ctx.Indents)-1]
		appendIndent(&sb, f.ListIndent, ctx.Indents, false, f.BracketColor)
	} else if length > 0 { // cannot show items, but actually has items
		sb.WriteString("...")
	}
	appendColoredString(&sb, "]", len(ctx.Indents), f.BracketColor, true)
	if f.ListShowAsTag {
		appendColoredString(&sb, ">", len(ctx.Indents), f.BracketColor, true)
	}

	return sb.String()
}

func (f *Formatter) mapToString(value reflect.Value, ctx gpfContext) string {
	var sb strings.Builder
	length := value.Len()
	displayLength := length
	if f.MapDisplayNum > 0 && length > f.MapDisplayNum {
		displayLength = f.MapDisplayNum
	}
	if f.MaxIndentLayer > 0 && len(ctx.Indents) >= f.MaxIndentLayer {
		displayLength = 0
	}

	if f.MapShowAsTag {
		appendColoredString(&sb, fmt.Sprint("<", f.typeToString(value.Type())), len(ctx.Indents), f.BracketColor, true)
		sb.WriteString(fmt.Sprintf(" :len=%d", value.Len()))
		if ctx.CurrPtr != 0 { // pointer
			sb.WriteString(fmt.Sprintf(" :ptr=%#x", ctx.CurrPtr))
		}
	} else {
		appendColoredString(&sb, "{", len(ctx.Indents), f.BracketColor, true)
	}
	if ctx.CurrPtr != 0 {
		if ctx.VistedPointers[ctx.CurrPtr] { // skip visited pointer
			displayLength = 0
		} else { // store unvisited pointer
			ctx.VistedPointers[ctx.CurrPtr] = true
		}
		// reset after use
		ctx.CurrPtr = 0
	}
	if displayLength > 0 { // can show items
		ctx.Indents = append(ctx.Indents, f.MapIndent)
		appendIndent(&sb, f.MapIndent, ctx.Indents, f.MapShowAsTag, f.BracketColor)
		for i, key := range value.MapKeys() {
			sb.WriteString(f.toString(key.Interface(), ctx))
			if f.MapShowAsTag {
				sb.WriteString("=")
			} else {
				sb.WriteString(": ")
			}
			sb.WriteString(f.toString(value.MapIndex(key).Interface(), ctx))
			if i < displayLength-1 { // before the last one
				if !f.MapShowAsTag {
					sb.WriteString(",")
				}
				appendIndent(&sb, f.MapIndent, ctx.Indents, true, f.BracketColor)
			} else { // last one
				if displayLength < length { // need ellipsis
					if !f.MapShowAsTag {
						sb.WriteString(",")
					}
					appendIndent(&sb, f.MapIndent, ctx.Indents, true, f.BracketColor)
					sb.WriteString("...")
				}
				break
			}
		}
		ctx.Indents = ctx.Indents[:len(ctx.Indents)-1]
		appendIndent(&sb, f.MapIndent, ctx.Indents, false, f.BracketColor)
	} else if length > 0 { // cannot show items, but actually has items
		if f.MapShowAsTag {
			sb.WriteRune(' ')
		}
		sb.WriteString("...")
	}
	if f.MapShowAsTag {
		appendColoredString(&sb, ">", len(ctx.Indents), f.BracketColor, true)
	} else {
		appendColoredString(&sb, "}", len(ctx.Indents), f.BracketColor, true)
	}

	return sb.String()
}

func (f *Formatter) typeToString(type_ reflect.Type) string {
	switch type_.Kind() {
	case reflect.Pointer:
		return "*" + f.typeToString(type_.Elem())
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", type_.Len(), f.typeToString(type_.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", f.typeToString(type_.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s, %s]", f.typeToString(type_.Key()), f.typeToString(type_.Elem()))
	case reflect.Chan:
		return fmt.Sprintf("chan[%s]", f.typeToString(type_.Elem()))
	case reflect.Func:
		return fmt.Sprintf("func[%s -> %s]", f.funcParamsTypeToString(type_), f.funcReturnsTypeString(type_))
	default:
		if type_.Name() == "" {
			return "?"
		} else {
			return type_.String()
		}
	}
}

func (f *Formatter) funcParamsTypeToString(type_ reflect.Type) string {
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

func (f *Formatter) funcReturnsTypeString(type_ reflect.Type) string {
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

func (f *Formatter) structToString(value reflect.Value, ctx gpfContext) string {
	type_ := value.Type()
	var fields = make(map[int]reflect.StructField)
	for i := 0; i < value.NumField(); i++ {
		field := type_.Field(i)
		if field.IsExported() {
			fields[i] = field
		}
	}
	length := len(fields)
	displayLength := length
	if f.StructDisplayNum > 0 && length > f.StructDisplayNum {
		displayLength = f.StructDisplayNum
	}
	if f.MaxIndentLayer > 0 && len(ctx.Indents) >= f.MaxIndentLayer {
		displayLength = 0
	}

	var sb strings.Builder

	appendColoredString(&sb, fmt.Sprint("<", f.typeToString(type_)), len(ctx.Indents), f.BracketColor, true)
	if ctx.CurrPtr != 0 { // pointer
		sb.WriteString(fmt.Sprintf(" :ptr=%#x", ctx.CurrPtr))

		if ctx.VistedPointers[ctx.CurrPtr] { // skip visited pointer
			displayLength = 0
		} else { // store unvisited pointer
			ctx.VistedPointers[ctx.CurrPtr] = true
		}
		// reset after use
		ctx.CurrPtr = 0
	}
	if displayLength > 0 { // can show items
		ctx.Indents = append(ctx.Indents, f.StructIndent)
		appendIndent(&sb, f.StructIndent, ctx.Indents, true, f.BracketColor)
		cnt := 0
		for i, field := range fields {
			sb.WriteString(field.Name)
			sb.WriteString("=")
			sb.WriteString(f.toString(value.Field(i).Interface(), ctx))
			if cnt < displayLength-1 { // before the last one
				appendIndent(&sb, f.StructIndent, ctx.Indents, true, f.BracketColor)
			} else { // last one
				if displayLength < length { // need ellipsis
					appendIndent(&sb, f.StructIndent, ctx.Indents, true, f.BracketColor)
					sb.WriteString("...")
				}
				break
			}
			cnt++
		}
		ctx.Indents = ctx.Indents[:len(ctx.Indents)-1]
		appendIndent(&sb, f.StructIndent, ctx.Indents, false, f.BracketColor)
	} else if length > 0 { // cannot show items, but actually has items
		sb.WriteString(" ...")
	}
	appendColoredString(&sb, ">", len(ctx.Indents), f.BracketColor, true)

	return sb.String()
}

func (f *Formatter) funcToString(value reflect.Value) string {
	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(f.typeToString(value.Type()))
	if value.IsNil() {
		sb.WriteString(" nil")
	} else {
		sb.WriteString(fmt.Sprintf(" :ptr=%#x", value.Pointer()))
	}
	sb.WriteString(">")
	return sb.String()
}
