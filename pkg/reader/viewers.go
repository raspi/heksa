package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/ascii"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/bit"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/bitWithAscii"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/bitWithDecimal"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/bitWithHex"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/block"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/decWithAscii"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/decimal"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/hex"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/hexWithAscii"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/octal"
	"sort"
	"strings"
)

type ByteFormatter uint8

const (
	ViewHex ByteFormatter = iota // Hexadecimal
	ViewDec                      // Decimal
	ViewOct                      // Octal
	ViewASCII
	ViewBit          // Bits 00000000-11111111
	ViewHexWithASCII // Displays hex and ascii at same time
	ViewDecWithASCII // Displays dec and ascii at same time
	ViewBitWithDec   // Displays bits and decimal at same time
	ViewBitWithHex   // Displays bits and hex at same time
	ViewBitWithAsc   // Displays bits and ASCII at same time
	ViewBlock        // Display block character, for visualing patterns
)

// Get enum from string
var formatterStringToEnumMap = map[string]ByteFormatter{
	`hex`:     ViewHex,
	`asc`:     ViewASCII,
	`bit`:     ViewBit,
	`dec`:     ViewDec,
	`oct`:     ViewOct,
	`hexwasc`: ViewHexWithASCII,
	`decwasc`: ViewDecWithASCII,
	`bitwdec`: ViewBitWithDec,
	`bitwhex`: ViewBitWithHex,
	`bitwasc`: ViewBitWithAsc,
	`blk`:     ViewBlock,
}

// getViewers returns viewers from string separated by ','
func GetViewers(viewers []string) (ds []ByteFormatter, err error) {
	for _, v := range viewers {
		en, ok := formatterStringToEnumMap[v]
		if !ok {
			return nil, fmt.Errorf(`invalid: %q, valid: %v`, v, strings.Join(GetViewerList(), `, `))
		}

		ds = append(ds, en)
	}

	if len(ds) == 0 {
		return nil, fmt.Errorf(`there has to be at least one viewer`)
	}

	return ds, nil
}

// GetViewerList lists byte formatters as strings for usage information
func GetViewerList() (viewers []string) {
	for s := range formatterStringToEnumMap {
		viewers = append(viewers, s)
	}

	sort.Strings(viewers)
	return viewers
}

// GetByteFormatter gets implementation of given formatter
func GetByteFormatter(formatter ByteFormatter, hilightBreak string, specialBreak string) base.ByteFormatter {
	switch formatter {
	case ViewASCII:
		return ascii.New()
	case ViewHex:
		return hex.New()
	case ViewBit:
		return bit.New()
	case ViewDec:
		return decimal.New()
	case ViewOct:
		return octal.New()
	case ViewHexWithASCII:
		return hexWithAscii.New(hilightBreak, specialBreak)
	case ViewDecWithASCII:
		return decWithAscii.New(hilightBreak, specialBreak)
	case ViewBitWithAsc:
		return bitWithAscii.New(hilightBreak, specialBreak)
	case ViewBitWithDec:
		return bitWithDecimal.New(hilightBreak, specialBreak)
	case ViewBitWithHex:
		return bitWithHex.New(hilightBreak, specialBreak)
	case ViewBlock:
		return block.New()
	default:
		return nil
	}
}
