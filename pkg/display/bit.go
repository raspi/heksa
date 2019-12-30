package display

import (
	clr "github.com/logrusorgru/aurora"
	"strings"
)

type Bit struct {
	palette map[uint8]clr.Color
	sb      strings.Builder
}

func NewBit() *Bit {
	return &Bit{
		sb: strings.Builder{},
	}
}

func (d *Bit) Format(b byte) string {
	d.sb.Reset()

	color, ok := d.palette[b]
	if !ok {
		color = clr.BrightFg
	}

	d.sb.WriteString(clr.Sprintf(`%08b `, clr.Colorize(b, color)))

	return d.sb.String()
}

func (d *Bit) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}

func (d *Bit) EofStr() string {
	return `         `
}
