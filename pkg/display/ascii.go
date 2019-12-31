package display

import (
	clr "github.com/logrusorgru/aurora"
	"strings"
)

type Ascii struct {
	sb strings.Builder
}

func NewAscii() *Ascii {
	return &Ascii{
		sb: strings.Builder{},
	}
}

func (d *Ascii) Format(b byte, color clr.Color) string {
	d.sb.Reset()
	d.sb.WriteString(clr.Sprintf(`%c`, clr.Colorize(d.toChar(b), color)))
	return d.sb.String()
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

func (d *Ascii) Header() string {
	return header(1)
}
