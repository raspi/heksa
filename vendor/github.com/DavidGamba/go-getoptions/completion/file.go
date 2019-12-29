// This file is part of go-getoptions.
//
// Copyright (C) 2015-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package completion

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// readDirNoSort - Same as ioutil/ReadDir but doesn't sort results.
//
//   Taken from https://golang.org/src/io/ioutil/ioutil.go
//   Copyright 2009 The Go Authors. All rights reserved.
//   Use of this source code is governed by a BSD-style
//   license that can be found in the LICENSE file.
func readDirNoSort(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// trimLeftDots - Given a string it trims the leading dots (".") and returns a count of how many were removed.
func trimLeftDots(s string) (int, string) {
	charFound := false
	count := 0
	return count, strings.TrimLeftFunc(s, func(r rune) bool {
		if !charFound && r == '.' {
			count++
			return true
		}
		return false
	})
}

// trimLeftDashes - Given a string it trims the leading dashes ("-") and returns a count of how many were removed.
func trimLeftDashes(s string) (int, string) {
	charFound := false
	count := 0
	return count, strings.TrimLeftFunc(s, func(r rune) bool {
		if !charFound && r == '-' {
			count++
			return true
		}
		return false
	})
}

// sortForCompletion - Places hidden files in the same sort possition as their non hidden counterparts.
// Also used for sorting options in the same fashion.
// Example:
//   file.txt
//   .file.txt.~
//   .hidden.txt
//   ..hidden.txt.~
//
//   -d
//   --debug
//   -h
//   --help
func sortForCompletion(list []string) {
	sort.Slice(list,
		func(i, j int) bool {
			var a, b string
			if filepath.Dir(list[i]) == filepath.Dir(list[j]) {
				a = filepath.Base(list[i])
				b = filepath.Base(list[j])
			} else {
				a = list[i]
				b = list[j]
			}

			// . always is less
			if filepath.Base(list[i]) == "." {
				return true
			}
			if filepath.Base(list[j]) == "." {
				return false
			}
			// .. is always less in any other case
			if filepath.Base(list[i]) == ".." {
				return true
			}
			if filepath.Base(list[j]) == ".." {
				return false
			}

			an, a := trimLeftDots(a)
			bn, b := trimLeftDots(b)
			if a == b {
				return an < bn
			}
			an, a = trimLeftDashes(a)
			bn, b = trimLeftDashes(b)
			if a == b {
				return an < bn
			}
			return a < b
		})
}

// listDir - Given a dir and a prefix returns a list of files in the dir filtered by their prefix.
// NOTE: dot (".") is a valid dirname.
func listDir(dirname string, prefix string) ([]string, error) {
	filenames := []string{}
	usedDirname := dirname
	dir := ""
	if strings.Contains(prefix, "/") {
		dir = filepath.Dir(prefix) + string(os.PathSeparator)
		prefix = strings.TrimPrefix(prefix, dir)
		usedDirname = filepath.Join(dirname, dir) + string(os.PathSeparator)
	}
	if prefix == "." {
		filenames = append(filenames, dir+"./")
		filenames = append(filenames, dir+"../")
	} else if prefix == ".." {
		filenames = append(filenames, dir+"../")
	}
	fileInfoList, err := readDirNoSort(usedDirname)
	if err != nil {
		Debug.Printf("listDir - dirname %s, prefix %s > files %v\n", dirname, prefix, filenames)
		return filenames, err
	}
	for _, fi := range fileInfoList {
		name := fi.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		if dirname != usedDirname {
			name = filepath.Join(dir, name)
		}
		if fi.IsDir() {
			filenames = append(filenames, name+"/")
		} else {
			filenames = append(filenames, name)
		}
	}
	sortForCompletion(filenames)
	if len(filenames) == 1 && strings.HasSuffix(filenames[0], "/") {
		filenames = append(filenames, filenames[0]+" ")
	}
	Debug.Printf("listDir - dirname %s, prefix %s > files %v\n", dirname, prefix, filenames)
	return filenames, err
}
