package check

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Checker struct {
	SolveError func(err error)
}

func (checker *Checker) Check(name string, data any) *Context {
	return &Context{
		name:       name,
		value:      data,
		solveError: checker.SolveError,
	}
}

func (checker *Checker) Assert(condition bool, errMsg string) {
	if !condition {
		checker.SolveError(errors.New(errMsg))
	}
}

type Context struct {
	name       string
	value      any
	err        error
	solveError func(err error)
}

func (ctx *Context) Error() error {
	return ctx.err
}

func (ctx *Context) printTypeWarning(checkMethod string) {
	printWarning("%s is invalid for the type of %s!", checkMethod, ctx.name)
}

func (ctx *Context) Assert(condition bool, errMsg string) *Context {
	if ctx.err != nil {
		return ctx
	}

	if !condition {
		ctx.err = errors.New(errMsg)
		ctx.solveError(ctx.err)
	}

	return ctx
}

func (ctx *Context) NotNil() *Context {
	if ctx.err != nil {
		return ctx
	}

	if ctx.value == nil {
		ctx.err = fmt.Errorf("%s is required!", ctx.name)
		ctx.solveError(ctx.err)
	}

	return ctx
}

func (ctx *Context) NotEmpty() *Context {
	if ctx.err != nil {
		return ctx
	}

	if length, ok := getLength(ctx.value); ok {
		if length == 0 {
			ctx.err = fmt.Errorf("%s must not be empty!", ctx.name)
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".NotEmpty()")
	}

	return ctx
}

func (ctx *Context) NotBlank() *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := toString(ctx.value); ok {
		if strings.TrimSpace(value) == "" {
			ctx.err = fmt.Errorf("%s must not be blank!", ctx.name)
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".NotBlank()")
	}

	return ctx
}

// check whether value consists of 0-9
func (ctx *Context) IsNumeric() *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := toString(ctx.value); ok {
		for _, ch := range value {
			if ch < '0' || ch > '9' {
				ctx.err = fmt.Errorf("%s must be numeric!", ctx.name)
				ctx.solveError(ctx.err)
				break
			}
		}
	} else {
		ctx.printTypeWarning(".IsNumeric()")
	}

	return ctx
}

// check whether value is a valid email
func (ctx *Context) IsEmail() *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := toString(ctx.value); ok {
		if _, err := mail.ParseAddress(value); err != nil {
			ctx.err = fmt.Errorf("%s is not a valid email!", ctx.name)
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".IsEmail()")
	}

	return ctx
}

// check whether value is a valid URL
func (ctx *Context) IsURL() *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := toString(ctx.value); ok {
		if _, err := url.ParseRequestURI(value); err != nil {
			ctx.err = fmt.Errorf("%s is not a valid URL!", ctx.name)
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".IsURL()")
	}

	return ctx
}

// check whether min <= value <= max
//
// Only valid for int/int8/.../uint/uint8/.../float32/float64 or pointer to them.
func (ctx *Context) Range(min, max float64) *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := toFloat64(ctx.value); ok {
		if value < min || value > max {
			ctx.err = fmt.Errorf("%s must be in range of [%v, %v]!", ctx.name, min, max)
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Range(min, max)")
	}

	return ctx
}

// check whether min <= len(value) <= max (value can be string, slice, array, or map)
//
// For string, what is checked is the number of bytes.
// If you want to check the number of characters(rune), use `Length`.
func (ctx *Context) Size(min, max int) *Context {
	if ctx.err != nil {
		return ctx
	}

	if length, ok := getLength(ctx.value); ok {
		if length < min || length > max {
			ctx.err = fmt.Errorf("%s must be in size of [%v, %v]!", ctx.name, min, max)
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Size(min, max)")
	}

	return ctx
}

// check whether min <= utf8.RuneCountInString(value) <= max
func (ctx *Context) Length(min, max int) *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := toString(ctx.value); ok {
		length := utf8.RuneCountInString(value)
		if length < min || length > max {
			ctx.err = fmt.Errorf("%s must be in length of [%v, %v]!", ctx.name, min, max)
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Length(min, max)")
	}

	return ctx
}

// check whether value == expect
//
// Only valid for int/int8/.../uint/uint8/.../float32/float64/string and their pointer.
//
// It will also be invalid if the types do not match. For example:
//
//	If value is int and expect is string, it is invalid;
//	If value is int and expect is float32, it is valid.
func (ctx *Context) Eq(expect any) *Context {
	if ctx.err != nil {
		return ctx
	}

	if equal, ok := basicEqual(ctx.value, expect); ok {
		if !equal {
			ctx.err = fmt.Errorf("%s must be equal to %s!", ctx.name, formatter.ToString(expect))
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Eq(expect)")
	}

	return ctx
}

// Check whether value != expect. The type handling mechanism is the same as `Eq`.
func (ctx *Context) Neq(expect any) *Context {
	if ctx.err != nil {
		return ctx
	}

	if equal, ok := basicEqual(ctx.value, expect); ok {
		if equal {
			ctx.err = fmt.Errorf("%s must not be equal to %s!", ctx.name, formatter.ToString(expect))
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Neq(expect)")
	}

	return ctx
}

// Check whether value > expect. The type handling mechanism is the same as `Eq`.
func (ctx *Context) Gt(expect any) *Context {
	if ctx.err != nil {
		return ctx
	}

	if greater, ok := basicGreater(ctx.value, expect); ok {
		if !greater {
			ctx.err = fmt.Errorf("%s must be greater than %s!", ctx.name, formatter.ToString(expect))
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Gt(expect)")
	}

	return ctx
}

// Check whether value >= expect. The type handling mechanism is the same as `Eq`.
func (ctx *Context) Ge(expect any) *Context {
	if ctx.err != nil {
		return ctx
	}

	if less, ok := basicLess(ctx.value, expect); ok {
		if less {
			ctx.err = fmt.Errorf("%s must be greater than or equal to %s!", ctx.name, formatter.ToString(expect))
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Ge(expect)")
	}

	return ctx
}

// Check whether value < expect. The type handling mechanism is the same as `Eq`.
func (ctx *Context) Lt(expect any) *Context {
	if ctx.err != nil {
		return ctx
	}

	if less, ok := basicLess(ctx.value, expect); ok {
		if !less {
			ctx.err = fmt.Errorf("%s must be less than %s!", ctx.name, formatter.ToString(expect))
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Lt(expect)")
	}

	return ctx
}

// Check whether value <= expect. The type handling mechanism is the same as `Eq`.
func (ctx *Context) Le(expect any) *Context {
	if ctx.err != nil {
		return ctx
	}

	if greater, ok := basicGreater(ctx.value, expect); ok {
		if greater {
			ctx.err = fmt.Errorf("%s must be less than or equal to %s!", ctx.name, formatter.ToString(expect))
			ctx.solveError(ctx.err)
		}
	} else {
		ctx.printTypeWarning(".Le(expect)")
	}

	return ctx
}

// Check whether value is in expect.
//
// When compare value with items in expect, it is only valid for int/int8/.../uint/uint8/.../float32/float64/string and their pointer.
func (ctx *Context) In(expect ...any) *Context {
	if ctx.err != nil {
		return ctx
	}

	if in := basicIn(ctx.value, expect); !in {
		ctx.err = fmt.Errorf("%s must be in %s!", ctx.name, formatter.ToString(expect))
		ctx.solveError(ctx.err)
	}

	return ctx
}

// Check whether value is not in expect. The type handling mechanism is the same as `In`.
func (ctx *Context) NotIn(expect ...any) *Context {
	if ctx.err != nil {
		return ctx
	}

	if in := basicIn(ctx.value, expect); in {
		ctx.err = fmt.Errorf("%s must not be in %s!", ctx.name, formatter.ToString(expect))
		ctx.solveError(ctx.err)
	}

	return ctx
}

// Check whether value matches expect (expect is a regexp)
//
// Only valid when value is string/*string
func (ctx *Context) Match(expect string) *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := toString(ctx.value); ok {
		if matched, err := regexp.MatchString(expect, value); err != nil {
			printWarning(err.Error())
			ctx.err = err
			ctx.solveError(ctx.err)
		} else {
			if !matched {
				ctx.err = fmt.Errorf("%s must match /%s/!", ctx.name, expect)
				ctx.solveError(ctx.err)
			}
		}
	} else {
		ctx.printTypeWarning(".Match(expect)")
	}

	return ctx
}
