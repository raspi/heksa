package display

import (
	clr "github.com/logrusorgru/aurora"
	"strings"
)

/*
Bit displays bytes as bits 00000000-11111111
*/
type Bit struct {
	sb strings.Builder
}

func NewBit() *Bit {
	return &Bit{
		sb: strings.Builder{},
	}
}

func (d *Bit) Format(b byte, color clr.Color) string {
	d.sb.Reset()
	d.sb.WriteString(clr.Sprintf(`%08b `, clr.Colorize(b, color)))
	return d.sb.String()
}

func (d *Bit) EofStr() string {
	return `         `
}

func (d *Bit) Header() string {
	return header(8)
}
