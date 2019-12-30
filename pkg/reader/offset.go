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

func GetOffsetViewer(viewerStr []string) (formatters []iface.OffsetFormatter, err error) {

	if len(viewerStr) > 2 {
		return nil, fmt.Errorf(`error: max two formatters, got: %v`, viewerStr)
	}

	formatters = make([]iface.OffsetFormatter, 0)

	for _, v := range viewerStr {
		if v == `` {
			continue
		}

		en, ok := offsetViewersStringToEnumMap[v]
		if !ok {
			return nil, fmt.Errorf(`invalid: %v`, viewerStr)
		}

		formatters = append(formatters, offsetViewers[en])

	}

	return formatters, nil

}
