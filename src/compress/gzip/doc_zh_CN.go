// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package gzip implements reading and writing of gzip format compressed files,
// as specified in RFC 1952.

// gzip 包实现了 gzip 格式压缩文件的读写, 参见RFC 1952.
package gzip

import (
    "bufio"
    "compress/flate"
    "errors"
    "fmt"
    "hash"
    "hash/crc32"
    "io"
    "time"
)

// These constants are copied from the flate package, so that code that imports
// "compress/gzip" does not also have to import "compress/flate".

// 这些常量都是拷贝自 flate 包, 因此导入 "compress/gzip" 后, 就不必再导入
// "compress/flate" 了.
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

// The gzip file stores a header giving metadata about the compressed file.
// That header is exposed as the fields of the Writer and Reader structs.
//
// Strings must be UTF-8 encoded and may only contain Unicode code points
// U+0001 through U+00FF, due to limitations of the GZIP file format.

// gzip 文件保存一个头域, 提供关于被压缩的文件的一些元数据.
// 该头域作为 Writer 和 Reader 类型的一个可导出字段, 可以提供给调用者访问.
type Header struct {
    Comment string    // comment
    Extra   []byte    // "extra data"
    ModTime time.Time // modification time
    Name    string    // file name
    OS      byte      // operating system type
}

// A Reader is an io.Reader that can be read to retrieve
// uncompressed data from a gzip-format compressed file.
//
// In general, a gzip file can be a concatenation of gzip files,
// each with its own header.  Reads from the Reader
// return the concatenation of the uncompressed data of each.
// Only the first header is recorded in the Reader fields.
//
// Gzip files store a length and checksum of the uncompressed data.
// The Reader will return a ErrChecksum when Read
// reaches the end of the uncompressed data if it does not
// have the expected length or checksum.  Clients should treat data
// returned by Read as tentative until they receive the io.EOF
// marking the end of the data.

// Reader 类型满足 io.Reader接口, 可以从 gzip 格式压缩文件读取并解压数据.
//
// 一般, 一个 gzip 文件可以是多个 gzip 文件的串联, 每一个都有自己的头域. 从
// Reader 读取数据会返回串联的每个文件的解压数据, 但只有第一个文件的头域被记录在
// Reader 的 Header 字段里.
//
// gzip 文件会保存未压缩数据的长度与校验和. 当读取到未压缩数据的结尾时, 如果数据
// 的长度或者校验和不正确, Reader 会返回 ErrCheckSum. 因此, 调用者应该将 Read 方
// 法返回的数据视为暂定的, 直到他们在数据结尾获得了一个 io.EOF.
type Reader struct {
    Header
}

// A Writer is an io.WriteCloser.
// Writes to a Writer are compressed and written to w.

// Writer 满足 io.WriteCloser接口. 它会将提供给它的数据压缩后写入下层 io.Writer
// 接口.
type Writer struct {
    Header
}

// NewReader creates a new Reader reading the given reader.
// If r does not also implement io.ByteReader,
// the decompressor may read more data than necessary from r.
//
// It is the caller's responsibility to call Close on the Reader when done.
//
// The Reader.Header fields will be valid in the Reader returned.

// NewReader 返回一个从 r 读取并解压数据的 *Reader.
// 其实现会缓冲输入流的数据, 并可能从 r 中读取比需要的更多的数据.
// 调用者有责任在读取完毕后调用返回值的 Close 方法.
func NewReader(r io.Reader) (*Reader, error)

// NewWriter returns a new Writer.
// Writes to the returned writer are compressed and written to w.
//
// It is the caller's responsibility to call Close on the WriteCloser when done.
// Writes may be buffered and not flushed until Close.
//
// Callers that wish to set the fields in Writer.Header must do so before
// the first call to Write, Flush, or Close.

// NewWriter 创建并返回一个 Writer. 写入返回值的数据都会在压缩后写入 w. 调用者有
// 责任在结束写入后调用返回值的 Close 方法. 因为写入的数据可能保存在缓冲中没有刷
// 新入下层.
//
// 如要设定 Writer.Header 字段, 调用者必须在第一次调用 Write 方法或者 Close 方法
// 之前设置. Header 字段的 Comment 和 Name 字段是 Go 的 utf-8 字符串, 但下层格式
// 要求为 NUL 中止的 ISO 8859-1 (Latin-1) 序列. 如果这两个字段的字符串包含 NUL
// 或非 Latin-1 字符, 将导致Write方法返回错误.
func NewWriter(w io.Writer) *Writer

// NewWriterLevel is like NewWriter but specifies the compression level instead
// of assuming DefaultCompression.
//
// The compression level can be DefaultCompression, NoCompression, or any
// integer value between BestSpeed and BestCompression inclusive. The error
// returned will be nil if the level is valid.

// NewWriterLevel 类似 NewWriter 但指定了压缩水平而不是采用默认的
// DefaultCompression.
//
// 参数 level 可以是 DefaultCompression/NoCompression/BestSpeed/BestCompression
// 之间包括二者的任何整数. 如果 level 合法, 返回的错误值为 nil.
func NewWriterLevel(w io.Writer, level int) (*Writer, error)

// Close closes the Reader. It does not close the underlying io.Reader.

// 调用 Close 会关闭 z, 但不会关闭下层 io.Reader 接口.
func (*Reader) Close() error

// Multistream controls whether the reader supports multistream files.
//
// If enabled (the default), the Reader expects the input to be a sequence of
// individually gzipped data streams, each with its own header and trailer,
// ending at EOF. The effect is that the concatenation of a sequence of gzipped
// files is treated as equivalent to the gzip of the concatenation of the
// sequence. This is standard behavior for gzip readers.
//
// Calling Multistream(false) disables this behavior; disabling the behavior can
// be useful when reading file formats that distinguish individual gzip data
// streams or mix gzip data streams with other data streams. In this mode, when
// the Reader reaches the end of the data stream, Read returns io.EOF. If the
// underlying reader implements io.ByteReader, it will be left positioned just
// after the gzip stream. To start the next stream, call z.Reset(r) followed by
// z.Multistream(false). If there is no next stream, z.Reset(r) will return
// io.EOF.
func (*Reader) Multistream(ok bool)

func (*Reader) Read(p []byte) (n int, err error)

// Reset discards the Reader z's state and makes it equivalent to the
// result of its original state from NewReader, but reading from r instead.
// This permits reusing a Reader rather than allocating a new one.

// Reset 将 z 重置, 丢弃当前的读取状态, 并将下层读取目标设为 r.
// 效果上等价于将 z 设为使用 r 重新调用 NewReader 返回的 Reader.
// 这让我们可以重用 z 而不是再申请一个新的.
func (*Reader) Reset(r io.Reader) error

// Close closes the Writer, flushing any unwritten data to the underlying
// io.Writer, but does not close the underlying io.Writer.

// 调用 Close 会关闭 z, 但不会关闭下层 io.Writer 接口.
func (*Writer) Close() error

// Flush flushes any pending compressed data to the underlying writer.
//
// It is useful mainly in compressed network protocols, to ensure that
// a remote reader has enough data to reconstruct a packet. Flush does
// not return until the data has been written. If the underlying
// writer returns an error, Flush returns that error.
//
// In the terminology of the zlib library, Flush is equivalent to Z_SYNC_FLUSH.

// Flush 将缓冲中的压缩数据刷新到下层 io.Writer 接口中.
//
// 本方法主要用在传输压缩数据的网络连接中, 以保证远端的接收者可以获得足够的数据
// 来重构数据报. Flush 会阻塞直到所有缓冲中的数据都写入下层 io.Writer 接口后才返
// 回. 如果下层的 io.Writetr 接口返回一个错误, Flush 也会返回该错误. 在 zlib 包
// 的术语中, Flush 方法等价于 Z_SYNC_FLUSH.
func (*Writer) Flush() error

// Reset discards the Writer z's state and makes it equivalent to the
// result of its original state from NewWriter or NewWriterLevel, but
// writing to w instead. This permits reusing a Writer rather than
// allocating a new one.

// Reset 将 z 重置, 丢弃当前的写入状态, 并将下层输出目标设为 dst. 效果上等价于将
// w 设为使用 dst 和 w 的压缩水平重新调用 NewWriterLevel 返回的 *Writer. 这让我
// 们可以重用 z 而不是再申请一个新的.
func (*Writer) Reset(w io.Writer)

// Write writes a compressed form of p to the underlying io.Writer. The
// compressed bytes are not necessarily flushed until the Writer is closed.

// Write 将 p 压缩后写入下层 io.Writer 接口.
// 压缩后的数据不一定会立刻刷新, 除非 Writer 被关闭或者显式的刷新.
func (*Writer) Write(p []byte) (int, error)

