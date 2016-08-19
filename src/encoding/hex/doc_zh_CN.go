// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package hex implements hexadecimal encoding and decoding.

// Package hex implements hexadecimal encoding and decoding.
package hex

import (
    "bytes"
    "errors"
    "fmt"
    "io"
)

// ErrLength results from decoding an odd length slice.
var ErrLength = errors.New("encoding/hex: odd length hex string")


// InvalidByteError values describe errors resulting from an invalid byte in a
// hex string.

// InvalidByteError values describe errors resulting from an invalid byte in a
// hex string.
type InvalidByteError byte


// Decode decodes src into DecodedLen(len(src)) bytes, returning the actual
// number of bytes written to dst.
//
// If Decode encounters invalid input, it returns an error describing the
// failure.
func Decode(dst, src []byte) (int, error)

// DecodeString returns the bytes represented by the hexadecimal string s.
func DecodeString(s string) ([]byte, error)

func DecodedLen(x int) int

// Dump returns a string that contains a hex dump of the given data. The format
// of the hex dump matches the output of `hexdump -C` on the command line.
func Dump(data []byte) string

// Dumper returns a WriteCloser that writes a hex dump of all written data to
// w. The format of the dump matches the output of `hexdump -C` on the command
// line.
func Dumper(w io.Writer) io.WriteCloser

// Encode encodes src into EncodedLen(len(src))
// bytes of dst.  As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements hexadecimal encoding.

// Encode encodes src into EncodedLen(len(src))
// bytes of dst. As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements hexadecimal encoding.
func Encode(dst, src []byte) int

// EncodeToString returns the hexadecimal encoding of src.
func EncodeToString(src []byte) string

// EncodedLen returns the length of an encoding of n source bytes.
func EncodedLen(n int) int

func (InvalidByteError) Error() string

