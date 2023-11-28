package main

import (
	"fmt"

	"github.com/dan-kuroto/gin-stronger/gs"
)

type Configuration struct {
	gs.Configuration
	Hello bool `yaml:"hello"`
}

var Config Configuration

// test
func main() {
	gs.Init(&gs.Config)
	// gs.InitDefault()
	// TODO: 如何让一个函数能接受一个struct和继承struct的struct作为参数？
	// gs.Init(&Config)
	fmt.Printf("GetConfig(): %v\n", gs.Config)
}
