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
	a := 1
	gs.Check("a", a).NotIn(2, 3).NotIn(1, 2)
	gs.Check("&a", &a).NotIn(2, 3).NotIn(1, 2)
	b := "1"
	gs.Check("b", b).NotIn(1, 2).NotIn("1", "2")
	gs.Check("&b", &b).NotIn(1, 2).NotIn("1", "2")
	c := 1.1
	gs.Check("c", c).NotIn(1.2, 1.3).NotIn(1.1, 1.2)
	gs.Check("&c", &c).NotIn(1.2, 1.3).NotIn(1.1, 1.2)
}
