package decimal

import "github.com/raspi/heksa/pkg/reader/byteFormatters/base"

// Check implementation
var _ base.ByteFormatter = DecimalPrinter{}

type DecimalPrinter struct {
}

func New() DecimalPrinter {
	return DecimalPrinter{}
}

func (p DecimalPrinter) Print(b byte) (o string) {
	return DecByteToString[b]
}

func (p DecimalPrinter) GetPrintSize() int {
	return 3
}

func (p DecimalPrinter) UseSplitter() bool {
	return true
}
