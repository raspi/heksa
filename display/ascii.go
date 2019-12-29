package display

import (
	clr "github.com/logrusorgru/aurora"
	"strings"
)

type Ascii struct {
	palette map[uint8]clr.Color
}

func NewAscii() *Ascii {
	return &Ascii{}
}

func (d Ascii) Display(a []byte) string {
	out := ``
	for idx, b := range a {
		if idx == 8 {
			out += ` `
		}
		color, ok := d.palette[b]
		if !ok {
			color = clr.BrightFg
		}

		out += clr.Sprintf(`%c`, clr.Colorize(d.toChar(b), color))
	}

	return strings.Trim(out, ` `)
}

func (d *Ascii) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}

// non-printable characters as dot ('.')
func (d *Ascii) toChar(b byte) byte {
	if b < 32 || b > 126 {
		return '.'
	}
	return b
}
