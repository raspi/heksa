package decimal

import (
	"fmt"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
)

// Check implementation
var _ base.OffsetFormatter = DecimalPrinter{}

// minimal size for padding zeroes
const minimalSize = 6

type DecimalPrinter struct {
	info   base.BaseInfo
	format string
	size   int
}

func New(info base.BaseInfo) DecimalPrinter {
	size := len(fmt.Sprintf(`%d`, info.FileSize))
	if size < minimalSize {
		size = minimalSize
	}

	p := DecimalPrinter{
		info: info,
		size: size,
	}

	p.format = fmt.Sprintf(`%%0%dd`, p.size)
	return p
}

func (p DecimalPrinter) GetFormatWidth() int {
	return p.size
}

func (p DecimalPrinter) Print(offset uint64) string {
	return fmt.Sprintf(p.format, offset)
}
