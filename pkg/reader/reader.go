package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/color"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"strings"
)

type Reader struct {
	r                    iface.ReadSeekerCloser
	charFormatters       []ByteFormatter // list of byte displayer(s) for data
	charFormatterCount   int
	offsetFormatter      []OffsetFormatter // offset formatters (max 2) first one is displayed on the left side and second one on the right side
	offsetFormatterCount int
	fileSize             int64  // file size reference
	ReadBytes            uint64 // How many bytes Reader has been reading so far (for limit)
	sb                   strings.Builder
	Splitter             string               // Splitter character for columns
	palette              [256]color.AnsiColor // color palette for each byte
	showHeader           bool                 //  Show formatter header?
	SplitterColor        color.AnsiColor
	OffsetColor          color.AnsiColor
	splitterBreak        string
	offsetBreak          string

	offsetFormatterFormat map[OffsetFormatter]string // Printf format
	offsetFormatterWidth  map[OffsetFormatter]int    // How much padding width needed
}

func New(r iface.ReadSeekerCloser, offsetFormatter []OffsetFormatter, formatters []ByteFormatter, palette [256]color.AnsiColor, showHeader bool, filesize int64) *Reader {
	if formatters == nil {
		panic(`nil formatter`)
	}

	reader := &Reader{
		r:                    r,
		fileSize:             filesize,
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

	reader.offsetFormatterFormat = make(map[OffsetFormatter]string, reader.offsetFormatterCount)
	reader.offsetFormatterWidth = make(map[OffsetFormatter]int, reader.offsetFormatterCount)

	reader.splitterBreak = fmt.Sprintf(`%s%s`, color.SetForeground, reader.SplitterColor)
	reader.offsetBreak = fmt.Sprintf(`%s%s`, color.SetForeground, reader.OffsetColor)

	for _, f := range reader.offsetFormatter {
		_, ok := reader.offsetFormatterWidth[f]

		if ok {
			continue
		}

		switch f {
		case OffsetHex:
			reader.offsetFormatterWidth[f] = len(fmt.Sprintf(`%x`, filesize))
		case OffsetDec:
			reader.offsetFormatterWidth[f] = len(fmt.Sprintf(`%x`, filesize))
		case OffsetOct:
			reader.offsetFormatterWidth[f] = len(fmt.Sprintf(`%x`, filesize))
		case OffsetPercent:
			reader.offsetFormatterWidth[f] = 8
		}

		width, ok := reader.offsetFormatterWidth[f]

		if !ok {
			panic(fmt.Errorf(`couldn't find width??`))
		}

		switch f {
		case OffsetHex:
			reader.offsetFormatterFormat[f] = fmt.Sprintf(`%%0%dx`, width)
		case OffsetDec:
			reader.offsetFormatterFormat[f] = fmt.Sprintf(`%%0%dd`, width)
		case OffsetOct:
			reader.offsetFormatterFormat[f] = fmt.Sprintf(`%%0%do`, width)
		case OffsetPercent:
			reader.offsetFormatterFormat[f] = `%07.3f%%`
		}

	}

	return reader
}

func (r *Reader) formatOffset(formatter OffsetFormatter, offset int64) {
	switch formatter {
	case OffsetPercent:
		percent := (float64(offset) * 100.0) / float64(r.fileSize)
		r.sb.WriteString(fmt.Sprintf(r.offsetFormatterFormat[formatter], percent))
	default:
		r.sb.WriteString(fmt.Sprintf(r.offsetFormatterFormat[formatter], offset))
	}
}

func (r *Reader) getoffsetLeft(offset int64) string {
	r.sb.Reset()
	if r.offsetFormatterCount > 0 {
		r.sb.WriteString(r.offsetBreak)
		// show offset on the left side
		r.formatOffset(r.offsetFormatter[0], offset)
		r.sb.WriteString(r.splitterBreak)
		r.sb.WriteString(r.Splitter)
	}

	return r.sb.String()
}

func (r *Reader) getoffsetRight(offset int64) string {
	r.sb.Reset()
	if r.offsetFormatterCount > 1 {
		// show offset on the right side
		r.sb.WriteString(r.splitterBreak)
		r.sb.WriteString(r.Splitter)
		r.sb.WriteString(r.offsetBreak)
		r.formatOffset(r.offsetFormatter[1], offset)
	}

	return r.sb.String()
}

// Read reads 16 bytes and provides string to display
func (r *Reader) Read() (string, error) {
	offset, err := r.r.Seek(0, io.SeekCurrent)
	if err != nil {
		return ``, err
	}

	offsetLeft := r.getoffsetLeft(offset)
	offsetRight := r.getoffsetRight(offset)
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
	for didx, byteFormatterType := range r.charFormatters {

		for i := 0; i < 16; i++ {
			if i == 8 {
				// Add pad for better visualization
				r.sb.WriteString(` `)
			}

			if rb > i {
				if i == 0 || (i > 0 && tmp[i] != tmp[i-1]) {
					// Only print on first and changed color
					r.sb.WriteString(fmt.Sprintf(`%s%s`, color.SetForeground, r.palette[tmp[i]]))
				}

				var s string

				switch byteFormatterType {
				case ViewHex:
					s = fmt.Sprintf(`%02x`, tmp[i])
				case ViewDec:
					s = fmt.Sprintf(`%03d`, tmp[i])
				case ViewOct:
					s = fmt.Sprintf(`%03o`, tmp[i])
				case ViewBit:
					s = fmt.Sprintf(`%08b`, tmp[i])
				case ViewASCII:
					s = fmt.Sprintf(`%c`, asciiByteToChar[tmp[i]])
				}

				if i < 15 {
					r.sb.WriteString(s)
					switch byteFormatterType {
					case ViewASCII:
						continue
					default:
						r.sb.WriteString(` `)
					}
				} else {
					switch byteFormatterType {
					case ViewHex, ViewBit, ViewOct, ViewDec:
					}
					r.sb.WriteString(s)
				}
			} else {
				// There is no data so we add padding
				if i < 15 {
					switch byteFormatterType {
					case ViewHex:
						r.sb.WriteString(`-- `)
					case ViewOct, ViewDec:
						r.sb.WriteString(`--- `)
					case ViewASCII:
						r.sb.WriteString(` `)
					}
				} else {
					switch byteFormatterType {
					case ViewHex:
						r.sb.WriteString(`--`)
					case ViewDec, ViewOct:
						r.sb.WriteString(`---`)
					case ViewASCII:
						r.sb.WriteString(` `)
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

func (r *Reader) offsetHeader(otype OffsetFormatter) string {
	width := r.offsetFormatterWidth[otype]
	return strings.Repeat(`_`, width)
}

func (r *Reader) header(l uint8) (out string) {
	format := fmt.Sprintf(`%%0%dx`, l)
	for i := uint8(0); i < 16; i++ {
		if i == 8 {
			out += ` `
		}
		out += fmt.Sprintf(format, i)
		if l > 1 && i < 15 {
			out += ` `
		}
	}

	return out
}

func (r *Reader) Header() (out string) {
	if !r.showHeader {
		return ``
	}

	if r.offsetFormatterCount > 0 {
		// show offset on the left side
		out += r.offsetBreak
		out += r.offsetHeader(r.offsetFormatter[0])
		out += r.splitterBreak
		out += r.Splitter
	}

	// iterate through every formatter which outputs it's own header
	for didx, dplay := range r.charFormatters {
		out += r.offsetBreak

		switch dplay {
		case ViewHex:
			out += r.header(2)
		case ViewASCII:
			out += r.header(1)
		case ViewDec, ViewOct:
			out += r.header(3)
		case ViewBit:
			out += r.header(8)
		}

		if didx < (r.charFormatterCount - 1) {
			out += r.splitterBreak
			out += r.Splitter
		}
	}

	if r.offsetFormatterCount > 1 {
		// show offset on the right side
		out += r.splitterBreak
		out += r.Splitter
		out += r.offsetBreak
		out += r.offsetHeader(r.offsetFormatter[1])
	}

	out += "\n"

	return out
}
