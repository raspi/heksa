package block

import (
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
)

type printer struct {
}

func New() base.ByteFormatter {

	return printer{}
}

func (p printer) Print(b byte) (o string) {
	base.HideVisualSplitter = true

	if base.ChangePalette {
		base.ChangePalette = false
		o += base.Palette[b]
	}
	o += "\u2588"

	return o
}

func (p printer) GetPrintSize() int {
	return 1
}
