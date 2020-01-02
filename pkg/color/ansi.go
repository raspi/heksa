package color

import "strconv"

// https://en.wikipedia.org/wiki/ANSI_escape_code

type AnsiColor struct {
	Color     Color
	Bold      bool
	Underline bool
	Blink     bool
	Crossed   bool
	Italic    bool
}

func (ac AnsiColor) String() string {
	out := strconv.Itoa(int(ac.Color)) + `m`
	return out
}
