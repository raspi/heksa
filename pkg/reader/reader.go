package reader

import (
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/pkg/iface"
	"strings"
)

type Reader struct {
	r                    iface.ReadSeekerCloser
	charFormatters       []iface.CharacterFormatter // displayer(s) for data
	charFormatterCount   int
	offsetFormatter      []iface.OffsetFormatter // offset formatters (max 2) first one is displayed on the left side and second one on the right side
	offsetFormatterCount int
	ReadBytes            uint64 // How many bytes Reader has been reading so far (for limit)
	sb                   strings.Builder
	Splitter             string         // Splitter character for columns
	palette              [256]clr.Color // color palette for each byte
}

func New(r iface.ReadSeekerCloser, offsetFormatter []iface.OffsetFormatter, formatters []iface.CharacterFormatter, palette [256]clr.Color) *Reader {
	if offsetFormatter == nil {
		panic(`nil offset formatter`)
	}

	if formatters == nil {
		panic(`nil formatter`)
	}

	reader := &Reader{
		r:                    r,
		charFormatters:       formatters,
		offsetFormatter:      offsetFormatter,
		ReadBytes:            0,
		sb:                   strings.Builder{},
		Splitter:             `|`,
		charFormatterCount:   len(formatters),
		offsetFormatterCount: len(offsetFormatter),
		palette:              palette,
	}

	return reader
}

// Read reads 16 bytes and provides string to display
func (r *Reader) Read() (string, error) {
	r.sb.Reset()
	r.sb.Grow(1024)

	if r.offsetFormatterCount > 0 {
		// show offset on the left side
		r.sb.WriteString(r.offsetFormatter[0].FormatOffset(r.r))
		r.sb.WriteString(r.Splitter)
	}

	tmp := make([]byte, 16)
	rb, err := r.r.Read(tmp)
	if err != nil {
		return ``, err
	}

	r.ReadBytes += uint64(rb)

	// iterate through every formatter which outputs it's own format
	for didx, dplay := range r.charFormatters {
		eof := []byte(dplay.EofStr())
		eofl := len(eof)

		for i := 0; i < 16; i++ {
			if i == 8 {
				// Add pad for better visualization
				r.sb.WriteString(` `)
			}

			if rb > i {
				s := dplay.Format(tmp[i], r.palette[tmp[i]])

				if i < 15 {
					r.sb.WriteString(s)
				} else {
					// No extra space for last
					r.sb.WriteString(strings.Trim(s, ` `))
				}
			} else {
				// There is no data so we add padding
				if i < 15 {
					r.sb.Write(eof)
				} else {
					// No extra spaces for last
					if eofl > 1 {
						r.sb.Write(eof[0 : eofl-1])
					} else {
						r.sb.Write(eof)
					}
				}
			}
		}

		if didx < (r.charFormatterCount - 1) {
			r.sb.WriteString(r.Splitter)
		}
	}

	if r.offsetFormatterCount > 1 {
		// show offset on the right side
		r.sb.WriteString(r.Splitter)
		r.sb.WriteString(r.offsetFormatter[1].FormatOffset(r.r))
	}

	return r.sb.String(), nil
}
