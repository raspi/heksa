package ascii

import "github.com/raspi/heksa/pkg/reader/byteFormatters/base"

// Check implementation
var _ base.ByteFormatter = printer{}

type printer struct {
}

func New() base.ByteFormatter {
	return printer{}
}

func (p printer) Print(b byte) (o string) {
	return string(AsciiByteToChar[b])
}

func (p printer) GetPrintSize() int {
	return 1
}

func (p printer) UseSplitter() bool {
	return true
}
