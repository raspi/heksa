package bit

import (
	"github.com/raspi/heksa/pkg/color"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
)

// Check implementation
var _ base.ByteFormatter = BitPrinter{}

type BitPrinter struct {
}

func New() BitPrinter {
	return BitPrinter{}
}

func (p BitPrinter) Print(b byte) (o string) {
	for idx, ru := range bitByteToString[b] {
		if idx == 0 {
			o += color.SetUnderlineOn
		}

		o += string(ru)

		if idx == 3 {
			o += color.SetUnderlineOff
		}
	}

	return o
}

func (p BitPrinter) GetPrintSize() int {
	return 8
}

func (p BitPrinter) UseSplitter() bool {
	return true
}
