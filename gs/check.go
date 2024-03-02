package gs

import (
	"github.com/dan-kuroto/gin-stronger/check"
)

var defaultChecker = &check.Checker{
	SolveError: func(err error) {
		panic(err.Error())
	},
}

func SetDefaultChecker(checker *check.Checker) {
	defaultChecker = checker
}

func Check(name string, data any) *check.Context {
	return defaultChecker.Check(name, data)
}

func Assert(condition bool, errMsg string) {
	defaultChecker.Assert(condition, errMsg)
}

func AssertError(err error) {
	defaultChecker.AssertError(err)
}
