package main

import "io"

type Reader struct {
	r           io.ReadSeeker
	displays    []Views     // displayer(s) for data
	offdisplats ShowsOffset // offset displayer
}

func New(r io.ReadSeeker, offdisplays ShowsOffset, displays []Views) *Reader {
	if offdisplays == nil {
		panic(`nil offset displayer`)
	}

	if displays == nil {
		panic(`nil displayer(s)`)
	}

	return &Reader{
		r:           r,
		displays:    displays,
		offdisplats: offdisplays,
	}
}

// Read reads 16 bytes and provides string to display
func (r Reader) Read() (string, error) {
	out := ``
	out += r.offdisplats.DisplayOffset(r.r)
	out += ` | `

	tmp := make([]byte, 16)
	rb, err := r.r.Read(tmp)
	if err != nil {
		return ``, err
	}

	for _, dplay := range r.displays {
		out += dplay.Display(tmp[0:rb])
		out += ` | `
	}

	return out, nil
}
