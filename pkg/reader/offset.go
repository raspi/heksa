package reader

import (
	"fmt"
	"sort"
	"strings"
)

type OffsetFormatter uint8

const (
	OffsetHex     OffsetFormatter = iota // Hexadecimal
	OffsetDec                            // Decimal
	OffsetOct                            // Octal
	OffsetPercent                        // Percentage 0-100 from offset and filesize
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

	for _, v := range viewerStr {
		if v == `no` {
			v = ``
		}

		if v == `` {
			continue
		}

		en, ok := offsetFormattersStringToEnumMap[v]
		if !ok {
			return nil, fmt.Errorf(`invalid: %q, valid: %v`, v, strings.Join(GetOffsetViewerList(), `, `))
		}

		formatters = append(formatters, en)

	}

	return formatters, nil
}

// GetOffsetViewerList lists offset formatters as strings for usage information
func GetOffsetViewerList() (viewers []string) {
	for s := range offsetFormattersStringToEnumMap {
		viewers = append(viewers, s)
	}

	sort.Strings(viewers)
	return viewers
}
