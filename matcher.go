/*
 * This file is part of Cashbook, a tool to analyze and report on sets of financial transactions.
 *
 * Copyright (C) 2014  Sourdough Labs Research and Development Corp.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"strings"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"regexp"
)


// Regexp to normalize incoming strings to match against.
var normalizer *regexp.Regexp = regexp.MustCompile("[^a-z A-Z]")

// Matchers look at an input string and return true if they match.  
// Matchers are expected to have no state
type Matcher interface {
	// Return true if this matcher matches the input string
	Match(input string) bool

	// Return a string suitable for use as a label (typically by removing 
	// the MatchThis string from the input)
	Label(input string) string

	GetGroupType() string
}

// A collection of Matchers
type MatchSet []Matcher

// Concrete matcher that compares a prefix string
type StartsWithMatcher struct {
	GroupType string
	MatchThis string
	GroupLabel string
	DontStripNumbers bool
}

func (matcher StartsWithMatcher) Match(input string) bool {
	return strings.HasPrefix(normalizer.ReplaceAllString(input, ""), matcher.MatchThis)
}

func (matcher StartsWithMatcher) GetGroupType() string {
	return matcher.GroupType
}

func (matcher StartsWithMatcher) Label(input string) string {	
	label := input

	if matcher.GroupLabel != "" {
		label = matcher.GroupLabel
	} else {
		if !matcher.DontStripNumbers {
			// Removes all non-alpha chars
			label = normalizer.ReplaceAllString(input, "") 			
		}

		label = strings.TrimSpace(strings.TrimPrefix(label, matcher.MatchThis))
		
		if label == "" {
			label = matcher.MatchThis
		}
	}

	return label
}

// Return the Matcher matching input, if the caller wishes a catch all matcher, they should
// make sure a NullMatcher is appended to the MatchSet
func (matchers MatchSet) FindMatcher(input string) (found_matcher Matcher) {

	// Iterate the set of matchers, calling match on each one.  
	// If true, break and return the found matcher
	for _, matcher := range matchers {
		if matcher.Match(input) {
			found_matcher = matcher
			break;
		}
	}

	return
}


// These are loaded from json and a MatchSet is assembled from them
// Doing it this way for now, but can do a custom marshaller/demarshaller
type MatcherDef struct {
	TypeName string
	GroupType string
	Matching string
	GroupLabel string
	DontStripNumbers bool
}

func match_set_from_file(filename string) MatchSet {

	file, e := ioutil.ReadFile(filename)
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }

	var profile []MatcherDef
	err := json.Unmarshal(file, &profile)
	
	if err != nil {
		log.Printf("Error reading json %v", err)
	}

	var set MatchSet = make(MatchSet, len(profile))

	for i, def := range profile {
		switch def.TypeName {
		case "StartsWithMatcher":
			set[i]  = &StartsWithMatcher{GroupType: def.GroupType, MatchThis: def.Matching, GroupLabel: def.GroupLabel, DontStripNumbers: def.DontStripNumbers}
		}
	}

	return set
}













