package gs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: 这里的使用方法还不够好，另外注释要改成英文

type RouterMapValue interface {
	isRouterMapValue()
}

// TODO: 果然还是要魔改，用vue-router那种搞法，直接用children来解决
type RouterHandler gin.HandlerFunc

func (RouterHandler) isRouterMapValue() {}

type RouterHandlers []RouterHandler

func (RouterHandlers) isRouterMapValue() {}

type RouterMap map[string]RouterMapValue

func (RouterMap) isRouterMapValue() {}

// Initial rootPath must be empty string
func walkRouterMap(routerMap RouterMap, rootPath string, callback func(rootPath string, subPath string, value RouterHandler)) {
	for subPath, value := range routerMap {
		switch value := value.(type) {
		case RouterHandler:
			callback(rootPath, subPath, value)
		case RouterHandlers:
			for _, handler := range value {
				callback(rootPath, subPath, handler)
			}
		case RouterMap:
			walkRouterMap(value, rootPath+subPath, callback)
		}
	}
}

func WalkRouterMap(routerMap RouterMap, callback func(rootPath string, subPath string, value RouterHandler)) {
	walkRouterMap(routerMap, "", callback)
}

// 用于RegisterRouterMap的辅助接口，使其能够同时支持gin.Engine和gin.RouterGroup
type ginEngineOrGroup interface {
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}

// TODO: 这里是否有必要用泛型？
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
func packageHandlerFunc(handler RouterHandler) gin.HandlerFunc {
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
		handler(c)
	}
}

func RegisterRouterMap(router ginEngineOrGroup, routerMap RouterMap) {
	for path, value := range routerMap {
		switch value := value.(type) {
		case RouterHandler:
			// TODO: 目前是强制GET+POST全部注册，以后要加上method字段控制注册哪个
			router.GET(path, packageHandlerFunc(value))
			router.POST(path, packageHandlerFunc(value))
		case RouterHandlers:
			handlers := make([]gin.HandlerFunc, 0, len(value))
			for _, handler := range value {
				handlers = append(handlers, packageHandlerFunc(handler))
			}
			router.GET(path, handlers...)
			router.POST(path, handlers...)
		case RouterMap:
			group := router.Group(path)
			RegisterRouterMap(group, value)
		}
	}
}

func RegisterStatic(router *gin.Engine, staticMap map[string]string) {
	for urlPath, filePath := range staticMap {
		router.Static(urlPath, filePath)
	}
}
