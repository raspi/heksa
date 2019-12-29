package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"io"
	"math/bits"
	"strings"
)

type Oct struct {
	fs        uint64 // File size
	bw        uint8  // Bit width calculated from file size
	palette   map[uint8]clr.Color
	offFormat string
}

func (d *Oct) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.bw = nearest(uint8(bits.Len64(d.fs)))
	d.offFormat = fmt.Sprintf(`%%0%vo`, d.bw)
}

func NewOct() *Oct {
	return &Oct{}
}

func (d Oct) Display(a []byte) string {
	out := ``
	for idx, b := range a {
		if idx == 8 {
			out += ` `
		}

		color, ok := d.palette[b]
		if !ok {
			color = clr.BrightFg
		}

		out += clr.Sprintf(`%03o `, clr.Colorize(b, color))
	}

	return strings.Trim(out, ` `)
}

// DisplayOffset displays offset as hexadecimal 0x00 - 0xFFFFFFFF....
func (d Oct) DisplayOffset(r io.ReadSeeker) string {
	off, _ := r.Seek(0, io.SeekCurrent)
	return fmt.Sprintf(d.offFormat, off)
}

func (d *Oct) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}
