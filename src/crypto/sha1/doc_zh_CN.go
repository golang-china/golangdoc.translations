// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package sha1 implements the SHA1 hash algorithm as defined in RFC 3174.

// Package sha1 implements the SHA1 hash
// algorithm as defined in RFC 3174.
package sha1

// The blocksize of SHA1 in bytes.

// The blocksize of SHA1 in bytes.
const BlockSize = 64

// The size of a SHA1 checksum in bytes.

// The size of a SHA1 checksum in bytes.
const Size = 20

// New returns a new hash.Hash computing the SHA1 checksum.

// New returns a new hash.Hash computing
// the SHA1 checksum.
func New() hash.Hash

// Sum returns the SHA1 checksum of the data.

// Sum returns the SHA1 checksum of the
// data.
func Sum(data []byte) [Size]byte
