// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package syscall contains an interface to the low-level operating system
// primitives. The details vary depending on the underlying system, and
// by default, godoc will display the syscall documentation for the current
// system. If you want godoc to display syscall documentation for another
// system, set $GOOS and $GOARCH to the desired system. For example, if
// you want to view documentation for freebsd/arm on linux/amd64, set $GOOS
// to freebsd and $GOARCH to arm.
// The primary use of syscall is inside other packages that provide a more
// portable interface to the system, such as "os", "time" and "net".  Use
// those packages rather than this one if you can.
// For details of the functions and data types in this package consult
// the manuals for the appropriate operating system.
// These calls return err == nil to indicate success; otherwise
// err is an operating system error describing the failure.
// On most systems, that error has type syscall.Errno.
//
// NOTE: This package is locked down. Code outside the standard
// Go repository should be migrated to use the corresponding
// package in the golang.org/x/sys repository. That is also where updates
// required by new systems or versions should be applied.
// See https://golang.org/s/go1.4-syscall for more information.
package syscall

import (
	"internal/race"
	"runtime"
	"sync"
	"unsafe"
)

const ImplementsGetwd = true

var ForkLock sync.RWMutex

// For testing: clients can set this flag to force
// creation of IPv6 sockets to return EAFNOSUPPORT.
var SocketDisableIPv6 bool

var (
	Stdin  = 0
	Stdout = 1
	Stderr = 2
)

// Credential holds user and group identities to be assumed
// by a child process started by StartProcess.
type Credential struct {
	Uid    uint32   // User ID.
	Gid    uint32   // Group ID.
	Groups []uint32 // Supplementary group IDs.
}

// An Errno is an unsigned number describing an error condition.
// It implements the error interface. The zero Errno is by convention
// a non-error, so code to convert from Errno to error should use:
// 	err = nil
// 	if errno != 0 {
// 		err = errno
// 	}
type Errno uintptr

// NetlinkMessage represents a netlink message.
type NetlinkMessage struct {
	Header NlMsghdr
	Data   []byte
}

// NetlinkRouteAttr represents a netlink route attribute.
type NetlinkRouteAttr struct {
	Attr  RtAttr
	Value []byte
}

// NetlinkRouteRequest represents a request message to receive routing
// and link states from the kernel.
type NetlinkRouteRequest struct {
	Header NlMsghdr
	Data   RtGenmsg
}

// ProcAttr holds attributes that will be applied to a new process started
// by StartProcess.
type ProcAttr struct {
	Dir   string    // Current working directory.
	Env   []string  // Environment.
	Files []uintptr // File descriptors.
	Sys   *SysProcAttr
}

// A Signal is a number describing a process signal.
// It implements the os.Signal interface.
type Signal int

type Sockaddr interface {
	sockaddr() (ptr unsafe.Pointer, len _Socklen, err error) // lowercase; only we can define Sockaddrs
}

type SockaddrInet4 struct {
	Port int
	Addr [4]byte
}

type SockaddrInet6 struct {
	Port   int
	ZoneId uint32
	Addr   [16]byte
}

type SockaddrLinklayer struct {
	Protocol uint16
	Ifindex  int
	Hatype   uint16
	Pkttype  uint8
	Halen    uint8
	Addr     [8]byte
}

type SockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

type SockaddrUnix struct {
	Name string
}

// SocketControlMessage represents a socket control message.
type SocketControlMessage struct {
	Header Cmsghdr
	Data   []byte
}

type SysProcAttr struct {
	Chroot       string         // Chroot.
	Credential   *Credential    // Credential.
	Ptrace       bool           // Enable tracing.
	Setsid       bool           // Create session.
	Setpgid      bool           // Set process group ID to Pgid, or, if Pgid == 0, to new pid.
	Setctty      bool           // Set controlling terminal to fd Ctty (only meaningful if Setsid is set)
	Noctty       bool           // Detach fd 0 from controlling terminal
	Ctty         int            // Controlling TTY fd
	Foreground   bool           // Place child's process group in foreground. (Implies Setpgid. Uses Ctty as fd of controlling TTY)
	Pgid         int            // Child's process group ID if Setpgid.
	Pdeathsig    Signal         // Signal that the process will get when its parent dies (Linux only)
	Cloneflags   uintptr        // Flags for clone calls (Linux only)
	Unshareflags uintptr        // Flags for unshare calls (Linux only)
	UidMappings  []SysProcIDMap // User ID mappings for user namespaces.
	GidMappings  []SysProcIDMap // Group ID mappings for user namespaces.

	// GidMappingsEnableSetgroups enabling setgroups syscall. If false, then
	// setgroups syscall will be disabled for the child process. This parameter
	// is no-op if GidMappings == nil. Otherwise for unprivileged users this
	// should be set to false for mappings work.
	GidMappingsEnableSetgroups bool
}

// SysProcIDMap holds Container ID to Host ID mappings used for User Namespaces
// in Linux. See user_namespaces(7).
type SysProcIDMap struct {
	ContainerID int // Container ID.
	HostID      int // Host ID.
	Size        int // Size.
}

type WaitStatus uint32

func Accept(fd int) (nfd int, sa Sockaddr, err error)

func Accept4(fd int, flags int) (nfd int, sa Sockaddr, err error)

func Access(path string, mode uint32) (err error)

// Deprecated: Use golang.org/x/net/bpf instead.
func AttachLsf(fd int, i []SockFilter) error

func Bind(fd int, sa Sockaddr) (err error)

// BindToDevice binds the socket associated with fd to device.
func BindToDevice(fd int, device string) (err error)

// BytePtrFromString returns a pointer to a NUL-terminated array of
// bytes containing the text of s. If s contains a NUL byte at any
// location, it returns (nil, EINVAL).
func BytePtrFromString(s string) (*byte, error)

// ByteSliceFromString returns a NUL-terminated slice of bytes
// containing the text of s. If s contains a NUL byte at any
// location, it returns (nil, EINVAL).
func ByteSliceFromString(s string) ([]byte, error)

func Chmod(path string, mode uint32) (err error)

func Chown(path string, uid int, gid int) (err error)

func Clearenv()

func CloseOnExec(fd int)

// CmsgLen returns the value to store in the Len field of the Cmsghdr
// structure, taking into account any necessary alignment.
func CmsgLen(datalen int) int

// CmsgSpace returns the number of bytes an ancillary element with
// payload of the passed data length occupies.
func CmsgSpace(datalen int) int

func Connect(fd int, sa Sockaddr) (err error)

func Creat(path string, mode uint32) (fd int, err error)

// Deprecated: Use golang.org/x/net/bpf instead.
func DetachLsf(fd int) error

func Environ() []string

// Ordinary exec.
func Exec(argv0 string, argv []string, envv []string) (err error)

// FcntlFlock performs a fcntl syscall for the F_GETLK, F_SETLK or F_SETLKW
// command.
func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error

// Combination of fork and exec, careful to be thread safe.
func ForkExec(argv0 string, argv []string, attr *ProcAttr) (pid int, err error)

func Fstat(fd int, s *Stat_t) (err error)

func Futimes(fd int, tv []Timeval) (err error)

func Futimesat(dirfd int, path string, tv []Timeval) (err error)

func Getenv(key string) (value string, found bool)

func Getgroups() (gids []int, err error)

func Getpagesize() int

func Getpeername(fd int) (sa Sockaddr, err error)

func Getpgrp() (pid int)

func Getsockname(fd int) (sa Sockaddr, err error)

func GetsockoptICMPv6Filter(fd, level, opt int) (*ICMPv6Filter, error)

func GetsockoptIPMreq(fd, level, opt int) (*IPMreq, error)

func GetsockoptIPMreqn(fd, level, opt int) (*IPMreqn, error)

func GetsockoptIPv6MTUInfo(fd, level, opt int) (*IPv6MTUInfo, error)

func GetsockoptIPv6Mreq(fd, level, opt int) (*IPv6Mreq, error)

func GetsockoptInet4Addr(fd, level, opt int) (value [4]byte, err error)

func GetsockoptInt(fd, level, opt int) (value int, err error)

func GetsockoptUcred(fd, level, opt int) (*Ucred, error)

func Gettimeofday(tv *Timeval) (err error)

func Getwd() (wd string, err error)

func Ioperm(from int, num int, on int) (err error)

func Iopl(level int) (err error)

func Link(oldpath string, newpath string) (err error)

// Deprecated: Use golang.org/x/net/bpf instead.
func LsfJump(code, k, jt, jf int) *SockFilter

// Deprecated: Use golang.org/x/net/bpf instead.
func LsfSocket(ifindex, proto int) (int, error)

// Deprecated: Use golang.org/x/net/bpf instead.
func LsfStmt(code, k int) *SockFilter

func Lstat(path string, s *Stat_t) (err error)

func Mkdir(path string, mode uint32) (err error)

func Mkfifo(path string, mode uint32) (err error)

func Mknod(path string, mode uint32, dev int) (err error)

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)

func Mount(source string, target string, fstype string, flags uintptr, data string) (err error)

func Munmap(b []byte) (err error)

// NetlinkRIB returns routing information base, as known as RIB, which
// consists of network facility information, states and parameters.
func NetlinkRIB(proto, family int) ([]byte, error)

func NsecToTimespec(nsec int64) (ts Timespec)

func NsecToTimeval(nsec int64) (tv Timeval)

func Open(path string, mode int, perm uint32) (fd int, err error)

func Openat(dirfd int, path string, flags int, mode uint32) (fd int, err error)

func ParseDirent(buf []byte, max int, names []string) (consumed int, count int, newnames []string)

// ParseNetlinkMessage parses b as an array of netlink messages and
// returns the slice containing the NetlinkMessage structures.
func ParseNetlinkMessage(b []byte) ([]NetlinkMessage, error)

// ParseNetlinkRouteAttr parses m's payload as an array of netlink
// route attributes and returns the slice containing the
// NetlinkRouteAttr structures.
func ParseNetlinkRouteAttr(m *NetlinkMessage) ([]NetlinkRouteAttr, error)

// ParseSocketControlMessage parses b as an array of socket control
// messages.
func ParseSocketControlMessage(b []byte) ([]SocketControlMessage, error)

// ParseUnixCredentials decodes a socket control message that contains
// credentials in a Ucred structure. To receive such a message, the
// SO_PASSCRED option must be enabled on the socket.
func ParseUnixCredentials(m *SocketControlMessage) (*Ucred, error)

// ParseUnixRights decodes a socket control message that contains an
// integer array of open file descriptors from another process.
func ParseUnixRights(m *SocketControlMessage) ([]int, error)

func Pipe(p []int) (err error)

func Pipe2(p []int, flags int) (err error)

func PtraceAttach(pid int) (err error)

func PtraceCont(pid int, signal int) (err error)

func PtraceDetach(pid int) (err error)

func PtraceGetEventMsg(pid int) (msg uint, err error)

func PtraceGetRegs(pid int, regsout *PtraceRegs) (err error)

func PtracePeekData(pid int, addr uintptr, out []byte) (count int, err error)

func PtracePeekText(pid int, addr uintptr, out []byte) (count int, err error)

func PtracePokeData(pid int, addr uintptr, data []byte) (count int, err error)

func PtracePokeText(pid int, addr uintptr, data []byte) (count int, err error)

func PtraceSetOptions(pid int, options int) (err error)

func PtraceSetRegs(pid int, regs *PtraceRegs) (err error)

func PtraceSingleStep(pid int) (err error)

func PtraceSyscall(pid int, signal int) (err error)

func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)

func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)

func Read(fd int, p []byte) (n int, err error)

func ReadDirent(fd int, buf []byte) (n int, err error)

func Readlink(path string, buf []byte) (n int, err error)

func Reboot(cmd int) (err error)

func Recvfrom(fd int, p []byte, flags int) (n int, from Sockaddr, err error)

func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error)

func Rename(oldpath string, newpath string) (err error)

func Rmdir(path string) error

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error)

func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error)

func SendmsgN(fd int, p, oob []byte, to Sockaddr, flags int) (n int, err error)

func Sendto(fd int, p []byte, flags int, to Sockaddr) (err error)

// Deprecated: Use golang.org/x/net/bpf instead.
func SetLsfPromisc(name string, m bool) error

func SetNonblock(fd int, nonblocking bool) (err error)

func Setenv(key, value string) error

func Setgid(gid int) (err error)

func Setgroups(gids []int) (err error)

func SetsockoptByte(fd, level, opt int, value byte) (err error)

func SetsockoptICMPv6Filter(fd, level, opt int, filter *ICMPv6Filter) error

func SetsockoptIPMreq(fd, level, opt int, mreq *IPMreq) (err error)

func SetsockoptIPMreqn(fd, level, opt int, mreq *IPMreqn) (err error)

func SetsockoptIPv6Mreq(fd, level, opt int, mreq *IPv6Mreq) (err error)

func SetsockoptInet4Addr(fd, level, opt int, value [4]byte) (err error)

func SetsockoptInt(fd, level, opt int, value int) (err error)

func SetsockoptLinger(fd, level, opt int, l *Linger) (err error)

func SetsockoptString(fd, level, opt int, s string) (err error)

func SetsockoptTimeval(fd, level, opt int, tv *Timeval) (err error)

func Setuid(uid int) (err error)

// SlicePtrFromStrings converts a slice of strings to a slice of
// pointers to NUL-terminated byte arrays. If any string contains
// a NUL byte, it returns (nil, EINVAL).
func SlicePtrFromStrings(ss []string) ([]*byte, error)

func Socket(domain, typ, proto int) (fd int, err error)

func Socketpair(domain, typ, proto int) (fd [2]int, err error)

// StartProcess wraps ForkExec for package os.
func StartProcess(argv0 string, argv []string, attr *ProcAttr) (pid int, handle uintptr, err error)

func Stat(path string, s *Stat_t) (err error)

// StringBytePtr returns a pointer to a NUL-terminated array of bytes.
// If s contains a NUL byte this function panics instead of returning
// an error.
//
// Deprecated: Use BytePtrFromString instead.
func StringBytePtr(s string) *byte

// StringByteSlice converts a string to a NUL-terminated []byte,
// If s contains a NUL byte this function panics instead of
// returning an error.
//
// Deprecated: Use ByteSliceFromString instead.
func StringByteSlice(s string) []byte

// StringSlicePtr converts a slice of strings to a slice of pointers
// to NUL-terminated byte arrays. If any string contains a NUL byte
// this function panics instead of returning an error.
//
// Deprecated: Use SlicePtrFromStrings instead.
func StringSlicePtr(ss []string) []*byte

func Symlink(oldpath string, newpath string) (err error)

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)

func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)

func Time(t *Time_t) (tt Time_t, err error)

func TimespecToNsec(ts Timespec) int64

func TimevalToNsec(tv Timeval) int64

// UnixCredentials encodes credentials into a socket control message
// for sending to another process. This can be used for
// authentication.
func UnixCredentials(ucred *Ucred) []byte

// UnixRights encodes a set of open file descriptors into a socket
// control message for sending to another process.
func UnixRights(fds ...int) []byte

func Unlink(path string) error

func Unlinkat(dirfd int, path string) error

func Unsetenv(key string) error

func Utimes(path string, tv []Timeval) (err error)

func UtimesNano(path string, ts []Timespec) (err error)

func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error)

func Write(fd int, p []byte) (n int, err error)

func (cmsg *Cmsghdr) SetLen(length int)

func (iov *Iovec) SetLen(length int)

func (msghdr *Msghdr) SetControllen(length int)

func (r *PtraceRegs) PC() uint64

func (r *PtraceRegs) SetPC(pc uint64)

func (ts *Timespec) Nano() int64

func (ts *Timespec) Unix() (sec int64, nsec int64)

func (tv *Timeval) Nano() int64

func (tv *Timeval) Unix() (sec int64, nsec int64)

func (e Errno) Error() string

func (e Errno) Temporary() bool

func (e Errno) Timeout() bool

func (s Signal) Signal()

func (s Signal) String() string

func (w WaitStatus) Continued() bool

func (w WaitStatus) CoreDump() bool

func (w WaitStatus) ExitStatus() int

func (w WaitStatus) Exited() bool

func (w WaitStatus) Signal() Signal

func (w WaitStatus) Signaled() bool

func (w WaitStatus) StopSignal() Signal

func (w WaitStatus) Stopped() bool

func (w WaitStatus) TrapCause() int

