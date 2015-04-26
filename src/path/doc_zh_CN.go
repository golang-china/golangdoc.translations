// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package path implements utility routines for manipulating slash-separated paths.

// path实现了对斜杠分隔的路径的实用操作函数。
package path

// ErrBadPattern indicates a globbing pattern was malformed.

// ErrBadPattern表示一个glob模式匹配字符串的格式错误。
var ErrBadPattern = errors.New("syntax error in pattern")

// Base returns the last element of path. Trailing slashes are removed before
// extracting the last element. If the path is empty, Base returns ".". If the path
// consists entirely of slashes, Base returns "/".

// Base函数返回路径的最后一个元素。在提取元素前会求掉末尾的斜杠。如果路径是""，会返回"."；如果路径是只有一个斜杆构成，会返回"/"。
func Base(path string) string

// Clean returns the shortest path name equivalent to path by purely lexical
// processing. It applies the following rules iteratively until no further
// processing can be done:
//
//	1. Replace multiple slashes with a single slash.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path.
//
// The returned path ends in a slash only if it is the root "/".
//
// If the result of this process is an empty string, Clean returns the string ".".
//
// See also Rob Pike, ``Lexical File Names in Plan 9 or Getting Dot-Dot Right,''
// http://plan9.bell-labs.com/sys/doc/lexnames.html

// Clean函数通过单纯的词法操作返回和path代表同一地址的最短路径。
//
// 它会不断的依次应用如下的规则，直到不能再进行任何处理：
//
//	1. 将连续的多个斜杠替换为单个斜杠
//	2. 剔除每一个.路径名元素（代表当前目录）
//	3. 剔除每一个路径内的..路径名元素（代表父目录）和它前面的非..路径名元素
//	4. 剔除开始一个根路径的..路径名元素，即将路径开始处的"/.."替换为"/"
//
// 只有路径代表根地址"/"时才会以斜杠结尾。如果处理的结果是空字符串，Clean会返回"."。
//
// 参见http://plan9.bell-labs.com/sys/doc/lexnames.html
func Clean(path string) string

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element using Split, the path is Cleaned and trailing
// slashes are removed. If the path is empty, Dir returns ".". If the path consists
// entirely of slashes followed by non-slash bytes, Dir returns a single slash. In
// any other case, the returned path does not end in a slash.

// Dir返回路径除去最后一个路径元素的部分，即该路径最后一个元素所在的目录。在使用Split去掉最后一个元素后，会简化路径并去掉末尾的斜杠。如果路径是空字符串，会返回"."；如果路径由1到多个斜杠后跟0到多个非斜杠字符组成，会返回"/"；其他任何情况下都不会返回以斜杠结尾的路径。
func Dir(path string) string

// Ext returns the file name extension used by path. The extension is the suffix
// beginning at the final dot in the final slash-separated element of path; it is
// empty if there is no dot.

// Ext函数返回path文件扩展名。返回值是路径最后一个斜杠分隔出的路径元素的最后一个'.'起始的后缀（包括'.'）。如果该元素没有'.'会返回空字符串。
func Ext(path string) string

// IsAbs returns true if the path is absolute.

// IsAbs返回路径是否是一个绝对路径。
func IsAbs(path string) bool

// Join joins any number of path elements into a single path, adding a separating
// slash if necessary. The result is Cleaned; in particular, all empty strings are
// ignored.

// Join函数可以将任意数量的路径元素放入一个单一路径里，会根据需要添加斜杠。结果是经过简化的，所有的空字符串元素会被忽略。
func Join(elem ...string) string

// Match returns true if name matches the shell file name pattern. The pattern
// syntax is:
//
//	pattern:
//		{ term }
//	term:
//		'*'         matches any sequence of non-/ characters
//		'?'         matches any single non-/ character
//		'[' [ '^' ] { character-range } ']'
//		            character class (must be non-empty)
//		c           matches character c (c != '*', '?', '\\', '[')
//		'\\' c      matches character c
//
//	character-range:
//		c           matches character c (c != '\\', '-', ']')
//		'\\' c      matches character c
//		lo '-' hi   matches character c for lo <= c <= hi
//
// Match requires pattern to match all of name, not just a substring. The only
// possible returned error is ErrBadPattern, when pattern is malformed.

// 如果name匹配shell文件名模式匹配字符串，Match函数返回真。该模式匹配字符串语法为：
//
//	pattern:
//		{ term }
//	term:
//		'*'                                  匹配0或多个非/的字符
//		'?'                                  匹配1个非/的字符
//		'[' [ '^' ] { character-range } ']'  字符组（必须非空）
//		c                                    匹配字符c（c != '*', '?', '\\', '['）
//		'\\' c                               匹配字符c
//	character-range:
//		c           匹配字符c（c != '\\', '-', ']'）
//		'\\' c      匹配字符c
//		lo '-' hi   匹配区间[lo, hi]内的字符
//
// Match要求匹配整个name字符串，而不是它的一部分。只有pattern语法错误时，会返回ErrBadPattern。
func Match(pattern, name string) (matched bool, err error)

// Split splits path immediately following the final slash. separating it into a
// directory and file name component. If there is no slash path, Split returns an
// empty dir and file set to path. The returned values have the property that path
// = dir+file.

// Split函数将路径从最后一个斜杠后面位置分隔为两个部分（dir和file）并返回。如果路径中没有斜杠，函数返回值dir会设为空字符串，file会设为path。两个返回值满足path
// == dir+file。
func Split(path string) (dir, file string)
