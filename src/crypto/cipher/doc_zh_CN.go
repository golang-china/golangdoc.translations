// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package cipher implements standard block cipher modes that can be wrapped
// around low-level block cipher implementations.
// See http://csrc.nist.gov/groups/ST/toolkit/BCM/current_modes.html
// and NIST Special Publication 800-38A.

// cipher包实现了多个标准的用于包装底层块加密算法的加密算法实现。
//
// 参见http://csrc.nist.gov/groups/ST/toolkit/BCM/current_modes.html和NIST
// Special Publication 800-38A。
package cipher

import (
    "crypto/subtle"
    "errors"
    "io"
    "runtime"
    "unsafe"
)

// AEAD is a cipher mode providing authenticated encryption with associated
// data. For a description of the methodology, see
//     https://en.wikipedia.org/wiki/Authenticated_encryption

// AEAD接口是一种提供了使用关联数据进行认证加密的功能的加密模式。
type AEAD interface {
    // NonceSize returns the size of the nonce that must be passed to Seal
    // and Open.
    NonceSize() int

    // Overhead returns the maximum difference between the lengths of a
    // plaintext and ciphertext.
    Overhead() int

    // Seal encrypts and authenticates plaintext, authenticates the
    // additional data and appends the result to dst, returning the updated
    // slice. The nonce must be NonceSize() bytes long and unique for all
    // time, for a given key.
    //
    // The plaintext and dst may alias exactly or not at all.
    Seal(dst, nonce, plaintext, data []byte) []byte

    // Open decrypts and authenticates ciphertext, authenticates the
    // additional data and, if successful, appends the resulting plaintext
    // to dst, returning the updated slice. The nonce must be NonceSize()
    // bytes long and both it and the additional data must match the
    // value passed to Seal.
    //
    // The ciphertext and dst may alias exactly or not at all.
    Open(dst, nonce, ciphertext, data []byte) ([]byte, error)
}

// A Block represents an implementation of block cipher
// using a given key.  It provides the capability to encrypt
// or decrypt individual blocks.  The mode implementations
// extend that capability to streams of blocks.

// Block接口代表一个使用特定密钥的底层块加/解密器。它提供了加密和解密独立数据块
// 的能力。
type Block interface {
    // BlockSize returns the cipher's block size.
    BlockSize() int

    // Encrypt encrypts the first block in src into dst.
    // Dst and src may point at the same memory.
    Encrypt(dst, src []byte)

    // Decrypt decrypts the first block in src into dst.
    // Dst and src may point at the same memory.
    Decrypt(dst, src []byte)
}

// A BlockMode represents a block cipher running in a block-based mode (CBC,
// ECB etc).

// BlockMode接口代表一个工作在块模式（如CBC、ECB等）的加/解密器。
type BlockMode interface {
    // BlockSize returns the mode's block size.
    BlockSize() int

    // CryptBlocks encrypts or decrypts a number of blocks. The length of
    // src must be a multiple of the block size. Dst and src may point to
    // the same memory.
    CryptBlocks(dst, src []byte)
}

// A Stream represents a stream cipher.

// Stream接口代表一个流模式的加/解密器。
type Stream interface {
    // XORKeyStream XORs each byte in the given slice with a byte from the
    // cipher's key stream. Dst and src may point to the same memory.
    XORKeyStream(dst, src []byte)
}

// StreamReader wraps a Stream into an io.Reader. It calls XORKeyStream
// to process each slice of data which passes through.

// 将一个Stream与一个io.Reader接口关联起来，Read方法会调用XORKeyStream方法来处理
// 获取的所有切片。
type StreamReader struct {
    S   Stream
    R   io.Reader
}

// StreamWriter wraps a Stream into an io.Writer. It calls XORKeyStream
// to process each slice of data which passes through. If any Write call
// returns short then the StreamWriter is out of sync and must be discarded.
// A StreamWriter has no internal buffering; Close does not need
// to be called to flush write data.

// 将一个Stream与一个io.Writer接口关联起来，Write方法会调用XORKeyStream方法来处
// 理提供的所有切片。如果Write方法返回的n小于提供的切片的长度，则表示
// StreamWriter不同步，必须丢弃。StreamWriter没有内建的缓存，不需要调用Close方法
// 去清空缓存。
type StreamWriter struct {
    S   Stream
    W   io.Writer
    Err error // unused
}

// NewCBCDecrypter returns a BlockMode which decrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the Block's
// block size and must match the iv used to encrypt the data.
func NewCBCDecrypter(b Block, iv []byte) BlockMode

// NewCBCEncrypter returns a BlockMode which encrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the
// Block's block size.

// 返回一个密码分组链接模式的、底层用b加密的BlockMode接口，初始向量iv的长度必须
// 等于b的块尺寸。
func NewCBCEncrypter(b Block, iv []byte) BlockMode

// NewCFBDecrypter returns a Stream which decrypts with cipher feedback mode,
// using the given Block. The iv must be the same length as the Block's block
// size.
func NewCFBDecrypter(block Block, iv []byte) Stream

// NewCFBEncrypter returns a Stream which encrypts with cipher feedback mode,
// using the given Block. The iv must be the same length as the Block's block
// size.

// 返回一个密码反馈模式的、底层用block加密的Stream接口，初始向量iv的长度必须等于
// block的块尺寸。
func NewCFBEncrypter(block Block, iv []byte) Stream

// NewCTR returns a Stream which encrypts/decrypts using the given Block in
// counter mode. The length of iv must be the same as the Block's block size.
func NewCTR(block Block, iv []byte) Stream

// NewGCM returns the given 128-bit, block cipher wrapped in Galois Counter Mode
// with the standard nonce length.

// 函数用迦洛瓦计数器模式包装提供的128位Block接口，并返回AEAD接口。
func NewGCM(cipher Block) (AEAD, error)

// NewOFB returns a Stream that encrypts or decrypts using the block cipher b in
// output feedback mode. The initialization vector iv's length must be equal to
// b's block size.
func NewOFB(b Block, iv []byte) Stream

func (StreamReader) Read(dst []byte) (n int, err error)

// Close closes the underlying Writer and returns its Close return value, if the
// Writer is also an io.Closer. Otherwise it returns nil.
func (StreamWriter) Close() error

func (StreamWriter) Write(src []byte) (n int, err error)

