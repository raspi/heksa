package base

type ByteFormatter interface {
	Print(byte) string
	// How many characters formatter will print (1-N)
	// Used for padding and grow hint
	GetPrintSize() int
}

var Palette [256]string

var ChangePalette bool

// Palettes for different splitters
var SpecialBreak string
var HilightBreak string
