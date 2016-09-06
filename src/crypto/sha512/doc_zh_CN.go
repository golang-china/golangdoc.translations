// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package sha512 implements the SHA-384, SHA-512, SHA-512/224, and SHA-512/256
// hash algorithms as defined in FIPS 180-4.

// sha512包实现了SHA384和SHA512哈希算法，参见FIPS 180-2。
package sha512

import (
	"crypto"
	"hash"
)

const (
	// Size is the size, in bytes, of a SHA-512 checksum.

	// The size of a SHA512 checksum in bytes.
	Size = 64

	// Size224 is the size, in bytes, of a SHA-512/224 checksum.
	Size224 = 28

	// Size256 is the size, in bytes, of a SHA-512/256 checksum.
	Size256 = 32

	// Size384 is the size, in bytes, of a SHA-384 checksum.
	Size384 = 48

	// BlockSize is the block size, in bytes, of the SHA-512/224,
	// SHA-512/256, SHA-384 and SHA-512 hash functions.
	BlockSize = 128
)

// New returns a new hash.Hash computing the SHA-512 checksum.

// 返回一个新的使用SHA512校验算法的hash.Hash接口。
func New() hash.Hash

// New384 returns a new hash.Hash computing the SHA-384 checksum.

// 返回一个新的使用SHA384校验算法的hash.Hash接口。
func New384() hash.Hash

// New512_224 returns a new hash.Hash computing the SHA-512/224 checksum.
func New512_224() hash.Hash

// New512_256 returns a new hash.Hash computing the SHA-512/256 checksum.
func New512_256() hash.Hash

// Sum384 returns the SHA384 checksum of the data.

// 返回数据的SHA384校验和。
func Sum384(data []byte) (sum384 [Size384]byte)

// Sum512 returns the SHA512 checksum of the data.

// 返回数据的SHA512校验和。
func Sum512(data []byte) [Size]byte

// Sum512_224 returns the Sum512/224 checksum of the data.
func Sum512_224(data []byte) (sum224 [Size224]byte)

// Sum512_256 returns the Sum512/256 checksum of the data.
func Sum512_256(data []byte) (sum256 [Size256]byte)

