package generator

import (
	"strconv"
	"sync"
	"time"

	"github.com/dan-kuroto/gin-stronger/config"
)

const (
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

type SnowFlakeGenerator struct {
	// init: 0
	sequence int64
	// init: -1
	lastStmp int64

	mutex  sync.Mutex
	config config.IConfiguration
}

func NewSnowFlake(config config.IConfiguration) *SnowFlakeGenerator {
	return &SnowFlakeGenerator{
		lastStmp: -1,
		config:   config,
	}
}

func (s *SnowFlakeGenerator) NextId() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	snowflake := s.config.GetSnowFlakeConfig()

	currStmp := s.getCurrStmp()
	if currStmp < s.lastStmp {
		panic("Clock moved backwards. Refusing to generate id!")
	} else if currStmp == s.lastStmp {
		s.sequence = (s.sequence + 1) & SNOW_FLAKE_MAX_SEQUENCE
		if s.sequence == 0 {
			currStmp = s.getNextMilli()
		}
	} else {
		s.sequence = 0
	}
	s.lastStmp = currStmp

	return (currStmp-snowflake.StartStmp)<<SNOW_FLAKE_TIMESTMP_LEFT |
		snowflake.DataCenterId<<SNOW_FLAKE_DATACENTER_LEFT |
		snowflake.MachineId<<SNOW_FLAKE_MACHINE_LEFT |
		s.sequence
}

func (s *SnowFlakeGenerator) NextStrId() string {
	return strconv.FormatInt(s.NextId(), 10)
}

func (s *SnowFlakeGenerator) NextShortId() string {
	return strconv.FormatInt(s.NextId(), 36)
}

func (s *SnowFlakeGenerator) getNextMilli() int64 {
	mill := s.getCurrStmp()
	for mill <= s.lastStmp {
		mill = s.getCurrStmp()
	}
	return mill
}

func (s *SnowFlakeGenerator) getCurrStmp() int64 {
	return time.Now().UnixMilli()
}
