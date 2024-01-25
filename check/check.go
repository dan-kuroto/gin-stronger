package check

type Checker struct {
	SolveError func(errMsg string)
}

func CheckParam[T any](checker *Checker, name string, value T, checkFuncs ...CheckFunc[T]) {
	for _, checkFunc := range checkFuncs {
		if errTpl := checkFunc(value); errTpl != "" {
			checker.SolveError(execErrTpl(errTpl, name, value))
		}
	}
}

// TODO: 每个函数加注释
func CheckParamCustom(checker *Checker, condition bool, errMsg string) {
	if !condition {
		checker.SolveError(errMsg)
	}
}

// Return error message template if value is invalid, otherwise return empty string.
type CheckFunc[T any] func(value T) (errTpl string)
