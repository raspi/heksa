package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"math/bits"
	"strings"
)

type Hex struct {
	fs        uint64 // File size
	bw        uint8  // Bit width calculated from file size
	offFormat string
	sb        strings.Builder
}

func (d *Hex) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.bw = nearest(uint8(bits.Len64(d.fs)))

	if d.bw == 0 {
		d.bw = 8
	}

	zeroes := (d.bw / 8) + 1

	if zeroes&1 != 0 {
		zeroes++
	}

	d.offFormat = fmt.Sprintf(`%%0%vx`, zeroes)
}

func NewHex() *Hex {
	return &Hex{
		sb: strings.Builder{},
	}
}

func (d *Hex) Format(b byte, color clr.Color) string {
	d.sb.Reset()
	d.sb.WriteString(clr.Sprintf(`%02x `, clr.Colorize(b, color)))
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
