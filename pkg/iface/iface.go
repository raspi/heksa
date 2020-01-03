package iface

import (
	"io"
)

type ReadSeekerCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
