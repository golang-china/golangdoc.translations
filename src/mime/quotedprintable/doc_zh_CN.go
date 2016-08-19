// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package quotedprintable implements quoted-printable encoding as specified by
// RFC 2045.

// Package quotedprintable implements quoted-printable encoding as specified by
// RFC 2045.
package quotedprintable

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
)

// Reader is a quoted-printable decoder.
type Reader struct {
	br   *bufio.Reader
	rerr error  // last read error
	line []byte // to be consumed before more of br
}


// A Writer is a quoted-printable writer that implements io.WriteCloser.
type Writer struct {
	// Binary mode treats the writer's input as pure binary and processes end of
	// line bytes as binary data.
	Binary bool

	w    io.Writer
	i    int
	line [78]byte
	cr   bool
}


// NewReader returns a quoted-printable reader, decoding from r.
func NewReader(r io.Reader) *Reader

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer

// Read reads and decodes quoted-printable data from the underlying reader.
func (*Reader) Read(p []byte) (n int, err error)

// Close closes the Writer, flushing any unwritten data to the underlying
// io.Writer, but does not close the underlying io.Writer.
func (*Writer) Close() error

// Write encodes p using quoted-printable encoding and writes it to the
// underlying io.Writer. It limits line length to 76 characters. The encoded
// bytes are not necessarily flushed until the Writer is closed.
func (*Writer) Write(p []byte) (n int, err error)

