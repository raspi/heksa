package color

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const (
	esc             = "\033["
	Clear           = esc + "0m"
	SetForeground   = esc + "38;5;"
	SetUnderlineOn  = esc + "4m"
	SetUnderlineOff = esc + "24m"
)

func GetColorGroupColorDefaults(src io.Reader, required []string) (groupcolors map[string]string, err error) {
	reqFound := make(map[string]bool)

	for _, name := range required {
		reqFound[name] = false
	}

	groupcolors = make(map[string]string)

	m, err := regexp.Compile(`^([^=]+)=([^\s]+)`)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, ";") {
			// Comment line, skip
			continue
		}
		found := m.FindStringSubmatch(line)
		if found == nil {
			continue
		}

		groupName, param := found[1], found[2]

		switch groupName {
		case `LineEven`, `LineOdd`: // Background
			reqFound[groupName] = true
			groupcolors[groupName] = esc + param + "m"
		default:
			reqFound[groupName] = true
			groupcolors[groupName] = SetForeground + param + "m"
		}
	}

	for name, found := range reqFound {
		if !found {
			return nil, fmt.Errorf(`group %+v was not found`, name)
		}
	}

	return groupcolors, nil
}

func GetColors(src io.Reader, groups map[string]string) (byteColors [256]string, err error) {
	for i := 0; i < 256; i++ {
		// Set default color for each byte
		byteColors[i] = groups[`Default`]
	}

	colorMatcher, err := regexp.Compile(`(\d+)`)
	if err != nil {
		panic(err)
	}

	m, err := regexp.Compile(`^([^=]+)=([^\r\n]+)`)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, ";") {
			// Comment line, skip
			continue
		}

		found := m.FindStringSubmatch(line)
		if found == nil {
			continue
		}

		groupName, params := found[1], found[2]

		byteList := colorMatcher.FindAllString(params, 1024)
		if byteList == nil {
			continue
		}

		for _, c := range byteList {
			n, err := strconv.Atoi(c)
			if err != nil {
				return byteColors, err
			}

			byteColors[n] = groups[groupName]
		}

	}

	return byteColors, nil
}
