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
	gs.Check("a", a).Eq(1).Eq(1.0).Eq(int64(1)).Eq("1").Eq(2.1)
	gs.Check("&a", &a).Eq(1).Eq(1.0).Eq(int64(1)).Eq("1").Eq(2)
	b := "1"
	gs.Check("b", b).Eq("1").Eq("2")
	gs.Check("&b", &b).Eq("1").Eq("2")
	c := true
	gs.Check("c", c).Eq(true).Eq(false)
	gs.Check("c", &c).Eq(true).Eq(false)
	d := check.Checker{}
	gs.Check("d", d).Eq(nil)
	gs.Check("&d", &d).Eq(nil)
}
