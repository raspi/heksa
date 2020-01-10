name: heksa
version: <VERSION>
origin: category/devel
comment: CLI hex dumper with colors
arch: <ARCH>
abi: freebsd:*:<ARCH>
www: https://github.com/raspi/heksa
maintainer: Pekka Järvinen
prefix: /usr/local
licenselogic: single
licenses: [Apache2]
desc: <<EOD
heksa is a command line hex binary dumper which uses ANSI colors
EOD
categories: [devel]
files: {
  /usr/local/bin/heksa: 'sha256sum',
}
directories: {
}
