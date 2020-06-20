package octal

import (
	"fmt"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
)

type printer struct {
	info   base.BaseInfo
	format string
	size   int
}

func New(info base.BaseInfo) base.OffsetFormatter {
	p := printer{
		info: info,
		size: len(fmt.Sprintf(`%o`, info.FileSize)),
	}

	p.format = fmt.Sprintf(`%%0%do`, p.size)

	return p
}

func (p printer) GetFormatWidth() int {
	return p.size
}
func (p printer) Print(offset uint64) string {
	return fmt.Sprintf(p.format, offset)
}
