// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package debug contains facilities for programs to debug themselves while
// they are running.

// debug 包含有程序在运行时调试其自身的功能.
package debug

import (
    "os"
    "runtime"
    "sort"
    "time"
)

// GCStats collect information about recent garbage collections.
type GCStats struct {
    LastGC         time.Time       // time of last collection
    NumGC          int64           // number of garbage collections
    PauseTotal     time.Duration   // total pause for all collections
    Pause          []time.Duration // pause history, most recent first
    PauseEnd       []time.Time     // pause end times history, most recent first
    PauseQuantiles []time.Duration
}

// FreeOSMemory forces a garbage collection followed by an attempt to return as
// much memory to the operating system as possible. (Even if this is not called,
// the runtime gradually returns memory to the operating system in a background
// task.)
func FreeOSMemory()

// PrintStack prints to standard error the stack trace returned by
// runtime.Stack.

// PrintStack 将 Stack
// 返回的栈跟踪信息打印到标准错误输出。
func PrintStack()

// ReadGCStats reads statistics about garbage collection into stats. The number
// of entries in the pause history is system-dependent; stats.Pause slice will
// be reused if large enough, reallocated otherwise. ReadGCStats may use the
// full capacity of the stats.Pause slice. If stats.PauseQuantiles is non-empty,
// ReadGCStats fills it with quantiles summarizing the distribution of pause
// time. For example, if len(stats.PauseQuantiles) is 5, it will be filled with
// the minimum, 25%, 50%, 75%, and maximum pause times.
func ReadGCStats(stats *GCStats)

// SetGCPercent sets the garbage collection target percentage: a collection is
// triggered when the ratio of freshly allocated data to live data remaining
// after the previous collection reaches this percentage. SetGCPercent returns
// the previous setting. The initial setting is the value of the GOGC
// environment variable at startup, or 100 if the variable is not set. A
// negative percentage disables garbage collection.
func SetGCPercent(percent int) int

// SetMaxStack sets the maximum amount of memory that can be used by a single
// goroutine stack. If any goroutine exceeds this limit while growing its stack,
// the program crashes. SetMaxStack returns the previous setting. The initial
// setting is 1 GB on 64-bit systems, 250 MB on 32-bit systems.
//
// SetMaxStack is useful mainly for limiting the damage done by goroutines that
// enter an infinite recursion. It only limits future stack growth.
func SetMaxStack(bytes int) int

// SetMaxThreads sets the maximum number of operating system threads that the Go
// program can use. If it attempts to use more than this many, the program
// crashes. SetMaxThreads returns the previous setting. The initial setting is
// 10,000 threads.
//
// The limit controls the number of operating system threads, not the number of
// goroutines. A Go program creates a new thread only when a goroutine is ready
// to run but all the existing threads are blocked in system calls, cgo calls,
// or are locked to other goroutines due to use of runtime.LockOSThread.
//
// SetMaxThreads is useful mainly for limiting the damage done by programs that
// create an unbounded number of threads. The idea is to take down the program
// before it takes down the operating system.
func SetMaxThreads(threads int) int

// SetPanicOnFault controls the runtime's behavior when a program faults at an
// unexpected (non-nil) address. Such faults are typically caused by bugs such
// as runtime memory corruption, so the default response is to crash the
// program. Programs working with memory-mapped files or unsafe manipulation of
// memory may cause faults at non-nil addresses in less dramatic situations;
// SetPanicOnFault allows such programs to request that the runtime trigger only
// a panic, not a crash. SetPanicOnFault applies only to the current goroutine.
// It returns the previous setting.
func SetPanicOnFault(enabled bool) bool

// Stack returns a formatted stack trace of the goroutine that calls it. It
// calls runtime.Stack with a large enough buffer to capture the entire trace.

// Stack 返回格式化的Go程调用的栈跟踪信息。
// 对于每一个例程，它包括来源行的信息和 PC 值，然后尝试获取，对于Go函数，
// 则是调用的函数或方法及其包含请求的行的文本。
//
// 此函数并不赞成使用。请使用 runtime 包中的 Stack 代替。
func Stack() []byte

// WriteHeapDump writes a description of the heap and the objects in
// it to the given file descriptor.
// The heap dump format is defined at https://golang.org/s/go13heapdump.

// WriteHeapDump writes a description of the heap and the objects in it to the
// given file descriptor. The heap dump format is defined at
// http://golang.org/s/go13heapdump.
func WriteHeapDump(fd uintptr)

