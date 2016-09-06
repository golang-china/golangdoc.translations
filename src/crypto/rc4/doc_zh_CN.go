// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package rc4 implements RC4 encryption, as defined in Bruce Schneier's
// Applied Cryptography.

// rc4包实现了RC4加密算法，参见Bruce Schneier's Applied Cryptography。
package rc4

import "strconv"

// A Cipher is an instance of RC4 using a particular key.

// Cipher是一个使用特定密钥的RC4实例，本类型实现了cipher.Stream接口。
type Cipher struct {
}

type KeySizeError int

// NewCipher creates and returns a new Cipher. The key argument should be the
// RC4 key, at least 1 byte and at most 256 bytes.

// NewCipher创建并返回一个新的Cipher。参数key是RC4密钥，至少1字节，最多256字节。
func NewCipher(key []byte) (*Cipher, error)

// Reset zeros the key data so that it will no longer appear in the
// process's memory.

// Reset方法会清空密钥数据，以便将其数据从程序内存中清除（以免被破解）
func (c *Cipher) Reset()

func (k KeySizeError) Error() string

