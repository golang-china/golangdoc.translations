// +build ingore

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
}

// A Writer is a quoted-printable writer that implements io.WriteCloser.
type Writer struct {
	// Binary mode treats the writer's input as pure binary and processes end of
	// line bytes as binary data.
	Binary bool
}

// NewReader returns a quoted-printable reader, decoding from r.
func NewReader(r io.Reader) *Reader

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer

// Read reads and decodes quoted-printable data from the underlying reader.
func (r *Reader) Read(p []byte) (n int, err error)

// Close closes the Writer, flushing any unwritten data to the underlying
// io.Writer, but does not close the underlying io.Writer.
func (w *Writer) Close() error

// Write encodes p using quoted-printable encoding and writes it to the
// underlying io.Writer. It limits line length to 76 characters. The encoded
// bytes are not necessarily flushed until the Writer is closed.
func (w *Writer) Write(p []byte) (n int, err error)

