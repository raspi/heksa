// This file is part of go-getoptions.
//
// Copyright (C) 2015-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// argList tracks passed arguments with an index and allows to peek ahead.

package getoptions

// argList - arguments list
type argList struct {
	list     []string // Original list
	listSize int      // Original list size
	idx      int
}

func newArgList(a []string) *argList {
	return &argList{list: a, listSize: len(a), idx: -1}
}

func (a *argList) size() int {
	return a.listSize
}

func (a *argList) index() int {
	return a.idx
}

func (a *argList) next() bool {
	a.idx++
	return a.idx < a.listSize
}

func (a *argList) existsNext() bool {
	return a.idx+1 < a.listSize
}

func (a *argList) value() string {
	return a.list[a.idx]
}

func (a *argList) peekNextValue() string {
	return a.list[a.idx+1]
}

func (a *argList) remaining() []string {
	return a.list[a.idx:]
}
