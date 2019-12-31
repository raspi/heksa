package display

import "fmt"

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
