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

type ReaderColors struct {
	LineOdd  string // background color
	LineEven string // background color
	Offset   string
	Splitter string
}

type Reader struct {
	r                      iface.ReadSeekerCloser
	offsetFormatters       []offFormatters.OffsetFormatter // offset formatters (max 2) first one is displayed on the left side and second one on the right side
	offsetFormatterCount   int                             // shorthand for len(offsetFormatters), for speeding up
	isStdin                bool                            // Are we reading from STDIN? if so, we can't ask for offset position from file
	readTotalBytes         uint64                          // How many bytes Reader has been reading so far (for limit)
	readRelativeTotalBytes uint64                          // How many bytes Reader has been reading relatively
	sb                     strings.Builder                 // Faster than concatenating strings
	Splitter               string                          // Splitter character for columns
	growHint               int                             // Grow hint for sb strings.Builder variable for speed
	formatterGroup         base.FormatterGroup
	colors                 ReaderColors
	isEven                 bool // change background color for printed line
	printRelativeOffset    bool // Print relative offset?
}

func New(r iface.ReadSeekerCloser, offsetFormatter []offFormatters.OffsetFormatter, colors ReaderColors, formatterGroup base.FormatterGroup, isStdin bool, useRelativeOffset bool) *Reader {

	if isStdin {
		useRelativeOffset = false
	}

	reader := &Reader{
		r:                      r,
		isStdin:                isStdin,
		offsetFormatters:       offsetFormatter,
		readTotalBytes:         0, // How many bytes we've read
		readRelativeTotalBytes: 0, // How many bytes we've read relatively (always starts from zero)
		sb:                     strings.Builder{},
		Splitter:               `┊`, // Splitter character between different columns
		offsetFormatterCount:   len(offsetFormatter),
		formatterGroup:         formatterGroup,
		colors:                 colors,
		isEven:                 false,
		printRelativeOffset:    useRelativeOffset,
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
		r.sb.WriteString(r.colors.Offset)
		// show offset on the left side
		r.sb.WriteString(r.offsetFormatters[0].Print(offset))
		r.sb.WriteString(r.colors.Splitter)
		r.sb.WriteString(r.Splitter)
	}

	return r.sb.String()
}

// getoffsetRight outputs the selected formatter on the right side after all the user selected byte formatters
func (r *Reader) getoffsetRight(offset uint64) string {
	r.sb.Reset()
	if r.offsetFormatterCount > 1 {
		// show offset on the right side
		r.sb.WriteString(r.colors.Splitter)
		r.sb.WriteString(r.Splitter)
		r.sb.WriteString(r.colors.Offset)
		r.sb.WriteString(r.offsetFormatters[1].Print(offset))
	}

	return r.sb.String()
}

// Read reads N (r.width) bytes and provides string to display
func (r *Reader) Read() (string, error) {
	var offset uint64

	if r.isStdin {
		// reading from STDIN, can't use seek
		offset = r.readTotalBytes
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

	offsetLeftRelative := r.getoffsetLeft(r.readRelativeTotalBytes)
	offsetRightRelative := r.getoffsetRight(r.readRelativeTotalBytes)

	r.sb.Reset()
	r.sb.Grow(r.growHint)

	// Fetch bytes with selected formatters
	tmp := make([]byte, r.formatterGroup.Width)
	bytesReadCount, err := r.r.Read(tmp)
	if err != nil {
		return ``, err
	}

	r.readTotalBytes += uint64(bytesReadCount)
	r.readRelativeTotalBytes += uint64(bytesReadCount)

	// Change between two background colors
	if r.isEven {
		r.sb.WriteString(r.colors.LineEven)
	} else {
		r.sb.WriteString(r.colors.LineOdd)
	}
	r.isEven = !r.isEven // Flip between true -> false | false -> true

	// Offset on the left
	r.sb.WriteString(offsetLeft)

	if r.printRelativeOffset {
		// Print relative offset
		r.sb.WriteString(offsetLeftRelative)
	}

	// Print the formatted bytes
	r.sb.WriteString(r.formatterGroup.Print(tmp[0:bytesReadCount]))

	if r.printRelativeOffset {
		// Print relative offset
		r.sb.WriteString(offsetRightRelative)
	}

	// Offset on the right
	r.sb.WriteString(offsetRight)

	// clear ANSI code so that terminal doesn't explode
	r.sb.WriteString(color.Clear)

	return r.sb.String(), nil
}

func (r *Reader) GetReadBytes() uint64 {
	return r.readTotalBytes
}
