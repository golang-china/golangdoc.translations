// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package exec runs external commands. It wraps os.StartProcess to make it
// easier to remap stdin and stdout, connect I/O with pipes, and do other
// adjustments.
//
// Note that the examples in this package assume a Unix system.
// They may not run on Windows, and they do not run in the Go Playground
// used by golang.org and godoc.org.

// exec包执行外部命令。它包装了os.StartProcess函数以便更容易的修正输入和输出，使
// 用管道连接I/O，以及作其它的一些调整。
package exec

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// ErrNotFound is the error resulting if a path search failed to find an
// executable file.

// 如果路径搜索没有找到可执行文件时，就会返回本错误。
var ErrNotFound = errors.New("executable file not found in $PATH")

// Cmd represents an external command being prepared or run.
//
// A Cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.

// Cmd代表一个正在准备或者在执行中的外部命令。
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

	// Stdin specifies the process's standard input. If Stdin is nil, the
	// process reads from the null device (os.DevNull). If Stdin is an *os.File,
	// the process's standard input is connected directly to that file.
	// Otherwise, during the execution of the command a separate goroutine reads
	// from Stdin and delivers that data to the command over a pipe. In this
	// case, Wait does not complete until the goroutine stops copying, either
	// because it has reached the end of Stdin (EOF or a read error) or because
	// writing to the pipe returned an error.
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

	// ExtraFiles specifies additional open files to be inherited by the new
	// process. It does not include standard input, standard output, or standard
	// error. If non-nil, entry i becomes file descriptor 3+i.
	//
	// BUG(rsc): On OS X 10.6, child processes may sometimes inherit unwanted
	// fds. https://golang.org/issue/2603

	// ExtraFiles specifies additional open files to be inherited by the new
	// process. It does not include standard input, standard output, or standard
	// error. If non-nil, entry i becomes file descriptor 3+i.
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
}

// Error records the name of a binary that failed to be executed
// and the reason it failed.

// Error类型记录执行失败的程序名和失败的原因。
type Error struct {
	Name string
	Err  error
}

// An ExitError reports an unsuccessful exit by a command.

// ExitError报告某个命令的一次未成功的返回。
type ExitError struct {
	*os.ProcessState

	// Stderr holds a subset of the standard error output from the
	// Cmd.Output method if standard error was not otherwise being
	// collected.
	//
	// If the error output is long, Stderr may contain only a prefix
	// and suffix of the output, with the middle replaced with
	// text about the number of omitted bytes.
	//
	// Stderr is provided for debugging, for inclusion in error messages.
	// Users with other needs should redirect Cmd.Stderr as needed.
	Stderr []byte
}

// Command returns the Cmd struct to execute the named program with
// the given arguments.
//
// It sets only the Path and Args in the returned structure.
//
// If name contains no path separators, Command uses LookPath to
// resolve the path to a complete name if possible. Otherwise it uses
// name directly.
//
// The returned Cmd's Args field is constructed from the command name
// followed by the elements of arg, so arg should not include the
// command name itself. For example, Command("echo", "hello")

// 函数返回一个*Cmd，用于使用给出的参数执行name指定的程序。返回值只设定了Path和
// Args两个参数。
//
// 如果name不含路径分隔符，将使用LookPath获取完整路径；否则直接使用name。参数arg
// 不应包含命令名。
func Command(name string, arg ...string) *Cmd

// CommandContext is like Command but includes a context.
//
// The provided context is used to kill the process (by calling
// os.Process.Kill) if the context becomes done before the command
// completes on its own.
func CommandContext(ctx context.Context, name string, arg ...string) *Cmd

// LookPath searches for an executable binary named file in the directories
// named by the PATH environment variable. If file contains a slash, it is tried
// directly and the PATH is not consulted. The result may be an absolute path or
// a path relative to the current directory.

// 在环境变量PATH指定的目录中搜索可执行文件，如file中有斜杠，则只在当前目录搜索
// 。返回完整路径或者相对于当前目录的一个相对路径。
func LookPath(file string) (string, error)

// CombinedOutput runs the command and returns its combined standard
// output and standard error.

// CombinedOutput runs the command and returns its combined standard output and
// standard error.
func (c *Cmd) CombinedOutput() ([]byte, error)

// Output runs the command and returns its standard output.
// Any returned error will usually be of type *ExitError.
// If c.Stderr was nil, Output populates ExitError.Stderr.

// Output runs the command and returns its standard output.
func (c *Cmd) Output() ([]byte, error)

// Run starts the specified command and waits for it to complete.
//
// The returned error is nil if the command runs, has no problems
// copying stdin, stdout, and stderr, and exits with a zero exit
// status.
//
// If the command fails to run or doesn't complete successfully, the
// error is of type *ExitError. Other error types may be
// returned for I/O problems.

// Run starts the specified command and waits for it to complete.
//
// The returned error is nil if the command runs, has no problems copying stdin,
// stdout, and stderr, and exits with a zero exit status.
//
// If the command fails to run or doesn't complete successfully, the error is of
// type *ExitError. Other error types may be returned for I/O problems.
func (c *Cmd) Run() error

// Start starts the specified command but does not wait for it to complete.
//
// The Wait method will return the exit code and release associated resources
// once the command exits.
func (c *Cmd) Start() error

// StderrPipe returns a pipe that will be connected to the command's standard
// error when the command starts.
//
// Wait will close the pipe after seeing the command exit, so most callers need
// not close the pipe themselves; however, an implication is that it is
// incorrect to call Wait before all reads from the pipe have completed. For the
// same reason, it is incorrect to use Run when using StderrPipe. See the
// StdoutPipe example for idiomatic usage.
func (c *Cmd) StderrPipe() (io.ReadCloser, error)

// StdinPipe returns a pipe that will be connected to the command's
// standard input when the command starts.
// The pipe will be closed automatically after Wait sees the command exit.
// A caller need only call Close to force the pipe to close sooner.
// For example, if the command being run will not exit until standard input
// is closed, the caller must close the pipe.

// StdinPipe returns a pipe that will be connected to the command's standard
// input when the command starts. The pipe will be closed automatically after
// Wait sees the command exit. A caller need only call Close to force the pipe
// to close sooner. For example, if the command being run will not exit until
// standard input is closed, the caller must close the pipe.
func (c *Cmd) StdinPipe() (io.WriteCloser, error)

// StdoutPipe returns a pipe that will be connected to the command's standard
// output when the command starts.
//
// Wait will close the pipe after seeing the command exit, so most callers need
// not close the pipe themselves; however, an implication is that it is
// incorrect to call Wait before all reads from the pipe have completed. For the
// same reason, it is incorrect to call Run when using StdoutPipe. See the
// example for idiomatic usage.
func (c *Cmd) StdoutPipe() (io.ReadCloser, error)

// Wait waits for the command to exit.
// It must have been started by Start.
//
// The returned error is nil if the command runs, has no problems
// copying stdin, stdout, and stderr, and exits with a zero exit
// status.
//
// If the command fails to run or doesn't complete successfully, the
// error is of type *ExitError. Other error types may be
// returned for I/O problems.
//
// If c.Stdin is not an *os.File, Wait also waits for the I/O loop
// copying from c.Stdin into the process's standard input
// to complete.
//
// Wait releases any resources associated with the Cmd.

// Wait waits for the command to exit. It must have been started by Start.
//
// The returned error is nil if the command runs, has no problems copying stdin,
// stdout, and stderr, and exits with a zero exit status.
//
// If the command fails to run or doesn't complete successfully, the error is of
// type *ExitError. Other error types may be returned for I/O problems.
//
// Wait releases any resources associated with the Cmd.
func (c *Cmd) Wait() error

func (e *Error) Error() string

func (e *ExitError) Error() string

