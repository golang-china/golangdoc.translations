// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package crypto collects common cryptographic constants.

// crypto包搜集了常用的密码（算法）常量。
package crypto

import (
	"hash"
	"io"
	"strconv"
)

const (
	MD4       Hash = 1 + iota // import golang.org/x/crypto/md4 // 导入code.google.com/p/go.crypto/md4
	MD5                       // import crypto/md5 // 导入crypto/md5
	SHA1                      // import crypto/sha1 // 导入crypto/sha1
	SHA224                    // import crypto/sha256 // 导入crypto/sha256
	SHA256                    // import crypto/sha256 // 导入crypto/sha256
	SHA384                    // import crypto/sha512 // 导入crypto/sha512
	SHA512                    // import crypto/sha512 // 导入crypto/sha512
	MD5SHA1                   // no implementation; MD5+SHA1 used for TLS RSA // 未实现；MD5+SHA1用于TLS RSA
	RIPEMD160                 // import golang.org/x/crypto/ripemd160 // 导入code.google.com/p/go.crypto/ripemd160
	SHA3_224                  // import golang.org/x/crypto/sha3
	SHA3_256                  // import golang.org/x/crypto/sha3
	SHA3_384                  // import golang.org/x/crypto/sha3
	SHA3_512                  // import golang.org/x/crypto/sha3

)

// Hash identifies a cryptographic hash function that is implemented in another
// package.

// Hash用来识别/标识另一个包里实现的加密函数。
type Hash uint

// PrivateKey represents a private key using an unspecified algorithm.

// 代表一个使用未指定算法的私钥。
type PrivateKey interface{}

// PublicKey represents a public key using an unspecified algorithm.

// 代表一个使用未指定算法的公钥。
type PublicKey interface{}

// Signer is an interface for an opaque private key that can be used for signing
// operations. For example, an RSA key kept in a hardware module.
type Signer interface {
	// Public returns the public key corresponding to the opaque,
	// private key.
	Public() PublicKey

	// Sign signs msg with the private key, possibly using entropy from
	// rand. For an RSA key, the resulting signature should be either a
	// PKCS#1 v1.5 or PSS signature (as indicated by opts). For an (EC)DSA
	// key, it should be a DER-serialised, ASN.1 signature structure.
	//
	// Hash implements the SignerOpts interface and, in most cases, one can
	// simply pass in the hash function used as opts. Sign may also attempt
	// to type assert opts to other types in order to obtain algorithm
	// specific values. See the documentation in each package for details.
	Sign(rand io.Reader, msg []byte, opts SignerOpts) (signature []byte, err error)
}

// SignerOpts contains options for signing with a Signer.
type SignerOpts interface {
	// HashFunc returns an identifier for the hash function used to produce
	// the message passed to Signer.Sign, or else zero to indicate that no
	// hashing was done.
	HashFunc() Hash
}

// RegisterHash registers a function that returns a new instance of the given
// hash function. This is intended to be called from the init function in
// packages that implement hash functions.

// 注册一个返回给定hash接口实例的函数，并指定其标识值，该函数应在实现hash接口的
// 包的init函数中调用。
func RegisterHash(h Hash, f func() hash.Hash)

// Available reports whether the given hash function is linked into the binary.

// 报告是否有hash函数注册到该标识值。
func (Hash) Available() bool

// HashFunc simply returns the value of h so that Hash implements SignerOpts.
func (Hash) HashFunc() Hash

// New returns a new hash.Hash calculating the given hash function. New panics
// if the hash function is not linked into the binary.

// 创建一个使用给定hash函数的hash.Hash接口，如果该标识值未注册hash函数，将会
// panic。
func (Hash) New() hash.Hash

// Size returns the length, in bytes, of a digest resulting from the given hash
// function. It doesn't require that the hash function in question be linked
// into the program.

// 返回给定hash函数返回值的字节长度。
func (Hash) Size() int
