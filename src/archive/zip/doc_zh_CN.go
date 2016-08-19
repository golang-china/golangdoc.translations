// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package zip provides support for reading and writing ZIP archives.
//
// See: http://www.pkware.com/documents/casestudies/APPNOTE.TXT
//
// This package does not support disk spanning.
//
// A note about ZIP64:
//
// To be backwards compatible the FileHeader has both 32 and 64 bit Size
// fields. The 64 bit fields will always contain the correct value and
// for normal archives both fields will be the same. For files requiring
// the ZIP64 format the 32 bit fields will be 0xffffffff and the 64 bit
// fields must be used instead.

// zip包提供了zip档案文件的读写服务.
//
// 参见http://www.pkware.com/documents/casestudies/APPNOTE.TXT
//
// 本包不支持跨硬盘的压缩。
//
// 关于ZIP64：
//
// 为了向下兼容，FileHeader同时拥有32位和64位的Size字段。
// 64位字段总是包含正确的值，对普通格式的档案未见它们的值是相同的。
// 对zip64格式的档案文件32位字段将是0xffffffff，必须使用64位字段。
package zip

import (
    "bufio"
    "compress/flate"
    "encoding/binary"
    "errors"
    "fmt"
    "hash"
    "hash/crc32"
    "io"
    "io/ioutil"
    "os"
    "path"
    "sync"
    "time"
)

// Compression methods.

// 预定义压缩算法。
const (
	Store   uint16 = 0
	Deflate uint16 = 8
)



var (
	ErrFormat    = errors.New("zip: not a valid zip file")
	ErrAlgorithm = errors.New("zip: unsupported compression algorithm")
	ErrChecksum  = errors.New("zip: checksum error")
)


// A Compressor returns a new compressing writer, writing to w.
// The WriteCloser's Close method must be used to flush pending data to w.
// The Compressor itself must be safe to invoke from multiple goroutines
// simultaneously, but each returned writer will be used only by
// one goroutine at a time.

// Compressor 返回一个新的压缩写入器，写入到 w 中。WriteCloser 的 Close 方法必须
// 必须被用于将等待的数据刷新到 w 中。Compressor 在多个Go程被同步调用时， 其自身
// 必须保证安全，但每个返回的写入器一次只会被一个Go程使用。
type Compressor func(w io.Writer) (io.WriteCloser, error)


// A Decompressor returns a new decompressing reader, reading from r.
// The ReadCloser's Close method must be used to release associated resources.
// The Decompressor itself must be safe to invoke from multiple goroutines
// simultaneously, but each returned reader will be used only by
// one goroutine at a time.

// Decompressor 返回一个新的解压读取器，从 r 中读取。ReadCloser 的 Close
// 方法必须被用于释放相关的资源。Decompressor 在多个Go程被同步调用时，
// 其自身必须保证安全，但每个返回的读取器一次只会被一个Go程使用。
type Decompressor func(r io.Reader) io.ReadCloser



type File struct {
	zip          *Reader
	zipr         io.ReaderAt
	zipsize      int64
	headerOffset int64
}


// FileHeader describes a file within a zip file.
// See the zip spec for details.

// FileHeader描述zip文件中的一个文件。
// 参见zip的定义获取细节。
type FileHeader struct {

	// Name是文件名，它必须是相对路径，
	// 不能以设备或斜杠开始，只接受'/'作为路径分隔符
	Name string

	CreatorVersion     uint16
	ReaderVersion      uint16
	Flags              uint16
	Method             uint16
	ModifiedTime       uint16 // MS-DOS time // MS-DOS时间
	ModifiedDate       uint16 // MS-DOS date // MS-DOS日期
	CRC32              uint32
	CompressedSize     uint32 // Deprecated: Use CompressedSize64 instead. // 已弃用；请使用CompressedSize64
	UncompressedSize   uint32 // Deprecated: Use UncompressedSize64 instead. // 已弃用；请使用UncompressedSize64
	CompressedSize64   uint64
	UncompressedSize64 uint64
	Extra              []byte
	ExternalAttrs      uint32 // Meaning depends on CreatorVersion // 其含义依赖于CreatorVersion
	Comment            string
}



type ReadCloser struct {
	f *os.File
}



type Reader struct {
	r             io.ReaderAt
	File          []*File
	Comment       string
	decompressors map[uint16]Decompressor
}


// Writer implements a zip file writer.

// Writer类型实现了zip文件的写入器。
type Writer struct {
	cw          *countWriter
	dir         []*header
	last        *fileWriter
	closed      bool
	compressors map[uint16]Compressor
}


// FileInfoHeader creates a partially-populated FileHeader from an
// os.FileInfo.
// Because os.FileInfo's Name method returns only the base name of
// the file it describes, it may be necessary to modify the Name field
// of the returned header to provide the full path name of the file.

// FileInfoHeader返回一个根据fi填写了部分字段的Header。
// 因为os.FileInfo接口的Name方法只返回它描述的文件的无路径名，
// 有可能需要将返回值的Name字段修改为文件的完整路径名。
func FileInfoHeader(fi os.FileInfo) (*FileHeader, error)

// NewReader returns a new Reader reading from r, which is assumed to
// have the given size in bytes.

// NewReader返回一个从r读取数据的*Reader，r被假设其大小为size字节。
func NewReader(r io.ReaderAt, size int64) (*Reader, error)

// NewWriter returns a new Writer writing a zip file to w.

// NewWriter创建并返回一个将zip文件写入w的*Writer。
func NewWriter(w io.Writer) *Writer

// OpenReader will open the Zip file specified by name and return a ReadCloser.

// OpenReader会打开name指定的zip文件并返回一个*ReadCloser。
func OpenReader(name string) (*ReadCloser, error)

// RegisterCompressor registers custom compressors for a specified method ID.
// The common methods Store and Deflate are built in.

// RegisterCompressor使用指定的方法ID注册一个Compressor类型函数。
// 常用的方法Store和Deflate是内建的。
func RegisterCompressor(method uint16, comp Compressor)

// RegisterDecompressor allows custom decompressors for a specified method ID.
// The common methods Store and Deflate are built in.

// RegisterDecompressor使用指定的方法ID注册一个Decompressor类型函数。
// 通用方法 Store 和 Deflate 是内建的。
func RegisterDecompressor(method uint16, dcomp Decompressor)

// DataOffset returns the offset of the file's possibly-compressed
// data, relative to the beginning of the zip file.
//
// Most callers should instead use Open, which transparently
// decompresses data and verifies checksums.

// DataOffset返回文件的可能存在的压缩数据相对于zip文件起始的偏移量。
// 大多数调用者应使用Open代替，该方法会主动解压缩数据并验证校验和。
func (*File) DataOffset() (offset int64, err error)

// Open returns a ReadCloser that provides access to the File's contents.
// Multiple files may be read concurrently.

// Open方法返回一个io.ReadCloser接口，提供读取文件内容的方法。
// 可以同时读取多个文件。
func (*File) Open() (io.ReadCloser, error)

// FileInfo returns an os.FileInfo for the FileHeader.

// FileInfo返回一个根据h的信息生成的os.FileInfo。
func (*FileHeader) FileInfo() os.FileInfo

// ModTime returns the modification time in UTC.
// The resolution is 2s.

// 返回最近一次修改的UTC时间。（精度2s）
func (*FileHeader) ModTime() time.Time

// Mode returns the permission and mode bits for the FileHeader.

// Mode返回h的权限和模式位。
func (*FileHeader) Mode() (mode os.FileMode)

// SetModTime sets the ModifiedTime and ModifiedDate fields to the given time in
// UTC. The resolution is 2s.

// 将ModifiedTime和ModifiedDate字段设置为给定的UTC时间。（精度2s）
func (*FileHeader) SetModTime(t time.Time)

// SetMode changes the permission and mode bits for the FileHeader.

// SetMode修改h的权限和模式位。
func (*FileHeader) SetMode(mode os.FileMode)

// Close closes the Zip file, rendering it unusable for I/O.

// Close关闭zip文件，使它不能用于I/O。
func (*ReadCloser) Close() error

// RegisterDecompressor registers or overrides a custom decompressor for a
// specific method ID. If a decompressor for a given method is not found,
// Reader will default to looking up the decompressor at the package level.
func (*Reader) RegisterDecompressor(method uint16, dcomp Decompressor)

// Close finishes writing the zip file by writing the central directory.
// It does not (and can not) close the underlying writer.

// Close方法通过写入中央目录关闭该*Writer。
// 本方法不会也没办法关闭下层的io.Writer接口。
func (*Writer) Close() error

// Create adds a file to the zip file using the provided name.
// It returns a Writer to which the file contents should be written.
// The name must be a relative path: it must not start with a drive
// letter (e.g. C:) or leading slash, and only forward slashes are
// allowed.
// The file's contents must be written to the io.Writer before the next
// call to Create, CreateHeader, or Close.

// 使用给出的文件名添加一个文件进zip文件。
// 本方法返回一个io.Writer接口（用于写入新添加文件的内容）。
// 文件名必须是相对路径，不能以设备或斜杠开始，只接受'/'作为路径分隔。
// 新增文件的内容必须在下一次调用CreateHeader、Create或Close方法之前全部写入。
func (*Writer) Create(name string) (io.Writer, error)

// CreateHeader adds a file to the zip file using the provided FileHeader
// for the file metadata.
// It returns a Writer to which the file contents should be written.
//
// The file's contents must be written to the io.Writer before the next
// call to Create, CreateHeader, or Close. The provided FileHeader fh
// must not be modified after a call to CreateHeader.

// CreateHeader 使用给出的*FileHeader来作为文件的元数据添加一个文件进zip文件。
// 本方法返回一个io.Writer接口（用于写入新添加文件的内容）。
//
// 新增文件的内容必须在下一次调用CreateHeader、Create或Close方法之前全部写入。
// 提供的 FileHeader fh 在调用 CreateHeader 后决不能修改。
func (*Writer) CreateHeader(fh *FileHeader) (io.Writer, error)

// Flush flushes any buffered data to the underlying writer.
// Calling Flush is not normally necessary; calling Close is sufficient.
func (*Writer) Flush() error

// RegisterCompressor registers or overrides a custom compressor for a specific
// method ID. If a compressor for a given method is not found, Writer will
// default to looking up the compressor at the package level.
func (*Writer) RegisterCompressor(method uint16, comp Compressor)

// SetOffset sets the offset of the beginning of the zip data within the
// underlying writer. It should be used when the zip data is appended to an
// existing file, such as a binary executable.
// It must be called before any data is written.
func (*Writer) SetOffset(n int64)

