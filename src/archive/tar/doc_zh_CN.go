// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package tar implements access to tar archives. It aims to cover most of the
// variations, including those produced by GNU and BSD tars.
//
// References:
//
//	http://www.freebsd.org/cgi/man.cgi?query=tar&sektion=5
//	http://www.gnu.org/software/tar/manual/html_node/Standard.html
//	http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html

// tar包实现了tar格式压缩文件的存取.
// 本包目标是覆盖大多数tar的变种，包括GNU和BSD生成的tar文件。
//
// 参见：
//
//	http://www.freebsd.org/cgi/man.cgi?query=tar&sektion=5
//	http://www.gnu.org/software/tar/manual/html_node/Standard.html
//	http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html
package tar

const (

	// Types
	TypeReg           = '0'    // regular file
	TypeRegA          = '\x00' // regular file
	TypeLink          = '1'    // hard link
	TypeSymlink       = '2'    // symbolic link
	TypeChar          = '3'    // character device node
	TypeBlock         = '4'    // block device node
	TypeDir           = '5'    // directory
	TypeFifo          = '6'    // fifo node
	TypeCont          = '7'    // reserved
	TypeXHeader       = 'x'    // extended header
	TypeXGlobalHeader = 'g'    // global extended header
	TypeGNULongName   = 'L'    // Next file has a long name
	TypeGNULongLink   = 'K'    // Next file symlinks to a file w/ a long name
	TypeGNUSparse     = 'S'    // sparse file
)

var (
	ErrWriteTooLong    = errors.New("archive/tar: write too long")
	ErrFieldTooLong    = errors.New("archive/tar: header field too long")
	ErrWriteAfterClose = errors.New("archive/tar: write after close")
)

var (
	ErrHeader = errors.New("archive/tar: invalid tar header")
)

// A Header represents a single header in a tar archive. Some fields may not be
// populated.

// Header代表tar档案文件里的单个头。 Header类型的某些字段可能未使用。
type Header struct {
	Name       string    // name of header file entry
	Mode       int64     // permission and mode bits
	Uid        int       // user id of owner
	Gid        int       // group id of owner
	Size       int64     // length in bytes
	ModTime    time.Time // modified time
	Typeflag   byte      // type of header entry
	Linkname   string    // target name of link
	Uname      string    // user name of owner
	Gname      string    // group name of owner
	Devmajor   int64     // major number of character or block device
	Devminor   int64     // minor number of character or block device
	AccessTime time.Time // access time
	ChangeTime time.Time // status change time
	Xattrs     map[string]string
}

// FileInfoHeader creates a partially-populated Header from fi. If fi describes a
// symlink, FileInfoHeader records link as the link target. If fi describes a
// directory, a slash is appended to the name. Because os.FileInfo's Name method
// returns only the base name of the file it describes, it may be necessary to
// modify the Name field of the returned header to provide the full path name of
// the file.

// FileInfoHeader返回一个根据fi填写了部分字段的Header。
// 如果fi描述一个符号链接，FileInfoHeader函数将link参数作为链接目标。
// 如果fi描述一个目录，会在名字后面添加斜杠。
// 因为os.FileInfo接口的Name方法只返回它描述的文件的无路径名，
// 有可能需要将返回值的Name字段修改为文件的完整路径名。
func FileInfoHeader(fi os.FileInfo, link string) (*Header, error)

// FileInfo returns an os.FileInfo for the Header.

// FileInfo返回该Header对应的文件信息。（os.FileInfo类型）
func (h *Header) FileInfo() os.FileInfo

// A Reader provides sequential access to the contents of a tar archive. A tar
// archive consists of a sequence of files. The Next method advances to the next
// file in the archive (including the first), and then it can be treated as an
// io.Reader to access the file's data.

// Reader提供了对一个tar档案文件的顺序读取。
// 一个tar档案文件包含一系列文件。
// Next方法返回档案中的下一个文件（包括第一个），
// 返回值可以被视为io.Reader来获取文件的数据。
type Reader struct {
	// contains filtered or unexported fields
}

// NewReader creates a new Reader reading from r.

// NewReader创建一个从r读取的Reader。
func NewReader(r io.Reader) *Reader

// Next advances to the next entry in the tar archive.

// Next 将前进到 tar 归档文件中的下一条记录，
//
// 在输入到达结尾时将返回 io.EOF。
func (tr *Reader) Next() (*Header, error)

// Read reads from the current entry in the tar archive. It returns 0, io.EOF when
// it reaches the end of that entry, until Next is called to advance to the next
// entry.

// 从档案文件的当前记录读取数据， 到达记录末端时返回(0,
// EOF)，直到调用Next方法转入下一记录。
func (tr *Reader) Read(b []byte) (n int, err error)

// A Writer provides sequential writing of a tar archive in POSIX.1 format. A tar
// archive consists of a sequence of files. Call WriteHeader to begin a new file,
// and then call Write to supply that file's data, writing at most hdr.Size bytes
// in total.

// Writer类型提供了POSIX.1格式的tar档案文件的顺序写入。
// 一个tar档案文件包含一系列文件。 调用WriteHeader来写入一个新的文件，
// 然后调用Write写入文件的数据，该记录写入的数据不能超过hdr.Size字节。
type Writer struct {
	// contains filtered or unexported fields
}

// NewWriter creates a new Writer writing to w.

// NewWriter创建一个写入w的*Writer。
func NewWriter(w io.Writer) *Writer

// Close closes the tar archive, flushing any unwritten data to the underlying
// writer.

// Close关闭tar档案文件，
// 会将缓冲中未写入下层的io.Writer接口的数据刷新到下层。
func (tw *Writer) Close() error

// Flush finishes writing the current file (optional).

// Flush结束当前文件的写入。（可选的）
func (tw *Writer) Flush() error

// Write writes to the current entry in the tar archive. Write returns the error
// ErrWriteTooLong if more than hdr.Size bytes are written after WriteHeader.

// Write向tar档案文件的当前记录中写入数据。
// 如果写入的数据总数超出上一次调用WriteHeader的参数hdr.Size字节， 返回ErrWriteTooLong错误。
func (tw *Writer) Write(b []byte) (n int, err error)

// WriteHeader writes hdr and prepares to accept the file's contents. WriteHeader
// calls Flush if it is not the first header. Calling after a Close will return
// ErrWriteAfterClose.

// WriteHeader写入hdr并准备接受文件内容。
// 如果不是第一次调用本方法，会调用Flush。
// 在Close之后调用本方法会返回ErrWriteAfterClose。
func (tw *Writer) WriteHeader(hdr *Header) error
