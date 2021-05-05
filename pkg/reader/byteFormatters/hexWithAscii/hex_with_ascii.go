package hexWithAscii

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/ascii"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/hex"
)

// Check implementation
var _ base.ByteFormatter = HexWithAsciiPrinter{}

type HexWithAsciiPrinter struct {
	p            base.ByteFormatter
	hilightBreak string
	specialBreak string
}

func New(hilightBreak string, specialBreak string) HexWithAsciiPrinter {
	return HexWithAsciiPrinter{
		p:            hex.New(),
		hilightBreak: hilightBreak,
		specialBreak: specialBreak,
	}
}

func (p HexWithAsciiPrinter) Print(b byte) (o string) {
	o += p.p.Print(b)
	o += ` ` + p.specialBreak + `[` + p.hilightBreak
	o += string(ascii.AsciiByteToChar[b])
	o += p.specialBreak + `]`
	return o
}

func (p HexWithAsciiPrinter) GetPrintSize() int {
	return 6
}

func (p HexWithAsciiPrinter) UseSplitter() bool {
	return true
}
