package display

import (
	clr "github.com/logrusorgru/aurora"
	"strings"
)

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
