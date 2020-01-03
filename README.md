# heksa

![Screenshot](https://github.com/raspi/heksa/blob/master/_assets/screenshot.png)

![GitHub All Releases](https://img.shields.io/github/downloads/raspi/heksa/total?style=for-the-badge)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/raspi/heksa?style=for-the-badge)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/raspi/heksa?style=for-the-badge)

Hex dumper with colors

```
heksa - hex file dumper v1.9.0 - (2020-01-02T07:55:55+02:00)
(c) Pekka Järvinen 2019- [ https://github.com/raspi/heksa ]
SYNOPSIS:
    heksa [--format|-f <fmt1,fmt2,..>] [--header|-H] [--help|-h|-?]
          [--limit|-l <[prefix]bytes>] [--offset-format|-o <fmt1[,fmt2]>]
          [--seek|-s <[prefix]offset>] [--version] <filename> or STDIN

OPTIONS:
    --format|-f <fmt1,fmt2,..>          One or multiple of: hex, dec, oct, bit (default: "hex,asc")

    --header|-H                         Show offset header (default: false)

    --help|-h|-?                        Show this help (default: false)

    --limit|-l <[prefix]bytes>          Read only N bytes (0 = no limit). See NOTES. (default: "0")

    --offset-format|-o <fmt1[,fmt2]>    One or two of: hex, dec, oct, per, no, ''. First one is displayed on the left side and second one on right side after formatters (default: "hex")

    --seek|-s <[prefix]offset>          Start reading from certain offset. See NOTES. (default: "0")

    --version                           Show version information (default: false)

NOTES:
    - You can use prefixes for seek and limit. 0x = hex, 0b = binary, 0o = octal.
    - Use 'no' or '' for offset formatter for disabling offset output.

EXAMPLES:
    heksa -f hex,asc,bit foo.dat
    heksa -o hex,per -f hex,asc foo.dat
    heksa -o hex -f hex,asc,bit foo.dat
    heksa -o no -f bit foo.dat
    heksa -l 0x1024 foo.dat
    heksa -s 0b1010 foo.dat
```

## Features

* ANSI colors for different bytes
* Output multiple formats at once (hexadecimal, decimal, octal, bits)
* Multiple offset formats (hexadecimal, decimal, octal)
* Read only N bytes
* Seek to given offset
* Read from stdin

## Requirements

* Terminal with ANSI color support
  * KDE's Konsole is currently used for development
* Operating system
  * GNU/Linux 
    * x64 arm arm64 ppc64 ppc64le
  * Microsoft Windows
    * x64
  * Darwin (Apple Mac)
    * x64
  * FreeBSD
    * x64 arm
  * NetBSD
    * x64 arm
  * OpenBSD
    * x64 arm arm64

## Get source

    git clone https://github.com/raspi/heksa

## Developing

1. Make changes
1. `make build` or just `go build .`

## Releasing new version:

Requirements:

* `upx` for compressing executables

1. Create new version tag
1. `make release`

If there's a lot of changes you can take a new screenshot with `make screenshot` helper
