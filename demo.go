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
	logger.Println("hello world")
	logger.Printf("hello %s", "world")
	ctx := &gin.Context{}
	ctx.Set("X-Trace-Id", gs.SnowFlake.NextShortId())
	logger.PrintWithCtx(ctx, "hello world")
	logger.PrintlnWithCtx(ctx, "hello world")
	logger.PrintfWithCtx(ctx, "hello %s", "world")
	logger.PanicfWithCtx(ctx, "hello %s", "world")
}
