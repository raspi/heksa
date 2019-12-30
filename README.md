# heksa

![Screenshot](https://github.com/raspi/heksa/blob/master/_assets/screenshot.png)

Hex dumper with colors

```
% heksa -h
heksa - hex file dumper v1.4.0 build 005bf566df570332cac04ab61e350dd3aacf79f3
(c) Pekka Järvinen 2019- [ https://github.com/raspi/heksa ]
SYNOPSIS:
    heksa [--format|-f <fmt1,fmt2,..>] [--help|-h|-?] [--limit|-l <bytes>]
          [--offset-format|-o <[fmt1][,fmt2]>] [--seek|-s <offset>] <filename>

OPTIONS:
    --format|-f <fmt1,fmt2,..>            One or multiple of: hex, dec, oct, bit (default: "hex,asc")

    --help|-h|-?                          Show this help (default: false)

    --limit|-l <bytes>                    Read only N bytes (0 = no limit) (default: 0)

    --offset-format|-o <[fmt1][,fmt2]>    Zero to two of: hex, dec, oct, per. First one is displayed on the left side and second one on right after formatters (default: "hex")

    --seek|-s <offset>                    Start reading from certain offset (default: 0)

EXAMPLES:
    heksa -f hex,asc,bit foo.dat
    heksa -o hex,per -f hex,asc foo.dat
    heksa -o hex -f hex,asc,bit foo.dat
    heksa -o '' -f bit foo.dat
    heksa -l 1024 foo.dat
    heksa -s 1234 foo.dat
```
## Get source

    go get -u github.com/raspi/heksa
