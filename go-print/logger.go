package gp

import (
	"github.com/gin-gonic/gin"
)

func getTraceId(ctx *gin.Context) string {
	if ctx != nil {
		if traceId := ctx.GetString("X-Trace-Id"); traceId != "" {
			return "[" + traceId + "]"
		}
	}
	return ""
}
