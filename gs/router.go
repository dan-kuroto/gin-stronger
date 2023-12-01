package gs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: 这里的使用方法还不够好，另外注释要改成英文

type Method string

const (
	GET     Method = http.MethodGet
	HEAD    Method = http.MethodHead
	POST    Method = http.MethodPost
	PUT     Method = http.MethodPut
	PATCH   Method = http.MethodPatch
	DELETE  Method = http.MethodDelete
	CONNECT Method = http.MethodConnect
	OPTIONS Method = http.MethodOptions
	TRACE   Method = http.MethodTrace
)

type Router struct {
	Path     string
	Methods  []Method
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

type response struct {
	Code int
	Msg  string
	Data any
}

// 默认的c.JSON能被多次调用,经封装后保证只会调用一次,返回{"msg": "ok", "data": data}
func Return(data any) {
	panic(response{Code: http.StatusOK, Msg: "ok", Data: data})
}

// 默认的c.JSON能被多次调用,经封装后保证只会调用一次
func Throw(code int, msg string, data any) {
	panic(response{Code: code, Msg: msg, Data: data})
}

// TODO: 重构一下,加上允许提供msg的功能,否则非throw情况下msg屁用没有,刚好可以给新增数据用
// Return的更省事调用方法: for string (format)
func ReturnS(format string, a ...any) {
	if len(a) == 0 {
		Return(format)
	} else {
		Return(fmt.Sprintf(format, a...))
	}
}

// Error的更省事调用法: for error
func ThrowE(err error) {
	Throw(http.StatusInternalServerError, err.Error(), nil)
}

// Error的更省事调用法: for string
func ThrowS(format string, a ...any) {
	if len(a) == 0 {
		Throw(http.StatusInternalServerError, format, nil)
	} else {
		Throw(http.StatusInternalServerError, fmt.Sprintf(format, a...), nil)
	}
}

// 包装一下handler，这样一来出错时直接panic就行了
func packageHandlerFunc(router *Router) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// 如果没有panic，认为用户使用了自己的处理逻辑，不做处理
			if r := recover(); r != nil {
				switch r := r.(type) {
				case response:
					c.JSON(r.Code, gin.H{
						"msg":  r.Msg,
						"data": r.Data,
					})
				default: // 用户自己用panic也可以,但默认code是StatusOK了
					c.JSON(http.StatusOK, r)
				}
			}
		}()
		for _, handler := range router.Handlers {
			handler(c)
		}
	}
}

// TODO: 魔改返回机制
func UseRouter(router ginEngineOrGroup, gsRouter *Router) {
	if len(gsRouter.Children) == 0 {
		for _, method := range gsRouter.Methods {
			router.Handle(string(method), gsRouter.Path, gsRouter.Handlers...)
		}
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

func RegisterStatic(router *gin.Engine, staticMap map[string]string) {
	for urlPath, filePath := range staticMap {
		router.Static(urlPath, filePath)
	}
}
