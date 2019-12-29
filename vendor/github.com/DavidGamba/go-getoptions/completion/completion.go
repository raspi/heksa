// This file is part of go-getoptions.
//
// Copyright (C) 2015-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package completion

import (
	"io/ioutil"
	"log"
	"strings"
)

// Debug - Debug logger set to ioutil.Discard by default
var Debug = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

/*
Node -

Example:

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
*/
type Node struct {
	Name     string // Name of the node. For StringNode Kinds, this is the completion.
	Kind     kind   // Kind of node.
	Children []*Node
	Entries  []string // Use as completions for OptionsNode and CustomNode Kind.
	// TODO: Maybe add sibling completion that gets activated with = for options
}

// CompletionType -
type kind int

const (
	// Root -
	Root kind = iota
	// StringNode -
	StringNode
	// FileListNode - Regular file completion you would expect.
	// Name used as the dir to start completing results from.
	// TODO: Allow ignore case.
	FileListNode
	// OptionsNode - Only enabled if prefix starts with -
	OptionsNode
	// OptionsWithCompletion - The completion completes to --option=
	OptionsWithCompletion
	// CustomNode -
	CustomNode
)

// NewNode -
func NewNode(name string, kind kind, entries []string) *Node {
	if entries == nil {
		entries = []string{}
	}
	return &Node{
		Name:    name,
		Kind:    kind,
		Entries: entries,
	}
}

// AddChild -
// TODO: Probably make sure that the name is not already in use since we find them by name.
func (n *Node) AddChild(node *Node) {
	n.Children = append(n.Children, node)
}

// SelfCompletions -
func (n *Node) SelfCompletions(prefix string) []string {
	switch n.Kind {
	case StringNode:
		if strings.HasPrefix(n.Name, prefix) {
			Debug.Printf("SelfCompletions - node: %s > %v\n", n.Name, []string{n.Name})
			return []string{n.Name}
		}
	case FileListNode:
		files, _ := listDir(n.Name, prefix)
		if strings.HasPrefix(prefix, ".") {
			Debug.Printf("SelfCompletions - node: %s > %v\n", n.Name, files)
			return files
		}
		// Don't return hidden files unless requested by the prefix
		ff := discardByPrefix(files, ".")
		Debug.Printf("SelfCompletions - node: %s > %v\n", n.Name, ff)
		return ff
	case OptionsNode:
		if strings.HasPrefix(prefix, "-") {
			sortForCompletion(n.Entries)
			ee := keepByPrefix(n.Entries, prefix)
			Debug.Printf("SelfCompletions - node: %s > %v\n", n.Name, ee)
			return ee
		}
	case CustomNode:
		sortForCompletion(n.Entries)
		ee := keepByPrefix(n.Entries, prefix)
		Debug.Printf("SelfCompletions - node: %s > %v\n", n.Name, ee)
		return ee
	}
	Debug.Printf("SelfCompletions - node: %s > %v\n", n.Name, []string{})
	return []string{}
}

// Completions -
func (n *Node) Completions(prefix string) []string {
	results := []string{}
	for _, child := range n.Children {
		results = append(results, child.SelfCompletions(prefix)...)
	}
	Debug.Printf("Completions - node: %s, prefix %s > %v\n", n.Name, prefix, results)
	return results
}

// GetChildByName - Traverses to the children and returns the first one to match name.
func (n *Node) GetChildByName(name string) *Node {
	for _, child := range n.Children {
		if child.Name == name {
			return child
		}
	}
	return NewNode("", Root, []string{})
}

func (n *Node) GetChildrenByKind(kind kind) []*Node {
	children := []*Node{}
	for _, child := range n.Children {
		if child.Kind == kind {
			children = append(children, child)
		}
	}
	return children
}

// keepByPrefix - Given a list and a prefix filter, it returns a list subset of the elements that start with the prefix.
func keepByPrefix(list []string, prefix string) []string {
	keepList := []string{}
	for _, e := range list {
		if strings.HasPrefix(e, prefix) {
			keepList = append(keepList, e)
		}
	}
	return keepList
}

// discardByPrefix - Given a list and a prefix filter, it returns a list subset of the elements that Do not start with the prefix.
func discardByPrefix(list []string, prefix string) []string {
	keepList := []string{}
	for _, e := range list {
		if !strings.HasPrefix(e, prefix) {
			keepList = append(keepList, e)
		}
	}
	return keepList
}

// CompLineComplete - Given a compLine (get it with os.Getenv("COMP_LINE")) it returns a list of completions.
func (n *Node) CompLineComplete(compLine string) []string {
	// TODO: This split might not consider files that have spaces in them.
	compLineParts := strings.Split(compLine, " ")
	// return compLineParts
	if len(compLineParts) == 0 || compLineParts[0] == "" {
		Debug.Printf("CompLineComplete - node: %s, compLine %s > %v - Empty compLineParts\n", n.Name, compLine, []string{})
		return []string{}
	}

	// Drop the executable or command
	compLineParts = compLineParts[1:]

	// We have a possibly partial request
	if len(compLineParts) >= 1 {
		current := compLineParts[0]

		cc := n.Completions(current)
		if len(compLineParts) == 1 && len(cc) > 1 {
			Debug.Printf("CompLineComplete - node: %s, compLine %s > %v - Multiple completions for this compLine\n", n.Name, compLine, cc)
			return cc
		}
		// Check if the current fully matches a command (child node)
		child := n.GetChildByName(current)
		if child.Name == current && child.Kind == StringNode {
			Debug.Printf("CompLineComplete - node: %s, compLine %s - Recursing into command %s\n", n.Name, compLine, current)
			// Recurse into the child node's completion
			return child.CompLineComplete(strings.Join(compLineParts, " "))
		}
		// Check if the current fully matches an option
		list := n.GetChildrenByKind(OptionsNode)
		list = append(list, n.GetChildrenByKind(CustomNode)...)
		for _, child := range list {
			for _, e := range child.Entries {
				if current == e {
					if len(compLineParts) == 1 {
						Debug.Printf("CompLineComplete - node: %s, compLine %s > %v - Fully Matched Option/Custom\n", n.Name, compLine, current)
						return []string{current}
					}
					Debug.Printf("CompLineComplete - node: %s, compLine %s - Fully matched Option/Custom %s, recursing to self\n", n.Name, compLine, current)
					// Recurse into the node self completion
					return n.CompLineComplete(strings.Join(compLineParts, " "))
				}
			}
		}
		// Get FileList completions after all other completions
		for _, child := range n.GetChildrenByKind(FileListNode) {
			cc := child.SelfCompletions(current)
			for _, e := range cc {
				if current == e {
					if len(compLineParts) == 1 {
						Debug.Printf("CompLineComplete - node: %s, compLine %s > %v - Fully matched File\n", n.Name, compLine, current)
						return []string{current}
					}
					Debug.Printf("CompLineComplete - node: %s, compLine %s - Fully matched File %s, recursing to self\n", n.Name, compLine, current)
					// Recurse into the node self completion
					return n.CompLineComplete(strings.Join(compLineParts, " "))
				}
			}
		}

		// Return a partial match
		Debug.Printf("CompLineComplete - node: %s, compLine %s - Partial match %s\n", n.Name, compLine, current)
		return n.Completions(current)
	}

	Debug.Printf("CompLineComplete - node: %s, compLine %s > [] - Return all results\n", n.Name, compLine)
	// No partial request, return all results
	return n.Completions("")
}
