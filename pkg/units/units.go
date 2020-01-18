package units

import (
	"strconv"
	"strings"
)

const (
	KB  = int64(1000)
	KiB = int64(1024)
	MB  = 1000 * KB
	MiB = 1024 * KiB
	GB  = 1000 * MB
	GiB = 1024 * MiB
	TB  = 1000 * GB
	TiB = 1024 * GiB
)

var units = map[string]int64{
	`KB`:  KB,
	`KiB`: KiB,
	`MB`:  MB,
	`MiB`: MiB,
	`GB`:  GB,
	`GiB`: GiB,
	`TB`:  TB,
	`TiB`: TiB,
}

// Parse parses string which may contain unit into int64
// Examples: 1000 -1000 123KiB 432MB 0x100 0o100 0b1111
func Parse(s string) (n int64, err error) {
	multiplier := int64(0)
	for u, v := range units {
		if strings.HasSuffix(s, u) {
			multiplier = v
			s = strings.TrimRight(s, u)
			break
		}
	}

	n, err = strconv.ParseInt(s, 0, 64)

	if multiplier != 0 {
		n = n * multiplier
	}

	return n, err
}
