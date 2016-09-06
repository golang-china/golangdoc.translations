// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package sha256 implements the SHA224 and SHA256 hash algorithms as defined
// in FIPS 180-4.

// sha256包实现了SHA224和SHA256哈希算法，参见FIPS 180-4。
package sha256

import (
	"crypto"
	"hash"
)

// The blocksize of SHA256 and SHA224 in bytes.

// SHA224和SHA256的字节块大小。
//
//     const Size = 32
//
// SHA256校验和的字节长度。
//
//     const Size224 = 28
//
// SHA224校验和的字节长度。
const BlockSize = 64

// The size of a SHA256 checksum in bytes.
const Size = 32

// The size of a SHA224 checksum in bytes.
const Size224 = 28

// New returns a new hash.Hash computing the SHA256 checksum.

// 返回一个新的使用SHA256校验算法的hash.Hash接口。
func New() hash.Hash

// New224 returns a new hash.Hash computing the SHA224 checksum.

// 返回一个新的使用SHA224校验算法的hash.Hash接口。
func New224() hash.Hash

// Sum224 returns the SHA224 checksum of the data.

// 返回数据的SHA224校验和。
func Sum224(data []byte) (sum224 [Size224]byte)

// Sum256 returns the SHA256 checksum of the data.

// 返回数据的SHA256校验和。
func Sum256(data []byte) [Size]byte

