// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package ioutil implements some I/O utility functions.

// ioutil 实现了一些I/O的工具函数。
package ioutil

import (
    "bytes"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "sync"
    "time"
)

// Discard is an io.Writer on which all Write calls succeed
// without doing anything.

// Discard 是一个 io.Writer，对它进行的任何 Write 调用都将无条件成功。
var Discard io.Writer = devNull(0)

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.

// NopCloser 将提供的 Reader r 用空操作 Close 方法包装后作为 ReadCloser 返回。
func NopCloser(r io.Reader) io.ReadCloser

// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.

// ReadAll 从 r 中读取，直至遇到错误或EOF，然后返回它所读取的数据。 一次成功的调
// 用应当返回 err == nil，而非 err == 因为 ReadAll 被定义为从 src 进行读取直至遇
// 到EOF，它并不会将来自 Read 的EOF视作错误来报告。
func ReadAll(r io.Reader) ([]byte, error)

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.

// ReadDir 读取名为 dirname
// 的目录并返回一个已排序的目录项列表。
func ReadDir(dirname string) ([]os.FileInfo, error)

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.

// ReadFile 读取名为 filename 的文件并返回其内容。 一次成功的调用应当返回 err ==
// nil，而非 err == EOF。因为 ReadFile 会读取整个文件， 它并不会将来自 Read 的
// EOF视作错误来报告。
func ReadFile(filename string) ([]byte, error)

// TempDir creates a new temporary directory in the directory dir
// with a name beginning with prefix and returns the path of the
// new directory.  If dir is the empty string, TempDir uses the
// default directory for temporary files (see os.TempDir).
// Multiple programs calling TempDir simultaneously
// will not choose the same directory.  It is the caller's responsibility
// to remove the directory when no longer needed.

// TempDir 在目录 dir 中创建一个名字以 prefix 开头的新的临时目录并返回该新目录的
// 路径。 若 dir 为空字符串，TempDir 就会为临时文件（Unix将目录也视作文件）使用
// 默认的目录（见 os.TempDir）。多程序同时调用 TempDir 将不会选择相同的目录。当
// 该目录不再被需要时， 调用者应负责将其移除。
func TempDir(dir, prefix string) (name string, err error)

// TempFile creates a new temporary file in the directory dir
// with a name beginning with prefix, opens the file for reading
// and writing, and returns the resulting *os.File.
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file.  The caller can use f.Name()
// to find the pathname of the file.  It is the caller's responsibility
// to remove the file when no longer needed.

// TempFile 在目录 dir 中创建一个名字以 prefix 开头的新的临时文件，打开该文件以
// 用于读写， 并返回其结果 *os.File。若 dir 为空字符串，TempFile 就会为临时文件
// 使用默认的目录（见 os.TempDir）。多程序同时调用 TempFile 将不会选择相同的文件
// 。调用者可使用 f.Name() 来查找该文件的路径名 pathname。当该文件不再被需要时，
// 调用者应负责将其移除。
func TempFile(dir, prefix string) (f *os.File, err error)

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it with permissions perm;
// otherwise WriteFile truncates it before writing.

// WriteFile 将数据写入到名为 filename 的文件中。 若该文件不存在，WriteFile 就会
// 按照权限 perm 创建它；否则 WriteFile 就会在写入前将其截断。
func WriteFile(filename string, data []byte, perm os.FileMode) error

