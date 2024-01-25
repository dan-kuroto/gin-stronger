package check

type Checker struct {
	SolveError func(errMsg string)
}

func CheckParam[T any](checker *Checker, name string, value T, checkFuncs ...CheckFunc[T]) {
	for _, checkFunc := range checkFuncs {
		if errTpl := checkFunc(value); errTpl != "" {
			checker.SolveError(execErrorTemplate(errTpl, name, value))
		}
	}
}

func CheckParamCustom(checker *Checker, condition bool, errMsg string) {
	if !condition {
		checker.SolveError(errMsg)
	}
}

// Return error message template if value is invalid, otherwise return empty string.
//
// This template uses the following placeholders:
// - {{.name}} means the parameter name.
// - {{.value}} means the parameter value.
type CheckFunc[T any] func(value T) (errTpl string)

func NotEmptyStr(value string) (errTpl string) {
	if len(value) == 0 {
		return "{{.name}} must not be empty!"
	} else {
		return ""
	}
}

func NotEmptySlice[T any](value []T) (errTpl string) {
	if len(value) == 0 {
		return "{{.name}} must not be empty!"
	} else {
		return ""
	}
}

// TODO: NotBlank, InRange, Eq, Neq, Gt, Ge, Lt, Le, In, NotIn
//       (InRange是指范围区间，In/NotIn是枚举值)
