// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package bytes implements functions for the manipulation of byte slices. It is
// analogous to the facilities of the strings package.

// Package bytes implements functions for
// the manipulation of byte slices. It is
// analogous to the facilities of the
// strings package.
package bytes

// MinRead is the minimum slice size passed to a Read call by Buffer.ReadFrom. As
// long as the Buffer has at least MinRead bytes beyond what is required to hold
// the contents of r, ReadFrom will not grow the underlying buffer.

// MinRead is the minimum slice size passed
// to a Read call by Buffer.ReadFrom. As
// long as the Buffer has at least MinRead
// bytes beyond what is required to hold
// the contents of r, ReadFrom will not
// grow the underlying buffer.
const MinRead = 512

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a
// buffer.

// ErrTooLarge is passed to panic if memory
// cannot be allocated to store data in a
// buffer.
var ErrTooLarge = errors.New("bytes.Buffer: too large")

// Compare returns an integer comparing two byte slices lexicographically. The
// result will be 0 if a==b, -1 if a < b, and +1 if a > b. A nil argument is
// equivalent to an empty slice.

// Compare returns an integer comparing two
// byte slices lexicographically. The
// result will be 0 if a==b, -1 if a < b,
// and +1 if a > b. A nil argument is
// equivalent to an empty slice.
func Compare(a, b []byte) int

// Contains reports whether subslice is within b.

// Contains reports whether subslice is
// within b.
func Contains(b, subslice []byte) bool

// Count counts the number of non-overlapping instances of sep in s.

// Count counts the number of
// non-overlapping instances of sep in s.
func Count(s, sep []byte) int

// Equal returns a boolean reporting whether a and b are the same length and
// contain the same bytes. A nil argument is equivalent to an empty slice.

// Equal returns a boolean reporting
// whether a and b are the same length and
// contain the same bytes. A nil argument
// is equivalent to an empty slice.
func Equal(a, b []byte) bool

// EqualFold reports whether s and t, interpreted as UTF-8 strings, are equal under
// Unicode case-folding.

// EqualFold reports whether s and t,
// interpreted as UTF-8 strings, are equal
// under Unicode case-folding.
func EqualFold(s, t []byte) bool

// Fields splits the slice s around each instance of one or more consecutive white
// space characters, returning a slice of subslices of s or an empty list if s
// contains only white space.

// Fields splits the slice s around each
// instance of one or more consecutive
// white space characters, returning a
// slice of subslices of s or an empty list
// if s contains only white space.
func Fields(s []byte) [][]byte

// FieldsFunc interprets s as a sequence of UTF-8-encoded Unicode code points. It
// splits the slice s at each run of code points c satisfying f(c) and returns a
// slice of subslices of s. If all code points in s satisfy f(c), or len(s) == 0,
// an empty slice is returned. FieldsFunc makes no guarantees about the order in
// which it calls f(c). If f does not return consistent results for a given c,
// FieldsFunc may crash.

// FieldsFunc interprets s as a sequence of
// UTF-8-encoded Unicode code points. It
// splits the slice s at each run of code
// points c satisfying f(c) and returns a
// slice of subslices of s. If all code
// points in s satisfy f(c), or len(s) ==
// 0, an empty slice is returned.
// FieldsFunc makes no guarantees about the
// order in which it calls f(c). If f does
// not return consistent results for a
// given c, FieldsFunc may crash.
func FieldsFunc(s []byte, f func(rune) bool) [][]byte

// HasPrefix tests whether the byte slice s begins with prefix.

// HasPrefix tests whether the byte slice s
// begins with prefix.
func HasPrefix(s, prefix []byte) bool

// HasSuffix tests whether the byte slice s ends with suffix.

// HasSuffix tests whether the byte slice s
// ends with suffix.
func HasSuffix(s, suffix []byte) bool

// Index returns the index of the first instance of sep in s, or -1 if sep is not
// present in s.

// Index returns the index of the first
// instance of sep in s, or -1 if sep is
// not present in s.
func Index(s, sep []byte) int

// IndexAny interprets s as a sequence of UTF-8-encoded Unicode code points. It
// returns the byte index of the first occurrence in s of any of the Unicode code
// points in chars. It returns -1 if chars is empty or if there is no code point in
// common.

// IndexAny interprets s as a sequence of
// UTF-8-encoded Unicode code points. It
// returns the byte index of the first
// occurrence in s of any of the Unicode
// code points in chars. It returns -1 if
// chars is empty or if there is no code
// point in common.
func IndexAny(s []byte, chars string) int

// IndexByte returns the index of the first instance of c in s, or -1 if c is not
// present in s.

// IndexByte returns the index of the first
// instance of c in s, or -1 if c is not
// present in s.
func IndexByte(s []byte, c byte) int

// IndexFunc interprets s as a sequence of UTF-8-encoded Unicode code points. It
// returns the byte index in s of the first Unicode code point satisfying f(c), or
// -1 if none do.

// IndexFunc interprets s as a sequence of
// UTF-8-encoded Unicode code points. It
// returns the byte index in s of the first
// Unicode code point satisfying f(c), or
// -1 if none do.
func IndexFunc(s []byte, f func(r rune) bool) int

// IndexRune interprets s as a sequence of UTF-8-encoded Unicode code points. It
// returns the byte index of the first occurrence in s of the given rune. It
// returns -1 if rune is not present in s.

// IndexRune interprets s as a sequence of
// UTF-8-encoded Unicode code points. It
// returns the byte index of the first
// occurrence in s of the given rune. It
// returns -1 if rune is not present in s.
func IndexRune(s []byte, r rune) int

// Join concatenates the elements of s to create a new byte slice. The separator
// sep is placed between elements in the resulting slice.

// Join concatenates the elements of s to
// create a new byte slice. The separator
// sep is placed between elements in the
// resulting slice.
func Join(s [][]byte, sep []byte) []byte

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is
// not present in s.

// LastIndex returns the index of the last
// instance of sep in s, or -1 if sep is
// not present in s.
func LastIndex(s, sep []byte) int

// LastIndexAny interprets s as a sequence of UTF-8-encoded Unicode code points. It
// returns the byte index of the last occurrence in s of any of the Unicode code
// points in chars. It returns -1 if chars is empty or if there is no code point in
// common.

// LastIndexAny interprets s as a sequence
// of UTF-8-encoded Unicode code points. It
// returns the byte index of the last
// occurrence in s of any of the Unicode
// code points in chars. It returns -1 if
// chars is empty or if there is no code
// point in common.
func LastIndexAny(s []byte, chars string) int

// LastIndexFunc interprets s as a sequence of UTF-8-encoded Unicode code points.
// It returns the byte index in s of the last Unicode code point satisfying f(c),
// or -1 if none do.

// LastIndexFunc interprets s as a sequence
// of UTF-8-encoded Unicode code points. It
// returns the byte index in s of the last
// Unicode code point satisfying f(c), or
// -1 if none do.
func LastIndexFunc(s []byte, f func(r rune) bool) int

// Map returns a copy of the byte slice s with all its characters modified
// according to the mapping function. If mapping returns a negative value, the
// character is dropped from the string with no replacement. The characters in s
// and the output are interpreted as UTF-8-encoded Unicode code points.

// Map returns a copy of the byte slice s
// with all its characters modified
// according to the mapping function. If
// mapping returns a negative value, the
// character is dropped from the string
// with no replacement. The characters in s
// and the output are interpreted as
// UTF-8-encoded Unicode code points.
func Map(mapping func(r rune) rune, s []byte) []byte

// Repeat returns a new byte slice consisting of count copies of b.

// Repeat returns a new byte slice
// consisting of count copies of b.
func Repeat(b []byte, count int) []byte

// Replace returns a copy of the slice s with the first n non-overlapping instances
// of old replaced by new. If old is empty, it matches at the beginning of the
// slice and after each UTF-8 sequence, yielding up to k+1 replacements for a
// k-rune slice. If n < 0, there is no limit on the number of replacements.

// Replace returns a copy of the slice s
// with the first n non-overlapping
// instances of old replaced by new. If old
// is empty, it matches at the beginning of
// the slice and after each UTF-8 sequence,
// yielding up to k+1 replacements for a
// k-rune slice. If n < 0, there is no
// limit on the number of replacements.
func Replace(s, old, new []byte, n int) []byte

// Runes returns a slice of runes (Unicode code points) equivalent to s.

// Runes returns a slice of runes (Unicode
// code points) equivalent to s.
func Runes(s []byte) []rune

// Split slices s into all subslices separated by sep and returns a slice of the
// subslices between those separators. If sep is empty, Split splits after each
// UTF-8 sequence. It is equivalent to SplitN with a count of -1.

// Split slices s into all subslices
// separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, Split splits after each
// UTF-8 sequence. It is equivalent to
// SplitN with a count of -1.
func Split(s, sep []byte) [][]byte

// SplitAfter slices s into all subslices after each instance of sep and returns a
// slice of those subslices. If sep is empty, SplitAfter splits after each UTF-8
// sequence. It is equivalent to SplitAfterN with a count of -1.

// SplitAfter slices s into all subslices
// after each instance of sep and returns a
// slice of those subslices. If sep is
// empty, SplitAfter splits after each
// UTF-8 sequence. It is equivalent to
// SplitAfterN with a count of -1.
func SplitAfter(s, sep []byte) [][]byte

// SplitAfterN slices s into subslices after each instance of sep and returns a
// slice of those subslices. If sep is empty, SplitAfterN splits after each UTF-8
// sequence. The count determines the number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices

// SplitAfterN slices s into subslices
// after each instance of sep and returns a
// slice of those subslices. If sep is
// empty, SplitAfterN splits after each
// UTF-8 sequence. The count determines the
// number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices
func SplitAfterN(s, sep []byte, n int) [][]byte

// SplitN slices s into subslices separated by sep and returns a slice of the
// subslices between those separators. If sep is empty, SplitN splits after each
// UTF-8 sequence. The count determines the number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices

// SplitN slices s into subslices separated
// by sep and returns a slice of the
// subslices between those separators. If
// sep is empty, SplitN splits after each
// UTF-8 sequence. The count determines the
// number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices
func SplitN(s, sep []byte, n int) [][]byte

// Title returns a copy of s with all Unicode letters that begin words mapped to
// their title case.
//
// BUG: The rule Title uses for word boundaries does not handle Unicode punctuation
// properly.

// Title returns a copy of s with all
// Unicode letters that begin words mapped
// to their title case.
//
// BUG: The rule Title uses for word
// boundaries does not handle Unicode
// punctuation properly.
func Title(s []byte) []byte

// ToLower returns a copy of the byte slice s with all Unicode letters mapped to
// their lower case.

// ToLower returns a copy of the byte slice
// s with all Unicode letters mapped to
// their lower case.
func ToLower(s []byte) []byte

// ToLowerSpecial returns a copy of the byte slice s with all Unicode letters
// mapped to their lower case, giving priority to the special casing rules.

// ToLowerSpecial returns a copy of the
// byte slice s with all Unicode letters
// mapped to their lower case, giving
// priority to the special casing rules.
func ToLowerSpecial(_case unicode.SpecialCase, s []byte) []byte

// ToTitle returns a copy of the byte slice s with all Unicode letters mapped to
// their title case.

// ToTitle returns a copy of the byte slice
// s with all Unicode letters mapped to
// their title case.
func ToTitle(s []byte) []byte

// ToTitleSpecial returns a copy of the byte slice s with all Unicode letters
// mapped to their title case, giving priority to the special casing rules.

// ToTitleSpecial returns a copy of the
// byte slice s with all Unicode letters
// mapped to their title case, giving
// priority to the special casing rules.
func ToTitleSpecial(_case unicode.SpecialCase, s []byte) []byte

// ToUpper returns a copy of the byte slice s with all Unicode letters mapped to
// their upper case.

// ToUpper returns a copy of the byte slice
// s with all Unicode letters mapped to
// their upper case.
func ToUpper(s []byte) []byte

// ToUpperSpecial returns a copy of the byte slice s with all Unicode letters
// mapped to their upper case, giving priority to the special casing rules.

// ToUpperSpecial returns a copy of the
// byte slice s with all Unicode letters
// mapped to their upper case, giving
// priority to the special casing rules.
func ToUpperSpecial(_case unicode.SpecialCase, s []byte) []byte

// Trim returns a subslice of s by slicing off all leading and trailing
// UTF-8-encoded Unicode code points contained in cutset.

// Trim returns a subslice of s by slicing
// off all leading and trailing
// UTF-8-encoded Unicode code points
// contained in cutset.
func Trim(s []byte, cutset string) []byte

// TrimFunc returns a subslice of s by slicing off all leading and trailing
// UTF-8-encoded Unicode code points c that satisfy f(c).

// TrimFunc returns a subslice of s by
// slicing off all leading and trailing
// UTF-8-encoded Unicode code points c that
// satisfy f(c).
func TrimFunc(s []byte, f func(r rune) bool) []byte

// TrimLeft returns a subslice of s by slicing off all leading UTF-8-encoded
// Unicode code points contained in cutset.

// TrimLeft returns a subslice of s by
// slicing off all leading UTF-8-encoded
// Unicode code points contained in cutset.
func TrimLeft(s []byte, cutset string) []byte

// TrimLeftFunc returns a subslice of s by slicing off all leading UTF-8-encoded
// Unicode code points c that satisfy f(c).

// TrimLeftFunc returns a subslice of s by
// slicing off all leading UTF-8-encoded
// Unicode code points c that satisfy f(c).
func TrimLeftFunc(s []byte, f func(r rune) bool) []byte

// TrimPrefix returns s without the provided leading prefix string. If s doesn't
// start with prefix, s is returned unchanged.

// TrimPrefix returns s without the
// provided leading prefix string. If s
// doesn't start with prefix, s is returned
// unchanged.
func TrimPrefix(s, prefix []byte) []byte

// TrimRight returns a subslice of s by slicing off all trailing UTF-8-encoded
// Unicode code points that are contained in cutset.

// TrimRight returns a subslice of s by
// slicing off all trailing UTF-8-encoded
// Unicode code points that are contained
// in cutset.
func TrimRight(s []byte, cutset string) []byte

// TrimRightFunc returns a subslice of s by slicing off all trailing UTF-8 encoded
// Unicode code points c that satisfy f(c).

// TrimRightFunc returns a subslice of s by
// slicing off all trailing UTF-8 encoded
// Unicode code points c that satisfy f(c).
func TrimRightFunc(s []byte, f func(r rune) bool) []byte

// TrimSpace returns a subslice of s by slicing off all leading and trailing white
// space, as defined by Unicode.

// TrimSpace returns a subslice of s by
// slicing off all leading and trailing
// white space, as defined by Unicode.
func TrimSpace(s []byte) []byte

// TrimSuffix returns s without the provided trailing suffix string. If s doesn't
// end with suffix, s is returned unchanged.

// TrimSuffix returns s without the
// provided trailing suffix string. If s
// doesn't end with suffix, s is returned
// unchanged.
func TrimSuffix(s, suffix []byte) []byte

// A Buffer is a variable-sized buffer of bytes with Read and Write methods. The
// zero value for Buffer is an empty buffer ready to use.

// A Buffer is a variable-sized buffer of
// bytes with Read and Write methods. The
// zero value for Buffer is an empty buffer
// ready to use.
type Buffer struct {
	// contains filtered or unexported fields
}

// NewBuffer creates and initializes a new Buffer using buf as its initial
// contents. It is intended to prepare a Buffer to read existing data. It can also
// be used to size the internal buffer for writing. To do that, buf should have the
// desired capacity but a length of zero.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is sufficient
// to initialize a Buffer.

// NewBuffer creates and initializes a new
// Buffer using buf as its initial
// contents. It is intended to prepare a
// Buffer to read existing data. It can
// also be used to size the internal buffer
// for writing. To do that, buf should have
// the desired capacity but a length of
// zero.
//
// In most cases, new(Buffer) (or just
// declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBuffer(buf []byte) *Buffer

// NewBufferString creates and initializes a new Buffer using string s as its
// initial contents. It is intended to prepare a buffer to read an existing string.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is sufficient
// to initialize a Buffer.

// NewBufferString creates and initializes
// a new Buffer using string s as its
// initial contents. It is intended to
// prepare a buffer to read an existing
// string.
//
// In most cases, new(Buffer) (or just
// declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBufferString(s string) *Buffer

// Bytes returns a slice of the contents of the unread portion of the buffer;
// len(b.Bytes()) == b.Len(). If the caller changes the contents of the returned
// slice, the contents of the buffer will change provided there are no intervening
// method calls on the Buffer.

// Bytes returns a slice of the contents of
// the unread portion of the buffer;
// len(b.Bytes()) == b.Len(). If the caller
// changes the contents of the returned
// slice, the contents of the buffer will
// change provided there are no intervening
// method calls on the Buffer.
func (b *Buffer) Bytes() []byte

// Grow grows the buffer's capacity, if necessary, to guarantee space for another n
// bytes. After Grow(n), at least n bytes can be written to the buffer without
// another allocation. If n is negative, Grow will panic. If the buffer can't grow
// it will panic with ErrTooLarge.

// Grow grows the buffer's capacity, if
// necessary, to guarantee space for
// another n bytes. After Grow(n), at least
// n bytes can be written to the buffer
// without another allocation. If n is
// negative, Grow will panic. If the buffer
// can't grow it will panic with
// ErrTooLarge.
func (b *Buffer) Grow(n int)

// Len returns the number of bytes of the unread portion of the buffer; b.Len() ==
// len(b.Bytes()).

// Len returns the number of bytes of the
// unread portion of the buffer; b.Len() ==
// len(b.Bytes()).
func (b *Buffer) Len() int

// Next returns a slice containing the next n bytes from the buffer, advancing the
// buffer as if the bytes had been returned by Read. If there are fewer than n
// bytes in the buffer, Next returns the entire buffer. The slice is only valid
// until the next call to a read or write method.

// Next returns a slice containing the next
// n bytes from the buffer, advancing the
// buffer as if the bytes had been returned
// by Read. If there are fewer than n bytes
// in the buffer, Next returns the entire
// buffer. The slice is only valid until
// the next call to a read or write method.
func (b *Buffer) Next(n int) []byte

// Read reads the next len(p) bytes from the buffer or until the buffer is drained.
// The return value n is the number of bytes read. If the buffer has no data to
// return, err is io.EOF (unless len(p) is zero); otherwise it is nil.

// Read reads the next len(p) bytes from
// the buffer or until the buffer is
// drained. The return value n is the
// number of bytes read. If the buffer has
// no data to return, err is io.EOF (unless
// len(p) is zero); otherwise it is nil.
func (b *Buffer) Read(p []byte) (n int, err error)

// ReadByte reads and returns the next byte from the buffer. If no byte is
// available, it returns error io.EOF.

// ReadByte reads and returns the next byte
// from the buffer. If no byte is
// available, it returns error io.EOF.
func (b *Buffer) ReadByte() (c byte, err error)

// ReadBytes reads until the first occurrence of delim in the input, returning a
// slice containing the data up to and including the delimiter. If ReadBytes
// encounters an error before finding a delimiter, it returns the data read before
// the error and the error itself (often io.EOF). ReadBytes returns err != nil if
// and only if the returned data does not end in delim.

// ReadBytes reads until the first
// occurrence of delim in the input,
// returning a slice containing the data up
// to and including the delimiter. If
// ReadBytes encounters an error before
// finding a delimiter, it returns the data
// read before the error and the error
// itself (often io.EOF). ReadBytes returns
// err != nil if and only if the returned
// data does not end in delim.
func (b *Buffer) ReadBytes(delim byte) (line []byte, err error)

// ReadFrom reads data from r until EOF and appends it to the buffer, growing the
// buffer as needed. The return value n is the number of bytes read. Any error
// except io.EOF encountered during the read is also returned. If the buffer
// becomes too large, ReadFrom will panic with ErrTooLarge.

// ReadFrom reads data from r until EOF and
// appends it to the buffer, growing the
// buffer as needed. The return value n is
// the number of bytes read. Any error
// except io.EOF encountered during the
// read is also returned. If the buffer
// becomes too large, ReadFrom will panic
// with ErrTooLarge.
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error)

// ReadRune reads and returns the next UTF-8-encoded Unicode code point from the
// buffer. If no bytes are available, the error returned is io.EOF. If the bytes
// are an erroneous UTF-8 encoding, it consumes one byte and returns U+FFFD, 1.

// ReadRune reads and returns the next
// UTF-8-encoded Unicode code point from
// the buffer. If no bytes are available,
// the error returned is io.EOF. If the
// bytes are an erroneous UTF-8 encoding,
// it consumes one byte and returns U+FFFD,
// 1.
func (b *Buffer) ReadRune() (r rune, size int, err error)

// ReadString reads until the first occurrence of delim in the input, returning a
// string containing the data up to and including the delimiter. If ReadString
// encounters an error before finding a delimiter, it returns the data read before
// the error and the error itself (often io.EOF). ReadString returns err != nil if
// and only if the returned data does not end in delim.

// ReadString reads until the first
// occurrence of delim in the input,
// returning a string containing the data
// up to and including the delimiter. If
// ReadString encounters an error before
// finding a delimiter, it returns the data
// read before the error and the error
// itself (often io.EOF). ReadString
// returns err != nil if and only if the
// returned data does not end in delim.
func (b *Buffer) ReadString(delim byte) (line string, err error)

// Reset resets the buffer so it has no content. b.Reset() is the same as
// b.Truncate(0).

// Reset resets the buffer so it has no
// content. b.Reset() is the same as
// b.Truncate(0).
func (b *Buffer) Reset()

// String returns the contents of the unread portion of the buffer as a string. If
// the Buffer is a nil pointer, it returns "<nil>".

// String returns the contents of the
// unread portion of the buffer as a
// string. If the Buffer is a nil pointer,
// it returns "<nil>".
func (b *Buffer) String() string

// Truncate discards all but the first n unread bytes from the buffer. It panics if
// n is negative or greater than the length of the buffer.

// Truncate discards all but the first n
// unread bytes from the buffer. It panics
// if n is negative or greater than the
// length of the buffer.
func (b *Buffer) Truncate(n int)

// UnreadByte unreads the last byte returned by the most recent read operation. If
// write has happened since the last read, UnreadByte returns an error.

// UnreadByte unreads the last byte
// returned by the most recent read
// operation. If write has happened since
// the last read, UnreadByte returns an
// error.
func (b *Buffer) UnreadByte() error

// UnreadRune unreads the last rune returned by ReadRune. If the most recent read
// or write operation on the buffer was not a ReadRune, UnreadRune returns an
// error. (In this regard it is stricter than UnreadByte, which will unread the
// last byte from any read operation.)

// UnreadRune unreads the last rune
// returned by ReadRune. If the most recent
// read or write operation on the buffer
// was not a ReadRune, UnreadRune returns
// an error. (In this regard it is stricter
// than UnreadByte, which will unread the
// last byte from any read operation.)
func (b *Buffer) UnreadRune() error

// Write appends the contents of p to the buffer, growing the buffer as needed. The
// return value n is the length of p; err is always nil. If the buffer becomes too
// large, Write will panic with ErrTooLarge.

// Write appends the contents of p to the
// buffer, growing the buffer as needed.
// The return value n is the length of p;
// err is always nil. If the buffer becomes
// too large, Write will panic with
// ErrTooLarge.
func (b *Buffer) Write(p []byte) (n int, err error)

// WriteByte appends the byte c to the buffer, growing the buffer as needed. The
// returned error is always nil, but is included to match bufio.Writer's WriteByte.
// If the buffer becomes too large, WriteByte will panic with ErrTooLarge.

// WriteByte appends the byte c to the
// buffer, growing the buffer as needed.
// The returned error is always nil, but is
// included to match bufio.Writer's
// WriteByte. If the buffer becomes too
// large, WriteByte will panic with
// ErrTooLarge.
func (b *Buffer) WriteByte(c byte) error

// WriteRune appends the UTF-8 encoding of Unicode code point r to the buffer,
// returning its length and an error, which is always nil but is included to match
// bufio.Writer's WriteRune. The buffer is grown as needed; if it becomes too
// large, WriteRune will panic with ErrTooLarge.

// WriteRune appends the UTF-8 encoding of
// Unicode code point r to the buffer,
// returning its length and an error, which
// is always nil but is included to match
// bufio.Writer's WriteRune. The buffer is
// grown as needed; if it becomes too
// large, WriteRune will panic with
// ErrTooLarge.
func (b *Buffer) WriteRune(r rune) (n int, err error)

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the buffer
// becomes too large, WriteString will panic with ErrTooLarge.

// WriteString appends the contents of s to
// the buffer, growing the buffer as
// needed. The return value n is the length
// of s; err is always nil. If the buffer
// becomes too large, WriteString will
// panic with ErrTooLarge.
func (b *Buffer) WriteString(s string) (n int, err error)

// WriteTo writes data to w until the buffer is drained or an error occurs. The
// return value n is the number of bytes written; it always fits into an int, but
// it is int64 to match the io.WriterTo interface. Any error encountered during the
// write is also returned.

// WriteTo writes data to w until the
// buffer is drained or an error occurs.
// The return value n is the number of
// bytes written; it always fits into an
// int, but it is int64 to match the
// io.WriterTo interface. Any error
// encountered during the write is also
// returned.
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error)

// A Reader implements the io.Reader, io.ReaderAt, io.WriterTo, io.Seeker,
// io.ByteScanner, and io.RuneScanner interfaces by reading from a byte slice.
// Unlike a Buffer, a Reader is read-only and supports seeking.

// A Reader implements the io.Reader,
// io.ReaderAt, io.WriterTo, io.Seeker,
// io.ByteScanner, and io.RuneScanner
// interfaces by reading from a byte slice.
// Unlike a Buffer, a Reader is read-only
// and supports seeking.
type Reader struct {
	// contains filtered or unexported fields
}

// NewReader returns a new Reader reading from b.

// NewReader returns a new Reader reading
// from b.
func NewReader(b []byte) *Reader

// Len returns the number of bytes of the unread portion of the slice.

// Len returns the number of bytes of the
// unread portion of the slice.
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
