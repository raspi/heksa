package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/DavidGamba/go-getoptions"
	"github.com/raspi/heksa/pkg/color"
	"github.com/raspi/heksa/pkg/iface"
	"github.com/raspi/heksa/pkg/reader"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
	offFormatters "github.com/raspi/heksa/pkg/reader/offsetFormatters/base"
	"github.com/raspi/heksa/pkg/units"
)

var (
	// These are set with Makefile -X=main.VERSION, etc
	VERSION   = `v0.0.0`
	BUILD     = `dev`
	BUILDDATE = `0000-00-00T00:00:00+00:00`
)

const (
	AUTHOR   = `Pekka Järvinen`
	HOMEPAGE = `https://github.com/raspi/heksa`
)

// These color group names MUST exist in config
var requiredColorGroupNames = []string{
	`LineEven`, `LineOdd`, `Splitter`, `Offset`, `Padding`, `Default`, `Special`, `Highlight`,
}

// Parse command line arguments
func getParams() (source iface.ReadSeekerCloser, offsetViewer []reader.OffsetFormatter, colorGroupings map[string]string, limit uint64, filesize int64, fg base.FormatterGroup, printRelative bool) {
	opt := getoptions.New()

	opt.HelpSynopsisArgs(`<filename> or STDIN`)

	opt.Bool(`help`, false,
		opt.Alias("h", "?"),
		opt.Description(`Show this help`),
	)

	opt.Bool(`version`, false,
		opt.Description(`Show version information`),
	)

	argOffset := opt.StringOptional(`offset-format`, `hex`,
		opt.Alias(`o`),
		opt.ArgName(`fmt1[,fmt2]`),
		opt.Description(
			`One or two of: `+strings.Join(reader.GetOffsetViewerList(), `, `)+`, no, ''.`+
				"\n"+
				`First one is displayed on the left side and second one on right side after formatters.`,
		),
	)

	argPrintRelativeOffset := opt.Bool(`print-relative-offset`, false,
		opt.Alias(`r`),
		opt.Description(`Print relative offset(s) starting from 0 (file only)`),
	)

	argFormat := opt.StringOptional(`format`, `hex,asc`,
		opt.Alias(`f`),
		opt.ArgName(`fmt1,fmt2,..`),
		opt.Description(`One or multiple of: `+strings.Join(reader.GetViewerList(), `, `)),
	)

	argLimit := opt.StringOptional(`limit`, `0`,
		opt.Alias("l"),
		opt.ArgName(`[prefix]bytes[unit]`),
		opt.Description(`Read only N bytes (0 = no limit). See NOTES.`),
	)

	argSeek := opt.StringOptional(`seek`, `0`,
		opt.Alias("s"),
		opt.ArgName(`[prefix]offset[unit]`),
		opt.Description(`Start reading from certain offset. See NOTES.`),
	)

	argWidth := opt.StringOptional(`width`, `16`,
		opt.Alias("w"),
		opt.ArgName(`[prefix]width`),
		opt.Description(`Width. See NOTES.`),
	)

	argSplitter := opt.IntOptional(`splitter`, 8,
		opt.Alias("S"),
		opt.ArgName(`size`),
		opt.Description(`Insert visual splitter every N bytes. Zero (0) disables.`),
	)

	remainingArgs, err := opt.Parse(os.Args[1:])

	if opt.Called("help") {
		_, _ = fmt.Fprintf(os.Stdout, `heksa - hex file dumper %v - (%v)`+"\n", VERSION, BUILDDATE)
		_, _ = fmt.Fprintf(os.Stdout, `(c) %v 2019- [ %v ]`+"\n", AUTHOR, HOMEPAGE)
		_, _ = fmt.Fprintln(os.Stdout, opt.Help())
		_, _ = fmt.Fprintln(os.Stdout, `NOTES:`)
		_, _ = fmt.Fprintln(os.Stdout, `    - You can use prefixes for seek, limit and width. 0x = hex, 0b = binary, 0o = octal`)
		_, _ = fmt.Fprintln(os.Stdout, `    - Use '--seek \-1234' for seeking from end of file`)
		_, _ = fmt.Fprintln(os.Stdout, `    - Limit and seek parameters supports units (KB, KiB, MB, MiB, GB, GiB, TB, TiB)`)
		_, _ = fmt.Fprintln(os.Stdout, `    - --print-relative-offset can be used when seeking to certain offset to also print extra offset position starting from zero`)
		_, _ = fmt.Fprintln(os.Stdout, `    - Offset formatters:`)
		_, _ = fmt.Fprintln(os.Stdout, `      - Disable formatter output with 'no' or ''`)
		_, _ = fmt.Fprintln(os.Stdout, `      - 'humiec' (IEC: 1024 B) and 'humsi' (SI: 1000 B) displays offset in human form (n KiB/KB)`)
		_, _ = fmt.Fprintln(os.Stdout, `    - Formatters:`)
		_, _ = fmt.Fprintln(os.Stdout, `      - 'blk' can be used to print simple color blocks which helps to visualize where data vs. human readable strings are`)
		_, _ = fmt.Fprintln(os.Stdout)
		_, _ = fmt.Fprintln(os.Stdout, `EXAMPLES:`)
		_, _ = fmt.Fprintln(os.Stdout, `    heksa -f hex,asc,bit foo.dat`)
		_, _ = fmt.Fprintln(os.Stdout, `    heksa -o hex,per -f hex,asc foo.dat`)
		_, _ = fmt.Fprintln(os.Stdout, `    heksa -o hex -f hex,asc,bit foo.dat`)
		_, _ = fmt.Fprintln(os.Stdout, `    heksa -o no -f bit foo.dat`)
		_, _ = fmt.Fprintln(os.Stdout, `    heksa -l 0x1024 foo.dat`)
		_, _ = fmt.Fprintln(os.Stdout, `    heksa -s 0b1010 foo.dat`)
		_, _ = fmt.Fprintln(os.Stdout, `    heksa -s 4321KiB foo.dat`)
		_, _ = fmt.Fprintln(os.Stdout, `    heksa -w 8 foo.dat`)
		_, _ = fmt.Fprintln(os.Stdout, `    echo "test" | heksa`)
		os.Exit(0)
	} else if opt.Called("version") {
		_, _ = fmt.Fprintf(os.Stdout, `%v build %v on %v`+"\n", VERSION, BUILD, BUILDDATE)
		os.Exit(0)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		_, _ = fmt.Fprintln(os.Stderr, opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}

	limitTmp, err := units.Parse(*argLimit)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error parsing limit: %v`, err)
		os.Exit(1)
	}
	limit = uint64(limitTmp)

	startOffset, err := units.Parse(strings.Replace(*argSeek, `\`, ``, -1))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error parsing seek: %v`, err)
		os.Exit(1)
	}

	widthTmp, err := units.Parse(*argWidth)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error parsing width: %v`, err)
		os.Exit(1)
	}
	width := uint16(widthTmp)
	if width == 0 {
		_, _ = fmt.Fprint(os.Stderr, `width must be > 0`)
		os.Exit(1)
	}

	offsetViewer, err = reader.GetOffsetFormatters(strings.Split(*argOffset, `,`))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error getting offset formatter: %v`, err)
		os.Exit(1)
	}

	displays, err := reader.GetViewers(strings.Split(*argFormat, `,`))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error getting formatter: %v`, err)
		os.Exit(1)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `couldn't stat stdin: %v`, err)
		os.Exit(1)
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Stdin has data
		source = os.Stdin
		filesize = -1
	} else {
		// Read file
		if len(remainingArgs) != 1 {
			_, _ = fmt.Fprintln(os.Stderr, `error: no file given as argument, see --help`)
			os.Exit(1)
		}

		fpath := remainingArgs[0]

		fhandle, err := os.Open(fpath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, `error opening file: %v`, err)
			os.Exit(1)
		}

		fi, err := fhandle.Stat()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, `error stat'ing file: %v`, err)
			os.Exit(1)
		}

		if fi.IsDir() {
			_, _ = fmt.Fprintf(os.Stderr, `error: %v is directory`, fpath)
			os.Exit(1)
		}

		// Seek to given offset
		if startOffset > 0 {
			_, err = fhandle.Seek(startOffset, io.SeekCurrent)
		} else if startOffset < 0 {
			_, err = fhandle.Seek(startOffset, io.SeekEnd)
		}

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, `couldn't seek to %v: %v`, startOffset, err)
			os.Exit(1)
		}

		filesize = fi.Size()

		if !fi.Mode().IsRegular() {
			// Not a regular file, so file size is unknown
			filesize = -1
		}

		source = fhandle
	}

	colorGroupings, err = color.GetColorGroupColorDefaults(strings.NewReader(DefaultGroupColors), requiredColorGroupNames)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error loading color group config: %v`, err)
		os.Exit(1)
	}

	palette, err := color.GetColors(strings.NewReader(DefaultByteColorGroups), colorGroupings)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error loading color config: %v`, err)
		os.Exit(1)
	}

	var formatters []base.ByteFormatter
	for _, f := range displays {
		fmter := reader.GetByteFormatter(f, colorGroupings[`Highlight`], colorGroupings[`Special`])
		if fmter == nil {
			_, _ = fmt.Fprintf(os.Stderr, `error: unknown formatter %v`, f)
			os.Exit(1)

		}

		formatters = append(formatters, fmter)
	}

	fGroup := base.New(formatters, palette, colorGroupings[`Splitter`], colorGroupings[`Padding`], width, uint8(*argSplitter))

	return source, offsetViewer, colorGroupings, limit, filesize, fGroup, *argPrintRelativeOffset
}

func main() {
	source, offViewer, colorGroupings, limit, filesize, fGroup, printRelative := getParams()
	usingLimit := limit > 0

	binfo := offFormatters.BaseInfo{
		FileSize: filesize,
	}

	var offormatters []offFormatters.OffsetFormatter
	for _, f := range offViewer {
		fmter := reader.GetFromOffsetFormatter(f, binfo)
		offormatters = append(offormatters, fmter)
	}

	colors := reader.ReaderColors{
		LineOdd:  colorGroupings[`LineOdd`],
		LineEven: colorGroupings[`LineEven`],
		Offset:   colorGroupings[`Offset`],
		Splitter: colorGroupings[`Splitter`],
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	isStdin := filesize == -1
	if isStdin {
		printRelative = false
	}

	r := reader.New(source, offormatters, colors, fGroup, isStdin, printRelative)

	// Dump hex
	for {
		select {
		case <-stop: // Kill or ctrl-C
			break
		default:
		}

		s, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(`error while reading file: %v`, err))
			os.Exit(1)
		}

		// Print formatted line
		// <optional offset formatter #1><split><format 1><split><format N...><optional split><optional offset formatter #2>
		_, _ = fmt.Println(s)

		if usingLimit && r.GetReadBytes() >= limit {
			// Limit is set and found
			break
		}
	}

	err := source.Close()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `couldn't close file: %v`, err)
		os.Exit(1)
	}
}
