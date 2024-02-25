package check

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
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
	fmt.Printf("[WARNING] %s is invalid for the type of %s!\n", checkMethod, ctx.name)
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
// Only valid for bool/int/int8/.../uint/uint8/.../float32/float64/string and their pointer.
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

/*
// generate a CheckFunc to check whether value is greater than expect
//
// (value > expect)
func Gt[T orderable](expect T) CheckFunc[T] {
	return func(data Context[T]) error {
		if data.Value <= expect {
			return fmt.Errorf("%s must be greater than %s!", data.Name, formatter.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is greater than or equal to expect
//
// (value >= expect)
func Ge[T orderable](expect T) CheckFunc[T] {
	return func(data Context[T]) error {
		if data.Value < expect {
			return fmt.Errorf("%s must be greater than or equal to %s!", data.Name, formatter.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is less than expect
//
// (value < expect)
func Lt[T orderable](expect T) CheckFunc[T] {
	return func(data Context[T]) error {
		if data.Value >= expect {
			return fmt.Errorf("%s must be less than %s!", data.Name, formatter.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is less than or equal to expect
//
// (value <= expect)
func Le[T orderable](expect T) CheckFunc[T] {
	return func(data Context[T]) error {
		if data.Value > expect {
			return fmt.Errorf("%s must be less than or equal to %s!", data.Name, formatter.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is in expect
func In[T comparable](expect ...T) CheckFunc[T] {
	return func(data Context[T]) error {
		for _, v := range expect {
			if data.Value == v {
				return nil
			}
		}
		return fmt.Errorf("%s must be in %s!", data.Name, formatter.ToString(expect))
	}
}

// generate a CheckFunc to check whether value is not in expect
func NotIn[T comparable](expect ...T) CheckFunc[T] {
	return func(data Context[T]) error {
		for _, v := range expect {
			if data.Value == v {
				return fmt.Errorf("%s must not be in %s!", data.Name, formatter.ToString(expect))
			}
		}
		return nil
	}
}

// generate a CheckFunc to check whether value matches expect (expect is a regexp)
func Match(expect string) CheckFunc[string] {
	return func(data Context[string]) error {
		matched, err := regexp.MatchString(expect, data.Value)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("%s must match /%s/!", data.Name, expect)
		}
		return nil
	}
}
*/
