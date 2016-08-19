// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package runtime contains operations that interact with Go's runtime system,
// such as functions to control goroutines. It also includes the low-level type
// information used by the reflect package; see reflect's documentation for the
// programmable interface to the run-time type system.
//
//
// Environment Variables
//
// The following environment variables ($name or %name%, depending on the host
// operating system) control the run-time behavior of Go programs. The meanings
// and use may change from release to release.
//
// The GOGC variable sets the initial garbage collection target percentage. A
// collection is triggered when the ratio of freshly allocated data to live data
// remaining after the previous collection reaches this percentage. The default
// is GOGC=100. Setting GOGC=off disables the garbage collector entirely. The
// runtime/debug package's SetGCPercent function allows changing this percentage
// at run time. See https://golang.org/pkg/runtime/debug/#SetGCPercent.
//
// The GODEBUG variable controls debugging variables within the runtime. It is a
// comma-separated list of name=val pairs setting these named variables:
//
//     allocfreetrace: setting allocfreetrace=1 causes every allocation to be
//     profiled and a stack trace printed on each object's allocation and free.
//
//     cgocheck: setting cgocheck=0 disables all checks for packages
//     using cgo to incorrectly pass Go pointers to non-Go code.
//     Setting cgocheck=1 (the default) enables relatively cheap
//     checks that may miss some errors.  Setting cgocheck=2 enables
//     expensive checks that should not miss any errors, but will
//     cause your program to run slower.
//
//     efence: setting efence=1 causes the allocator to run in a mode
//     where each object is allocated on a unique page and addresses are
//     never recycled.
//
//     gccheckmark: setting gccheckmark=1 enables verification of the
//     garbage collector's concurrent mark phase by performing a
//     second mark pass while the world is stopped.  If the second
//     pass finds a reachable object that was not found by concurrent
//     mark, the garbage collector will panic.
//
//     gcpacertrace: setting gcpacertrace=1 causes the garbage collector to
//     print information about the internal state of the concurrent pacer.
//
//     gcshrinkstackoff: setting gcshrinkstackoff=1 disables moving goroutines
//     onto smaller stacks. In this mode, a goroutine's stack can only grow.
//
//     gcstackbarrieroff: setting gcstackbarrieroff=1 disables the use of stack barriers
//     that allow the garbage collector to avoid repeating a stack scan during the
//     mark termination phase.
//
//     gcstackbarrierall: setting gcstackbarrierall=1 installs stack barriers
//     in every stack frame, rather than in exponentially-spaced frames.
//
//     gcstoptheworld: setting gcstoptheworld=1 disables concurrent garbage collection,
//     making every garbage collection a stop-the-world event. Setting gcstoptheworld=2
//     also disables concurrent sweeping after the garbage collection finishes.
//
//     gctrace: setting gctrace=1 causes the garbage collector to emit a single line to standard
//     error at each collection, summarizing the amount of memory collected and the
//     length of the pause. Setting gctrace=2 emits the same summary but also
//     repeats each collection. The format of this line is subject to change.
//     Currently, it is:
//     	gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
//     where the fields are as follows:
//     	gc #        the GC number, incremented at each GC
//     	@#s         time in seconds since program start
//     	#%          percentage of time spent in GC since program start
//     	#+...+#     wall-clock/CPU times for the phases of the GC
//     	#->#-># MB  heap size at GC start, at GC end, and live heap
//     	# MB goal   goal heap size
//     	# P         number of processors used
//     The phases are stop-the-world (STW) sweep termination, concurrent
//     mark and scan, and STW mark termination. The CPU times
//     for mark/scan are broken down in to assist time (GC performed in
//     line with allocation), background GC time, and idle GC time.
//     If the line ends with "(forced)", this GC was forced by a
//     runtime.GC() call and all phases are STW.
//
//     memprofilerate: setting memprofilerate=X will update the value of runtime.MemProfileRate.
//     When set to 0 memory profiling is disabled.  Refer to the description of
//     MemProfileRate for the default value.
//
//     invalidptr: defaults to invalidptr=1, causing the garbage collector and stack
//     copier to crash the program if an invalid pointer value (for example, 1)
//     is found in a pointer-typed location. Setting invalidptr=0 disables this check.
//     This should only be used as a temporary workaround to diagnose buggy code.
//     The real fix is to not store integers in pointer-typed locations.
//
//     sbrk: setting sbrk=1 replaces the memory allocator and garbage collector
//     with a trivial allocator that obtains memory from the operating system and
//     never reclaims any memory.
//
//     scavenge: scavenge=1 enables debugging mode of heap scavenger.
//
//     scheddetail: setting schedtrace=X and scheddetail=1 causes the scheduler to emit
//     detailed multiline info every X milliseconds, describing state of the scheduler,
//     processors, threads and goroutines.
//
//     schedtrace: setting schedtrace=X causes the scheduler to emit a single line to standard
//     error every X milliseconds, summarizing the scheduler state.
//
// The net and net/http packages also refer to debugging variables in GODEBUG.
// See the documentation for those packages for details.
//
// The GOMAXPROCS variable limits the number of operating system threads that
// can execute user-level Go code simultaneously. There is no limit to the
// number of threads that can be blocked in system calls on behalf of Go code;
// those do not count against the GOMAXPROCS limit. This package's GOMAXPROCS
// function queries and changes the limit.
//
// The GOTRACEBACK variable controls the amount of output generated when a Go
// program fails due to an unrecovered panic or an unexpected runtime condition.
// By default, a failure prints a stack trace for the current goroutine, eliding
// functions internal to the run-time system, and then exits with exit code 2.
// The failure prints stack traces for all goroutines if there is no current
// goroutine or the failure is internal to the run-time. GOTRACEBACK=none omits
// the goroutine stack traces entirely. GOTRACEBACK=single (the default) behaves
// as described above. GOTRACEBACK=all adds stack traces for all user-created
// goroutines. GOTRACEBACK=system is like ``all'' but adds stack frames for
// run-time functions and shows goroutines created internally by the run-time.
// GOTRACEBACK=crash is like ``system'' but crashes in an operating
// system-specific manner instead of exiting. For example, on Unix systems, the
// crash raises SIGABRT to trigger a core dump. For historical reasons, the
// GOTRACEBACK settings 0, 1, and 2 are synonyms for none, all, and system,
// respectively. The runtime/debug package's SetTraceback function allows
// increasing the amount of output at run time, but it cannot reduce the amount
// below that specified by the environment variable. See
// https://golang.org/pkg/runtime/debug/#SetTraceback.
//
// The GOARCH, GOOS, GOPATH, and GOROOT environment variables complete the set
// of Go environment variables. They influence the building of Go programs (see
// https://golang.org/cmd/go and https://golang.org/pkg/go/build). GOARCH, GOOS,
// and GOROOT are recorded at compile time and made available by constants or
// functions in this package, but they do not influence the execution of the
// run-time system.

// TODO(osc): 需更新 runtime 包含与Go的运行时系统进行交互的操作，例如用于控制Go
// 程的函数. 它也包括用于 reflect 包的底层类型信息；运行时类型系统的可编程接口见
// reflect 文档。
//
// 环境变量
//
// 以下环境变量（$name 或 %name%, 取决于宿主操作系统）控制了Go程序的运行时行为。
// 其意义与使用方法在发行版之间可能有所不同。
//
// GOGC 变量用于设置初始垃圾回收的目标百分比。从上次回收后开始，当新分配数据的比
// 例占到剩余实时数据的此百分比时， 就会再次触发回收。默认为 GOGC=100。要完全关
// 闭垃圾回收器，需设置 GOGC=off。runtime/debug 包的 SetGCPercent 函数允许在运行
// 时更改此百分比。 详见
// http://zh.golanger.com/pkg/runtime/debug/#SetGCPercent。
//
// GOGCTRACE 变量用于控制来自垃圾回收器的调试输出。设置 GOGCTRACE=1 会使垃圾回收
// 器发出 每一次回收所产生的单行标准错误输出、概述回收的内存量以及暂停的时长。设
// 置 GOGCTRACE=2 不仅会发出同样的概述，还会重复每一次回收。
//
// GOMAXPROCS 变量用于限制可同时执行的用户级Go代码所产生的操作系统线程数。对于Go
// 代码所代表的系统调用而言， 可被阻塞的线程则没有限制；它们不计入 GOMAXPROCS 的
// 限制。本包中的 GOMAXPROCS 函数可查询并更改此限制。
//
// GOTRACEBACK 用于控制因未恢复的恐慌或意外的运行时状况导致Go程序运行失败时所产
// 生的输出量。 默认情况下，失败会为每个现有的Go程打印出栈跟踪，省略运行时系统的
// 内部函数，并以退出码 2 退出。 若 GOTRACEBACK=0，则每个Go程的栈跟踪都会完全省
// 略。 若 GOTRACEBACK=1，则采用默认的行为。 若 GOTRACEBACK=2，则每个Go程的栈跟
// 踪，包括运行时函数都会输出。 若 GOTRACEBACK=crash，则每个Go程的栈跟踪，包括运
// 行时函数，都会输出， 此外程序可能以操作系统特定的方式崩溃而非退出。例如，在
// Unix系统上，程序会发出 SIGABRT 信号，从而触发内核转储。
//
// GOARCH、GOOS、GOPATH 和 GOROOT 环境变量均为Go的环境变量。它们影响了Go程序的构
// 建 （详见 http://golang.org/cmd/go 和 http://golang.org/pkg/go/build）。
// GOARCH、GOOS 和 GOROOT 会在编译时被记录，并使该包中的常量或函数变得可用， 但
// 它们并不影响运行时系统的执行。
//
// 公共的竞争检测API，当且仅当使用 -race 构建时才会出现。
package runtime

import (
    "C"
    "runtime/internal/atomic"
    "runtime/internal/sys"
    "unsafe"
)

// Compiler is the name of the compiler toolchain that built the
// running binary.  Known toolchains are:
//
//     gc      Also known as cmd/compile.
//     gccgo   The gccgo front end, part of the GCC compiler suite.

// Compiler 为构建了可运行二进制文件的编译工具链。已知的工具链为：
//     go       code.google.com/p/go 上的 5g/6g/8g 编译器套件。
//     gccgo    gccgo前端，GCC编译器条件的一部分。
const Compiler = "gc"



const (
	EINTR          = C.EINTR
	EAGAIN         = C.EAGAIN
	ENOMEM         = C.ENOMEM
	PROT_NONE      = C.PROT_NONE
	PROT_READ      = C.PROT_READ
	PROT_WRITE     = C.PROT_WRITE
	PROT_EXEC      = C.PROT_EXEC
	MAP_ANON       = C.MAP_ANONYMOUS
	MAP_PRIVATE    = C.MAP_PRIVATE
	MAP_FIXED      = C.MAP_FIXED
	MADV_DONTNEED  = C.MADV_DONTNEED
	SA_RESTART     = C.SA_RESTART
	SA_ONSTACK     = C.SA_ONSTACK
	SA_RESTORER    = C.SA_RESTORER
	SA_SIGINFO     = C.SA_SIGINFO
	SIGHUP         = C.SIGHUP
	SIGINT         = C.SIGINT
	SIGQUIT        = C.SIGQUIT
	SIGILL         = C.SIGILL
	SIGTRAP        = C.SIGTRAP
	SIGABRT        = C.SIGABRT
	SIGBUS         = C.SIGBUS
	SIGFPE         = C.SIGFPE
	SIGKILL        = C.SIGKILL
	SIGUSR1        = C.SIGUSR1
	SIGSEGV        = C.SIGSEGV
	SIGUSR2        = C.SIGUSR2
	SIGPIPE        = C.SIGPIPE
	SIGALRM        = C.SIGALRM
	SIGSTKFLT      = C.SIGSTKFLT
	SIGCHLD        = C.SIGCHLD
	SIGCONT        = C.SIGCONT
	SIGSTOP        = C.SIGSTOP
	SIGTSTP        = C.SIGTSTP
	SIGTTIN        = C.SIGTTIN
	SIGTTOU        = C.SIGTTOU
	SIGURG         = C.SIGURG
	SIGXCPU        = C.SIGXCPU
	SIGXFSZ        = C.SIGXFSZ
	SIGVTALRM      = C.SIGVTALRM
	SIGPROF        = C.SIGPROF
	SIGWINCH       = C.SIGWINCH
	SIGIO          = C.SIGIO
	SIGPWR         = C.SIGPWR
	SIGSYS         = C.SIGSYS
	FPE_INTDIV     = C.FPE_INTDIV
	FPE_INTOVF     = C.FPE_INTOVF
	FPE_FLTDIV     = C.FPE_FLTDIV
	FPE_FLTOVF     = C.FPE_FLTOVF
	FPE_FLTUND     = C.FPE_FLTUND
	FPE_FLTRES     = C.FPE_FLTRES
	FPE_FLTINV     = C.FPE_FLTINV
	FPE_FLTSUB     = C.FPE_FLTSUB
	BUS_ADRALN     = C.BUS_ADRALN
	BUS_ADRERR     = C.BUS_ADRERR
	BUS_OBJERR     = C.BUS_OBJERR
	SEGV_MAPERR    = C.SEGV_MAPERR
	SEGV_ACCERR    = C.SEGV_ACCERR
	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF
	O_RDONLY       = C.O_RDONLY
	O_CLOEXEC      = C.O_CLOEXEC
	EPOLLIN        = C.POLLIN
	EPOLLOUT       = C.POLLOUT
	EPOLLERR       = C.POLLERR
	EPOLLHUP       = C.POLLHUP
	EPOLLRDHUP     = C.POLLRDHUP
	EPOLLET        = C.EPOLLET
	EPOLL_CLOEXEC  = C.EPOLL_CLOEXEC
	EPOLL_CTL_ADD  = C.EPOLL_CTL_ADD
	EPOLL_CTL_DEL  = C.EPOLL_CTL_DEL
	EPOLL_CTL_MOD  = C.EPOLL_CTL_MOD
)



const (
	EINTR          = C.EINTR
	EAGAIN         = C.EAGAIN
	ENOMEM         = C.ENOMEM
	PROT_NONE      = C.PROT_NONE
	PROT_READ      = C.PROT_READ
	PROT_WRITE     = C.PROT_WRITE
	PROT_EXEC      = C.PROT_EXEC
	MAP_ANON       = C.MAP_ANONYMOUS
	MAP_PRIVATE    = C.MAP_PRIVATE
	MAP_FIXED      = C.MAP_FIXED
	MADV_DONTNEED  = C.MADV_DONTNEED
	SA_RESTART     = C.SA_RESTART
	SA_ONSTACK     = C.SA_ONSTACK
	SA_SIGINFO     = C.SA_SIGINFO
	SIGHUP         = C.SIGHUP
	SIGINT         = C.SIGINT
	SIGQUIT        = C.SIGQUIT
	SIGILL         = C.SIGILL
	SIGTRAP        = C.SIGTRAP
	SIGABRT        = C.SIGABRT
	SIGBUS         = C.SIGBUS
	SIGFPE         = C.SIGFPE
	SIGKILL        = C.SIGKILL
	SIGUSR1        = C.SIGUSR1
	SIGSEGV        = C.SIGSEGV
	SIGUSR2        = C.SIGUSR2
	SIGPIPE        = C.SIGPIPE
	SIGALRM        = C.SIGALRM
	SIGSTKFLT      = C.SIGSTKFLT
	SIGCHLD        = C.SIGCHLD
	SIGCONT        = C.SIGCONT
	SIGSTOP        = C.SIGSTOP
	SIGTSTP        = C.SIGTSTP
	SIGTTIN        = C.SIGTTIN
	SIGTTOU        = C.SIGTTOU
	SIGURG         = C.SIGURG
	SIGXCPU        = C.SIGXCPU
	SIGXFSZ        = C.SIGXFSZ
	SIGVTALRM      = C.SIGVTALRM
	SIGPROF        = C.SIGPROF
	SIGWINCH       = C.SIGWINCH
	SIGIO          = C.SIGIO
	SIGPWR         = C.SIGPWR
	SIGSYS         = C.SIGSYS
	FPE_INTDIV     = C.FPE_INTDIV
	FPE_INTOVF     = C.FPE_INTOVF
	FPE_FLTDIV     = C.FPE_FLTDIV
	FPE_FLTOVF     = C.FPE_FLTOVF
	FPE_FLTUND     = C.FPE_FLTUND
	FPE_FLTRES     = C.FPE_FLTRES
	FPE_FLTINV     = C.FPE_FLTINV
	FPE_FLTSUB     = C.FPE_FLTSUB
	BUS_ADRALN     = C.BUS_ADRALN
	BUS_ADRERR     = C.BUS_ADRERR
	BUS_OBJERR     = C.BUS_OBJERR
	SEGV_MAPERR    = C.SEGV_MAPERR
	SEGV_ACCERR    = C.SEGV_ACCERR
	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF
	EPOLLIN        = C.POLLIN
	EPOLLOUT       = C.POLLOUT
	EPOLLERR       = C.POLLERR
	EPOLLHUP       = C.POLLHUP
	EPOLLRDHUP     = C.POLLRDHUP
	EPOLLET        = C.EPOLLET
	EPOLL_CLOEXEC  = C.EPOLL_CLOEXEC
	EPOLL_CTL_ADD  = C.EPOLL_CTL_ADD
	EPOLL_CTL_DEL  = C.EPOLL_CTL_DEL
	EPOLL_CTL_MOD  = C.EPOLL_CTL_MOD
)


// GOARCH is the running program's architecture target:
// 386, amd64, or arm.

// GOARCH 为所运行程序的目标架构：
// 386、amd64、arm 或 s390x。
const GOARCH string = sys.GOARCH


// GOOS is the running program's operating system target:
// one of darwin, freebsd, linux, and so on.

// GOOS 为所运行程序的目标操作系统：
// darwin、freebsd或linux等等。
const GOOS string = sys.GOOS



const (
	O_RDONLY    = C.O_RDONLY
	O_CLOEXEC   = C.O_CLOEXEC
	SA_RESTORER = C.SA_RESTORER
)



const (
	O_RDONLY    = C.O_RDONLY
	O_CLOEXEC   = C.O_CLOEXEC
	SA_RESTORER = 0 // unused

)



const (
	PROT_NONE      = C.PROT_NONE
	PROT_READ      = C.PROT_READ
	PROT_WRITE     = C.PROT_WRITE
	PROT_EXEC      = C.PROT_EXEC
	MAP_ANON       = C.MAP_ANONYMOUS
	MAP_PRIVATE    = C.MAP_PRIVATE
	MAP_FIXED      = C.MAP_FIXED
	MADV_DONTNEED  = C.MADV_DONTNEED
	SA_RESTART     = C.SA_RESTART
	SA_ONSTACK     = C.SA_ONSTACK
	SA_RESTORER    = C.SA_RESTORER
	SA_SIGINFO     = C.SA_SIGINFO
	SIGHUP         = C.SIGHUP
	SIGINT         = C.SIGINT
	SIGQUIT        = C.SIGQUIT
	SIGILL         = C.SIGILL
	SIGTRAP        = C.SIGTRAP
	SIGABRT        = C.SIGABRT
	SIGBUS         = C.SIGBUS
	SIGFPE         = C.SIGFPE
	SIGKILL        = C.SIGKILL
	SIGUSR1        = C.SIGUSR1
	SIGSEGV        = C.SIGSEGV
	SIGUSR2        = C.SIGUSR2
	SIGPIPE        = C.SIGPIPE
	SIGALRM        = C.SIGALRM
	SIGSTKFLT      = C.SIGSTKFLT
	SIGCHLD        = C.SIGCHLD
	SIGCONT        = C.SIGCONT
	SIGSTOP        = C.SIGSTOP
	SIGTSTP        = C.SIGTSTP
	SIGTTIN        = C.SIGTTIN
	SIGTTOU        = C.SIGTTOU
	SIGURG         = C.SIGURG
	SIGXCPU        = C.SIGXCPU
	SIGXFSZ        = C.SIGXFSZ
	SIGVTALRM      = C.SIGVTALRM
	SIGPROF        = C.SIGPROF
	SIGWINCH       = C.SIGWINCH
	SIGIO          = C.SIGIO
	SIGPWR         = C.SIGPWR
	SIGSYS         = C.SIGSYS
	FPE_INTDIV     = C.FPE_INTDIV & 0xFFFF
	FPE_INTOVF     = C.FPE_INTOVF & 0xFFFF
	FPE_FLTDIV     = C.FPE_FLTDIV & 0xFFFF
	FPE_FLTOVF     = C.FPE_FLTOVF & 0xFFFF
	FPE_FLTUND     = C.FPE_FLTUND & 0xFFFF
	FPE_FLTRES     = C.FPE_FLTRES & 0xFFFF
	FPE_FLTINV     = C.FPE_FLTINV & 0xFFFF
	FPE_FLTSUB     = C.FPE_FLTSUB & 0xFFFF
	BUS_ADRALN     = C.BUS_ADRALN & 0xFFFF
	BUS_ADRERR     = C.BUS_ADRERR & 0xFFFF
	BUS_OBJERR     = C.BUS_OBJERR & 0xFFFF
	SEGV_MAPERR    = C.SEGV_MAPERR & 0xFFFF
	SEGV_ACCERR    = C.SEGV_ACCERR & 0xFFFF
	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_PROF    = C.ITIMER_PROF
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
)



const (
	_ selectDir = iota
)


// MemProfileRate controls the fraction of memory allocations
// that are recorded and reported in the memory profile.
// The profiler aims to sample an average of
// one allocation per MemProfileRate bytes allocated.
//
// To include every allocated block in the profile, set MemProfileRate to 1.
// To turn off profiling entirely, set MemProfileRate to 0.
//
// The tools that process the memory profiles assume that the
// profile rate is constant across the lifetime of the program
// and equal to the current value.  Programs that change the
// memory profiling rate should do so just once, as early as
// possible in the execution of the program (for example,
// at the beginning of main).

// MemProfileRate controls the fraction of memory allocations
// that are recorded and reported in the memory profile.
// The profiler aims to sample an average of
// one allocation per MemProfileRate bytes allocated.
//
// To include every allocated block in the profile, set MemProfileRate to 1.
// To turn off profiling entirely, set MemProfileRate to 0.
//
// The tools that process the memory profiles assume that the
// profile rate is constant across the lifetime of the program
// and equal to the current value. Programs that change the
// memory profiling rate should do so just once, as early as
// possible in the execution of the program (for example,
// at the beginning of main).
var MemProfileRate int = 512 * 1024


// BlockProfileRecord describes blocking events originated
// at a particular call sequence (stack trace).
type BlockProfileRecord struct {
	Count  int64
	Cycles int64
}



type EpollEvent C.struct_epoll_event



type EpollEvent C.struct_epoll_event


// The Error interface identifies a run time error.

// Error 接口用于标识运行时错误。
type Error interface {
	error

	// RuntimeError is a no-op function but
	// serves to distinguish types that are run time
	// errors from ordinary errors: a type is a
	// run time error if it has a RuntimeError method.
	//
	// RuntimeError 是一个无操作函数，它只用于区分是运行时错误还是一般错误：
	// 若一个类型拥有 RuntimeError 方法，它就是运行时错误。
	RuntimeError()
}



type FPregset C.elf_fpregset_t



type Fpreg C.struct__fpreg



type Fpreg1 C.struct__fpreg



type Fpstate C.struct__fpstate



type Fpstate C.struct__libc_fpstate



type Fpstate1 C.struct__fpstate



type Fpxreg C.struct__fpxreg



type Fpxreg C.struct__libc_fpxreg



type Fpxreg1 C.struct__fpxreg


// Frame is the information returned by Frames for each call frame.
type Frame struct {
	// Program counter for this frame; multiple frames may have
	// the same PC value.
	PC uintptr

	// Func for this frame; may be nil for non-Go code or fully
	// inlined functions.
	Func *Func

	// Function name, file name, and line number for this call frame.
	// May be the empty string or zero if not known.
	// If Func is not nil then Function == Func.Name().
	Function string
	File     string
	Line     int

	// Entry point for the function; may be zero if not known.
	// If Func is not nil then Entry == Func.Entry().
	Entry uintptr
}


// Frames may be used to get function/file/line information for a
// slice of PC values returned by Callers.
type Frames struct {
	callers []uintptr

	// If previous caller in iteration was a panic, then
	// ci.callers[0] is the address of the faulting instruction
	// instead of the return address of the call.
	wasPanic bool

	// Frames to return for subsequent calls to the Next method.
	// Used for non-Go frames.
	frames *[]Frame
}


// A Func represents a Go function in the running binary.
type Func struct {
	opaque struct{} // unexported field to disallow conversions
}



type Gregset C.elf_gregset_t



type Itimerval C.struct_itimerval



type Itimerval C.struct_itimerval



type Itimerval C.struct_itimerval



type Mcontext C.mcontext_t


// A MemProfileRecord describes the live objects allocated
// by a particular call sequence (stack trace).
type MemProfileRecord struct {
	AllocBytes, FreeBytes     int64       // number of bytes allocated, freed
	AllocObjects, FreeObjects int64       // number of objects allocated, freed
	Stack0                    [32]uintptr // stack trace for this record; ends at first 0 entry
}


// A MemStats records statistics about the memory allocator.
type MemStats struct {
	// General statistics.
	Alloc      uint64 // bytes allocated and not yet freed
	TotalAlloc uint64 // bytes allocated (even if freed)
	Sys        uint64 // bytes obtained from system (sum of XxxSys below)
	Lookups    uint64 // number of pointer lookups
	Mallocs    uint64 // number of mallocs
	Frees      uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    uint64 // bytes allocated and not yet freed (same as Alloc above)
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
	NextGC        uint64 // next collection will happen when HeapAlloc ≥ this amount
	LastGC        uint64 // end time of last collection (nanoseconds since 1970)
	PauseTotalNs  uint64
	PauseNs       [256]uint64 // circular buffer of recent GC pause durations, most recent at [(NumGC+255)%256]
	PauseEnd      [256]uint64 // circular buffer of recent GC pause end times
	NumGC         uint32
	GCCPUFraction float64 // fraction of CPU time used by GC
	EnableGC      bool
	DebugGC       bool

	// Per-size allocation statistics.
	// 61 is NumSizeClasses in the C code.
	BySize [61]struct {
		Size    uint32
		Mallocs uint64
		Frees   uint64
	}
}


// types used in sigcontext
type Ptregs C.struct_pt_regs



type Sigaction C.struct_sigaction



type Sigaction C.struct_xsigaction



type Sigaction C.struct_kernel_sigaction



type SigaltstackT C.struct_sigaltstack



type SigaltstackT C.struct_sigaltstack



type SigaltstackT C.struct_sigaltstack



type SigaltstackT C.struct_sigaltstack


// PPC64 uses sigcontext in place of mcontext in ucontext. see
// http://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/tree/arch/powerpc/include/uapi/asm/ucontext.h

// PPC64 uses sigcontext in place of mcontext in ucontext. see
// http://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/tree/arch/powerpc/include/uapi/asm/ucontext.h
type Sigcontext C.struct_sigcontext



type Sigcontext C.struct_sigcontext



type Sigcontext C.struct_sigcontext



type Sigcontext C.struct_sigcontext



type Siginfo C.siginfo_t



type Siginfo C.siginfo_t



type Siginfo C.struct_xsiginfo



type Sigset C.sigset_t


// A StackRecord describes a single execution stack.
type StackRecord struct {
	Stack0 [32]uintptr // stack trace for this record; ends at first 0 entry
}



type Timespec C.struct_timespec



type Timespec C.struct_timespec



type Timespec C.struct_timespec



type Timeval C.struct_timeval



type Timeval C.struct_timeval



type Timeval C.struct_timeval


// A TypeAssertionError explains a failed type assertion.

// TypeAssertionError 用于阐明失败的类型断言。
type TypeAssertionError struct {
	interfaceString string
	concreteString  string
	assertedString  string
	missingMethod   string // one method needed by Interface, missing from Concrete

}



type Ucontext C.struct_ucontext



type Ucontext C.struct_ucontext



type Ucontext C.ucontext_t



type Ucontext C.struct_ucontext



type Usigset C.__sigset_t



type Usigset C.__sigset_t



type Vreg C.elf_vrreg_t



type Xmmreg C.struct__libc_xmmreg



type Xmmreg C.struct__xmmreg



type Xmmreg1 C.struct__xmmreg


// BlockProfile returns n, the number of records in the current blocking
// profile. If len(p) >= n, BlockProfile copies the profile into p and returns
// n, true. If len(p) < n, BlockProfile does not change p and returns n, false.
//
// Most clients should use the runtime/pprof package or the testing package's
// -test.blockprofile flag instead of calling BlockProfile directly.
func BlockProfile(p []BlockProfileRecord) (n int, ok bool)

// Breakpoint executes a breakpoint trap.
func Breakpoint()

// CPUProfile returns the next chunk of binary CPU profiling stack trace data,
// blocking until data is available. If profiling is turned off and all the
// profile data accumulated while it was on has been returned, CPUProfile
// returns nil. The caller must save the returned data before calling CPUProfile
// again.
//
// Most clients should use the runtime/pprof package or the testing package's
// -test.cpuprofile flag instead of calling CPUProfile directly.

// CPUProfile returns the next chunk of binary CPU profiling stack trace data,
// blocking until data is available. If profiling is turned off and all the
// profile data accumulated while it was on has been returned, CPUProfile
// returns nil. The caller must save the returned data before calling CPUProfile
// again.
//
// Most clients should use the runtime/pprof package or the testing package's
// -test.cpuprofile flag instead of calling CPUProfile directly.
func CPUProfile() []byte

// Caller reports file and line number information about function invocations on
// the calling goroutine's stack. The argument skip is the number of stack
// frames to ascend, with 0 identifying the caller of Caller. (For historical
// reasons the meaning of skip differs between Caller and Callers.) The return
// values report the program counter, file name, and line number within the file
// of the corresponding call. The boolean ok is false if it was not possible to
// recover the information.

// Caller 报告关于调用Go程的栈上的函数调用的文件和行号信息。
// 实参 skip 为占用的栈帧数，若为0则表示 Caller 的调用者。（由于历史原因，skip
// 的意思在 Caller 和 Callers 中并不相同。）返回值报告程序计数器，
// 文件名及对应调用的文件中的行号。若无法获得信息，布尔值 ok 即为 false。
func Caller(skip int) (pc uintptr, file string, line int, ok bool)

// Callers fills the slice pc with the return program counters of function
// invocations on the calling goroutine's stack. The argument skip is the number
// of stack frames to skip before recording in pc, with 0 identifying the frame
// for Callers itself and 1 identifying the caller of Callers. It returns the
// number of entries written to pc.
//
// Note that since each slice entry pc[i] is a return program counter, looking
// up the file and line for pc[i] (for example, using (*Func).FileLine) will
// return the file and line number of the instruction immediately following the
// call. To look up the file and line number of the call itself, use pc[i]-1. As
// an exception to this rule, if pc[i-1] corresponds to the function
// runtime.sigpanic, then pc[i] is the program counter of a faulting instruction
// and should be used without any subtraction.

// Callers 把调用它的Go程栈上函数请求的返回程序计数器填充到切片 pc 中。 实参
// skip 为开始在 pc 中记录之前所要跳过的栈帧数，若为 0 则表示 Callers 自身的栈
// 帧， 若为 1 则表示 Callers 的调用者。它返回写入到 pc 中的项数。
//
// 注意，由于每个切片项 pc[i] 都是一个返回程序计数器，因此查找 pc[i] 的文件和
// 行（例如，使用 (*Func).FileLine）将会在该调用之后立即返回该指令所在的文件和行
// 号。 要想在调用序列中方便地查看文件/行号信息，请使用 Frames。
func Callers(skip int, pc []uintptr) int

// CallersFrames takes a slice of PC values returned by Callers and
// prepares to return function/file/line information.
// Do not change the slice until you are done with the Frames.
func CallersFrames(callers []uintptr) *Frames

// FuncForPC returns a *Func describing the function that contains the
// given program counter address, or else nil.
func FuncForPC(pc uintptr) *Func

// GC runs a garbage collection and blocks the caller until the
// garbage collection is complete. It may also block the entire
// program.
func GC()

// GOMAXPROCS sets the maximum number of CPUs that can be executing
// simultaneously and returns the previous setting.  If n < 1, it does not
// change the current setting.
// The number of logical CPUs on the local machine can be queried with NumCPU.
// This call will go away when the scheduler improves.

// GOMAXPROCS 设置可同时使用执行的最大CPU数，并返回先前的设置。
// 若 n < 1，它就不会更改当前设置。本地机器的逻辑CPU数可通过 NumCPU 查询。
// 当调度器改进后，此调用将会消失。
func GOMAXPROCS(n int) int

// GOROOT returns the root of the Go tree.
// It uses the GOROOT environment variable, if set,
// or else the root used during the Go build.

// GOROOT 返回Go目录树的根目录。
// 若设置了GOROOT环境变量，就会使用它，否则就会将Go的构建目录作为根目录
func GOROOT() string

// Goexit terminates the goroutine that calls it. No other goroutine is
// affected. Goexit runs all deferred calls before terminating the goroutine.
// Because Goexit is not panic, however, any recover calls in those deferred
// functions will return nil.
//
// Calling Goexit from the main goroutine terminates that goroutine without func
// main returning. Since func main has not returned, the program continues
// execution of other goroutines. If all other goroutines exit, the program
// crashes.

// Goexit terminates the goroutine that calls it. No other goroutine is
// affected. Goexit runs all deferred calls before terminating the goroutine.
// Because Goexit is not panic, however, any recover calls in those deferred
// functions will return nil.
//
// Calling Goexit from the main goroutine terminates that goroutine without func
// main returning. Since func main has not returned, the program continues
// execution of other goroutines. If all other goroutines exit, the program
// crashes.
func Goexit()

// GoroutineProfile returns n, the number of records in the active goroutine
// stack profile. If len(p) >= n, GoroutineProfile copies the profile into p and
// returns n, true. If len(p) < n, GoroutineProfile does not change p and
// returns n, false.
//
// Most clients should use the runtime/pprof package instead of calling
// GoroutineProfile directly.
func GoroutineProfile(p []StackRecord) (n int, ok bool)

// Gosched yields the processor, allowing other goroutines to run.  It does not
// suspend the current goroutine, so execution resumes automatically.

// Gosched yields the processor, allowing other goroutines to run. It does not
// suspend the current goroutine, so execution resumes automatically.
func Gosched()

// KeepAlive marks its argument as currently reachable.
// This ensures that the object is not freed, and its finalizer is not run,
// before the point in the program where KeepAlive is called.
//
// A very simplified example showing where KeepAlive is required:
//     type File struct { d int }
//     d, err := syscall.Open("/file/path", syscall.O_RDONLY, 0)
//     // ... do something if err != nil ...
//     p := &FILE{d}
//     runtime.SetFinalizer(p, func(p *File) { syscall.Close(p.d) })
//     var buf [10]byte
//     n, err := syscall.Read(p.d, buf[:])
//     // Ensure p is not finalized until Read returns.
//     runtime.KeepAlive(p)
//     // No more uses of p after this point.
//
// Without the KeepAlive call, the finalizer could run at the start of
// syscall.Read, closing the file descriptor before syscall.Read makes
// the actual system call.
func KeepAlive(interface{})

// LockOSThread wires the calling goroutine to its current operating system
// thread. Until the calling goroutine exits or calls UnlockOSThread, it will
// always execute in that thread, and no other goroutine can.
func LockOSThread()

func MSanRead(addr unsafe.Pointer, len int)

func MSanWrite(addr unsafe.Pointer, len int)

// MemProfile returns a profile of memory allocated and freed per allocation
// site.
//
// MemProfile returns n, the number of records in the current memory profile.
// If len(p) >= n, MemProfile copies the profile into p and returns n, true.
// If len(p) < n, MemProfile does not change p and returns n, false.
//
// If inuseZero is true, the profile includes allocation records
// where r.AllocBytes > 0 but r.AllocBytes == r.FreeBytes.
// These are sites where memory was allocated, but it has all
// been released back to the runtime.
//
// The returned profile may be up to two garbage collection cycles old.
// This is to avoid skewing the profile toward allocations; because
// allocations happen in real time but frees are delayed until the garbage
// collector performs sweeping, the profile only accounts for allocations
// that have had a chance to be freed by the garbage collector.
//
// Most clients should use the runtime/pprof package or
// the testing package's -test.memprofile flag instead
// of calling MemProfile directly.
func MemProfile(p []MemProfileRecord, inuseZero bool) (n int, ok bool)

// NumCPU returns the number of logical CPUs usable by the current process.
//
// The set of available CPUs is checked by querying the operating system
// at process startup. Changes to operating system CPU allocation after
// process startup are not reflected.

// NumCPU 返回当前处理器的可用逻辑CPU数。
//
// 可用的 CPU 的设置通过在进程启动时通过向操作系统查询来获取。在进程启动后更改操
// 作系统的 CPU 分配并不会反映出来。
func NumCPU() int

// NumCgoCall returns the number of cgo calls made by the current process.

// NumCgoCall 返回由当前进程创建的cgo调用数。
func NumCgoCall() int64

// NumGoroutine returns the number of goroutines that currently exist.

// NumGoroutine 返回当前存在的Go程数。
func NumGoroutine() int

func RaceAcquire(addr unsafe.Pointer)

// RaceDisable disables handling of race events in the current goroutine.
func RaceDisable()

// RaceEnable re-enables handling of race events in the current goroutine.
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
func ReadMemStats(m *MemStats)

// ReadTrace returns the next chunk of binary tracing data, blocking until data
// is available. If tracing is turned off and all the data accumulated while it
// was on has been returned, ReadTrace returns nil. The caller must copy the
// returned data before calling ReadTrace again.
// ReadTrace must be called from one goroutine at a time.
func ReadTrace() []byte

// SetBlockProfileRate controls the fraction of goroutine blocking events
// that are reported in the blocking profile.  The profiler aims to sample
// an average of one blocking event per rate nanoseconds spent blocked.
//
// To include every blocking event in the profile, pass rate = 1.
// To turn off profiling entirely, pass rate <= 0.

// SetBlockProfileRate controls the fraction of goroutine blocking events
// that are reported in the blocking profile. The profiler aims to sample
// an average of one blocking event per rate nanoseconds spent blocked.
//
// To include every blocking event in the profile, pass rate = 1.
// To turn off profiling entirely, pass rate <= 0.
func SetBlockProfileRate(rate int)

// SetCPUProfileRate sets the CPU profiling rate to hz samples per second. If hz
// <= 0, SetCPUProfileRate turns off profiling. If the profiler is on, the rate
// cannot be changed without first turning it off.
//
// Most clients should use the runtime/pprof package or the testing package's
// -test.cpuprofile flag instead of calling SetCPUProfileRate directly.
func SetCPUProfileRate(hz int)

// SetCgoTraceback records three C functions to use to gather
// traceback information from C code and to convert that traceback
// information into symbolic information. These are used when printing
// stack traces for a program that uses cgo.
//
// The traceback and context functions may be called from a signal
// handler, and must therefore use only async-signal safe functions.
// The symbolizer function may be called while the program is
// crashing, and so must be cautious about using memory.  None of the
// functions may call back into Go.
//
// The context function will be called with a single argument, a
// pointer to a struct:
//
//     struct {
//         Context uintptr
//     }
//
// In C syntax, this struct will be
//
//     struct {
//         uintptr_t Context;
//     };
//
// If the Context field is 0, the context function is being called to
// record the current traceback context. It should record in the
// Context field whatever information is needed about the current
// point of execution to later produce a stack trace, probably the
// stack pointer and PC. In this case the context function will be
// called from C code.
//
// If the Context field is not 0, then it is a value returned by a
// previous call to the context function. This case is called when the
// context is no longer needed; that is, when the Go code is returning
// to its C code caller. This permits permits the context function to
// release any associated resources.
//
// While it would be correct for the context function to record a
// complete a stack trace whenever it is called, and simply copy that
// out in the traceback function, in a typical program the context
// function will be called many times without ever recording a
// traceback for that context. Recording a complete stack trace in a
// call to the context function is likely to be inefficient.
//
// The traceback function will be called with a single argument, a
// pointer to a struct:
//
//     struct {
//         Context uintptr
//         Buf     *uintptr
//         Max     uintptr
//     }
//
// In C syntax, this struct will be
//
//     struct {
//         uintptr_t  Context;
//         uintptr_t* Buf;
//         uintptr_t  Max;
//     };
//
// The Context field will be zero to gather a traceback from the
// current program execution point. In this case, the traceback
// function will be called from C code.
//
// Otherwise Context will be a value previously returned by a call to
// the context function. The traceback function should gather a stack
// trace from that saved point in the program execution. The traceback
// function may be called from an execution thread other than the one
// that recorded the context, but only when the context is known to be
// valid and unchanging. The traceback function may also be called
// deeper in the call stack on the same thread that recorded the
// context. The traceback function may be called multiple times with
// the same Context value; it will usually be appropriate to cache the
// result, if possible, the first time this is called for a specific
// context value.
//
// Buf is where the traceback information should be stored. It should
// be PC values, such that Buf[0] is the PC of the caller, Buf[1] is
// the PC of that function's caller, and so on.  Max is the maximum
// number of entries to store.  The function should store a zero to
// indicate the top of the stack, or that the caller is on a different
// stack, presumably a Go stack.
//
// Unlike runtime.Callers, the PC values returned should, when passed
// to the symbolizer function, return the file/line of the call
// instruction.  No additional subtraction is required or appropriate.
//
// The symbolizer function will be called with a single argument, a
// pointer to a struct:
//
//     struct {
//         PC      uintptr // program counter to fetch information for
//         File    *byte   // file name (NUL terminated)
//         Lineno  uintptr // line number
//         Func    *byte   // function name (NUL terminated)
//         Entry   uintptr // function entry point
//         More    uintptr // set non-zero if more info for this PC
//         Data    uintptr // unused by runtime, available for function
//     }
//
// In C syntax, this struct will be
//
//     struct {
//         uintptr_t PC;
//         char*     File;
//         uintptr_t Lineno;
//         char*     Func;
//         uintptr_t Entry;
//         uintptr_t More;
//         uintptr_t Data;
//     };
//
// The PC field will be a value returned by a call to the traceback
// function.
//
// The first time the function is called for a particular traceback,
// all the fields except PC will be 0. The function should fill in the
// other fields if possible, setting them to 0/nil if the information
// is not available. The Data field may be used to store any useful
// information across calls. The More field should be set to non-zero
// if there is more information for this PC, zero otherwise. If More
// is set non-zero, the function will be called again with the same
// PC, and may return different information (this is intended for use
// with inlined functions). If More is zero, the function will be
// called with the next PC value in the traceback. When the traceback
// is complete, the function will be called once more with PC set to
// zero; this may be used to free any information. Each call will
// leave the fields of the struct set to the same values they had upon
// return, except for the PC field when the More field is zero. The
// function must not keep a copy of the struct pointer between calls.
//
// When calling SetCgoTraceback, the version argument is the version
// number of the structs that the functions expect to receive.
// Currently this must be zero.
//
// The symbolizer function may be nil, in which case the results of
// the traceback function will be displayed as numbers. If the
// traceback function is nil, the symbolizer function will never be
// called. The context function may be nil, in which case the
// traceback function will only be called with the context field set
// to zero.  If the context function is nil, then calls from Go to C
// to Go will not show a traceback for the C portion of the call stack.
func SetCgoTraceback(version int, traceback, context, symbolizer unsafe.Pointer)

// SetFinalizer sets the finalizer associated with x to f.
// When the garbage collector finds an unreachable block
// with an associated finalizer, it clears the association and runs
// f(x) in a separate goroutine.  This makes x reachable again, but
// now without an associated finalizer.  Assuming that SetFinalizer
// is not called again, the next time the garbage collector sees
// that x is unreachable, it will free x.
//
// SetFinalizer(x, nil) clears any finalizer associated with x.
//
// The argument x must be a pointer to an object allocated by
// calling new or by taking the address of a composite literal.
// The argument f must be a function that takes a single argument
// to which x's type can be assigned, and can have arbitrary ignored return
// values. If either of these is not true, SetFinalizer aborts the
// program.
//
// Finalizers are run in dependency order: if A points at B, both have
// finalizers, and they are otherwise unreachable, only the finalizer
// for A runs; once A is freed, the finalizer for B can run.
// If a cyclic structure includes a block with a finalizer, that
// cycle is not guaranteed to be garbage collected and the finalizer
// is not guaranteed to run, because there is no ordering that
// respects the dependencies.
//
// The finalizer for x is scheduled to run at some arbitrary time after
// x becomes unreachable.
// There is no guarantee that finalizers will run before a program exits,
// so typically they are useful only for releasing non-memory resources
// associated with an object during a long-running program.
// For example, an os.File object could use a finalizer to close the
// associated operating system file descriptor when a program discards
// an os.File without calling Close, but it would be a mistake
// to depend on a finalizer to flush an in-memory I/O buffer such as a
// bufio.Writer, because the buffer would not be flushed at program exit.
//
// It is not guaranteed that a finalizer will run if the size of *x is
// zero bytes.
//
// It is not guaranteed that a finalizer will run for objects allocated
// in initializers for package-level variables. Such objects may be
// linker-allocated, not heap-allocated.
//
// A single goroutine runs all finalizers for a program, sequentially.
// If a finalizer must run for a long time, it should do so by starting
// a new goroutine.

// SetFinalizer 为 f 设置与 x 相关联的终结器。 当垃圾回收器找到一个无法访问的块
// 及与其相关联的终结器时，就会清理该关联， 并在一个独立的Go程中运行f(x)。这会使
// x 再次变得可访问，但现在没有了相关联的终结器。 假设 SetFinalizer 未被再次调
// 用，当下一次垃圾回收器发现 x 无法访问时，就会释放 x。
//
// SetFinalizer(x, nil) 会清理任何与 x 相关联的终结器。
//
// 实参 x 必须是一个对象的指针，该对象通过调用新的或获取一个复合字面地址来分配。
// 实参 f 必须是一个函数，该函数获取一个 x 的类型的单一实参，并拥有可任意忽略的
// 返回值。 只要这些条件有一个不满足，SetFinalizer 就会跳过该程序。
//
// 终结器按照依赖顺序运行：若 A 指向 B，则二者都有终结器，当只有 A 的终结器运行
// 时， 它们才无法访问；一旦 A 被释放，则 B 的终结器便可运行。若循环依赖的结构包
// 含块及其终结器， 则该循环并不能保证被垃圾回收，而其终结器并不能保证运行，这是
// 因为其依赖没有顺序。
//
// x 的终结器预定为在 x 无法访问后的任意时刻运行。无法保证终结器会在程序退出前运
// 行， 因此它们通常只在长时间运行的程序中释放一个关联至对象的非内存资源时使用。
// 例如，当程序丢弃 os.File 而没有调用 Close 时，该 os.File 对象便可使用一个终结
// 器 来关闭与其相关联的操作系统文件描述符，但依赖终结器去刷新一个内存中的I/O缓
// 存是错误的， 因为该缓存不会在程序退出时被刷新。
//
// 一个程序的单个Go程会按顺序运行所有的终结器。若某个终结器需要长时间运行， 它应
// 当通过开始一个新的Go程来继续。 TODO(osc): 仍需校对及语句优化
func SetFinalizer(obj interface{}, finalizer interface{})

// Stack formats a stack trace of the calling goroutine into buf
// and returns the number of bytes written to buf.
// If all is true, Stack formats stack traces of all other goroutines
// into buf after the trace for the current goroutine.
func Stack(buf []byte, all bool) int

// StartTrace enables tracing for the current process.
// While tracing, the data will be buffered and available via ReadTrace.
// StartTrace returns an error if tracing is already enabled.
// Most clients should use the runtime/trace package or the testing package's
// -test.trace flag instead of calling StartTrace directly.
func StartTrace() error

// StopTrace stops tracing, if it was previously enabled.
// StopTrace only returns after all the reads for the trace have completed.
func StopTrace()

// ThreadCreateProfile returns n, the number of records in the thread creation
// profile. If len(p) >= n, ThreadCreateProfile copies the profile into p and
// returns n, true. If len(p) < n, ThreadCreateProfile does not change p and
// returns n, false.
//
// Most clients should use the runtime/pprof package instead of calling
// ThreadCreateProfile directly.
func ThreadCreateProfile(p []StackRecord) (n int, ok bool)

// UnlockOSThread unwires the calling goroutine from its fixed operating system
// thread. If the calling goroutine has not called LockOSThread, UnlockOSThread
// is a no-op.
func UnlockOSThread()

// Version returns the Go tree's version string.
// It is either the commit hash and date at the time of the build or,
// when possible, a release tag like "go1.3".

// Version 返回Go目录树的版本字符串。
// 它一般是一个提交散列值及其构建时间，也可能是一个类似于 "go1.3" 的发行标注。
func Version() string

// Next returns frame information for the next caller.
// If more is false, there are no more callers (the Frame value is valid).
func (*Frames) Next() (frame Frame, more bool)

// Entry returns the entry address of the function.
func (*Func) Entry() uintptr

// FileLine returns the file name and line number of the
// source code corresponding to the program counter pc.
// The result will not be accurate if pc is not a program
// counter within f.
func (*Func) FileLine(pc uintptr) (file string, line int)

// Name returns the name of the function.
func (*Func) Name() string

// InUseBytes returns the number of bytes in use (AllocBytes - FreeBytes).
func (*MemProfileRecord) InUseBytes() int64

// InUseObjects returns the number of objects in use (AllocObjects -
// FreeObjects).
func (*MemProfileRecord) InUseObjects() int64

// Stack returns the stack trace associated with the record,
// a prefix of r.Stack0.
func (*MemProfileRecord) Stack() []uintptr

// Stack returns the stack trace associated with the record,
// a prefix of r.Stack0.
func (*StackRecord) Stack() []uintptr

func (*TypeAssertionError) Error() string

func (*TypeAssertionError) RuntimeError()

