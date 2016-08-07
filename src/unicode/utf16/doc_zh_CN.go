// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package utf16 implements encoding and decoding of UTF-16 sequences.

// utf16 包实现了对UTF-16序列的编码和解码。
package utf16

// Decode returns the Unicode code point sequence represented
// by the UTF-16 encoding s.

// Decode 返回由UTF-16编码 s 所表示的Unicode码点序列。
func Decode(s []uint16) []rune

// DecodeRune returns the UTF-16 decoding of a surrogate pair.
// If the pair is not a valid UTF-16 surrogate pair, DecodeRune returns
// the Unicode replacement code point U+FFFD.

// DecodeRune 返回替代值对的UTF-16解码。
// 若该值对并非有效的UTF-16替代值对，DecodeRune
// 就会返回Unicode的替换码点U+FFFD。
func DecodeRune(r1, r2 rune) rune

// Encode returns the UTF-16 encoding of the Unicode code point sequence s.

// Encode 返回Unicode码点序列 s 的UTF-16编码。
func Encode(s []rune) []uint16

// EncodeRune returns the UTF-16 surrogate pair r1, r2 for the given rune.
// If the rune is not a valid Unicode code point or does not need encoding,
// EncodeRune returns U+FFFD, U+FFFD.

// EncodeRune 返回给定符文的UTF-16替代值对 r1, r2。
// 若该符文并非有效的Unicode码点或无需编码，EncodeRune 就会返回 U+FFFD, U+FFFD。
func EncodeRune(r rune) (r1, r2 rune)

// IsSurrogate reports whether the specified Unicode code point
// can appear in a surrogate pair.

// IsSurrogate
// 在指定的Unicode码点可出现在替代值对中时返回 true。
func IsSurrogate(r rune) bool

