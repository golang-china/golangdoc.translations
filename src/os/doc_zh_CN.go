// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package os provides a platform-independent interface to operating system
// functionality. The design is Unix-like, although the error handling is
// Go-like; failing calls return values of type error rather than error numbers.
// Often, more information is available within the error. For example, if a call
// that takes a file name fails, such as Open or Stat, the error will include
// the failing file name when printed and will be of type *PathError, which may
// be unpacked for more information.
//
// The os interface is intended to be uniform across all operating systems.
// Features not generally available appear in the system-specific package
// syscall.
//
// Here is a simple example, opening a file and reading some of it.
//
//     file, err := os.Open("file.go") // For read access.
//     if err != nil {
//         log.Fatal(err)
//     }
//
// If the open fails, the error string will be self-explanatory, like
//
//     open file.go: no such file or directory
//
// The file's data can then be read into a slice of bytes. Read and Write take
// their byte counts from the length of the argument slice.
//
//     data := make([]byte, 100)
//     count, err := file.Read(data)
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Printf("read %d bytes: %q\n", count, data[:count])

// os包提供了操作系统函数的不依赖平台的接口。设计为Unix风格的，虽然错误处理是go
// 风格的；失败的调用会返回错误值而非错误码。通常错误值里包含更多信息。例如，如
// 果某个使用一个文件名的调用（如Open、Stat）失败了，打印错误时会包含该文件名，
// 错误类型将为*PathError，其内部可以解包获得更多信息。
//
// os包的接口规定为在所有操作系统中都是一致的。非公用的属性可以从操作系统特定的
// syscall包获取。
//
// 下面是一个简单的例子，打开一个文件并从中读取一些数据：
//
//     file, err := os.Open("file.go") // For read access.
//     if err != nil {
//         log.Fatal(err)
//     }
//
// 如果打开失败，错误字符串是自解释的，例如：
//
//     open file.go: no such file or directory
//
// 文件的信息可以读取进一个[]byte切片。Read和Write方法从切片参数获取其内的字节数
// 。
//
//     data := make([]byte, 100)
//     count, err := file.Read(data)
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Printf("read %d bytes: %q\n", count, data[:count])
package os

import (
	"errors"
	"internal/syscall/windows"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

// DevNull is the name of the operating system's ``null device.'' On Unix-like
// systems, it is "/dev/null"; on Windows, "NUL".

// DevNull是操作系统空设备的名字。在类似Unix的操作系统中，是"/dev/null"；在
// Windows中，为"NUL"。
const DevNull = "/dev/null"

// The defined file mode bits are the most significant bits of the FileMode. The
// nine least-significant bits are the standard Unix rwxrwxrwx permissions. The
// values of these bits should be considered part of the public API and may be
// used in wire protocols or disk representations: they must not be changed,
// although new bits might be added.

// 这些被定义的位是FileMode最重要的位。另外9个不重要的位为标准Unix rwxrwxrwx权限
// （任何人都可读、写、运行）。这些（重要）位的值应被视为公共API的一部分，可能会
// 用于线路协议或硬盘标识：它们不能被修改，但可以添加新的位。
const (
	// The single letters are the abbreviations
	// used by the String method's formatting.
	ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory // d: 目录
	ModeAppend                                     // a: append-only // a: 只能写入，且只能写入到末尾
	ModeExclusive                                  // l: exclusive use // l: 用于执行
	ModeTemporary                                  // T: temporary file (not backed up) // T: 临时文件（非备份文件）
	ModeSymlink                                    // L: symbolic link // L: 符号链接（不是快捷方式文件）
	ModeDevice                                     // D: device file // D: 设备
	ModeNamedPipe                                  // p: named pipe (FIFO) // p: 命名管道（FIFO）
	ModeSocket                                     // S: Unix domain socket // S: Unix域socket
	ModeSetuid                                     // u: setuid // u: 表示文件具有其创建者用户id权限
	ModeSetgid                                     // g: setgid // g: 表示文件具有其创建者组id的权限
	ModeCharDevice                                 // c: Unix character device, when ModeDevice is set // c: 字符设备，需已设置ModeDevice
	ModeSticky                                     // t: sticky // t: 只有root/创建者能删除/移动文件

	// Mask for the type bits. For regular files, none will be set.

	// 覆盖所有类型位（用于通过&获取类型位），对普通文件，所有这些位都不应被设置
	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice

	ModePerm FileMode = 0777 // permission bits // 覆盖所有Unix权限位（用于通过&获取类型位）
)

const (
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // if possible, truncate file when opened.
)

const (
	PathSeparator     = '/' // OS-specific path separator // 操作系统指定的路径分隔符
	PathListSeparator = ':' // OS-specific path list separator // 操作系统指定的表分隔符
)

// Seek whence values.
const (
	SEEK_SET int = 0 // seek relative to the origin of the file // 相对于文件起始位置seek
	SEEK_CUR int = 1 // seek relative to the current offset // 相对于文件当前位置seek
	SEEK_END int = 2 // seek relative to the end // 相对于文件结尾位置seek
)

// Args hold the command-line arguments, starting with the program name.

// Args保管了命令行参数，第一个是程序名。
var Args []string

// Portable analogs of some common system call errors.
var (
	ErrInvalid    = errors.New("invalid argument")
	ErrPermission = errors.New("permission denied")
	ErrExist      = errors.New("file already exists")
	ErrNotExist   = errors.New("file does not exist")
)

// The only signal values guaranteed to be present on all systems are Interrupt
// (send the process an interrupt) and Kill (force the process to exit).

// 仅有的肯定会被所有操作系统提供的信号，Interrupt（中断信号）和Kill（强制退出信
// 号）。
var (
	Interrupt Signal = syscall.SIGINT
	Kill      Signal = syscall.SIGKILL
)

// Stdin, Stdout, and Stderr are open Files pointing to the standard input,
// standard output, and standard error file descriptors.

// Stdin、Stdout和Stderr是指向标准输入、标准输出、标准错误输出的文件描述符。
var (
	Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
	Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)

// File represents an open file descriptor.

// File代表一个打开的文件对象。
type File struct {
}

// A FileInfo describes a file and is returned by Stat and Lstat.

// FileInfo用来描述一个文件对象。
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}

// A FileMode represents a file's mode and permission bits.
// The bits have the same definition on all systems, so that
// information about files can be moved from one system
// to another portably.  Not all bits apply to all systems.
// The only required bit is ModeDir for directories.

// FileMode 代表文件的模式和权限位。这些字位在所有的操作系统都有相同的含义，因此
// 文件的信息可以在不同的操作系统之间安全的移植。不是所有的位都能用于所有的系
// 统，唯一共有的是用于表示目录的ModeDir位。
type FileMode uint32

// LinkError records an error during a link or symlink or rename
// system call and the paths that caused it.

// LinkError 记录在 Link、Symlink、Rename 系统调用时出现的错误，以及导致错误的
// 路径。
type LinkError struct {
	Op  string
	Old string
	New string
	Err error
}

// PathError records an error and the operation and file path that caused it.

// PathError记录一个错误，以及导致错误的路径。
type PathError struct {
	Op   string
	Path string
	Err  error
}

// ProcAttr holds the attributes that will be applied to a new process
// started by StartProcess.

// ProcAttr保管将被StartProcess函数用于一个新进程的属性。
type ProcAttr struct {
	// If Dir is non-empty, the child changes into the directory before
	// creating the process.
	Dir string
	// If Env is non-nil, it gives the environment variables for the
	// new process in the form returned by Environ.
	// If it is nil, the result of Environ will be used.
	Env []string
	// Files specifies the open files inherited by the new process.  The
	// first three entries correspond to standard input, standard output, and
	// standard error.  An implementation may support additional entries,
	// depending on the underlying operating system.  A nil entry corresponds
	// to that file being closed when the process starts.
	Files []*File

	// Operating system-specific process creation attributes.
	// Note that setting this field means that your program
	// may not execute properly or even compile on some
	// operating systems.
	Sys *syscall.SysProcAttr
}

// Process stores the information about a process created by StartProcess.

// Process保管一个被StarProcess创建的进程的信息。
type Process struct {
	Pid int
}

// ProcessState stores information about a process, as reported by Wait.

// ProcessState stores information about a process, as reported by Wait.

// ProcessState保管Wait函数报告的某个已退出进程的信息。
type ProcessState struct {
}

// A Signal represents an operating system signal.
// The usual underlying implementation is operating system-dependent:
// on Unix it is syscall.Signal.

// Signal 代表一个操作系统信号。一般其底层实现是依赖于操作系统的：在Unix中，它是
// syscall.Signal类型。
type Signal interface {
	String() string
	Signal() // to distinguish from other Stringers
}

// SyscallError records an error from a specific system call.

// SyscallError记录某个系统调用出现的错误。
type SyscallError struct {
	Syscall string
	Err     error
}

// Chdir changes the current working directory to the named directory.
// If there is an error, it will be of type *PathError.

// Chdir将当前工作目录修改为dir指定的目录。如果出错，会返回*PathError底层类型的
// 错误。
func Chdir(dir string) error

// Chmod changes the mode of the named file to mode.
// If the file is a symbolic link, it changes the mode of the link's target.
// If there is an error, it will be of type *PathError.

// Chmod修改name指定的文件对象的mode。如果name指定的文件是一个符号链接，它会修改
// 该链接的目的地文件的mode。如果出错，会返回*PathError底层类型的错误。
func Chmod(name string, mode FileMode) error

// Chown changes the numeric uid and gid of the named file. If the file is a
// symbolic link, it changes the uid and gid of the link's target. If there is
// an error, it will be of type *PathError.

// Chmod修改name指定的文件对象的用户id和组id。如果name指定的文件是一个符号链接，
// 它会修改该链接的目的地文件的用户id和组id。如果出错，会返回*PathError底层类型
// 的错误。
func Chown(name string, uid, gid int) error

// Chtimes changes the access and modification times of the named
// file, similar to the Unix utime() or utimes() functions.
//
// The underlying filesystem may truncate or round the values to a
// less precise time unit.
// If there is an error, it will be of type *PathError.

// Chtimes修改name指定的文件对象的访问时间和修改时间，类似Unix的utime()或
// utimes()函数。底层的文件系统可能会截断/舍入时间单位到更低的精确度。如果出错，
// 会返回*PathError底层类型的错误。
func Chtimes(name string, atime time.Time, mtime time.Time) error

// Clearenv deletes all environment variables.

// Clearenv删除所有环境变量。
func Clearenv()

// Create creates the named file with mode 0666 (before umask), truncating
// it if it already exists. If successful, methods on the returned
// File can be used for I/O; the associated file descriptor has mode
// O_RDWR.
// If there is an error, it will be of type *PathError.

// Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文
// 件已存在会截断它（为空文件）。如果成功，返回的文件对象可用于I/O；对应的文件描
// 述符具有O_RDWR模式。如果出错，错误底层类型是*PathError。
func Create(name string) (file *File, err error)

// Environ returns a copy of strings representing the environment,
// in the form "key=value".

// Environ返回表示环境变量的格式为"key=value"的字符串的切片拷贝。
func Environ() []string

// Exit causes the current program to exit with the given status code.
// Conventionally, code zero indicates success, non-zero an error.
// The program terminates immediately; deferred functions are not run.

// Exit让当前程序以给出的状态码code退出。一般来说，状态码0表示成功，非0表示出错
// 。程序会立刻终止，defer的函数不会被执行。
func Exit(code int)

// Expand replaces ${var} or $var in the string based on the mapping function.
// For example, os.ExpandEnv(s) is equivalent to os.Expand(s, os.Getenv).

// Expand函数替换s中的${var}或$var为mapping(var)。例如，os.ExpandEnv(s)等价于
// os.Expand(s, os.Getenv)。
func Expand(s string, mapping func(string) string) string

// ExpandEnv replaces ${var} or $var in the string according to the values
// of the current environment variables.  References to undefined
// variables are replaced by the empty string.

// ExpandEnv函数替换s中的${var}或$var为名为var
// 的环境变量的值。引用未定义环境变量会被替换为空字符串。
func ExpandEnv(s string) string

// FindProcess looks for a running process by its pid.
//
// The Process it returns can be used to obtain information
// about the underlying operating system process.
//
// On Unix systems, FindProcess always succeeds and returns a Process
// for the given pid, regardless of whether the process exists.

// FindProcess根据进程id查找一个运行中的进程。函数返回的进程对象可以用于获取其关
// 于底层操作系统进程的信息。
func FindProcess(pid int) (p *Process, err error)

// Getegid returns the numeric effective group id of the caller.

// Getegid返回调用者的有效组ID。
func Getegid() int

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.

// Getenv检索并返回名为key的环境变量的值。如果不存在该环境变量会返回空字符串。
func Getenv(key string) string

// Geteuid returns the numeric effective user id of the caller.

// Geteuid返回调用者的有效用户ID。
func Geteuid() int

// Getgid returns the numeric group id of the caller.

// Getgid返回调用者的组ID。
func Getgid() int

// Getgroups returns a list of the numeric ids of groups that the caller belongs
// to.

// Getgroups返回调用者所属的所有用户组的组ID。
func Getgroups() ([]int, error)

// Getpagesize returns the underlying system's memory page size.

// Getpagesize返回底层的系统内存页的尺寸。
func Getpagesize() int

// Getpid returns the process id of the caller.

// Getpid返回调用者所在进程的进程ID。
func Getpid() int

// Getppid returns the process id of the caller's parent.

// Getppid返回调用者所在进程的父进程的进程ID。
func Getppid() int

// Getuid returns the numeric user id of the caller.

// Getuid返回调用者的用户ID。
func Getuid() int

// Getwd returns a rooted path name corresponding to the
// current directory.  If the current directory can be
// reached via multiple paths (due to symbolic links),
// Getwd may return any one of them.

// Getwd返回一个对应当前工作目录的根路径。如果当前目录可以经过多条路径抵达（因为
// 硬链接），Getwd会返回其中一个。
func Getwd() (dir string, err error)

// Hostname returns the host name reported by the kernel.

// Hostname返回内核提供的主机名。
func Hostname() (name string, err error)

// IsExist returns a boolean indicating whether the error is known to report
// that a file or directory already exists. It is satisfied by ErrExist as
// well as some syscall errors.

// 返回一个布尔值说明该错误是否表示一个文件或目录已经存在。ErrExist和一些系统调
// 用错误会使它返回真。
func IsExist(err error) bool

// IsNotExist returns a boolean indicating whether the error is known to
// report that a file or directory does not exist. It is satisfied by
// ErrNotExist as well as some syscall errors.

// 返回一个布尔值说明该错误是否表示一个文件或目录不存在。ErrNotExist和一些系统调
// 用错误会使它返回真。
func IsNotExist(err error) bool

// IsPathSeparator reports whether c is a directory separator character.

// IsPathSeparator返回字符c是否是一个路径分隔符。
func IsPathSeparator(c uint8) bool

// IsPermission returns a boolean indicating whether the error is known to
// report that permission is denied. It is satisfied by ErrPermission as well
// as some syscall errors.

// 返回一个布尔值说明该错误是否表示因权限不足要求被拒绝。ErrPermission和一些系统
// 调用错误会使它返回真。
func IsPermission(err error) bool

// Lchown changes the numeric uid and gid of the named file. If the file is a
// symbolic link, it changes the uid and gid of the link itself. If there is an
// error, it will be of type *PathError.

// Chmod修改name指定的文件对象的用户id和组id。如果name指定的文件是一个符号链接，
// 它会修改该符号链接自身的用户id和组id。如果出错，会返回*PathError底层类型的错
// 误。
func Lchown(name string, uid, gid int) error

// Link creates newname as a hard link to the oldname file.
// If there is an error, it will be of type *LinkError.

// Link创建一个名为newname指向oldname的硬链接。如果出错，会返回*
// LinkError底层类型的错误。
func Link(oldname, newname string) error

// Lstat returns a FileInfo describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link.  Lstat makes no attempt to follow the link.
// If there is an error, it will be of type *PathError.

// Lstat返回一个描述name指定的文件对象的FileInfo。如果指定的文件对象是一个符号链
// 接，返回的FileInfo描述该符号链接的信息，本函数不会试图跳转该链接。如果出错，
// 返回的错误值为*PathError类型。
func Lstat(name string) (fi FileInfo, err error)

// Mkdir creates a new directory with the specified name and permission bits.
// If there is an error, it will be of type *PathError.

// Mkdir使用指定的权限和名称创建一个目录。如果出错，会返回*PathError底层类型的错
// 误。
func Mkdir(name string, perm FileMode) error

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.

// MkdirAll使用指定的权限和名称创建一个目录，包括任何必要的上级目录，并返回nil，
// 否则返回错误。权限位perm会应用在每一个被本函数创建的目录上。如果path指定了一
// 个已经存在的目录，MkdirAll不做任何操作并返回nil。
func MkdirAll(path string, perm FileMode) error

// NewFile returns a new File with the given file descriptor and name.

// NewFile使用给出的Unix文件描述符和名称创建一个文件。
func NewFile(fd uintptr, name string) *File

// NewSyscallError returns, as an error, a new SyscallError
// with the given system call name and error details.
// As a convenience, if err is nil, NewSyscallError returns nil.

// NewSyscallError返回一个指定系统调用名称和错误细节的SyscallError。如果err为nil
// ，本函数会返回nil。
func NewSyscallError(syscall string, err error) error

// Open opens the named file for reading.  If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.

// Open打开一个文件用于读取。如果操作成功，返回的文件对象的方法可用于读取数据；
// 对应的文件描述符具有O_RDONLY模式。如果出错，错误底层类型是*PathError。
func Open(name string) (file *File, err error)

// OpenFile is the generalized open call; most users will use Open
// or Create instead.  It opens the named file with specified flag
// (O_RDONLY etc.) and perm, (0666 etc.) if applicable.  If successful,
// methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.

// OpenFile是一个更一般性的文件打开函数，大多数调用者都应用Open或Create代替本函
// 数。它会使用指定的选项（如O_RDONLY等）、指定的模式（如0666等）打开指定名称的
// 文件。如果操作成功，返回的文件对象可用于I/O。如果出错，错误底层类型是
// *PathError。
func OpenFile(name string, flag int, perm FileMode) (file *File, err error)

// Pipe returns a connected pair of Files; reads from r return bytes
// written to w. It returns the files and an error, if any.

// Pipe返回一对关联的文件对象。从r的读取将返回写入w的数据。本函数会返回两个文件
// 对象和可能的错误。
func Pipe() (r *File, w *File, err error)

// Readlink returns the destination of the named symbolic link.
// If there is an error, it will be of type *PathError.

// Readlink获取name指定的符号链接文件指向的文件的路径。如果出错，会返回
// *PathError底层类型的错误。
func Readlink(name string) (string, error)

// Remove removes the named file or directory.
// If there is an error, it will be of type *PathError.

// Remove删除name指定的文件或目录。如果出错，会返回*PathError底层类型的错误。
func Remove(name string) error

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters.  If the path does not exist, RemoveAll
// returns nil (no error).

// RemoveAll删除path指定的文件，或目录及它包含的任何下级对象。它会尝试删除所有东
// 西，除非遇到错误并返回。如果path指定的对象不存在，RemoveAll会返回nil而不返回
// 错误。
func RemoveAll(path string) error

// Rename renames (moves) oldpath to newpath. If newpath already exists, Rename
// replaces it. OS-specific restrictions may apply when oldpath and newpath are
// in different directories. If there is an error, it will be of type
// *LinkError.

// Rename修改一个文件的名字，移动一个文件。可能会有一些个操作系统特定的限制。
func Rename(oldpath, newpath string) error

// SameFile reports whether fi1 and fi2 describe the same file.
// For example, on Unix this means that the device and inode fields
// of the two underlying structures are identical; on other systems
// the decision may be based on the path names.
// SameFile only applies to results returned by this package's Stat.
// It returns false in other cases.

// SameFile返回fi1和fi2是否在描述同一个文件。例如，在Unix这表示二者底层结构的设
// 备和索引节点是相同的；在其他系统中可能是根据路径名确定的。SameFile应只使用本
// 包Stat函数返回的FileInfo类型值为参数，其他情况下，它会返回假。
func SameFile(fi1, fi2 FileInfo) bool

// Setenv sets the value of the environment variable named by the key.
// It returns an error, if any.

// Setenv设置名为key的环境变量。如果出错会返回该错误。
func Setenv(key, value string) error

// StartProcess starts a new process with the program, arguments and attributes
// specified by name, argv and attr.
//
// StartProcess is a low-level interface. The os/exec package provides
// higher-level interfaces.
//
// If there is an error, it will be of type *PathError.

// StartProcess使用提供的属性、程序名、命令行参数开始一个新进程。StartProcess函
// 数是一个低水平的接口。os/exec包提供了高水平的接口，应该尽量使用该包。如果出错
// ，错误的底层类型会是*PathError。
func StartProcess(name string, argv []string, attr *ProcAttr) (*Process, error)

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.

// Stat返回一个描述name指定的文件对象的FileInfo。如果指定的文件对象是一个符号链
// 接，返回的FileInfo描述该符号链接指向的文件的信息，本函数会尝试跳转该链接。如
// 果出错，返回的错误值为*PathError类型。
func Stat(name string) (fi FileInfo, err error)

// Symlink creates newname as a symbolic link to oldname.
// If there is an error, it will be of type *LinkError.

// Symlink创建一个名为newname指向oldname的符号链接。如果出错，会返回*
// LinkError底层类型的错误。
func Symlink(oldname, newname string) error

// TempDir returns the default directory to use for temporary files.

// TempDir返回一个用于保管临时文件的默认目录。
func TempDir() string

// Truncate changes the size of the named file.
// If the file is a symbolic link, it changes the size of the link's target.
// If there is an error, it will be of type *PathError.

// Truncate修改name指定的文件的大小。如果该文件为一个符号链接，将修改链接指向的
// 文件的大小。如果出错，会返回*PathError底层类型的错误。
func Truncate(name string, size int64) error

// Unsetenv unsets a single environment variable.
func Unsetenv(key string) error

// Chdir changes the current working directory to the file,
// which must be a directory.
// If there is an error, it will be of type *PathError.

// Chdir将当前工作目录修改为f，f必须是一个目录。如果出错，错误底层类型是
// *PathError。
func (*File) Chdir() error

// Chmod changes the mode of the file to mode.
// If there is an error, it will be of type *PathError.

// Chmod修改文件的模式。如果出错，错误底层类型是*PathError。
func (*File) Chmod(mode FileMode) error

// Chown changes the numeric uid and gid of the named file.
// If there is an error, it will be of type *PathError.

// Chown修改文件的用户ID和组ID。如果出错，错误底层类型是*PathError。
func (*File) Chown(uid, gid int) error

// Close closes the File, rendering it unusable for I/O.
// It returns an error, if any.

// Close关闭文件f，使文件不能用于读写。它返回可能出现的错误。
func (*File) Close() error

// Fd returns the integer Plan 9 file descriptor referencing the open file. The
// file descriptor is valid only until f.Close is called or f is garbage
// collected.

// Fd返回与文件f对应的整数类型的Unix文件描述符。
func (*File) Fd() uintptr

// Name returns the name of the file as presented to Open.

// Name方法返回（提供给Open/Create等方法的）文件名称。
func (*File) Name() string

// Read reads up to len(b) bytes from the File.
// It returns the number of bytes read and an error, if any.
// EOF is signaled by a zero count with err set to io.EOF.

// Read方法从f中读取最多len(b)字节数据并写入b。它返回读取的字节数和可能遇到的任
// 何错误。文件终止标志是读取0个字节且返回值err为io.EOF。
func (*File) Read(b []byte) (n int, err error)

// ReadAt reads len(b) bytes from the File starting at byte offset off.
// It returns the number of bytes read and the error, if any.
// ReadAt always returns a non-nil error when n < len(b).
// At end of file, that error is io.EOF.

// ReadAt从指定的位置（相对于文件开始位置）读取len(b)字节数据并写入b。它返回读取
// 的字节数和可能遇到的任何错误。当n<len(b)时，本方法总是会返回错误；如果是因为
// 到达文件结尾，返回值err会是io.EOF。
func (*File) ReadAt(b []byte, off int64) (n int, err error)

// Readdir reads the contents of the directory associated with file and
// returns a slice of up to n FileInfo values, as would be returned
// by Lstat, in directory order. Subsequent calls on the same file will yield
// further FileInfos.
//
// If n > 0, Readdir returns at most n FileInfo structures. In this case, if
// Readdir returns an empty slice, it will return a non-nil error
// explaining why. At the end of a directory, the error is io.EOF.
//
// If n <= 0, Readdir returns all the FileInfo from the directory in
// a single slice. In this case, if Readdir succeeds (reads all
// the way to the end of the directory), it returns the slice and a
// nil error. If it encounters an error before the end of the
// directory, Readdir returns the FileInfo read until that point
// and a non-nil error.

// Readdir读取目录f的内容，返回一个有n个成员的[]FileInfo，这些FileInfo是被Lstat
// 返回的，采用目录顺序。对本函数的下一次调用会返回上一次调用剩余未读取的内容的
// 信息。
//
// 如果n>0，Readdir函数会返回一个最多n个成员的切片。这时，如果Readdir返回一个空
// 切片，它会返回一个非nil的错误说明原因。如果到达了目录f的结尾，返回值err会是
// io.EOF。
//
// 如果n<=0，Readdir函数返回目录中剩余所有文件对象的FileInfo构成的切片。此时，如
// 果Readdir调用成功（读取所有内容直到结尾），它会返回该切片和nil的错误值。如果
// 在到达结尾前遇到错误，会返回之前成功读取的FileInfo构成的切片和该错误。
func (*File) Readdir(n int) (fi []FileInfo, err error)

// Readdirnames reads and returns a slice of names from the directory f.
//
// If n > 0, Readdirnames returns at most n names. In this case, if
// Readdirnames returns an empty slice, it will return a non-nil error
// explaining why. At the end of a directory, the error is io.EOF.
//
// If n <= 0, Readdirnames returns all the names from the directory in
// a single slice. In this case, if Readdirnames succeeds (reads all
// the way to the end of the directory), it returns the slice and a
// nil error. If it encounters an error before the end of the
// directory, Readdirnames returns the names read until that point and
// a non-nil error.

// Readdir读取目录f的内容，返回一个有n个成员的[]string，切片成员为目录中文件对象
// 的名字，采用目录顺序。对本函数的下一次调用会返回上一次调用剩余未读取的内容的
// 信息。
//
// 如果n>0，Readdir函数会返回一个最多n个成员的切片。这时，如果Readdir返回一个空
// 切片，它会返回一个非nil的错误说明原因。如果到达了目录f的结尾，返回值err会是
// io.EOF。
//
// 如果n<=0，Readdir函数返回目录中剩余所有文件对象的名字构成的切片。此时，如果
// Readdir调用成功（读取所有内容直到结尾），它会返回该切片和nil的错误值。如果在
// 到达结尾前遇到错误，会返回之前成功读取的名字构成的切片和该错误。
func (*File) Readdirnames(n int) (names []string, err error)

// Seek sets the offset for the next Read or Write on file to offset,
// interpreted according to whence: 0 means relative to the origin of the file,
// 1 means relative to the current offset, and 2 means relative to the end. It
// returns the new offset and an error, if any. The behavior of Seek on a file
// opened with O_APPEND is not specified.

// Seek设置下一次读/写的位置。offset为相对偏移量，而whence决定相对位置：0为相对
// 文件开头，1为相对当前位置，2为相对文件结尾。它返回新的偏移量（相对开头）和可
// 能的错误。
func (*File) Seek(offset int64, whence int) (ret int64, err error)

// Stat returns the FileInfo structure describing file.
// If there is an error, it will be of type *PathError.

// Stat返回描述文件f的FileInfo类型值。如果出错，错误底层类型是*PathError。
func (*File) Stat() (fi FileInfo, err error)

// Sync commits the current contents of the file to stable storage.
// Typically, this means flushing the file system's in-memory copy
// of recently written data to disk.

// Sync递交文件的当前内容进行稳定的存储。一般来说，这表示将文件系统的最近写入的
// 数据在内存中的拷贝刷新到硬盘中稳定保存。
func (*File) Sync() (err error)

// Truncate changes the size of the file.
// It does not change the I/O offset.
// If there is an error, it will be of type *PathError.

// Truncate改变文件的大小，它不会改变I/O的当前位置。
// 如果截断文件，多出的部分就会被丢弃。如果出错，错误底层类型是*PathError。
func (*File) Truncate(size int64) error

// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).

// Write向文件中写入len(b)字节数据。它返回写入的字节数和可能遇到的任何错误。如果
// 返回值n!=len(b)，本方法会返回一个非nil的错误。
func (*File) Write(b []byte) (n int, err error)

// WriteAt writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).

// WriteAt在指定的位置（相对于文件开始位置）写入len(b)字节数据。它返回写入的字节
// 数和可能遇到的任何错误。如果返回值n!=len(b)，本方法会返回一个非nil的错误。
func (*File) WriteAt(b []byte, off int64) (n int, err error)

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.

// WriteString类似Write，但接受一个字符串参数。
func (*File) WriteString(s string) (ret int, err error)

func (*LinkError) Error() string

func (*PathError) Error() string

// Kill causes the Process to exit immediately.

// Kill让进程立刻退出。
func (*Process) Kill() error

// Release releases any resources associated with the Process p,
// rendering it unusable in the future.
// Release only needs to be called if Wait is not.

// Release释放进程p绑定的所有资源，
// 使它们（资源）不能再被（进程p）使用。只有没有调用Wait方法时才需要调用本方法。
func (*Process) Release() error

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.

// Signal方法向进程发送一个信号。在windows中向进程发送Interrupt信号尚未实现。
func (*Process) Signal(sig Signal) error

// Wait waits for the Process to exit, and then returns a
// ProcessState describing its status and an error, if any.
// Wait releases any resources associated with the Process.
// On most operating systems, the Process must be a child
// of the current process or an error will be returned.

// Wait方法阻塞直到进程退出，然后返回一个描述ProcessState描述进程的状态和可能的
// 错误。Wait方法会释放绑定到进程p的所有资源。在大多数操作系统中，进程p必须是当
// 前进程的子进程，否则会返回错误。
func (*Process) Wait() (*ProcessState, error)

// Exited reports whether the program has exited.

// Exited报告进程是否已退出。
func (*ProcessState) Exited() bool

// Pid returns the process id of the exited process.

// Pi返回一个已退出的进程的进程id。
func (*ProcessState) Pid() int

func (*ProcessState) String() string

// Success reports whether the program exited successfully,
// such as with exit status 0 on Unix.

// Success报告进程是否成功退出，如在Unix里以状态码0退出。
func (*ProcessState) Success() bool

// Sys returns system-dependent exit information about
// the process.  Convert it to the appropriate underlying
// type, such as syscall.WaitStatus on Unix, to access its contents.

// Sys返回该已退出进程系统特定的退出信息。需要将其类型转换为适当的底层类型，如
// Unix里转换为*syscall.WaitStatus类型以获取其内容。
func (*ProcessState) Sys() interface{}

// SysUsage returns system-dependent resource usage information about
// the exited process.  Convert it to the appropriate underlying
// type, such as *syscall.Rusage on Unix, to access its contents.
// (On Unix, *syscall.Rusage matches struct rusage as defined in the
// getrusage(2) manual page.)

// SysUsage返回该已退出进程系统特定的资源使用信息。需要将其类型转换为适当的底层
// 类型，如Unix里转换为*syscall.Rusage类型以获取其内容。
func (*ProcessState) SysUsage() interface{}

// SystemTime returns the system CPU time of the exited process and its
// children.

// SystemTime返回已退出进程及其子进程耗费的系统CPU时间。
func (*ProcessState) SystemTime() time.Duration

// UserTime returns the user CPU time of the exited process and its children.

// UserTime返回已退出进程及其子进程耗费的用户CPU时间。
func (*ProcessState) UserTime() time.Duration

func (*SyscallError) Error() string

// IsDir reports whether m describes a directory.
// That is, it tests for the ModeDir bit being set in m.

// IsDir报告m是否是一个目录。
func (FileMode) IsDir() bool

// IsRegular reports whether m describes a regular file.
// That is, it tests that no mode type bits are set.

// IsRegular报告m是否是一个普通文件。
func (FileMode) IsRegular() bool

// Perm returns the Unix permission bits in m.

// Perm方法返回m的Unix权限位。
func (FileMode) Perm() FileMode

func (FileMode) String() string
