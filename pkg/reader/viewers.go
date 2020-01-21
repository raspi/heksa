package reader

import (
	"fmt"
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
}

// Padding when we are at EOF
var formatterPaddingMap = map[ByteFormatter]string{
	ViewASCII:        `‡`,
	ViewHex:          `‡‡`,
	ViewDec:          `‡‡‡`,
	ViewOct:          `‡‡‡`,
	ViewHexWithASCII: `‡‡‡‡‡‡`,
	ViewDecWithASCII: `‡‡‡‡‡‡‡`,
	ViewBit:          `‡‡‡‡‡‡‡‡`,
	ViewBitWithHex:   `‡‡‡‡‡‡‡‡‡‡‡`,
	ViewBitWithDec:   `‡‡‡‡‡‡‡‡‡‡‡‡`,
	ViewBitWithAsc:   `‡‡‡‡‡‡‡‡‡‡‡‡`,
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
