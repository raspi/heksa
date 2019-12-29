package display

/*
nearest returns nearest bits to eight bits

1-7 = 8
9-15 = 16

and so on
*/
func nearest(bitWidth uint8) uint8 {
	return (bitWidth + (8 - 1)) & ^(bitWidth - 1)
}
