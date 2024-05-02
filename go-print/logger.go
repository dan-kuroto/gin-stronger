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

func (l *Logger) Print(v ...any) {
	l.Logger.Print(l.convertArgs(v)...)
}

func (l *Logger) Println(v ...any) {
	l.Logger.Println(l.convertArgs(v)...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.Logger.Printf(format, l.convertArgs(v)...)
}

func (l *Logger) Fatal(v ...any) {
	l.Logger.Fatal(l.convertArgs(v)...)
}

func (l *Logger) Fatalln(v ...any) {
	l.Logger.Fatalln(l.convertArgs(v)...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.Logger.Fatalf(format, l.convertArgs(v)...)
}

func (l *Logger) Panic(v ...any) {
	l.Logger.Panic(l.convertArgs(v)...)
}

func (l *Logger) Panicln(v ...any) {
	l.Logger.Panicln(l.convertArgs(v)...)
}

func (l *Logger) Panicf(format string, v ...any) {
	l.Logger.Panicf(format, l.convertArgs(v)...)
}
