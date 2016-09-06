// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package bio implements common I/O abstractions used within the Go toolchain.
package bio

import (
	"bufio"
	"io"
	"log"
	"os"
)

// Reader implements a seekable buffered io.Reader.
type Reader struct {
	*bufio.Reader
}

// Writer implements a seekable buffered io.Writer.
type Writer struct {
	*bufio.Writer
}

// Create creates the file named name and returns a Writer
// for that file.
func Create(name string) (*Writer, error)

// MustClose closes Closer c and calls log.Fatal if it returns a non-nil error.
func MustClose(c io.Closer)

// MustWriter returns a Writer that wraps the provided Writer,
// except that it calls log.Fatal instead of returning a non-nil error.
func MustWriter(w io.Writer) io.Writer

// Open returns a Reader for the file named name.
func Open(name string) (*Reader, error)

func (r *Reader) Close() error

func (r *Reader) Offset() int64

func (r *Reader) Seek(offset int64, whence int) int64

func (w *Writer) Close() error

func (w *Writer) Offset() int64

func (w *Writer) Seek(offset int64, whence int) int64

