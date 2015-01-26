// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package des implements the Data Encryption Standard (DES) and the Triple Data
// Encryption Algorithm (TDEA) as defined in U.S. Federal Information Processing
// Standards Publication 46-3.

// Package des implements the Data
// Encryption Standard (DES) and the Triple
// Data Encryption Algorithm (TDEA) as
// defined in U.S. Federal Information
// Processing Standards Publication 46-3.
package des

// The DES block size in bytes.

// The DES block size in bytes.
const BlockSize = 8

// NewCipher creates and returns a new cipher.Block.

// NewCipher creates and returns a new
// cipher.Block.
func NewCipher(key []byte) (cipher.Block, error)

// NewTripleDESCipher creates and returns a new cipher.Block.

// NewTripleDESCipher creates and returns a
// new cipher.Block.
func NewTripleDESCipher(key []byte) (cipher.Block, error)

type KeySizeError int

func (k KeySizeError) Error() string
