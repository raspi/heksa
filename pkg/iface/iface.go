package iface

import (
	clr "github.com/logrusorgru/aurora"
	"io"
)

// OffsetFormatter is interface for displaying file offset in X format (where X might be hex, decimal, octal, ..)
type OffsetFormatter interface {
	DisplayOffset(r ReadSeekerCloser) string
	SetFileSize(int64) // For leading zeros information
}

// CharacterFormatter displays bytes in X format
type CharacterFormatter interface {
	Format(b byte) string // Get the colorized representation
	SetPalette(map[uint8]clr.Color)
	EofStr() string // String if EOF has been reached. for lining output.
}

type ReadSeekerCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
