// test file
package main

import (
	"fmt"

	"github.com/dan-kuroto/gin-stronger/gs"
)

func init() {
	gs.InitConfigDefault()
}

func main() {
	fmt.Printf("config: %v\n", gs.Config)
}
