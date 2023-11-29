// test file
package main

import (
	"net/http"

	"github.com/dan-kuroto/gin-stronger/gs"
	"github.com/gin-gonic/gin"
)

func HandlerFuncDemo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
}

func GetRouterMap() gs.RouterMap {
	return gs.RouterMap{
		"/api": gs.RouterMap{
			"/hello": gs.RouterHandler(HandlerFuncDemo),
		},
	}
}
