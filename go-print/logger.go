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

func (l *Logger) convertArgsWithCtx(ctx *gin.Context, args []any) []any {
	args = l.convertArgs(args)
	if ctx != nil {
		if traceId := ctx.GetString("X-Trace-Id"); traceId != "" {
			args = append([]any{"[" + traceId + "] "}, args...)
		}
	}
	return args
}

// region overrides

func (l *Logger) PrintWithCtx(ctx *gin.Context, v ...any) {
	l.Output(2, fmt.Sprint(l.convertArgsWithCtx(ctx, v)...))
}

func (l *Logger) Print(v ...any) {
	l.Output(2, fmt.Sprint(l.convertArgs(v)...))
}

func (l *Logger) Println(v ...any) {
	l.Output(2, fmt.Sprintln(l.convertArgs(v)...))
}

func (l *Logger) Printf(format string, v ...any) {
	l.Output(2, fmt.Sprintf(format, l.convertArgs(v)...))
}

func (l *Logger) Fatal(v ...any) {
	l.Output(2, fmt.Sprint(l.convertArgs(v)...))
	os.Exit(1)
}

func (l *Logger) Fatalln(v ...any) {
	l.Output(2, fmt.Sprintln(l.convertArgs(v)...))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.Output(2, fmt.Sprintf(format, l.convertArgs(v)...))
	os.Exit(1)
}

func (l *Logger) Panic(v ...any) {
	s := fmt.Sprint(l.convertArgs(v)...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) Panicln(v ...any) {
	s := fmt.Sprintln(l.convertArgs(v)...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, l.convertArgs(v)...)
	l.Output(2, s)
	panic(s)
}

// regionend
