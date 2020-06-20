package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/color"
	"github.com/raspi/heksa/pkg/iface"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"io"
	"strings"
)

type Colors struct {
	splitterBreak string
	offsetBreak   string
}

type Reader struct {
	r                     iface.ReadSeekerCloser
	charFormatters        []base.ByteFormatter       // list of byte displayer(s) for data
	charFormatterCount    int                        // shorthand for len(charFormatters), for speeding up
	offsetFormatter       []OffsetFormatter          // offset formatters (max 2) first one is displayed on the left side and second one on the right side
	offsetFormatterCount  int                        // shorthand for len(offsetFormatter), for speeding up
	fileSize              int64                      // file size reference, -1 means STDIN. Hint for offset formatter(s) for how many padding characters to use.
	ReadBytes             uint64                     // How many bytes Reader has been reading so far (for limit)
	sb                    strings.Builder            // Faster than concatenating strings
	Splitter              string                     // Splitter character for columns
	offsetFormatterFormat map[OffsetFormatter]string // Printf format for offset format
	offsetFormatterWidth  map[OffsetFormatter]int    // How much padding width needed, calculated from fileSize variable
	Colors                Colors                     // Colors
	growHint              int                        // Grow hint for sb strings.Builder variable for speed
	width                 int                        // Width
	visualSplitterSize    int                        // Size of visual splitter (2 = XX XX XX, 3 = XXX XXX XXX, etc)
	visualSplitter        string                     // Visual splitter that gets inserted every visualSplitterSize bytes
}

func New(r iface.ReadSeekerCloser, offsetFormatter []OffsetFormatter, formatters []base.ByteFormatter, formatterWidth uint16, filesize int64) *Reader {
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
		fileSize:             filesize,
		charFormatters:       formatters,
		offsetFormatter:      offsetFormatter,
		ReadBytes:            0, // How many bytes we've read
		sb:                   strings.Builder{},
		Splitter:             `┊`, // Splitter character between different columns
		charFormatterCount:   len(formatters),
		offsetFormatterCount: len(offsetFormatter),
		Colors: Colors{
			splitterBreak: fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey93_eeeeee}),
			offsetBreak:   fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey93_eeeeee}),
		},
	}

	reader.offsetFormatterFormat = make(map[OffsetFormatter]string, reader.offsetFormatterCount)
	reader.offsetFormatterWidth = make(map[OffsetFormatter]int, reader.offsetFormatterCount)

	for _, f := range reader.charFormatters {
		reader.growHint += int(formatterWidth)
		reader.growHint += int(formatterWidth) * f.GetPrintSize()
	}

	for _, f := range reader.offsetFormatter {
		switch f {
		case OffsetHex:
			width := len(fmt.Sprintf(`%x`, filesize))
			reader.growHint += width + 1
			reader.offsetFormatterWidth[f] = width
			reader.offsetFormatterFormat[f] = fmt.Sprintf(`%%0%dx`, width)
		case OffsetDec:
			width := len(fmt.Sprintf(`%d`, filesize))
			reader.growHint += width + 1
			reader.offsetFormatterWidth[f] = width
			reader.offsetFormatterFormat[f] = fmt.Sprintf(`%%0%dd`, width)
		case OffsetOct:
			width := len(fmt.Sprintf(`%o`, filesize))
			reader.growHint += width + 1
			reader.offsetFormatterWidth[f] = width
			reader.offsetFormatterFormat[f] = fmt.Sprintf(`%%0%do`, width)
		case OffsetPercent:
			width := 9
			reader.growHint += width
			reader.offsetFormatterWidth[f] = width
			reader.offsetFormatterFormat[f] = `%07.3f%%`
		}
	}

	return reader
}

// formatOffset generates offset output such as "0123"
func (r *Reader) formatOffset(formatter OffsetFormatter, offset uint64) {
	switch formatter {
	case OffsetPercent:
		percent := (float64(offset) * 100.0) / float64(r.fileSize)
		r.sb.WriteString(fmt.Sprintf(r.offsetFormatterFormat[formatter], percent))
	default:
		r.sb.WriteString(fmt.Sprintf(r.offsetFormatterFormat[formatter], offset))
	}
}

// getoffsetLeft outputs the selected formatter on the left side
func (r *Reader) getoffsetLeft(offset uint64) string {
	r.sb.Reset()
	if r.offsetFormatterCount > 0 {
		r.sb.WriteString(r.Colors.offsetBreak)
		// show offset on the left side
		r.formatOffset(r.offsetFormatter[0], offset)
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
		r.formatOffset(r.offsetFormatter[1], offset)
	}

	return r.sb.String()
}

// Read reads N (r.width) bytes and provides string to display
func (r *Reader) Read() (string, error) {
	var offset uint64

	if r.fileSize == -1 {
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
