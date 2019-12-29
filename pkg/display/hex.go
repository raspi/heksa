package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
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

	if d.bw == 0 {
		d.bw = 8
	}

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

// DisplayOffset displays offset as hexadecimal 0x00 - 0xFFFFFFFF....
func (d Hex) DisplayOffset(r iface.ReadSeekerCloser) string {
	off, _ := r.Seek(0, io.SeekCurrent)
	return fmt.Sprintf(d.offFormat, off)
}

func (d *Hex) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}
