package hex

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
)

// Check implementation
var _ base.ByteFormatter = HexPrinter{}

type HexPrinter struct {
}

func New() HexPrinter {
	return HexPrinter{}
}

func (p HexPrinter) Print(b byte) (o string) {
	return HexByteToString[b]
}

func (p HexPrinter) GetPrintSize() int {
	return 2
}

func (p HexPrinter) UseSplitter() bool {
	return true
}
