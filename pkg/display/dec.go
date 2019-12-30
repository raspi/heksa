package display

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"math/bits"
	"strings"
)

type Dec struct {
	fs        uint64
	bw        uint8 // Bit width calculated from file size
	palette   map[uint8]clr.Color
	offFormat string
	sb        strings.Builder
}

func (d *Dec) SetFileSize(s int64) {
	d.fs = uint64(s)
	d.bw = nearest(uint8(bits.Len64(d.fs)))
	d.offFormat = fmt.Sprintf(`%%0%vd`, d.bw)
}

func NewDec() *Dec {
	return &Dec{
		fs: 8,
		sb: strings.Builder{},
	}
}

func (d *Dec) Format(b byte) string {
	d.sb.Reset()

	color, ok := d.palette[b]
	if !ok {
		color = clr.BrightFg
	}

	d.sb.WriteString(clr.Sprintf(`%03d `, clr.Colorize(b, color)))

	return d.sb.String()

}

// DisplayOffset displays offset as decimal 0 - 9999999....
func (d *Dec) DisplayOffset(r iface.ReadSeekerCloser) string {
	d.sb.Reset()
	off, _ := r.Seek(0, io.SeekCurrent)
	d.sb.WriteString(fmt.Sprintf(d.offFormat, off))
	return d.sb.String()

}

func (d *Dec) SetPalette(p map[uint8]clr.Color) {
	d.palette = p
}

func (d *Dec) EofStr() string {
	return `    `
}
