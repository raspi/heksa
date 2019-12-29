package main

import (
	"fmt"
	"github.com/raspi/heksa/display"
	"io"
)

type offsetViewer uint8

const (
	OffsetHex offsetViewer = iota
	OffsetDec
)

var offsetViewers = map[offsetViewer]ShowsOffset{
	OffsetHex: display.NewHex(),
	OffsetDec: display.NewDec(),
}

var viewersOffS = map[string]offsetViewer{
	`hex`: OffsetHex,
	`dec`: OffsetDec,
}

// ShowsOffset is interface for displaying file offset in X format (where X might be hex, decimal, octal, ..)
type ShowsOffset interface {
	DisplayOffset(r io.ReadSeeker) string
	SetBitWidthSize(uint8) // For leading zeros information
}

func getOffsetViewer(viewerStr string) (ShowsOffset, error) {
	en, ok := viewersOffS[viewerStr]
	if !ok {
		return nil, fmt.Errorf(`invalid: %v`, viewerStr)
	}

	return offsetViewers[en], nil
}
