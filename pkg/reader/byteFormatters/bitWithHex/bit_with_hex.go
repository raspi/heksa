package bitWithHex

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/bit"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/hex"
)

type printer struct {
	p            base.ByteFormatter
	hilightBreak string
	specialBreak string
}

func New(hilightBreak string, specialBreak string) base.ByteFormatter {
	return printer{
		p:            bit.New(),
		hilightBreak: hilightBreak,
		specialBreak: specialBreak,
	}
}

func (p printer) Print(b byte) (o string) {
	o += p.p.Print(b)
	o += ` ` + p.specialBreak + `[` + p.hilightBreak
	o += hex.HexByteToString[b]
	o += p.specialBreak + `]`
	return o
}

func (p printer) GetPrintSize() int {
	return 13
}

func (p printer) UseSplitter() bool {
	return true
}
