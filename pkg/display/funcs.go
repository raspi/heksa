package display

import "fmt"

/*
nearest returns nearest bits to uint8-64 len

1-7 = 8
9-15 = 16

and so on
*/
func nearest(bitWidth uint8) uint8 {
	if bitWidth > 64 {
		return 64
	} else if bitWidth > 32 {
		return 64
	} else if bitWidth > 16 {
		return 32
	} else if bitWidth > 8 {
		return 16
	}
	return 8
}

// header returns header for formatters of different lengths
// for example:
// - bit formatter displays 8 characters (00000000-11111111)
// - dec formatter displays 3 characters (000-255)
func header(l uint8) (out string) {
	format := fmt.Sprintf(`%%0%vx`, l)
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
