// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package des implements the Data Encryption Standard (DES) and the
// Triple Data Encryption Algorithm (TDEA) as defined
// in U.S. Federal Information Processing Standards Publication 46-3.

// des包实现了DES标准和TDEA算法，参见U.S. Federal Information Processing
// Standards Publication 46-3。
package des

import (
    "crypto/cipher"
    "encoding/binary"
    "strconv"
)

// The DES block size in bytes.

// DES字节块的大小。
const BlockSize = 8

type KeySizeError int

// NewCipher creates and returns a new cipher.Block.

// 创建并返回一个使用DES算法的cipher.Block接口。
func NewCipher(key []byte) (cipher.Block, error)

// NewTripleDESCipher creates and returns a new cipher.Block.

// 创建并返回一个使用TDEA算法的cipher.Block接口。
func NewTripleDESCipher(key []byte) (cipher.Block, error)

func (KeySizeError) Error() string

