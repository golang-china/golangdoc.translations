// Copyright The Go Authors. All rights reserved.
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
// To be backwards compatible the FileHeader has both 32 and 64 bit Size fields.
// The 64 bit fields will always contain the correct value and for normal archives
// both fields will be the same. For files requiring the ZIP64 format the 32 bit
// fields will be 0xffffffff and the 64 bit fields must be used instead.

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

// RegisterCompressor registers custom compressors for a specified method ID. The
// common methods Store and Deflate are built in.

// RegisterCompressor使用指定的方法ID注册一个Compressor类型函数。
// 常用的方法Store和Deflate是内建的。
func RegisterCompressor(method uint16, comp Compressor)

// RegisterDecompressor allows custom decompressors for a specified method ID.

// RegisterDecompressor使用指定的方法ID注册一个Decompressor类型函数。
func RegisterDecompressor(method uint16, d Decompressor)

// A Compressor returns a compressing writer, writing to the provided writer. On
// Close, any pending data should be flushed.

// Compressor函数类型会返回一个io.WriteCloser，该接口会将数据压缩后写入提供的接口。
// 关闭时，应将缓冲中的数据刷新到下层接口中。
type Compressor func(io.Writer) (io.WriteCloser, error)

// Decompressor is a function that wraps a Reader with a decompressing Reader. The
// decompressed ReadCloser is returned to callers who open files from within the
// archive. These callers are responsible for closing this reader when they're
// finished reading.

// Decompressor函数类型会把一个io.Reader包装成具有decompressing特性的io.Reader.
// Decompressor函数类型会返回一个io.ReadCloser，
// 该接口的Read方法会将读取自提供的接口的数据提前解压缩。
// 程序员有责任在读取结束时关闭该io.ReadCloser。
type Decompressor func(io.Reader) io.ReadCloser

type File struct {
	FileHeader
	// contains filtered or unexported fields
}

// DataOffset returns the offset of the file's possibly-compressed data, relative
// to the beginning of the zip file.
//
// Most callers should instead use Open, which transparently decompresses data and
// verifies checksums.

// DataOffset返回文件的可能存在的压缩数据相对于zip文件起始的偏移量。
// 大多数调用者应使用Open代替，该方法会主动解压缩数据并验证校验和。
func (f *File) DataOffset() (offset int64, err error)

// Open returns a ReadCloser that provides access to the File's contents. Multiple
// files may be read concurrently.

// Open方法返回一个io.ReadCloser接口，提供读取文件内容的方法。
// 可以同时读取多个文件。
func (f *File) Open() (rc io.ReadCloser, err error)

// FileHeader describes a file within a zip file. See the zip spec for details.

// FileHeader描述zip文件中的一个文件。 参见zip的定义获取细节。
type FileHeader struct {
	// Name is the name of the file.
	// It must be a relative path: it must not start with a drive
	// letter (e.g. C:) or leading slash, and only forward slashes
	// are allowed.
	Name string

	CreatorVersion     uint16
	ReaderVersion      uint16
	Flags              uint16
	Method             uint16
	ModifiedTime       uint16 // MS-DOS time
	ModifiedDate       uint16 // MS-DOS date
	CRC32              uint32
	CompressedSize     uint32 // deprecated; use CompressedSize64
	UncompressedSize   uint32 // deprecated; use UncompressedSize64
	CompressedSize64   uint64
	UncompressedSize64 uint64
	Extra              []byte
	ExternalAttrs      uint32 // Meaning depends on CreatorVersion
	Comment            string
}

// FileInfoHeader creates a partially-populated FileHeader from an os.FileInfo.
// Because os.FileInfo's Name method returns only the base name of the file it
// describes, it may be necessary to modify the Name field of the returned header
// to provide the full path name of the file.

// FileInfoHeader返回一个根据fi填写了部分字段的Header。
// 因为os.FileInfo接口的Name方法只返回它描述的文件的无路径名，
// 有可能需要将返回值的Name字段修改为文件的完整路径名。
func FileInfoHeader(fi os.FileInfo) (*FileHeader, error)

// FileInfo returns an os.FileInfo for the FileHeader.

// FileInfo返回一个根据h的信息生成的os.FileInfo。
func (h *FileHeader) FileInfo() os.FileInfo

// ModTime returns the modification time in UTC. The resolution is 2s.

// 返回最近一次修改的UTC时间。（精度2s）
func (h *FileHeader) ModTime() time.Time

// Mode returns the permission and mode bits for the FileHeader.

// Mode返回h的权限和模式位。
func (h *FileHeader) Mode() (mode os.FileMode)

// SetModTime sets the ModifiedTime and ModifiedDate fields to the given time in
// UTC. The resolution is 2s.

// 将ModifiedTime和ModifiedDate字段设置为给定的UTC时间。（精度2s）
func (h *FileHeader) SetModTime(t time.Time)

// SetMode changes the permission and mode bits for the FileHeader.

// SetMode修改h的权限和模式位。
func (h *FileHeader) SetMode(mode os.FileMode)

type ReadCloser struct {
	Reader
	// contains filtered or unexported fields
}

// OpenReader will open the Zip file specified by name and return a ReadCloser.

// OpenReader会打开name指定的zip文件并返回一个*ReadCloser。
func OpenReader(name string) (*ReadCloser, error)

// Close closes the Zip file, rendering it unusable for I/O.

// Close关闭zip文件，使它不能用于I/O。
func (rc *ReadCloser) Close() error

type Reader struct {
	File    []*File
	Comment string
	// contains filtered or unexported fields
}

// NewReader returns a new Reader reading from r, which is assumed to have the
// given size in bytes.

// NewReader返回一个从r读取数据的*Reader，r被假设其大小为size字节。
func NewReader(r io.ReaderAt, size int64) (*Reader, error)

// Writer implements a zip file writer.

// Writer类型实现了zip文件的写入器。
type Writer struct {
	// contains filtered or unexported fields
}

// NewWriter returns a new Writer writing a zip file to w.

// NewWriter创建并返回一个将zip文件写入w的*Writer。
func NewWriter(w io.Writer) *Writer

// Close finishes writing the zip file by writing the central directory. It does
// not (and can not) close the underlying writer.

// Close方法通过写入中央目录关闭该*Writer。
// 本方法不会也没办法关闭下层的io.Writer接口。
func (w *Writer) Close() error

// Create adds a file to the zip file using the provided name. It returns a Writer
// to which the file contents should be written. The name must be a relative path:
// it must not start with a drive letter (e.g. C:) or leading slash, and only
// forward slashes are allowed. The file's contents must be written to the
// io.Writer before the next call to Create, CreateHeader, or Close.

// 使用给出的文件名添加一个文件进zip文件。
// 本方法返回一个io.Writer接口（用于写入新添加文件的内容）。
// 文件名必须是相对路径，不能以设备或斜杠开始，只接受'/'作为路径分隔。
// 新增文件的内容必须在下一次调用CreateHeader、Create或Close方法之前全部写入。
func (w *Writer) Create(name string) (io.Writer, error)

// CreateHeader adds a file to the zip file using the provided FileHeader for the
// file metadata. It returns a Writer to which the file contents should be written.
// The file's contents must be written to the io.Writer before the next call to
// Create, CreateHeader, or Close.

// 使用给出的*FileHeader来作为文件的元数据添加一个文件进zip文件。
// 本方法返回一个io.Writer接口（用于写入新添加文件的内容）。
// 新增文件的内容必须在下一次调用CreateHeader、Create或Close方法之前全部写入。
func (w *Writer) CreateHeader(fh *FileHeader) (io.Writer, error)

// Flush flushes any buffered data to the underlying writer. Calling Flush is not
// normally necessary; calling Close is sufficient.
func (w *Writer) Flush() error
