// This file is part of go-getoptions.
//
// Copyright (C) 2015-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package option - internal option struct and methods.
package option

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/DavidGamba/go-getoptions/text"
)

// Debug Logger instance set to `ioutil.Discard` by default.
// Enable debug logging by setting: `Debug.SetOutput(os.Stderr)`.
var Debug = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

// Handler - Signature for the function that handles saving to the option.
type Handler func(optName string, argument string, usedAlias string) error

// Type - Indicates the type of option.
type Type int

// Option Types
const (
	BoolType Type = iota
	StringType
	IntType
	Float64Type
	StringRepeatType
	IntRepeatType
	StringMapType
)

// Option - main object
type Option struct {
	Name           string
	Aliases        []string
	Called         bool    // Indicates if the option was passed on the command line
	UsedAlias      string  // Alias used when the option was called
	Handler        Handler // method used to handle the option
	IsOptional     bool    // Indicates if an option has an optional argument
	MapKeysToLower bool    // Indicates if the option of map type has it keys set ToLower
	OptType        Type    // Option Type
	MinArgs        int     // minimum args when using multi
	MaxArgs        int     // maximum args when using multi

	IsRequired    bool   // Indicates if the option is required
	IsRequiredErr string // Error message for the required option

	IsPassedToCommand bool // Indicates if the option was passed from root to child

	// Help
	DefaultStr   string // String representation of default value
	Description  string // Optional description used for help
	HelpArgName  string // Optional arg name used for help
	HelpSynopsis string // Help synopsis

	// Pointer receivers:
	value    interface{}        // Value without type safety
	pBool    *bool              // receiver for bool pointer
	pString  *string            // receiver for string pointer
	pInt     *int               // receiver for int pointer
	pFloat64 *float64           // receiver for float64 pointer
	pStringS *[]string          // receiver for string slice pointer
	pIntS    *[]int             // receiver for int slice pointer
	pStringM *map[string]string // receiver for string map pointer
}

// New - Returns a new option object
func New(name string, optType Type) *Option {
	opt := &Option{
		Name:    name,
		OptType: optType,
		Aliases: []string{name},
	}
	switch optType {
	case StringType, StringRepeatType:
		opt.HelpArgName = "string"
	case IntType, IntRepeatType:
		opt.HelpArgName = "int"
	case Float64Type:
		opt.HelpArgName = "float64"
	case StringMapType:
		opt.HelpArgName = "key=value"
	}
	opt.synopsis()
	return opt
}

func (opt *Option) synopsis() {
	aliases := []string{}
	for _, e := range opt.Aliases {
		if len(e) > 1 {
			e = "--" + e
		} else {
			e = "-" + e
		}
		aliases = append(aliases, e)
	}
	opt.HelpSynopsis = strings.Join(aliases, "|")
	if opt.OptType != BoolType {
		opt.HelpSynopsis = fmt.Sprintf("%s <%s>", opt.HelpSynopsis, opt.HelpArgName)
	}
	if opt.MaxArgs > 1 {
		opt.HelpSynopsis = opt.HelpSynopsis + "..."
	}
}

// Value - Get untyped option value
func (opt *Option) Value() interface{} {
	return opt.value
}

// SetAlias - Adds aliases to an option.
func (opt *Option) SetAlias(alias ...string) *Option {
	opt.Aliases = append(opt.Aliases, alias...)
	opt.synopsis()
	return opt
}

// SetDescription - Updates the Description.
func (opt *Option) SetDescription(s string) *Option {
	opt.Description = s
	return opt
}

// SetHelpArgName - Updates the HelpArgName.
func (opt *Option) SetHelpArgName(s string) *Option {
	opt.HelpArgName = s
	opt.synopsis()
	return opt
}

// SetDefaultStr - Updates the DefaultStr.
func (opt *Option) SetDefaultStr(s string) *Option {
	opt.DefaultStr = s
	return opt
}

// SetRequired - Marks an option as required.
func (opt *Option) SetRequired(msg string) *Option {
	opt.IsRequired = true
	opt.IsRequiredErr = msg
	return opt
}

// CheckRequired - Returns error if the option is required.
func (opt *Option) CheckRequired() error {
	if opt.IsRequired {
		if !opt.Called {
			if opt.IsRequiredErr != "" {
				return fmt.Errorf(opt.IsRequiredErr)
			}
			return fmt.Errorf(text.ErrorMissingRequiredOption, opt.Name)
		}
	}
	return nil
}

// SetCalled - Marks the option as called and records the alias used to call it.
func (opt *Option) SetCalled(usedAlias string) *Option {
	opt.Called = true
	opt.UsedAlias = usedAlias
	return opt
}

// SetBool - Set the option's data.
func (opt *Option) SetBool(b bool) *Option {
	opt.value = b
	*opt.pBool = b
	return opt
}

// SetBoolPtr - Set the option's data.
func (opt *Option) SetBoolPtr(b *bool) *Option {
	opt.value = *b
	opt.pBool = b
	return opt
}

// SetString - Set the option's data.
func (opt *Option) SetString(s string) *Option {
	opt.value = s
	*opt.pString = s
	return opt
}

// SetStringPtr - Set the option's data.
func (opt *Option) SetStringPtr(s *string) *Option {
	opt.value = *s
	opt.pString = s
	return opt
}

// SetInt - Set the option's data.
func (opt *Option) SetInt(i int) *Option {
	opt.value = i
	*opt.pInt = i
	return opt
}

// SetIntPtr - Set the option's data.
func (opt *Option) SetIntPtr(i *int) *Option {
	opt.value = *i
	opt.pInt = i
	return opt
}

// Int - Get the option's data.
// Exposed due to handle increment. Maybe there is a better way.
func (opt *Option) Int() int {
	return *opt.pInt
}

// SetFloat64 - Set the option's data.
func (opt *Option) SetFloat64(f float64) *Option {
	opt.value = f
	*opt.pFloat64 = f
	return opt
}

// SetFloat64Ptr - Set the option's data.
func (opt *Option) SetFloat64Ptr(f *float64) *Option {
	opt.value = *f
	opt.pFloat64 = f
	return opt
}

// SetStringSlice - Set the option's data.
func (opt *Option) SetStringSlice(s []string) *Option {
	opt.value = s
	*opt.pStringS = s
	return opt
}

// SetStringSlicePtr - Set the option's data.
func (opt *Option) SetStringSlicePtr(s *[]string) *Option {
	opt.value = *s
	opt.pStringS = s
	return opt
}

// SetIntSlice - Set the option's data.
func (opt *Option) SetIntSlice(s []int) *Option {
	opt.value = s
	*opt.pIntS = s
	return opt
}

// SetIntSlicePtr - Set the option's data.
func (opt *Option) SetIntSlicePtr(s *[]int) *Option {
	opt.value = *s
	opt.pIntS = s
	return opt
}

// SetStringMapPtr - Set the option's data.
func (opt *Option) SetStringMapPtr(m *map[string]string) *Option {
	opt.value = *m
	opt.pStringM = m
	return opt
}

// SetKeyValueToStringMap - Set the option's data.
func (opt *Option) SetKeyValueToStringMap(k, v string) *Option {
	if opt.MapKeysToLower {
		(*opt.pStringM)[strings.ToLower(k)] = v
	} else {
		(*opt.pStringM)[k] = v
	}
	opt.value = *opt.pStringM
	return opt
}

// Save - Saves the data provided into the option
func (opt *Option) Save(a ...string) error {
	Debug.Printf("optType: %d\n", opt.OptType)
	switch opt.OptType {
	case StringType:
		opt.SetString(a[0])
		return nil
	case IntType:
		i, err := strconv.Atoi(a[0])
		if err != nil {
			return fmt.Errorf(text.ErrorConvertToInt, opt.UsedAlias, a[0])
		}
		opt.SetInt(i)
		return nil
	case Float64Type:
		// TODO: Read the different errors when parsing float
		i, err := strconv.ParseFloat(a[0], 64)
		if err != nil {
			return fmt.Errorf(text.ErrorConvertToFloat64, opt.UsedAlias, a[0])
		}
		opt.SetFloat64(i)
		return nil
	case StringRepeatType:
		opt.SetStringSlice(append(*opt.pStringS, a...))
		return nil
	case IntRepeatType:
		var is []int
		for _, e := range a {
			if strings.Contains(e, "..") {
				Debug.Printf("e: %s\n", e)
				n := strings.SplitN(e, "..", 2)
				Debug.Printf("n: %v\n", n)
				n1, n2 := n[0], n[1]
				in1, err := strconv.Atoi(n1)
				if err != nil {
					// TODO: Create new error description for this error.
					return fmt.Errorf(text.ErrorConvertToInt, opt.UsedAlias, e)
				}
				in2, err := strconv.Atoi(n2)
				if err != nil {
					// TODO: Create new error description for this error.
					return fmt.Errorf(text.ErrorConvertToInt, opt.UsedAlias, e)
				}
				if in1 < in2 {
					for j := in1; j <= in2; j++ {
						is = append(is, j)
					}
				} else {
					// TODO: Create new error description for this error.
					return fmt.Errorf(text.ErrorConvertToInt, opt.UsedAlias, e)
				}
			} else {
				i, err := strconv.Atoi(e)
				if err != nil {
					return fmt.Errorf(text.ErrorConvertToInt, opt.UsedAlias, e)
				}
				is = append(is, i)
			}
		}
		opt.SetIntSlice(append(*opt.pIntS, is...))
		return nil
	case StringMapType:
		keyValue := strings.Split(a[0], "=")
		if len(keyValue) < 2 {
			return fmt.Errorf(text.ErrorArgumentIsNotKeyValue, opt.UsedAlias)
		}
		opt.SetKeyValueToStringMap(keyValue[0], keyValue[1])
		return nil
	default: // BoolType
		opt.SetBool(!*opt.pBool)
		return nil
	}
}

// Sort Interface
func Sort(list []*Option) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
}
