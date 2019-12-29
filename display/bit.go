package display

import (
	clr "github.com/logrusorgru/aurora"
	"strings"
)

type Bit struct {
	palette map[uint8]clr.Color
}

func NewBit() *Bit {
	return &Bit{}
}

func (d Bit) Display(a []byte) string {
	out := ``
	for idx, b := range a {
		if idx == 8 {
			out += ` `
		}

		color, ok := d.palette[b]
		if !ok {
			color = clr.BrightFg
		}

		out += clr.Sprintf(`%08b `, clr.Colorize(b, color))
	}

	return strings.Trim(out, ` `)
}

func (d *Bit) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}
