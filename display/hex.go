package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"io"
	"math/bits"
	"strings"
)

type Hex struct {
	fs        uint64 // File size
	bw        uint8  // Bit width calculated from file size
	palette   map[uint8]clr.Color
	offFormat string
}

func (d *Hex) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.bw = nearest(uint8(bits.Len64(d.fs)))
	d.offFormat = fmt.Sprintf(`%%0%vx`, d.bw)
}

func NewHex() *Hex {
	return &Hex{}
}

func (d Hex) Display(a []byte) string {
	out := ``
	for idx, b := range a {
		if idx == 8 {
			out += ` `
		}

		color, ok := d.palette[b]
		if !ok {
			color = clr.BrightFg
		}

		out += clr.Sprintf(`%02x `, clr.Colorize(b, color))
	}

	return strings.Trim(out, ` `)
}

func (d Hex) leading(i int64) string {
	out := fmt.Sprintf(`%02x`, i)
	out = strings.Repeat(`0`, int(d.fs-2)-len(out)) + out
	return out
}

// DisplayOffset displays offset as hexadecimal 0x00 - 0xFFFFFFFF....
func (d Hex) DisplayOffset(r io.ReadSeeker) string {
	off, _ := r.Seek(0, io.SeekCurrent)
	return fmt.Sprintf(d.offFormat, off)
}

func (d *Hex) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}
