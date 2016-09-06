// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package profile provides a representation of profile.proto and
// methods to encode/decode profiles in this format.

// Implements methods to filter samples from profiles.
//
// Package profile provides a representation of profile.proto and methods to
// encode/decode profiles in this format.
package profile

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	// LegacyHeapAllocated instructs the heapz parsers to use the
	// allocated memory stats instead of the default in-use memory. Note
	// that tcmalloc doesn't provide all allocated memory, only in-use
	// stats.
	LegacyHeapAllocated bool
)

// Demangler maps symbol names to a human-readable form. This may
// include C++ demangling and additional simplification. Names that
// are not demangled may be missing from the resulting map.

// Demangler maps symbol names to a human-readable form. This may include C++
// demangling and additional simplification. Names that are not demangled may be
// missing from the resulting map.
type Demangler func(name []string) (map[string]string, error)

// Function corresponds to Profile.Function
type Function struct {
	ID         uint64
	Name       string
	SystemName string
	Filename   string
	StartLine  int64
}

// Label corresponds to Profile.Label
type Label struct {
}

// Line corresponds to Profile.Line
type Line struct {
	Function *Function
	Line     int64
}

// Location corresponds to Profile.Location
type Location struct {
	ID      uint64
	Mapping *Mapping
	Address uint64
	Line    []Line
}

// Mapping corresponds to Profile.Mapping
type Mapping struct {
	ID              uint64
	Start           uint64
	Limit           uint64
	Offset          uint64
	File            string
	BuildID         string
	HasFunctions    bool
	HasFilenames    bool
	HasLineNumbers  bool
	HasInlineFrames bool
}

// Profile is an in-memory representation of profile.proto.
type Profile struct {
	SampleType    []*ValueType
	Sample        []*Sample
	Mapping       []*Mapping
	Location      []*Location
	Function      []*Function
	DropFrames    string
	KeepFrames    string
	TimeNanos     int64
	DurationNanos int64
	PeriodType    *ValueType
	Period        int64
}

// Sample corresponds to Profile.Sample
type Sample struct {
	Location []*Location
	Value    []int64
	Label    map[string][]string
	NumLabel map[string][]int64
}

// TagMatch selects tags for filtering
type TagMatch func(key, val string, nval int64) bool

// ValueType corresponds to Profile.ValueType
type ValueType struct {
	Type string // cpu, wall, inuse_space, etc
	Unit string // seconds, nanoseconds, bytes, etc
}

// Parse parses a profile and checks for its validity. The input
// may be a gzip-compressed encoded protobuf or one of many legacy
// profile formats which may be unsupported in the future.

// Parse parses a profile and checks for its validity. The input may be a
// gzip-compressed encoded protobuf or one of many legacy profile formats which
// may be unsupported in the future.
func Parse(r io.Reader) (*Profile, error)

// ParseTracebacks parses a set of tracebacks and returns a newly
// populated profile. It will accept any text file and generate a
// Profile out of it with any hex addresses it can identify, including
// a process map if it can recognize one. Each sample will include a
// tag "source" with the addresses recognized in string format.

// ParseTracebacks parses a set of tracebacks and returns a newly populated
// profile. It will accept any text file and generate a Profile out of it with
// any hex addresses it can identify, including a process map if it can
// recognize one. Each sample will include a tag "source" with the addresses
// recognized in string format.
func ParseTracebacks(b []byte) (*Profile, error)

// Aggregate merges the locations in the profile into equivalence
// classes preserving the request attributes. It also updates the
// samples to point to the merged locations.

// Aggregate merges the locations in the profile into equivalence classes
// preserving the request attributes. It also updates the samples to point to
// the merged locations.
func (p *Profile) Aggregate(inlineFrame, function, filename, linenumber, address bool) error

// CheckValid tests whether the profile is valid. Checks include, but are
// not limited to:
//   - len(Profile.Sample[n].value) == len(Profile.value_unit)
//   - Sample.id has a corresponding Profile.Location

// CheckValid tests whether the profile is valid. Checks include, but are not
// limited to:
//
//     - len(Profile.Sample[n].value) == len(Profile.value_unit)
//     - Sample.id has a corresponding Profile.Location
func (p *Profile) CheckValid() error

// Compatible determines if two profiles can be compared/merged.
// returns nil if the profiles are compatible; otherwise an error with
// details on the incompatibility.

// Compatible determines if two profiles can be compared/merged. returns nil if
// the profiles are compatible; otherwise an error with details on the
// incompatibility.
func (p *Profile) Compatible(pb *Profile) error

// Copy makes a fully independent copy of a profile.
func (p *Profile) Copy() *Profile

// Demangle attempts to demangle and optionally simplify any function
// names referenced in the profile. It works on a best-effort basis:
// it will silently preserve the original names in case of any errors.

// Demangle attempts to demangle and optionally simplify any function names
// referenced in the profile. It works on a best-effort basis: it will silently
// preserve the original names in case of any errors.
func (p *Profile) Demangle(d Demangler) error

// Empty returns true if the profile contains no samples.
func (p *Profile) Empty() bool

// FilterSamplesByName filters the samples in a profile and only keeps
// samples where at least one frame matches focus but none match ignore.
// Returns true is the corresponding regexp matched at least one sample.

// FilterSamplesByName filters the samples in a profile and only keeps samples
// where at least one frame matches focus but none match ignore. Returns true is
// the corresponding regexp matched at least one sample.
func (p *Profile) FilterSamplesByName(focus, ignore, hide *regexp.Regexp) (fm, im, hm bool)

// FilterSamplesByTag removes all samples from the profile, except
// those that match focus and do not match the ignore regular
// expression.

// FilterSamplesByTag removes all samples from the profile, except those that
// match focus and do not match the ignore regular expression.
func (p *Profile) FilterSamplesByTag(focus, ignore TagMatch) (fm, im bool)

// HasFileLines determines if all locations in this profile have
// symbolized file and line number information.

// HasFileLines determines if all locations in this profile have symbolized file
// and line number information.
func (p *Profile) HasFileLines() bool

// HasFunctions determines if all locations in this profile have
// symbolized function information.

// HasFunctions determines if all locations in this profile have symbolized
// function information.
func (p *Profile) HasFunctions() bool

// Merge adds profile p adjusted by ratio r into profile p. Profiles
// must be compatible (same Type and SampleType).
// TODO(rsilvera): consider normalizing the profiles based on the
// total samples collected.

// Merge adds profile p adjusted by ratio r into profile p. Profiles must be
// compatible (same Type and SampleType). TODO(rsilvera): consider normalizing
// the profiles based on the total samples collected.
func (p *Profile) Merge(pb *Profile, r float64) error

// ParseMemoryMap parses a memory map in the format of
// /proc/self/maps, and overrides the mappings in the current profile.
// It renumbers the samples and locations in the profile correspondingly.

// ParseMemoryMap parses a memory map in the format of /proc/self/maps, and
// overrides the mappings in the current profile. It renumbers the samples and
// locations in the profile correspondingly.
func (p *Profile) ParseMemoryMap(rd io.Reader) error

// Prune removes all nodes beneath a node matching dropRx, and not
// matching keepRx. If the root node of a Sample matches, the sample
// will have an empty stack.

// Prune removes all nodes beneath a node matching dropRx, and not matching
// keepRx. If the root node of a Sample matches, the sample will have an empty
// stack.
func (p *Profile) Prune(dropRx, keepRx *regexp.Regexp)

// RemoveUninteresting prunes and elides profiles using built-in
// tables of uninteresting function names.

// RemoveUninteresting prunes and elides profiles using built-in tables of
// uninteresting function names.
func (p *Profile) RemoveUninteresting() error

// Print dumps a text representation of a profile. Intended mainly
// for debugging purposes.

// Print dumps a text representation of a profile. Intended mainly for debugging
// purposes.
func (p *Profile) String() string

// Write writes the profile as a gzip-compressed marshaled protobuf.
func (p *Profile) Write(w io.Writer) error

