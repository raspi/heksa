// This file is part of go-getoptions.
//
// Copyright (C) 2015-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package help

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/DavidGamba/go-getoptions/option"
	"github.com/DavidGamba/go-getoptions/text"
)

// Indentation - Number of spaces used for indentation.
var Indentation = 4

func indent(s string) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", Indentation), s)
}

func wrapFn(wrap bool, open, close string) func(s string) string {
	if wrap {
		return func(s string) string {
			return fmt.Sprintf("%s%s%s", open, s, close)
		}
	}
	return func(s string) string {
		return s
	}
}

// Name -
func Name(scriptName, name, description string) string {
	out := scriptName
	if scriptName != "" {
		out += fmt.Sprintf(" %s", name)
	} else {
		out += name
	}
	if description != "" {
		out += fmt.Sprintf(" - %s", description)
	}
	return fmt.Sprintf("%s:\n%s\n", text.HelpNameHeader, indent(out))
}

// Synopsis - Return a default synopsis.
func Synopsis(scriptName, name, args string, options []*option.Option, commands []string) string {
	synopsisName := scriptName
	if scriptName != "" {
		synopsisName += fmt.Sprintf(" %s", name)
	} else {
		synopsisName += name
	}
	synopsisName = indent(synopsisName)
	if name != "" {
		scriptName += " " + name
	}
	normalOptions := []*option.Option{}
	requiredOptions := []*option.Option{}
	for _, option := range options {
		if option.IsRequired {
			requiredOptions = append(requiredOptions, option)
		} else {
			normalOptions = append(normalOptions, option)
		}
	}
	option.Sort(normalOptions)
	option.Sort(requiredOptions)
	optSynopsis := func(opt *option.Option) string {
		txt := ""
		wrap := wrapFn(!opt.IsRequired, "[", "]")
		switch opt.OptType {
		case option.BoolType, option.StringType, option.IntType, option.Float64Type:
			txt += wrap(opt.HelpSynopsis)
		case option.StringRepeatType, option.IntRepeatType, option.StringMapType:
			if opt.IsRequired {
				wrap = wrapFn(opt.IsRequired, "<", ">")
			}
			txt += wrap(opt.HelpSynopsis) + "..."
		}
		return txt
	}
	var out string
	line := synopsisName
	for _, option := range append(requiredOptions, normalOptions...) {
		syn := optSynopsis(option)
		// fmt.Printf("%d - %d - %d | %s | %s\n", len(line), len(syn), len(line)+len(syn), syn, line)
		if len(line)+len(syn) > 80 {
			out += line + "\n"
			line = fmt.Sprintf("%s %s", strings.Repeat(" ", len(synopsisName)), syn)
		} else {
			line += fmt.Sprintf(" %s", syn)
		}
	}
	syn := ""
	if len(commands) > 0 {
		syn += "<command> "
	}
	if args == "" {
		syn += "[<args>]"
	} else {
		syn += args
	}
	if len(line)+len(syn) > 80 {
		out += line + "\n"
		line = fmt.Sprintf("%s %s", strings.Repeat(" ", len(synopsisName)), syn)
	} else {
		line += fmt.Sprintf(" %s", syn)
	}
	out += line
	return fmt.Sprintf("%s:\n%s\n", text.HelpSynopsisHeader, out)
}

// CommandList -
// commandMap => name: description
func CommandList(commandMap map[string]string) string {
	if len(commandMap) <= 0 {
		return ""
	}
	names := []string{}
	for name := range commandMap {
		names = append(names, name)
	}
	sort.Strings(names)
	factor := longestStringLen(names)
	out := ""
	for _, command := range names {
		out += indent(fmt.Sprintf("%s    %s\n", pad(true, command, factor), commandMap[command]))
	}
	return fmt.Sprintf("%s:\n%s", text.HelpCommandsHeader, out)
}

// longestStringLen - Given a slice of strings it returns the length of the longest string in the slice
func longestStringLen(s []string) int {
	i := 0
	for _, e := range s {
		if len(e) > i {
			i = len(e)
		}
	}
	return i
}

// pad - Given a string and a padding factor it will return the string padded with spaces.
//
// Example:
//     pad(true, "--flag", 8) -> '--flag  '
func pad(do bool, s string, factor int) string {
	if do {
		return fmt.Sprintf("%-"+strconv.Itoa(factor)+"s", s)
	}
	return s
}

// OptionList - Return a formatted list of options and their descriptions.
func OptionList(options []*option.Option) string {
	synopsisLength := 0
	normalOptions := []*option.Option{}
	requiredOptions := []*option.Option{}
	for _, opt := range options {
		l := len(opt.HelpSynopsis)
		if l > synopsisLength {
			synopsisLength = l
		}
		if opt.IsRequired {
			requiredOptions = append(requiredOptions, opt)
		} else {
			normalOptions = append(normalOptions, opt)
		}
	}
	option.Sort(normalOptions)
	option.Sort(requiredOptions)
	helpString := func(opt *option.Option) string {
		txt := ""
		factor := synopsisLength + 4
		padding := strings.Repeat(" ", factor)
		txt += indent(pad(!opt.IsRequired || opt.Description != "", opt.HelpSynopsis, factor))
		if opt.Description != "" {
			description := strings.Replace(opt.Description, "\n", "\n    "+padding, -1)
			txt += fmt.Sprintf("%s", description)
		}
		if !opt.IsRequired {
			if opt.Description != "" {
				txt += " "
			}
			txt += fmt.Sprintf("(default: %s)\n\n", opt.DefaultStr)
		} else {
			txt += "\n\n"
		}
		return txt
	}
	out := ""
	if len(requiredOptions) > 0 {
		out += fmt.Sprintf("%s:\n", text.HelpRequiredOptionsHeader)
		for _, option := range requiredOptions {
			out += helpString(option)
		}
	}
	if len(normalOptions) > 0 {
		out += fmt.Sprintf("%s:\n", text.HelpOptionsHeader)
		for _, option := range normalOptions {
			out += helpString(option)
		}
	}
	return out
}
