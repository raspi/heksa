package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"io"
	"strings"
)

type Dec struct {
	fs      uint8
	palette map[uint8]clr.Color
}

func (d Dec) SetBitWidthSize(s uint8) {
	d.fs = s
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

func (d Dec) leading(i int64) string {
	out := fmt.Sprintf(`%02x`, i)
	out = strings.Repeat(`0`, int(d.fs-2)-len(out)) + out
	return out
}

// DisplayOffset displays offset as decimal 0 - 9999999....
func (d Dec) DisplayOffset(r io.ReadSeeker) string {
	off, _ := r.Seek(0, io.SeekCurrent)
	return d.leading(off)
}

func (d *Dec) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}
