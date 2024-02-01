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
	ListIndent:   2,
	MapIndent:    2,
	StructIndent: 2,
}

func main() {
	a := A{}
	fmt.Println(f.ToString(a))
}
