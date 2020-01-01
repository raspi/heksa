package display

import (
	"fmt"
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

func (d *Bit) Format(b byte) string {
	d.sb.Reset()
	d.sb.WriteString(fmt.Sprintf(`%08b `, b))
	return d.sb.String()
}

func (d *Bit) EofStr() string {
	return `         `
}

func (d *Bit) Header() string {
	return header(8)
}
