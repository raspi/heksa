package display

import (
	"fmt"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"strings"
)

type Percent struct {
	fs uint64 // File size
	sb strings.Builder
}

func NewPercent() *Percent {
	return &Percent{
		fs: 0,
		sb: strings.Builder{},
	}
}

func (d *Percent) SetFileSize(s int64) {
	d.fs = uint64(s)
}

// FormatOffset displays offset as percentage 0% - 100%
func (d *Percent) FormatOffset(r iface.ReadSeekerCloser) string {
	if d.fs == 0 {
		// No clue when file size is zero
		// it is a stream from stdin probably
		return `?%`
	}

	d.sb.Reset()
	off, _ := r.Seek(0, io.SeekCurrent)
	percent := (float64(off) * 100.0) / float64(d.fs)
	d.sb.WriteString(fmt.Sprintf(`%07.3f%%`, percent))
	return d.sb.String()
}

func (d *Percent) OffsetHeader() string {
	return `________`
}
