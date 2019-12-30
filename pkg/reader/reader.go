package reader

import (
	"github.com/raspi/heksa/pkg/iface"
	"strings"
)

type Reader struct {
	r                     iface.ReadSeekerCloser
	displays              []iface.CharacterFormatter // displayer(s) for data
	displayFormatterCount int
	offsetFormatter       iface.OffsetFormatter // offset displayer
	ReadBytes             uint64
	sb                    strings.Builder
	Splitter              string
}

func New(r iface.ReadSeekerCloser, offsetFormatter iface.OffsetFormatter, formatters []iface.CharacterFormatter) *Reader {
	if offsetFormatter == nil {
		panic(`nil offset displayer`)
	}

	if formatters == nil {
		panic(`nil displayer(s)`)
	}

	reader := &Reader{
		r:                     r,
		displays:              formatters,
		offsetFormatter:       offsetFormatter,
		ReadBytes:             0,
		sb:                    strings.Builder{},
		Splitter:              `|`,
		displayFormatterCount: len(formatters),
	}

	return reader
}

// Read reads 16 bytes and provides string to display
func (r *Reader) Read() (string, error) {
	r.sb.Reset()
	r.sb.Grow(1024)

	r.sb.WriteString(r.offsetFormatter.FormatOffset(r.r))
	r.sb.WriteString(r.Splitter)

	tmp := make([]byte, 16)
	rb, err := r.r.Read(tmp)
	if err != nil {
		return ``, err
	}

	r.ReadBytes += uint64(rb)

	for didx, dplay := range r.displays {

		eof := []byte(dplay.EofStr())
		eofl := len(eof)

		for i := 0; i < 16; i++ {
			if i == 8 {
				r.sb.WriteString(` `)
			}

			if rb > i {
				s := dplay.Format(tmp[i])
				if i < 15 {
					r.sb.WriteString(s)
				} else {
					r.sb.WriteString(strings.Trim(s, ` `))
				}
			} else {
				if i < 15 {
					r.sb.Write(eof)
				} else {
					if eofl > 1 {
						r.sb.Write(eof[0 : eofl-1])
					} else {
						r.sb.Write(eof)
					}
				}
			}
		}

		if didx < (r.displayFormatterCount - 1) {
			r.sb.WriteString(r.Splitter)
		}
	}

	return r.sb.String(), nil
}
