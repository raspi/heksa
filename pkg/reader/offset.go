package reader

import (
	"fmt"
	offFormatters "github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/decimal"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/hex"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/human"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/octal"
	"github.com/raspi/heksa/pkg/reader/offsetFormatters/percent"
	"sort"
	"strings"
)

type OffsetFormatter uint8

const (
	OffsetHex      OffsetFormatter = iota // Hexadecimal
	OffsetDec                             // Decimal
	OffsetOct                             // Octal
	OffsetPercent                         // Percentage 0-100 from offset and filesize
	OffsetHumanSI                         // Offset in human form (SI) 1000
	OffsetHumanIEC                        // Offset in human form (IEC) 1024
)

// Get enum from string
var offsetFormattersStringToEnumMap = map[string]OffsetFormatter{
	`hex`:    OffsetHex,
	`dec`:    OffsetDec,
	`oct`:    OffsetOct,
	`per`:    OffsetPercent,
	`humiec`: OffsetHumanIEC,
	`humsi`:  OffsetHumanSI,
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

func GetFromOffsetFormatter(formatter OffsetFormatter, info offFormatters.BaseInfo) offFormatters.OffsetFormatter {
	var f offFormatters.OffsetFormatter
	switch formatter {
	case OffsetDec:
		f = decimal.New(info)
	case OffsetHex:
		f = hex.New(info)
	case OffsetOct:
		f = octal.New(info)
	case OffsetPercent:
		f = percent.New(info)
	case OffsetHumanSI: // 1000
		f = human.New(info, 1000)
	case OffsetHumanIEC: // 1024
		f = human.New(info, 1024)
	default:
		return nil
	}

	return f
}
