package main

import (
	"github.com/dan-kuroto/gin-stronger/config"
	gp "github.com/dan-kuroto/gin-stronger/go-print"
	"github.com/dan-kuroto/gin-stronger/gs"
	"github.com/gin-gonic/gin"
)

func main() {
	gs.Config = &config.Configuration{}
	gs.InitIdGenerators()

	logger := gp.NewLogger("demo")
	logger.Print("hello world")
	ctx := &gin.Context{}
	// TODO 临时用的雪花算法 但我觉得这里用uuid比较好 顺便也加一下id-generator支持
	// TODO 后面真实情况下让用户自己用中间件生成traceid，然后我就在logger里获取!
	ctx.Set("X-Trace-Id", gs.SnowFlake.NextShortId())
	logger.PrintWithCtx(ctx, "hello world")
	logger.PrintWithCtx(ctx, "hello world")
}
