package gs

import (
	"errors"

	"github.com/dan-kuroto/gin-stronger/check"
)

var defaultChecker = check.Checker{
	SolveError: func(err error) {
		panic(err.Error())
	},
}

func SetDefaultChecker(checker *check.Checker) {
	defaultChecker = *checker
}

func CheckParam[T any](name string, value T, checkFuncs ...check.CheckFunc[T]) {
	check.CheckParam(&defaultChecker, name, value, checkFuncs...)
}

func SimpleCheck(condition bool, errMsg string) {
	check.SimpleCheck(&defaultChecker, condition, errors.New(errMsg))
}
