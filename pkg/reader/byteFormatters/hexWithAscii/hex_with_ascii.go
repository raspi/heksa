package hexWithAscii

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/ascii"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/hex"
)

type printer struct {
	p base.ByteFormatter
}

func New() base.ByteFormatter {
	return printer{
		p: hex.New(),
	}
}

func (p printer) Print(b byte) (o string) {
	base.ChangePalette = true

	o += p.p.Print(b)
	o += ` ` + base.SpecialBreak + `[` + base.HilightBreak
	o += string(ascii.AsciiByteToChar[b])
	o += base.SpecialBreak + `]`
	return o
}

func (p printer) GetPrintSize() int {
	return 6
}
