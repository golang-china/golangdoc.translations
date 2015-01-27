// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package flate implements the DEFLATE compressed data format, described in RFC
// 1951. The gzip and zlib packages implement access to DEFLATE-based file formats.
package flate

const (
	NoCompression = 0
	BestSpeed     = 1

	BestCompression    = 9
	DefaultCompression = -1
)

// NewReader returns a new ReadCloser that can be used to read the uncompressed
// version of r. If r does not also implement io.ByteReader, the decompressor may
// read more data than necessary from r. It is the caller's responsibility to call
// Close on the ReadCloser when finished reading.
//
// The ReadCloser returned by NewReader also implements Resetter.
func NewReader(r io.Reader) io.ReadCloser

// NewReaderDict is like NewReader but initializes the reader with a preset
// dictionary. The returned Reader behaves as if the uncompressed data stream
// started with the given dictionary, which has already been read. NewReaderDict is
// typically used to read data compressed by NewWriterDict.
//
// The ReadCloser returned by NewReader also implements Resetter.
func NewReaderDict(r io.Reader, dict []byte) io.ReadCloser

// A CorruptInputError reports the presence of corrupt input at a given offset.
type CorruptInputError int64

func (e CorruptInputError) Error() string

// An InternalError reports an error in the flate code itself.
type InternalError string

func (e InternalError) Error() string

// A ReadError reports an error encountered while reading input.
type ReadError struct {
	Offset int64 // byte offset where error occurred
	Err    error // error returned by underlying Read
}

func (e *ReadError) Error() string

// The actual read interface needed by NewReader. If the passed in io.Reader does
// not also have ReadByte, the NewReader will introduce its own buffering.
type Reader interface {
	io.Reader
	io.ByteReader
}

// Resetter resets a ReadCloser returned by NewReader or NewReaderDict to to switch
// to a new underlying Reader. This permits reusing a ReadCloser instead of
// allocating a new one.
type Resetter interface {
	// Reset discards any buffered data and resets the Resetter as if it was
	// newly initialized with the given reader.
	Reset(r io.Reader, dict []byte) error
}

// A WriteError reports an error encountered while writing output.
type WriteError struct {
	Offset int64 // byte offset where error occurred
	Err    error // error returned by underlying Write
}

func (e *WriteError) Error() string

// A Writer takes data written to it and writes the compressed form of that data to
// an underlying writer (see NewWriter).
type Writer struct {
	// contains filtered or unexported fields
}

// NewWriter returns a new Writer compressing data at the given level. Following
// zlib, levels range from 1 (BestSpeed) to 9 (BestCompression); higher levels
// typically run slower but compress more. Level 0 (NoCompression) does not attempt
// any compression; it only adds the necessary DEFLATE framing. Level -1
// (DefaultCompression) uses the default compression level.
//
// If level is in the range [-1, 9] then the error returned will be nil. Otherwise
// the error returned will be non-nil.
func NewWriter(w io.Writer, level int) (*Writer, error)

// NewWriterDict is like NewWriter but initializes the new Writer with a preset
// dictionary. The returned Writer behaves as if the dictionary had been written to
// it without producing any compressed output. The compressed data written to w can
// only be decompressed by a Reader initialized with the same dictionary.
func NewWriterDict(w io.Writer, level int, dict []byte) (*Writer, error)

// Close flushes and closes the writer.
func (w *Writer) Close() error

// Flush flushes any pending compressed data to the underlying writer. It is useful
// mainly in compressed network protocols, to ensure that a remote reader has
// enough data to reconstruct a packet. Flush does not return until the data has
// been written. If the underlying writer returns an error, Flush returns that
// error.
//
// In the terminology of the zlib library, Flush is equivalent to Z_SYNC_FLUSH.
func (w *Writer) Flush() error

// Reset discards the writer's state and makes it equivalent to the result of
// NewWriter or NewWriterDict called with dst and w's level and dictionary.
func (w *Writer) Reset(dst io.Writer)

// Write writes data to w, which will eventually write the compressed form of data
// to its underlying writer.
func (w *Writer) Write(data []byte) (n int, err error)
