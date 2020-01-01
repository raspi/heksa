package display

import (
	"fmt"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"strings"
)

/*
Hex displays bytes in hexadecimal format 00-ff.
*/
type Hex struct {
	fs        uint64 // File size
	offFormat string // Format for offset column
	sb        strings.Builder
	zeroes    int
}

func NewHex() *Hex {
	return &Hex{
		fs: 0,
		sb: strings.Builder{},
	}
}

func (d *Hex) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.zeroes = len(fmt.Sprintf(`%x`, d.fs))
	if d.zeroes&1 != 0 {
		d.zeroes++
	}

	d.offFormat = fmt.Sprintf(`%%0%vx`, d.zeroes)
}

func (d *Hex) Format(b byte) string {
	d.sb.Reset()
	d.sb.WriteString(fmt.Sprintf(`%02x `, b))
	return d.sb.String()
}

// FormatOffset displays offset as hexadecimal 0x00 - 0xFFFFFFFF....
func (d *Hex) FormatOffset(r iface.ReadSeekerCloser) string {
	d.sb.Reset()
	off, _ := r.Seek(0, io.SeekCurrent)
	d.sb.WriteString(fmt.Sprintf(d.offFormat, off))
	return d.sb.String()
}

func (d *Hex) EofStr() string {
	return `   `
}

func (d *Hex) OffsetHeader() string {
	return strings.Repeat(`_`, d.zeroes)
}

func (d *Hex) Header() string {
	return header(2)
}
