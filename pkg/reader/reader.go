package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/color"
	"github.com/raspi/heksa/pkg/iface"
	"io"
	"strings"
)

type Colors struct {
	palette       [256]string // precalculated color palette for each byte
	SplitterColor color.AnsiColor
	OffsetColor   color.AnsiColor
	splitterBreak string
	offsetBreak   string
	specialBreak  string
	hilightBreak  string
}

type Reader struct {
	r                     iface.ReadSeekerCloser
	charFormatters        []ByteFormatter            // list of byte displayer(s) for data
	charFormatterCount    int                        // shorthand for len(charFormatters), for speeding up
	offsetFormatter       []OffsetFormatter          // offset formatters (max 2) first one is displayed on the left side and second one on the right side
	offsetFormatterCount  int                        // shorhand for len(offsetFormatter), for speeding up
	fileSize              int64                      // file size reference, -1 means STDIN. Hint for offset formatter(s) for how many padding characters to use.
	ReadBytes             uint64                     // How many bytes Reader has been reading so far (for limit)
	sb                    strings.Builder            // Faster than concatenating strings
	Splitter              string                     // Splitter character for columns
	offsetFormatterFormat map[OffsetFormatter]string // Printf format for offset format
	offsetFormatterWidth  map[OffsetFormatter]int    // How much padding width needed, calculated from fileSize variable
	Colors                Colors                     // Colors
	growHint              int                        // Grow hint for sb strings.Builder variable for speed
}

func New(r iface.ReadSeekerCloser, offsetFormatter []OffsetFormatter, formatters []ByteFormatter, palette [256]color.AnsiColor, filesize int64) *Reader {
	if formatters == nil {
		panic(`nil formatter`)
	}

	var calcpalette [256]string

	for idx := range palette {
		calcpalette[idx] = fmt.Sprintf(`%s%s`, color.SetForeground, palette[idx].String())
	}

	reader := &Reader{
		r:                    r,
		fileSize:             filesize,
		charFormatters:       formatters,
		offsetFormatter:      offsetFormatter,
		ReadBytes:            0, // How many byte's we've read
		sb:                   strings.Builder{},
		Splitter:             `┊`, // Splitter character between different columns
		charFormatterCount:   len(formatters),
		offsetFormatterCount: len(offsetFormatter),
		Colors: Colors{
			palette:       calcpalette,
			SplitterColor: color.AnsiColor{Color: color.ColorGrey93_eeeeee},
			OffsetColor:   color.AnsiColor{Color: color.ColorGrey93_eeeeee},
			specialBreak:  fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey35_585858}),
			hilightBreak:  fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey100_ffffff}),
		},
	}

	reader.offsetFormatterFormat = make(map[OffsetFormatter]string, reader.offsetFormatterCount)
	reader.offsetFormatterWidth = make(map[OffsetFormatter]int, reader.offsetFormatterCount)

	reader.Colors.splitterBreak = fmt.Sprintf(`%s%s`, color.SetForeground, reader.Colors.SplitterColor)
	reader.Colors.offsetBreak = fmt.Sprintf(`%s%s`, color.SetForeground, reader.Colors.OffsetColor)

	for _, f := range reader.charFormatters {
		switch f {
		case ViewHex:
			reader.growHint += 49
		case ViewDec, ViewOct:
			reader.growHint += 65
		case ViewASCII:
			reader.growHint += 18
		case ViewDecWithASCII:
			reader.growHint += 129
		case ViewHexWithASCII:
			reader.growHint += 113
		}
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

func (r *Reader) formatOffset(formatter OffsetFormatter, offset uint64) {
	switch formatter {
	case OffsetPercent:
		percent := (float64(offset) * 100.0) / float64(r.fileSize)
		r.sb.WriteString(fmt.Sprintf(r.offsetFormatterFormat[formatter], percent))
	default:
		r.sb.WriteString(fmt.Sprintf(r.offsetFormatterFormat[formatter], offset))
	}
}

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

// Read reads 16 bytes and provides string to display
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
					r.sb.WriteString(r.Colors.palette[tmp[i]])
				}

				switch byteFormatterType {
				case ViewHex:
					r.sb.WriteString(fmt.Sprintf(`%02x`, tmp[i]))
				case ViewDec:
					r.sb.WriteString(fmt.Sprintf(`%03d`, tmp[i]))
				case ViewOct:
					r.sb.WriteString(fmt.Sprintf(`%03o`, tmp[i]))
				case ViewBit, ViewBitWithDec, ViewBitWithHex, ViewBitWithAsc:
					if byteFormatterType == ViewBitWithAsc {
						r.sb.WriteString(r.Colors.palette[tmp[i]])
					}

					for idx, ru := range fmt.Sprintf(`%08b`, tmp[i]) {
						if idx == 0 {
							r.sb.WriteString(color.SetUnderlineOn)
						}

						r.sb.WriteRune(ru)

						if idx == 3 {
							r.sb.WriteString(color.SetUnderlineOff)
						}
					}

					switch byteFormatterType {
					case ViewBitWithDec:
						r.sb.WriteString(` `)
						r.sb.WriteString(fmt.Sprintf(`%03d`, tmp[i]))
					case ViewBitWithHex:
						r.sb.WriteString(` `)
						r.sb.WriteString(fmt.Sprintf(`%02x`, tmp[i]))
					case ViewBitWithAsc:
						r.sb.WriteString(` `)
						r.sb.WriteString(r.Colors.specialBreak)
						r.sb.WriteString(`[`)
						r.sb.WriteString(r.Colors.hilightBreak)
						r.sb.WriteString(fmt.Sprintf(`%c`, asciiByteToChar[tmp[i]]))
						r.sb.WriteString(r.Colors.specialBreak)
						r.sb.WriteString(`]`)
					}

				case ViewASCII:
					r.sb.WriteString(fmt.Sprintf(`%c`, asciiByteToChar[tmp[i]]))
				case ViewHexWithASCII:
					r.sb.WriteString(r.Colors.palette[tmp[i]])
					r.sb.WriteString(fmt.Sprintf(`%02x `, tmp[i]))
					r.sb.WriteString(r.Colors.specialBreak)
					r.sb.WriteString(`[`)
					r.sb.WriteString(r.Colors.hilightBreak)
					r.sb.WriteString(fmt.Sprintf(`%c`, asciiByteToChar[tmp[i]]))
					r.sb.WriteString(r.Colors.specialBreak)
					r.sb.WriteString(`]`)
				case ViewDecWithASCII:
					r.sb.WriteString(r.Colors.palette[tmp[i]])
					r.sb.WriteString(fmt.Sprintf(`%03d `, tmp[i]))
					r.sb.WriteString(r.Colors.specialBreak)
					r.sb.WriteString(`[`)
					r.sb.WriteString(r.Colors.hilightBreak)
					r.sb.WriteString(fmt.Sprintf(`%c`, asciiByteToChar[tmp[i]]))
					r.sb.WriteString(r.Colors.specialBreak)
					r.sb.WriteString(`]`)
				}

				if i < 15 {
					switch byteFormatterType {
					case ViewASCII:
						continue
					default:
						r.sb.WriteString(` `)
					}
				}
			} else {
				// There is no data so we add padding
				if i == 0 || (i > 0 && tmp[i] != tmp[i-1]) {
					// Only print on first and changed color
					r.sb.WriteString(r.Colors.palette[0])
				}

				r.sb.WriteString(formatterPaddingMap[byteFormatterType])

				if i < 15 {
					switch byteFormatterType {
					case ViewASCII:
						continue
					default:
						r.sb.WriteString(` `)
					}
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
