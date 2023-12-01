// test file
package main

import (
	"github.com/dan-kuroto/gin-stronger/gs"
	"github.com/gin-gonic/gin"
)

func init() {
	gs.InitConfigDefault()
}

func main() {
	engine := gin.Default()

	gs.UseRouters(engine, GetRouters())

	engine.Run(gs.Config.GetGinAddr())
}
