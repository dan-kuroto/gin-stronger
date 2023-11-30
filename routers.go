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

func GetRouterMap() gs.RouterMap {
	return gs.RouterMap{
		"/api": gs.RouterMap{
			"/hello": gs.RouterHandlers{HandlerFuncDemo1, HandlerFuncDemo2},
		},
	}
}
