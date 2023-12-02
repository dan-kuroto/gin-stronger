// test file
package main

import (
	"fmt"
	"net/http"

	"github.com/dan-kuroto/gin-stronger/gs"
	"github.com/gin-gonic/gin"
)

func HandlerFuncDemo1(c *gin.Context) {
	fmt.Println("demo1")
}

func HandlerFuncDemo2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
}

func HandlerFuncDemo3(c *gin.Context) {
	panic("test")
}

func GetRouters() []gs.Router {
	return []gs.Router{
		{
			Path: "/api",
			Children: []gs.Router{
				{
					Path:     "/hello1",
					Method:   gs.GET | gs.POST,
					Handlers: []gin.HandlerFunc{HandlerFuncDemo1, HandlerFuncDemo2},
				},
				{
					Path:     "/hello2",
					Method:   gs.GET | gs.HEAD | gs.POST | gs.PUT | gs.PATCH | gs.DELETE | gs.CONNECT | gs.OPTIONS | gs.TRACE,
					Handlers: []gin.HandlerFunc{HandlerFuncDemo1, HandlerFuncDemo2},
				},
			},
		},
		{
			Path:     "/test1",
			Method:   gs.Any,
			Handlers: []gin.HandlerFunc{HandlerFuncDemo2},
		},
		{
			Path:     "/test2",
			Handlers: []gin.HandlerFunc{HandlerFuncDemo2},
		},
		{
			Path:     "/test3",
			Handlers: []gin.HandlerFunc{HandlerFuncDemo3},
		},
	}
}
