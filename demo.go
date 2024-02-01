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
	ListIndent:     0,
	MapIndent:      2,
	StructIndent:   2,
	ListDisplayNum: 0,
	MapDisplayNum:  0,
}

func main() {
	fmt.Println(f.ToString(map[string]any{"1": 1, "2": map[string]any{"3": "#"}}))
	a := &A{Data: struct {
		Age     int
		Items   []string
		private int
	}{Items: []string{"a", "b", "c"}}}
	fmt.Println(f.ToString(a))
}
