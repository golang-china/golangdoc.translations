// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package exec runs external commands. It wraps os.StartProcess to make it easier
// to remap stdin and stdout, connect I/O with pipes, and do other adjustments.

// Package exec runs external commands. It
// wraps os.StartProcess to make it easier
// to remap stdin and stdout, connect I/O
// with pipes, and do other adjustments.
package exec

// ErrNotFound is the error resulting if a path search failed to find an executable
// file.

// ErrNotFound is the error resulting if a
// path search failed to find an executable
// file.
var ErrNotFound = errors.New("executable file not found in $path")

// ErrNotFound is the error resulting if a path search failed to find an executable
// file.

// ErrNotFound is the error resulting if a
// path search failed to find an executable
// file.
var ErrNotFound = errors.New("executable file not found in $PATH")

// ErrNotFound is the error resulting if a path search failed to find an executable
// file.

// ErrNotFound is the error resulting if a
// path search failed to find an executable
// file.
var ErrNotFound = errors.New("executable file not found in %PATH%")

// LookPath searches for an executable binary named file in the directories named
// by the path environment variable. If file begins with "/", "#", "./", or "../",
// it is tried directly and the path is not consulted. The result may be an
// absolute path or a path relative to the current directory.

// LookPath searches for an executable
// binary named file in the directories
// named by the path environment variable.
// If file begins with "/", "#", "./", or
// "../", it is tried directly and the path
// is not consulted. The result may be an
// absolute path or a path relative to the
// current directory.
func LookPath(file string) (string, error)

// Cmd represents an external command being prepared or run.

// Cmd represents an external command being
// prepared or run.
type Cmd struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// Env specifies the environment of the process.
	// If Env is nil, Run uses the current process's environment.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// Stdin specifies the process's standard input.
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error) or because writing to the pipe returned an error.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If Stdout and Stderr are the same writer, at most one
	// goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	// ExtraFiles specifies additional open files to be inherited by the
	// new process. It does not include standard input, standard output, or
	// standard error. If non-nil, entry i becomes file descriptor 3+i.
	//
	// BUG: on OS X 10.6, child processes may sometimes inherit unwanted fds.
	// http://golang.org/issue/2603
	ExtraFiles []*os.File

	// SysProcAttr holds optional, operating system-specific attributes.
	// Run passes it to os.StartProcess as the os.ProcAttr's Sys field.
	SysProcAttr *syscall.SysProcAttr

	// Process is the underlying process, once started.
	Process *os.Process

	// ProcessState contains information about an exited process,
	// available after a call to Wait or Run.
	ProcessState *os.ProcessState
	// contains filtered or unexported fields
}

// Command returns the Cmd struct to execute the named program with the given
// arguments.
//
// It sets only the Path and Args in the returned structure.
//
// If name contains no path separators, Command uses LookPath to resolve the path
// to a complete name if possible. Otherwise it uses name directly.
//
// The returned Cmd's Args field is constructed from the command name followed by
// the elements of arg, so arg should not include the command name itself. For
// example, Command("echo", "hello")

// Command returns the Cmd struct to
// execute the named program with the given
// arguments.
//
// It sets only the Path and Args in the
// returned structure.
//
// If name contains no path separators,
// Command uses LookPath to resolve the
// path to a complete name if possible.
// Otherwise it uses name directly.
//
// The returned Cmd's Args field is
// constructed from the command name
// followed by the elements of arg, so arg
// should not include the command name
// itself. For example, Command("echo",
// "hello")
func Command(name string, arg ...string) *Cmd

// CombinedOutput runs the command and returns its combined standard output and
// standard error.

// CombinedOutput runs the command and
// returns its combined standard output and
// standard error.
func (c *Cmd) CombinedOutput() ([]byte, error)

// Output runs the command and returns its standard output.

// Output runs the command and returns its
// standard output.
func (c *Cmd) Output() ([]byte, error)

// Run starts the specified command and waits for it to complete.
//
// The returned error is nil if the command runs, has no problems copying stdin,
// stdout, and stderr, and exits with a zero exit status.
//
// If the command fails to run or doesn't complete successfully, the error is of
// type *ExitError. Other error types may be returned for I/O problems.

// Run starts the specified command and
// waits for it to complete.
//
// The returned error is nil if the command
// runs, has no problems copying stdin,
// stdout, and stderr, and exits with a
// zero exit status.
//
// If the command fails to run or doesn't
// complete successfully, the error is of
// type *ExitError. Other error types may
// be returned for I/O problems.
func (c *Cmd) Run() error

// Start starts the specified command but does not wait for it to complete.
//
// The Wait method will return the exit code and release associated resources once
// the command exits.

// Start starts the specified command but
// does not wait for it to complete.
//
// The Wait method will return the exit
// code and release associated resources
// once the command exits.
func (c *Cmd) Start() error

// StderrPipe returns a pipe that will be connected to the command's standard error
// when the command starts.
//
// Wait will close the pipe after seeing the command exit, so most callers need not
// close the pipe themselves; however, an implication is that it is incorrect to
// call Wait before all reads from the pipe have completed. For the same reason, it
// is incorrect to use Run when using StderrPipe. See the StdoutPipe example for
// idiomatic usage.

// StderrPipe returns a pipe that will be
// connected to the command's standard
// error when the command starts.
//
// Wait will close the pipe after seeing
// the command exit, so most callers need
// not close the pipe themselves; however,
// an implication is that it is incorrect
// to call Wait before all reads from the
// pipe have completed. For the same
// reason, it is incorrect to use Run when
// using StderrPipe. See the StdoutPipe
// example for idiomatic usage.
func (c *Cmd) StderrPipe() (io.ReadCloser, error)

// StdinPipe returns a pipe that will be connected to the command's standard input
// when the command starts. The pipe will be closed automatically after Wait sees
// the command exit. A caller need only call Close to force the pipe to close
// sooner. For example, if the command being run will not exit until standard input
// is closed, the caller must close the pipe.

// StdinPipe returns a pipe that will be
// connected to the command's standard
// input when the command starts. The pipe
// will be closed automatically after Wait
// sees the command exit. A caller need
// only call Close to force the pipe to
// close sooner. For example, if the
// command being run will not exit until
// standard input is closed, the caller
// must close the pipe.
func (c *Cmd) StdinPipe() (io.WriteCloser, error)

// StdoutPipe returns a pipe that will be connected to the command's standard
// output when the command starts.
//
// Wait will close the pipe after seeing the command exit, so most callers need not
// close the pipe themselves; however, an implication is that it is incorrect to
// call Wait before all reads from the pipe have completed. For the same reason, it
// is incorrect to call Run when using StdoutPipe. See the example for idiomatic
// usage.

// StdoutPipe returns a pipe that will be
// connected to the command's standard
// output when the command starts.
//
// Wait will close the pipe after seeing
// the command exit, so most callers need
// not close the pipe themselves; however,
// an implication is that it is incorrect
// to call Wait before all reads from the
// pipe have completed. For the same
// reason, it is incorrect to call Run when
// using StdoutPipe. See the example for
// idiomatic usage.
func (c *Cmd) StdoutPipe() (io.ReadCloser, error)

// Wait waits for the command to exit. It must have been started by Start.
//
// The returned error is nil if the command runs, has no problems copying stdin,
// stdout, and stderr, and exits with a zero exit status.
//
// If the command fails to run or doesn't complete successfully, the error is of
// type *ExitError. Other error types may be returned for I/O problems.
//
// Wait releases any resources associated with the Cmd.

// Wait waits for the command to exit. It
// must have been started by Start.
//
// The returned error is nil if the command
// runs, has no problems copying stdin,
// stdout, and stderr, and exits with a
// zero exit status.
//
// If the command fails to run or doesn't
// complete successfully, the error is of
// type *ExitError. Other error types may
// be returned for I/O problems.
//
// Wait releases any resources associated
// with the Cmd.
func (c *Cmd) Wait() error

// Error records the name of a binary that failed to be executed and the reason it
// failed.

// Error records the name of a binary that
// failed to be executed and the reason it
// failed.
type Error struct {
	Name string
	Err  error
}

func (e *Error) Error() string

// An ExitError reports an unsuccessful exit by a command.

// An ExitError reports an unsuccessful
// exit by a command.
type ExitError struct {
	*os.ProcessState
}

func (e *ExitError) Error() string
