package block

import "github.com/raspi/heksa/pkg/reader/byteFormatters/base"

// Check implementation
var _ base.ByteFormatter = printer{}

type printer struct {
	useSplitter bool
}

func New() base.ByteFormatter {
	return printer{
		useSplitter: false,
	}
}

func (p printer) Print(b byte) string {
	return "\u2588"
}

func (p printer) GetPrintSize() int {
	return 1
}

func (p printer) UseSplitter() bool {
	return p.useSplitter
}
