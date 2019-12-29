package main

import (
	"fmt"
	clr "github.com/logrusorgru/aurora"
	"github.com/raspi/heksa/display"
)

type dataViewer uint8

const (
	ViewHex dataViewer = iota
	ViewDec
	ViewOct
	ViewASCII
	ViewBit
)

var viewerEnumMap = map[dataViewer]Views{
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

type Views interface {
	Display([]byte) string
	SetPalette(map[uint8]clr.Color)
}

// getViewers returns viewers from string separated by ','
func getViewers(viewers []string) (ds []Views, err error) {

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
