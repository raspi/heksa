package iface

import (
	clr "github.com/logrusorgru/aurora"
	"io"
)

// ShowsOffset is interface for displaying file offset in X format (where X might be hex, decimal, octal, ..)
type ShowsOffset interface {
	DisplayOffset(r ReadSeekerCloser) string
	SetFileSize(int64) // For leading zeros information
}

type Views interface {
	Display([]byte) string
	SetPalette(map[uint8]clr.Color)
}

type ReadSeekerCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
