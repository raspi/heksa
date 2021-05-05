package octal

import "github.com/raspi/heksa/pkg/reader/byteFormatters/base"

// Check implementation
var _ base.ByteFormatter = OctalPrinter{}

type OctalPrinter struct {
}

func New() OctalPrinter {
	return OctalPrinter{}
}

func (p OctalPrinter) Print(b byte) (o string) {
	return octByteToString[b]
}

func (p OctalPrinter) GetPrintSize() int {
	return 3
}

func (p OctalPrinter) UseSplitter() bool {
	return true
}
