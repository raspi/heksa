package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/color"
	"github.com/raspi/heksa/pkg/iface"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	offFormatters "github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
	"io"
	"strings"
)

type Colors struct {
	splitterBreak string
	offsetBreak   string
}

type Reader struct {
	r                    iface.ReadSeekerCloser
	charFormatters       []base.ByteFormatter            // list of byte displayer(s) for data
	charFormatterCount   int                             // shorthand for len(charFormatters), for speeding up
	offsetFormatters     []offFormatters.OffsetFormatter // offset formatters (max 2) first one is displayed on the left side and second one on the right side
	offsetFormatterCount int                             // shorthand for len(offsetFormatters), for speeding up
	isStdin              bool
	ReadBytes            uint64          // How many bytes Reader has been reading so far (for limit)
	sb                   strings.Builder // Faster than concatenating strings
	Splitter             string          // Splitter character for columns
	Colors               Colors          // Colors
	growHint             int             // Grow hint for sb strings.Builder variable for speed
	width                int             // Width
	visualSplitterSize   int             // Size of visual splitter (2 = XX XX XX, 3 = XXX XXX XXX, etc)
	visualSplitter       string          // Visual splitter that gets inserted every visualSplitterSize bytes
}

func New(r iface.ReadSeekerCloser, offsetFormatter []offFormatters.OffsetFormatter, formatters []base.ByteFormatter, formatterWidth uint16, isStdin bool) *Reader {
	if formatters == nil {
		panic(`nil formatter`)
	}

	if formatterWidth == 0 {
		panic(`zero formatterWidth`)
	}

	reader := &Reader{
		r:                    r,
		visualSplitterSize:   8,   // Insert extra space after every N bytes
		visualSplitter:       ` `, // insert visualSplitter every visualSplitterSize bytes
		width:                int(formatterWidth),
		isStdin:              isStdin,
		charFormatters:       formatters,
		offsetFormatters:     offsetFormatter,
		ReadBytes:            0, // How many bytes we've read
		sb:                   strings.Builder{},
		Splitter:             `┊`, // Splitter character between different columns
		charFormatterCount:   len(formatters),
		offsetFormatterCount: len(offsetFormatter),
		Colors: Colors{
			splitterBreak: fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey93_eeeeee}),
			offsetBreak:   fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey82_d0d0d0}),
		},
	}

	for _, f := range reader.charFormatters {
		reader.growHint += int(formatterWidth)
		reader.growHint += int(formatterWidth) * f.GetPrintSize()
	}

	for _, f := range reader.offsetFormatters {
		reader.growHint += f.GetFormatWidth()
	}

	return reader
}

// getoffsetLeft outputs the selected formatter on the left side
func (r *Reader) getoffsetLeft(offset uint64) string {
	r.sb.Reset()
	if r.offsetFormatterCount > 0 {
		r.sb.WriteString(r.Colors.offsetBreak)
		// show offset on the left side
		r.sb.WriteString(r.offsetFormatters[0].Print(offset))
		r.sb.WriteString(r.Colors.splitterBreak)
		r.sb.WriteString(r.Splitter)
	}

	return r.sb.String()
}

// getoffsetRight outputs the selected formatter on the right side after all the user selected byte formatters
func (r *Reader) getoffsetRight(offset uint64) string {
	r.sb.Reset()
	if r.offsetFormatterCount > 1 {
		// show offset on the right side
		r.sb.WriteString(r.Colors.splitterBreak)
		r.sb.WriteString(r.Splitter)
		r.sb.WriteString(r.Colors.offsetBreak)
		r.sb.WriteString(r.offsetFormatters[1].Print(offset))
	}

	return r.sb.String()
}

// Read reads N (r.width) bytes and provides string to display
func (r *Reader) Read() (string, error) {
	var offset uint64

	if r.isStdin {
		// reading from STDIN, can't use seek
		offset = r.ReadBytes
	} else {
		// Reading from file
		offsettmp, err := r.r.Seek(0, io.SeekCurrent)
		if err != nil {
			return ``, fmt.Errorf(`couldn't seek: %w`, err)
		}

		offset = uint64(offsettmp)
	}

	offsetLeft := r.getoffsetLeft(offset)
	offsetRight := r.getoffsetRight(offset)
	r.sb.Reset()
	r.sb.Grow(r.growHint)

	r.sb.WriteString(offsetLeft)

	tmp := make([]byte, r.width)
	bytesReadCount, err := r.r.Read(tmp)
	if err != nil {
		return ``, err
	}

	r.ReadBytes += uint64(bytesReadCount)

	// iterate through every formatter which outputs it's own format
	for didx, byteFormatterType := range r.charFormatters {
		// First character to print, so always true
		base.ChangePalette = true

		for i := 0; i < r.width; i++ {
			if i != 0 && i%r.visualSplitterSize == 0 {
				// Add pad for better visualization
				r.sb.WriteString(r.visualSplitter)
			}

			if bytesReadCount > i {
				if i == 0 || (i > 0 && tmp[i] != tmp[i-1] && base.Palette[tmp[i]] != base.Palette[tmp[i-1]]) {
					base.ChangePalette = true
				}

				r.sb.WriteString(byteFormatterType.Print(tmp[i]))

				if i < (r.width-1) && byteFormatterType.GetPrintSize() > 1 {
					r.sb.WriteString(` `)
				}
			} else {
				// There is no data so we add padding
				if i == 0 || (i > 0 && tmp[i] != tmp[i-1] && base.Palette[tmp[i]] != base.Palette[tmp[i-1]]) {
					// Only print on first and changed color
					r.sb.WriteString(base.Palette[0])
				}

				r.sb.WriteString(strings.Repeat(`‡`, byteFormatterType.GetPrintSize()))

				if i < (r.width-1) && byteFormatterType.GetPrintSize() > 1 {
					r.sb.WriteString(` `)
				}
			}
		}

		if didx < (r.charFormatterCount - 1) {
			r.sb.WriteString(r.Colors.splitterBreak)
			r.sb.WriteString(r.Splitter)
		}
	}

	r.sb.WriteString(offsetRight)

	return r.sb.String(), nil
}
