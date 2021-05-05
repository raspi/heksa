package bitWithAscii

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/ascii"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/bit"
)

// Check implementation
var _ base.ByteFormatter = BitWithAsciiPrinter{}

type BitWithAsciiPrinter struct {
	p            base.ByteFormatter
	hilightBreak string
	specialBreak string
}

func New(hilightBreak string, specialBreak string) BitWithAsciiPrinter {
	return BitWithAsciiPrinter{
		p:            bit.New(),
		hilightBreak: hilightBreak,
		specialBreak: specialBreak,
	}
}

func (p BitWithAsciiPrinter) Print(b byte) (o string) {
	return p.p.Print(b) + ` ` + ascii.PrintSpecial(p.specialBreak, p.hilightBreak, b)
}

func (p BitWithAsciiPrinter) GetPrintSize() int {
	return 12
}

func (p BitWithAsciiPrinter) UseSplitter() bool {
	return true
}
