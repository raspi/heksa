# heksa
Hex dumper with colors

```
% heksa -h
heksa - hex file dumper v1.1.1 build e52b556b08e65e4f76dd3db17789299ac047b23b
(c) Pekka Järvinen 2019- [ https://github.com/raspi/heksa ]
SYNOPSIS:
    heksa [--format|-f <fmt1,fmt2,..>] [--help|-h|-?] [--limit|-l <bytes>]
          [--offset-display|-o <offset format>] [--seek|-s <offset>] <filename>

OPTIONS:
    --format|-f <fmt1,fmt2,..>             One or multiple of: hex, dec, oct, bit (default: "hex,asc")

    --help|-h|-?                           Show this help (default: false)

    --limit|-l <bytes>                     Read only N bytes (0 = no limit) (default: 0)

    --offset-display|-o <offset format>    One of: hex, dec (default: "hex")

    --seek|-s <offset>                     Start reading from certain offset (default: 0)

EXAMPLES:
    heksa -f hex,asc,bit foo.dat
    heksa -o hex -f hex,asc,bit foo.dat
    heksa -o hex -f bit foo.dat
    heksa -l 1024 foo.dat
    heksa -s 1234 foo.dat
```

![Screenshot](https://github.com/raspi/heksa/blob/master/_assets/screenshot.png)
