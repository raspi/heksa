package PKGBUILD

import (
	"encoding/json"
	"time"
)

type Commands struct {
	Prepare []string `json:"prepare,omitempty"` // 1. prepare(){...}  What is needed before building (for example apply source patches here)
	Build   []string `json:"build,omitempty"`   // 2. build(){...} How it's built from source?
	Test    []string `json:"test,omitempty"`    // 3. check(){...} Tests before install
	Install []string `json:"install,omitempty"` // 4. package(){...} How package is installed?
}

type OptionalPackage struct {
	Package string `json:"package"` // Package name
	Reason  string `json:"reason"`  // Description why it's needed
}

type Depends struct {
	Packages      []string `json:"packages,omitempty"`       // PKGBUILD $Dependencies package names this package Dependencies for running it
	BuildPackages []string `json:"build_packages,omitempty"` // PKGBUILD $makedepends package names when building this package from source
	TestPackages  []string `json:"test_packages,omitempty"`  // PKGBUILD $checkdepends package names test suite Dependencies on
}

// Source files
type Source struct {
	URL       string            `json:"url"`             // URL to file
	Alias     string            `json:"alias,omitempty"` // [OPTIONAL] Rename file to this, leave empty for no aliasing
	Checksums map[string]string `json:"checksums"`       // [type]checksums for file, for example [sha256]1234..
}

type Files map[string][]Source // [architecture][]{file1, file2, ..}

type ProviderOperator uint8
type Provider struct {
	Name     string
	Operator ProviderOperator
	Version  string
}

type Meta struct {
	Version string `json:"ver"`
}

type Template struct {
	Meta            Meta   `json:"_meta"`            // Used by this library for possible meta extension(s), not used in PKGBUILD file
	Maintainer      string `json:"maintainer"`       // Maintainer's name (# comment)
	MaintainerEmail string `json:"maintainer_email"` // Maintainer's email address (# comment)

	// Either the name of the package or an array of names for split packages. Valid characters for members of
	// this array are alphanumerics, and any of the following characters: “@ . _ + -”. Additionally, names
	// are not allowed to start with hyphens or dots.
	Name []string `json:"name"` // $pkgname Package name(s)

	// The version of the software as released from the author (e.g., 2.7.1). The variable is
	// not allowed to contain colons, forward slashes, hyphens or whitespace.
	//
	// The pkgver variable can be automatically updated by providing a pkgver() function in the PKGBUILD that
	// outputs the new package version. This is run after downloading and extracting the Files and
	// running the prepare() function (if present), so it can use those files in determining the new pkgver.
	// This is most useful when used with Files from version control systems.
	Version          string                       `json:"version"`                // $pkgver Package version
	Release          uint64                       `json:"release"`                // $pkgrel Increment after each PKGBUILD release
	ReleaseTime      time.Time                    `json:"release_time,omitempty"` // [OPTIONAL] $epoch ReleaseTime for hinting newer release
	ShortDescription string                       `json:"short_description"`      // $pkgdesc Short one line description about the package
	Licenses         []string                     `json:"licenses"`               // $license License(s)
	URL              string                       `json:"url"`                    // $url Package homepage URL
	ChangeLogFile    string                       `json:"changelog_file"`         // $changelog [OPTIONAL]
	Groups           []string                     `json:"groups"`
	Dependencies     map[string]Depends           `json:"dependencies"`                // [OPTIONAL] [architecture]Dependencies..
	OptionalPackages map[string][]OptionalPackage `json:"optional_packages,omitempty"` // PKGBUILD $optdepends

	// An array of “virtual provisions” this package provides. This allows a package to provide dependencies other
	// than its own package name. For example, the dcron package can provide cron, which allows packages to depend
	// on cron rather than dcron OR fcron.
	//
	// Versioned provisions are also possible, in the name=version format. For example, dcron can provide cron=2.0 to
	// satisfy the cron>=2.0 dependency of other packages. Provisions
	// involving the > and < operators are invalid as only specific versions of a package may be provided.
	//
	// If the provision name appears to be a library (ends with .so), makepkg will try to find the library in the
	// built package and append the correct version. Appending the version yourself disables automatic detection.
	//
	// Additional architecture-specific provides can be added by appending an underscore and the architecture name e.g., provides_x86_64=().
	Provides  map[string]Provider `json:"provides"`            // [OPTIONAL]
	Conflicts map[string][]string `json:"conflicts,omitempty"` // [OPTIONAL] [arch][]{pkg1,pkg2, ..}

	// An array of packages this package should replace. This can be used to handle renamed/combined packages.
	// For example, if the j2re package is renamed to jre, this directive allows future upgrades to continue as
	// expected even though the package has moved. Versioned replaces are supported using the operators as described
	// in Dependencies.
	//
	// Sysupgrade is currently the only pacman operation that utilizes this field. A normal sync or upgrade will not use its value.
	//
	// Additional architecture-specific replaces can be added by appending an underscore and the architecture name e.g., replaces_x86_64=().
	Replaces map[string][]string `json:"replaces,omitempty"` // [arch][]{pkg1,pkg2, ..}

	// This allows you to override some of makepkg’s default behavior when building packages.
	Options []string `json:"options,omitempty"`

	// Specifies a special install script that is to be included in the package. This file should reside in the same
	// directory as the PKGBUILD and will be copied into the package by makepkg. It does not need to be included in
	// the source array (e.g., install=$pkgname.install).
	Install string `json:"install"`

	// An array of file names corresponding to those from the source array. Files listed here will not be extracted
	// with the rest of the source files. This is useful for packages that use compressed data directly.
	NoExtractFiles []string `json:"no_extract_files,omitempty"`

	// An array of PGP fingerprints. If this array is non-empty, makepkg will only accept signatures from the keys
	// listed here and will ignore the trust values from the keyring. If the source file was signed with a subkey,
	// makepkg will still use the primary key for comparison.
	// Only full fingerprints are accepted. They must be uppercase and must not contain whitespace characters.
	ValidPGPKeys []string `json:"valid_pgp_keys,omitempty"`

	// An array of file names, without preceding slashes, that should be backed up if the package is removed or upgraded.
	// This is commonly used for packages placing configuration files in /etc.
	// See "Handling Config Files" in pacman(8) for more information.
	Backup []string `json:"backup,omitempty"`
	Files  Files    `json:"files"` // [architecture][]{file1, file2, ..}

	Commands Commands `json:"commands"`
}

func New(sources Files, cmds Commands, depends map[string]Depends, optional map[string][]OptionalPackage, options []string) Template {
	return Template{
		Meta: Meta{
			Version: "v1.0.0",
		},
		Files:            sources,
		Dependencies:     depends,
		OptionalPackages: optional,
		Commands:         cmds,
		Release:          1,
		ReleaseTime:      time.Unix(0, 0), // Default to zero and use $pkgver version as reference
		Options:          options,
	}
}

// Read template from JSON file
func FromJson(source []byte) (tpl Template, err error) {
	err = json.Unmarshal(source, &tpl)
	return tpl, err
}
