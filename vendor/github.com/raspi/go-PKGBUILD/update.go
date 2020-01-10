package PKGBUILD

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path"
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

// Update checksums to file(s)
// File must be in format
// <checksum> <file path>
func GetChecksumsFromFile(chtype checksumType, path string, fn func(fpath string) (url, arch, alias string)) (f Files) {
	f = make(Files)
	lines, err := GetLinesFromFile(path)

	if err != nil {
		panic(err)
	}

	for _, line := range lines {
		checksumAndFilename := regexp.MustCompile(`^([^\s]+)\s+([^\s]+)$`)
		matches := checksumAndFilename.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		checksum := matches[1]
		fname := matches[2]

		url, arch, alias := fn(fname)

		if url == `` {
			continue
		}

		newSource := Source{
			URL: url,
			Checksums: map[string]string{
				chtype.String(): checksum,
			},
		}

		if alias != `` {
			newSource.Alias = alias
		}

		f[arch] = append(f[arch], newSource)
	}

	return f
}

// Read a file and split with new line separator
func GetLinesFromFile(path string) (lines []string, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(bytes.NewReader(b))

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines, nil
}

// How architecture can be found from file name
var DefaultArchRegEx = regexp.MustCompile(`linux-([^\.]+)\.`)

func (t Template) DefaultChecksumFilesFunc(fpath string) (url, arch, alias string) {
	fpath = strings.TrimLeft(fpath, `.`)
	fpath = strings.TrimLeft(fpath, `/`)
	filename := path.Base(fpath)

	if !strings.Contains(filename, `linux`) {
		return ``, ``, ``
	}

	match := DefaultArchRegEx.FindStringSubmatch(filename)
	if len(match) == 0 {
		return ``, ``, ``
	}

	filename = strings.Replace(filename, t.Name[0], `$pkgname`, 1)
	filename = strings.Replace(filename, t.Version, `$pkgver`, 1)

	arch, ok := GoArchToLinuxArch[match[1]]
	if !ok {
		return ``, ``, ``
	}

	url = path.Join(t.PackageURLPrefix, filename)
	return url, arch, alias
}
