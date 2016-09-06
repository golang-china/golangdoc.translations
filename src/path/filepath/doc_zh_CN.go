// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package filepath implements utility routines for manipulating filename paths
// in a way compatible with the target operating system-defined file paths.

// Package filepath implements utility routines for manipulating filename paths
// in a way compatible with the target operating system-defined file paths.
//
// Functions in this package replace any occurrences of the slash ('/')
// character with os.PathSeparator when returning paths unless otherwise
// specified.
package filepath

import (
	"errors"
	"os"
	"runtime"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	Separator     = os.PathSeparator
	ListSeparator = os.PathListSeparator
)

// ErrBadPattern indicates a globbing pattern was malformed.

// ErrBadPattern表示一个glob模式匹配字符串的格式错误。
var ErrBadPattern = errors.New("syntax error in pattern")

// SkipDir is used as a return value from WalkFuncs to indicate that
// the directory named in the call is to be skipped. It is not returned
// as an error by any function.

// 用作WalkFunc类型的返回值，表示该次调用的path参数指定的目录应被跳过。本错误不
// 应被任何其他函数返回。
var SkipDir = errors.New("skip this directory")

// WalkFunc is the type of the function called for each file or directory
// visited by Walk. The path argument contains the argument to Walk as a prefix;
// that is, if Walk is called with "dir", which is a directory containing the
// file "a", the walk function will be called with argument "dir/a". The info
// argument is the os.FileInfo for the named path.
//
// If there was a problem walking to the file or directory named by path, the
// incoming error will describe the problem and the function can decide how to
// handle that error (and Walk will not descend into that directory). If an
// error is returned, processing stops. The sole exception is when the function
// returns the special value SkipDir. If the function returns SkipDir when
// invoked on a directory, Walk skips the directory's contents entirely. If the
// function returns SkipDir when invoked on a non-directory file, Walk skips the
// remaining files in the containing directory.

// Walk函数对每一个文件/目录都会调用WalkFunc函数类型值。调用时path参数会包含Walk
// 的root参数作为前缀；就是说，如果Walk函数的root为"dir"，该目录下有文件"a"，将
// 会使用"dir/a"调用walkFn参数。walkFn参数被调用时的info参数是path指定的地址（文
// 件/目录）的文件信息，类型为os.FileInfo。
//
// 如果遍历path指定的文件或目录时出现了问题，传入的参数err会描述该问题，WalkFunc
// 类型函数可以决定如何去处理该错误（Walk函数将不会深入该目录）；如果该函数返回
// 一个错误，Walk函数的执行会中止；只有一个例外，如果Walk的walkFn返回值是SkipDir
// ，将会跳过该目录的内容而Walk函数照常执行处理下一个文件。
type WalkFunc func(path string, info os.FileInfo, err error) error

// Abs returns an absolute representation of path.
// If the path is not absolute it will be joined with the current
// working directory to turn it into an absolute path. The absolute
// path name for a given file is not guaranteed to be unique.
// Abs calls Clean on the result.

// Abs returns an absolute representation of path.
// If the path is not absolute it will be joined with the current
// working directory to turn it into an absolute path. The absolute
// path name for a given file is not guaranteed to be unique.
func Abs(path string) (string, error)

// Base returns the last element of path. Trailing path separators are removed
// before extracting the last element. If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.

// Base函数返回路径的最后一个元素。在提取元素前会求掉末尾的路径分隔符。如果路径
// 是""，会返回"."；如果路径是只有一个斜杆构成，会返回单个路径分隔符。
func Base(path string) string

// Clean returns the shortest path name equivalent to path
// by purely lexical processing. It applies the following rules
// iteratively until no further processing can be done:
//
// 	1. Replace multiple Separator elements with a single one.
// 	2. Eliminate each . path name element (the current directory).
// 	3. Eliminate each inner .. path name element (the parent directory)
// 	   along with the non-.. element that precedes it.
// 	4. Eliminate .. elements that begin a rooted path:
// 	   that is, replace "/.." by "/" at the beginning of a path,
// 	   assuming Separator is '/'.
//
// The returned path ends in a slash only if it represents a root directory,
// such as "/" on Unix or `C:\` on Windows.
//
// Finally, any occurrences of slash are replaced by Separator.
//
// If the result of this process is an empty string, Clean
// returns the string ".".
//
// See also Rob Pike, ``Lexical File Names in Plan 9 or
// Getting Dot-Dot Right,''
// https://9p.io/sys/doc/lexnames.html

// Clean returns the shortest path name equivalent to path
// by purely lexical processing. It applies the following rules
// iteratively until no further processing can be done:
//
// 	1. Replace multiple Separator elements with a single one.
// 	2. Eliminate each . path name element (the current directory).
// 	3. Eliminate each inner .. path name element (the parent directory)
// 	   along with the non-.. element that precedes it.
// 	4. Eliminate .. elements that begin a rooted path:
// 	   that is, replace "/.." by "/" at the beginning of a path,
// 	   assuming Separator is '/'.
//
// The returned path ends in a slash only if it represents a root directory,
// such as "/" on Unix or `C:\` on Windows.
//
// If the result of this process is an empty string, Clean
// returns the string ".".
//
// See also Rob Pike, ``Lexical File Names in Plan 9 or
// Getting Dot-Dot Right,''
// https://9p.io/sys/doc/lexnames.html
func Clean(path string) string

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element, Dir calls Clean on the path and trailing
// slashes are removed. If the path is empty, Dir returns ".". If the path
// consists entirely of separators, Dir returns a single separator. The returned
// path does not end in a separator unless it is the root directory.

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element, the path is Cleaned and trailing slashes
// are removed. If the path is empty, Dir returns ".". If the path consists
// entirely of separators, Dir returns a single separator. The returned path
// does not end in a separator unless it is the root directory.
func Dir(path string) string

// EvalSymlinks returns the path name after the evaluation of any symbolic
// links.
// If path is relative the result will be relative to the current directory,
// unless one of the components is an absolute symbolic link.
// EvalSymlinks calls Clean on the result.

// EvalSymlinks returns the path name after the evaluation of any symbolic
// links.
// If path is relative the result will be relative to the current directory,
// unless one of the components is an absolute symbolic link.
func EvalSymlinks(path string) (string, error)

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final element of path; it is empty if there is
// no dot.

// Ext函数返回path文件扩展名。返回值是路径最后一个路径元素的最后一个'.'起始的后
// 缀（包括'.'）。如果该元素没有'.'会返回空字符串。
func Ext(path string) string

// FromSlash returns the result of replacing each slash ('/') character
// in path with a separator character. Multiple slashes are replaced
// by multiple separators.

// FromSlash函数将path中的斜杠（'/'）替换为路径分隔符并返回替换结果，多个斜杠会
// 替换为多个路径分隔符。
func FromSlash(path string) string

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.

// Glob函数返回所有匹配模式匹配字符串pattern的文件或者nil（如果没有匹配的文件）
// 。pattern的语法和Match函数相同。pattern可以描述多层的名字，如/usr/*/bin/ed（
// 假设路径分隔符是'/'）。
func Glob(pattern string) (matches []string, err error)

// HasPrefix exists for historical compatibility and should not be used.

// HasPrefix函数出于历史兼容问题保留，不应被使用。
func HasPrefix(p, prefix string) bool

// IsAbs reports whether the path is absolute.

// IsAbs返回路径是否是一个绝对路径。
func IsAbs(path string) bool

// Join joins any number of path elements into a single path, adding
// a Separator if necessary. Join calls Clean on the result; in particular,
// all empty strings are ignored.
// On Windows, the result is a UNC path if and only if the first path
// element is a UNC path.

// Join joins any number of path elements into a single path, adding
// a Separator if necessary. The result is Cleaned, in particular
// all empty strings are ignored.
// On Windows, the result is a UNC path if and only if the first path
// element is a UNC path.
func Join(elem ...string) string

// Match reports whether name matches the shell file name pattern.
// The pattern syntax is:
//
// 	pattern:
// 		{ term }
// 	term:
// 		'*'         matches any sequence of non-Separator characters
// 		'?'         matches any single non-Separator character
// 		'[' [ '^' ] { character-range } ']'
// 		            character class (must be non-empty)
// 		c           matches character c (c != '*', '?', '\\', '[')
// 		'\\' c      matches character c
//
// 	character-range:
// 		c           matches character c (c != '\\', '-', ']')
// 		'\\' c      matches character c
// 		lo '-' hi   matches character c for lo <= c <= hi
//
// Match requires pattern to match all of name, not just a substring.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
//
// On Windows, escaping is disabled. Instead, '\\' is treated as
// path separator.

// Match returns true if name matches the shell file name pattern.
//
// 	    The pattern syntax is:
//
// 	pattern:
// 	    { term }
// 	term:
// 	    '*'                                  匹配0或多个非路径分隔符的字符
// 	    '?'                                  匹配1个非路径分隔符的字符
// 	    '[' [ '^' ] { character-range } ']'  字符组（必须非空）
// 	    c                                    匹配字符c（c != '*', '?', '\\', '['）
// 	    '\\' c                               匹配字符c
// 	character-range:
// 	    c           匹配字符c（c != '\\', '-', ']'）
// 	    '\\' c      匹配字符c
// 	    lo '-' hi   匹配区间[lo, hi]内的字符
//
// Match要求匹配整个name字符串，而不是它的一部分。只有pattern语法错误时，会返回
// ErrBadPattern。
//
// Windows系统中，不能进行转义：'\\'被视为路径分隔符。
func Match(pattern, name string) (matched bool, err error)

// Rel returns a relative path that is lexically equivalent to targpath when
// joined to basepath with an intervening separator. That is,
// Join(basepath, Rel(basepath, targpath)) is equivalent to targpath itself.
// On success, the returned path will always be relative to basepath,
// even if basepath and targpath share no elements.
// An error is returned if targpath can't be made relative to basepath or if
// knowing the current working directory would be necessary to compute it.
// Rel calls Clean on the result.

// Rel returns a relative path that is lexically equivalent to targpath when
// joined to basepath with an intervening separator. That is,
// Join(basepath, Rel(basepath, targpath)) is equivalent to targpath itself.
// On success, the returned path will always be relative to basepath,
// even if basepath and targpath share no elements.
// An error is returned if targpath can't be made relative to basepath or if
// knowing the current working directory would be necessary to compute it.
func Rel(basepath, targpath string) (string, error)

// Split splits path immediately following the final Separator,
// separating it into a directory and file name component.
// If there is no Separator in path, Split returns an empty dir
// and file set to path.
// The returned values have the property that path = dir+file.

// Split函数将路径从最后一个路径分隔符后面位置分隔为两个部分（dir和file）并返回
// 。如果路径中没有路径分隔符，函数返回值dir会设为空字符串，file会设为path。两个
// 返回值满足path == dir+file。
func Split(path string) (dir, file string)

// SplitList splits a list of paths joined by the OS-specific ListSeparator,
// usually found in PATH or GOPATH environment variables. Unlike strings.Split,
// SplitList returns an empty slice when passed an empty string. SplitList does
// not replace slash characters in the returned paths.

// 将PATH或GOPATH等环境变量里的多个路径分割开（这些路径被OS特定的表分隔符连接起
// 来）。与strings.Split函数的不同之处是：对""，SplitList返回[]string{}，而
// strings.Split返回[]string{""}。
func SplitList(path string) []string

// ToSlash returns the result of replacing each separator character
// in path with a slash ('/') character. Multiple separators are
// replaced by multiple slashes.

// ToSlash函数将path中的路径分隔符替换为斜杠（'/'）并返回替换结果，多个路径分隔
// 符会替换为多个斜杠。
func ToSlash(path string) string

// VolumeName returns leading volume name.
// Given "C:\foo\bar" it returns "C:" on Windows.
// Given "\\host\share\foo" it returns "\\host\share".
// On other platforms it returns "".

// VolumeName函数返回最前面的卷名。如Windows系统里提供参数"C:\foo\bar"会返回"C:"
// ；Unix/linux系统的"\\host\share\foo"会返回"\\host\share"；其他平台会返回""。
func VolumeName(path string) string

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order, which makes the output deterministic but means that for very large
// directories Walk can be inefficient. Walk does not follow symbolic links.

// Walk函数会遍历root指定的目录下的文件树，对每一个该文件树中的目录和文件都会调
// 用walkFn，包括root自身。所有访问文件/目录时遇到的错误都会传递给walkFn过滤。文
// 件是按词法顺序遍历的，这让输出更漂亮，但也导致处理非常大的目录时效率会降低。
// Walk函数不会遍历文件树中的符号链接（快捷方式）文件包含的路径。
func Walk(root string, walkFn WalkFunc) error

