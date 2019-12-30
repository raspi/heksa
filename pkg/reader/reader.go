package reader

import (
	"github.com/raspi/heksa/pkg/iface"
	"strings"
)

type Reader struct {
	r               iface.ReadSeekerCloser
	displays        []iface.Views     // displayer(s) for data
	offsetFormatter iface.ShowsOffset // offset displayer
	ReadBytes       uint64
	sb              strings.Builder
	Splitter        string
}

func New(r iface.ReadSeekerCloser, offsetFormatter iface.ShowsOffset, formatters []iface.Views) *Reader {
	if offsetFormatter == nil {
		panic(`nil offset displayer`)
	}

	if formatters == nil {
		panic(`nil displayer(s)`)
	}

	reader := &Reader{
		r:               r,
		displays:        formatters,
		offsetFormatter: offsetFormatter,
		ReadBytes:       0,
		sb:              strings.Builder{},
		Splitter:        `|`,
	}

	return reader
}

// Read reads 16 bytes and provides string to display
func (r *Reader) Read() (string, error) {
	r.sb.Reset()
	r.sb.Grow(1024)

	r.sb.WriteString(r.offsetFormatter.DisplayOffset(r.r))
	r.sb.WriteString(r.Splitter)

	tmp := make([]byte, 16)
	rb, err := r.r.Read(tmp)
	if err != nil {
		return ``, err
	}

	r.ReadBytes += uint64(rb)

	for _, dplay := range r.displays {
		for i := 0; i < 16; i++ {
			if i == 8 {
				r.sb.WriteString(` `)
			}

			if rb > i {
				r.sb.WriteString(dplay.Display(tmp[i]))
			} else {
				r.sb.WriteString(dplay.EofStr())
			}
		}

		r.sb.WriteString(r.Splitter)
	}

	return r.sb.String(), nil
}
