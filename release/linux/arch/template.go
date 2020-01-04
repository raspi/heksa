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
	"sort"
	"strings"
	"text/template"
)

var masterTemplate = TemplateRaw{
	Maintainer:     `Pekka Järvinen`,
	PkgName:        `heksa`,
	Version:        `v0.0.0`, // Dynamic
	PkgDescription: `hex dumper with colors`,
	PkgLicense:     "Apache 2.0",
	Url:            `https://github.com/raspi/heksa`,
	SourceDyn:      "https://github.com/raspi/heksa/releases/download/$pkgver/$pkgname-$pkgver-linux-<ARCH>.tar.gz", // Dynamic
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
	Maintainer     string // Package maintainer
	PkgName        string // Package name (app name)
	Version        string // App version
	PkgDescription string // App description
	PkgLicense     string // License
	Url            string // Homepage URL
	SourceDyn      string
	Source         Sources  // Source URL
	PrepareSteps   []string // prepare(){}
	BuildSteps     []string // build(){}
	CheckSteps     []string // check(){}
	PackageSteps   []string // package(){}
}

type Template struct {
	Maintainer     string // Package maintainer
	PkgName        string // Package name (app name)
	Version        string // App version
	PkgDescription string // App description
	PkgArch        string // CPU arch
	PkgLicense     string // License
	Url            string // Homepage URL
	Source         string // Source URL
	PrepareSteps   string // prepare(){}
	BuildSteps     string // build(){}
	CheckSteps     string // check(){}
	PackageSteps   string // package(){}
}

// Rewrite Go's arch names to Linux ones
var GoArchToLinuxArch = map[string]string{
	`amd64`:   `x86_64`,
	`arm`:     `arm`,
	`arm64`:   `aarch64`,
	`ppc64`:   `ppc64`,
	`ppc64le`: `ppc64le`,
}

type Source struct {
	Url      string
	Checksum string
	File     string
}

type Sources = map[string]Source

func (rt TemplateRaw) ToTemplate() Template {
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

	source := ``

	archList := []string{}

	for k, v := range rt.Source {
		log.Printf(`adding source arch %v %v %v`, k, v.Checksum, v.File)
		archList = append(archList, k)
		source += fmt.Sprintf(`source_%v=("`, k)
		source += v.Url
		source += `")` + "\n"
		source += fmt.Sprintf(`sha256sums_%v=('`, k)
		source += v.Checksum
		source += `')` + "\n"

	}

	sort.Strings(archList)

	archStrlist := fmt.Sprintf(`%q`, archList)
	archStrlist = strings.ReplaceAll(archStrlist, `[`, ``)
	archStrlist = strings.ReplaceAll(archStrlist, `]`, ``)
	archStrlist = strings.ReplaceAll(archStrlist, `"`, `'`)

	return Template{
		Maintainer:     rt.Maintainer,
		PkgName:        rt.PkgName,
		Version:        rt.Version,
		PkgDescription: rt.PkgDescription,
		PkgArch:        archStrlist,
		PkgLicense:     rt.PkgLicense,
		Url:            rt.Url,
		Source:         source,
		PrepareSteps:   prepare,
		BuildSteps:     build,
		CheckSteps:     check,
		PackageSteps:   pack,
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

	sumsFile, err := ioutil.ReadFile(path.Join(releasePath, fmt.Sprintf(`%s-%s.shasums`, masterTemplate.PkgName, *versionArg)))
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(bytes.NewReader(sumsFile))

	sources := make(map[string]Source)

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
		linuxarch, ok := GoArchToLinuxArch[goarch]
		if !ok {
			log.Fatalf(`architecture %v not found`, goarch)
		}

		sources[linuxarch] = Source{
			Checksum: checksum,
			File:     fname,
			Url:      strings.ReplaceAll(masterTemplate.SourceDyn, `<ARCH>`, goarch),
		}

	}

	masterTemplate.Source = sources
	masterTemplate.Version = *versionArg

	var tmpw bytes.Buffer
	tpl, err := template.New(``).ParseFiles(`./PKGBUILD.txt`)
	if err != nil {
		panic(err)
	}

	err = tpl.ExecuteTemplate(&tmpw, `tpl`, masterTemplate.ToTemplate())
	if err != nil {
		panic(err)
	}

	pkgfile, err := os.Create(path.Join(releasePath, fmt.Sprintf(`%v-%v-linux-Arch.PKGBUILD`, masterTemplate.PkgName, masterTemplate.Version)))
	if err != nil {
		log.Fatal(err)
	}

	pkgfile.Write(tmpw.Bytes())
	pkgfile.Close()

	log.Printf(`wrote %v`, pkgfile.Name())
}
