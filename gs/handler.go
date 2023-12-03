package gs

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

var ginContextType = reflect.TypeOf(&gin.Context{})

type IResponse interface {
	ToMap() map[string]any
}

// TODO: 1. 利用泛型机制使参数和返回值支持直接为结构体(就像SpringBoot一样)
// TODO: 2. 测了下gin有自带的panic recover机制，查一下能不能像SpringBoot一样自己加拦截器

func getFunctionParamTypes(funcType reflect.Type) []reflect.Type {
	numIn := funcType.NumIn()
	types := make([]reflect.Type, 0, numIn)
	for i := 0; i < numIn; i++ {
		types = append(types, funcType.In(i))
	}
	return types
}

func getFunctionResultTypes(funcType reflect.Type) []reflect.Type {
	numOut := funcType.NumOut()
	types := make([]reflect.Type, 0, numOut)
	for i := 0; i < numOut; i++ {
		types = append(types, funcType.Out(i))
	}
	return types
}

// functions:
// Parameters can include *gin.Context and request struct.
// No result or return gs.IResponse
func PackageHandlers(functions ...any) []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, 0, len(functions))
	for _, function := range functions {
		funcType := reflect.TypeOf(function)
		paramTypes := getFunctionParamTypes(funcType)
		resultTypes := getFunctionResultTypes(funcType)

		// if function is gin.HandlerFunc, packaging is unnecessary
		if len(paramTypes) == 1 && len(resultTypes) == 0 && paramTypes[0] == ginContextType {
			handlers = append(handlers, gin.HandlerFunc(function.(func(*gin.Context))))
			continue
		}
	}
	return handlers
}
