package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"math/bits"
	"strings"
)

type Dec struct {
	fs        uint64
	bw        uint8 // Bit width calculated from file size
	palette   map[uint8]clr.Color
	offFormat string
}

func (d *Dec) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.bw = nearest(uint8(bits.Len64(d.fs)))
	d.offFormat = fmt.Sprintf(`%%0%vd`, d.bw)
}

func NewDec() *Dec {
	return &Dec{
		fs: 8,
	}
}

func (d Dec) Display(a []byte) string {
	out := ``
	for idx, b := range a {
		if idx == 8 {
			out += ` `
		}

		color, ok := d.palette[b]
		if !ok {
			color = clr.BrightFg
		}

		out += clr.Sprintf(`%03d `, clr.Colorize(b, color))
	}

	return strings.Trim(out, ` `)
}

// DisplayOffset displays offset as decimal 0 - 9999999....
func (d Dec) DisplayOffset(r iface.ReadSeekerCloser) string {
	off, _ := r.Seek(0, io.SeekCurrent)
	return fmt.Sprintf(d.offFormat, off)
}

func (d *Dec) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}
