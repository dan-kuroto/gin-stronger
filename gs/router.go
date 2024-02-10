package gs

import (
	"net/http"

	"github.com/dan-kuroto/gin-stronger/config"
	"github.com/gin-gonic/gin"
)

type HttpMethod uint16
type StaticMapFunc func() map[string]string

const (
	GET HttpMethod = 1 << iota
	HEAD
	POST
	PUT
	PATCH
	DELETE
	CONNECT
	OPTIONS
	TRACE
	Any = GET | HEAD | POST | PUT | PATCH | DELETE | CONNECT | OPTIONS | TRACE
)

var rootRouter = Router{Path: ""}
var staticMapFunc StaticMapFunc

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

func AddRouter(router ginEngineOrGroup, gsRouter *Router) {
	if len(gsRouter.Children) == 0 {
		handleRouter(router, gsRouter)
	} else {
		group := router.Group(gsRouter.Path)
		if len(gsRouter.MiddleWares) != 0 {
			group.Use(gsRouter.MiddleWares...)
		}
		for _, subRouter := range gsRouter.Children {
			AddRouter(group, &subRouter)
		}
	}
}

func InitStatic(engine *gin.Engine) {
	for urlPath, filePath := range staticMapFunc() {
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
func SetStatic(getter StaticMapFunc) {
	staticMapFunc = getter
}

func RunApp[T config.IConfiguration](config T) {
	PrintBanner()
	InitConfig(config)
	InitIdGenerators()

	if Config.GetGinRelease() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	AddRouter(engine, &rootRouter)
	InitStatic(engine)

	engine.Run(Config.GetGinAddr())
}

// It is shorthand for gs.RunApp(&gs.Configuration{})
func RunAppDefault() {
	RunApp(&config.Configuration{})
}
