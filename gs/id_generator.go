package gs

import (
	"github.com/dan-kuroto/gin-stronger/generator"
)

var SnowFlake *generator.SnowFlakeGenerator

func InitIdGenerators() {
	SnowFlake = generator.NewSnowFlake(Config)
}
