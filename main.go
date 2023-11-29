package main

import (
	"fmt"

	"github.com/dan-kuroto/gin-stronger/gs"
)

type Configuration2 struct {
	gs.Configuration `yaml:",inline"`
	Hello            string `yaml:"hello"`
}

var Config Configuration2

// test
func init() {
	// gs.InitConfigDefault()
	// gs.InitConfig(&gs.Config)
	gs.InitConfig(&Config)
}

func main() {
	fmt.Printf("GetConfig(): %v\n", gs.Config)
	fmt.Printf("GetConfig(): %v\n", Config)
}
