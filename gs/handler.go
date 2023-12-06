package gs

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

var ginContextType = reflect.TypeOf(&gin.Context{})

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

func callFunction(funcValue reflect.Value, inputs ...reflect.Value) []any {
	output := funcValue.Call(inputs)

	results := make([]any, 0, len(output))
	for _, value := range output {
		results = append(results, value.Interface())
	}
	return results
}

func packageHandler(function any, paramTypes []reflect.Type, resultTypes []reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make([]reflect.Value, 0, len(paramTypes))
		for _, paramType := range paramTypes {
			if paramType == ginContextType {
				params = append(params, reflect.ValueOf(c))
			} else {
				var param reflect.Value
				if paramType.Kind() == reflect.Ptr {
					param = reflect.New(paramType.Elem())
				} else {
					param = reflect.New(paramType)
				}
				if err := c.ShouldBind(param.Interface()); err != nil {
					panic(err)
				}
				if paramType.Kind() == reflect.Ptr {
					params = append(params, param)
				} else {
					params = append(params, param.Elem())
				}
			}
		}
		results := callFunction(reflect.ValueOf(function), params...)
		if len(results) == 1 {
			c.JSON(http.StatusOK, results[0])
		}
	}
}

// functions need to meet some conditions:
// (1) Parameters can include *gin.Context and request struct.
// (2) No result or return gs.IResponse
func PackageHandlers(functions ...any) []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, 0, len(functions))
	for _, function := range functions {
		funcType := reflect.TypeOf(function)
		paramTypes := getFunctionParamTypes(funcType)
		resultTypes := getFunctionResultTypes(funcType)

		if len(paramTypes) > 2 {
			panic("function parameter type is not supported")
		}
		if len(resultTypes) > 1 {
			panic("function result type is not supported")
		}
		if len(paramTypes) == 2 {
			if (paramTypes[0] != ginContextType && paramTypes[1] != ginContextType) ||
				(paramTypes[0] == ginContextType && paramTypes[1] == ginContextType) {
				panic("function parameter type is not supported")
			}
		}
		// if function is gin.HandlerFunc, packaging is unnecessary
		if len(paramTypes) == 1 && len(resultTypes) == 0 && paramTypes[0] == ginContextType {
			handlers = append(handlers, gin.HandlerFunc(function.(func(*gin.Context))))
		} else {
			handlers = append(handlers, packageHandler(function, paramTypes, resultTypes))
		}
	}
	return handlers
}

// Register recovery function for specific type.
//
// 'handlers' element must be func(c *gin.Context, err T), T is a certain type.
//
// ATTENTION: panic handler is based on defer&recover, so it's always the last
// to be called, regardless of the order in 'gs.Router.MiddleWares'.
func PackagePanicHandler(handlers ...any) gin.HandlerFunc {
	errTypes := make([]reflect.Type, 0, len(handlers))
	for _, handler := range handlers {
		errTypes = append(errTypes, getFunctionParamTypes(reflect.TypeOf(handler))[1])
	}
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		errValue := reflect.ValueOf(err)
		for i, errType := range errTypes {
			if errValue.CanConvert(errType) {
				handlerValue := reflect.ValueOf(handlers[i])
				callFunction(handlerValue, reflect.ValueOf(c), errValue.Convert(errType))
				break
			}
		}
		c.Next()
	})
}
