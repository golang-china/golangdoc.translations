// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package utf8 implements functions and constants to support text encoded in
// UTF-8. It includes functions to translate between runes and UTF-8 byte
// sequences.

// utf8 包实现了支持UTF-8文本编码的函数和常量.
// 其中包括了在符文和UTF-8字节序列之间进行转译的函数。
package utf8

// Numbers fundamental to the encoding.

// 用于编码的基本数值。
const (
	RuneError = '\uFFFD'     // the "error" Rune or "Unicode replacement character"
	RuneSelf  = 0x80         // characters below Runeself are represented as themselves in a single byte.
	MaxRune   = '\U0010FFFF' // Maximum valid Unicode code point.
	UTFMax    = 4            // maximum number of bytes of a UTF-8 encoded Unicode character.
)

// DecodeLastRune unpacks the last UTF-8 encoding in p and returns the rune and its
// width in bytes. If p is empty it returns (RuneError, 0). Otherwise, if the
// encoding is invalid, it returns (RuneError, 1). Both are impossible results for
// correct UTF-8.
//
// An encoding is invalid if it is incorrect UTF-8, encodes a rune that is out of
// range, or is not the shortest possible UTF-8 encoding for the value. No other
// validation is performed.

// DecodeLastRune 解包 p
// 中的最后一个UTF-8编码，并返回该符文及其字节宽度。 若此编码无效，它就会返回
// (RuneError, 1)，即一个对于正确的UTF-8来说不可能的值。
// 若一个编码为错误的UTF-8值，或该符文的编码超出范围，或不是该值可能的最短UTF-8编码，
// 那么它就是无效的。除此之外，并不进行其它的验证。
func DecodeLastRune(p []byte) (r rune, size int)

// DecodeLastRuneInString is like DecodeLastRune but its input is a string. If s is
// empty it returns (RuneError, 0). Otherwise, if the encoding is invalid, it
// returns (RuneError, 1). Both are impossible results for correct UTF-8.
//
// An encoding is invalid if it is incorrect UTF-8, encodes a rune that is out of
// range, or is not the shortest possible UTF-8 encoding for the value. No other
// validation is performed.

// DecodeLastRuneInString 类似于
// DecodeLastRune，但其输入为字符串。 若此编码无效，它就会返回
// (RuneError, 1)，即一个对于正确的UTF-8来说不可能的值。
// 若一个编码为错误的UTF-8值，或该符文的编码超出范围，或不是该值可能的最短UTF-8编码，
// 那么它就是无效的。除此之外，并不进行其它的验证。
func DecodeLastRuneInString(s string) (r rune, size int)

// DecodeRune unpacks the first UTF-8 encoding in p and returns the rune and its
// width in bytes. If p is empty it returns (RuneError, 0). Otherwise, if the
// encoding is invalid, it returns (RuneError, 1). Both are impossible results for
// correct UTF-8.
//
// An encoding is invalid if it is incorrect UTF-8, encodes a rune that is out of
// range, or is not the shortest possible UTF-8 encoding for the value. No other
// validation is performed.

// DecodeRune 解包 p
// 中的第一个UTF-8编码，并返回该符文及其字节宽度。 若此编码无效，它就会返回
// (RuneError, 1)，即一个对于正确的UTF-8来说不可能的值。
// 若一个编码为错误的UTF-8值，或该符文的编码超出范围，或不是该值可能的最短UTF-8编码，
// 那么它就是无效的。除此之外，并不进行其它的验证。
func DecodeRune(p []byte) (r rune, size int)

// DecodeRuneInString is like DecodeRune but its input is a string. If s is empty
// it returns (RuneError, 0). Otherwise, if the encoding is invalid, it returns
// (RuneError, 1). Both are impossible results for correct UTF-8.
//
// An encoding is invalid if it is incorrect UTF-8, encodes a rune that is out of
// range, or is not the shortest possible UTF-8 encoding for the value. No other
// validation is performed.

// DecodeRuneInString 类似于
// DecodeRune，但其输入为字符串。 若此编码无效，它就会返回
// (RuneError, 1)，即一个对于正确的UTF-8来说不可能的值。
// 若一个编码为错误的UTF-8值，或该符文的编码超出范围，或不是该值可能的最短UTF-8编码，
// 那么它就是无效的。除此之外，并不进行其它的验证。
func DecodeRuneInString(s string) (r rune, size int)

// EncodeRune writes into p (which must be large enough) the UTF-8 encoding of the
// rune. It returns the number of bytes written.

// EncodeRune 将该符文的UTF-8编码写入到 p 中（它必须足够大）。
// 它返回写入的字节数。
func EncodeRune(p []byte, r rune) int

// FullRune reports whether the bytes in p begin with a full UTF-8 encoding of a
// rune. An invalid encoding is considered a full Rune since it will convert as a
// width-1 error rune.

// FullRune 报告 p 中的字节是否以全UTF-8编码的符文开始。
// 无效的编码取被视作一个完整的符文，因为它会转换成宽度为1的错误符文。
func FullRune(p []byte) bool

// FullRuneInString is like FullRune but its input is a string.

// FullRuneInString 类似于 FullRune，但其输入为字符串。
func FullRuneInString(s string) bool

// RuneCount returns the number of runes in p. Erroneous and short encodings are
// treated as single runes of width 1 byte.

// RuneCount 返回 p 中的符文数。
// 错误编码和短编码将被视作宽度为1字节的单个符文。
func RuneCount(p []byte) int

// RuneCountInString is like RuneCount but its input is a string.

// RuneCountInString 类似于
// RuneCount，但其输入为字符串。
func RuneCountInString(s string) (n int)

// RuneLen returns the number of bytes required to encode the rune. It returns -1
// if the rune is not a valid value to encode in UTF-8.

// RuneLen 返回编码该符文所需的字节数。
// 若该符文并非有效的UTF-8编码值，就返回 -1。
func RuneLen(r rune) int

// RuneStart reports whether the byte could be the first byte of an encoded rune.
// Second and subsequent bytes always have the top two bits set to 10.

// RuneStart 报告该字节是否为符文编码的第一个字节。
// 第二个及后续字节的最高两位总是置为 10。
func RuneStart(b byte) bool

// Valid reports whether p consists entirely of valid UTF-8-encoded runes.

// Valid 报告 p 是否完全由有效的，UTF-8编码的符文构成。
func Valid(p []byte) bool

// ValidRune reports whether r can be legally encoded as UTF-8. Code points that
// are out of range or a surrogate half are illegal.

// ValidRune 报告 r 是否能合法地作为UTF-8编码。
// 超出返回或半替代值的码点是非法的。
func ValidRune(r rune) bool

// ValidString reports whether s consists entirely of valid UTF-8-encoded runes.

// ValidString 报告 s 是否完全由有效的，UTF-8编码的符文构成。
func ValidString(s string) bool
