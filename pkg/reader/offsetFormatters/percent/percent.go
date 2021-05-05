package percent

import (
	"fmt"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
)

// Check implementation
var _ base.OffsetFormatter = PercentPrinter{}

type PercentPrinter struct {
	info        base.BaseInfo
	format      string
	size        int
	unknownSize bool // Unknown file size? (reading from STDIN)
}

func New(info base.BaseInfo) PercentPrinter {
	p := PercentPrinter{
		info:        info,
		size:        9,
		format:      `%07.3f%%`,
		unknownSize: false,
	}

	if p.info.FileSize == -1 {
		// Unknown file size
		p.unknownSize = true
		p.size = 3
	}

	return p
}

func (p PercentPrinter) GetFormatWidth() int {
	return p.size
}

func (p PercentPrinter) Print(offset uint64) string {
	if p.unknownSize {
		// Can't know percentage for STDIN
		return `??%`
	}

	return fmt.Sprintf(p.format, (float64(offset)*100.0)/float64(p.info.FileSize))
}
