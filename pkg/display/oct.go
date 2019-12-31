package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"math/bits"
	"strings"
)

type Oct struct {
	fs        uint64 // File size
	bw        uint8  // Bit width calculated from file size
	offFormat string
	sb        strings.Builder
}

func (d *Oct) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.bw = nearest(uint8(bits.Len64(d.fs)))
	d.offFormat = fmt.Sprintf(`%%0%vo`, d.bw)
}

func NewOct() *Oct {
	return &Oct{
		sb: strings.Builder{},
	}
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
