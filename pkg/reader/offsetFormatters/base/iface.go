package base

type OffsetFormatter interface {
	GetFormatWidth() int
	Print(offset uint64) string
}

type BaseInfo struct {
	FileSize int64
}
