package main

import (
	"fmt"
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
	var a any = map[string]int{"a": 1}
	fmt.Println(gs.ToString(a))
	b, ok := a.(map[string]any)
	fmt.Println(ok, gs.ToString(b))
}
