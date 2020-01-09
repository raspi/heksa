package main

import (
	"flag"
	"fmt"
	"github.com/raspi/go-PKGBUILD"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	versionArg := flag.String("version", ``, `version (v1.2.3)`)
	flag.Parse()

	if *versionArg == `` {
		log.Fatal(`no version given`)
	}

	releasePath := path.Join(`..`, `..`, *versionArg)

	jb, err := ioutil.ReadFile(`package.json`)
	if err != nil {
		log.Fatal(err)
	}

	basetpl, err := PKGBUILD.FromJson(jb)
	if err != nil {
		log.Fatal(err)
	}

	packageName := basetpl.Name[0]
	basetpl.Version = *versionArg
	basetpl.Files = PKGBUILD.GetChecksumsFromFile(
		PKGBUILD.Sha256,
		path.Join(releasePath, fmt.Sprintf(`%s-%s.shasums`, packageName, basetpl.Version)),
		`https://github.com/raspi/heksa/releases/download/$pkgver/$pkgname-$pkgver-linux-`,
		PKGBUILD.ReplaceFromChecksumFilename+`.tar.gz`,
	)

	pkgfile, err := os.Create(path.Join(releasePath, fmt.Sprintf(`%v-%v-linux-Arch.PKGBUILD`, packageName, basetpl.Version)))
	if err != nil {
		log.Fatal(err)
	}

	pkgfile.WriteString(basetpl.String())
	pkgfile.Close()

	log.Printf(`wrote %v`, pkgfile.Name())

}
