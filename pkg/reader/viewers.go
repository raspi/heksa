package reader

import (
	"fmt"
	"sort"
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
}

var formatterPaddingMap = map[ByteFormatter]string{
	ViewASCII:        `‡`,
	ViewHex:          `‡‡`,
	ViewDec:          `‡‡‡`,
	ViewOct:          `‡‡‡`,
	ViewHexWithASCII: `‡‡‡‡‡‡`,
	ViewDecWithASCII: `‡‡‡‡‡‡‡`,
	ViewBit:          `‡‡‡‡‡‡‡‡`,
	ViewBitWithDec:   `‡‡‡‡‡‡‡‡‡‡‡‡`,
}

// getViewers returns viewers from string separated by ','
func GetViewers(viewers []string) (ds []ByteFormatter, err error) {
	for _, v := range viewers {
		en, ok := formatterStringToEnumMap[v]
		if !ok {
			return nil, fmt.Errorf(`invalid: %v`, v)
		}

		ds = append(ds, en)
	}

	if len(ds) == 0 {
		return nil, fmt.Errorf(`there has to be at least one viewer`)
	}

	return ds, nil
}

func GetViewerList() (viewers []string) {
	for s, _ := range formatterStringToEnumMap {
		viewers = append(viewers, s)
	}

	sort.Strings(viewers)
	return viewers
}
