package gs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpMethod uint16

const (
	GET     HttpMethod = 0b000000001
	HEAD    HttpMethod = 0b000000010
	POST    HttpMethod = 0b000000100
	PUT     HttpMethod = 0b000001000
	PATCH   HttpMethod = 0b000010000
	DELETE  HttpMethod = 0b000100000
	CONNECT HttpMethod = 0b001000000
	OPTIONS HttpMethod = 0b010000000
	TRACE   HttpMethod = 0b100000000
	Any     HttpMethod = 0b111111111
)

type Router struct {
	Path string
	// default value is gs.GET
	Method   HttpMethod
	Handlers []gin.HandlerFunc
	Children []Router
}

type ginEngineOrGroup interface {
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
	Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}

func handleRouter(router ginEngineOrGroup, gsRouter *Router) {
	if gsRouter.Method == Any {
		router.Any(gsRouter.Path, gsRouter.Handlers...)
	} else if gsRouter.Method == 0 {
		router.GET(gsRouter.Path, gsRouter.Handlers...)
	} else {
		if gsRouter.Method&GET != 0 {
			router.Handle(http.MethodGet, gsRouter.Path, gsRouter.Handlers...)
		}
		if gsRouter.Method&HEAD != 0 {
			router.Handle(http.MethodHead, gsRouter.Path, gsRouter.Handlers...)
		}
		if gsRouter.Method&POST != 0 {
			router.Handle(http.MethodPost, gsRouter.Path, gsRouter.Handlers...)
		}
		if gsRouter.Method&PUT != 0 {
			router.Handle(http.MethodPut, gsRouter.Path, gsRouter.Handlers...)
		}
		if gsRouter.Method&PATCH != 0 {
			router.Handle(http.MethodPatch, gsRouter.Path, gsRouter.Handlers...)
		}
		if gsRouter.Method&DELETE != 0 {
			router.Handle(http.MethodDelete, gsRouter.Path, gsRouter.Handlers...)
		}
		if gsRouter.Method&CONNECT != 0 {
			router.Handle(http.MethodConnect, gsRouter.Path, gsRouter.Handlers...)
		}
		if gsRouter.Method&OPTIONS != 0 {
			router.Handle(http.MethodOptions, gsRouter.Path, gsRouter.Handlers...)
		}
		if gsRouter.Method&TRACE != 0 {
			router.Handle(http.MethodTrace, gsRouter.Path, gsRouter.Handlers...)
		}
	}
}

func UseRouter(router ginEngineOrGroup, gsRouter *Router) {
	if len(gsRouter.Children) == 0 {
		handleRouter(router, gsRouter)
	} else {
		group := router.Group(gsRouter.Path)
		for _, subRouter := range gsRouter.Children {
			UseRouter(group, &subRouter)
		}
	}
}

func UseRouters(router ginEngineOrGroup, gsRouters []Router) {
	for _, gsRouter := range gsRouters {
		UseRouter(router, &gsRouter)
	}
}

// TODO: 1. 利用泛型机制使参数和返回值支持直接为结构体(就像SpringBoot一样)
// TODO: 2. 测了下gin有自带的panic recover机制，查一下能不能像SpringBoot一样自己加拦截器

func RegisterStatic(router *gin.Engine, staticMap map[string]string) {
	for urlPath, filePath := range staticMap {
		router.Static(urlPath, filePath)
	}
}
