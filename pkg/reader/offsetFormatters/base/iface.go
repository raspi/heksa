package base

type OffsetFormatter interface {
	GetFormatWidth() int
	Print(offset uint64) string
}

// BaseInfo contains meta information about a file which is read for offset formatters
type BaseInfo struct {
	// File size hint for offset formatter(s), for example how many padding zeroes are needed when printing out position
	FileSize int64
}
