package reader

import (
	"fmt"
)

type OffsetFormatter uint8

const (
	OffsetHex OffsetFormatter = iota
	OffsetDec
	OffsetOct
	OffsetPercent
)

// Get enum from string
var offsetFormattersStringToEnumMap = map[string]OffsetFormatter{
	`hex`: OffsetHex,
	`dec`: OffsetDec,
	`oct`: OffsetOct,
	`per`: OffsetPercent,
}

// GetOffsetFormatters parses string and returns proper formatter(s)
func GetOffsetFormatters(viewerStr []string) (formatters []OffsetFormatter, err error) {

	if len(viewerStr) > 2 {
		return nil, fmt.Errorf(`error: max two formatters, got: %v`, viewerStr)
	}

	//formatters = make([]iface.OffsetFormatter, 0)

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

		formatters = append(formatters, en)

	}

	return formatters, nil
}
