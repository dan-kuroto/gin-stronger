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

var rootRouter = Router{Path: ""}
var staticMap = make(map[string]string)

type Router struct {
	Path string
	// invalid for router group. default value is gs.GET
	Method HttpMethod
	// valid for router group
	MiddleWares []gin.HandlerFunc
	// invalid for router group
	Handlers []gin.HandlerFunc
	// if len(Children) != 0, it is a router group
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
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
}

type Controller interface {
	GetRouter() Router
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

func addRouter(router ginEngineOrGroup, gsRouter *Router) {
	if len(gsRouter.Children) == 0 {
		handleRouter(router, gsRouter)
	} else {
		group := router.Group(gsRouter.Path)
		if len(gsRouter.MiddleWares) != 0 {
			group.Use(gsRouter.MiddleWares...)
		}
		for _, subRouter := range gsRouter.Children {
			addRouter(group, &subRouter)
		}
	}
}

func initStatic(engine *gin.Engine) {
	for urlPath, filePath := range staticMap {
		engine.Static(urlPath, filePath)
	}
}

// Set global URL preffix.
// Has no effect on the prefix of `RegisterStatic`.
func SetGlobalPreffix(preffix string) {
	rootRouter.Path = preffix
}

func AddGlobalMiddleware(middlewares ...gin.HandlerFunc) {
	rootRouter.MiddleWares = append(rootRouter.MiddleWares, middlewares...)
}

func UseController(controller Controller) {
	rootRouter.Children = append(rootRouter.Children, controller.GetRouter())
}

// Register static files. `url2path` is the mapping of url to file path
// (directory path is supported).
//
// It won't be affected by `SetGlobalPreffix`.
func RegisterStatic(url2path map[string]string) {
	for urlPath, filePath := range url2path {
		staticMap[urlPath] = filePath
	}
}

func RunApp[T IConfiguration](config T) {
	printBanner()
	initConfig(config)

	if Config.GetGinRelease() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	addRouter(engine, &rootRouter)
	initStatic(engine)

	engine.Run(Config.GetGinAddr())
}

// It is shorthand for gs.RunApp(&gs.Configuration{})
func RunAppDefault() {
	RunApp(&Configuration{})
}
