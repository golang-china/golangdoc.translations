// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package bzip2 implements bzip2 decompression.

// Package bzip2 implements bzip2 decompression.
package bzip2

import (
    "bufio"
    "io"
    "sort"
)

// A StructuralError is returned when the bzip2 data is found to be
// syntactically invalid.
type StructuralError string


// NewReader returns an io.Reader which decompresses bzip2 data from r.
// If r does not also implement io.ByteReader,
// the decompressor may read more data than necessary from r.
func NewReader(r io.Reader) io.Reader

func (StructuralError) Error() string

