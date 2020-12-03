package main

// Generate Arch Linux PKGBUILD file from package.json

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/raspi/go-PKGBUILD"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	versionArg := flag.String("version", ``, `version (v1.2.3)`)
	flag.Parse()

	if *versionArg == `` {
		fmt.Fprintf(os.Stderr, `no version given`)
		os.Exit(1)
	}

	releasePath := path.Join(`..`, `..`, *versionArg)

	log.Printf(`Reading package.json`)
	jb, err := ioutil.ReadFile(`package.json`)
	if err != nil {
		fmt.Fprintf(os.Stderr, `couldn't load package.json: %v`, err)
		os.Exit(1)
	}

	basetpl, err := PKGBUILD.FromJson(jb)
	if err != nil {
		fmt.Fprintf(os.Stderr, `couldn't build template from JSON: %v`, err)
		os.Exit(1)
	}

	packageName := basetpl.Name[0]
	basetpl.Version = *versionArg

	flist := make(PKGBUILD.Files)
	shasumsfile := fmt.Sprintf(`%s-%s.shasums`, packageName, basetpl.Version)

	log.Printf(`Reading checksums`)
	cf, err := os.Open(path.Join(releasePath, shasumsfile))
	if err != nil {
		fmt.Fprintf(os.Stderr, `couldn't open shasums file %q: %v`, shasumsfile, err)
		os.Exit(1)
	}
	defer cf.Close()

	scanner := bufio.NewScanner(cf)

	for scanner.Scan() {
		l := strings.Split(scanner.Text(), "  ")
		csum, fname := l[0], path.Clean(l[1])

		if !strings.Contains(fname, `linux`) {
			continue
		}

		match := PKGBUILD.DefaultArchRegEx.FindStringSubmatch(fname)
		if len(match) == 0 {
			continue
		}

		arch, ok := PKGBUILD.GoArchToLinuxArch[match[1]]
		if !ok {
			continue
		}

		u, err := url.Parse(basetpl.PackageURLPrefix)
		if err != nil {
			continue
		}

		if strings.HasPrefix(fname, packageName) {
			fname = strings.Replace(fname, packageName, `$pkgname`, 1)
		}

		if strings.Contains(fname, basetpl.Version) {
			fname = strings.Replace(fname, basetpl.Version, `$pkgver`, 1)
		}

		u.Path = path.Join(u.Path, fname)

		flist[arch] = append(flist[arch], PKGBUILD.Source{
			URL:   u.String(),
			Alias: ``,
			Checksums: map[string]string{
				PKGBUILD.Sha256.String(): csum,
			},
		})
	}

	basetpl.Files = flist

	basetpl.Commands.Install, err = PKGBUILD.GetLinesFromFile(`install.sh`)
	if err != nil {
		log.Fatal(err)
	}

	verrs := basetpl.Validate()
	if len(verrs) > 0 {
		fmt.Fprintf(os.Stderr, `error:`)
		for _, E := range verrs {
			fmt.Fprintf(os.Stderr, `- %v`, E)
		}

		fmt.Fprintf(os.Stderr, `invalid template`)
		os.Exit(1)
	}

	pf, err := ioutil.TempFile(`.`, `tmp-PKGBUILD-`)
	if err != nil {
		fmt.Fprintf(os.Stderr, `error: %v`, err)
		os.Exit(1)
	}
	defer pf.Close()

	log.Printf(`Writing temp PKGBUILD file %#v`, pf.Name())
	pf.WriteString(basetpl.String())

	fpath, err := filepath.Abs(path.Join(`..`, `..`, basetpl.Version, fmt.Sprintf(`%s-%s-linux-Arch.PKGBUILD`, packageName, basetpl.Version)))
	if err != nil {
		fmt.Fprintf(os.Stderr, `error: %v`, err)
		os.Exit(1)
	}

	err = os.Rename(pf.Name(), fpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, `error: %v`, err)
		os.Exit(1)
	}

	log.Printf(`Renamed %#v -> %#v`, pf.Name(), fpath)
	log.Print(`PKGBUILD Done.`)
}
