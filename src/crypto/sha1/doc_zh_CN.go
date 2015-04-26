// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package sha1 implements the SHA1 hash algorithm as defined in RFC 3174.

// sha1包实现了SHA1哈希算法，参见RFC 3174。
package sha1

// The blocksize of SHA1 in bytes.

// SHA1的块大小。
//
//	const Size = 20
//
// SHA1校验和的字节数。
const BlockSize = 64

// The size of a SHA1 checksum in bytes.
const Size = 20

// New returns a new hash.Hash computing the SHA1 checksum.

// 返回一个新的使用SHA1校验的hash.Hash接口。
func New() hash.Hash

// Sum returns the SHA1 checksum of the data.

// 返回数据data的SHA1校验和。
func Sum(data []byte) [Size]byte
