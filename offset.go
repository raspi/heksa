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
	OffsetPercent
)

var offsetViewers = map[offsetViewer]ShowsOffset{
	OffsetHex:     display.NewHex(),
	OffsetDec:     display.NewDec(),
	OffsetPercent: display.NewPercent(),
}

var offsetViewersStringToEnumMap = map[string]offsetViewer{
	`hex`: OffsetHex,
	`dec`: OffsetDec,
	`per`: OffsetPercent,
}

// ShowsOffset is interface for displaying file offset in X format (where X might be hex, decimal, octal, ..)
type ShowsOffset interface {
	DisplayOffset(r io.ReadSeeker) string
	SetFileSize(int64) // For leading zeros information
}

func getOffsetViewer(viewerStr string) (ShowsOffset, error) {
	en, ok := offsetViewersStringToEnumMap[viewerStr]
	if !ok {
		return nil, fmt.Errorf(`invalid: %v`, viewerStr)
	}

	return offsetViewers[en], nil
}
