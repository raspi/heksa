package main

import (
	"./display"
	"fmt"
	clr "github.com/logrusorgru/aurora"
)

type dataViewer uint8

const (
	ViewHex dataViewer = iota
	ViewDec
	ViewASCII
	ViewBit
)

var viewerEnumMap = map[dataViewer]Views{
	ViewHex:   display.NewHex(),
	ViewASCII: display.NewAscii(),
	ViewBit:   display.NewBit(),
	ViewDec:   display.NewDec(),
}

var viewersS = map[string]dataViewer{
	`hex`: ViewHex,
	`asc`: ViewASCII,
	`bit`: ViewBit,
	`dec`: ViewDec,
}

type Views interface {
	Display([]byte) string
	SetPalette(map[uint8]clr.Color)
}

// getViewers returns viewers from string separated by ','
func getViewers(viewers []string) (ds []Views, err error) {

	for _, v := range viewers {
		en, ok := viewersS[v]
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
