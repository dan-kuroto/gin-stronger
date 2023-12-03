// test file
package main

import (
	"net/http"

	"github.com/dan-kuroto/gin-stronger/gs"
	"github.com/gin-gonic/gin"
)

func HandlerFuncDemo1(c *gin.Context) {
	c.JSON(http.StatusAccepted, map[string]any{"hello": "world"})
}

type HelloStruct struct {
	Hello string `json:"hello"`
}

func HandlerFuncDemo2(c *gin.Context) HelloStruct {
	return HelloStruct{Hello: "world"}
}

func GetRouters() []gs.Router {
	return []gs.Router{
		{
			Path: "/api",
			Children: []gs.Router{
				{
					Path:     "/test1",
					Handlers: gs.PackageHandlers(HandlerFuncDemo1),
				},
			},
		},
	}
}
