// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package strings implements simple functions to manipulate UTF-8 encoded
// strings.
//
// For information about UTF-8 strings in Go, see
// https://blog.golang.org/strings.

// Package strings implements simple functions to manipulate UTF-8 encoded
// strings.
//
// For information about UTF-8 strings in Go, see
// https://blog.golang.org/strings.
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
type Reader struct {
	s        string
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}


// Replacer replaces a list of strings with replacements.
// It is safe for concurrent use by multiple goroutines.
type Replacer struct {
	r replacer
}


// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
//
// Compare is included only for symmetry with package bytes.
// It is usually clearer and always faster to use the built-in
// string comparison operators ==, <, >, and so on.
func Compare(a, b string) int

// Contains reports whether substr is within s.
func Contains(s, substr string) bool

// ContainsAny reports whether any Unicode code points in chars are within s.
func ContainsAny(s, chars string) bool

// ContainsRune reports whether the Unicode code point r is within s.
func ContainsRune(s string, r rune) bool

// Count counts the number of non-overlapping instances of sep in s. If sep is
// an empty string, Count returns 1 + the number of Unicode code points in s.
func Count(s, sep string) int

// EqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding.
func EqualFold(s, t string) bool

// Fields splits the string s around each instance of one or more consecutive
// white space characters, as defined by unicode.IsSpace, returning an array of
// substrings of s or an empty list if s contains only white space.
func Fields(s string) []string

// FieldsFunc splits the string s at each run of Unicode code points c
// satisfying f(c) and returns an array of slices of s. If all code points in s
// satisfy f(c) or the string is empty, an empty slice is returned. FieldsFunc
// makes no guarantees about the order in which it calls f(c). If f does not
// return consistent results for a given c, FieldsFunc may crash.
func FieldsFunc(s string, f func(rune) bool) []string

// HasPrefix tests whether the string s begins with prefix.
func HasPrefix(s, prefix string) bool

// HasSuffix tests whether the string s ends with suffix.
func HasSuffix(s, suffix string) bool

// Index returns the index of the first instance of sep in s, or -1 if sep is
// not present in s.
func Index(s, sep string) int

// IndexAny returns the index of the first instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.
func IndexAny(s, chars string) int

// IndexByte returns the index of the first instance of c in s, or -1 if c is
// not present in s.
func IndexByte(s string, c byte) int

// IndexFunc returns the index into s of the first Unicode
// code point satisfying f(c), or -1 if none do.
func IndexFunc(s string, f func(rune) bool) int

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in s.
func IndexRune(s string, r rune) int

// Join concatenates the elements of a to create a single string. The separator
// string sep is placed between elements in the resulting string.

// Join concatenates the elements of a to create a single string. The separator
// string sep is placed between elements in the resulting string.
func Join(a []string, sep string) string

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is
// not present in s.
func LastIndex(s, sep string) int

// LastIndexAny returns the index of the last instance of any Unicode code
// point from chars in s, or -1 if no Unicode code point from chars is
// present in s.
func LastIndexAny(s, chars string) int

// LastIndexByte returns the index of the last instance of c in s, or -1 if c is
// not present in s.
func LastIndexByte(s string, c byte) int

// LastIndexFunc returns the index into s of the last
// Unicode code point satisfying f(c), or -1 if none do.
func LastIndexFunc(s string, f func(rune) bool) int

// Map returns a copy of the string s with all its characters modified according
// to the mapping function. If mapping returns a negative value, the character
// is dropped from the string with no replacement.
func Map(mapping func(rune) rune, s string) string

// NewReader returns a new Reader reading from s.
// It is similar to bytes.NewBufferString but more efficient and read-only.
func NewReader(s string) *Reader

// NewReplacer returns a new Replacer from a list of old, new string pairs.
// Replacements are performed in order, without overlapping matches.
func NewReplacer(oldnew ...string) *Replacer

// Repeat returns a new string consisting of count copies of the string s.
func Repeat(s string, count int) string

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune string.
// If n < 0, there is no limit on the number of replacements.
func Replace(s, old, new string, n int) string

// Split slices s into all substrings separated by sep and returns a slice of
// the substrings between those separators.
// If sep is empty, Split splits after each UTF-8 sequence.
// It is equivalent to SplitN with a count of -1.
func Split(s, sep string) []string

// SplitAfter slices s into all substrings after each instance of sep and
// returns a slice of those substrings.
// If sep is empty, SplitAfter splits after each UTF-8 sequence.
// It is equivalent to SplitAfterN with a count of -1.
func SplitAfter(s, sep string) []string

// SplitAfterN slices s into substrings after each instance of sep and returns a
// slice of those substrings. If sep is empty, SplitAfterN splits after each
// UTF-8 sequence. The count determines the number of substrings to return:
//
//     n > 0: at most n substrings; the last substring will be the unsplit remainder.
//     n == 0: the result is nil (zero substrings)
//     n < 0: all substrings

// SplitAfterN slices s into substrings after each instance of sep and returns a
// slice of those substrings. If sep is empty, SplitAfterN splits after each
// UTF-8 sequence. The count determines the number of substrings to return:
//
//     n > 0: at most n substrings; the last substring will be the unsplit remainder.
//     n == 0: the result is nil (zero substrings)
//     n < 0: all substrings
func SplitAfterN(s, sep string, n int) []string

// SplitN slices s into substrings separated by sep and returns a slice of the
// substrings between those separators. If sep is empty, SplitN splits after
// each UTF-8 sequence. The count determines the number of substrings to return:
//
//     n > 0: at most n substrings; the last substring will be the unsplit remainder.
//     n == 0: the result is nil (zero substrings)
//     n < 0: all substrings

// SplitN slices s into substrings separated by sep and returns a slice of the
// substrings between those separators. If sep is empty, SplitN splits after
// each UTF-8 sequence. The count determines the number of substrings to return:
//
//     n > 0: at most n substrings; the last substring will be the unsplit remainder.
//     n == 0: the result is nil (zero substrings)
//     n < 0: all substrings
func SplitN(s, sep string, n int) []string

// Title returns a copy of the string s with all Unicode letters that begin
// words mapped to their title case.
//
// BUG(rsc): The rule Title uses for word boundaries does not handle Unicode
// punctuation properly.
func Title(s string) string

// ToLower returns a copy of the string s with all Unicode letters mapped to
// their lower case.
func ToLower(s string) string

// ToLowerSpecial returns a copy of the string s with all Unicode letters mapped
// to their lower case, giving priority to the special casing rules.
func ToLowerSpecial(_case unicode.SpecialCase, s string) string

// ToTitle returns a copy of the string s with all Unicode letters mapped to
// their title case.
func ToTitle(s string) string

// ToTitleSpecial returns a copy of the string s with all Unicode letters mapped
// to their title case, giving priority to the special casing rules.
func ToTitleSpecial(_case unicode.SpecialCase, s string) string

// ToUpper returns a copy of the string s with all Unicode letters mapped to
// their upper case.
func ToUpper(s string) string

// ToUpperSpecial returns a copy of the string s with all Unicode letters mapped
// to their upper case, giving priority to the special casing rules.
func ToUpperSpecial(_case unicode.SpecialCase, s string) string

// Trim returns a slice of the string s with all leading and
// trailing Unicode code points contained in cutset removed.
func Trim(s string, cutset string) string

// TrimFunc returns a slice of the string s with all leading
// and trailing Unicode code points c satisfying f(c) removed.
func TrimFunc(s string, f func(rune) bool) string

// TrimLeft returns a slice of the string s with all leading
// Unicode code points contained in cutset removed.
func TrimLeft(s string, cutset string) string

// TrimLeftFunc returns a slice of the string s with all leading
// Unicode code points c satisfying f(c) removed.
func TrimLeftFunc(s string, f func(rune) bool) string

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func TrimPrefix(s, prefix string) string

// TrimRight returns a slice of the string s, with all trailing
// Unicode code points contained in cutset removed.
func TrimRight(s string, cutset string) string

// TrimRightFunc returns a slice of the string s with all trailing
// Unicode code points c satisfying f(c) removed.
func TrimRightFunc(s string, f func(rune) bool) string

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func TrimSpace(s string) string

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.
func TrimSuffix(s, suffix string) string

// Len returns the number of bytes of the unread portion of the
// string.
func (*Reader) Len() int

func (*Reader) Read(b []byte) (n int, err error)

func (*Reader) ReadAt(b []byte, off int64) (n int, err error)

func (*Reader) ReadByte() (byte, error)

func (*Reader) ReadRune() (ch rune, size int, err error)

// Reset resets the Reader to be reading from s.
func (*Reader) Reset(s string)

// Seek implements the io.Seeker interface.
func (*Reader) Seek(offset int64, whence int) (int64, error)

// Size returns the original length of the underlying string.
// Size is the number of bytes available for reading via ReadAt.
// The returned value is always the same and is not affected by calls
// to any other method.
func (*Reader) Size() int64

func (*Reader) UnreadByte() error

func (*Reader) UnreadRune() error

// WriteTo implements the io.WriterTo interface.
func (*Reader) WriteTo(w io.Writer) (n int64, err error)

// Replace returns a copy of s with all replacements performed.
func (*Replacer) Replace(s string) string

// WriteString writes s to w with all replacements performed.
func (*Replacer) WriteString(w io.Writer, s string) (n int, err error)

