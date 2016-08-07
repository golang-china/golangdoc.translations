// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package hex implements hexadecimal encoding and decoding.

// hex包实现了16进制字符表示的编解码。
package hex

import (
    "bytes"
    "errors"
    "fmt"
    "io"
)

// ErrLength results from decoding an odd length slice.

// 解码一个长度为奇数的切片时，将返回此错误。
var ErrLength = errors.New("encoding/hex: odd length hex string")

// InvalidByteError values describe errors resulting from an invalid byte in a
// hex string.

// 描述一个hex编码字符串中的非法字符。
type InvalidByteError byte

// Decode decodes src into DecodedLen(len(src)) bytes, returning the actual
// number of bytes written to dst.
//
// If Decode encounters invalid input, it returns an error describing the
// failure.

// 将src解码为DecodedLen(len(src))字节，返回实际写入dst的字节数；如遇到非法字符
// ，返回描述错误的error。
func Decode(dst, src []byte) (int, error)

// DecodeString returns the bytes represented by the hexadecimal string s.

// 返回hex编码的字符串s代表的数据。
func DecodeString(s string) ([]byte, error)

// 长度x的编码数据解码后的明文数据的长度
func DecodedLen(x int) int

// Dump returns a string that contains a hex dump of the given data. The format
// of the hex dump matches the output of `hexdump -C` on the command line.

// 返回给定数据的hex
// dump格式的字符串，这个字符串与控制台下`hexdump
// -C`对该数据的输出是一致的。
func Dump(data []byte) string

// Dumper returns a WriteCloser that writes a hex dump of all written data to
// w. The format of the dump matches the output of `hexdump -C` on the command
// line.

// 返回一个io.WriteCloser接口，将写入的数据的hex
// dump格式写入w，具体格式为'hexdump -C'。
func Dumper(w io.Writer) io.WriteCloser

// Encode encodes src into EncodedLen(len(src))
// bytes of dst.  As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements hexadecimal encoding.

// 将src的数据解码为EncodedLen(len(src))字节，返回实际写入dst的字节数：
// EncodedLen(len(src))。
func Encode(dst, src []byte) int

// EncodeToString returns the hexadecimal encoding of src.

// 将数据src编码为字符串s。
func EncodeToString(src []byte) string

// EncodedLen returns the length of an encoding of n source bytes.

// 长度x的明文数据编码后的编码数据的长度。
func EncodedLen(n int) int

func (InvalidByteError) Error() string

