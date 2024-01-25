package main

import (
	"github.com/dan-kuroto/gin-stronger/check"
	"github.com/dan-kuroto/gin-stronger/gs"
)

func main() {
	var a []string
	gs.CheckParam("name", a, check.NotEmptySlice)
}
