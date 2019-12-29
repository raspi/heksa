// This file is part of go-getoptions.
//
// Copyright (C) 2015-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package completion - provides a Tree structure that can be used to define a program's completions.

Example Tree:

  mygit
  ├ log
  │ ├ sublog
  │ │ ├ --help
  │ │ ├ <file-completion>
  │ │ └ <custom-completion (sha1 list)>
  │ ├ --help
  │ └ <file-completion>
  ├ show
  │ ├ --help
  │ ├ --dir=<dir-completion>
  │ └ <file-completion>
  ├ --help
  └ --version

A tree node can have children and leaves.
The children are commands, the leaves can be options, file completions, custom completions and options that trigger custom file completions (--dir=<dir-comletion>).

Completions have a hierachy, commands are shown before file completions, and options are only shown if `-` is passed as part of the COMPLINE.

For custom completions a full list of completions must be passed as leaves to the node.
However, there file and dir completions are provided as a convenience.

Custom completions for options are triggered with the `=` sing after the full option test has been provided.

*/
package completion
