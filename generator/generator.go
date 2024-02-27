package generator

type IdType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string
}

type IdGenerator[T IdType] interface {
	NextId() T
}
