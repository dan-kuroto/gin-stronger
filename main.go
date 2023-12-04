// test file
package main

import (
	"fmt"

	"github.com/dan-kuroto/gin-stronger/gs"
	"github.com/gin-gonic/gin"
)

func init() {
	gs.InitConfigDefault()
}

func main() {
	engine := gin.Default()

	// TODO: 中间件要在Routers之前注册，试试能不能整个活把逻辑交给我来解决
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		fmt.Printf("err: %v\n", err)
		c.AbortWithStatusJSON(404, gin.H{"message": "Not found"})
	}))
	gs.UsePanicHandler("", func(c *gin.Context, err string) {
		c.AbortWithStatusJSON(404, gin.H{"message": err})
	})
	gs.UseRouters(engine, GetRouters())

	engine.Run(gs.Config.GetGinAddr())
}
