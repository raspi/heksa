package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/color"
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
	Splitter             string               // Splitter character for columns
	palette              [256]color.AnsiColor // color palette for each byte
	showHeader           bool                 //  Show formatter header?
	SplitterColor        color.AnsiColor
	OffsetColor          color.AnsiColor
	splitterBreak        string
	offsetBreak          string
}

func New(r iface.ReadSeekerCloser, offsetFormatter []iface.OffsetFormatter, formatters []iface.CharacterFormatter, palette [256]color.AnsiColor, showHeader bool) *Reader {
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
		Splitter:             `┊`,
		charFormatterCount:   len(formatters),
		offsetFormatterCount: len(offsetFormatter),
		palette:              palette,
		showHeader:           showHeader,
		SplitterColor:        color.AnsiColor{Color: color.ColorGrey93_eeeeee},
		OffsetColor:          color.AnsiColor{Color: color.ColorGrey93_eeeeee},
	}

	reader.splitterBreak = fmt.Sprintf(`%s%s`, color.SetForeground, reader.SplitterColor)
	reader.offsetBreak = fmt.Sprintf(`%s%s`, color.SetForeground, reader.OffsetColor)

	return reader
}

func (r *Reader) getoffsetLeft() string {
	r.sb.Reset()
	if r.offsetFormatterCount > 0 {
		r.sb.WriteString(r.offsetBreak)
		// show offset on the left side
		r.sb.WriteString(r.offsetFormatter[0].FormatOffset(r.r))
		r.sb.WriteString(r.splitterBreak)
		r.sb.WriteString(r.Splitter)
	}
	return r.sb.String()
}

func (r *Reader) getoffsetRight() string {
	r.sb.Reset()
	if r.offsetFormatterCount > 1 {
		// show offset on the right side
		r.sb.WriteString(r.splitterBreak)
		r.sb.WriteString(r.Splitter)
		r.sb.WriteString(r.offsetBreak)
		r.sb.WriteString(r.offsetFormatter[1].FormatOffset(r.r))
	}

	return r.sb.String()
}

// Read reads 16 bytes and provides string to display
func (r *Reader) Read() (string, error) {
	offsetLeft := r.getoffsetLeft()
	offsetRight := r.getoffsetRight()
	r.sb.Reset()
	r.sb.Grow(256)

	r.sb.WriteString(offsetLeft)

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
				s := dplay.Format(tmp[i])

				if i == 0 || (i > 0 && tmp[i] != tmp[i-1]) {
					// Only print on first and changed color
					r.sb.WriteString(fmt.Sprintf(`%s%s`, color.SetForeground, r.palette[tmp[i]]))
				}

				if i < 15 {
					r.sb.WriteString(s)
				} else {
					if eofl == 1 {
						r.sb.WriteString(s)
					} else {
						// No extra space for last
						r.sb.WriteString(strings.TrimRight(s, ` `))
					}
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
			r.sb.WriteString(r.splitterBreak)
			r.sb.WriteString(r.Splitter)
		}
	}

	r.sb.WriteString(offsetRight)

	return r.sb.String(), nil
}

func (r *Reader) Header() string {
	if !r.showHeader {
		return ``
	}

	r.sb.Reset()

	if r.offsetFormatterCount > 0 {
		// show offset on the left side
		r.sb.WriteString(r.offsetBreak)
		r.sb.WriteString(r.offsetFormatter[0].OffsetHeader())
		r.sb.WriteString(r.splitterBreak)
		r.sb.WriteString(r.Splitter)
	}

	// iterate through every formatter which outputs it's own header
	for didx, dplay := range r.charFormatters {
		r.sb.WriteString(r.offsetBreak)
		r.sb.WriteString(dplay.Header())
		if didx < (r.charFormatterCount - 1) {
			r.sb.WriteString(r.splitterBreak)
			r.sb.WriteString(r.Splitter)
		}
	}

	if r.offsetFormatterCount > 1 {
		// show offset on the right side
		r.sb.WriteString(r.splitterBreak)
		r.sb.WriteString(r.Splitter)
		r.sb.WriteString(r.offsetBreak)
		r.sb.WriteString(r.offsetFormatter[1].OffsetHeader())
	}

	r.sb.WriteString("\n")

	return r.sb.String()
}
