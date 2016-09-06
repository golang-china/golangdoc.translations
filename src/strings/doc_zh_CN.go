// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package strings implements simple functions to manipulate UTF-8 encoded
// strings.
//
// For information about UTF-8 strings in Go, see
// https://blog.golang.org/strings.

// strings包实现了用于操作字符的简单函数。
package strings

import (
	"errors"
	"io"
	"unicode"
	"unicode/utf8"
)

// A Reader implements the io.Reader, io.ReaderAt, io.Seeker, io.WriterTo,
// io.ByteScanner, and io.RuneScanner interfaces by reading
// from a string.

// Reader类型通过从一个字符串读取数据，实现了io.Reader、io.Seeker、io.ReaderAt、
// io.WriterTo、io.ByteScanner、io.RuneScanner接口。
type Reader struct {
}

// Replacer replaces a list of strings with replacements.
// It is safe for concurrent use by multiple goroutines.

// Replacer类型进行一系列字符串的替换。
type Replacer struct {
}

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
//
// Compare is included only for symmetry with package bytes.
// It is usually clearer and always faster to use the built-in
// string comparison operators ==, <, >, and so on.
func Compare(a, b string) int

// Contains reports whether substr is within s.

// 判断字符串s是否包含子串substr。
func Contains(s, substr string) bool

// ContainsAny reports whether any Unicode code points in chars are within s.

// 判断字符串s是否包含字符串chars中的任一字符。
func ContainsAny(s, chars string) bool

// ContainsRune reports whether the Unicode code point r is within s.

// 判断字符串s是否包含utf-8码值r。
func ContainsRune(s string, r rune) bool

// Count counts the number of non-overlapping instances of sep in s. If sep is
// an empty string, Count returns 1 + the number of Unicode code points in s.

// 返回字符串s中有几个不重复的sep子串。
func Count(s, sep string) int

// EqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding.

// 判断两个utf-8编码字符串（将unicode大写、小写、标题三种格式字符视为相同）是否
// 相同。
func EqualFold(s, t string) bool

// Fields splits the string s around each instance of one or more consecutive
// white space characters, as defined by unicode.IsSpace, returning an array of
// substrings of s or an empty list if s contains only white space.

// 返回将字符串按照空白（unicode.IsSpace确定，可以是一到多个连续的空白字符）分割
// 的多个字符串。如果字符串全部是空白或者是空字符串的话，会返回空切片。
func Fields(s string) []string

// FieldsFunc splits the string s at each run of Unicode code points c
// satisfying f(c) and returns an array of slices of s. If all code points in s
// satisfy f(c) or the string is empty, an empty slice is returned. FieldsFunc
// makes no guarantees about the order in which it calls f(c). If f does not
// return consistent results for a given c, FieldsFunc may crash.

// 类似Fields，但使用函数f来确定分割符（满足f的unicode码值）。如果字符串全部是分
// 隔符或者是空字符串的话，会返回空切片。
func FieldsFunc(s string, f func(rune) bool) []string

// HasPrefix tests whether the string s begins with prefix.

// 判断s是否有前缀字符串prefix。
func HasPrefix(s, prefix string) bool

// HasSuffix tests whether the string s ends with suffix.

// 判断s是否有后缀字符串suffix。
func HasSuffix(s, suffix string) bool

// Index returns the index of the first instance of sep in s, or -1 if sep is
// not present in s.

// 子串sep在字符串s中第一次出现的位置，不存在则返回-1。
func Index(s, sep string) int

// IndexAny returns the index of the first instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.

// 字符串chars中的任一utf-8码值在s中第一次出现的位置，如果不存在或者chars为空字
// 符串则返回-1。
func IndexAny(s, chars string) int

// IndexByte returns the index of the first instance of c in s, or -1 if c is
// not present in s.

// 字符c在s中第一次出现的位置，不存在则返回-1。
func IndexByte(s string, c byte) int

// IndexFunc returns the index into s of the first Unicode
// code point satisfying f(c), or -1 if none do.

// s中第一个满足函数f的位置i（该处的utf-8码值r满足f(r)==true），不存在则返回-1。
func IndexFunc(s string, f func(rune) bool) int

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in s.

// unicode码值r在s中第一次出现的位置，不存在则返回-1。
func IndexRune(s string, r rune) int

// Join concatenates the elements of a to create a single string. The separator
// string sep is placed between elements in the resulting string.

// 将一系列字符串连接为一个字符串，之间用sep来分隔。
func Join(a []string, sep string) string

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is
// not present in s.

// 子串sep在字符串s中最后一次出现的位置，不存在则返回-1。
func LastIndex(s, sep string) int

// LastIndexAny returns the index of the last instance of any Unicode code
// point from chars in s, or -1 if no Unicode code point from chars is
// present in s.

// 字符串chars中的任一utf-8码值在s中最后一次出现的位置，如不存在或者chars为空字
// 符串则返回-1。
func LastIndexAny(s, chars string) int

// LastIndexByte returns the index of the last instance of c in s, or -1 if c is
// not present in s.
func LastIndexByte(s string, c byte) int

// LastIndexFunc returns the index into s of the last
// Unicode code point satisfying f(c), or -1 if none do.

// s中最后一个满足函数f的unicode码值的位置i，不存在则返回-1。
func LastIndexFunc(s string, f func(rune) bool) int

// Map returns a copy of the string s with all its characters modified according
// to the mapping function. If mapping returns a negative value, the character
// is dropped from the string with no replacement.

// 将s的每一个unicode码值r都替换为mapping(r)，返回这些新码值组成的字符串拷贝。如
// 果mapping返回一个负值，将会丢弃该码值而不会被替换。（返回值中对应位置将没有码
// 值）
func Map(mapping func(rune) rune, s string) string

// NewReader returns a new Reader reading from s.
// It is similar to bytes.NewBufferString but more efficient and read-only.

// NewReader创建一个从s读取数据的Reader。本函数类似bytes.NewBufferString，但是更
// 有效率，且为只读的。
func NewReader(s string) *Reader

// NewReplacer returns a new Replacer from a list of old, new string pairs.
// Replacements are performed in order, without overlapping matches.

// 使用提供的多组old、new字符串对创建并返回一个*Replacer。替换是依次进行的，匹配
// 时不会重叠。
func NewReplacer(oldnew ...string) *Replacer

// Repeat returns a new string consisting of count copies of the string s.

// 返回count个s串联的字符串。
func Repeat(s string, count int) string

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune string.
// If n < 0, there is no limit on the number of replacements.

// 返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串。
func Replace(s, old, new string, n int) string

// Split slices s into all substrings separated by sep and returns a slice of
// the substrings between those separators.
// If sep is empty, Split splits after each UTF-8 sequence.
// It is equivalent to SplitN with a count of -1.

// 用去掉s中出现的sep的方式进行分割，会分割到结尾，并返回生成的所有片段组成的切
// 片（每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割）。如果sep为
// 空字符，Split会将s切分成每一个unicode码值一个字符串。
func Split(s, sep string) []string

// SplitAfter slices s into all substrings after each instance of sep and
// returns a slice of those substrings.
// If sep is empty, SplitAfter splits after each UTF-8 sequence.
// It is equivalent to SplitAfterN with a count of -1.

// 用从s中出现的sep后面切断的方式进行分割，会分割到结尾，并返回生成的所有片段组
// 成的切片（每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割）。如
// 果sep为空字符，Split会将s切分成每一个unicode码值一个字符串。
func SplitAfter(s, sep string) []string

// SplitAfterN slices s into substrings after each instance of sep and returns a
// slice of those substrings. If sep is empty, SplitAfterN splits after each
// UTF-8 sequence. The count determines the number of substrings to return:
//
// 	n > 0: at most n substrings; the last substring will be the unsplit remainder.
// 	n == 0: the result is nil (zero substrings)
// 	n < 0: all substrings

// 用从s中出现的sep后面切断的方式进行分割，会分割到结尾，并返回生成的所有片段组
// 成的切片（每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割）。如
// 果sep为空字符，Split会将s切分成每一个unicode码值一个字符串。参数n决定返回的切
// 片的数目：
//
// 	n > 0 : 返回的切片最多n个子字符串；最后一个子字符串包含未进行切割的部分。
// 	n == 0: 返回nil
// 	n < 0 : 返回所有的子字符串组成的切
func SplitAfterN(s, sep string, n int) []string

// SplitN slices s into substrings separated by sep and returns a slice of the
// substrings between those separators. If sep is empty, SplitN splits after
// each UTF-8 sequence. The count determines the number of substrings to return:
//
// 	n > 0: at most n substrings; the last substring will be the unsplit remainder.
// 	n == 0: the result is nil (zero substrings)
// 	n < 0: all substrings

// 用去掉s中出现的sep的方式进行分割，会分割到结尾，并返回生成的所有片段组成的切
// 片（每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割）。如果sep为
// 空字符，Split会将s切分成每一个unicode码值一个字符串。参数n决定返回的切片的数
// 目：
//
// 	n > 0 : 返回的切片最多n个子字符串；最后一个子字符串包含未进行切割的部分。
// 	n == 0: 返回nil
// 	n < 0 : 返回所有的子字符串组成的切片
func SplitN(s, sep string, n int) []string

// Title returns a copy of the string s with all Unicode letters that begin
// words mapped to their title case.
//
// BUG(rsc): The rule Title uses for word boundaries does not handle Unicode
// punctuation properly.

// 返回s中每个单词的首字母都改为标题格式的字符串拷贝。
//
// BUG:
// Title用于划分单词的规则不能很好的处理Unicode标点符号。
func Title(s string) string

// ToLower returns a copy of the string s with all Unicode letters mapped to
// their lower case.

// 返回将所有字母都转为对应的小写版本的拷贝。
func ToLower(s string) string

// ToLowerSpecial returns a copy of the string s with all Unicode letters mapped
// to their lower case, giving priority to the special casing rules.

// 使用_case规定的字符映射，返回将所有字母都转为对应的小写版本的拷贝。
func ToLowerSpecial(_case unicode.SpecialCase, s string) string

// ToTitle returns a copy of the string s with all Unicode letters mapped to
// their title case.

// 返回将所有字母都转为对应的标题版本的拷贝。
func ToTitle(s string) string

// ToTitleSpecial returns a copy of the string s with all Unicode letters mapped
// to their title case, giving priority to the special casing rules.

// 使用_case规定的字符映射，返回将所有字母都转为对应的标题版本的拷贝。
func ToTitleSpecial(_case unicode.SpecialCase, s string) string

// ToUpper returns a copy of the string s with all Unicode letters mapped to
// their upper case.

// 返回将所有字母都转为对应的大写版本的拷贝。
func ToUpper(s string) string

// ToUpperSpecial returns a copy of the string s with all Unicode letters mapped
// to their upper case, giving priority to the special casing rules.

// 使用_case规定的字符映射，返回将所有字母都转为对应的大写版本的拷贝。
func ToUpperSpecial(_case unicode.SpecialCase, s string) string

// Trim returns a slice of the string s with all leading and
// trailing Unicode code points contained in cutset removed.

// 返回将s前后端所有cutset包含的utf-8码值都去掉的字符串。
func Trim(s string, cutset string) string

// TrimFunc returns a slice of the string s with all leading
// and trailing Unicode code points c satisfying f(c) removed.

// 返回将s前后端所有满足f的unicode码值都去掉的字符串。
func TrimFunc(s string, f func(rune) bool) string

// TrimLeft returns a slice of the string s with all leading
// Unicode code points contained in cutset removed.

// 返回将s前端所有cutset包含的utf-8码值都去掉的字符串。
func TrimLeft(s string, cutset string) string

// TrimLeftFunc returns a slice of the string s with all leading
// Unicode code points c satisfying f(c) removed.

// 返回将s前端所有满足f的unicode码值都去掉的字符串。
func TrimLeftFunc(s string, f func(rune) bool) string

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.

// 返回去除s可能的前缀prefix的字符串。
func TrimPrefix(s, prefix string) string

// TrimRight returns a slice of the string s, with all trailing
// Unicode code points contained in cutset removed.

// 返回将s后端所有cutset包含的utf-8码值都去掉的字符串。
func TrimRight(s string, cutset string) string

// TrimRightFunc returns a slice of the string s with all trailing
// Unicode code points c satisfying f(c) removed.

// 返回将s后端所有满足f的unicode码值都去掉的字符串。
func TrimRightFunc(s string, f func(rune) bool) string

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.

// 返回将s前后端所有空白（unicode.IsSpace指定）都去掉的字符串。
func TrimSpace(s string) string

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.

// 返回去除s可能的后缀suffix的字符串。
func TrimSuffix(s, suffix string) string

// Len returns the number of bytes of the unread portion of the
// string.

// Len返回r包含的字符串还没有被读取的部分。
func (r *Reader) Len() int

func (r *Reader) Read(b []byte) (n int, err error)

func (r *Reader) ReadAt(b []byte, off int64) (n int, err error)

func (r *Reader) ReadByte() (byte, error)

func (r *Reader) ReadRune() (ch rune, size int, err error)

// Reset resets the Reader to be reading from s.
func (r *Reader) Reset(s string)

// Seek implements the io.Seeker interface.

// Seek实现了io.Seeker接口。
func (r *Reader) Seek(offset int64, whence int) (int64, error)

// Size returns the original length of the underlying string.
// Size is the number of bytes available for reading via ReadAt.
// The returned value is always the same and is not affected by calls
// to any other method.
func (r *Reader) Size() int64

func (r *Reader) UnreadByte() error

func (r *Reader) UnreadRune() error

// WriteTo implements the io.WriterTo interface.

// WriteTo实现了io.WriterTo接口。
func (r *Reader) WriteTo(w io.Writer) (n int64, err error)

// Replace returns a copy of s with all replacements performed.
func (r *Replacer) Replace(s string) string

// WriteString writes s to w with all replacements performed.
func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)

