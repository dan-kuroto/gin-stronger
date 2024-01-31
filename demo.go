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
	fmt.Println(gp.ToString(struct{ A string }{}))
	fmt.Println(gp.ToString(A{}))
	fmt.Println(gp.ToString(func(s string) A {
		return A{}
	}))
}
