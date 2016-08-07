// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package pprof writes runtime profiling data in the format expected
// by the pprof visualization tool.
// For more information about pprof, see
// http://code.google.com/p/google-perftools/.

// pprof 包按照可视化工具 pprof 所要求的格式写出运行时分析数据. 更多有关 pprof
// 的信息见 http://code.google.com/p/google-perftools/。
package pprof

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "runtime"
    "sort"
    "strings"
    "sync"
    "text/tabwriter"
)

// A Profile is a collection of stack traces showing the call sequences that led
// to instances of a particular event, such as allocation. Packages can create
// and maintain their own profiles; the most common use is for tracking
// resources that must be explicitly closed, such as files or network
// connections.
//
// A Profile's methods can be called from multiple goroutines simultaneously.
//
// Each Profile has a unique name. A few profiles are predefined:
//
//     goroutine    - stack traces of all current goroutines
//     heap         - a sampling of all heap allocations
//     threadcreate - stack traces that led to the creation of new OS threads
//     block        - stack traces that led to blocking on synchronization primitives
//
// These predefined profiles maintain themselves and panic on an explicit Add or
// Remove method call.
//
// The heap profile reports statistics as of the most recently completed garbage
// collection; it elides more recent allocation to avoid skewing the profile
// away from live data and toward garbage. If there has been no garbage
// collection at all, the heap profile reports all known allocations. This
// exception helps mainly in programs running without garbage collection
// enabled, usually for debugging purposes.
//
// The CPU profile is not available as a Profile. It has a special API, the
// StartCPUProfile and StopCPUProfile functions, because it streams output to a
// writer during profiling.

// Profile 是一个栈跟踪的集合，它显示了引导特定事件实例的调用序列，例如分配。 包
// 可以创建并维护它们自己的分析，它一般用于跟踪必须被显式关闭的资源，例如文件或
// 网络连接。
//
// 一个 Profile 的方法可被多个Go程同时调用。
//
// 每个 Profile 都有唯一的名称。有些 Profile 是预定义的：
//
//     goroutine    - 所有当前Go程的栈跟踪
//     heap         - 所有堆分配的采样
//     threadcreate - 引导新OS的线程创建的栈跟踪
//     block        - 引导同步原语中阻塞的栈跟踪
//
// 这些预声明分析并不能作为 Profile 使用。它有专门的API，即 StartCPUProfile 和
// StopCPUProfile 函数，因为它在分析时是以流的形式输出到写入器的。
type Profile struct {
}

// Lookup returns the profile with the given name, or nil if no such profile
// exists.

// Lookup
// 返回给定名称的分析，若不存在该分析，则返回 nil。
func Lookup(name string) *Profile

// NewProfile creates a new profile with the given name.
// If a profile with that name already exists, NewProfile panics.
// The convention is to use a 'import/path.' prefix to create
// separate name spaces for each package.

// NewProfile 以给定的名称创建一个新的分析。 若拥有该名称的分析已存在，
// NewProfile 就会引起恐慌。 约定使用一个 'import/path' 导入路径前缀来为每个包创
// 建单独的命名空间。
func NewProfile(name string) *Profile

// Profiles returns a slice of all the known profiles, sorted by name.

// Profiles
// 返回所有已知分析的切片，按名称排序。
func Profiles() []*Profile

// StartCPUProfile enables CPU profiling for the current process.
// While profiling, the profile will be buffered and written to w.
// StartCPUProfile returns an error if profiling is already enabled.
//
// On Unix-like systems, StartCPUProfile does not work by default for
// Go code built with -buildmode=c-archive or -buildmode=c-shared.
// StartCPUProfile relies on the SIGPROF signal, but that signal will
// be delivered to the main program's SIGPROF signal handler (if any)
// not to the one used by Go.  To make it work, call os/signal.Notify
// for syscall.SIGPROF, but note that doing so may break any profiling
// being done by the main program.

// StartCPUProfile 为当前进程开启CPU分析。 在分析时，分析报告会缓存并写入到 w 中
// 。若分析已经开启，StartCPUProfile 就会返回错误。
func StartCPUProfile(w io.Writer) error

// StopCPUProfile stops the current CPU profile, if any.
// StopCPUProfile only returns after all the writes for the
// profile have completed.

// StopCPUProfile 会停止当前的CPU分析，如果有的话。 StopCPUProfile
// 只会在所有的分析报告写入完毕后才会返回。
func StopCPUProfile()

// WriteHeapProfile is shorthand for Lookup("heap").WriteTo(w, 0).
// It is preserved for backwards compatibility.

// WriteHeapProfile 是 Lookup("heap").WriteTo(w, 0) 的简写。
// 它是为了保持向后兼容性而存在的。
func WriteHeapProfile(w io.Writer) error

// Add adds the current execution stack to the profile, associated with value.
// Add stores value in an internal map, so value must be suitable for use as a
// map key and will not be garbage collected until the corresponding call to
// Remove. Add panics if the profile already contains a stack for value.
//
// The skip parameter has the same meaning as runtime.Caller's skip and controls
// where the stack trace begins. Passing skip=0 begins the trace in the function
// calling Add. For example, given this execution stack:
//
//     Add
//     called from rpc.NewClient
//     called from mypkg.Run
//     called from main.main
//
// Passing skip=0 begins the stack trace at the call to Add inside
// rpc.NewClient. Passing skip=1 begins the stack trace at the call to NewClient
// inside mypkg.Run.

// Add 将当前与值相关联的执行栈添加到该分析中。 Add 在一个内部映射中存储值，因此
// 值必须适于用作映射键，且在对应的 Remove 调用之前不会被垃圾收集。若分析已经包
// 含了值的栈，Add 就会引发恐慌。
//
// skip 形参与 runtime.Caller 的 skip 意思相同，它用于控制栈跟踪从哪里开始。 传
// 入 skip=0 会从函数调用 Add 处开始跟踪。例如，给定以下执行栈：
//
//     Add
//     调用自 rpc.NewClient
//     调用自 mypkg.Run
//     调用自 main.main
//
// 传入 skip=0 会从 rpc.NewClient 中的 Add 调用处开始栈跟踪。 传入 skip=1 会从
// mypkg.Run 中的 NewClient 调用处开始栈跟踪。
func (*Profile) Add(value interface{}, skip int)

// Count returns the number of execution stacks currently in the profile.

// Count 返回该分析中当前执行栈的数量。
func (*Profile) Count() int

// Name returns this profile's name, which can be passed to Lookup to reobtain
// the profile.

// Name 返回该分析的名称，它可被传入 Lookup 来重新获取该分析。
func (*Profile) Name() string

// Remove removes the execution stack associated with value from the profile.
// It is a no-op if the value is not in the profile.

// Remove 从该分析中移除与值 value 相关联的执行栈。 若值 value
// 不在此分析中，则为空操作。
func (*Profile) Remove(value interface{})

// WriteTo writes a pprof-formatted snapshot of the profile to w.
// If a write to w returns an error, WriteTo returns that error.
// Otherwise, WriteTo returns nil.
//
// The debug parameter enables additional output.
// Passing debug=0 prints only the hexadecimal addresses that pprof needs.
// Passing debug=1 adds comments translating addresses to function names
// and line numbers, so that a programmer can read the profile without tools.
//
// The predefined profiles may assign meaning to other debug values;
// for example, when printing the "goroutine" profile, debug=2 means to
// print the goroutine stacks in the same form that a Go program uses
// when dying due to an unrecovered panic.

// WriteTo 将pprof格式的分析快照写入 w 中。 若一个向 w 的写入返回一个错误，
// WriteTo 就会返回该错误。 否则，WriteTo 就会返回 nil。
//
// debug 形参用于开启附加的输出。 传入 debug=0 只会打印pprof所需要的十六进制地址
// 。 传入 debug=1 会将地址翻译为函数名和行号并添加注释，以便让程序员无需工具阅
// 读分析报告。
//
// 预声明分析报告可为其它 debug 值赋予含义；例如，当打印“Go程”的分析报告时，
// debug=2 意为：由于不可恢复的恐慌而濒临崩溃时，使用与Go程序相同的形式打印Go程
// 的栈信息。
func (*Profile) WriteTo(w io.Writer, debug int) error

