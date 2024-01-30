package main

import (
	"fmt"

	"github.com/dan-kuroto/gin-stronger/utils"
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
	fmt.Println(utils.ToString(struct{ A string }{}))
	fmt.Println(utils.ToString(A{}))
	fmt.Println(utils.ToString(func(s string) A {
		return A{}
	}))
}
