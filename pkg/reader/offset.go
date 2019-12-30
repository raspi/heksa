package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/display"
	"github.com/raspi/heksa/pkg/iface"
)

type offsetViewer uint8

const (
	OffsetHex offsetViewer = iota
	OffsetDec
	OffsetOct
	OffsetPercent
)

var offsetViewers = map[offsetViewer]iface.OffsetFormatter{
	OffsetHex:     display.NewHex(),
	OffsetDec:     display.NewDec(),
	OffsetOct:     display.NewOct(),
	OffsetPercent: display.NewPercent(),
}

var offsetViewersStringToEnumMap = map[string]offsetViewer{
	`hex`: OffsetHex,
	`dec`: OffsetDec,
	`oct`: OffsetOct,
	`per`: OffsetPercent,
}

func GetOffsetViewer(viewerStr string) (iface.OffsetFormatter, error) {
	en, ok := offsetViewersStringToEnumMap[viewerStr]
	if !ok {
		return nil, fmt.Errorf(`invalid: %v`, viewerStr)
	}

	return offsetViewers[en], nil
}
