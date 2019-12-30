package reader

import (
	"fmt"
	"github.com/raspi/heksa/pkg/display"
	"github.com/raspi/heksa/pkg/iface"
)

type dataViewer uint8

const (
	ViewHex dataViewer = iota
	ViewDec
	ViewOct
	ViewASCII
	ViewBit
)

var viewerEnumMap = map[dataViewer]iface.CharacterFormatter{
	ViewHex:   display.NewHex(),
	ViewASCII: display.NewAscii(),
	ViewBit:   display.NewBit(),
	ViewDec:   display.NewDec(),
	ViewOct:   display.NewOct(),
}

// Get enum from string
var viewersStringToEnumMap = map[string]dataViewer{
	`hex`: ViewHex,
	`asc`: ViewASCII,
	`bit`: ViewBit,
	`dec`: ViewDec,
	`oct`: ViewOct,
}

// getViewers returns viewers from string separated by ','
func GetViewers(viewers []string) (ds []iface.CharacterFormatter, err error) {

	for _, v := range viewers {
		en, ok := viewersStringToEnumMap[v]
		if !ok {
			return nil, fmt.Errorf(`invalid: %v`, v)
		}

		ds = append(ds, viewerEnumMap[en])
	}

	if len(ds) == 0 {
		return nil, fmt.Errorf(`there has to be at least one viewer`)
	}

	return ds, nil
}
