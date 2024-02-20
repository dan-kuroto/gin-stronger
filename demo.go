package main

import (
	"log"

	"github.com/dan-kuroto/gin-stronger/check"
	"github.com/dan-kuroto/gin-stronger/gs"
)

func main() {
	gs.SetDefaultChecker(&check.Checker{
		SolveError: func(err error) {
			log.Println(err)
		},
	})
	// gs.Check("a", "123").NotEmpty().IsNumeric().NotBlank().Range(1, 2)
	gs.Check("b", 1).Range(0.1, 1.1).Range(1.0, 2).Range(1.0, 1).Range(2, 2)
}
