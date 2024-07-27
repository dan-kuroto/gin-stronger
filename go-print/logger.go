package gp

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	log.Logger

	formatter *Formatter
}

func NewLogger(pkgName string) *Logger {
	logger := new(Logger)
	logger.SetOutput(os.Stderr)
	if pkgName != "" {
		logger.SetPrefix("[" + pkgName + "] ")
	}
	logger.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	logger.formatter = &DefaultFormatter

	return logger
}

func (l *Logger) UseFormatter(formatter *Formatter) *Logger {
	l.formatter = formatter
	return l
}

func (l *Logger) ClearFormatter() *Logger {
	return l.UseFormatter(nil)
}

func (l *Logger) convertArgs(args []any) []any {
	if l.formatter != nil {
		for i, arg := range args {
			args[i] = l.formatter.ToString(arg)
		}
	}
	return args
}

func (l *Logger) getTraceId(ctx *gin.Context) string {
	if ctx != nil {
		if traceId := ctx.GetString("X-Trace-Id"); traceId != "" {
			return "[" + traceId + "]"
		}
	}
	return ""
}

// region overrides

func (l *Logger) Print(args ...any) {
	args = l.convertArgs(args)
	l.Output(2, fmt.Sprint(args...))
}

func (l *Logger) PrintWithCtx(ctx *gin.Context, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		args = append([]any{traceId + " "}, args...)
	}
	l.Output(2, fmt.Sprint(args...))
}

func (l *Logger) Println(args ...any) {
	args = l.convertArgs(args)
	l.Output(2, fmt.Sprintln(args...))
}

func (l *Logger) PrintlnWithCtx(ctx *gin.Context, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		args = append([]any{traceId}, args...)
	}
	l.Output(2, fmt.Sprintln(args...))
}

func (l *Logger) Printf(format string, args ...any) {
	args = l.convertArgs(args)
	l.Output(2, fmt.Sprintf(format, args...))
}

func (l *Logger) PrintfWithCtx(ctx *gin.Context, format string, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		format = traceId + " " + format
	}
	l.Output(2, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(args ...any) {
	args = l.convertArgs(args)
	l.Output(2, fmt.Sprint(args...))
	os.Exit(1)
}

func (l *Logger) FatalWithCtx(ctx *gin.Context, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		args = append([]any{traceId + " "}, args...)
	}
	l.Output(2, fmt.Sprint(args...))
	os.Exit(1)
}

func (l *Logger) Fatalln(args ...any) {
	args = l.convertArgs(args)
	l.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

func (l *Logger) FatallnWithCtx(ctx *gin.Context, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		args = append([]any{traceId}, args...)
	}
	l.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...any) {
	args = l.convertArgs(args)
	l.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) FatalfWithCtx(ctx *gin.Context, format string, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		format = traceId + " " + format
	}
	l.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) Panic(args ...any) {
	args = l.convertArgs(args)
	s := fmt.Sprint(args...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) PanicWithCtx(ctx *gin.Context, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		args = append([]any{traceId + " "}, args...)
	}
	s := fmt.Sprint(args...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) Panicln(args ...any) {
	args = l.convertArgs(args)
	s := fmt.Sprintln(args...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) PaniclnWithCtx(ctx *gin.Context, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		args = append([]any{traceId}, args...)
	}
	s := fmt.Sprintln(args...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) Panicf(format string, args ...any) {
	args = l.convertArgs(args)
	s := fmt.Sprintf(format, args...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) PanicfWithCtx(ctx *gin.Context, format string, args ...any) {
	args = l.convertArgs(args)
	if traceId := l.getTraceId(ctx); traceId != "" {
		format = traceId + " " + format
	}
	s := fmt.Sprintf(format, args...)
	l.Output(2, s)
	panic(s)
}

// regionend
