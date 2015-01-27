// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package gzip implements reading and writing of gzip format compressed files, as
// specified in RFC 1952.
package gzip

// These constants are copied from the flate package, so that code that imports
// "compress/gzip" does not also have to import "compress/flate".
const (
	NoCompression      = flate.NoCompression
	BestSpeed          = flate.BestSpeed
	BestCompression    = flate.BestCompression
	DefaultCompression = flate.DefaultCompression
)

var (
	// ErrChecksum is returned when reading GZIP data that has an invalid checksum.
	ErrChecksum = errors.New("gzip: invalid checksum")
	// ErrHeader is returned when reading GZIP data that has an invalid header.
	ErrHeader = errors.New("gzip: invalid header")
)

// The gzip file stores a header giving metadata about the compressed file. That
// header is exposed as the fields of the Writer and Reader structs.
type Header struct {
	Comment string    // comment
	Extra   []byte    // "extra data"
	ModTime time.Time // modification time
	Name    string    // file name
	OS      byte      // operating system type
}

// A Reader is an io.Reader that can be read to retrieve uncompressed data from a
// gzip-format compressed file.
//
// In general, a gzip file can be a concatenation of gzip files, each with its own
// header. Reads from the Reader return the concatenation of the uncompressed data
// of each. Only the first header is recorded in the Reader fields.
//
// Gzip files store a length and checksum of the uncompressed data. The Reader will
// return a ErrChecksum when Read reaches the end of the uncompressed data if it
// does not have the expected length or checksum. Clients should treat data
// returned by Read as tentative until they receive the io.EOF marking the end of
// the data.
type Reader struct {
	Header
	// contains filtered or unexported fields
}

// NewReader creates a new Reader reading the given reader. If r does not also
// implement io.ByteReader, the decompressor may read more data than necessary from
// r. It is the caller's responsibility to call Close on the Reader when done.
func NewReader(r io.Reader) (*Reader, error)

// Close closes the Reader. It does not close the underlying io.Reader.
func (z *Reader) Close() error

// Multistream controls whether the reader supports multistream files.
//
// If enabled (the default), the Reader expects the input to be a sequence of
// individually gzipped data streams, each with its own header and trailer, ending
// at EOF. The effect is that the concatenation of a sequence of gzipped files is
// treated as equivalent to the gzip of the concatenation of the sequence. This is
// standard behavior for gzip readers.
//
// Calling Multistream(false) disables this behavior; disabling the behavior can be
// useful when reading file formats that distinguish individual gzip data streams
// or mix gzip data streams with other data streams. In this mode, when the Reader
// reaches the end of the data stream, Read returns io.EOF. If the underlying
// reader implements io.ByteReader, it will be left positioned just after the gzip
// stream. To start the next stream, call z.Reset(r) followed by
// z.Multistream(false). If there is no next stream, z.Reset(r) will return io.EOF.
func (z *Reader) Multistream(ok bool)

func (z *Reader) Read(p []byte) (n int, err error)

// Reset discards the Reader z's state and makes it equivalent to the result of its
// original state from NewReader, but reading from r instead. This permits reusing
// a Reader rather than allocating a new one.
func (z *Reader) Reset(r io.Reader) error

// A Writer is an io.WriteCloser. Writes to a Writer are compressed and written to
// w.
type Writer struct {
	Header
	// contains filtered or unexported fields
}

// NewWriter returns a new Writer. Writes to the returned writer are compressed and
// written to w.
//
// It is the caller's responsibility to call Close on the WriteCloser when done.
// Writes may be buffered and not flushed until Close.
//
// Callers that wish to set the fields in Writer.Header must do so before the first
// call to Write or Close. The Comment and Name header fields are UTF-8 strings in
// Go, but the underlying format requires NUL-terminated ISO 8859-1 (Latin-1). NUL
// or non-Latin-1 runes in those strings will lead to an error on Write.
func NewWriter(w io.Writer) *Writer

// NewWriterLevel is like NewWriter but specifies the compression level instead of
// assuming DefaultCompression.
//
// The compression level can be DefaultCompression, NoCompression, or any integer
// value between BestSpeed and BestCompression inclusive. The error returned will
// be nil if the level is valid.
func NewWriterLevel(w io.Writer, level int) (*Writer, error)

// Close closes the Writer, flushing any unwritten data to the underlying
// io.Writer, but does not close the underlying io.Writer.
func (z *Writer) Close() error

// Flush flushes any pending compressed data to the underlying writer.
//
// It is useful mainly in compressed network protocols, to ensure that a remote
// reader has enough data to reconstruct a packet. Flush does not return until the
// data has been written. If the underlying writer returns an error, Flush returns
// that error.
//
// In the terminology of the zlib library, Flush is equivalent to Z_SYNC_FLUSH.
func (z *Writer) Flush() error

// Reset discards the Writer z's state and makes it equivalent to the result of its
// original state from NewWriter or NewWriterLevel, but writing to w instead. This
// permits reusing a Writer rather than allocating a new one.
func (z *Writer) Reset(w io.Writer)

// Write writes a compressed form of p to the underlying io.Writer. The compressed
// bytes are not necessarily flushed until the Writer is closed.
func (z *Writer) Write(p []byte) (int, error)
