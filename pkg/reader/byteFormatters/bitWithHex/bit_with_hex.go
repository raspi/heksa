package bitWithHex

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/bit"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/hex"
)

// Check implementation
var _ base.ByteFormatter = BitWithHexPrinter{}

type BitWithHexPrinter struct {
	p            base.ByteFormatter
	hilightBreak string
	specialBreak string
}

func New(hilightBreak string, specialBreak string) BitWithHexPrinter {
	return BitWithHexPrinter{
		p:            bit.New(),
		hilightBreak: hilightBreak,
		specialBreak: specialBreak,
	}
}

func (p BitWithHexPrinter) Print(b byte) (o string) {
	o += p.p.Print(b)
	o += ` ` + p.specialBreak + `[` + p.hilightBreak
	o += hex.HexByteToString[b]
	o += p.specialBreak + `]`
	return o
}

func (p BitWithHexPrinter) GetPrintSize() int {
	return 13
}

func (p BitWithHexPrinter) UseSplitter() bool {
	return true
}
