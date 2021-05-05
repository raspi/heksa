package octal

import (
	"fmt"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
)

// Check implementation
var _ base.OffsetFormatter = OctalPrinter{}

// minimal size for padding zeroes
const minimalSize = 6

type OctalPrinter struct {
	info   base.BaseInfo
	format string
	size   int
}

func New(info base.BaseInfo) OctalPrinter {
	size := len(fmt.Sprintf(`%o`, info.FileSize))
	if size < minimalSize {
		size = minimalSize
	}

	p := OctalPrinter{
		info: info,
		size: size,
	}

	p.format = fmt.Sprintf(`%%0%do`, p.size)

	return p
}

func (p OctalPrinter) GetFormatWidth() int {
	return p.size
}

func (p OctalPrinter) Print(offset uint64) string {
	return fmt.Sprintf(p.format, offset)
}
