package ascii

import "github.com/raspi/heksa/pkg/reader/byteFormatters/base"

// Check implementation
var _ base.ByteFormatter = AsciiPrinter{}

type AsciiPrinter struct {
}

func New() AsciiPrinter {
	return AsciiPrinter{}
}

func (p AsciiPrinter) Print(b byte) (o string) {
	return string(AsciiByteToChar[b])
}

func (p AsciiPrinter) GetPrintSize() int {
	return 1
}

func (p AsciiPrinter) UseSplitter() bool {
	return true
}
