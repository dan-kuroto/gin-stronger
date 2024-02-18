package gs

import (
	"github.com/dan-kuroto/gin-stronger/check"
)

type Checker check.Checker

var defaultChecker = Checker{
	SolveError: func(err error) {
		panic(err.Error())
	},
}

func SetDefaultChecker(checker *Checker) {
	defaultChecker = *checker
}
