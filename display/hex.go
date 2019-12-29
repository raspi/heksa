package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"io"
	"strings"
)

type HexOffset struct {
	fs      uint8
	palette map[uint8]clr.Color
}

func (d HexOffset) SetBitWidthSize(s uint8) {
	d.fs = s
}

func NewHex() *HexOffset {
	return &HexOffset{
		fs: 8,
	}
}

func (d HexOffset) Display(a []byte) string {
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

func (d HexOffset) leading(i int64) string {
	out := fmt.Sprintf(`%02x`, i)
	out = strings.Repeat(`0`, int(d.fs-2)-len(out)) + out
	return out
}

func (d HexOffset) DisplayOffset(r io.ReadSeeker) string {
	off, _ := r.Seek(0, io.SeekCurrent)
	return d.leading(off)
}

func (d *HexOffset) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}
