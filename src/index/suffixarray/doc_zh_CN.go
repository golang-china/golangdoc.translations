// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package suffixarray implements substring search in logarithmic time using an
// in-memory suffix array.
//
// Example use:
//
//     // create index for some data
//     index := suffixarray.New(data)
//
//     // lookup byte slice s
//     offsets1 := index.Lookup(s, -1) // the list of all indices where s occurs in data
//     offsets2 := index.Lookup(s, 3)  // the list of at most 3 indices where s occurs in data

// Package suffixarray implements substring search in logarithmic time using an
// in-memory suffix array.
//
// Example use:
//
//     // create index for some data
//     index := suffixarray.New(data)
//
//     // lookup byte slice s
//     offsets1 := index.Lookup(s, -1) // the list of all indices where s occurs in data
//     offsets2 := index.Lookup(s, 3)  // the list of at most 3 indices where s occurs in data
package suffixarray

import (
    "bytes"
    "encoding/binary"
    "io"
    "regexp"
    "sort"
)

// Index implements a suffix array for fast substring search.
type Index struct {
	data []byte
	sa   []int // suffix array for data; len(sa) == len(data)
}


// New creates a new Index for data.
// Index creation time is O(N*log(N)) for N = len(data).
func New(data []byte) *Index

// Bytes returns the data over which the index was created.
// It must not be modified.
func (*Index) Bytes() []byte

// FindAllIndex returns a sorted list of non-overlapping matches of the
// regular expression r, where a match is a pair of indices specifying
// the matched slice of x.Bytes(). If n < 0, all matches are returned
// in successive order. Otherwise, at most n matches are returned and
// they may not be successive. The result is nil if there are no matches,
// or if n == 0.
func (*Index) FindAllIndex(r *regexp.Regexp, n int) (result [][]int)

// Lookup returns an unsorted list of at most n indices where the byte string s
// occurs in the indexed data. If n < 0, all occurrences are returned.
// The result is nil if s is empty, s is not found, or n == 0.
// Lookup time is O(log(N)*len(s) + len(result)) where N is the
// size of the indexed data.
func (*Index) Lookup(s []byte, n int) (result []int)

// Read reads the index from r into x; x must not be nil.
func (*Index) Read(r io.Reader) error

// Write writes the index x to w.
func (*Index) Write(w io.Writer) error

