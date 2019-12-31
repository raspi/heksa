package iface

import (
	clr "github.com/logrusorgru/aurora"
	"io"
)

// OffsetFormatter is interface for displaying file offset in X format (where X might be hex, decimal, octal, ..)
type OffsetFormatter interface {
	FormatOffset(r ReadSeekerCloser) string
	SetFileSize(int64) // For leading zeros information
	OffsetHeader() string
}

// CharacterFormatter displays bytes in X format
type CharacterFormatter interface {
	Format(b byte, c clr.Color) string // Get the colorized representation
	EofStr() string                    // String if EOF has been reached. for lining output.
	Header() string
}

type ReadSeekerCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
