// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package zlib implements reading and writing of zlib format compressed data,
// as specified in RFC 1950.
//
// The implementation provides filters that uncompress during reading
// and compress during writing.  For example, to write compressed data
// to a buffer:
//
//     var b bytes.Buffer
//     w := zlib.NewWriter(&b)
//     w.Write([]byte("hello, world\n"))
//     w.Close()
//
// and to read that data back:
//
//     r, err := zlib.NewReader(&b)
//     io.Copy(os.Stdout, r)
//     r.Close()

// Package zlib implements reading and writing of zlib format compressed data,
// as specified in RFC 1950.
//
// The implementation provides filters that uncompress during reading
// and compress during writing.  For example, to write compressed data
// to a buffer:
//
//     var b bytes.Buffer
//     w := zlib.NewWriter(&b)
//     w.Write([]byte("hello, world\n"))
//     w.Close()
//
// and to read that data back:
//
//     r, err := zlib.NewReader(&b)
//     io.Copy(os.Stdout, r)
//     r.Close()
package zlib

import (
    "bufio"
    "compress/flate"
    "errors"
    "fmt"
    "hash"
    "hash/adler32"
    "io"
)

// These constants are copied from the flate package, so that code that imports
// "compress/zlib" does not also have to import "compress/flate".
const (
	NoCompression      = flate.NoCompression
	BestSpeed          = flate.BestSpeed
	BestCompression    = flate.BestCompression
	DefaultCompression = flate.DefaultCompression
)



var (
	// ErrChecksum is returned when reading ZLIB data that has an invalid
	// checksum.
	ErrChecksum = errors.New("zlib: invalid checksum")
	// ErrDictionary is returned when reading ZLIB data that has an invalid
	// dictionary.
	ErrDictionary = errors.New("zlib: invalid dictionary")
	// ErrHeader is returned when reading ZLIB data that has an invalid header.
	ErrHeader = errors.New("zlib: invalid header")
)


// Resetter resets a ReadCloser returned by NewReader or NewReaderDict to
// to switch to a new underlying Reader. This permits reusing a ReadCloser
// instead of allocating a new one.
type Resetter interface {
	// Reset discards any buffered data and resets the Resetter as if it was
	// newly initialized with the given reader.
	Reset(r io.Reader, dict []byte) error
}


// A Writer takes data written to it and writes the compressed
// form of that data to an underlying writer (see NewWriter).
type Writer struct {
	w           io.Writer
	level       int
	dict        []byte
	compressor  *flate.Writer
	digest      hash.Hash32
	err         error
	scratch     [4]byte
	wroteHeader bool
}


// NewReader creates a new ReadCloser. Reads from the returned ReadCloser read
// and decompress data from r. The implementation buffers input and may read
// more data than necessary from r. It is the caller's responsibility to call
// Close on the ReadCloser when done.
//
// The ReadCloser returned by NewReader also implements Resetter.

// NewReader creates a new ReadCloser.
// Reads from the returned ReadCloser read and decompress data from r.
// If r does not implement io.ByteReader, the decompressor may read more
// data than necessary from r.
// It is the caller's responsibility to call Close on the ReadCloser when done.
//
// The ReadCloser returned by NewReader also implements Resetter.
func NewReader(r io.Reader) (io.ReadCloser, error)

// NewReaderDict is like NewReader but uses a preset dictionary. NewReaderDict
// ignores the dictionary if the compressed data does not refer to it. If the
// compressed data refers to a different dictionary, NewReaderDict returns
// ErrDictionary.
//
// The ReadCloser returned by NewReaderDict also implements Resetter.
func NewReaderDict(r io.Reader, dict []byte) (io.ReadCloser, error)

// NewWriter creates a new Writer.
// Writes to the returned Writer are compressed and written to w.
//
// It is the caller's responsibility to call Close on the WriteCloser when done.
// Writes may be buffered and not flushed until Close.
func NewWriter(w io.Writer) *Writer

// NewWriterLevel is like NewWriter but specifies the compression level instead
// of assuming DefaultCompression.
//
// The compression level can be DefaultCompression, NoCompression, or any
// integer value between BestSpeed and BestCompression inclusive. The error
// returned will be nil if the level is valid.
func NewWriterLevel(w io.Writer, level int) (*Writer, error)

// NewWriterLevelDict is like NewWriterLevel but specifies a dictionary to
// compress with.
//
// The dictionary may be nil. If not, its contents should not be modified until
// the Writer is closed.
func NewWriterLevelDict(w io.Writer, level int, dict []byte) (*Writer, error)

// Close closes the Writer, flushing any unwritten data to the underlying
// io.Writer, but does not close the underlying io.Writer.
func (*Writer) Close() error

// Flush flushes the Writer to its underlying io.Writer.
func (*Writer) Flush() error

// Reset clears the state of the Writer z such that it is equivalent to its
// initial state from NewWriterLevel or NewWriterLevelDict, but instead writing
// to w.
func (*Writer) Reset(w io.Writer)

// Write writes a compressed form of p to the underlying io.Writer. The
// compressed bytes are not necessarily flushed until the Writer is closed or
// explicitly flushed.
func (*Writer) Write(p []byte) (n int, err error)

