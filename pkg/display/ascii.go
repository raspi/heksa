package display

import (
	clr "github.com/logrusorgru/aurora"
	"strings"
)

type Ascii struct {
	palette map[uint8]clr.Color
	sb      strings.Builder
}

func NewAscii() *Ascii {
	return &Ascii{
		sb: strings.Builder{},
	}
}

func (d *Ascii) Format(b byte) string {
	d.sb.Reset()

	color, ok := d.palette[b]
	if !ok {
		color = clr.BrightFg
	}

	d.sb.WriteString(clr.Sprintf(`%c`, clr.Colorize(d.toChar(b), color)))
	return d.sb.String()

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

func (d *Ascii) EofStr() string {
	return ` `
}
