package PKGBUILD

import "strings"

func wrapString(source string, prefix string, suffix string) string {
	return prefix + source + suffix
}

func wrapStrings(source []string, join string, prefix string, suffix string) string {
	var wrapped []string
	for _, s := range source {
		wrapped = append(wrapped, wrapString(s, prefix, suffix))
	}
	return strings.Join(wrapped, join)
}

// Rewrite Go's arch names to Arch Linux ones
var GoArchToLinuxArch = map[string]string{
	`amd64`:   `x86_64`,
	`arm`:     `arm`,
	`arm64`:   `aarch64`,
	`ppc64`:   `ppc64`,
	`ppc64le`: `ppc64le`,
	`x86_64`:  `x86_64`,
	`aarch64`: `aarch64`,
}
