package main

import (
	"fmt"
	"github.com/DavidGamba/go-getoptions"
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
	opt := getoptions.New()

	opt.HelpSynopsisArgs(`<filename>`)

	offsetDisplayS := opt.StringOptional(`offset-display`, `hex`,
		opt.Alias(`o`),
		opt.ArgName(`offset format`),
		opt.Description(`One of: hex, dec`),
	)

	formatS := opt.StringOptional(`format`, `hex,asc`,
		opt.Alias(`f`),
		opt.ArgName(`fmt1,fmt2,..`),
		opt.Description(`One or multiple of: hex, dec, oct, bit`),
	)

	opt.Bool(`help`, false,
		opt.Alias("h", "?"),
		opt.Description(`Show this help`),
	)

	limitS := opt.IntOptional(`limit`, 0,
		opt.Alias("l"),
		opt.ArgName(`bytes`),
		opt.Description(`Read only N bytes (0 = no limit)`),
	)

	startOffsetS := opt.IntOptional(`seek`, 0,
		opt.Alias("s"),
		opt.ArgName(`offset`),
		opt.Description(`Start reading from certain offset`),
	)

	remaining, err := opt.Parse(os.Args[1:])

	if opt.Called("help") {
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`heksa - hex file dumper %v build %v`+"\n", VERSION, BUILD))
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`(c) %v 2019- [ %v ]`+"\n", AUTHOR, HOMEPAGE))
		fmt.Fprintf(os.Stdout, opt.Help())
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`EXAMPLES:`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -f hex,asc,bit foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -o hex -f hex,asc,bit foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -o hex -f bit foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -l 1024 foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -s 1234 foo.dat`)+"\n")
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		fmt.Fprintf(os.Stderr, opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}

	/*
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			fmt.Println("data is being piped to stdin")
		}
	*/

	offViewer, err := getOffsetViewer(*offsetDisplayS)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error getting offset displayer: %v`, err))
		os.Exit(1)
	}

	displays, err := getViewers(strings.Split(*formatS, `,`))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error getting data displayer: %v`, err))
		os.Exit(1)
	}

	if len(remaining) != 1 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error: no file given as argument`))
		os.Exit(1)
	}

	fpath := remaining[0]

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

	startOffset := int64(*startOffsetS)

	if startOffset != 0 {
		_, err = f.Seek(startOffset, io.SeekCurrent)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(`couldn't seek: %v`, err))
			os.Exit(1)
		}
	}

	bitWidth := uint8(bits.Len64(uint64(fi.Size())))
	bitWidth = (bitWidth + (8 - 1)) & ^(bitWidth - 1)

	offViewer.SetFileSize(fi.Size())

	for idx, _ := range displays {
		displays[idx].SetPalette(defaultCharacterColors)
	}

	limit := uint64(*limitS)

	r := New(f, offViewer, displays)

	// Dump hex
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

		if limit > 0 && r.ReadBytes >= limit {
			break
		}

	}

}
