// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package adler32 implements the Adler-32 checksum.
//
// It is defined in RFC 1950:
//     Adler-32 is composed of two sums accumulated per byte: s1 is
//     the sum of all bytes, s2 is the sum of all s1 values. Both sums
//     are done modulo 65521. s1 is initialized to 1, s2 to zero.  The
//     Adler-32 checksum is stored as s2*65536 + s1 in most-
//     significant-byte first (network) order.

// Package adler32 implements the Adler-32 checksum.
//
// It is defined in RFC 1950:
//     Adler-32 is composed of two sums accumulated per byte: s1 is
//     the sum of all bytes, s2 is the sum of all s1 values. Both sums
//     are done modulo 65521. s1 is initialized to 1, s2 to zero.  The
//     Adler-32 checksum is stored as s2*65536 + s1 in most-
//     significant-byte first (network) order.
package adler32

import "hash"

// The size of an Adler-32 checksum in bytes.
const Size = 4


// Checksum returns the Adler-32 checksum of data.
func Checksum(data []byte) uint32

// New returns a new hash.Hash32 computing the Adler-32 checksum.
// Its Sum method will lay the value out in big-endian byte order.
func New() hash.Hash32

