package gs

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
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
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件路径
		MaxSize:    10,             // 每个日志文件的最大大小（MB）
		MaxBackups: 10,             // 保留旧日志文件的最大数量
		MaxAge:     30,             // 保留旧日志文件的最大天数
		Compress:   true,           // 是否压缩旧日志文件
	}
	defer logFile.Close()
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(
		zerolog.LevelWriterAdapter{Writer: os.Stderr},
		logFile,
	)).With().Timestamp().Caller().Logger()
}

// retrieves a logger from the gin.Context or creates a new one if it doesn't exist.
func GetLoggerByGinCtx(ctx *gin.Context) *zerolog.Logger {
	if ctx == nil {
		return &log.Logger
	}
	value, exists := ctx.Get("gs-logger")
	if exists {
		return value.(*zerolog.Logger)
	}
	logger := log.With().Str("traceID", shortuuid.New()).Logger()
	ctx.Set("gs-logger", &logger)
	return &logger
}
