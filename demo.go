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
	// 测试要求下面每一个都报错,且都错在最后一个校验点上
	a := 1
	gs.Check("a", a).NotNil().Size(1, 2).Range(2, 3)
	gs.Check("a2", &a).NotNil().Size(1, 2).Range(2, 3)
	b := "1"
	gs.Check("b", b).Size(1, 2).Size(0, 0)
	gs.Check("b2", &b).Size(1, 2).Size(0, 0)
	c := []string{"1", "2"}
	gs.Check("c", c).Size(1, 2).Size(0, 0)
	gs.Check("c2", &c).Size(1, 2).Size(0, 0)
	var d []int
	gs.Check("d", d).Size(1, 2)
	gs.Check("d2", &d).Size(1, 2)
	e := [2]int{}
	gs.Check("e", e).Size(1, 2).Size(0, 0)
	gs.Check("e2", &e).Size(1, 2).Size(0, 0)
	f := map[string]int{"1": 1, "2": 2}
	gs.Check("f", f).Size(1, 2).Size(3, 4)
	gs.Check("f2", &f).Size(1, 2).Size(3, 4)
	var g map[string]int
	gs.Check("g", g).Size(1, 2)
	gs.Check("g2", g).Size(1, 2)
}
