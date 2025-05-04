package gs

import (
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		funcObj := runtime.FuncForPC(pc)
		if funcObj == nil {
			return filepath.Base(file) + ":" + strconv.Itoa(line)
		}
		return funcObj.Name() + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
}

// retrieves a logger from the gin.Context or creates a new one if it doesn't exist.
func GetLoggerByGinCtx(ctx *gin.Context) *zerolog.Logger {
	value, exists := ctx.Get("gs-logger")
	if exists {
		return value.(*zerolog.Logger)
	}
	logger := log.With().Str("traceID", shortuuid.New()).Logger()
	ctx.Set("gs-logger", &logger)
	return &logger
}
