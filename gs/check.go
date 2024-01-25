package gs

import "github.com/dan-kuroto/gin-stronger/check"

var defaultChecker = check.Checker{
	SolveError: func(errMsg string) {
		panic(errMsg)
	},
}

func SetDefaultChecker(checker *check.Checker) {
	defaultChecker = *checker
}

func CheckParam[T any](name string, value T, checkFuncs ...check.CheckFunc[T]) {
	check.CheckParam(&defaultChecker, name, value, checkFuncs...)
}

func CheckParamCustom(condition bool, errMsg string) {
	check.CheckParamCustom(&defaultChecker, condition, errMsg)
}
