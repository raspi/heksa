package display

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
