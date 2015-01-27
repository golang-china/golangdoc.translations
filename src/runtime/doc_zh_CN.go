// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package runtime contains operations that interact with Go's runtime system, such
// as functions to control goroutines. It also includes the low-level type
// information used by the reflect package; see reflect's documentation for the
// programmable interface to the run-time type system.
//
//
// Environment Variables
//
// The following environment variables ($name or %name%, depending on the host
// operating system) control the run-time behavior of Go programs. The meanings and
// use may change from release to release.
//
// The GOGC variable sets the initial garbage collection target percentage. A
// collection is triggered when the ratio of freshly allocated data to live data
// remaining after the previous collection reaches this percentage. The default is
// GOGC=100. Setting GOGC=off disables the garbage collector entirely. The
// runtime/debug package's SetGCPercent function allows changing this percentage at
// run time. See http://golang.org/pkg/runtime/debug/#SetGCPercent.
//
// The GODEBUG variable controls debug output from the runtime. GODEBUG value is a
// comma-separated list of name=val pairs. Supported names are:
//
//	allocfreetrace: setting allocfreetrace=1 causes every allocation to be
//	profiled and a stack trace printed on each object's allocation and free.
//
//	efence: setting efence=1 causes the allocator to run in a mode
//	where each object is allocated on a unique page and addresses are
//	never recycled.
//
//	gctrace: setting gctrace=1 causes the garbage collector to emit a single line to standard
//	error at each collection, summarizing the amount of memory collected and the
//	length of the pause. Setting gctrace=2 emits the same summary but also
//	repeats each collection.
//
//	gcdead: setting gcdead=1 causes the garbage collector to clobber all stack slots
//	that it thinks are dead.
//
//	invalidptr: defaults to invalidptr=1, causing the garbage collector and stack
//	copier to crash the program if an invalid pointer value (for example, 1)
//	is found in a pointer-typed location. Setting invalidptr=0 disables this check.
//	This should only be used as a temporary workaround to diagnose buggy code.
//	The real fix is to not store integers in pointer-typed locations.
//
//	scheddetail: setting schedtrace=X and scheddetail=1 causes the scheduler to emit
//	detailed multiline info every X milliseconds, describing state of the scheduler,
//	processors, threads and goroutines.
//
//	schedtrace: setting schedtrace=X causes the scheduler to emit a single line to standard
//	error every X milliseconds, summarizing the scheduler state.
//
//	scavenge: scavenge=1 enables debugging mode of heap scavenger.
//
// The GOMAXPROCS variable limits the number of operating system threads that can
// execute user-level Go code simultaneously. There is no limit to the number of
// threads that can be blocked in system calls on behalf of Go code; those do not
// count against the GOMAXPROCS limit. This package's GOMAXPROCS function queries
// and changes the limit.
//
// The GOTRACEBACK variable controls the amount of output generated when a Go
// program fails due to an unrecovered panic or an unexpected runtime condition. By
// default, a failure prints a stack trace for every extant goroutine, eliding
// functions internal to the run-time system, and then exits with exit code 2. If
// GOTRACEBACK=0, the per-goroutine stack traces are omitted entirely. If
// GOTRACEBACK=1, the default behavior is used. If GOTRACEBACK=2, the per-goroutine
// stack traces include run-time functions. If GOTRACEBACK=crash, the per-goroutine
// stack traces include run-time functions, and if possible the program crashes in
// an operating-specific manner instead of exiting. For example, on Unix systems,
// the program raises SIGABRT to trigger a core dump.
//
// The GOARCH, GOOS, GOPATH, and GOROOT environment variables complete the set of
// Go environment variables. They influence the building of Go programs (see
// http://golang.org/cmd/go and http://golang.org/pkg/go/build). GOARCH, GOOS, and
// GOROOT are recorded at compile time and made available by constants or functions
// in this package, but they do not influence the execution of the run-time system.

// TODO(osc): 需更新 runtime
// 包含与Go的运行时系统进行交互的操作，例如用于控制Go程的函数. 它也包括用于 reflect
// 包的底层类型信息；运行时类型系统的可编程接口见 reflect 文档。
//
// 环境变量
//
// 以下环境变量（$name 或 %name%,
// 取决于宿主操作系统）控制了Go程序的运行时行为。
// 其意义与使用方法在发行版之间可能有所不同。
//
// GOGC
// 变量用于设置初始垃圾回收的目标百分比。从上次回收后开始，当新分配数据的比例占到剩余实时数据的此百分比时，
// 就会再次触发回收。默认为
// GOGC=100。要完全关闭垃圾回收器，需设置 GOGC=off。runtime/debug 包的 SetGCPercent
// 函数允许在运行时更改此百分比。 详见 http://zh.golanger.com/pkg/runtime/debug/#SetGCPercent。
//
// GOGCTRACE
// 变量用于控制来自垃圾回收器的调试输出。设置 GOGCTRACE=1 会使垃圾回收器发出
// 每一次回收所产生的单行标准错误输出、概述回收的内存量以及暂停的时长。设置 GOGCTRACE=2
// 不仅会发出同样的概述，还会重复每一次回收。
//
// GOMAXPROCS
// 变量用于限制可同时执行的用户级Go代码所产生的操作系统线程数。对于Go代码所代表的系统调用而言，
// 可被阻塞的线程则没有限制；它们不计入 GOMAXPROCS 的限制。本包中的 GOMAXPROCS
// 函数可查询并更改此限制。
//
// GOTRACEBACK
// 用于控制因未恢复的恐慌或意外的运行时状况导致Go程序运行失败时所产生的输出量。
// 默认情况下，失败会为每个现有的Go程打印出栈跟踪，省略运行时系统的内部函数，并以退出码 2 退出。 若
// GOTRACEBACK=0，则每个Go程的栈跟踪都会完全省略。 若
// GOTRACEBACK=1，则采用默认的行为。 若
// GOTRACEBACK=2，则每个Go程的栈跟踪，包括运行时函数都会输出。 若
// GOTRACEBACK=crash，则每个Go程的栈跟踪，包括运行时函数，都会输出，
// 此外程序可能以操作系统特定的方式崩溃而非退出。例如，在Unix系统上，程序会发出 SIGABRT
// 信号，从而触发内核转储。
//
// GOARCH、GOOS、GOPATH 和 GOROOT
// 环境变量均为Go的环境变量。它们影响了Go程序的构建 （详见 http://golang.org/cmd/go 和
// http://golang.org/pkg/go/build）。 GOARCH、GOOS 和 GOROOT
// 会在编译时被记录，并使该包中的常量或函数变得可用，
// 但它们并不影响运行时系统的执行。
//
// 公共的竞争检测API，当且仅当使用 -race 构建时才会出现。
package runtime

const (
	O_RDONLY  = C.O_RDONLY
	O_CLOEXEC = C.O_CLOEXEC
)

const (
	EINTR  = C.EINTR
	EAGAIN = C.EAGAIN
	ENOMEM = C.ENOMEM

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANONYMOUS
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_DONTNEED = C.MADV_DONTNEED

	SA_RESTART  = C.SA_RESTART
	SA_ONSTACK  = C.SA_ONSTACK
	SA_RESTORER = C.SA_RESTORER
	SA_SIGINFO  = C.SA_SIGINFO

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGBUS    = C.SIGBUS
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGUSR1   = C.SIGUSR1
	SIGSEGV   = C.SIGSEGV
	SIGUSR2   = C.SIGUSR2
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGSTKFLT = C.SIGSTKFLT
	SIGCHLD   = C.SIGCHLD
	SIGCONT   = C.SIGCONT
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGURG    = C.SIGURG
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGIO     = C.SIGIO
	SIGPWR    = C.SIGPWR
	SIGSYS    = C.SIGSYS

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	O_RDONLY  = C.O_RDONLY
	O_CLOEXEC = C.O_CLOEXEC

	EPOLLIN       = C.POLLIN
	EPOLLOUT      = C.POLLOUT
	EPOLLERR      = C.POLLERR
	EPOLLHUP      = C.POLLHUP
	EPOLLRDHUP    = C.POLLRDHUP
	EPOLLET       = C.EPOLLET
	EPOLL_CLOEXEC = C.EPOLL_CLOEXEC
	EPOLL_CTL_ADD = C.EPOLL_CTL_ADD
	EPOLL_CTL_DEL = C.EPOLL_CTL_DEL
	EPOLL_CTL_MOD = C.EPOLL_CTL_MOD
)

const (
	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANONYMOUS
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_DONTNEED = C.MADV_DONTNEED

	SA_RESTART  = C.SA_RESTART
	SA_ONSTACK  = C.SA_ONSTACK
	SA_RESTORER = C.SA_RESTORER
	SA_SIGINFO  = C.SA_SIGINFO

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGBUS    = C.SIGBUS
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGUSR1   = C.SIGUSR1
	SIGSEGV   = C.SIGSEGV
	SIGUSR2   = C.SIGUSR2
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGSTKFLT = C.SIGSTKFLT
	SIGCHLD   = C.SIGCHLD
	SIGCONT   = C.SIGCONT
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGURG    = C.SIGURG
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGIO     = C.SIGIO
	SIGPWR    = C.SIGPWR
	SIGSYS    = C.SIGSYS

	FPE_INTDIV = C.FPE_INTDIV & 0xFFFF
	FPE_INTOVF = C.FPE_INTOVF & 0xFFFF
	FPE_FLTDIV = C.FPE_FLTDIV & 0xFFFF
	FPE_FLTOVF = C.FPE_FLTOVF & 0xFFFF
	FPE_FLTUND = C.FPE_FLTUND & 0xFFFF
	FPE_FLTRES = C.FPE_FLTRES & 0xFFFF
	FPE_FLTINV = C.FPE_FLTINV & 0xFFFF
	FPE_FLTSUB = C.FPE_FLTSUB & 0xFFFF

	BUS_ADRALN = C.BUS_ADRALN & 0xFFFF
	BUS_ADRERR = C.BUS_ADRERR & 0xFFFF
	BUS_OBJERR = C.BUS_OBJERR & 0xFFFF

	SEGV_MAPERR = C.SEGV_MAPERR & 0xFFFF
	SEGV_ACCERR = C.SEGV_ACCERR & 0xFFFF

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_PROF    = C.ITIMER_PROF
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
)

const (
	EINTR  = C.EINTR
	EFAULT = C.EFAULT

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANON
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_DONTNEED = C.MADV_DONTNEED
	MADV_FREE     = C.MADV_FREE

	MACH_MSG_TYPE_MOVE_RECEIVE   = C.MACH_MSG_TYPE_MOVE_RECEIVE
	MACH_MSG_TYPE_MOVE_SEND      = C.MACH_MSG_TYPE_MOVE_SEND
	MACH_MSG_TYPE_MOVE_SEND_ONCE = C.MACH_MSG_TYPE_MOVE_SEND_ONCE
	MACH_MSG_TYPE_COPY_SEND      = C.MACH_MSG_TYPE_COPY_SEND
	MACH_MSG_TYPE_MAKE_SEND      = C.MACH_MSG_TYPE_MAKE_SEND
	MACH_MSG_TYPE_MAKE_SEND_ONCE = C.MACH_MSG_TYPE_MAKE_SEND_ONCE
	MACH_MSG_TYPE_COPY_RECEIVE   = C.MACH_MSG_TYPE_COPY_RECEIVE

	MACH_MSG_PORT_DESCRIPTOR         = C.MACH_MSG_PORT_DESCRIPTOR
	MACH_MSG_OOL_DESCRIPTOR          = C.MACH_MSG_OOL_DESCRIPTOR
	MACH_MSG_OOL_PORTS_DESCRIPTOR    = C.MACH_MSG_OOL_PORTS_DESCRIPTOR
	MACH_MSG_OOL_VOLATILE_DESCRIPTOR = C.MACH_MSG_OOL_VOLATILE_DESCRIPTOR

	MACH_MSGH_BITS_COMPLEX = C.MACH_MSGH_BITS_COMPLEX

	MACH_SEND_MSG  = C.MACH_SEND_MSG
	MACH_RCV_MSG   = C.MACH_RCV_MSG
	MACH_RCV_LARGE = C.MACH_RCV_LARGE

	MACH_SEND_TIMEOUT   = C.MACH_SEND_TIMEOUT
	MACH_SEND_INTERRUPT = C.MACH_SEND_INTERRUPT
	MACH_SEND_ALWAYS    = C.MACH_SEND_ALWAYS
	MACH_SEND_TRAILER   = C.MACH_SEND_TRAILER
	MACH_RCV_TIMEOUT    = C.MACH_RCV_TIMEOUT
	MACH_RCV_NOTIFY     = C.MACH_RCV_NOTIFY
	MACH_RCV_INTERRUPT  = C.MACH_RCV_INTERRUPT
	MACH_RCV_OVERWRITE  = C.MACH_RCV_OVERWRITE

	NDR_PROTOCOL_2_0      = C.NDR_PROTOCOL_2_0
	NDR_INT_BIG_ENDIAN    = C.NDR_INT_BIG_ENDIAN
	NDR_INT_LITTLE_ENDIAN = C.NDR_INT_LITTLE_ENDIAN
	NDR_FLOAT_IEEE        = C.NDR_FLOAT_IEEE
	NDR_CHAR_ASCII        = C.NDR_CHAR_ASCII

	SA_SIGINFO   = C.SA_SIGINFO
	SA_RESTART   = C.SA_RESTART
	SA_ONSTACK   = C.SA_ONSTACK
	SA_USERTRAMP = C.SA_USERTRAMP
	SA_64REGSET  = C.SA_64REGSET

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGEMT    = C.SIGEMT
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGBUS    = C.SIGBUS
	SIGSEGV   = C.SIGSEGV
	SIGSYS    = C.SIGSYS
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGTERM   = C.SIGTERM
	SIGURG    = C.SIGURG
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGCONT   = C.SIGCONT
	SIGCHLD   = C.SIGCHLD
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGIO     = C.SIGIO
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGINFO   = C.SIGINFO
	SIGUSR1   = C.SIGUSR1
	SIGUSR2   = C.SIGUSR2

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	EV_ADD       = C.EV_ADD
	EV_DELETE    = C.EV_DELETE
	EV_CLEAR     = C.EV_CLEAR
	EV_RECEIPT   = C.EV_RECEIPT
	EV_ERROR     = C.EV_ERROR
	EVFILT_READ  = C.EVFILT_READ
	EVFILT_WRITE = C.EVFILT_WRITE
)

const (
	EINTR  = C.EINTR
	EFAULT = C.EFAULT
	EBUSY  = C.EBUSY
	EAGAIN = C.EAGAIN

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANON
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_FREE = C.MADV_FREE

	SA_SIGINFO = C.SA_SIGINFO
	SA_RESTART = C.SA_RESTART
	SA_ONSTACK = C.SA_ONSTACK

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGEMT    = C.SIGEMT
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGBUS    = C.SIGBUS
	SIGSEGV   = C.SIGSEGV
	SIGSYS    = C.SIGSYS
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGTERM   = C.SIGTERM
	SIGURG    = C.SIGURG
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGCONT   = C.SIGCONT
	SIGCHLD   = C.SIGCHLD
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGIO     = C.SIGIO
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGINFO   = C.SIGINFO
	SIGUSR1   = C.SIGUSR1
	SIGUSR2   = C.SIGUSR2

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	EV_ADD       = C.EV_ADD
	EV_DELETE    = C.EV_DELETE
	EV_CLEAR     = C.EV_CLEAR
	EV_ERROR     = C.EV_ERROR
	EVFILT_READ  = C.EVFILT_READ
	EVFILT_WRITE = C.EVFILT_WRITE
)

const (
	EINTR  = C.EINTR
	EFAULT = C.EFAULT

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANON
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_FREE = C.MADV_FREE

	SA_SIGINFO = C.SA_SIGINFO
	SA_RESTART = C.SA_RESTART
	SA_ONSTACK = C.SA_ONSTACK

	UMTX_OP_WAIT_UINT         = C.UMTX_OP_WAIT_UINT
	UMTX_OP_WAIT_UINT_PRIVATE = C.UMTX_OP_WAIT_UINT_PRIVATE
	UMTX_OP_WAKE              = C.UMTX_OP_WAKE
	UMTX_OP_WAKE_PRIVATE      = C.UMTX_OP_WAKE_PRIVATE

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGEMT    = C.SIGEMT
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGBUS    = C.SIGBUS
	SIGSEGV   = C.SIGSEGV
	SIGSYS    = C.SIGSYS
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGTERM   = C.SIGTERM
	SIGURG    = C.SIGURG
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGCONT   = C.SIGCONT
	SIGCHLD   = C.SIGCHLD
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGIO     = C.SIGIO
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGINFO   = C.SIGINFO
	SIGUSR1   = C.SIGUSR1
	SIGUSR2   = C.SIGUSR2

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	EV_ADD       = C.EV_ADD
	EV_DELETE    = C.EV_DELETE
	EV_CLEAR     = C.EV_CLEAR
	EV_RECEIPT   = C.EV_RECEIPT
	EV_ERROR     = C.EV_ERROR
	EVFILT_READ  = C.EVFILT_READ
	EVFILT_WRITE = C.EVFILT_WRITE
)

const (
	EINTR  = C.EINTR
	EAGAIN = C.EAGAIN
	ENOMEM = C.ENOMEM

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANONYMOUS
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_DONTNEED = C.MADV_DONTNEED

	SA_RESTART  = C.SA_RESTART
	SA_ONSTACK  = C.SA_ONSTACK
	SA_RESTORER = C.SA_RESTORER
	SA_SIGINFO  = C.SA_SIGINFO

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGBUS    = C.SIGBUS
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGUSR1   = C.SIGUSR1
	SIGSEGV   = C.SIGSEGV
	SIGUSR2   = C.SIGUSR2
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGSTKFLT = C.SIGSTKFLT
	SIGCHLD   = C.SIGCHLD
	SIGCONT   = C.SIGCONT
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGURG    = C.SIGURG
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGIO     = C.SIGIO
	SIGPWR    = C.SIGPWR
	SIGSYS    = C.SIGSYS

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	EPOLLIN       = C.POLLIN
	EPOLLOUT      = C.POLLOUT
	EPOLLERR      = C.POLLERR
	EPOLLHUP      = C.POLLHUP
	EPOLLRDHUP    = C.POLLRDHUP
	EPOLLET       = C.EPOLLET
	EPOLL_CLOEXEC = C.EPOLL_CLOEXEC
	EPOLL_CTL_ADD = C.EPOLL_CTL_ADD
	EPOLL_CTL_DEL = C.EPOLL_CTL_DEL
	EPOLL_CTL_MOD = C.EPOLL_CTL_MOD
)

const (
	EINTR  = C.EINTR
	EFAULT = C.EFAULT

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANON
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_FREE = C.MADV_FREE

	SA_SIGINFO = C.SA_SIGINFO
	SA_RESTART = C.SA_RESTART
	SA_ONSTACK = C.SA_ONSTACK

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGEMT    = C.SIGEMT
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGBUS    = C.SIGBUS
	SIGSEGV   = C.SIGSEGV
	SIGSYS    = C.SIGSYS
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGTERM   = C.SIGTERM
	SIGURG    = C.SIGURG
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGCONT   = C.SIGCONT
	SIGCHLD   = C.SIGCHLD
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGIO     = C.SIGIO
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGINFO   = C.SIGINFO
	SIGUSR1   = C.SIGUSR1
	SIGUSR2   = C.SIGUSR2

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	EV_ADD       = C.EV_ADD
	EV_DELETE    = C.EV_DELETE
	EV_CLEAR     = C.EV_CLEAR
	EV_RECEIPT   = 0
	EV_ERROR     = C.EV_ERROR
	EVFILT_READ  = C.EVFILT_READ
	EVFILT_WRITE = C.EVFILT_WRITE
)

const (
	REG_GS     = C._REG_GS
	REG_FS     = C._REG_FS
	REG_ES     = C._REG_ES
	REG_DS     = C._REG_DS
	REG_EDI    = C._REG_EDI
	REG_ESI    = C._REG_ESI
	REG_EBP    = C._REG_EBP
	REG_ESP    = C._REG_ESP
	REG_EBX    = C._REG_EBX
	REG_EDX    = C._REG_EDX
	REG_ECX    = C._REG_ECX
	REG_EAX    = C._REG_EAX
	REG_TRAPNO = C._REG_TRAPNO
	REG_ERR    = C._REG_ERR
	REG_EIP    = C._REG_EIP
	REG_CS     = C._REG_CS
	REG_EFL    = C._REG_EFL
	REG_UESP   = C._REG_UESP
	REG_SS     = C._REG_SS
)

const (
	REG_RDI    = C._REG_RDI
	REG_RSI    = C._REG_RSI
	REG_RDX    = C._REG_RDX
	REG_RCX    = C._REG_RCX
	REG_R8     = C._REG_R8
	REG_R9     = C._REG_R9
	REG_R10    = C._REG_R10
	REG_R11    = C._REG_R11
	REG_R12    = C._REG_R12
	REG_R13    = C._REG_R13
	REG_R14    = C._REG_R14
	REG_R15    = C._REG_R15
	REG_RBP    = C._REG_RBP
	REG_RBX    = C._REG_RBX
	REG_RAX    = C._REG_RAX
	REG_GS     = C._REG_GS
	REG_FS     = C._REG_FS
	REG_ES     = C._REG_ES
	REG_DS     = C._REG_DS
	REG_TRAPNO = C._REG_TRAPNO
	REG_ERR    = C._REG_ERR
	REG_RIP    = C._REG_RIP
	REG_CS     = C._REG_CS
	REG_RFLAGS = C._REG_RFLAGS
	REG_RSP    = C._REG_RSP
	REG_SS     = C._REG_SS
)

const (
	REG_R0   = C._REG_R0
	REG_R1   = C._REG_R1
	REG_R2   = C._REG_R2
	REG_R3   = C._REG_R3
	REG_R4   = C._REG_R4
	REG_R5   = C._REG_R5
	REG_R6   = C._REG_R6
	REG_R7   = C._REG_R7
	REG_R8   = C._REG_R8
	REG_R9   = C._REG_R9
	REG_R10  = C._REG_R10
	REG_R11  = C._REG_R11
	REG_R12  = C._REG_R12
	REG_R13  = C._REG_R13
	REG_R14  = C._REG_R14
	REG_R15  = C._REG_R15
	REG_CPSR = C._REG_CPSR
)

const (
	EINTR  = C.EINTR
	EFAULT = C.EFAULT

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANON
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_FREE = C.MADV_FREE

	SA_SIGINFO = C.SA_SIGINFO
	SA_RESTART = C.SA_RESTART
	SA_ONSTACK = C.SA_ONSTACK

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGEMT    = C.SIGEMT
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGBUS    = C.SIGBUS
	SIGSEGV   = C.SIGSEGV
	SIGSYS    = C.SIGSYS
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGTERM   = C.SIGTERM
	SIGURG    = C.SIGURG
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGCONT   = C.SIGCONT
	SIGCHLD   = C.SIGCHLD
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGIO     = C.SIGIO
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGINFO   = C.SIGINFO
	SIGUSR1   = C.SIGUSR1
	SIGUSR2   = C.SIGUSR2

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	EV_ADD       = C.EV_ADD
	EV_DELETE    = C.EV_DELETE
	EV_CLEAR     = C.EV_CLEAR
	EV_ERROR     = C.EV_ERROR
	EVFILT_READ  = C.EVFILT_READ
	EVFILT_WRITE = C.EVFILT_WRITE
)

const (
	EINTR       = C.EINTR
	EBADF       = C.EBADF
	EFAULT      = C.EFAULT
	EAGAIN      = C.EAGAIN
	ETIMEDOUT   = C.ETIMEDOUT
	EWOULDBLOCK = C.EWOULDBLOCK
	EINPROGRESS = C.EINPROGRESS

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANON
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_FREE = C.MADV_FREE

	SA_SIGINFO = C.SA_SIGINFO
	SA_RESTART = C.SA_RESTART
	SA_ONSTACK = C.SA_ONSTACK

	SIGHUP    = C.SIGHUP
	SIGINT    = C.SIGINT
	SIGQUIT   = C.SIGQUIT
	SIGILL    = C.SIGILL
	SIGTRAP   = C.SIGTRAP
	SIGABRT   = C.SIGABRT
	SIGEMT    = C.SIGEMT
	SIGFPE    = C.SIGFPE
	SIGKILL   = C.SIGKILL
	SIGBUS    = C.SIGBUS
	SIGSEGV   = C.SIGSEGV
	SIGSYS    = C.SIGSYS
	SIGPIPE   = C.SIGPIPE
	SIGALRM   = C.SIGALRM
	SIGTERM   = C.SIGTERM
	SIGURG    = C.SIGURG
	SIGSTOP   = C.SIGSTOP
	SIGTSTP   = C.SIGTSTP
	SIGCONT   = C.SIGCONT
	SIGCHLD   = C.SIGCHLD
	SIGTTIN   = C.SIGTTIN
	SIGTTOU   = C.SIGTTOU
	SIGIO     = C.SIGIO
	SIGXCPU   = C.SIGXCPU
	SIGXFSZ   = C.SIGXFSZ
	SIGVTALRM = C.SIGVTALRM
	SIGPROF   = C.SIGPROF
	SIGWINCH  = C.SIGWINCH
	SIGUSR1   = C.SIGUSR1
	SIGUSR2   = C.SIGUSR2

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	PTHREAD_CREATE_DETACHED = C.PTHREAD_CREATE_DETACHED

	FORK_NOSIGCHLD = C.FORK_NOSIGCHLD
	FORK_WAITPID   = C.FORK_WAITPID

	MAXHOSTNAMELEN = C.MAXHOSTNAMELEN

	O_NONBLOCK = C.O_NONBLOCK
	FD_CLOEXEC = C.FD_CLOEXEC
	F_GETFL    = C.F_GETFL
	F_SETFL    = C.F_SETFL
	F_SETFD    = C.F_SETFD

	POLLIN  = C.POLLIN
	POLLOUT = C.POLLOUT
	POLLHUP = C.POLLHUP
	POLLERR = C.POLLERR

	PORT_SOURCE_FD = C.PORT_SOURCE_FD
)

const (
	REG_RDI    = C.REG_RDI
	REG_RSI    = C.REG_RSI
	REG_RDX    = C.REG_RDX
	REG_RCX    = C.REG_RCX
	REG_R8     = C.REG_R8
	REG_R9     = C.REG_R9
	REG_R10    = C.REG_R10
	REG_R11    = C.REG_R11
	REG_R12    = C.REG_R12
	REG_R13    = C.REG_R13
	REG_R14    = C.REG_R14
	REG_R15    = C.REG_R15
	REG_RBP    = C.REG_RBP
	REG_RBX    = C.REG_RBX
	REG_RAX    = C.REG_RAX
	REG_GS     = C.REG_GS
	REG_FS     = C.REG_FS
	REG_ES     = C.REG_ES
	REG_DS     = C.REG_DS
	REG_TRAPNO = C.REG_TRAPNO
	REG_ERR    = C.REG_ERR
	REG_RIP    = C.REG_RIP
	REG_CS     = C.REG_CS
	REG_RFLAGS = C.REG_RFL
	REG_RSP    = C.REG_RSP
	REG_SS     = C.REG_SS
)

const (
	PROT_NONE  = 0
	PROT_READ  = 1
	PROT_WRITE = 2
	PROT_EXEC  = 4

	MAP_ANON    = 1
	MAP_PRIVATE = 2

	DUPLICATE_SAME_ACCESS   = C.DUPLICATE_SAME_ACCESS
	THREAD_PRIORITY_HIGHEST = C.THREAD_PRIORITY_HIGHEST

	SIGPROF          = 0 // dummy value for badsignal
	SIGINT           = C.SIGINT
	CTRL_C_EVENT     = C.CTRL_C_EVENT
	CTRL_BREAK_EVENT = C.CTRL_BREAK_EVENT

	CONTEXT_CONTROL = C.CONTEXT_CONTROL
	CONTEXT_FULL    = C.CONTEXT_FULL

	EXCEPTION_ACCESS_VIOLATION     = C.STATUS_ACCESS_VIOLATION
	EXCEPTION_BREAKPOINT           = C.STATUS_BREAKPOINT
	EXCEPTION_FLT_DENORMAL_OPERAND = C.STATUS_FLOAT_DENORMAL_OPERAND
	EXCEPTION_FLT_DIVIDE_BY_ZERO   = C.STATUS_FLOAT_DIVIDE_BY_ZERO
	EXCEPTION_FLT_INEXACT_RESULT   = C.STATUS_FLOAT_INEXACT_RESULT
	EXCEPTION_FLT_OVERFLOW         = C.STATUS_FLOAT_OVERFLOW
	EXCEPTION_FLT_UNDERFLOW        = C.STATUS_FLOAT_UNDERFLOW
	EXCEPTION_INT_DIVIDE_BY_ZERO   = C.STATUS_INTEGER_DIVIDE_BY_ZERO
	EXCEPTION_INT_OVERFLOW         = C.STATUS_INTEGER_OVERFLOW

	INFINITE     = C.INFINITE
	WAIT_TIMEOUT = C.WAIT_TIMEOUT

	EXCEPTION_CONTINUE_EXECUTION = C.EXCEPTION_CONTINUE_EXECUTION
	EXCEPTION_CONTINUE_SEARCH    = C.EXCEPTION_CONTINUE_SEARCH
)

// Compiler is the name of the compiler toolchain that built the running binary.
// Known toolchains are:
//
//	gc      The 5g/6g/8g compiler suite at code.google.com/p/go.
//	gccgo   The gccgo front end, part of the GCC compiler suite.

// Compiler
// 为构建了可运行二进制文件的编译工具链。已知的工具链为：
//
//	go       code.google.com/p/go 上的 5g/6g/8g 编译器套件。
//	gccgo    gccgo前端，GCC编译器条件的一部分。
const Compiler = "gc"

// GOARCH is the running program's architecture target: 386, amd64, or arm.

// GOARCH 为所运行程序的目标架构： 386、amd64 或 arm。
const GOARCH string = theGoarch

// GOOS is the running program's operating system target: one of darwin, freebsd,
// linux, and so on.

// GOOS 为所运行程序的目标操作系统： darwin、freebsd或linux等等。
const GOOS string = theGoos

const (
	_ selectDir = iota
)

// MemProfileRate controls the fraction of memory allocations that are recorded and
// reported in the memory profile. The profiler aims to sample an average of one
// allocation per MemProfileRate bytes allocated.
//
// To include every allocated block in the profile, set MemProfileRate to 1. To
// turn off profiling entirely, set MemProfileRate to 0.
//
// The tools that process the memory profiles assume that the profile rate is
// constant across the lifetime of the program and equal to the current value.
// Programs that change the memory profiling rate should do so just once, as early
// as possible in the execution of the program (for example, at the beginning of
// main).
var MemProfileRate int = 512 * 1024

// BlockProfile returns n, the number of records in the current blocking profile.
// If len(p) >= n, BlockProfile copies the profile into p and returns n, true. If
// len(p) < n, BlockProfile does not change p and returns n, false.
//
// Most clients should use the runtime/pprof package or the testing package's
// -test.blockprofile flag instead of calling BlockProfile directly.
func BlockProfile(p []BlockProfileRecord) (n int, ok bool)

// Breakpoint executes a breakpoint trap.
func Breakpoint()

// CPUProfile returns the next chunk of binary CPU profiling stack trace data,
// blocking until data is available. If profiling is turned off and all the profile
// data accumulated while it was on has been returned, CPUProfile returns nil. The
// caller must save the returned data before calling CPUProfile again.
//
// Most clients should use the runtime/pprof package or the testing package's
// -test.cpuprofile flag instead of calling CPUProfile directly.
func CPUProfile() []byte

// Caller reports file and line number information about function invocations on
// the calling goroutine's stack. The argument skip is the number of stack frames
// to ascend, with 0 identifying the caller of Caller. (For historical reasons the
// meaning of skip differs between Caller and Callers.) The return values report
// the program counter, file name, and line number within the file of the
// corresponding call. The boolean ok is false if it was not possible to recover
// the information.

// Caller
// 报告关于调用Go程的栈上的函数调用的文件和行号信息。 实参 skip
// 为占用的栈帧数，若为0则表示 Caller 的调用者。（由于历史原因，skip 的意思在 Caller 和 Callers
// 中并不相同。）返回值报告程序计数器，
// 文件名及对应调用的文件中的行号。若无法获得信息，布尔值 ok 即为 false。
func Caller(skip int) (pc uintptr, file string, line int, ok bool)

// Callers fills the slice pc with the return program counters of function
// invocations on the calling goroutine's stack. The argument skip is the number of
// stack frames to skip before recording in pc, with 0 identifying the frame for
// Callers itself and 1 identifying the caller of Callers. It returns the number of
// entries written to pc.
//
// Note that since each slice entry pc[i] is a return program counter, looking up
// the file and line for pc[i] (for example, using (*Func).FileLine) will return
// the file and line number of the instruction immediately following the call. To
// look up the file and line number of the call itself, use pc[i]-1. As an
// exception to this rule, if pc[i-1] corresponds to the function runtime.sigpanic,
// then pc[i] is the program counter of a faulting instruction and should be used
// without any subtraction.

// Callers
// 把调用它的Go程栈上函数请求的返回程序计数器填充到切片 pc 中。 实参 skip 为开始在 pc
// 中记录之前所要跳过的栈帧数，若为 0 则表示 Callers 自身的栈帧， 若为 1 则表示 Callers
// 的调用者。它返回写入到 pc 中的项数。
//
// 注意，由于每个切片项 pc[i]
// 都是一个返回程序计数器，因此查找 pc[i] 的文件和行（例如，使用
// (*Func).FileLine）将会在该调用之后立即返回该指令所在的文件和行号。
// 要查找该调用本身所在的文件和行号，请使用 pc[i]-1。此规则的一个例外是，若 pc[i-1] 对应于函数
// runtime.sigpanic，那么 pc[i]
// 就是失败指令的程序计数器，因此应当不通过任何减法来使用。
func Callers(skip int, pc []uintptr) int

// GC runs a garbage collection.

// GC 运行一次垃圾回收。
func GC()

// GOMAXPROCS sets the maximum number of CPUs that can be executing simultaneously
// and returns the previous setting. If n < 1, it does not change the current
// setting. The number of logical CPUs on the local machine can be queried with
// NumCPU. This call will go away when the scheduler improves.

// GOMAXPROCS
// 设置可同时使用执行的最大CPU数，并返回先前的设置。 若 n <
// 1，它就不会更改当前设置。本地机器的逻辑CPU数可通过 NumCPU 查询。
// 当调度器改进后，此调用将会消失。
func GOMAXPROCS(n int) int

// GOROOT returns the root of the Go tree. It uses the GOROOT environment variable,
// if set, or else the root used during the Go build.

// GOROOT 返回Go目录树的根目录。
// 若设置了GOROOT环境变量，就会使用它，否则就会将Go的构建目录作为根目录
func GOROOT() string

// Goexit terminates the goroutine that calls it. No other goroutine is affected.
// Goexit runs all deferred calls before terminating the goroutine. Because Goexit
// is not panic, however, any recover calls in those deferred functions will return
// nil.
//
// Calling Goexit from the main goroutine terminates that goroutine without func
// main returning. Since func main has not returned, the program continues
// execution of other goroutines. If all other goroutines exit, the program
// crashes.
func Goexit()

// GoroutineProfile returns n, the number of records in the active goroutine stack
// profile. If len(p) >= n, GoroutineProfile copies the profile into p and returns
// n, true. If len(p) < n, GoroutineProfile does not change p and returns n, false.
//
// Most clients should use the runtime/pprof package instead of calling
// GoroutineProfile directly.
func GoroutineProfile(p []StackRecord) (n int, ok bool)

// Gosched yields the processor, allowing other goroutines to run. It does not
// suspend the current goroutine, so execution resumes automatically.
func Gosched()

// LockOSThread wires the calling goroutine to its current operating system thread.
// Until the calling goroutine exits or calls UnlockOSThread, it will always
// execute in that thread, and no other goroutine can.
func LockOSThread()

// MemProfile returns n, the number of records in the current memory profile. If
// len(p) >= n, MemProfile copies the profile into p and returns n, true. If len(p)
// < n, MemProfile does not change p and returns n, false.
//
// If inuseZero is true, the profile includes allocation records where r.AllocBytes
// > 0 but r.AllocBytes == r.FreeBytes. These are sites where memory was allocated,
// but it has all been released back to the runtime.
//
// Most clients should use the runtime/pprof package or the testing package's
// -test.memprofile flag instead of calling MemProfile directly.
func MemProfile(p []MemProfileRecord, inuseZero bool) (n int, ok bool)

// NumCPU returns the number of logical CPUs on the local machine.

// NumCPU 返回本地机器的逻辑CPU数。
func NumCPU() int

// NumCgoCall returns the number of cgo calls made by the current process.

// NumCgoCall 返回由当前进程创建的cgo调用数。
func NumCgoCall() int64

// NumGoroutine returns the number of goroutines that currently exist.

// NumGoroutine 返回当前存在的Go程数。
func NumGoroutine() int

func RaceAcquire(addr unsafe.Pointer)

// RaceDisable disables handling of race events in the current goroutine.

// RaceEnable re-enables handling of race events in the current goroutine.
func RaceDisable()

// RaceEnable re-enables handling of race events in the current goroutine.

// RaceDisable disables handling of race events in the current goroutine.
func RaceEnable()

func RaceRead(addr unsafe.Pointer)

func RaceReadRange(addr unsafe.Pointer, len int)

func RaceRelease(addr unsafe.Pointer)

func RaceReleaseMerge(addr unsafe.Pointer)

func RaceSemacquire(s *uint32)

func RaceSemrelease(s *uint32)

func RaceWrite(addr unsafe.Pointer)

func RaceWriteRange(addr unsafe.Pointer, len int)

// ReadMemStats populates m with memory allocator statistics.

// ReadMemStats 将内存分配器的统计填充到 m 中。
func ReadMemStats(m *MemStats)

// SetBlockProfileRate controls the fraction of goroutine blocking events that are
// reported in the blocking profile. The profiler aims to sample an average of one
// blocking event per rate nanoseconds spent blocked.
//
// To include every blocking event in the profile, pass rate = 1. To turn off
// profiling entirely, pass rate <= 0.
func SetBlockProfileRate(rate int)

// SetCPUProfileRate sets the CPU profiling rate to hz samples per second. If hz <=
// 0, SetCPUProfileRate turns off profiling. If the profiler is on, the rate cannot
// be changed without first turning it off.
//
// Most clients should use the runtime/pprof package or the testing package's
// -test.cpuprofile flag instead of calling SetCPUProfileRate directly.
func SetCPUProfileRate(hz int)

// SetFinalizer sets the finalizer associated with x to f. When the garbage
// collector finds an unreachable block with an associated finalizer, it clears the
// association and runs f(x) in a separate goroutine. This makes x reachable again,
// but now without an associated finalizer. Assuming that SetFinalizer is not
// called again, the next time the garbage collector sees that x is unreachable, it
// will free x.
//
// SetFinalizer(x, nil) clears any finalizer associated with x.
//
// The argument x must be a pointer to an object allocated by calling new or by
// taking the address of a composite literal. The argument f must be a function
// that takes a single argument to which x's type can be assigned, and can have
// arbitrary ignored return values. If either of these is not true, SetFinalizer
// aborts the program.
//
// Finalizers are run in dependency order: if A points at B, both have finalizers,
// and they are otherwise unreachable, only the finalizer for A runs; once A is
// freed, the finalizer for B can run. If a cyclic structure includes a block with
// a finalizer, that cycle is not guaranteed to be garbage collected and the
// finalizer is not guaranteed to run, because there is no ordering that respects
// the dependencies.
//
// The finalizer for x is scheduled to run at some arbitrary time after x becomes
// unreachable. There is no guarantee that finalizers will run before a program
// exits, so typically they are useful only for releasing non-memory resources
// associated with an object during a long-running program. For example, an os.File
// object could use a finalizer to close the associated operating system file
// descriptor when a program discards an os.File without calling Close, but it
// would be a mistake to depend on a finalizer to flush an in-memory I/O buffer
// such as a bufio.Writer, because the buffer would not be flushed at program exit.
//
// It is not guaranteed that a finalizer will run if the size of *x is zero bytes.
//
// It is not guaranteed that a finalizer will run for objects allocated in
// initializers for package-level variables. Such objects may be linker-allocated,
// not heap-allocated.
//
// A single goroutine runs all finalizers for a program, sequentially. If a
// finalizer must run for a long time, it should do so by starting a new goroutine.

// SetFinalizer 为 f 设置与 x 相关联的终结器。
// 当垃圾回收器找到一个无法访问的块及与其相关联的终结器时，就会清理该关联，
// 并在一个独立的Go程中运行f(x)。这会使 x
// 再次变得可访问，但现在没有了相关联的终结器。 假设 SetFinalizer
// 未被再次调用，当下一次垃圾回收器发现 x 无法访问时，就会释放 x。
//
// SetFinalizer(x, nil) 会清理任何与 x 相关联的终结器。
//
// 实参 x
// 必须是一个对象的指针，该对象通过调用新的或获取一个复合字面地址来分配。 实参 f
// 必须是一个函数，该函数获取一个 x
// 的类型的单一实参，并拥有可任意忽略的返回值。
// 只要这些条件有一个不满足，SetFinalizer 就会跳过该程序。
//
// 终结器按照依赖顺序运行：若 A 指向 B，则二者都有终结器，当只有 A 的终结器运行时，
// 它们才无法访问；一旦 A 被释放，则 B
// 的终结器便可运行。若循环依赖的结构包含块及其终结器，
// 则该循环并不能保证被垃圾回收，而其终结器并不能保证运行，这是因为其依赖没有顺序。
//
// x 的终结器预定为在 x
// 无法访问后的任意时刻运行。无法保证终结器会在程序退出前运行，
// 因此它们通常只在长时间运行的程序中释放一个关联至对象的非内存资源时使用。 例如，当程序丢弃 os.File 而没有调用 Close 时，该 os.File
// 对象便可使用一个终结器
// 来关闭与其相关联的操作系统文件描述符，但依赖终结器去刷新一个内存中的I/O缓存是错误的，
// 因为该缓存不会在程序退出时被刷新。
//
// 一个程序的单个Go程会按顺序运行所有的终结器。若某个终结器需要长时间运行，
// 它应当通过开始一个新的Go程来继续。 TODO(osc): 仍需校对及语句优化
func SetFinalizer(obj interface{}, finalizer interface{})

// Stack formats a stack trace of the calling goroutine into buf and returns the
// number of bytes written to buf. If all is true, Stack formats stack traces of
// all other goroutines into buf after the trace for the current goroutine.
func Stack(buf []byte, all bool) int

// ThreadCreateProfile returns n, the number of records in the thread creation
// profile. If len(p) >= n, ThreadCreateProfile copies the profile into p and
// returns n, true. If len(p) < n, ThreadCreateProfile does not change p and
// returns n, false.
//
// Most clients should use the runtime/pprof package instead of calling
// ThreadCreateProfile directly.
func ThreadCreateProfile(p []StackRecord) (n int, ok bool)

// UnlockOSThread unwires the calling goroutine from its fixed operating system
// thread. If the calling goroutine has not called LockOSThread, UnlockOSThread is
// a no-op.
func UnlockOSThread()

// Version returns the Go tree's version string. It is either the commit hash and
// date at the time of the build or, when possible, a release tag like "go1.3".

// Version 返回Go目录树的版本字符串。
// 它一般是一个提交散列值及其构建时间，也可能是一个类似于 "go1.3" 的发行标注。
func Version() string

// BlockProfileRecord describes blocking events originated at a particular call
// sequence (stack trace).
type BlockProfileRecord struct {
	Count  int64
	Cycles int64
	StackRecord
}

type Context C.CONTEXT

type EpollEvent C.struct_epoll_event

// The Error interface identifies a run time error.

// Error 接口用于标识运行时错误。
type Error interface {
	error

	// RuntimeError is a no-op function but
	// serves to distinguish types that are runtime
	// errors from ordinary errors: a type is a
	// runtime error if it has a RuntimeError method.
	RuntimeError()
}

type ExceptionRecord C.EXCEPTION_RECORD

type ExceptionState32 C.struct_i386_exception_state

type ExceptionState64 C.struct_x86_exception_state64

type FPControl C.struct_fp_control

type FPStatus C.struct_fp_status

type FloatState32 C.struct_i386_float_state

type FloatState64 C.struct_x86_float_state64

type FloatingSaveArea C.FLOATING_SAVE_AREA

type Fpreg C.struct__fpreg

type Fpreg1 C.struct__fpreg

type Fpregset C.fpregset_t

type Fpstate C.struct__fpstate

type Fpstate1 C.struct__fpstate

type Fpxreg C.struct__fpxreg

type Fpxreg1 C.struct__fpxreg

// A Func represents a Go function in the running binary.
type Func struct {
	// contains filtered or unexported fields
}

// FuncForPC returns a *Func describing the function that contains the given
// program counter address, or else nil.
func FuncForPC(pc uintptr) *Func

// Entry returns the entry address of the function.
func (f *Func) Entry() uintptr

// FileLine returns the file name and line number of the source code corresponding
// to the program counter pc. The result will not be accurate if pc is not a
// program counter within f.
func (f *Func) FileLine(pc uintptr) (file string, line int)

// Name returns the name of the function.
func (f *Func) Name() string

type Itimerval C.struct_itimerval

type Kevent C.struct_kevent

type KeventT C.struct_kevent

type Lwpparams C.struct_lwp_params

type M128a C.M128A

type MachBody C.mach_msg_body_t

type MachHeader C.mach_msg_header_t

type MachNDR C.NDR_record_t

type MachPort C.mach_msg_port_descriptor_t

type Mcontext C.mcontext_t

type Mcontext32 C.struct_mcontext32

type Mcontext64 C.struct_mcontext64

type McontextT C.mcontext_t

// A MemProfileRecord describes the live objects allocated by a particular call
// sequence (stack trace).
type MemProfileRecord struct {
	AllocBytes, FreeBytes     int64       // number of bytes allocated, freed
	AllocObjects, FreeObjects int64       // number of objects allocated, freed
	Stack0                    [32]uintptr // stack trace for this record; ends at first 0 entry
}

// InUseBytes returns the number of bytes in use (AllocBytes - FreeBytes).
func (r *MemProfileRecord) InUseBytes() int64

// InUseObjects returns the number of objects in use (AllocObjects - FreeObjects).
func (r *MemProfileRecord) InUseObjects() int64

// Stack returns the stack trace associated with the record, a prefix of r.Stack0.
func (r *MemProfileRecord) Stack() []uintptr

// A MemStats records statistics about the memory allocator.

// MemStats 用于记录内存分配器的统计量。
type MemStats struct {
	// General statistics.
	Alloc      uint64 // bytes allocated and still in use
	TotalAlloc uint64 // bytes allocated (even if freed)
	Sys        uint64 // bytes obtained from system (sum of XxxSys below)
	Lookups    uint64 // number of pointer lookups
	Mallocs    uint64 // number of mallocs
	Frees      uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    uint64 // bytes allocated and still in use
	HeapSys      uint64 // bytes obtained from system
	HeapIdle     uint64 // bytes in idle spans
	HeapInuse    uint64 // bytes in non-idle span
	HeapReleased uint64 // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  uint64 // bytes used by stack allocator
	StackSys    uint64
	MSpanInuse  uint64 // mspan structures
	MSpanSys    uint64
	MCacheInuse uint64 // mcache structures
	MCacheSys   uint64
	BuckHashSys uint64 // profiling bucket hash table
	GCSys       uint64 // GC metadata
	OtherSys    uint64 // other system allocations

	// Garbage collector statistics.
	NextGC       uint64 // next collection will happen when HeapAlloc ≥ this amount
	LastGC       uint64 // end time of last collection (nanoseconds since 1970)
	PauseTotalNs uint64
	PauseNs      [256]uint64 // circular buffer of recent GC pause durations, most recent at [(NumGC+255)%256]
	PauseEnd     [256]uint64 // circular buffer of recent GC pause end times
	NumGC        uint32
	EnableGC     bool
	DebugGC      bool

	// Per-size allocation statistics.
	// 61 is NumSizeClasses in the C code.
	BySize [61]struct {
		Size    uint32
		Mallocs uint64
		Frees   uint64
	}
}

type Overlapped C.OVERLAPPED

type PortEvent C.port_event_t

type Pthread C.pthread_t

type PthreadAttr C.pthread_attr_t

type RegMMST C.struct_mmst_reg

type RegXMM C.struct_xmm_reg

type Regs32 C.struct_i386_thread_state

type Regs64 C.struct_x86_thread_state64

type Rtprio C.struct_rtprio

type SemT C.sem_t

type Sigaction C.struct_sigaction

type SigaltstackT C.struct_sigaltstack

type Sigcontext C.struct_sigcontext

type Sighandler C.union___sigaction_u

type Siginfo C.siginfo_t

type Sigset C.sigset_t

type Sigval C.union_sigval

// A StackRecord describes a single execution stack.
type StackRecord struct {
	Stack0 [32]uintptr // stack trace for this record; ends at first 0 entry
}

// Stack returns the stack trace associated with the record, a prefix of r.Stack0.
func (r *StackRecord) Stack() []uintptr

type StackT C.stack_t

// depends on Timespec, must appear below
type Stat C.struct_stat

type SystemInfo C.SYSTEM_INFO

type TforkT C.struct___tfork

type ThrParam C.struct_thr_param

type Timespec C.struct_timespec

type Timeval C.struct_timeval

// A TypeAssertionError explains a failed type assertion.

// TypeAssertionError 用于阐明失败的类型断言。
type TypeAssertionError struct {
	// contains filtered or unexported fields
}

func (e *TypeAssertionError) Error() string

func (*TypeAssertionError) RuntimeError()

type Ucontext C.ucontext_t

type UcontextT C.ucontext_t

type Usigset C.__sigset_t

type Xmmreg C.struct__xmmreg

type Xmmreg1 C.struct__xmmreg
