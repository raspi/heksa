package decWithAscii

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/ascii"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/decimal"
)

// Check implementation
var _ base.ByteFormatter = DecimalWithAsciiPrinter{}

type DecimalWithAsciiPrinter struct {
	p            base.ByteFormatter
	hilightBreak string
	specialBreak string
}

func New(hilightBreak string, specialBreak string) DecimalWithAsciiPrinter {
	return DecimalWithAsciiPrinter{
		p:            decimal.New(),
		hilightBreak: hilightBreak,
		specialBreak: specialBreak,
	}
}

func (p DecimalWithAsciiPrinter) Print(b byte) (o string) {
	o += p.p.Print(b)
	o += ` ` + p.specialBreak + `[` + p.hilightBreak
	o += string(ascii.AsciiByteToChar[b])
	o += p.specialBreak + `]`
	return o
}

func (p DecimalWithAsciiPrinter) GetPrintSize() int {
	return 7
}

func (p DecimalWithAsciiPrinter) UseSplitter() bool {
	return true
}
