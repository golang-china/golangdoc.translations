// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package trace // import "internal/trace"

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
    "testing"
)

// Event types in the trace.
// Verbatim copy from src/runtime/trace.go.
const (
    EvNone           = 0  // unused
    EvBatch          = 1  // start of per-P batch of events [pid, timestamp]
    EvFrequency      = 2  // contains tracer timer frequency [frequency (ticks per second)]
    EvStack          = 3  // stack [stack id, number of PCs, array of PCs]
    EvGomaxprocs     = 4  // current value of GOMAXPROCS [timestamp, GOMAXPROCS, stack id]
    EvProcStart      = 5  // start of P [timestamp, thread id]
    EvProcStop       = 6  // stop of P [timestamp]
    EvGCStart        = 7  // GC start [timestamp, stack id]
    EvGCDone         = 8  // GC done [timestamp]
    EvGCScanStart    = 9  // GC scan start [timestamp]
    EvGCScanDone     = 10 // GC scan done [timestamp]
    EvGCSweepStart   = 11 // GC sweep start [timestamp, stack id]
    EvGCSweepDone    = 12 // GC sweep done [timestamp]
    EvGoCreate       = 13 // goroutine creation [timestamp, new goroutine id, start PC, stack id]
    EvGoStart        = 14 // goroutine starts running [timestamp, goroutine id]
    EvGoEnd          = 15 // goroutine ends [timestamp]
    EvGoStop         = 16 // goroutine stops (like in select{}) [timestamp, stack]
    EvGoSched        = 17 // goroutine calls Gosched [timestamp, stack]
    EvGoPreempt      = 18 // goroutine is preempted [timestamp, stack]
    EvGoSleep        = 19 // goroutine calls Sleep [timestamp, stack]
    EvGoBlock        = 20 // goroutine blocks [timestamp, stack]
    EvGoUnblock      = 21 // goroutine is unblocked [timestamp, goroutine id, stack]
    EvGoBlockSend    = 22 // goroutine blocks on chan send [timestamp, stack]
    EvGoBlockRecv    = 23 // goroutine blocks on chan recv [timestamp, stack]
    EvGoBlockSelect  = 24 // goroutine blocks on select [timestamp, stack]
    EvGoBlockSync    = 25 // goroutine blocks on Mutex/RWMutex [timestamp, stack]
    EvGoBlockCond    = 26 // goroutine blocks on Cond [timestamp, stack]
    EvGoBlockNet     = 27 // goroutine blocks on network [timestamp, stack]
    EvGoSysCall      = 28 // syscall enter [timestamp, stack]
    EvGoSysExit      = 29 // syscall exit [timestamp, goroutine id, real timestamp]
    EvGoSysBlock     = 30 // syscall blocks [timestamp]
    EvGoWaiting      = 31 // denotes that goroutine is blocked when tracing starts [goroutine id]
    EvGoInSyscall    = 32 // denotes that goroutine is in syscall when tracing starts [goroutine id]
    EvHeapAlloc      = 33 // memstats.heap_alloc change [timestamp, heap_alloc]
    EvNextGC         = 34 // memstats.next_gc change [timestamp, next_gc]
    EvTimerGoroutine = 35 // denotes timer goroutine [timer goroutine id]
    EvFutileWakeup   = 36 // denotes that the previous wakeup of this goroutine was futile [timestamp]
    EvCount          = 37
)

const (
    // Special P identifiers:
    FakeP    = 1000000 + iota
    TimerP   // depicts timer unblocks
    NetpollP // depicts network unblocks
    SyscallP // depicts returns from syscalls
)

// ErrTimeOrder is returned by Parse when the trace contains
// time stamps that do not respect actual event ordering.
var ErrTimeOrder = fmt.Errorf("time stamps out of order")

var EventDescriptions = [EvCount]struct {
    Name  string
    Stack bool
    Args  []string
}{
    EvNone:           {"None", false, []string{}},
    EvBatch:          {"Batch", false, []string{"p", "seq", "ticks"}},
    EvFrequency:      {"Frequency", false, []string{"freq", "unused"}},
    EvStack:          {"Stack", false, []string{"id", "siz"}},
    EvGomaxprocs:     {"Gomaxprocs", true, []string{"procs"}},
    EvProcStart:      {"ProcStart", false, []string{"thread"}},
    EvProcStop:       {"ProcStop", false, []string{}},
    EvGCStart:        {"GCStart", true, []string{}},
    EvGCDone:         {"GCDone", false, []string{}},
    EvGCScanStart:    {"GCScanStart", false, []string{}},
    EvGCScanDone:     {"GCScanDone", false, []string{}},
    EvGCSweepStart:   {"GCSweepStart", true, []string{}},
    EvGCSweepDone:    {"GCSweepDone", false, []string{}},
    EvGoCreate:       {"GoCreate", true, []string{"g", "pc"}},
    EvGoStart:        {"GoStart", false, []string{"g"}},
    EvGoEnd:          {"GoEnd", false, []string{}},
    EvGoStop:         {"GoStop", true, []string{}},
    EvGoSched:        {"GoSched", true, []string{}},
    EvGoPreempt:      {"GoPreempt", true, []string{}},
    EvGoSleep:        {"GoSleep", true, []string{}},
    EvGoBlock:        {"GoBlock", true, []string{}},
    EvGoUnblock:      {"GoUnblock", true, []string{"g"}},
    EvGoBlockSend:    {"GoBlockSend", true, []string{}},
    EvGoBlockRecv:    {"GoBlockRecv", true, []string{}},
    EvGoBlockSelect:  {"GoBlockSelect", true, []string{}},
    EvGoBlockSync:    {"GoBlockSync", true, []string{}},
    EvGoBlockCond:    {"GoBlockCond", true, []string{}},
    EvGoBlockNet:     {"GoBlockNet", true, []string{}},
    EvGoSysCall:      {"GoSysCall", true, []string{}},
    EvGoSysExit:      {"GoSysExit", false, []string{"g", "seq", "ts"}},
    EvGoSysBlock:     {"GoSysBlock", false, []string{}},
    EvGoWaiting:      {"GoWaiting", false, []string{"g"}},
    EvGoInSyscall:    {"GoInSyscall", false, []string{"g"}},
    EvHeapAlloc:      {"HeapAlloc", false, []string{"mem"}},
    EvNextGC:         {"NextGC", false, []string{"mem"}},
    EvTimerGoroutine: {"TimerGoroutine", false, []string{"g", "unused"}},
    EvFutileWakeup:   {"FutileWakeup", false, []string{}},
}

// Event describes one event in the trace.
type Event struct {
    Off   int       // offset in input file (for debugging and error reporting)
    Type  byte      // one of Ev*
    Seq   int64     // sequence number
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
func Parse(r io.Reader) ([]*Event, error)

// Print dumps events to stdout. For debugging.
func Print(events []*Event)

// RelatedGoroutines finds a set of goroutines related to goroutine goid.
func RelatedGoroutines(events []*Event, goid uint64) map[uint64]bool

// symbolizeTrace attaches func/file/line info to stack traces.
func Symbolize(events []*Event, bin string) error

func TestCorruptedInputs(t *testing.T)

