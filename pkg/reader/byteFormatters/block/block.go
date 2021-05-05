package block

import "github.com/raspi/heksa/pkg/reader/byteFormatters/base"

// Check implementation
var _ base.ByteFormatter = BlockPrinter{}

type BlockPrinter struct {
	useSplitter bool
}

func New() BlockPrinter {
	return BlockPrinter{
		useSplitter: false,
	}
}

func (p BlockPrinter) Print(b byte) string {
	return "\u2588"
}

func (p BlockPrinter) GetPrintSize() int {
	return 1
}

func (p BlockPrinter) UseSplitter() bool {
	return p.useSplitter
}
