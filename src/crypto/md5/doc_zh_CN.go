// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package md5 implements the MD5 hash algorithm as defined in RFC 1321.

// md5 包实现了在 RFC 1321 中定义的 MD5 哈希算法.
package md5

// The blocksize of MD5 in bytes.

// MD5 块大小，以字节为单位.
const BlockSize = 64

// The size of an MD5 checksum in bytes.

// MD5 校验和的大小，以字节为单位.
const Size = 16

// New returns a new hash.Hash computing the MD5 checksum.

// New 返回一个新的计算 MD5 校验和的 hash.Hash 接口.
func New() hash.Hash

// Sum returns the MD5 checksum of the data.

// Sum 返回 data 的 MD5 校验和.
func Sum(data []byte) [Size]byte
