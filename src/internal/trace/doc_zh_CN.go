// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package trace

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
    "unsafe"
)

// Event types in the trace.
// Verbatim copy from src/runtime/trace.go.
const (
	EvNone           = 0  // unused
	EvBatch          = 1  // start of per-P batch of events [pid, timestamp]
	EvFrequency      = 2  // contains tracer timer frequency [frequency (ticks per second)]
	EvStack          = 3  // stack [stack id, number of PCs, array of {PC, func string ID, file string ID, line}]
	EvGomaxprocs     = 4  // current value of GOMAXPROCS [timestamp, GOMAXPROCS, stack id]
	EvProcStart      = 5  // start of P [timestamp, thread id]
	EvProcStop       = 6  // stop of P [timestamp]
	EvGCStart        = 7  // GC start [timestamp, seq, stack id]
	EvGCDone         = 8  // GC done [timestamp]
	EvGCScanStart    = 9  // GC scan start [timestamp]
	EvGCScanDone     = 10 // GC scan done [timestamp]
	EvGCSweepStart   = 11 // GC sweep start [timestamp, stack id]
	EvGCSweepDone    = 12 // GC sweep done [timestamp]
	EvGoCreate       = 13 // goroutine creation [timestamp, new goroutine id, new stack id, stack id]
	EvGoStart        = 14 // goroutine starts running [timestamp, goroutine id, seq]
	EvGoEnd          = 15 // goroutine ends [timestamp]
	EvGoStop         = 16 // goroutine stops (like in select{}) [timestamp, stack]
	EvGoSched        = 17 // goroutine calls Gosched [timestamp, stack]
	EvGoPreempt      = 18 // goroutine is preempted [timestamp, stack]
	EvGoSleep        = 19 // goroutine calls Sleep [timestamp, stack]
	EvGoBlock        = 20 // goroutine blocks [timestamp, stack]
	EvGoUnblock      = 21 // goroutine is unblocked [timestamp, goroutine id, seq, stack]
	EvGoBlockSend    = 22 // goroutine blocks on chan send [timestamp, stack]
	EvGoBlockRecv    = 23 // goroutine blocks on chan recv [timestamp, stack]
	EvGoBlockSelect  = 24 // goroutine blocks on select [timestamp, stack]
	EvGoBlockSync    = 25 // goroutine blocks on Mutex/RWMutex [timestamp, stack]
	EvGoBlockCond    = 26 // goroutine blocks on Cond [timestamp, stack]
	EvGoBlockNet     = 27 // goroutine blocks on network [timestamp, stack]
	EvGoSysCall      = 28 // syscall enter [timestamp, stack]
	EvGoSysExit      = 29 // syscall exit [timestamp, goroutine id, seq, real timestamp]
	EvGoSysBlock     = 30 // syscall blocks [timestamp]
	EvGoWaiting      = 31 // denotes that goroutine is blocked when tracing starts [timestamp, goroutine id]
	EvGoInSyscall    = 32 // denotes that goroutine is in syscall when tracing starts [timestamp, goroutine id]
	EvHeapAlloc      = 33 // memstats.heap_live change [timestamp, heap_alloc]
	EvNextGC         = 34 // memstats.next_gc change [timestamp, next_gc]
	EvTimerGoroutine = 35 // denotes timer goroutine [timer goroutine id]
	EvFutileWakeup   = 36 // denotes that the previous wakeup of this goroutine was futile [timestamp]
	EvString         = 37 // string dictionary entry [ID, length, string]
	EvGoStartLocal   = 38 // goroutine starts running on the same P as the last event [timestamp, goroutine id]
	EvGoUnblockLocal = 39 // goroutine is unblocked on the same P as the last event [timestamp, goroutine id, stack]
	EvGoSysExitLocal = 40 // syscall exit on the same P as the last event [timestamp, goroutine id, real timestamp]
	EvCount          = 41
)



const (
	// Special P identifiers:
	FakeP    = 1000000 + iota
	TimerP   // depicts timer unblocks
	NetpollP // depicts network unblocks
	SyscallP // depicts returns from syscalls

)


// BreakTimestampsForTesting causes the parser to randomly alter timestamps (for
// testing of broken cputicks).
var BreakTimestampsForTesting bool


// ErrTimeOrder is returned by Parse when the trace contains
// time stamps that do not respect actual event ordering.
var ErrTimeOrder = fmt.Errorf("time stamps out of order")



var EventDescriptions = [EvCount]struct {
	Name       string
	minVersion int
	Stack      bool
	Args       []string
}{
	EvNone:           {"None", 1005, false, []string{}},
	EvBatch:          {"Batch", 1005, false, []string{"p", "ticks"}},
	EvFrequency:      {"Frequency", 1005, false, []string{"freq"}},
	EvStack:          {"Stack", 1005, false, []string{"id", "siz"}},
	EvGomaxprocs:     {"Gomaxprocs", 1005, true, []string{"procs"}},
	EvProcStart:      {"ProcStart", 1005, false, []string{"thread"}},
	EvProcStop:       {"ProcStop", 1005, false, []string{}},
	EvGCStart:        {"GCStart", 1005, true, []string{"seq"}},
	EvGCDone:         {"GCDone", 1005, false, []string{}},
	EvGCScanStart:    {"GCScanStart", 1005, false, []string{}},
	EvGCScanDone:     {"GCScanDone", 1005, false, []string{}},
	EvGCSweepStart:   {"GCSweepStart", 1005, true, []string{}},
	EvGCSweepDone:    {"GCSweepDone", 1005, false, []string{}},
	EvGoCreate:       {"GoCreate", 1005, true, []string{"g", "stack"}},
	EvGoStart:        {"GoStart", 1005, false, []string{"g", "seq"}},
	EvGoEnd:          {"GoEnd", 1005, false, []string{}},
	EvGoStop:         {"GoStop", 1005, true, []string{}},
	EvGoSched:        {"GoSched", 1005, true, []string{}},
	EvGoPreempt:      {"GoPreempt", 1005, true, []string{}},
	EvGoSleep:        {"GoSleep", 1005, true, []string{}},
	EvGoBlock:        {"GoBlock", 1005, true, []string{}},
	EvGoUnblock:      {"GoUnblock", 1005, true, []string{"g", "seq"}},
	EvGoBlockSend:    {"GoBlockSend", 1005, true, []string{}},
	EvGoBlockRecv:    {"GoBlockRecv", 1005, true, []string{}},
	EvGoBlockSelect:  {"GoBlockSelect", 1005, true, []string{}},
	EvGoBlockSync:    {"GoBlockSync", 1005, true, []string{}},
	EvGoBlockCond:    {"GoBlockCond", 1005, true, []string{}},
	EvGoBlockNet:     {"GoBlockNet", 1005, true, []string{}},
	EvGoSysCall:      {"GoSysCall", 1005, true, []string{}},
	EvGoSysExit:      {"GoSysExit", 1005, false, []string{"g", "seq", "ts"}},
	EvGoSysBlock:     {"GoSysBlock", 1005, false, []string{}},
	EvGoWaiting:      {"GoWaiting", 1005, false, []string{"g"}},
	EvGoInSyscall:    {"GoInSyscall", 1005, false, []string{"g"}},
	EvHeapAlloc:      {"HeapAlloc", 1005, false, []string{"mem"}},
	EvNextGC:         {"NextGC", 1005, false, []string{"mem"}},
	EvTimerGoroutine: {"TimerGoroutine", 1005, false, []string{"g"}},
	EvFutileWakeup:   {"FutileWakeup", 1005, false, []string{}},
	EvString:         {"String", 1007, false, []string{}},
	EvGoStartLocal:   {"GoStartLocal", 1007, false, []string{"g"}},
	EvGoUnblockLocal: {"GoUnblockLocal", 1007, true, []string{"g"}},
	EvGoSysExitLocal: {"GoSysExitLocal", 1007, false, []string{"g", "ts"}},
}


// Event describes one event in the trace.
type Event struct {
	Off   int       // offset in input file (for debugging and error reporting)
	Type  byte      // one of Ev*
	seq   int64     // sequence number
	Ts    int64     // timestamp in nanoseconds
	P     int       // P on which the event happened (can be one of TimerP, NetpollP, SyscallP)
	G     uint64    // G on which the event happened
	StkID uint64    // unique stack ID
	Stk   []*Frame  // stack trace (can be empty)
	Args  [3]uint64 // event-type-specific arguments
	// linked event (can be nil), depends on event type:
	// for GCStart: the GCStop
	// for GCScanStart: the GCScanDone
	// for GCSweepStart: the GCSweepDone
	// for GoCreate: first GoStart of the created goroutine
	// for GoStart: the associated GoEnd, GoBlock or other blocking event
	// for GoSched/GoPreempt: the next GoStart
	// for GoBlock and other blocking events: the unblock event
	// for GoUnblock: the associated GoStart
	// for blocking GoSysCall: the associated GoSysExit
	// for GoSysExit: the next GoStart
	Link *Event
}


// Frame is a frame in stack traces.
type Frame struct {
	PC   uint64
	Fn   string
	File string
	Line int
}


// GDesc contains statistics about execution of a single goroutine.
type GDesc struct {
	ID           uint64
	Name         string
	PC           uint64
	CreationTime int64
	StartTime    int64
	EndTime      int64

	ExecTime      int64
	SchedWaitTime int64
	IOTime        int64
	BlockTime     int64
	SyscallTime   int64
	GCTime        int64
	SweepTime     int64
	TotalTime     int64
}


// GoroutineStats generates statistics for all goroutines in the trace.
func GoroutineStats(events []*Event) map[uint64]*GDesc

// Parse parses, post-processes and verifies the trace.
func Parse(r io.Reader, bin string) ([]*Event, error)

// Print dumps events to stdout. For debugging.
func Print(events []*Event)

// PrintEvent dumps the event to stdout. For debugging.
func PrintEvent(ev *Event)

// RelatedGoroutines finds a set of goroutines related to goroutine goid.
func RelatedGoroutines(events []*Event, goid uint64) map[uint64]bool

