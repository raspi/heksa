package human

import (
	"fmt"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
)

type printer struct {
	format     string
	formatSize int
	unit       uint64
}

func New(unit int) base.OffsetFormatter {

	p := printer{
		unit: uint64(unit),
	}

	switch unit {
	case 1000: // SI
		p.format = `% 8.3f %cB`
		p.formatSize = 11
	case 1024: // IEC
		p.format = `% 8.3f %ciB`
		p.formatSize = 12
	default:
		panic(fmt.Sprintf(`invalid unit formatSize %v`, unit))
	}

	return p
}

func (p printer) GetFormatWidth() int {
	return p.formatSize
}

func (p printer) Print(b uint64) string {
	if b < p.unit {
		switch p.unit {
		case 1000: // SI
			return fmt.Sprintf(`% 8d B `, b)
		case 1024: // IEC
			return fmt.Sprintf(`% 8d B  `, b)
		}
	}

	div, exp := p.unit, uint8(0)

	for n := b / p.unit; n >= p.unit; n /= p.unit {
		div *= p.unit
		exp++
	}

	return fmt.Sprintf(p.format, float64(b)/float64(div), "KMGTPE"[exp])
}
