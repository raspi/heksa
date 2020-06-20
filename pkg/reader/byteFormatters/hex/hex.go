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
	if base.ChangePalette {
		base.ChangePalette = false
		o += base.Palette[b]
	}

	o += HexByteToString[b]
	return o
}

func (p printer) GetPrintSize() int {
	return 2
}
