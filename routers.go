// test file
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dan-kuroto/gin-stronger/gs"
	"github.com/gin-gonic/gin"
)

type HelloStruct struct {
	Hello string `json:"hello" form:"hello"`
}

func HandlerFuncDemo() {
	log.Println("api")
}

func HandlerFuncDemo1(c *gin.Context) {
	hello := HelloStruct{}
	if err := c.ShouldBind(&hello); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	c.JSON(http.StatusOK, hello)
}

func HandlerFuncDemo2(c *gin.Context) HelloStruct {
	return HelloStruct{Hello: "world"}
}

func HandlerFuncDemo3(c *gin.Context, hello HelloStruct) HelloStruct {
	return HelloStruct{Hello: hello.Hello}
}

func HandlerFuncDemo4(hello *HelloStruct) HelloStruct {
	return HelloStruct{Hello: hello.Hello}
}

func GetRouters() []gs.Router {
	return []gs.Router{
		{
			Path: "/api",
			Children: []gs.Router{
				{
					Path:     "/test1",
					Method:   gs.GET | gs.POST,
					Handlers: gs.PackageHandlers(HandlerFuncDemo, HandlerFuncDemo1),
				},
				{
					Path:     "/test2",
					Method:   gs.GET | gs.POST,
					Handlers: gs.PackageHandlers(HandlerFuncDemo2),
				},
				{
					Path:     "/test3",
					Method:   gs.GET | gs.POST,
					Handlers: gs.PackageHandlers(HandlerFuncDemo3),
				},
				{
					Path:     "/test4",
					Method:   gs.GET | gs.POST,
					Handlers: gs.PackageHandlers(HandlerFuncDemo4),
				},
			},
		},
	}
}
