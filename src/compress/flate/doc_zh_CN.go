// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package flate implements the DEFLATE compressed data format, described in RFC
// 1951. The gzip and zlib packages implement access to DEFLATE-based file formats.

// flate 包实现了 deflate 压缩数据格式, 参见RFC 1951.
// gzip 包和 zlib 包实现了对基于 deflate 的文件格式的访问.
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

// NewReader 返回一个从 r 读取并解压数据的 io.ReadCloser.
// 调用者有责任在读取完毕后调用返回值的Close方法.
func NewReader(r io.Reader) io.ReadCloser

// NewReaderDict is like NewReader but initializes the reader with a preset
// dictionary. The returned Reader behaves as if the uncompressed data stream
// started with the given dictionary, which has already been read. NewReaderDict is
// typically used to read data compressed by NewWriterDict.
//
// The ReadCloser returned by NewReader also implements Resetter.

// NewReaderDict 类似 NewReader, 但会使用预设的字典初始化返回的 Reader.
//
// 返回的 Reader 表现的好像原始未压缩的数据流以该字典起始.
// NewReaderDict 用于读取 NewWriterDict 压缩的数据.
func NewReaderDict(r io.Reader, dict []byte) io.ReadCloser

// A CorruptInputError reports the presence of corrupt input at a given offset.

// CorruptInputError 表示在输入的指定偏移量位置存在损坏.
type CorruptInputError int64

func (e CorruptInputError) Error() string

// An InternalError reports an error in the flate code itself.

// InternalError 表示flate数据自身的错误.
type InternalError string

func (e InternalError) Error() string

// A ReadError reports an error encountered while reading input.

// ReadError 代表在读取输入流时遇到的错误.
type ReadError struct {
	Offset int64 // byte offset where error occurred
	Err    error // error returned by underlying Read
}

func (e *ReadError) Error() string

// The actual read interface needed by NewReader. If the passed in io.Reader does
// not also have ReadByte, the NewReader will introduce its own buffering.

// Reader 是 NewReader 真正需要的接口. 如果提供的 io.Reader 没有提供 ReadByte 方法,
// NewReader 函数会自行添加缓冲.
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

// WriteError 代表在写入输出流时遇到的错误.
type WriteError struct {
	Offset int64 // byte offset where error occurred
	Err    error // error returned by underlying Write
}

func (e *WriteError) Error() string

// A Writer takes data written to it and writes the compressed form of that data to
// an underlying writer (see NewWriter).

// Writer 将提供给它的数据压缩后写入下层的 io.Writer 接口.
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

// NewWriter 返回一个压缩水平为 level 的 Writer.
//
// 和 zlib 包一样, level 的范围是 1 (BestSpeed) 到9 (BestCompression);
// 值越大, 压缩效果越好, 但也越慢; level 为 0 表示不尝试做任何压缩,
// 只添加必需的 deflate 框架; level 为 -1 时会使用默认的压缩水平;
// 如果 level 在 [-1, 9] 范围内, error 返回值将是 nil, 否则将返回非 nil 的错误值.
func NewWriter(w io.Writer, level int) (*Writer, error)

// NewWriterDict is like NewWriter but initializes the new Writer with a preset
// dictionary. The returned Writer behaves as if the dictionary had been written to
// it without producing any compressed output. The compressed data written to w can
// only be decompressed by a Reader initialized with the same dictionary.

// NewWriterDict 类似NewWriter，但会使用预设的字典初始化返回的 Writer.
//
// 返回的 Writer 表现的好像已经将原始/未压缩数据 dict 写入 w 了,
// 使用 w 压缩的数据只能被使用同样的字典初始化生成的 Reader 接口解压缩.
func NewWriterDict(w io.Writer, level int, dict []byte) (*Writer, error)

// Close flushes and closes the writer.

// Close 刷新缓冲并关闭 w.
func (w *Writer) Close() error

// Flush flushes any pending compressed data to the underlying writer. It is useful
// mainly in compressed network protocols, to ensure that a remote reader has
// enough data to reconstruct a packet. Flush does not return until the data has
// been written. If the underlying writer returns an error, Flush returns that
// error.
//
// In the terminology of the zlib library, Flush is equivalent to Z_SYNC_FLUSH.

// Flush 将缓冲中的压缩数据刷新到下层 io.Writer 接口中.
//
// 本方法主要用在传输压缩数据的网络连接中, 以保证远端的接收者可以获得足够的数据来重构数据报.
// Flush 会阻塞直到所有缓冲中的数据都写入下层 io.Writer 接口后才返回.
// 如果下层的 io.Writetr 接口返回一个错误, Flush 也会返回该错误.
// 在zlib包的术语中, Flush 方法等价于 Z_SYNC_FLUSH.
func (w *Writer) Flush() error

// Reset discards the writer's state and makes it equivalent to the result of
// NewWriter or NewWriterDict called with dst and w's level and dictionary.

// Reset 将 w 重置, 丢弃当前的写入状态, 并将下层输出目标设为 dst.
// 效果上等价于将 w 设为使用 dst 和 w 的压缩水平/字典重新调用 NewWriter 或 NewWriterDict 返回的 *Writer.
func (w *Writer) Reset(dst io.Writer)

// Write writes data to w, which will eventually write the compressed form of data
// to its underlying writer.

// Write 向 w写入数据, 最终会将压缩后的数据写入下层 io.Writer接口.
func (w *Writer) Write(data []byte) (n int, err error)
