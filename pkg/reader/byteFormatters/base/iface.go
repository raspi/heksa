package base

import (
	"fmt"
	"github.com/raspi/heksa/pkg/color"
	"strings"
)

type ByteFormatter interface {
	Print(byte) string
	// How many characters formatter will print (1-N)
	// Used for padding and grow hint
	GetPrintSize() int
}

type FormatterGroup struct {
	palette            [256]string
	specialBreak       string
	hilightBreak       string
	changePalette      bool
	formatters         []ByteFormatter
	Width              int
	sb                 strings.Builder
	visualSplitterSize int
	visualSplitter     string // Visual splitter that gets inserted every visualSplitterSize bytes
	formatterCount     int
	Splitter           string
	splitterBreak      string
}

func New(formatters []ByteFormatter, palette [256]string, SpecialBreak string, HilightBreak string, width uint16, visualSplitterSize uint8) FormatterGroup {
	if formatters == nil {
		panic(`nil formatter`)
	}

	if width == 0 {
		panic(`zero width`)
	}

	return FormatterGroup{
		palette:            palette,
		specialBreak:       SpecialBreak,
		hilightBreak:       HilightBreak,
		formatters:         formatters,
		changePalette:      true,
		Width:              int(width),
		sb:                 strings.Builder{},
		visualSplitterSize: int(visualSplitterSize),
		visualSplitter:     ` `,
		formatterCount:     len(formatters),
		Splitter:           `┊`, // Splitter character between different columns
		splitterBreak:      fmt.Sprintf(`%s%s`, color.SetForeground, color.AnsiColor{Color: color.ColorGrey93_eeeeee}),
	}
}

func (fg *FormatterGroup) Print(tmp []byte) string {
	fg.sb.Reset()

	paddingIndex := len(tmp)

	// iterate through every formatter which outputs it's own format
	for didx, byteFormatterType := range fg.formatters {
		// First character to print, so always true
		fg.changePalette = true

		for i := 0; i < fg.Width; i++ {
			if fg.visualSplitterSize != 0 && i != 0 && i%fg.visualSplitterSize == 0 {
				// Add pad for better visualization every visualSplitterSize bytes
				fg.sb.WriteString(fg.visualSplitter)
			}

			if paddingIndex > i {
				if i == 0 || (i > 0 && tmp[i] != tmp[i-1] && fg.palette[tmp[i]] != fg.palette[tmp[i-1]]) {
					fg.changePalette = true
				}

				if fg.changePalette {
					fg.sb.WriteString(fg.palette[tmp[i]])
				}

				fg.sb.WriteString(byteFormatterType.Print(tmp[i]))

				if i < (fg.Width-1) && byteFormatterType.GetPrintSize() > 1 {
					fg.sb.WriteString(` `)
				}
			} else {
				// No data available, add padding

				if i == paddingIndex {
					// We're at start of padding
					fg.sb.WriteString(fg.palette[0])
				}

				fg.sb.WriteString(strings.Repeat(`‡`, byteFormatterType.GetPrintSize()))

				if i < (fg.Width-1) && (byteFormatterType.GetPrintSize() > 1) {
					fg.sb.WriteString(` `)
				}
			}
		}

		if didx < (fg.formatterCount - 1) {
			fg.sb.WriteString(fg.splitterBreak)
			fg.sb.WriteString(fg.Splitter)
		}
	}

	return fg.sb.String()
}
