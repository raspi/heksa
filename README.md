# heksa

![Screenshot](https://github.com/raspi/heksa/blob/master/_assets/screenshot.png)

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

## Get source

    go get -u github.com/raspi/heksa
