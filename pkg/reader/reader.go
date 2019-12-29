package reader

import "github.com/raspi/heksa/pkg/iface"

type Reader struct {
	r               iface.ReadSeekerCloser
	displays        []iface.Views     // displayer(s) for data
	offsetFormatter iface.ShowsOffset // offset displayer
	ReadBytes       uint64
}

func New(r iface.ReadSeekerCloser, offsetFormatter iface.ShowsOffset, formatters []iface.Views) *Reader {
	if offsetFormatter == nil {
		panic(`nil offset displayer`)
	}

	if formatters == nil {
		panic(`nil displayer(s)`)
	}

	return &Reader{
		r:               r,
		displays:        formatters,
		offsetFormatter: offsetFormatter,
		ReadBytes:       0,
	}
}

// Read reads 16 bytes and provides string to display
func (r *Reader) Read() (string, error) {
	out := ``
	out += r.offsetFormatter.DisplayOffset(r.r)
	out += ` | `

	tmp := make([]byte, 16)
	rb, err := r.r.Read(tmp)
	if err != nil {
		return ``, err
	}

	r.ReadBytes += uint64(rb)

	for _, dplay := range r.displays {
		out += dplay.Display(tmp[0:rb])
		out += ` | `
	}

	return out, nil
}
