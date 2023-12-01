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

func GetRouters() []gs.Router {
	return []gs.Router{
		{
			Path: "/api",
			Children: []gs.Router{
				{
					Path:     "/hello",
					Methods:  []gs.Method{gs.GET, gs.POST},
					Handlers: []gin.HandlerFunc{HandlerFuncDemo1, HandlerFuncDemo2},
				},
			},
		},
	}
}
