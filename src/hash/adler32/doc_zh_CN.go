// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package adler32 implements the Adler-32 checksum.
//
// It is defined in RFC 1950:
//
//	Adler-32 is composed of two sums accumulated per byte: s1 is
//	the sum of all bytes, s2 is the sum of all s1 values. Both sums
//	are done modulo 65521. s1 is initialized to 1, s2 to zero.  The
//	Adler-32 checksum is stored as s2*65536 + s1 in most-
//	significant-byte first (network) order.

// adler32包实现了Adler-32校验和算法，参见RFC 1950：
//
//	Adler-32由两个每字节累积的和组成：
//	s1是所有字节的累积，s2是所有s1的累积。两个累积值都取65521的余数。s1初始为1，s2初始为0。
//	Afler-32校验和保存为s2*65536 + s1。（最高有效字节在前/大端在前）
package adler32

// The size of an Adler-32 checksum in bytes.

// Adler-32校验和的字节数。
const Size = 4

// Checksum returns the Adler-32 checksum of data.

// 返回数据data的Adler-32校验和。
func Checksum(data []byte) uint32

// New returns a new hash.Hash32 computing the Adler-32 checksum.

// 返回一个计算Adler-32校验和的hash.Hash32接口。
func New() hash.Hash32
