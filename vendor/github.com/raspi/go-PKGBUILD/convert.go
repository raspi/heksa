package PKGBUILD

// Convert template to PKGBUILD file

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

func (t Template) getArchList() (l []string) {
	for arch, _ := range t.Files {
		l = append(l, arch)
	}

	sort.Strings(l)
	return l
}

func (t Template) sourceToMap() (sources map[string][]string) {
	sources = make(map[string][]string)

	for _, arch := range t.getArchList() {
		for _, src := range t.Files[arch] {
			key := `source`
			if arch != `` {
				key = fmt.Sprintf(`%s_%s`, key, arch)
			}

			dest := src.URL
			if src.Alias != `` {
				dest = fmt.Sprintf(`%s::%s`, src.Alias, dest)
			}

			sources[key] = append(sources[key], dest)
		}

		for _, src := range t.Files[arch] {
			var checksumtypes []string
			for ctype := range src.Checksums {
				checksumtypes = append(checksumtypes, ctype)
			}

			for _, ctype := range checksumtypes {
				key := fmt.Sprintf(`%ssums_%s`, ctype, arch)
				sources[key] = append(sources[key], src.Checksums[ctype])
			}
		}
	}

	return sources
}

func (t Template) sourceToString() string {
	srcs := t.sourceToMap()
	var sorted []string
	for k, _ := range srcs {
		sorted = append(sorted, k)
	}

	sort.Strings(sorted)

	var arr []string
	for _, v := range sorted {
		presuf := `'`
		if strings.HasPrefix(v, `source`) {
			presuf = `"`
		}
		m := srcs[v]
		arr = append(arr, fmt.Sprintf(`%s=(%s)`, v, wrapStrings(m, ` `, presuf, presuf)))
	}

	return strings.Join(arr, "\n")
}

func (t Template) getDependsArchList() (l []string) {
	for arch, _ := range t.Dependencies {
		l = append(l, arch)
	}
	return l
}

func (t Template) getOptPackageArchList() (l []string) {
	for arch, _ := range t.OptionalPackages {
		l = append(l, arch)
	}
	return l
}

// Optional packages
func (t Template) getOptionalPackages() (m map[string][]string) {
	m = make(map[string][]string)
	for _, arch := range t.getOptPackageArchList() {
		for _, opt := range t.OptionalPackages[arch] {
			key := `optdepends`
			if arch != `` {
				key = fmt.Sprintf(`%v_%v`, key, arch)
			}
			m[key] = append(m[key], fmt.Sprintf(`%s: %s`, opt.Package, opt.Reason))
		}
	}

	return m
}

func (t Template) optionalToString() (out string) {
	var arr []string
	for arch, opt := range t.getOptionalPackages() {
		arr = append(arr, fmt.Sprintf(`%s=(%s)`, arch, wrapStrings(opt, ` `, `'`, `'`)))
	}
	return strings.Join(arr, "\n")
}

// Dependencies needed
func (t Template) getDepends() (m map[string][]string) {
	m = make(map[string][]string)
	for _, arch := range t.getDependsArchList() {
		// Packages needed for running
		key := `depends`
		if arch != `` {
			key = fmt.Sprintf(`%s_%s`, key, arch)
		}

		m[key] = t.Dependencies[arch].Packages

		// Needed to make package from source
		key = `makedepends`
		if arch != `` {
			key = fmt.Sprintf(`%s_%s`, key, arch)
		}

		m[key] = t.Dependencies[arch].BuildPackages

		// Needed for running test(s)
		key = `checkdepends`
		if arch != `` {
			key = fmt.Sprintf(`%s_%s`, key, arch)
		}

		m[key] = t.Dependencies[arch].TestPackages
	}

	return m
}

func (t Template) dependsToString() string {
	var arr []string
	for k, v := range t.getDepends() {
		if len(v) == 0 {
			continue
		}
		arr = append(arr, fmt.Sprintf(`%s=(%s)`+"\n", k, wrapStrings(v, ` `, `'`, `'`)))
	}
	return strings.Join(arr, "\n")
}

// Convert to PKGBUILD file
func (t Template) String() string {
	var out bytes.Buffer

	_, _ = fmt.Fprintf(&out, `# Maintainer: %s <%s>`+"\n", t.Maintainer, t.MaintainerEmail)
	_, _ = fmt.Fprintf(&out, `# Generated at: %s `+"\n", time.Now())
	_, _ = fmt.Fprintln(&out)

	if len(t.Name) == 1 {
		_, _ = fmt.Fprintf(&out, `pkgname=%v`+"\n", t.Name[0])
	} else {
		_, _ = fmt.Fprintf(&out, `pkgname=%v`+"\n", wrapStrings(t.Name, ` `, ``, ``))
	}

	_, _ = fmt.Fprintf(&out, `pkgver=%v`+"\n", t.Version)
	_, _ = fmt.Fprintf(&out, `pkgrel=%v`+"\n", t.Release)

	epoch := t.ReleaseTime.Unix()

	if epoch > 0 {
		_, _ = fmt.Fprintf(&out, `epoch=%v`+"\n", epoch)
	}

	_, _ = fmt.Fprintf(&out, `pkgdesc=%q`+"\n", t.ShortDescription)
	_, _ = fmt.Fprintf(&out, `url=%q`+"\n", t.URL)
	_, _ = fmt.Fprintf(&out, `license=(%v)`+"\n", wrapStrings(t.Licenses, ` `, `'`, `'`))
	_, _ = fmt.Fprintf(&out, `arch=(%v)`+"\n", wrapStrings(t.getArchList(), ` `, `'`, `'`))

	if len(t.Options) > 0 {
		_, _ = fmt.Fprintf(&out, `options=(%v)`+"\n", strings.Join(t.Options, ` `))
	}

	if t.Install != `` {
		_, _ = fmt.Fprintf(&out, `install=%v`+"\n", t.Install)
	}

	if t.ChangeLogFile != `` {
		_, _ = fmt.Fprintf(&out, `changelog=%v`+"\n", t.ChangeLogFile)
	}

	if len(t.ValidPGPKeys) > 0 {
		_, _ = fmt.Fprintf(&out, `validpgpkeys=(%v)`+"\n", wrapStrings(t.ValidPGPKeys, ` `, `'`, `'`))
	}

	if len(t.NoExtractFiles) > 0 {
		_, _ = fmt.Fprintf(&out, `noextract=(%v)`+"\n", wrapStrings(t.NoExtractFiles, ` `, `'`, `'`))
	}

	if len(t.Groups) > 0 {
		_, _ = fmt.Fprintf(&out, `groups=(%v)`+"\n", wrapStrings(t.Groups, ` `, `'`, `'`))
	}

	if len(t.Backup) > 0 {
		_, _ = fmt.Fprintf(&out, `backup=(%v)`+"\n", wrapStrings(t.Backup, ` `, `'`, `'`))
	}

	_, _ = fmt.Fprint(&out, t.dependsToString())
	if len(t.OptionalPackages) > 0 {
		_, _ = fmt.Fprint(&out, t.optionalToString()+"\n")
	}

	_, _ = fmt.Fprintln(&out, t.sourceToString())

	// Calculate version from source package
	if len(t.Commands.Version) > 0 {
		_, _ = fmt.Fprintln(&out, "\n"+`pkgver() {`)
		_, _ = fmt.Fprint(&out, `  `+strings.Join(t.Commands.Version, "\n  "))
		_, _ = fmt.Fprintln(&out, "\n}")
	}

	if len(t.Commands.Prepare) > 0 {
		_, _ = fmt.Fprintln(&out, "\n"+`prepare() {`)
		_, _ = fmt.Fprint(&out, `  `+strings.Join(t.Commands.Prepare, "\n  "))
		_, _ = fmt.Fprintln(&out, "\n}")
	}

	if len(t.Commands.Build) > 0 {
		_, _ = fmt.Fprintln(&out, "\n"+`build() {`)
		_, _ = fmt.Fprint(&out, `  `+strings.Join(t.Commands.Build, "\n  "))
		_, _ = fmt.Fprintln(&out, "\n}")
	}

	if len(t.Commands.Test) > 0 {
		_, _ = fmt.Fprintln(&out, "\n"+`check() {`)
		_, _ = fmt.Fprint(&out, `  `+strings.Join(t.Commands.Test, "\n  "))
		_, _ = fmt.Fprintln(&out, "\n}")
	}

	if len(t.Commands.Install) > 0 {
		_, _ = fmt.Fprintln(&out, "\n"+`package() {`)
		_, _ = fmt.Fprint(&out, `  `+strings.Join(t.Commands.Install, "\n  "))
		_, _ = fmt.Fprintln(&out, "\n}")
	}

	return out.String()
}
