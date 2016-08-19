// Copyright 2009 The Go Authors. All rights reserved.
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
//     	log.Fatal(err)
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
//     	log.Fatal(err)
//     }
//     fmt.Printf("read %d bytes: %q\n", count, data[:count])

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
//     	log.Fatal(err)
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
//     	log.Fatal(err)
//     }
//     fmt.Printf("read %d bytes: %q\n", count, data[:count])
package os

import (
    "errors"
    "io"
    "runtime"
    "sync"
    "sync/atomic"
    "syscall"
    "time"
)

// DevNull is the name of the operating system's ``null device.''
// On Unix-like systems, it is "/dev/null"; on Windows, "NUL".
const DevNull = "/dev/null"


// The defined file mode bits are the most significant bits of the FileMode.
// The nine least-significant bits are the standard Unix rwxrwxrwx permissions.
// The values of these bits should be considered part of the public API and
// may be used in wire protocols or disk representations: they must not be
// changed, although new bits might be added.
const (
	// The single letters are the abbreviations
	// used by the String method's formatting.
	ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
	ModeAppend                                     // a: append-only
	ModeExclusive                                  // l: exclusive use
	ModeTemporary                                  // T: temporary file (not backed up)
	ModeSymlink                                    // L: symbolic link
	ModeDevice                                     // D: device file
	ModeNamedPipe                                  // p: named pipe (FIFO)
	ModeSocket                                     // S: Unix domain socket
	ModeSetuid                                     // u: setuid
	ModeSetgid                                     // g: setgid
	ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
	ModeSticky                                     // t: sticky
	// Mask for the type bits. For regular files, none will be set.
	ModeType          = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice
	ModePerm FileMode = 0777 // Unix permission bits

)


// Flags to OpenFile wrapping those of the underlying system. Not all
// flags may be implemented on a given system.
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
	PathSeparator     = '/' // OS-specific path separator
	PathListSeparator = ':' // OS-specific path list separator

)


// Seek whence values.

// Seek whence values.
//
// Deprecated: Use io.SeekStart, io.SeekCurrent, and io.SeekEnd.
const (
	SEEK_SET int = 0 // seek relative to the origin of the file
	SEEK_CUR int = 1 // seek relative to the current offset
	SEEK_END int = 2 // seek relative to the end

)


// Args hold the command-line arguments, starting with the program name.
var Args []string


// Portable analogs of some common system call errors.
var (
	ErrInvalid    = errors.New("invalid argument") // methods on File will return this error when the receiver is nil
	ErrPermission = errors.New("permission denied")
	ErrExist      = errors.New("file already exists")
	ErrNotExist   = errors.New("file does not exist")
)


// The only signal values guaranteed to be present on all systems
// are Interrupt (send the process an interrupt) and Kill (force
// the process to exit).
var (
	Interrupt Signal = syscall.SIGINT
	Kill      Signal = syscall.SIGKILL
)


// Stdin, Stdout, and Stderr are open Files pointing to the standard input,
// standard output, and standard error file descriptors.
var (
	Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
	Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)


// File represents an open file descriptor.
type File struct {
}


// A FileInfo describes a file and is returned by Stat and Lstat.
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

// A FileMode represents a file's mode and permission bits.
// The bits have the same definition on all systems, so that
// information about files can be moved from one system
// to another portably. Not all bits apply to all systems.
// The only required bit is ModeDir for directories.
type FileMode uint32


// LinkError records an error during a link or symlink or rename
// system call and the paths that caused it.
type LinkError struct {
	Op  string
	Old string
	New string
	Err error
}


// PathError records an error and the operation and file path that caused it.
type PathError struct {
	Op   string
	Path string
	Err  error
}


// ProcAttr holds the attributes that will be applied to a new process
// started by StartProcess.
type ProcAttr struct {
	// If Dir is non-empty, the child changes into the directory before
	// creating the process.
	Dir string
	// If Env is non-nil, it gives the environment variables for the
	// new process in the form returned by Environ.
	// If it is nil, the result of Environ will be used.
	Env []string
	// Files specifies the open files inherited by the new process. The
	// first three entries correspond to standard input, standard output, and
	// standard error. An implementation may support additional entries,
	// depending on the underlying operating system. A nil entry corresponds
	// to that file being closed when the process starts.
	Files []*File

	// Operating system-specific process creation attributes.
	// Note that setting this field means that your program
	// may not execute properly or even compile on some
	// operating systems.
	Sys *syscall.SysProcAttr
}


// Process stores the information about a process created by StartProcess.
type Process struct {
	Pid    int
	handle uintptr // handle is accessed atomically on Windows
	isdone uint32  // process has been successfully waited on, non zero if true
}


// ProcessState stores information about a process, as reported by Wait.
type ProcessState struct {
	pid    int                // The process's id.
	status syscall.WaitStatus // System-dependent status info.
	rusage *syscall.Rusage
}


// A Signal represents an operating system signal.
// The usual underlying implementation is operating system-dependent:
// on Unix it is syscall.Signal.
type Signal interface {
	String() string
	Signal() // to distinguish from other Stringers
}


// SyscallError records an error from a specific system call.
type SyscallError struct {
	Syscall string
	Err     error
}


// Chdir changes the current working directory to the named directory.
// If there is an error, it will be of type *PathError.
func Chdir(dir string) error

// Chmod changes the mode of the named file to mode.
// If the file is a symbolic link, it changes the mode of the link's target.
// If there is an error, it will be of type *PathError.
func Chmod(name string, mode FileMode) error

// Chown changes the numeric uid and gid of the named file. If the file is a
// symbolic link, it changes the uid and gid of the link's target. If there is
// an error, it will be of type *PathError.
func Chown(name string, uid, gid int) error

// Chtimes changes the access and modification times of the named
// file, similar to the Unix utime() or utimes() functions.
//
// The underlying filesystem may truncate or round the values to a
// less precise time unit.
// If there is an error, it will be of type *PathError.
func Chtimes(name string, atime time.Time, mtime time.Time) error

// Clearenv deletes all environment variables.
func Clearenv()

// Create creates the named file with mode 0666 (before umask), truncating
// it if it already exists. If successful, methods on the returned
// File can be used for I/O; the associated file descriptor has mode
// O_RDWR.
// If there is an error, it will be of type *PathError.
func Create(name string) (*File, error)

// Environ returns a copy of strings representing the environment,
// in the form "key=value".
func Environ() []string

// Exit causes the current program to exit with the given status code.
// Conventionally, code zero indicates success, non-zero an error.
// The program terminates immediately; deferred functions are not run.
func Exit(code int)

// Expand replaces ${var} or $var in the string based on the mapping function.
// For example, os.ExpandEnv(s) is equivalent to os.Expand(s, os.Getenv).
func Expand(s string, mapping func(string) string) string

// ExpandEnv replaces ${var} or $var in the string according to the values
// of the current environment variables.  References to undefined
// variables are replaced by the empty string.

// ExpandEnv replaces ${var} or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
func ExpandEnv(s string) string

// FindProcess looks for a running process by its pid.
//
// The Process it returns can be used to obtain information
// about the underlying operating system process.
//
// On Unix systems, FindProcess always succeeds and returns a Process
// for the given pid, regardless of whether the process exists.
func FindProcess(pid int) (*Process, error)

// Getegid returns the numeric effective group id of the caller.
func Getegid() int

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
func Getenv(key string) string

// Geteuid returns the numeric effective user id of the caller.
func Geteuid() int

// Getgid returns the numeric group id of the caller.
func Getgid() int

// Getgroups returns a list of the numeric ids of groups that the caller belongs
// to.
func Getgroups() ([]int, error)

// Getpagesize returns the underlying system's memory page size.
func Getpagesize() int

// Getpid returns the process id of the caller.
func Getpid() int

// Getppid returns the process id of the caller's parent.
func Getppid() int

// Getuid returns the numeric user id of the caller.
func Getuid() int

// Getwd returns a rooted path name corresponding to the
// current directory.  If the current directory can be
// reached via multiple paths (due to symbolic links),
// Getwd may return any one of them.

// Getwd returns a rooted path name corresponding to the
// current directory. If the current directory can be
// reached via multiple paths (due to symbolic links),
// Getwd may return any one of them.
func Getwd() (dir string, err error)

// Hostname returns the host name reported by the kernel.
func Hostname() (name string, err error)

// IsExist returns a boolean indicating whether the error is known to report
// that a file or directory already exists. It is satisfied by ErrExist as
// well as some syscall errors.
func IsExist(err error) bool

// IsNotExist returns a boolean indicating whether the error is known to
// report that a file or directory does not exist. It is satisfied by
// ErrNotExist as well as some syscall errors.
func IsNotExist(err error) bool

// IsPathSeparator reports whether c is a directory separator character.
func IsPathSeparator(c uint8) bool

// IsPermission returns a boolean indicating whether the error is known to
// report that permission is denied. It is satisfied by ErrPermission as well
// as some syscall errors.
func IsPermission(err error) bool

// Lchown changes the numeric uid and gid of the named file. If the file is a
// symbolic link, it changes the uid and gid of the link itself. If there is an
// error, it will be of type *PathError.
func Lchown(name string, uid, gid int) error

// Link creates newname as a hard link to the oldname file.
// If there is an error, it will be of type *LinkError.
func Link(oldname, newname string) error

// LookupEnv retrieves the value of the environment variable named
// by the key. If the variable is present in the environment the
// value (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be empty and the boolean will
// be false.
func LookupEnv(key string) (string, bool)

// Lstat returns a FileInfo describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link.  Lstat makes no attempt to follow the link.
// If there is an error, it will be of type *PathError.

// Lstat returns a FileInfo describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link. Lstat makes no attempt to follow the link.
// If there is an error, it will be of type *PathError.
func Lstat(name string) (FileInfo, error)

// Mkdir creates a new directory with the specified name and permission bits.
// If there is an error, it will be of type *PathError.
func Mkdir(name string, perm FileMode) error

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func MkdirAll(path string, perm FileMode) error

// NewFile returns a new File with the given file descriptor and name.
func NewFile(fd uintptr, name string) *File

// NewSyscallError returns, as an error, a new SyscallError
// with the given system call name and error details.
// As a convenience, if err is nil, NewSyscallError returns nil.
func NewSyscallError(syscall string, err error) error

// Open opens the named file for reading.  If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.

// Open opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func Open(name string) (*File, error)

// OpenFile is the generalized open call; most users will use Open
// or Create instead.  It opens the named file with specified flag
// (O_RDONLY etc.) and perm, (0666 etc.) if applicable.  If successful,
// methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.

// OpenFile is the generalized open call; most users will use Open
// or Create instead. It opens the named file with specified flag
// (O_RDONLY etc.) and perm, (0666 etc.) if applicable. If successful,
// methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.
func OpenFile(name string, flag int, perm FileMode) (*File, error)

// Pipe returns a connected pair of Files; reads from r return bytes written to
// w. It returns the files and an error, if any.
func Pipe() (r *File, w *File, err error)

// Readlink returns the destination of the named symbolic link.
// If there is an error, it will be of type *PathError.
func Readlink(name string) (string, error)

// Remove removes the named file or directory.
// If there is an error, it will be of type *PathError.
func Remove(name string) error

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters.  If the path does not exist, RemoveAll
// returns nil (no error).

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
func RemoveAll(path string) error

// Rename renames (moves) oldpath to newpath. If newpath already exists, Rename
// replaces it. OS-specific restrictions may apply when oldpath and newpath are
// in different directories. If there is an error, it will be of type
// *LinkError.
func Rename(oldpath, newpath string) error

// SameFile reports whether fi1 and fi2 describe the same file.
// For example, on Unix this means that the device and inode fields
// of the two underlying structures are identical; on other systems
// the decision may be based on the path names.
// SameFile only applies to results returned by this package's Stat.
// It returns false in other cases.
func SameFile(fi1, fi2 FileInfo) bool

// Setenv sets the value of the environment variable named by the key.
// It returns an error, if any.
func Setenv(key, value string) error

// StartProcess starts a new process with the program, arguments and attributes
// specified by name, argv and attr.
//
// StartProcess is a low-level interface. The os/exec package provides
// higher-level interfaces.
//
// If there is an error, it will be of type *PathError.
func StartProcess(name string, argv []string, attr *ProcAttr) (*Process, error)

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func Stat(name string) (FileInfo, error)

// Symlink creates newname as a symbolic link to oldname.
// If there is an error, it will be of type *LinkError.
func Symlink(oldname, newname string) error

// TempDir returns the default directory to use for temporary files.
func TempDir() string

// Truncate changes the size of the named file.
// If the file is a symbolic link, it changes the size of the link's target.
// If there is an error, it will be of type *PathError.
func Truncate(name string, size int64) error

// Unsetenv unsets a single environment variable.
func Unsetenv(key string) error

// Chdir changes the current working directory to the file,
// which must be a directory.
// If there is an error, it will be of type *PathError.
func (*File) Chdir() error

// Chmod changes the mode of the file to mode.
// If there is an error, it will be of type *PathError.
func (*File) Chmod(mode FileMode) error

// Chown changes the numeric uid and gid of the named file.
// If there is an error, it will be of type *PathError.
func (*File) Chown(uid, gid int) error

// Close closes the File, rendering it unusable for I/O.
// It returns an error, if any.
func (*File) Close() error

// Fd returns the integer Unix file descriptor referencing the open file. The
// file descriptor is valid only until f.Close is called or f is garbage
// collected.
func (*File) Fd() uintptr

// Name returns the name of the file as presented to Open.
func (*File) Name() string

// Read reads up to len(b) bytes from the File.
// It returns the number of bytes read and an error, if any.
// EOF is signaled by a zero count with err set to io.EOF.
func (*File) Read(b []byte) (n int, err error)

// ReadAt reads len(b) bytes from the File starting at byte offset off.
// It returns the number of bytes read and the error, if any.
// ReadAt always returns a non-nil error when n < len(b).
// At end of file, that error is io.EOF.
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
func (*File) Readdir(n int) ([]FileInfo, error)

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
func (*File) Readdirnames(n int) (names []string, err error)

// Seek sets the offset for the next Read or Write on file to offset,
// interpreted according to whence: 0 means relative to the origin of the file,
// 1 means relative to the current offset, and 2 means relative to the end. It
// returns the new offset and an error, if any. The behavior of Seek on a file
// opened with O_APPEND is not specified.
func (*File) Seek(offset int64, whence int) (ret int64, err error)

// Stat returns the FileInfo structure describing file.
// If there is an error, it will be of type *PathError.
func (*File) Stat() (FileInfo, error)

// Sync commits the current contents of the file to stable storage.
// Typically, this means flushing the file system's in-memory copy
// of recently written data to disk.
func (*File) Sync() error

// Truncate changes the size of the file.
// It does not change the I/O offset.
// If there is an error, it will be of type *PathError.
func (*File) Truncate(size int64) error

// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (*File) Write(b []byte) (n int, err error)

// WriteAt writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
func (*File) WriteAt(b []byte, off int64) (n int, err error)

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (*File) WriteString(s string) (n int, err error)

func (*LinkError) Error() string

func (*PathError) Error() string

// Kill causes the Process to exit immediately.
func (*Process) Kill() error

// Release releases any resources associated with the Process p,
// rendering it unusable in the future.
// Release only needs to be called if Wait is not.
func (*Process) Release() error

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.
func (*Process) Signal(sig Signal) error

// Wait waits for the Process to exit, and then returns a
// ProcessState describing its status and an error, if any.
// Wait releases any resources associated with the Process.
// On most operating systems, the Process must be a child
// of the current process or an error will be returned.
func (*Process) Wait() (*ProcessState, error)

// Exited reports whether the program has exited.
func (*ProcessState) Exited() bool

// Pid returns the process id of the exited process.
func (*ProcessState) Pid() int

func (*ProcessState) String() string

// Success reports whether the program exited successfully,
// such as with exit status 0 on Unix.
func (*ProcessState) Success() bool

// Sys returns system-dependent exit information about
// the process.  Convert it to the appropriate underlying
// type, such as syscall.WaitStatus on Unix, to access its contents.

// Sys returns system-dependent exit information about
// the process. Convert it to the appropriate underlying
// type, such as syscall.WaitStatus on Unix, to access its contents.
func (*ProcessState) Sys() interface{}

// SysUsage returns system-dependent resource usage information about
// the exited process.  Convert it to the appropriate underlying
// type, such as *syscall.Rusage on Unix, to access its contents.
// (On Unix, *syscall.Rusage matches struct rusage as defined in the
// getrusage(2) manual page.)

// SysUsage returns system-dependent resource usage information about
// the exited process. Convert it to the appropriate underlying
// type, such as *syscall.Rusage on Unix, to access its contents.
// (On Unix, *syscall.Rusage matches struct rusage as defined in the
// getrusage(2) manual page.)
func (*ProcessState) SysUsage() interface{}

// SystemTime returns the system CPU time of the exited process and its
// children.
func (*ProcessState) SystemTime() time.Duration

// UserTime returns the user CPU time of the exited process and its children.
func (*ProcessState) UserTime() time.Duration

func (*SyscallError) Error() string

// IsDir reports whether m describes a directory.
// That is, it tests for the ModeDir bit being set in m.
func (FileMode) IsDir() bool

// IsRegular reports whether m describes a regular file.
// That is, it tests that no mode type bits are set.
func (FileMode) IsRegular() bool

// Perm returns the Unix permission bits in m.
func (FileMode) Perm() FileMode

func (FileMode) String() string

