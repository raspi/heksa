package main

// Generate Arch Linux PKGBUILD file from package.json

import (
	"flag"
	"fmt"
	"github.com/raspi/go-PKGBUILD"
	"io/ioutil"
	"log"
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
		basetpl.DefaultChecksumFilesFunc,
	)

	basetpl.Commands.Install, err = PKGBUILD.GetLinesFromFile(`install.sh`)
	if err != nil {
		log.Fatal(err)
	}

	verrs := basetpl.Validate()
	if len(verrs) > 0 {
		for _, E := range verrs {
			log.Print(E)
		}

		log.Fatal(`error(s) found. aborting.`)
	}

	fmt.Println(basetpl)
}
