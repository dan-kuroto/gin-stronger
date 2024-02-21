package check

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
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

	if length, ok := getLength(ctx.value); ok && length == 0 {
		ctx.err = fmt.Errorf("%s must not be empty!", ctx.name)
		ctx.solveError(ctx.err)
	}

	return ctx
}

func (ctx *Context) NotBlank() *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := ctx.value.(string); ok && strings.TrimSpace(value) == "" {
		ctx.err = fmt.Errorf("%s must not be blank!", ctx.name)
		ctx.solveError(ctx.err)
	}

	return ctx
}

// check whether value consists of 0-9
func (ctx *Context) IsNumeric() *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := ctx.value.(string); ok {
		for _, ch := range value {
			if ch < '0' || ch > '9' {
				ctx.err = fmt.Errorf("%s must be numeric!", ctx.name)
				ctx.solveError(ctx.err)
				break
			}
		}
	}

	return ctx
}

// check whether value is a valid email
func (ctx *Context) IsEmail() *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := ctx.value.(string); ok {
		if _, err := mail.ParseAddress(value); err != nil {
			ctx.err = fmt.Errorf("%s is not a valid email!", ctx.name)
			ctx.solveError(ctx.err)
		}
	}

	return ctx
}

// check whether value is a valid URL
func (ctx *Context) IsURL() *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := ctx.value.(string); ok {
		if _, err := url.ParseRequestURI(value); err != nil {
			ctx.err = fmt.Errorf("%s is not a valid URL!", ctx.name)
			ctx.solveError(ctx.err)
		}
	}

	return ctx
}

// check whether min <= value <= max
func (ctx *Context) Range(min, max float64) *Context {
	if ctx.err != nil {
		return ctx
	}

	if value, ok := toFloat64(ctx.value); ok {
		if value < min || value > max {
			ctx.err = fmt.Errorf("%s must be in range of [%v, %v]!", ctx.name, min, max)
			ctx.solveError(ctx.err)
		}
	}

	return ctx
}

/*
// generate a CheckFunc to check whether the size of value(slice, array, or map)
// in range of [min, max]
//
// (min <= len(value) <= max)
func Size[T any](min, max int) CheckFunc[[]T] {
	return func(data Context[[]T]) error {
		if len(data.Value) < min || len(data.Value) > max {
			return fmt.Errorf("%s must be in size of [%v, %v]!", data.Name, min, max)
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether the length of value(string) in range of [min, max]
//
// (min <= utf8.RuneCountInString((value) <= max)
func Length(min, max int) CheckFunc[string] {
	return func(data Context[string]) error {
		if utf8.RuneCountInString(data.Value) < min || utf8.RuneCountInString(data.Value) > max {
			return fmt.Errorf("%s must be in length of [%v, %v]!", data.Name, min, max)
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is equal to expect
//
// (value == expect)
func Eq[T comparable](expect T) CheckFunc[T] {
	return func(data Context[T]) error {
		if data.Value != expect {
			return fmt.Errorf("%s must be equal to %s!", data.Name, formatter.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is not equal to expect
//
// (value != expect)
func Neq[T comparable](expect T) CheckFunc[T] {
	return func(data Context[T]) error {
		if data.Value == expect {
			return fmt.Errorf("%s must not be equal to %s!", data.Name, formatter.ToString(expect))
		} else {
			return nil
		}
	}
}

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
