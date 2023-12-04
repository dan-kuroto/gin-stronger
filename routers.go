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

func HandlerFuncDemo1(hello *HelloStruct) HelloStruct {
	panic("panic demo~")
	return HelloStruct{Hello: hello.Hello}
}

func HandlerFuncDemo2(c *gin.Context, hello HelloStruct) HelloStruct {
	return HelloStruct{Hello: hello.Hello}
}
func HandlerFuncDemoLog() {
	fmt.Println("log ...")
}

func GetRouters() []gs.Router {
	return []gs.Router{
		{
			Path: "/api",
			Children: []gs.Router{
				{
					Path:   "/test1",
					Method: gs.GET | gs.POST,
					Handlers: gs.PackageHandlers(
						HandlerFuncDemoLog,
						HandlerFuncDemo1,
						HandlerFuncDemoLog,
					),
				},
				{
					Path:     "/test2",
					Method:   gs.GET | gs.POST,
					Handlers: gs.PackageHandlers(HandlerFuncDemo2),
				},
			},
		},
	}
}
