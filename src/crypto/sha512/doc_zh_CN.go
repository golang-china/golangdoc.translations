// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package sha512 implements the SHA384 and SHA512 hash algorithms as defined in
// FIPS 180-2.

// sha512包实现了SHA384和SHA512哈希算法，参见FIPS 180-2。
package sha512

// The blocksize of SHA512 and SHA384 in bytes.

// SHA384和SHA512的字节块大小。
//
//	const Size = 64
//
// SHA512校验和的字节长度。
//
//	const Size384 = 48
//
// SHA384校验和的字节长度。
const BlockSize = 128

// The size of a SHA512 checksum in bytes.
const Size = 64

// The size of a SHA384 checksum in bytes.
const Size384 = 48

// New returns a new hash.Hash computing the SHA512 checksum.

// 返回一个新的使用SHA512校验算法的hash.Hash接口。
func New() hash.Hash

// New384 returns a new hash.Hash computing the SHA384 checksum.

// 返回一个新的使用SHA384校验算法的hash.Hash接口。
func New384() hash.Hash

// Sum384 returns the SHA384 checksum of the data.

// 返回数据的SHA384校验和。
func Sum384(data []byte) (sum384 [Size384]byte)

// Sum512 returns the SHA512 checksum of the data.

// 返回数据的SHA512校验和。
func Sum512(data []byte) [Size]byte
