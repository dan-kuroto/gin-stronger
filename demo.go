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

func main() {
	var a func() A
	b := func(s string, b func(int) map[string]A) (A, error) {
		return A{}, nil
	}
	fmt.Println(gp.ToString(a))
	fmt.Println(gp.ToString(b))
}
