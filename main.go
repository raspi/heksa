package main

import (
	"fmt"
	"github.com/DavidGamba/go-getoptions"
	"github.com/raspi/heksa/pkg/iface"
	"github.com/raspi/heksa/pkg/reader"
	"io"
	"os"
	"strconv"
	"strings"
)

var VERSION = `v0.0.0`
var BUILD = `dev`
var BUILDDATE = `0000-00-00T00:00:00+00:00`

const AUTHOR = `Pekka Järvinen`
const HOMEPAGE = `https://github.com/raspi/heksa`

func getParams() (source iface.ReadSeekerCloser, displays []iface.CharacterFormatter, offsetViewer []iface.OffsetFormatter, limit uint64, startOffset int64) {
	opt := getoptions.New()

	opt.HelpSynopsisArgs(`<filename>`)

	offsetDisplayS := opt.StringOptional(`offset-format`, `hex`,
		opt.Alias(`o`),
		opt.ArgName(`[fmt1][,fmt2]`),
		opt.Description(`Zero to two of: hex, dec, oct, per. First one is displayed on the left side and second one on right after formatters`),
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

	opt.Bool(`version`, false,
		opt.Description(`Show version information`),
	)

	limitS := opt.StringOptional(`limit`, `0`,
		opt.Alias("l"),
		opt.ArgName(`[prefix]bytes`),
		opt.Description(`Read only N bytes (0 = no limit). See NOTES.`),
	)

	startOffsetS := opt.StringOptional(`seek`, `0`,
		opt.Alias("s"),
		opt.ArgName(`[prefix]offset`),
		opt.Description(`Start reading from certain offset. See NOTES.`),
	)

	remaining, err := opt.Parse(os.Args[1:])

	if opt.Called("help") {
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`heksa - hex file dumper %v - (%v)`+"\n", VERSION, BUILDDATE))
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`(c) %v 2019- [ %v ]`+"\n", AUTHOR, HOMEPAGE))
		fmt.Fprintf(os.Stdout, opt.Help())
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`NOTES:`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    You can use prefixes for seek and limit. 0x = hex, 0b = binary, 0o = octal`)+"\n")
		fmt.Fprintf(os.Stdout, "\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`EXAMPLES:`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -f hex,asc,bit foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -o hex,per -f hex,asc foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -o hex -f hex,asc,bit foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -o '' -f bit foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -l 0x1024 foo.dat`)+"\n")
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`    heksa -s 0b1010 foo.dat`)+"\n")
		os.Exit(0)
	} else if opt.Called("version") {
		fmt.Fprintf(os.Stdout, fmt.Sprintf(`%v build %v on %v`+"\n", VERSION, BUILD, BUILDDATE))
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		fmt.Fprintf(os.Stderr, opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}

	offsetViewer, err = reader.GetOffsetViewer(strings.Split(*offsetDisplayS, `,`))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error getting offset displayer: %v`, err))
		os.Exit(1)
	}

	displays, err = reader.GetViewers(strings.Split(*formatS, `,`))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error getting data displayer: %v`, err))
		os.Exit(1)
	}

	limit, err = strconv.ParseUint(*limitS, 0, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error parsing limit: %v`, err))
		os.Exit(1)
	}

	startOffset, err = strconv.ParseInt(*startOffsetS, 0, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`error parsing seek: %v`, err))
		os.Exit(1)
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Stdin has data
		source = os.Stdin

		// No clue of file size when streaming from stdin
		for idx, _ := range offsetViewer {
			offsetViewer[idx].SetFileSize(0)
		}
	} else {
		// Read file
		if len(remaining) != 1 {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(`error: no file given as argument, see --help`))
			os.Exit(1)
		}

		fpath := remaining[0]

		fhandle, err := os.Open(fpath)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(`error opening file: %v`, err))
			os.Exit(1)
		}

		fi, err := fhandle.Stat()
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(`error stat'ing file: %v`, err))
			os.Exit(1)
		}

		if fi.IsDir() {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(`error: %v is directory`, fpath))
			os.Exit(1)
		}

		// Hint offset viewer
		for idx, _ := range offsetViewer {
			offsetViewer[idx].SetFileSize(fi.Size())
		}

		source = fhandle

	}

	return source, displays, offsetViewer, limit, startOffset
}

func main() {
	source, displays, offViewer, limit, startOffset := getParams()
	palette := defaultCharacterColors

	for i := uint8(0); i < 255; i++ {
		_, ok := palette[i]
		if !ok {
			// Fall back
			palette[i] = defaultColor
		}
	}

	if startOffset != 0 {
		_, err := source.Seek(startOffset, io.SeekCurrent)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(`couldn't seek: %v`, err))
			os.Exit(1)
		}
	}

	r := reader.New(source, offViewer, displays, palette)

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
			// Limit is set and found
			break
		}

	}

	source.Close()

}
