// This file is part of go-getoptions.
//
// Copyright (C) 2015-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package getoptions - Go option parser inspired on the flexibility of Perl’s
GetOpt::Long.

It will operate on any given slice of strings and return the remaining (non
used) command line arguments. This allows to easily subcommand.


Usage

The following is a basic example:

	package main

	import (
		"fmt"
		"io/ioutil"
		"log"
		"os"

		"github.com/DavidGamba/go-getoptions"
	)

	var logger = log.New(os.Stderr, "DEBUG: ", log.LstdFlags)

	func main() {
		// Declare the variables you want your options to update
		var debug bool
		var greetCount int

		// Declare the GetOptions object
		opt := getoptions.New()

		// Options definition
		opt.Bool("help", false, opt.Alias("h", "?")) // Aliases can be defined
		opt.BoolVar(&debug, "debug", false)
		opt.IntVar(&greetCount, "greet", 0,
			opt.Required(), // Mark option as required
			opt.Description("Number of times to greet."), // Set the automated help description
			opt.ArgName("number"),                        // Change the help synopsis arg from <int> to <number>
		)
		greetings := opt.StringMap("list", 1, 99,
			opt.Description("Greeting list by language."),
			opt.ArgName("lang=msg"), // Change the help synopsis arg from <key=value> to <lang=msg>
		)

		// Parse cmdline arguments or any provided []string
		remaining, err := opt.Parse(os.Args[1:])

		// Handle help before handling user errors
		if opt.Called("help") {
			fmt.Fprintf(os.Stderr, opt.Help())
			os.Exit(1)
		}

		// Handle user errors
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
			fmt.Fprintf(os.Stderr, opt.HelpSynopsis())
			os.Exit(1)
		}
		if !debug {
			logger.SetOutput(ioutil.Discard)
		}
		logger.Printf("Remaining: %v\n", remaining)

		for i := 0; i < greetCount; i++ {
			fmt.Println("Hello World, from go-getoptions!")
		}
		if len(greetings) > 0 {
			fmt.Printf("Greeting List:\n")
			for k, v := range greetings {
				fmt.Printf("\t%s=%s\n", k, v)
			}
		}
	}

Features

• Allow passing options and non-options in any order.

• Support for `--long` options.

• Support for short (`-s`) options with flexible behaviour (see https://github.com/DavidGamba/go-getoptions#operation_modes for details):

 - Normal (default)
 - Bundling
 - SingleDash

• `Called()` method indicates if the option was passed on the command line.

• Multiple aliases for the same option. e.g. `help`, `man`.

• `CalledAs()` method indicates what alias was used to call the option on the command line.

• Simple synopsis and option list automated help.

• Boolean, String, Int and Float64 type options.

• Negatable Boolean options.
For example: `--verbose`, `--no-verbose` or `--noverbose`.

• Options with Array arguments.
The same option can be used multiple times with different arguments.
The list of arguments will be saved into an Array like structure inside the program.

• Options with array arguments and multiple entries.
For example: `color --rgb 10 20 30 --next-option`

• When using integer array options with multiple arguments, positive integer ranges are allowed.
For example: `1..3` to indicate `1 2 3`.

• Options with key value arguments and multiple entries.

• Options with Key Value arguments.
This allows the same option to be used multiple times with arguments of key value type.
For example: `rpmbuild --define name=myrpm --define version=123`.

• Supports passing `--` to stop parsing arguments (everything after will be left in the `remaining []string`).

• Supports command line options with '='.
For example: You can use `--string=mystring` and `--string mystring`.

• Allows passing arguments to options that start with dash `-` when passed after equal.
For example: `--string=--hello` and `--int=-123`.

• Options with optional arguments.
If the default argument is not passed the default is set.
For example: You can call `--int 123` which yields `123` or `--int` which yields the given default.

• Allows abbreviations when the provided option is not ambiguous.
For example: An option called `build` can be called with `--b`, `--bu`, `--bui`, `--buil` and `--build` as long as there is no ambiguity.
In the case of ambiguity, the shortest non ambiguous combination is required.

• Support for the lonesome dash "-".
To indicate, for example, when to read input from STDIO.

• Incremental options.
Allows the same option to be called multiple times to increment a counter.

• Supports case sensitive options.
For example, you can use `v` to define `verbose` and `V` to define `Version`.

• Support indicating if an option is required and allows overriding default error message.

• Errors exposed as public variables to allow overriding them for internationalization.

• Supports subcommands (stop parsing arguments when non option is passed).

• Multiple ways of managing unknown options:
  - Fail on unknown (default).
  - Warn on unknown.
  - Pass through, allows for subcommands and can be combined with Require Order.

• Require order: Allows for subcommands. Stop parsing arguments when the first non-option is found.
When mixed with Pass through, it also stops parsing arguments when the first unmatched option is found.

Panic

The library will panic if it finds that the programmer (not end user):

• Defined the same alias twice.

• Defined wrong min and max values for SliceMulti methods.
*/
package getoptions
