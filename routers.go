// test file
package main

import (
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
	c.JSON(http.StatusAccepted, map[string]any{"hello": "world"})
}

func HandlerFuncDemo2(c *gin.Context) HelloStruct {
	return HelloStruct{Hello: "world"}
}

func HandlerFuncDemo3(hello HelloStruct) HelloStruct {
	return HelloStruct{Hello: hello.Hello}
}

func GetRouters() []gs.Router {
	return []gs.Router{
		{
			Path: "/api",
			Children: []gs.Router{
				{
					Path:     "/test1",
					Handlers: gs.PackageHandlers(HandlerFuncDemo, HandlerFuncDemo1),
				},
				{
					Path:     "/test2",
					Handlers: gs.PackageHandlers(HandlerFuncDemo2),
				},
				{
					Path:     "/test3",
					Handlers: gs.PackageHandlers(HandlerFuncDemo3),
				},
			},
		},
	}
}
