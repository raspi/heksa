package hex

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
)

type printer struct {
}

func New() base.ByteFormatter {
	return printer{}
}

func (p printer) Print(b byte) (o string) {
	return HexByteToString[b]
}

func (p printer) GetPrintSize() int {
	return 2
}

func (p printer) UseSplitter() bool {
	return true
}
