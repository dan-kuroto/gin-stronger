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
	// TODO: 最后的bug,struct变量为指针时,如果没有请求体,就会报错,但我认为应该给一个零值的!
	engine := gin.Default()

	gs.UseRouters(engine, GetRouters())

	engine.Run(gs.Config.GetGinAddr())
}
