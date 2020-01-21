# heksa

![Screenshot](https://github.com/raspi/heksa/blob/master/_assets/screenshot.png)

![GitHub All Releases](https://img.shields.io/github/downloads/raspi/heksa/total?style=for-the-badge)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/raspi/heksa?style=for-the-badge)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/raspi/heksa?style=for-the-badge)

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

## heksa --help

```
heksa - hex file dumper v1.11.0 - (2020-01-18T19:33:11+02:00)
(c) Pekka Järvinen 2019- [ https://github.com/raspi/heksa ]
SYNOPSIS:
    heksa [--format|-f <fmt1,fmt2,..>] [--help|-h|-?]
          [--limit|-l <[prefix]bytes[unit]>] [--offset-format|-o <fmt1[,fmt2]>]
          [--seek|-s <[prefix]offset[unit]>] [--version] <filename> or STDIN

OPTIONS:
    --format|-f <fmt1,fmt2,..>          One or multiple of: asc, bit, bitwasc, bitwdec, bitwhex, dec, decwasc, hex, hexwasc, oct (default: "hex,asc")

    --help|-h|-?                        Show this help (default: false)

    --limit|-l <[prefix]bytes[unit]>    Read only N bytes (0 = no limit). See NOTES. (default: "0")

    --offset-format|-o <fmt1[,fmt2]>    One or two of: dec, hex, oct, per, no, ''.
                                        First one is displayed on the left side and second one on right side after formatters. (default: "hex")

    --seek|-s <[prefix]offset[unit]>    Start reading from certain offset. See NOTES. (default: "0")

    --version                           Show version information (default: false)


NOTES:
    - You can use prefixes for seek and limit. 0x = hex, 0b = binary, 0o = octal
    - Use 'no' or '' for offset formatter for disabling offset output
    - Use '--seek \-[prefix]1000' for seeking to end of file
    - Offset and seek parameters supports units (KB, KiB, MB, MiB, GB, GiB, TB, TiB)

EXAMPLES:
    heksa -f hex,asc,bit foo.dat
    heksa -o hex,per -f hex,asc foo.dat
    heksa -o hex -f hex,asc,bit foo.dat
    heksa -o no -f bit foo.dat
    heksa -l 0x1024 foo.dat
    heksa -s 0b1010 foo.dat
    heksa -s 4321KiB foo.dat
```

![Screenshot](https://github.com/raspi/heksa/blob/master/_assets/screenshot2.png)

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

If there's a lot of visual changes you can take new screenshots with `screenshot.sh` script in [_assets](_assets) directory
