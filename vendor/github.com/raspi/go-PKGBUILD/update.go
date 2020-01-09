package PKGBUILD

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type checksumType uint8

const (
	Sha1 checksumType = iota
	Sha224
	Sha256
	Sha384
	Sha512
	B2
	Md5
)

func (ct checksumType) String() string {
	switch ct {
	case Sha1:
		return `sha1`
	case Sha224:
		return `sha224`
	case Sha256:
		return `sha256`
	case Sha384:
		return `sha384`
	case Sha512:
		return `sha512`
	case B2:
		return `b2`
	case Md5:
		return `md5`
	default:
		return `?unknown?`
	}
}

const (
	ReplaceFromChecksumFilename = `<FNAMEARCH>`
)

// Update checksums to file(s)
// File must be in format
// <checksum> <file path>
// Filename in checksum file must be in format
// something-linux-<ARCH>.something
//
// String ReplaceFromChecksumFilename is replaced with architecture name from checksum filename's architecture
func GetChecksumsFromFile(chtype checksumType, path string, prefix string, suffix string) (f Files) {
	f = make(Files)

	sumsFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(bytes.NewReader(sumsFile))

	for sc.Scan() {
		line := sc.Text()
		if !strings.Contains(line, `linux`) {
			continue
		}

		checksumAndFilename := regexp.MustCompile(`([^\s]+)\s+([^\s]+)`)
		matches := checksumAndFilename.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		checksum := matches[1]
		fname := matches[2]

		cpuArchitecture := regexp.MustCompile(`linux-([^.]+)\.`)
		cpuArchFromFilename := cpuArchitecture.FindStringSubmatch(fname)

		if cpuArchFromFilename == nil {
			continue
		}

		goarch := cpuArchFromFilename[1]
		linuxarch, ok := GoArchToLinuxArch[goarch]
		if !ok {
			log.Fatalf(`architecture %v not found`, goarch)
		}

		newSuffix := strings.ReplaceAll(suffix, ReplaceFromChecksumFilename, goarch)

		f[linuxarch] = append(f[linuxarch],
			Source{
				URL: fmt.Sprintf(`%s%s`, prefix, newSuffix),
				Checksums: map[string]string{
					chtype.String(): checksum,
				},
			},
		)
	}

	return f
}
