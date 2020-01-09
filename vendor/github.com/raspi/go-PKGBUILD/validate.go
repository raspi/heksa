package PKGBUILD

import (
	"fmt"
	"regexp"
	"strings"
)

func (t Template) Validate() (errs []error) {

	if t.Maintainer == `` {
		errs = append(errs, fmt.Errorf(`maintainer is empty`))
	}

	if t.MaintainerEmail == `` {
		errs = append(errs, fmt.Errorf(`maintainer email is empty`))
	}

	if len(t.Name) == 0 {
		errs = append(errs, fmt.Errorf(`name is empty`))
	} else {
		for idx, name := range t.Name {
			E := t.validateName(name)
			if E != nil {
				for _, ER := range E {
					errs = append(errs, fmt.Errorf(`error in name #%d: %v`, idx, ER))
				}
			}
		}
	}

	if t.Version == `` {
		errs = append(errs, fmt.Errorf(`version is empty`))
	} else {
		for _, nv := range []string{`/`, `-`, ` `, `:`} {
			if strings.Contains(t.Version, nv) {
				errs = append(errs, fmt.Errorf(`invalid character in version %q: %q`, t.Version, nv))
			}
		}
	}

	if t.Release == 0 {
		errs = append(errs, fmt.Errorf(`release should be at least 1`))
	}

	if len(t.Licenses) == 0 {
		errs = append(errs, fmt.Errorf(`no licence(s) given`))
	}

	if t.ShortDescription == `` {
		errs = append(errs, fmt.Errorf(`short description is empty`))
	} else {
		for _, nv := range []string{"\t", "\n", "\r"} {
			if strings.Contains(t.ShortDescription, nv) {
				errs = append(errs, fmt.Errorf(`invalid character in short description %q: %q`, t.ShortDescription, nv))
			}
		}
	}

	if t.URL == `` {
		errs = append(errs, fmt.Errorf(`url is empty`))
	}

	return errs
}

func (t Template) validateName(name string) (errs []error) {
	if strings.HasPrefix(name, `-`) {
		errs = append(errs, fmt.Errorf(`%q can't start with '-'`, name))
	}

	if strings.HasPrefix(name, `.`) {
		errs = append(errs, fmt.Errorf(`%q can't start with '.'`, name))
	}

	validRe := regexp.MustCompile(`^[a-z0-9_+@\.\-]+$`)

	if !validRe.MatchString(name) {
		errs = append(errs, fmt.Errorf(`invalid characters in name %q`, name))
	}

	return errs
}
