package bitWithDecimal

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/bit"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/decimal"
)

// Check implementation
var _ base.ByteFormatter = BitWithDecimalPrinter{}

type BitWithDecimalPrinter struct {
	p            base.ByteFormatter
	hilightBreak string
	specialBreak string
}

func New(hilightBreak string, specialBreak string) BitWithDecimalPrinter {
	return BitWithDecimalPrinter{
		p:            bit.New(),
		hilightBreak: hilightBreak,
		specialBreak: specialBreak,
	}
}

func (p BitWithDecimalPrinter) Print(b byte) (o string) {
	o += p.p.Print(b)
	o += ` ` + p.specialBreak + `[` + p.hilightBreak
	o += decimal.DecByteToString[b]
	o += p.specialBreak + `]`
	return o
}

func (p BitWithDecimalPrinter) GetPrintSize() int {
	return 14
}

func (p BitWithDecimalPrinter) UseSplitter() bool {
	return true
}
