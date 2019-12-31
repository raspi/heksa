package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"strings"
)

type Dec struct {
	fs        uint64
	offFormat string // Format for offset column
	sb        strings.Builder
	zeroes    int
}

func NewDec() *Dec {
	return &Dec{
		fs: 0,
		sb: strings.Builder{},
	}
}

func (d *Dec) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.zeroes = len(string(d.fs))
	d.offFormat = fmt.Sprintf(`%%0%vd`, d.zeroes)
}

func (d *Dec) Format(b byte, color clr.Color) string {
	d.sb.Reset()
	d.sb.WriteString(clr.Sprintf(`%03d `, clr.Colorize(b, color)))
	return d.sb.String()
}

// FormatOffset displays offset as decimal 0 - 9999999....
func (d *Dec) FormatOffset(r iface.ReadSeekerCloser) string {
	d.sb.Reset()
	off, _ := r.Seek(0, io.SeekCurrent)
	d.sb.WriteString(fmt.Sprintf(d.offFormat, off))
	return d.sb.String()
}

func (d *Dec) EofStr() string {
	return `    `
}

func (d *Dec) OffsetHeader() string {
	return strings.Repeat(`_`, d.zeroes)
}

func (d *Dec) Header() string {
	return header(3)
}
