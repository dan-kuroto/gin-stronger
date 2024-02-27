package gs

import (
	gp "github.com/dan-kuroto/gin-stronger/go-print"
)

var formatter = gp.Formatter{BracketColor: true}

func ToString(data any) string {
	return formatter.ToString(data)
}
