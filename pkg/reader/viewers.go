package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/display"
	"github.com/raspi/heksa/pkg/iface"
)

type byteFormatter uint8

const (
	ViewHex byteFormatter = iota
	ViewDec
	ViewOct
	ViewASCII
	ViewBit
)

var formatterEnumToImplMap = map[byteFormatter]iface.CharacterFormatter{
	ViewHex:   display.NewHex(),
	ViewASCII: display.NewAscii(),
	ViewBit:   display.NewBit(),
	ViewDec:   display.NewDec(),
	ViewOct:   display.NewOct(),
}

// Get enum from string
var formatterStringToEnumMap = map[string]byteFormatter{
	`hex`: ViewHex,
	`asc`: ViewASCII,
	`bit`: ViewBit,
	`dec`: ViewDec,
	`oct`: ViewOct,
}

// getViewers returns viewers from string separated by ','
func GetViewers(viewers []string) (ds []iface.CharacterFormatter, err error) {

	for _, v := range viewers {
		en, ok := formatterStringToEnumMap[v]
		if !ok {
			return nil, fmt.Errorf(`invalid: %v`, v)
		}

		ds = append(ds, formatterEnumToImplMap[en])
	}

	if len(ds) == 0 {
		return nil, fmt.Errorf(`there has to be at least one viewer`)
	}

	return ds, nil
}
