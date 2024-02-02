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
	ListShowAsTag: true,
}

func main() {
	fmt.Println(f.ToString(map[string]any{
		"1": 1,
		"2": map[string]any{
			"3": "4",
			"5": []any{
				6,
				7,
				[3]string{"8", "9"},
				[]int{},
			},
		},
	}))
}
