package main

import (
	"fmt"

	gp "github.com/dan-kuroto/gin-stronger/go-print"
)

type A struct {
	Name string
	Data struct {
		Age     int
		Items   []string
		private int
	}
}

var f = gp.Formatter{
	ListShowAsTag:  true,
	MapShowAsTag:   true,
	ListIndent:     0,
	MapIndent:      3,
	StructIndent:   4,
	ListDisplayNum: 10,
}

func main() {
	a := make([]int, 0, 4)
	for i := 0; i < 10; i++ {
		a = append(a, i)
		fmt.Println(f.ToString(a))
	}
}
