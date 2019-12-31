package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"strings"
)

/*
Oct displays bytes in octal format 000-377
*/
type Oct struct {
	fs        uint64 // File size
	offFormat string // Format for offset column
	sb        strings.Builder
	zeroes    int
}

func NewOct() *Oct {
	return &Oct{
		fs: 0,
		sb: strings.Builder{},
	}
}

func (d *Oct) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.zeroes = len(fmt.Sprintf(`%o`, d.fs))
	d.offFormat = fmt.Sprintf(`%%0%vo`, d.zeroes)
}

func (d *Oct) Format(b byte, color clr.Color) string {
	d.sb.Reset()
	d.sb.WriteString(clr.Sprintf(`%03o `, clr.Colorize(b, color)))
	return d.sb.String()
}

// FormatOffset displays offset as hexadecimal 0x00 - 0xFFFFFFFF....
func (d *Oct) FormatOffset(r iface.ReadSeekerCloser) string {
	d.sb.Reset()
	off, _ := r.Seek(0, io.SeekCurrent)
	d.sb.WriteString(fmt.Sprintf(d.offFormat, off))
	return d.sb.String()
}

func (d *Oct) EofStr() string {
	return `    `
}

func (d *Oct) OffsetHeader() string {
	return strings.Repeat(`_`, d.zeroes)
}

func (d *Oct) Header() string {
	return header(3)
}
