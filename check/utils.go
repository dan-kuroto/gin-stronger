package check

import (
	gp "github.com/dan-kuroto/gin-stronger/go-print"
)

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type orderable interface {
	number | rune | string
}

var formatter = gp.Formatter{}

func getLength(data any) (length int, ok bool) {
	// TODO
}
