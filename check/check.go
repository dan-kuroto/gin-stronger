package check

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/dan-kuroto/gin-stronger/utils"
)

type Checker struct {
	SolveError func(err error)
}

type Data[T any] struct {
	Name  string
	Value T
}

type CheckFunc[T any] func(data Data[T]) error

type number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

type orderable interface {
	number | rune | string
}

func CheckParam[T any](checker *Checker, name string, value T, checkFuncs ...CheckFunc[T]) {
	data := Data[T]{Name: name, Value: value}
	for _, checkFunc := range checkFuncs {
		if err := checkFunc(data); err != nil {
			checker.SolveError(err)
		}
	}
}

func SimpleCheck(checker *Checker, condition bool, err error) {
	if !condition {
		checker.SolveError(err)
	}
}

func NotNil[T any](data Data[*T]) error {
	if data.Value == nil {
		return fmt.Errorf("`%s` is required!", data.Name)
	} else {
		return nil
	}
}

func NotEmptyStr(data Data[string]) error {
	if len(data.Value) == 0 {
		return fmt.Errorf("`%s` must not be empty!", data.Name)
	} else {
		return nil
	}
}

func NotEmptyList[T any](data Data[[]T]) error {
	if len(data.Value) == 0 {
		return fmt.Errorf("`%s` must not be empty!", data.Name)
	} else {
		return nil
	}
}

func NotEmptyMap[K comparable, V any](data Data[map[K]V]) error {
	if len(data.Value) == 0 {
		return fmt.Errorf("`%s` must not be empty!", data.Name)
	} else {
		return nil
	}
}

func NotBlank(data Data[string]) error {
	if strings.TrimSpace(data.Value) == "" {
		return fmt.Errorf("`%s` must not be blank!", data.Name)
	} else {
		return nil
	}
}

// check whether value consists of 0-9
func IsNumeric(data Data[string]) error {
	for _, ch := range data.Value {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("`%s` must be numeric!", data.Name)
		}
	}
	return nil
}

// check whether value is a valid email
func IsEmail(data Data[string]) error {
	if _, err := mail.ParseAddress(data.Value); err != nil {
		return fmt.Errorf("`%s` is not a valid email!", data.Name)
	}
	return nil
}

// check whether value is a valid URL
func IsURL(data Data[string]) error {
	if _, err := url.ParseRequestURI(data.Value); err != nil {
		return fmt.Errorf("`%s` is not a valid URL!", data.Name)
	}
	return nil
}

// generate a CheckFunc to check whether value in range of [min, max]
//
// (min <= value <= max)
func Range[T number](min, max T) CheckFunc[T] {
	return func(data Data[T]) error {
		if data.Value < min || data.Value > max {
			return fmt.Errorf("`%s` must be in range of [%v, %v]!", data.Name, min, max)
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether the size of value(slice, array, or map)
// in range of [min, max]
//
// (min <= len(value) <= max)
func Size[T any](min, max int) CheckFunc[[]T] {
	return func(data Data[[]T]) error {
		if len(data.Value) < min || len(data.Value) > max {
			return fmt.Errorf("`%s` must be in size of [%v, %v]!", data.Name, min, max)
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether the length of value(string) in range of [min, max]
//
// (min <= utf8.RuneCountInString((value) <= max)
func Length(min, max int) CheckFunc[string] {
	return func(data Data[string]) error {
		if utf8.RuneCountInString(data.Value) < min || utf8.RuneCountInString(data.Value) > max {
			return fmt.Errorf("`%s` must be in length of [%v, %v]!", data.Name, min, max)
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is equal to expect
//
// (value == expect)
func Eq[T comparable](expect T) CheckFunc[T] {
	return func(data Data[T]) error {
		if data.Value != expect {
			return fmt.Errorf("`%s` must be equal to %s!", data.Name, utils.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is not equal to expect
//
// (value != expect)
func Neq[T comparable](expect T) CheckFunc[T] {
	return func(data Data[T]) error {
		if data.Value == expect {
			return fmt.Errorf("`%s` must not be equal to %s!", data.Name, utils.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is greater than expect
//
// (value > expect)
func Gt[T orderable](expect T) CheckFunc[T] {
	return func(data Data[T]) error {
		if data.Value <= expect {
			return fmt.Errorf("`%s` must be greater than %s!", data.Name, utils.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is greater than or equal to expect
//
// (value >= expect)
func Ge[T orderable](expect T) CheckFunc[T] {
	return func(data Data[T]) error {
		if data.Value < expect {
			return fmt.Errorf("`%s` must be greater than or equal to %s!", data.Name, utils.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is less than expect
//
// (value < expect)
func Lt[T orderable](expect T) CheckFunc[T] {
	return func(data Data[T]) error {
		if data.Value >= expect {
			return fmt.Errorf("`%s` must be less than %s!", data.Name, utils.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is less than or equal to expect
//
// (value <= expect)
func Le[T orderable](expect T) CheckFunc[T] {
	return func(data Data[T]) error {
		if data.Value > expect {
			return fmt.Errorf("`%s` must be less than or equal to %s!", data.Name, utils.ToString(expect))
		} else {
			return nil
		}
	}
}

// generate a CheckFunc to check whether value is in expect
func In[T comparable](expect ...T) CheckFunc[T] {
	return func(data Data[T]) error {
		for _, v := range expect {
			if data.Value == v {
				return nil
			}
		}
		return fmt.Errorf("`%s` must be in %s!", data.Name, utils.ToString(expect))
	}
}

// generate a CheckFunc to check whether value is not in expect
func NotIn[T comparable](expect ...T) CheckFunc[T] {
	return func(data Data[T]) error {
		for _, v := range expect {
			if data.Value == v {
				return fmt.Errorf("`%s` must not be in %s!", data.Name, utils.ToString(expect))
			}
		}
		return nil
	}
}

// generate a CheckFunc to check whether value matches expect (expect is a regexp)
func Match(expect string) CheckFunc[string] {
	return func(data Data[string]) error {
		matched, err := regexp.MatchString(expect, data.Value)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("`%s` must match /%s/!", data.Name, expect)
		}
		return nil
	}
}
