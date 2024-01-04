package gs

import "fmt"

type keyType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string
}

type IdGenerator[WorkerKeyType keyType, IDType keyType] interface {
	NextId(worker WorkerKeyType) IDType
	NextStrId(worker WorkerKeyType) string
}

type snowFlakeWorker struct{}

func (w *snowFlakeWorker) NextId() int64 {
	// TODO ...
	return 0
}

type snowFlake struct {
	workers map[string]*snowFlakeWorker
}

var SnowFlake snowFlake

func (s *snowFlake) NextId(worker string) int64 {
	// TODO ...
	return s.workers[worker].NextId()
}

func (s *snowFlake) NextStrId(worker string) string {
	return fmt.Sprint(s.NextId(worker))
}

// TODO: UUID
