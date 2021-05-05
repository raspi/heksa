package hex

import (
	"fmt"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
)

// Check implementation
var _ base.OffsetFormatter = HexPrinter{}

// minimal size for padding zeroes
const minimalSize = 8 // 8 = 0xFFFFFFFF = 4294967295 bytes = ~4 GiB

type HexPrinter struct {
	info   base.BaseInfo
	format string
	size   int
}

func New(info base.BaseInfo) HexPrinter {
	size := len(fmt.Sprintf(`%x`, info.FileSize))
	if size < minimalSize {
		size = minimalSize
	}

	p := HexPrinter{
		info: info,
		size: size,
	}

	p.format = fmt.Sprintf(`%%0%dx`, p.size)
	return p
}

func (p HexPrinter) GetFormatWidth() int {
	return p.size
}

func (p HexPrinter) Print(offset uint64) string {
	return fmt.Sprintf(p.format, offset)
}
