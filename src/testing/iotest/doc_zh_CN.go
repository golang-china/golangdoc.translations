// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package iotest implements Readers and Writers useful mainly for testing.
package iotest

import (
	"errors"
	"io"
	"log"
)

var ErrTimeout = errors.New("timeout")

// DataErrReader changes the way errors are handled by a Reader. Normally, a
// Reader returns an error (typically EOF) from the first Read call after the
// last piece of data is read. DataErrReader wraps a Reader and changes its
// behavior so the final error is returned along with the final data, instead of
// in the first call after the final data.
func DataErrReader(r io.Reader) io.Reader

// HalfReader returns a Reader that implements Read
// by reading half as many requested bytes from r.

// HalfReader returns a Reader that implements Read by reading half as many
// requested bytes from r.
func HalfReader(r io.Reader) io.Reader

// NewReadLogger returns a reader that behaves like r except
// that it logs (using log.Print) each read to standard error,
// printing the prefix and the hexadecimal data read.

// NewReadLogger returns a reader that behaves like r except that it logs (using
// log.Print) each read to standard error, printing the prefix and the
// hexadecimal data read.
func NewReadLogger(prefix string, r io.Reader) io.Reader

// NewWriteLogger returns a writer that behaves like w except
// that it logs (using log.Printf) each write to standard error,
// printing the prefix and the hexadecimal data written.

// NewWriteLogger returns a writer that behaves like w except that it logs
// (using log.Printf) each write to standard error, printing the prefix and the
// hexadecimal data written.
func NewWriteLogger(prefix string, w io.Writer) io.Writer

// OneByteReader returns a Reader that implements
// each non-empty Read by reading one byte from r.

// OneByteReader returns a Reader that implements each non-empty Read by reading
// one byte from r.
func OneByteReader(r io.Reader) io.Reader

// TimeoutReader returns ErrTimeout on the second read
// with no data. Subsequent calls to read succeed.

// TimeoutReader returns ErrTimeout on the second read with no data. Subsequent
// calls to read succeed.
func TimeoutReader(r io.Reader) io.Reader

// TruncateWriter returns a Writer that writes to w
// but stops silently after n bytes.

// TruncateWriter returns a Writer that writes to w but stops silently after n
// bytes.
func TruncateWriter(w io.Writer, n int64) io.Writer

