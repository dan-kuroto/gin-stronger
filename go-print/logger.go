package gp

import (
	"fmt"
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

func (l *Logger) convertArgs(v []any) []any {
	if l.formatter != nil {
		for i, vv := range v {
			v[i] = l.formatter.ToString(vv)
		}
	}
	return v
}

// region overrides

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
