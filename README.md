# heksa

![Screenshot](https://github.com/raspi/heksa/blob/master/_assets/screenshot.png)

![GitHub All Releases](https://img.shields.io/github/downloads/raspi/heksa/total?style=for-the-badge)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/raspi/heksa?style=for-the-badge)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/raspi/heksa?style=for-the-badge)

[![Packaging status](https://repology.org/badge/vertical-allrepos/heksa.svg)](https://repology.org/project/heksa/versions)

Hex dumper with colors

## Features

* ANSI colors for different byte groups such as 
  * Printable: A-Z, a-z, 0-9
  * Spaces: space, tab, new line
  * Special: 0x00, 0xFF
* Output multiple formats at once (hexadecimal, decimal, octal, bits or special combination formats)
* Multiple offset formats (hexadecimal, decimal, octal, percentage)
  * First one is displayed on left side and second one on the right side
* Read only N bytes
* Seek to given offset
  * also reads from end of file when using minus sign
* Seek and limit supports 
  * Prefixes hex (`0x`), octal (`0o`) and binary (`0b`)
  * Units (KB, KiB, MB, MiB, GB, GiB, TB, TiB)
* Read from stdin

![Screenshot](https://github.com/raspi/heksa/blob/master/_assets/screenshot2.png)

## heksa --help

```
heksa - hex file dumper v1.14.0 - (2021-05-18T16:20:59+03:00)
(c) Pekka Järvinen 2019- [ https://github.com/raspi/heksa ]
SYNOPSIS:
    heksa [--format|-f <fmt1,fmt2,..>] [--help|-h|-?]
          [--limit|-l <[prefix]bytes[unit]>] [--offset-format|-o <fmt1[,fmt2]>]
          [--print-relative-offset|-r] [--seek|-s <[prefix]offset[unit]>]
          [--splitter|-S <size>] [--version] [--width|-w <[prefix]width>]
          <filename> or STDIN

OPTIONS:
    --format|-f <fmt1,fmt2,..>          One or multiple of: asc, bit, bitwasc, bitwdec, bitwhex, blk, dec, decwasc, hex, hexwasc, oct (default: "hex,asc")

    --help|-h|-?                        Show this help (default: false)

    --limit|-l <[prefix]bytes[unit]>    Read only N bytes (0 = no limit). See NOTES. (default: "0")

    --offset-format|-o <fmt1[,fmt2]>    One or two of: dec, hex, humiec, humsi, oct, per, no, ''.
                                        First one is displayed on the left side and second one on right side after formatters. (default: "hex")

    --print-relative-offset|-r          Print relative offset(s) starting from 0 (file only) (default: false)

    --seek|-s <[prefix]offset[unit]>    Start reading from certain offset. See NOTES. (default: "0")

    --splitter|-S <size>                Insert visual splitter every N bytes. Zero (0) disables. (default: 8)

    --version                           Show version information (default: false)

    --width|-w <[prefix]width>          Width. See NOTES. (default: "16")


NOTES:
    - You can use prefixes for seek, limit and width. 0x = hex, 0b = binary, 0o = octal
    - Use '--seek \-1234' for seeking from end of file
    - Limit and seek parameters supports units (KB, KiB, MB, MiB, GB, GiB, TB, TiB)
    - --print-relative-offset can be used when seeking to certain offset to also print extra offset position starting from zero
    - Offset formatters:
      - Disable formatter output with 'no' or ''
      - 'humiec' (IEC: 1024 B) and 'humsi' (SI: 1000 B) displays offset in human form (n KiB/KB)
    - Formatters:
      - 'blk' can be used to print simple color blocks which helps to visualize where data vs. human readable strings are

EXAMPLES:
    heksa -f hex,asc,bit foo.dat
    heksa -o hex,per -f hex,asc foo.dat
    heksa -o hex -f hex,asc,bit foo.dat
    heksa -o no -f bit foo.dat
    heksa -l 0x1024 foo.dat
    heksa -s 0b1010 foo.dat
    heksa -s 4321KiB foo.dat
    heksa -w 8 foo.dat
    echo "test" | heksa
```

## Requirements

* Terminal with ANSI color support
  * [KDE](https://kde.org/)'s [Konsole](https://konsole.kde.org/) is currently used for development
* Operating system
  * [GNU/Linux](https://www.gnu.org/distros/distros.html)
    * x64 arm arm64 ppc64 ppc64le
  * [Microsoft Windows](https://www.microsoft.com/en-us/windows)
    * x64
  * [Darwin](https://www.apple.com/macos/) (Apple Mac)
    * x64
  * [FreeBSD](https://www.freebsd.org/)
    * x64 arm
  * [NetBSD](https://www.netbsd.org/)
    * x64 arm
  * [OpenBSD](https://www.openbsd.org/)
    * x64 arm arm64
  * Other OSes supported by [Go](https://golang.org)
    * For full list, see: https://golang.org/doc/install/source#environment

## Get source

    git clone https://github.com/raspi/heksa

## Contributing and helping with the project

See [CONTRIBUTING.md](CONTRIBUTING.md) and [current issues](https://github.com/raspi/heksa/issues) that might need help.

## Developing

1. Make changes
1. `make build` or just `go build .`

## Releasing new version:

Requirements:

* `upx` for compressing executables

1. Create new version tag
1. `make release`

If there's a lot of visual changes you can take new screenshots with `screenshot.sh` script in [_assets](_assets) directory
