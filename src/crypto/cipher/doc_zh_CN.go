// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package cipher implements standard block cipher modes that can be wrapped around
// low-level block cipher implementations. See
// http://csrc.nist.gov/groups/ST/toolkit/BCM/current_modes.html and NIST Special
// Publication 800-38A.
package cipher

// AEAD is a cipher mode providing authenticated encryption with associated data.
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

// NewGCM returns the given 128-bit, block cipher wrapped in Galois Counter Mode.
func NewGCM(cipher Block) (AEAD, error)

// A Block represents an implementation of block cipher using a given key. It
// provides the capability to encrypt or decrypt individual blocks. The mode
// implementations extend that capability to streams of blocks.
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

// A BlockMode represents a block cipher running in a block-based mode (CBC, ECB
// etc).
type BlockMode interface {
	// BlockSize returns the mode's block size.
	BlockSize() int

	// CryptBlocks encrypts or decrypts a number of blocks. The length of
	// src must be a multiple of the block size. Dst and src may point to
	// the same memory.
	CryptBlocks(dst, src []byte)
}

// NewCBCDecrypter returns a BlockMode which decrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the Block's
// block size and must match the iv used to encrypt the data.
func NewCBCDecrypter(b Block, iv []byte) BlockMode

// NewCBCEncrypter returns a BlockMode which encrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the Block's
// block size.
func NewCBCEncrypter(b Block, iv []byte) BlockMode

// A Stream represents a stream cipher.
type Stream interface {
	// XORKeyStream XORs each byte in the given slice with a byte from the
	// cipher's key stream. Dst and src may point to the same memory.
	XORKeyStream(dst, src []byte)
}

// NewCFBDecrypter returns a Stream which decrypts with cipher feedback mode, using
// the given Block. The iv must be the same length as the Block's block size.
func NewCFBDecrypter(block Block, iv []byte) Stream

// NewCFBEncrypter returns a Stream which encrypts with cipher feedback mode, using
// the given Block. The iv must be the same length as the Block's block size.
func NewCFBEncrypter(block Block, iv []byte) Stream

// NewCTR returns a Stream which encrypts/decrypts using the given Block in counter
// mode. The length of iv must be the same as the Block's block size.
func NewCTR(block Block, iv []byte) Stream

// NewOFB returns a Stream that encrypts or decrypts using the block cipher b in
// output feedback mode. The initialization vector iv's length must be equal to b's
// block size.
func NewOFB(b Block, iv []byte) Stream

// StreamReader wraps a Stream into an io.Reader. It calls XORKeyStream to process
// each slice of data which passes through.
type StreamReader struct {
	S Stream
	R io.Reader
}

func (r StreamReader) Read(dst []byte) (n int, err error)

// StreamWriter wraps a Stream into an io.Writer. It calls XORKeyStream to process
// each slice of data which passes through. If any Write call returns short then
// the StreamWriter is out of sync and must be discarded. A StreamWriter has no
// internal buffering; Close does not need to be called to flush write data.
type StreamWriter struct {
	S   Stream
	W   io.Writer
	Err error // unused
}

// Close closes the underlying Writer and returns its Close return value, if the
// Writer is also an io.Closer. Otherwise it returns nil.
func (w StreamWriter) Close() error

func (w StreamWriter) Write(src []byte) (n int, err error)
