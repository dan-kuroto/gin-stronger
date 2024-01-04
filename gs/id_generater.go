package gs

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type IdType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string
}

type IdGenerator[T IdType] interface {
	NextId() T
}

const (
	// TODO: 把这个也改成yml文件配置,但是gs-cli自动配置为诗夜初次发动态时间,就可以了
	SNOW_FLAKE_START_STMP         int64 = 1626779686000
	SNOW_FLAKE_SEQUENCE_BIT       int64 = 12
	SNOW_FLAKE_MACHINE_BIT        int64 = 5
	SNOW_FLAKE_DATACENTER_BIT     int64 = 5
	SNOW_FLAKE_MAX_DATACENTER_NUM int64 = ^(-1 << SNOW_FLAKE_DATACENTER_BIT)
	SNOW_FLAKE_MAX_MACHINE_NUM    int64 = ^(-1 << SNOW_FLAKE_MACHINE_BIT)
	SNOW_FLAKE_MAX_SEQUENCE       int64 = ^(-1 << SNOW_FLAKE_SEQUENCE_BIT)
	SNOW_FLAKE_MACHINE_LEFT       int64 = SNOW_FLAKE_SEQUENCE_BIT
	SNOW_FLAKE_DATACENTER_LEFT    int64 = SNOW_FLAKE_SEQUENCE_BIT + SNOW_FLAKE_MACHINE_BIT
	SNOW_FLAKE_TIMESTMP_LEFT      int64 = SNOW_FLAKE_DATACENTER_LEFT + SNOW_FLAKE_DATACENTER_BIT
)

type snowFlake struct {
	// init: 0
	sequence int64
	// init: -1
	lastStmp int64

	mutex sync.Mutex
}

var SnowFlake snowFlake

func (s *snowFlake) NextId() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	snowflake := Config.GetSnowFlakeConfig()

	currStmp := getNewStmp()
	if currStmp < s.lastStmp {
		log.Println("Clock moved backwards. Refusing to generate id")
		os.Exit(1)
	} else if currStmp == s.lastStmp {
		s.sequence = (s.sequence + 1) & SNOW_FLAKE_MAX_SEQUENCE
		if s.sequence == 0 {
			currStmp = s.getNextMilli()
		}
	} else {
		s.sequence = 0
	}
	s.lastStmp = currStmp

	return (currStmp-SNOW_FLAKE_START_STMP)<<SNOW_FLAKE_TIMESTMP_LEFT |
		snowflake.DataCenterId<<SNOW_FLAKE_DATACENTER_LEFT |
		snowflake.MachineId<<SNOW_FLAKE_MACHINE_LEFT |
		s.sequence
}

func (s *snowFlake) NextStrId() string {
	return fmt.Sprint(s.NextId())
}

func (s *snowFlake) getNextMilli() int64 {
	mill := getNewStmp()
	for mill <= s.lastStmp {
		mill = getNewStmp()
	}
	return mill
}

func getNewStmp() int64 {
	return time.Now().UnixMilli()
}
