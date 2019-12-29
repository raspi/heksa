// This file is part of go-getoptions.
//
// Copyright (C) 2015-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package text - User facing strings
package text

// ErrorMissingArgument holds the text for missing argument error.
// It has a string placeholder '%s' for the name of the option missing the argument.
var ErrorMissingArgument = "Missing argument for option '%s'!"

// ErrorAmbiguousArgument holds the text for ambiguous argument error.
// It has a string placeholder '%s' for the passed option and a []string list of matches.
var ErrorAmbiguousArgument = "Ambiguous option '%s', matches %v!"

// ErrorMissingRequiredOption holds the text for missing required option error.
// It has a string placeholder '%s' for the name of the missing option.
var ErrorMissingRequiredOption = "Missing required option '%s'!"

// ErrorArgumentIsNotKeyValue holds the text for Map type options where the argument is not of key=value type.
// It has a string placeholder '%s' for the name of the option missing the argument.
var ErrorArgumentIsNotKeyValue = "Argument error for option '%s': Should be of type 'key=value'!"

// ErrorArgumentWithDash holds the text for missing argument error in cases where the next argument looks like an option (starts with '-').
// It has a string placeholder '%s' for the name of the option missing the argument.
var ErrorArgumentWithDash = "Missing argument for option '%s'!\n" +
	"If passing arguments that start with '-' use --option=-argument"

// ErrorConvertToInt holds the text for Int Coversion argument error.
// It has two string placeholders ('%s'). The first one for the name of the option with the wrong argument and the second one for the argument that could not be converted.
var ErrorConvertToInt = "Argument error for option '%s': Can't convert string to int: '%s'"

// ErrorConvertToFloat64 holds the text for Float64 Coversion argument error.
// It has two string placeholders ('%s'). The first one for the name of the option with the wrong argument and the second one for the argument that could not be converted.
var ErrorConvertToFloat64 = "Argument error for option '%s': Can't convert string to float64: '%s'"

// MessageOnUnknown holds the text for the unknown option message.
// It has a string placeholder '%s' for the name of the option missing the argument.
var MessageOnUnknown = "Unknown option '%s'"

// MessageOnInterrupt holds the text for the message to be printed when an interrupt is received.
var MessageOnInterrupt = "Interrupt signal received"

// HelpNameHeader holds the header text for the command name
var HelpNameHeader = "NAME"

// HelpSynopsisHeader holds the header text for the synopsis
var HelpSynopsisHeader = "SYNOPSIS"

// HelpCommandsHeader holds the header text for the command list
var HelpCommandsHeader = "COMMANDS"

// HelpRequiredOptionsHeader holds the header text for the required parameters
var HelpRequiredOptionsHeader = "REQUIRED PARAMETERS"

// HelpOptionsHeader holds the header text for the option list
var HelpOptionsHeader = "OPTIONS"
