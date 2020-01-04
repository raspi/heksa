package main

/*
Helper for creating Arch Linux PKGBUILD packages
See Makefile in the root directory
*/

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

var masterTemplate = TemplateRaw{
	Maintainer:     `Pekka Järvinen`,
	PkgName:        `heksa`,
	Version:        `v0.0.0`, // Dynamic
	PkgDescription: `hex dumper with colors`,
	PkgArch:        "???", // Dynamic
	PkgLicense:     "Apache 2.0",
	Url:            `https://github.com/raspi/heksa`,
	Source:         "https://github.com/raspi/heksa/releases/download/$pkgver/$pkgname-$pkgver-linux-<ARCH>.tar.gz", // Dynamic
	PrepareSteps:   []string{},
	BuildSteps:     []string{},
	CheckSteps:     []string{},

	// Install instructions
	PackageSteps: []string{
		`cd "$srcdir"`,
		`install -Dm644 "LICENSE" -t "$pkgdir/usr/share/licenses/$pkgname"`,
		`install -Dm644 "README.md" -t "$pkgdir/usr/share/doc/$pkgname"`,
		`install -Dm755 "bin/$pkgname" -t "$pkgdir/usr/bin"`,
	},
}

type TemplateRaw struct {
	Maintainer      string   // Package maintainer
	PkgName         string   // Package name (app name)
	Version         string   // App version
	PkgDescription  string   // App description
	PkgArch         string   // CPU arch
	PkgLicense      string   // License
	Url             string   // Homepage URL
	Source          string   // Source URL
	PrepareSteps    []string // prepare(){}
	BuildSteps      []string // build(){}
	CheckSteps      []string // check(){}
	PackageSteps    []string // package(){}
	Sha256Checksums []string // SHA256 checksums
}

type Template struct {
	Maintainer      string // Package maintainer
	PkgName         string // Package name (app name)
	Version         string // App version
	PkgDescription  string // App description
	PkgArch         string // CPU arch
	PkgLicense      string // License
	Url             string // Homepage URL
	Source          string // Source URL
	PrepareSteps    string // prepare(){}
	BuildSteps      string // build(){}
	CheckSteps      string // check(){}
	PackageSteps    string // package(){}
	Sha256Checksums string // SHA256 checksums
}

// Rewrite Go's arch names to Linux ones
var GoArchToLinuxArch = map[string]string{
	`amd64`:   `x86_64`,
	`arm`:     `arm`,
	`arm64`:   `aarch64`,
	`ppc64`:   `ppc64`,
	`ppc64le`: `ppc64le`,
}

func (rt TemplateRaw) ToTemplate() Template {
	sha256checksums := fmt.Sprintf(`%q`, rt.Sha256Checksums)
	sha256checksums = strings.ReplaceAll(sha256checksums, `[`, ``)
	sha256checksums = strings.ReplaceAll(sha256checksums, `]`, ``)

	prepare := ``
	if len(rt.PrepareSteps) > 0 {
		prepare += `prepare() {` + "\n"
		prepare += `  `
		prepare += strings.Join(rt.PrepareSteps, "\n  ") + "\n"
		prepare += `}` + "\n"
	}

	build := ``
	if len(rt.BuildSteps) > 0 {
		build += `build() {` + "\n"
		build += `  `
		build += strings.Join(rt.BuildSteps, "\n  ") + "\n"
		build += `}` + "\n"
	}

	check := ``
	if len(rt.CheckSteps) > 0 {
		check += `check() {` + "\n"
		check += `  `
		check += strings.Join(rt.CheckSteps, "\n  ") + "\n"
		check += `}` + "\n"
	}

	pack := ``
	if len(rt.PackageSteps) > 0 {
		pack += `package() {` + "\n"
		pack += `  `
		pack += strings.Join(rt.PackageSteps, "\n  ") + "\n"
		pack += `}` + "\n"
	}

	return Template{
		Maintainer:      rt.Maintainer,
		PkgName:         rt.PkgName,
		Version:         rt.Version,
		PkgDescription:  rt.PkgDescription,
		PkgArch:         rt.PkgArch,
		PkgLicense:      rt.PkgLicense,
		Url:             rt.Url,
		Source:          rt.Source,
		PrepareSteps:    prepare,
		BuildSteps:      build,
		CheckSteps:      check,
		PackageSteps:    pack,
		Sha256Checksums: sha256checksums,
	}
}

/* Generate PKGBUILD files to release/<version> directory */
func main() {

	versionArg := flag.String("version", ``, `version (v1.2.3)`)

	flag.Parse()

	if *versionArg == `` {
		log.Fatal(`no version given`)
	}

	releasePath := path.Join(`..`, `..`, *versionArg)

	files, err := ioutil.ReadDir(releasePath)
	if err != nil {
		log.Fatal(err)
	}

	sumsFile, err := ioutil.ReadFile(path.Join(releasePath, fmt.Sprintf(`%s-%s.shasums`, masterTemplate.PkgName, *versionArg)))
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(bytes.NewReader(sumsFile))

	var checksums = make(map[string]string)

	for sc.Scan() {
		line := sc.Text()
		if !strings.Contains(line, `linux`) {
			continue
		}

		if !strings.Contains(line, *versionArg) {
			continue
		}

		if !strings.Contains(line, masterTemplate.PkgName) {
			continue
		}

		mre := regexp.MustCompile(`([^\s]*)\s*([^\s]*)`)
		matches := mre.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		checksum := matches[1]
		fname := matches[2]

		fmre := regexp.MustCompile(`linux-([^.]*)\.`)
		fmatches := fmre.FindStringSubmatch(fname)

		if fmatches == nil {
			continue
		}

		goarch := fmatches[1]

		checksums[goarch] = checksum
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), `linux`) {
			continue
		}

		mre := regexp.MustCompile(`linux-([^.]*)\.`)
		matches := mre.FindStringSubmatch(file.Name())

		if matches == nil {
			continue
		}

		goarch := matches[1]

		tplcopy := masterTemplate

		arch, ok := GoArchToLinuxArch[goarch]

		if !ok {
			continue
		}

		tplcopy.Source = strings.ReplaceAll(tplcopy.Source, `<ARCH>`, goarch)
		tplcopy.PkgArch = arch
		tplcopy.Sha256Checksums = append(tplcopy.Sha256Checksums, checksums[goarch])
		tplcopy.Version = *versionArg

		var tmpw bytes.Buffer
		tpl, err := template.New(``).ParseFiles(`./PKGBUILD.txt`)
		if err != nil {
			panic(err)
		}

		err = tpl.ExecuteTemplate(&tmpw, `tpl`, tplcopy.ToTemplate())
		if err != nil {
			panic(err)
		}

		pkgfile, err := os.Create(path.Join(releasePath, fmt.Sprintf(`%v-%v-linux-Arch-%s.PKGBUILD`, tplcopy.PkgName, tplcopy.Version, arch)))
		if err != nil {
			log.Fatal(err)
		}

		pkgfile.Write(tmpw.Bytes())
		pkgfile.Close()

		log.Printf(`wrote %v`, pkgfile.Name())
	}
}
