// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package strings implements simple functions to manipulate strings.

// Package strings implements simple
// functions to manipulate strings.
package strings

// Contains returns true if substr is within s.

// Contains returns true if substr is
// within s.
func Contains(s, substr string) bool

// ContainsAny returns true if any Unicode code points in chars are within s.

// ContainsAny returns true if any Unicode
// code points in chars are within s.
func ContainsAny(s, chars string) bool

// ContainsRune returns true if the Unicode code point r is within s.

// ContainsRune returns true if the Unicode
// code point r is within s.
func ContainsRune(s string, r rune) bool

// Count counts the number of non-overlapping instances of sep in s.

// Count counts the number of
// non-overlapping instances of sep in s.
func Count(s, sep string) int

// EqualFold reports whether s and t, interpreted as UTF-8 strings, are equal under
// Unicode case-folding.

// EqualFold reports whether s and t,
// interpreted as UTF-8 strings, are equal
// under Unicode case-folding.
func EqualFold(s, t string) bool

// Fields splits the string s around each instance of one or more consecutive white
// space characters, as defined by unicode.IsSpace, returning an array of
// substrings of s or an empty list if s contains only white space.

// Fields splits the string s around each
// instance of one or more consecutive
// white space characters, as defined by
// unicode.IsSpace, returning an array of
// substrings of s or an empty list if s
// contains only white space.
func Fields(s string) []string

// FieldsFunc splits the string s at each run of Unicode code points c satisfying
// f(c) and returns an array of slices of s. If all code points in s satisfy f(c)
// or the string is empty, an empty slice is returned. FieldsFunc makes no
// guarantees about the order in which it calls f(c). If f does not return
// consistent results for a given c, FieldsFunc may crash.

// FieldsFunc splits the string s at each
// run of Unicode code points c satisfying
// f(c) and returns an array of slices of
// s. If all code points in s satisfy f(c)
// or the string is empty, an empty slice
// is returned. FieldsFunc makes no
// guarantees about the order in which it
// calls f(c). If f does not return
// consistent results for a given c,
// FieldsFunc may crash.
func FieldsFunc(s string, f func(rune) bool) []string

// HasPrefix tests whether the string s begins with prefix.

// HasPrefix tests whether the string s
// begins with prefix.
func HasPrefix(s, prefix string) bool

// HasSuffix tests whether the string s ends with suffix.

// HasSuffix tests whether the string s
// ends with suffix.
func HasSuffix(s, suffix string) bool

// Index returns the index of the first instance of sep in s, or -1 if sep is not
// present in s.

// Index returns the index of the first
// instance of sep in s, or -1 if sep is
// not present in s.
func Index(s, sep string) int

// IndexAny returns the index of the first instance of any Unicode code point from
// chars in s, or -1 if no Unicode code point from chars is present in s.

// IndexAny returns the index of the first
// instance of any Unicode code point from
// chars in s, or -1 if no Unicode code
// point from chars is present in s.
func IndexAny(s, chars string) int

// IndexByte returns the index of the first instance of c in s, or -1 if c is not
// present in s.

// IndexByte returns the index of the first
// instance of c in s, or -1 if c is not
// present in s.
func IndexByte(s string, c byte) int

// IndexFunc returns the index into s of the first Unicode code point satisfying
// f(c), or -1 if none do.

// IndexFunc returns the index into s of
// the first Unicode code point satisfying
// f(c), or -1 if none do.
func IndexFunc(s string, f func(rune) bool) int

// IndexRune returns the index of the first instance of the Unicode code point r,
// or -1 if rune is not present in s.

// IndexRune returns the index of the first
// instance of the Unicode code point r, or
// -1 if rune is not present in s.
func IndexRune(s string, r rune) int

// Join concatenates the elements of a to create a single string. The separator
// string sep is placed between elements in the resulting string.

// Join concatenates the elements of a to
// create a single string. The separator
// string sep is placed between elements in
// the resulting string.
func Join(a []string, sep string) string

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is
// not present in s.

// LastIndex returns the index of the last
// instance of sep in s, or -1 if sep is
// not present in s.
func LastIndex(s, sep string) int

// LastIndexAny returns the index of the last instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.

// LastIndexAny returns the index of the
// last instance of any Unicode code point
// from chars in s, or -1 if no Unicode
// code point from chars is present in s.
func LastIndexAny(s, chars string) int

// LastIndexFunc returns the index into s of the last Unicode code point satisfying
// f(c), or -1 if none do.

// LastIndexFunc returns the index into s
// of the last Unicode code point
// satisfying f(c), or -1 if none do.
func LastIndexFunc(s string, f func(rune) bool) int

// Map returns a copy of the string s with all its characters modified according to
// the mapping function. If mapping returns a negative value, the character is
// dropped from the string with no replacement.

// Map returns a copy of the string s with
// all its characters modified according to
// the mapping function. If mapping returns
// a negative value, the character is
// dropped from the string with no
// replacement.
func Map(mapping func(rune) rune, s string) string

// Repeat returns a new string consisting of count copies of the string s.

// Repeat returns a new string consisting
// of count copies of the string s.
func Repeat(s string, count int) string

// Replace returns a copy of the string s with the first n non-overlapping
// instances of old replaced by new. If old is empty, it matches at the beginning
// of the string and after each UTF-8 sequence, yielding up to k+1 replacements for
// a k-rune string. If n < 0, there is no limit on the number of replacements.

// Replace returns a copy of the string s
// with the first n non-overlapping
// instances of old replaced by new. If old
// is empty, it matches at the beginning of
// the string and after each UTF-8
// sequence, yielding up to k+1
// replacements for a k-rune string. If n <
// 0, there is no limit on the number of
// replacements.
func Replace(s, old, new string, n int) string

// Split slices s into all substrings separated by sep and returns a slice of the
// substrings between those separators. If sep is empty, Split splits after each
// UTF-8 sequence. It is equivalent to SplitN with a count of -1.

// Split slices s into all substrings
// separated by sep and returns a slice of
// the substrings between those separators.
// If sep is empty, Split splits after each
// UTF-8 sequence. It is equivalent to
// SplitN with a count of -1.
func Split(s, sep string) []string

// SplitAfter slices s into all substrings after each instance of sep and returns a
// slice of those substrings. If sep is empty, SplitAfter splits after each UTF-8
// sequence. It is equivalent to SplitAfterN with a count of -1.

// SplitAfter slices s into all substrings
// after each instance of sep and returns a
// slice of those substrings. If sep is
// empty, SplitAfter splits after each
// UTF-8 sequence. It is equivalent to
// SplitAfterN with a count of -1.
func SplitAfter(s, sep string) []string

// SplitAfterN slices s into substrings after each instance of sep and returns a
// slice of those substrings. If sep is empty, SplitAfterN splits after each UTF-8
// sequence. The count determines the number of substrings to return:
//
//	n > 0: at most n substrings; the last substring will be the unsplit remainder.
//	n == 0: the result is nil (zero substrings)
//	n < 0: all substrings

// SplitAfterN slices s into substrings
// after each instance of sep and returns a
// slice of those substrings. If sep is
// empty, SplitAfterN splits after each
// UTF-8 sequence. The count determines the
// number of substrings to return:
//
//	n > 0: at most n substrings; the last substring will be the unsplit remainder.
//	n == 0: the result is nil (zero substrings)
//	n < 0: all substrings
func SplitAfterN(s, sep string, n int) []string

// SplitN slices s into substrings separated by sep and returns a slice of the
// substrings between those separators. If sep is empty, SplitN splits after each
// UTF-8 sequence. The count determines the number of substrings to return:
//
//	n > 0: at most n substrings; the last substring will be the unsplit remainder.
//	n == 0: the result is nil (zero substrings)
//	n < 0: all substrings

// SplitN slices s into substrings
// separated by sep and returns a slice of
// the substrings between those separators.
// If sep is empty, SplitN splits after
// each UTF-8 sequence. The count
// determines the number of substrings to
// return:
//
//	n > 0: at most n substrings; the last substring will be the unsplit remainder.
//	n == 0: the result is nil (zero substrings)
//	n < 0: all substrings
func SplitN(s, sep string, n int) []string

// Title returns a copy of the string s with all Unicode letters that begin words
// mapped to their title case.
//
// BUG: The rule Title uses for word boundaries does not handle Unicode punctuation
// properly.

// Title returns a copy of the string s
// with all Unicode letters that begin
// words mapped to their title case.
//
// BUG: The rule Title uses for word
// boundaries does not handle Unicode
// punctuation properly.
func Title(s string) string

// ToLower returns a copy of the string s with all Unicode letters mapped to their
// lower case.

// ToLower returns a copy of the string s
// with all Unicode letters mapped to their
// lower case.
func ToLower(s string) string

// ToLowerSpecial returns a copy of the string s with all Unicode letters mapped to
// their lower case, giving priority to the special casing rules.

// ToLowerSpecial returns a copy of the
// string s with all Unicode letters mapped
// to their lower case, giving priority to
// the special casing rules.
func ToLowerSpecial(_case unicode.SpecialCase, s string) string

// ToTitle returns a copy of the string s with all Unicode letters mapped to their
// title case.

// ToTitle returns a copy of the string s
// with all Unicode letters mapped to their
// title case.
func ToTitle(s string) string

// ToTitleSpecial returns a copy of the string s with all Unicode letters mapped to
// their title case, giving priority to the special casing rules.

// ToTitleSpecial returns a copy of the
// string s with all Unicode letters mapped
// to their title case, giving priority to
// the special casing rules.
func ToTitleSpecial(_case unicode.SpecialCase, s string) string

// ToUpper returns a copy of the string s with all Unicode letters mapped to their
// upper case.

// ToUpper returns a copy of the string s
// with all Unicode letters mapped to their
// upper case.
func ToUpper(s string) string

// ToUpperSpecial returns a copy of the string s with all Unicode letters mapped to
// their upper case, giving priority to the special casing rules.

// ToUpperSpecial returns a copy of the
// string s with all Unicode letters mapped
// to their upper case, giving priority to
// the special casing rules.
func ToUpperSpecial(_case unicode.SpecialCase, s string) string

// Trim returns a slice of the string s with all leading and trailing Unicode code
// points contained in cutset removed.

// Trim returns a slice of the string s
// with all leading and trailing Unicode
// code points contained in cutset removed.
func Trim(s string, cutset string) string

// TrimFunc returns a slice of the string s with all leading and trailing Unicode
// code points c satisfying f(c) removed.

// TrimFunc returns a slice of the string s
// with all leading and trailing Unicode
// code points c satisfying f(c) removed.
func TrimFunc(s string, f func(rune) bool) string

// TrimLeft returns a slice of the string s with all leading Unicode code points
// contained in cutset removed.

// TrimLeft returns a slice of the string s
// with all leading Unicode code points
// contained in cutset removed.
func TrimLeft(s string, cutset string) string

// TrimLeftFunc returns a slice of the string s with all leading Unicode code
// points c satisfying f(c) removed.

// TrimLeftFunc returns a slice of the
// string s with all leading Unicode code
// points c satisfying f(c) removed.
func TrimLeftFunc(s string, f func(rune) bool) string

// TrimPrefix returns s without the provided leading prefix string. If s doesn't
// start with prefix, s is returned unchanged.

// TrimPrefix returns s without the
// provided leading prefix string. If s
// doesn't start with prefix, s is returned
// unchanged.
func TrimPrefix(s, prefix string) string

// TrimRight returns a slice of the string s, with all trailing Unicode code points
// contained in cutset removed.

// TrimRight returns a slice of the string
// s, with all trailing Unicode code points
// contained in cutset removed.
func TrimRight(s string, cutset string) string

// TrimRightFunc returns a slice of the string s with all trailing Unicode code
// points c satisfying f(c) removed.

// TrimRightFunc returns a slice of the
// string s with all trailing Unicode code
// points c satisfying f(c) removed.
func TrimRightFunc(s string, f func(rune) bool) string

// TrimSpace returns a slice of the string s, with all leading and trailing white
// space removed, as defined by Unicode.

// TrimSpace returns a slice of the string
// s, with all leading and trailing white
// space removed, as defined by Unicode.
func TrimSpace(s string) string

// TrimSuffix returns s without the provided trailing suffix string. If s doesn't
// end with suffix, s is returned unchanged.

// TrimSuffix returns s without the
// provided trailing suffix string. If s
// doesn't end with suffix, s is returned
// unchanged.
func TrimSuffix(s, suffix string) string

// A Reader implements the io.Reader, io.ReaderAt, io.Seeker, io.WriterTo,
// io.ByteScanner, and io.RuneScanner interfaces by reading from a string.

// A Reader implements the io.Reader,
// io.ReaderAt, io.Seeker, io.WriterTo,
// io.ByteScanner, and io.RuneScanner
// interfaces by reading from a string.
type Reader struct {
	// contains filtered or unexported fields
}

// NewReader returns a new Reader reading from s. It is similar to
// bytes.NewBufferString but more efficient and read-only.

// NewReader returns a new Reader reading
// from s. It is similar to
// bytes.NewBufferString but more efficient
// and read-only.
func NewReader(s string) *Reader

// Len returns the number of bytes of the unread portion of the string.

// Len returns the number of bytes of the
// unread portion of the string.
func (r *Reader) Len() int

func (r *Reader) Read(b []byte) (n int, err error)

func (r *Reader) ReadAt(b []byte, off int64) (n int, err error)

func (r *Reader) ReadByte() (b byte, err error)

func (r *Reader) ReadRune() (ch rune, size int, err error)

// Seek implements the io.Seeker interface.

// Seek implements the io.Seeker interface.
func (r *Reader) Seek(offset int64, whence int) (int64, error)

func (r *Reader) UnreadByte() error

func (r *Reader) UnreadRune() error

// WriteTo implements the io.WriterTo interface.

// WriteTo implements the io.WriterTo
// interface.
func (r *Reader) WriteTo(w io.Writer) (n int64, err error)

// Replacer replaces a list of strings with replacements. It is safe for concurrent
// use by multiple goroutines.

// Replacer replaces a list of strings with
// replacements. It is safe for concurrent
// use by multiple goroutines.
type Replacer struct {
	// contains filtered or unexported fields
}

// NewReplacer returns a new Replacer from a list of old, new string pairs.
// Replacements are performed in order, without overlapping matches.

// NewReplacer returns a new Replacer from
// a list of old, new string pairs.
// Replacements are performed in order,
// without overlapping matches.
func NewReplacer(oldnew ...string) *Replacer

// Replace returns a copy of s with all replacements performed.

// Replace returns a copy of s with all
// replacements performed.
func (r *Replacer) Replace(s string) string

// WriteString writes s to w with all replacements performed.

// WriteString writes s to w with all
// replacements performed.
func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)
