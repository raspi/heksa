package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"math/bits"
)

type Percent struct {
	fs      uint64 // File size
	bw      uint8  // Bit width calculated from file size
	palette map[uint8]clr.Color
}

func (d *Percent) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.bw = nearest(uint8(bits.Len64(d.fs)))
}

func NewPercent() *Percent {
	return &Percent{
		fs: 8,
	}
}

// DisplayOffset displays offset as percentage 0% - 100%
func (d Percent) DisplayOffset(r iface.ReadSeekerCloser) string {
	off, _ := r.Seek(0, io.SeekCurrent)
	percent := float64(off) * 100.0 / float64(d.fs)
	return fmt.Sprintf(`%07.3f`, percent)
}
