package main

import (
	"flag"
	"fmt"
	"io"
	"math/bits"
	"os"
	"strings"
)

var VERSION = `v0.0.0`
var BUILD = `dev`

const AUTHOR = `Pekka Järvinen`
const HOMEPAGE = `https://github.com/raspi/heksa`

func main() {

	offsetDisplayS := flag.String(`o`, `hex`, `offset displayer`)
	displayS := flag.String(`d`, `hex,asc`, `displayer`)

	flag.Usage = func() {
		fmt.Printf(`heksa - hex file dumper %v build %v`+"\n", VERSION, BUILD)
		fmt.Printf(`(c) %v 2019- - %v`+"\n", AUTHOR, HOMEPAGE)
		fmt.Println()
		fmt.Println(`Usage:`)
		fmt.Printf(`  %v <file>`+"\n", os.Args[0])
		fmt.Printf(`  %v -o <offset viewer> <file>`+"\n", os.Args[0])
		fmt.Printf(`  %v -d <data viewer(s)> <file>`+"\n", os.Args[0])
		fmt.Println()
		fmt.Println(`Example:`)
		fmt.Printf(`  %v example.dat`+"\n", os.Args[0])
		fmt.Printf(`  %v -o hex -d hex,asc,bit example.dat`+"\n", os.Args[0])
		fmt.Println()
		fmt.Println(`Offset viewers, defaults to hex:`)
		fmt.Println(`  hex  Hex (0x00 - 0xFF)`)
		fmt.Println(`  dec  Decimal (000 - 255)`)
		fmt.Println()
		fmt.Println(`Data viewers (can be combined with ','), defaults to hex,asc:`)
		fmt.Println(`  hex  Hex (0x00 - 0xFF)`)
		fmt.Println(`  dec  Decimal (000 - 255)`)
		fmt.Println(`  asc  ASCII`)
		fmt.Println(`  bit  Bits (00000000 - 11111111)`)
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	offViewer, err := getOffsetViewer(*offsetDisplayS)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error getting offset displayer: %v`, err))
		os.Exit(1)
	}

	displays, err := getViewers(strings.Split(*displayS, `,`))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error getting data displayer: %v`, err))
		os.Exit(1)
	}

	fpath := flag.Arg(0)

	f, err := os.Open(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error opening file: %v`, err))
		os.Exit(1)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error stat'ing file: %v`, err))
		os.Exit(1)
	}

	bitWidth := uint8(bits.Len64(uint64(fi.Size())))
	bitWidth = (bitWidth + (8 - 1)) & ^(bitWidth - 1)

	offViewer.SetBitWidthSize(bitWidth)

	for idx, _ := range displays {
		displays[idx].SetPalette(defaultCharacterColors)
	}

	r := New(f, offViewer, displays)

	for {
		s, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Fprintln(os.Stderr, fmt.Sprintf(`error while reading file: %v`, err))
			os.Exit(1)
		}

		fmt.Println(s)
	}
}
