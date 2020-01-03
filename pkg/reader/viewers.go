package reader

import (
	"fmt"
)

type ByteFormatter uint8

const (
	ViewHex ByteFormatter = iota
	ViewDec
	ViewOct
	ViewASCII
	ViewBit
	ViewHexWithASCII
	//ViewDecWithASCII
)

// Get enum from string
var formatterStringToEnumMap = map[string]ByteFormatter{
	`hex`:     ViewHex,
	`asc`:     ViewASCII,
	`bit`:     ViewBit,
	`dec`:     ViewDec,
	`oct`:     ViewOct,
	`hexwasc`: ViewHexWithASCII,
	//`decwasc`: ViewDecWithASCII,
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
