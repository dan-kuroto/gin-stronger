// test file
package main

import (
	"fmt"

	"github.com/dan-kuroto/gin-stronger/gs"
	"github.com/gin-gonic/gin"
)

type HelloStruct struct {
	Hello string `json:"hello" form:"hello"`
}

func HandlerDemo1(hello *HelloStruct) HelloStruct {
	panic("panic demo~")
	return HelloStruct{Hello: hello.Hello}
}

func HandlerDemo2(c *gin.Context, hello HelloStruct) HelloStruct {
	panic(1)
	return HelloStruct{Hello: hello.Hello}
}

func LogHandlerStart(c *gin.Context) {
	fmt.Println("log start ...")
}

func LogHandlerEnd(c *gin.Context) {
	fmt.Println("log end ...")
}

func PanicStringHandler(c *gin.Context, err string) {
	c.JSON(500, gin.H{"errMessage": err})
}

func PanicAnyHandler(c *gin.Context, err any) {
	c.JSON(500, gin.H{"errMessage": "unknown error"})
}

func PanicIntHandler(c *gin.Context, err int) {
	c.JSON(500, gin.H{"errMessage": err})
}

func GetRouters() []gs.Router {
	return []gs.Router{
		{
			Path: "/api",
			MiddleWares: []gin.HandlerFunc{
				LogHandlerStart,
				gs.PackagePanicHandler(
					PanicStringHandler,
					PanicAnyHandler,
				),
				LogHandlerEnd,
			},
			Children: []gs.Router{
				{
					Path:     "/test1",
					Method:   gs.GET | gs.POST,
					Handlers: gs.PackageHandlers(HandlerDemo1),
				},
			},
		},
		{
			Path: "/api2",
			MiddleWares: []gin.HandlerFunc{
				gs.PackagePanicHandler(PanicIntHandler),
				func(c *gin.Context) {
					fmt.Println("test")
				},
			},
			Children: []gs.Router{
				{
					Path:     "/test2",
					Method:   gs.GET | gs.POST,
					Handlers: gs.PackageHandlers(HandlerDemo2),
				},
			},
		},
	}
}
