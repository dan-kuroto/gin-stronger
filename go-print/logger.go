package gp

import (
	"log"
	"os"
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

	// TODO: 在Println等情况真正用上gs.Formatter
	logger.formatter = &DefaultFormatter

	return logger
}

func (l *Logger) UseFormatter(formatter *Formatter) *Logger {
	l.formatter = formatter
	return l
}

func (l *Logger) UseStdFormatter() *Logger {
	l.formatter = nil
	return l
}
