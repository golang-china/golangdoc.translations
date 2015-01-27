// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package aes implements AES encryption (formerly Rijndael), as defined in U.S.
// Federal Information Processing Standards Publication 197.
package aes

// The AES block size in bytes.
const BlockSize = 16

// NewCipher creates and returns a new cipher.Block. The key argument should be the
// AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func NewCipher(key []byte) (cipher.Block, error)

type KeySizeError int

func (k KeySizeError) Error() string
