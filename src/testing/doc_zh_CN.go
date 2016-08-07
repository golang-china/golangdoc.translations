// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package testing provides support for automated testing of Go packages. It is
// intended to be used in concert with the ``go test'' command, which automates
// execution of any function of the form
//
//     func TestXxx(*testing.T)
//
// where Xxx can be any alphanumeric string (but the first letter must not be in
// [a-z]) and serves to identify the test routine.
//
// Within these functions, use the Error, Fail or related methods to signal
// failure.
//
// To write a new test suite, create a file whose name ends _test.go that
// contains the TestXxx functions as described here. Put the file in the same
// package as the one being tested. The file will be excluded from regular
// package builds but will be included when the ``go test'' command is run. For
// more detail, run ``go help test'' and ``go help testflag''.
//
// Tests and benchmarks may be skipped if not applicable with a call to the Skip
// method of *T and *B:
//
//     func TestTimeConsuming(t *testing.T) {
//         if testing.Short() {
//             t.Skip("skipping test in short mode.")
//         }
//         ...
//     }
//
//
// Benchmarks
//
// Functions of the form
//
//     func BenchmarkXxx(*testing.B)
//
// are considered benchmarks, and are executed by the "go test" command when its
// -bench flag is provided. Benchmarks are run sequentially.
//
// For a description of the testing flags, see
// https://golang.org/cmd/go/#hdr-Description_of_testing_flags.
//
// A sample benchmark function looks like this:
//
//     func BenchmarkHello(b *testing.B) {
//         for i := 0; i < b.N; i++ {
//             fmt.Sprintf("hello")
//         }
//     }
//
// The benchmark function must run the target code b.N times. During benchmark
// execution, b.N is adjusted until the benchmark function lasts long enough to
// be timed reliably. The output
//
//     BenchmarkHello    10000000    282 ns/op
//
// means that the loop ran 10000000 times at a speed of 282 ns per loop.
//
// If a benchmark needs some expensive setup before running, the timer may be
// reset:
//
//     func BenchmarkBigLen(b *testing.B) {
//         big := NewBig()
//         b.ResetTimer()
//         for i := 0; i < b.N; i++ {
//             big.Len()
//         }
//     }
//
// If a benchmark needs to test performance in a parallel setting, it may use
// the RunParallel helper function; such benchmarks are intended to be used with
// the go test -cpu flag:
//
//     func BenchmarkTemplateParallel(b *testing.B) {
//         templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
//         b.RunParallel(func(pb *testing.PB) {
//             var buf bytes.Buffer
//             for pb.Next() {
//                 buf.Reset()
//                 templ.Execute(&buf, "World")
//             }
//         })
//     }
//
//
// Examples
//
// The package also runs and verifies example code. Example functions may
// include a concluding line comment that begins with "Output:" and is compared
// with the standard output of the function when the tests are run. (The
// comparison ignores leading and trailing space.) These are examples of an
// example:
//
//     func ExampleHello() {
//             fmt.Println("hello")
//             // Output: hello
//     }
//
//     func ExampleSalutations() {
//             fmt.Println("hello, and")
//             fmt.Println("goodbye")
//             // Output:
//             // hello, and
//             // goodbye
//     }
//
// Example functions without output comments are compiled but not executed.
//
// The naming convention to declare examples for the package, a function F, a
// type T and method M on type T are:
//
//     func Example() { ... }
//     func ExampleF() { ... }
//     func ExampleT() { ... }
//     func ExampleT_M() { ... }
//
// Multiple example functions for a package/type/function/method may be provided
// by appending a distinct suffix to the name. The suffix must start with a
// lower-case letter.
//
//     func Example_suffix() { ... }
//     func ExampleF_suffix() { ... }
//     func ExampleT_suffix() { ... }
//     func ExampleT_M_suffix() { ... }
//
// The entire test file is presented as the example when it contains a single
// example function, at least one other function, type, variable, or constant
// declaration, and no test or benchmark functions.
//
//
// Main
//
// It is sometimes necessary for a test program to do extra setup or teardown
// before or after testing. It is also sometimes necessary for a test to control
// which code runs on the main thread. To support these and other cases, if a
// test file contains a function:
//
//     func TestMain(m *testing.M)
//
// then the generated test will call TestMain(m) instead of running the tests
// directly. TestMain runs in the main goroutine and can do whatever setup and
// teardown is necessary around a call to m.Run. It should then call os.Exit
// with the result of m.Run. When TestMain is called, flag.Parse has not been
// run. If TestMain depends on command-line flags, including those of the
// testing package, it should call flag.Parse explicitly.
//
// A simple implementation of TestMain is:
//
//     func TestMain(m *testing.M) {
//         flag.Parse()
//         os.Exit(m.Run())
//     }

// Package testing provides support for automated testing of Go packages. It is
// intended to be used in concert with the ``go test'' command, which automates
// execution of any function of the form
//
//     func TestXxx(*testing.T)
//
// where Xxx can be any alphanumeric string (but the first letter must not be in
// [a-z]) and serves to identify the test routine.
//
// Within these functions, use the Error, Fail or related methods to signal
// failure.
//
// To write a new test suite, create a file whose name ends _test.go that
// contains the TestXxx functions as described here. Put the file in the same
// package as the one being tested. The file will be excluded from regular
// package builds but will be included when the ``go test'' command is run. For
// more detail, run ``go help test'' and ``go help testflag''.
//
// Tests and benchmarks may be skipped if not applicable with a call to the Skip
// method of *T and *B:
//
//     func TestTimeConsuming(t *testing.T) {
//         if testing.Short() {
//             t.Skip("skipping test in short mode.")
//         }
//         ...
//     }
//
// Benchmarks
//
// Functions of the form
//
//     func BenchmarkXxx(*testing.B)
//
// are considered benchmarks, and are executed by the "go test" command when its
// -bench flag is provided. Benchmarks are run sequentially.
//
// For a description of the testing flags, see
// http://golang.org/cmd/go/#hdr-Description_of_testing_flags.
//
// A sample benchmark function looks like this:
//
//     func BenchmarkHello(b *testing.B) {
//         for i := 0; i < b.N; i++ {
//             fmt.Sprintf("hello")
//         }
//     }
//
// The benchmark function must run the target code b.N times. During benchark
// execution, b.N is adjusted until the benchmark function lasts long enough to
// be timed reliably. The output
//
//     BenchmarkHello    10000000    282 ns/op
//
// means that the loop ran 10000000 times at a speed of 282 ns per loop.
//
// If a benchmark needs some expensive setup before running, the timer may be
// reset:
//
//     func BenchmarkBigLen(b *testing.B) {
//         big := NewBig()
//         b.ResetTimer()
//         for i := 0; i < b.N; i++ {
//             big.Len()
//         }
//     }
//
// If a benchmark needs to test performance in a parallel setting, it may use
// the RunParallel helper function; such benchmarks are intended to be used with
// the go test -cpu flag:
//
//     func BenchmarkTemplateParallel(b *testing.B) {
//         templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
//         b.RunParallel(func(pb *testing.PB) {
//             var buf bytes.Buffer
//             for pb.Next() {
//                 buf.Reset()
//                 templ.Execute(&buf, "World")
//             }
//         })
//     }
//
// Examples
//
// The package also runs and verifies example code. Example functions may
// include a concluding line comment that begins with "Output:" and is compared
// with the standard output of the function when the tests are run. (The
// comparison ignores leading and trailing space.) These are examples of an
// example:
//
//     func ExampleHello() {
//             fmt.Println("hello")
//             // Output: hello
//     }
//
//     func ExampleSalutations() {
//             fmt.Println("hello, and")
//             fmt.Println("goodbye")
//             // Output:
//             // hello, and
//             // goodbye
//     }
//
// Example functions without output comments are compiled but not executed.
//
// The naming convention to declare examples for the package, a function F, a
// type T and method M on type T are:
//
//     func Example() { ... }
//     func ExampleF() { ... }
//     func ExampleT() { ... }
//     func ExampleT_M() { ... }
//
// Multiple example functions for a package/type/function/method may be provided
// by appending a distinct suffix to the name. The suffix must start with a
// lower-case letter.
//
//     func Example_suffix() { ... }
//     func ExampleF_suffix() { ... }
//     func ExampleT_suffix() { ... }
//     func ExampleT_M_suffix() { ... }
//
// The entire test file is presented as the example when it contains a single
// example function, at least one other function, type, variable, or constant
// declaration, and no test or benchmark functions.
//
// Main
//
// It is sometimes necessary for a test program to do extra setup or teardown
// before or after testing. It is also sometimes necessary for a test to control
// which code runs on the main thread. To support these and other cases, if a
// test file contains a function:
//
//     func TestMain(m *testing.M)
//
// then the generated test will call TestMain(m) instead of running the tests
// directly. TestMain runs in the main goroutine and can do whatever setup and
// teardown is necessary around a call to m.Run. It should then call os.Exit
// with the result of m.Run.
//
// The minimal implementation of TestMain is:
//
//     func TestMain(m *testing.M) { os.Exit(m.Run()) }
//
// In effect, that is the implementation used when no TestMain is explicitly
// defined.
package testing

import (
    "bytes"
    "flag"
    "fmt"
    "io"
    "os"
    "runtime"
    "runtime/debug"
    "runtime/pprof"
    "runtime/trace"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

// B is a type passed to Benchmark functions to manage benchmark timing and to
// specify the number of iterations to run.
//
// A benchmark ends when its Benchmark function returns or calls any of the
// methods FailNow, Fatal, Fatalf, SkipNow, Skip, or Skipf. Those methods must
// be called only from the goroutine running the Benchmark function. The other
// reporting methods, such as the variations of Log and Error, may be called
// simultaneously from multiple goroutines.
//
// Like in tests, benchmark logs are accumulated during execution and dumped to
// standard error when done. Unlike in tests, benchmark logs are always printed,
// so as not to hide output whose existence may be affecting benchmark results.

// B is a type passed to Benchmark functions to manage benchmark timing and to
// specify the number of iterations to run.
type B struct {
    N int
}

// The results of a benchmark run.
type BenchmarkResult struct {
    N         int           // The number of iterations.
    T         time.Duration // The total time taken.
    Bytes     int64         // Bytes processed in one iteration.
    MemAllocs uint64        // The total number of memory allocations.
    MemBytes  uint64        // The total number of bytes allocated.
}

// Cover records information about test coverage checking. NOTE: This struct is
// internal to the testing infrastructure and may change. It is not covered
// (yet) by the Go 1 compatibility guidelines.
type Cover struct {
    Mode            string
    Counters        map[string][]uint32
    Blocks          map[string][]CoverBlock
    CoveredPackages string
}

// CoverBlock records the coverage data for a single basic block. NOTE: This
// struct is internal to the testing infrastructure and may change. It is not
// covered (yet) by the Go 1 compatibility guidelines.
type CoverBlock struct {
    Line0 uint32
    Col0  uint16
    Line1 uint32
    Col1  uint16
    Stmts uint16
}

// An internal type but exported because it is cross-package; part of the
// implementation of the "go test" command.
type InternalBenchmark struct {
    Name string
    F    func(b *B)
}

type InternalExample struct {
    Name   string
    F      func()
    Output string
}

// An internal type but exported because it is cross-package; part of the
// implementation of the "go test" command.
type InternalTest struct {
    Name string
    F    func(*T)
}

// M is a type passed to a TestMain function to run the actual tests.
type M struct {
}

// A PB is used by RunParallel for running parallel benchmarks.
type PB struct {
}

// T is a type passed to Test functions to manage test state and support
// formatted test logs. Logs are accumulated during execution and dumped to
// standard error when done.
//
// A test ends when its Test function returns or calls any of the methods
// FailNow, Fatal, Fatalf, SkipNow, Skip, or Skipf. Those methods, as well as
// the Parallel method, must be called only from the goroutine running the Test
// function.
//
// The other reporting methods, such as the variations of Log and Error, may be
// called simultaneously from multiple goroutines.

// T is a type passed to Test functions to manage test state and support
// formatted test logs. Logs are accumulated during execution and dumped to
// standard error when done.
type T struct {
}

// TB is the interface common to T and B.
type TB interface {
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    Fail()
    FailNow()
    Failed() bool
    Fatal(args ...interface{})
    Fatalf(format string, args ...interface{})
    Log(args ...interface{})
    Logf(format string, args ...interface{})
    Skip(args ...interface{})
    SkipNow()
    Skipf(format string, args ...interface{})
    Skipped() bool
}

// AllocsPerRun returns the average number of allocations during calls to f.
// Although the return value has type float64, it will always be an integral
// value.
//
// To compute the number of allocations, the function will first be run once as
// a warm-up. The average number of allocations over the specified number of
// runs will then be measured and returned.
//
// AllocsPerRun sets GOMAXPROCS to 1 during its measurement and will restore it
// before returning.

// AllocsPerRun returns the average number of allocations during calls to f.
// Although the return value has type float64, it will always be an integral
// value.
//
// To compute the number of allocations, the function will first be run once as
// a warm-up. The average number of allocations over the specified number of
// runs will then be measured and returned.
//
// AllocsPerRun sets GOMAXPROCS to 1 during its measurement and will restore it
// before returning.
func AllocsPerRun(runs int, f func()) (avg float64)

// Benchmark benchmarks a single function. Useful for creating custom benchmarks
// that do not use the "go test" command.
func Benchmark(f func(b *B)) BenchmarkResult

// Coverage reports the current code coverage as a fraction in the range [0, 1].
// If coverage is not enabled, Coverage returns 0.
//
// When running a large set of sequential test cases, checking Coverage after
// each one can be useful for identifying which test cases exercise new code
// paths. It is not a replacement for the reports generated by 'go test -cover'
// and 'go tool cover'.
func Coverage() float64

// An internal function but exported because it is cross-package; part of the
// implementation of the "go test" command.
func Main(matchString func(pat, str string) (bool, error), tests []InternalTest, benchmarks []InternalBenchmark, examples []InternalExample)

// MainStart is meant for use by tests generated by 'go test'. It is not meant
// to be called directly and is not subject to the Go 1 compatibility document.
// It may change signature from release to release.
func MainStart(matchString func(pat, str string) (bool, error), tests []InternalTest, benchmarks []InternalBenchmark, examples []InternalExample) *M

// RegisterCover records the coverage data accumulators for the tests. NOTE:
// This function is internal to the testing infrastructure and may change. It is
// not covered (yet) by the Go 1 compatibility guidelines.
func RegisterCover(c Cover)

// An internal function but exported because it is cross-package; part of the
// implementation of the "go test" command.
func RunBenchmarks(matchString func(pat, str string) (bool, error), benchmarks []InternalBenchmark)

func RunExamples(matchString func(pat, str string) (bool, error), examples []InternalExample) (ok bool)

func RunTests(matchString func(pat, str string) (bool, error), tests []InternalTest) (ok bool)

// Short reports whether the -test.short flag is set.
func Short() bool

// Verbose reports whether the -test.v flag is set.
func Verbose() bool

// Error is equivalent to Log followed by Fail.
func (*B) Error(args ...interface{})

// Errorf is equivalent to Logf followed by Fail.
func (*B) Errorf(format string, args ...interface{})

// Fail marks the function as having failed but continues execution.
func (*B) Fail()

// FailNow marks the function as having failed and stops its execution.
// Execution will continue at the next test or benchmark. FailNow must be called
// from the goroutine running the test or benchmark function, not from other
// goroutines created during the test. Calling FailNow does not stop those other
// goroutines.
func (*B) FailNow()

// Failed reports whether the function has failed.
func (*B) Failed() bool

// Fatal is equivalent to Log followed by FailNow.
func (*B) Fatal(args ...interface{})

// Fatalf is equivalent to Logf followed by FailNow.
func (*B) Fatalf(format string, args ...interface{})

// Log formats its arguments using default formatting, analogous to Println, and
// records the text in the error log. The text will be printed only if the test
// fails or the -test.v flag is set.
func (*B) Log(args ...interface{})

// Logf formats its arguments according to the format, analogous to Printf, and
// records the text in the error log. The text will be printed only if the test
// fails or the -test.v flag is set.
func (*B) Logf(format string, args ...interface{})

// ReportAllocs enables malloc statistics for this benchmark. It is equivalent
// to setting -test.benchmem, but it only affects the benchmark function that
// calls ReportAllocs.
func (*B) ReportAllocs()

// ResetTimer zeros the elapsed benchmark time and memory allocation counters.
// It does not affect whether the timer is running.
func (*B) ResetTimer()

// RunParallel runs a benchmark in parallel. It creates multiple goroutines and
// distributes b.N iterations among them. The number of goroutines defaults to
// GOMAXPROCS. To increase parallelism for non-CPU-bound benchmarks, call
// SetParallelism before RunParallel. RunParallel is usually used with the go
// test -cpu flag.
//
// The body function will be run in each goroutine. It should set up any
// goroutine-local state and then iterate until pb.Next returns false. It should
// not use the StartTimer, StopTimer, or ResetTimer functions, because they have
// global effect.
func (*B) RunParallel(body func(*PB))

// SetBytes records the number of bytes processed in a single operation. If this
// is called, the benchmark will report ns/op and MB/s.
func (*B) SetBytes(n int64)

// SetParallelism sets the number of goroutines used by RunParallel to
// p*GOMAXPROCS. There is usually no need to call SetParallelism for CPU-bound
// benchmarks. If p is less than 1, this call will have no effect.
func (*B) SetParallelism(p int)

// Skip is equivalent to Log followed by SkipNow.
func (*B) Skip(args ...interface{})

// SkipNow marks the test as having been skipped and stops its execution.
// Execution will continue at the next test or benchmark. See also FailNow.
// SkipNow must be called from the goroutine running the test, not from other
// goroutines created during the test. Calling SkipNow does not stop those other
// goroutines.
func (*B) SkipNow()

// Skipf is equivalent to Logf followed by SkipNow.
func (*B) Skipf(format string, args ...interface{})

// Skipped reports whether the test was skipped.
func (*B) Skipped() bool

// StartTimer starts timing a test.  This function is called automatically
// before a benchmark starts, but it can also used to resume timing after
// a call to StopTimer.

// StartTimer starts timing a test. This function is called automatically before
// a benchmark starts, but it can also used to resume timing after a call to
// StopTimer.
func (*B) StartTimer()

// StopTimer stops timing a test.  This can be used to pause the timer
// while performing complex initialization that you don't
// want to measure.

// StopTimer stops timing a test. This can be used to pause the timer while
// performing complex initialization that you don't want to measure.
func (*B) StopTimer()

// Run runs the tests. It returns an exit code to pass to os.Exit.
func (*M) Run() int

// Next reports whether there are more iterations to execute.
func (*PB) Next() bool

// Error is equivalent to Log followed by Fail.
func (*T) Error(args ...interface{})

// Errorf is equivalent to Logf followed by Fail.
func (*T) Errorf(format string, args ...interface{})

// Fail marks the function as having failed but continues execution.
func (*T) Fail()

// FailNow marks the function as having failed and stops its execution.
// Execution will continue at the next test or benchmark. FailNow must be called
// from the goroutine running the test or benchmark function, not from other
// goroutines created during the test. Calling FailNow does not stop those other
// goroutines.
func (*T) FailNow()

// Failed reports whether the function has failed.
func (*T) Failed() bool

// Fatal is equivalent to Log followed by FailNow.
func (*T) Fatal(args ...interface{})

// Fatalf is equivalent to Logf followed by FailNow.
func (*T) Fatalf(format string, args ...interface{})

// Log formats its arguments using default formatting, analogous to Println, and
// records the text in the error log. The text will be printed only if the test
// fails or the -test.v flag is set.
func (*T) Log(args ...interface{})

// Logf formats its arguments according to the format, analogous to Printf, and
// records the text in the error log. The text will be printed only if the test
// fails or the -test.v flag is set.
func (*T) Logf(format string, args ...interface{})

// Parallel signals that this test is to be run in parallel with (and only with)
// other parallel tests.
func (*T) Parallel()

// Skip is equivalent to Log followed by SkipNow.
func (*T) Skip(args ...interface{})

// SkipNow marks the test as having been skipped and stops its execution.
// Execution will continue at the next test or benchmark. See also FailNow.
// SkipNow must be called from the goroutine running the test, not from other
// goroutines created during the test. Calling SkipNow does not stop those other
// goroutines.
func (*T) SkipNow()

// Skipf is equivalent to Logf followed by SkipNow.
func (*T) Skipf(format string, args ...interface{})

// Skipped reports whether the test was skipped.
func (*T) Skipped() bool

func (BenchmarkResult) AllocedBytesPerOp() int64

func (BenchmarkResult) AllocsPerOp() int64

func (BenchmarkResult) MemString() string

func (BenchmarkResult) NsPerOp() int64

func (BenchmarkResult) String() string

