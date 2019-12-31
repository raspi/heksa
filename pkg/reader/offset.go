package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/display"
	"github.com/raspi/heksa/pkg/iface"
)

type offsetFormatter uint8

const (
	OffsetHex offsetFormatter = iota
	OffsetDec
	OffsetOct
	OffsetPercent
)

var offsetFormattersEnumToImpl = map[offsetFormatter]iface.OffsetFormatter{
	OffsetHex:     display.NewHex(),
	OffsetDec:     display.NewDec(),
	OffsetOct:     display.NewOct(),
	OffsetPercent: display.NewPercent(),
}

// Get enum from string
var offsetFormattersStringToEnumMap = map[string]offsetFormatter{
	`hex`: OffsetHex,
	`dec`: OffsetDec,
	`oct`: OffsetOct,
	`per`: OffsetPercent,
}

// GetOffsetFormatters parses string and returns proper formatter(s)
func GetOffsetFormatters(viewerStr []string) (formatters []iface.OffsetFormatter, err error) {

	if len(viewerStr) > 2 {
		return nil, fmt.Errorf(`error: max two formatters, got: %v`, viewerStr)
	}

	formatters = make([]iface.OffsetFormatter, 0)

	for _, v := range viewerStr {
		if v == `no` {
			v = ``
		}

		if v == `` {
			continue
		}

		en, ok := offsetFormattersStringToEnumMap[v]
		if !ok {
			return nil, fmt.Errorf(`invalid: %v`, viewerStr)
		}

		formatters = append(formatters, offsetFormattersEnumToImpl[en])

	}

	return formatters, nil
}
