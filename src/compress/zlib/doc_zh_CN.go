// Copyright The Go Authors. All rights reserved.
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

// zlib 包实现了对 zlib 格式压缩数据的读写, 参见 RFC 1950.
//
// 本包的实现提供了在读取时解压和写入时压缩的滤镜.
// 例如, 将压缩数据写入一个 bytes.Buffer:
//
//     var b bytes.Buffer
//     w := zlib.NewWriter(&b)
//     w.Write([]byte("hello, world\n"))
//     w.Close()
//
// 然后将数据读取回来:
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

// 这些常量都是拷贝自 flate 包, 因此导入 "compress/zlib" 后,
// 就不必再导入 "compress/flate" 了.
const (
    NoCompression      = flate.NoCompression
    BestSpeed          = flate.BestSpeed
    BestCompression    = flate.BestCompression
    DefaultCompression = flate.DefaultCompression
)

var (
    // ErrChecksum is returned when reading ZLIB data that has an invalid checksum.
    ErrChecksum = errors.New("zlib: invalid checksum")
    // ErrDictionary is returned when reading ZLIB data that has an invalid dictionary.
    ErrDictionary = errors.New("zlib: invalid dictionary")
    // ErrHeader is returned when reading ZLIB data that has an invalid header.
    ErrHeader = errors.New("zlib: invalid header")
)

// Resetter resets a ReadCloser returned by NewReader or NewReaderDict to
// to switch to a new underlying Reader. This permits reusing a ReadCloser
// instead of allocating a new one.

// Resetter 复位由 NewReader/NewReaderDict 返回的 ReadCloser, 底层以切换到新的
// Reader. 允许重新使用 ReadCloser 而非分配一个新的.
type Resetter interface {
    // Reset discards any buffered data and resets the Resetter as if it was
    // newly initialized with the given reader.
    Reset(r io.Reader, dict []byte) error
}

// A Writer takes data written to it and writes the compressed
// form of that data to an underlying writer (see NewWriter).

// Writer 将提供给它的数据压缩后写入下层 io.Writer 接口.
type Writer struct {
}

// NewReader creates a new ReadCloser. Reads from the returned ReadCloser read
// and decompress data from r. The implementation buffers input and may read
// more data than necessary from r. It is the caller's responsibility to call
// Close on the ReadCloser when done.
//
// The ReadCloser returned by NewReader also implements Resetter.

// NewReader 返回一个从 r 读取并解压数据的 io.ReadCloser.
// 其实现会缓冲输入流的数据, 并可能从 r 中读取比需要的更多的数据.
// 调用者有责任在读取完毕后调用返回值的 Close 方法.
func NewReader(r io.Reader) (io.ReadCloser, error)

// NewReaderDict is like NewReader but uses a preset dictionary. NewReaderDict
// ignores the dictionary if the compressed data does not refer to it. If the
// compressed data refers to a different dictionary, NewReaderDict returns
// ErrDictionary.
//
// The ReadCloser returned by NewReaderDict also implements Resetter.

// NewReaderDict 类似 NewReader, 但会使用预设的字典初始化返回的 Reader.
//
// 如果压缩数据没有采用字典, 本函数会忽略该参数.
func NewReaderDict(r io.Reader, dict []byte) (io.ReadCloser, error)

// NewWriter creates a new Writer.
// Writes to the returned Writer are compressed and written to w.
//
// It is the caller's responsibility to call Close on the WriteCloser when done.
// Writes may be buffered and not flushed until Close.

// NewWriter 创建并返回一个 Writer. 写入返回值的数据都会在压缩后写入 w.
//
// 调用者有责任在结束写入后调用返回值的 Close 方法.
// 因为写入的数据可能保存在缓冲中没有刷新入下层.
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

// NewWriterLevelDict is like NewWriterLevel but specifies a dictionary to
// compress with.
//
// The dictionary may be nil. If not, its contents should not be modified until
// the Writer is closed.

// NewWriterLevelDict 类似 NewWriterLevel 但还指定了用于压缩的字典.
// dict 参数可以为 nil; 否则, 在返回的 Writer 关闭之前, 其内容不可被修改.
func NewWriterLevelDict(w io.Writer, level int, dict []byte) (*Writer, error)

// Close closes the Writer, flushing any unwritten data to the underlying
// io.Writer, but does not close the underlying io.Writer.

// 调用 Close 会刷新缓冲并关闭 w, 但不会关闭下层 io.Writer 接口.
func (*Writer) Close() error

// Flush flushes the Writer to its underlying io.Writer.

// Flush 将缓冲中的压缩数据刷新到下层 io.Writer 接口中.
func (*Writer) Flush() error

// Reset clears the state of the Writer z such that it is equivalent to its
// initial state from NewWriterLevel or NewWriterLevelDict, but instead writing
// to w.

// Reset 将 w 重置, 丢弃当前的写入状态, 并将下层输出目标设为 dst. 效果上等价于将
// w 设为使用 dst 和 w 的压缩水平, 字典重新调用
// NewWriterLevel/NewWriterLevelDict 返回的 *Writer.
func (*Writer) Reset(w io.Writer)

// Write writes a compressed form of p to the underlying io.Writer. The
// compressed bytes are not necessarily flushed until the Writer is closed or
// explicitly flushed.

// Write 将 p 压缩后写入下层 io.Writer 接口.
// 压缩后的数据不一定会立刻刷新, 除非 Writer 被关闭或者显式的刷新.
func (*Writer) Write(p []byte) (n int, err error)

