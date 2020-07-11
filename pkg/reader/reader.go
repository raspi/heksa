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
	offsetBreak   string
	splitterBreak string
}

type Reader struct {
	r                    iface.ReadSeekerCloser
	offsetFormatters     []offFormatters.OffsetFormatter // offset formatters (max 2) first one is displayed on the left side and second one on the right side
	offsetFormatterCount int                             // shorthand for len(offsetFormatters), for speeding up
	isStdin              bool
	ReadBytes            uint64          // How many bytes Reader has been reading so far (for limit)
	sb                   strings.Builder // Faster than concatenating strings
	Splitter             string          // Splitter character for columns
	Colors               Colors          // Colors
	growHint             int             // Grow hint for sb strings.Builder variable for speed
	formatterGroup       base.FormatterGroup
}

func New(r iface.ReadSeekerCloser, offsetFormatter []offFormatters.OffsetFormatter, formatterGroup base.FormatterGroup, isStdin bool) *Reader {
	reader := &Reader{
		r:                    r,
		isStdin:              isStdin,
		offsetFormatters:     offsetFormatter,
		ReadBytes:            0, // How many bytes we've read
		sb:                   strings.Builder{},
		Splitter:             `┊`, // Splitter character between different columns
		offsetFormatterCount: len(offsetFormatter),
		formatterGroup:       formatterGroup,
		Colors: Colors{
			offsetBreak:   fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey82_d0d0d0}),
			splitterBreak: fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey93_eeeeee}),
		},
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

	// Offset on the left
	r.sb.WriteString(offsetLeft)

	// Fetch bytes with selected formatters
	tmp := make([]byte, r.formatterGroup.Width)
	bytesReadCount, err := r.r.Read(tmp)
	if err != nil {
		return ``, err
	}

	r.ReadBytes += uint64(bytesReadCount)

	r.sb.WriteString(r.formatterGroup.Print(tmp[0:bytesReadCount]))

	// Offset on the right
	r.sb.WriteString(offsetRight)

	return r.sb.String(), nil
}
