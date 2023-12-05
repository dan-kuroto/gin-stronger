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

	// TODO: 中间件要在Routers之前注册，试试能不能整个活把逻辑交给我来解决
	gs.UsePanicHandler(engine, func(c *gin.Context, err string) {
		c.JSON(500, gin.H{"errMessage": err})
	})
	gs.UsePanicHandler(engine, func(c *gin.Context, err any) {
		c.JSON(500, gin.H{"errMessage": "unknown err"})
	})
	gs.UseRouters(engine, GetRouters())

	engine.Run(gs.Config.GetGinAddr())
}
