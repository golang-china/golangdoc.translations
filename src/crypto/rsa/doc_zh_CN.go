// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package rsa implements RSA encryption as specified in PKCS#1.
//
// RSA is a single, fundamental operation that is used in this package to
// implement either public-key encryption or public-key signatures.
//
// The original specification for encryption and signatures with RSA is PKCS#1
// and the terms "RSA encryption" and "RSA signatures" by default refer to
// PKCS#1 version 1.5. However, that specification has flaws and new designs
// should use version two, usually called by just OAEP and PSS, where possible.
//
// Two sets of interfaces are included in this package. When a more abstract
// interface isn't necessary, there are functions for encrypting/decrypting with
// v1.5/OAEP and signing/verifying with v1.5/PSS. If one needs to abstract over
// the public-key primitive, the PrivateKey struct implements the Decrypter and
// Signer interfaces from the crypto package.

// rsa包实现了PKCS#1规定的RSA加密算法。
package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"hash"
	"io"
	"math/big"
)

const (
	// PSSSaltLengthAuto causes the salt in a PSS signature to be as large
	// as possible when signing, and to be auto-detected when verifying.
	PSSSaltLengthAuto = 0

	// PSSSaltLengthEqualsHash causes the salt length to equal the length
	// of the hash used in the signature.
	PSSSaltLengthEqualsHash = -1
)

// ErrDecryption represents a failure to decrypt a message.
// It is deliberately vague to avoid adaptive attacks.

// ErrDecryption 代表解密数据失败。它故意写的语焉不详，以避免适应性攻击。
var ErrDecryption = errors.New("crypto/rsa: decryption error")

// ErrMessageTooLong is returned when attempting to encrypt a message which is
// too large for the size of the public key.

// 当试图用公钥加密尺寸过大的数据时，就会返回ErrMessageTooLong。
var ErrMessageTooLong = errors.New("crypto/rsa: message too long for RSA public key size")

// ErrVerification represents a failure to verify a signature.
// It is deliberately vague to avoid adaptive attacks.

// ErrVerification代表认证签名失败。它故意写的语焉不详，以避免适应性攻击。
var ErrVerification = errors.New("crypto/rsa: verification error")

// CRTValue contains the precomputed Chinese remainder theorem values.

// CRTValue包含预先计算的中国剩余定理的值。
type CRTValue struct {
	Exp   *big.Int // D mod (prime-1).
	Coeff *big.Int // R·Coeff ≡ 1 mod Prime.
	R     *big.Int // product of primes prior to this (inc p and q).
}

// OAEPOptions is an interface for passing options to OAEP decryption using the
// crypto.Decrypter interface.
type OAEPOptions struct {
	// Hash is the hash function that will be used when generating the mask.
	Hash crypto.Hash

	// Label is an arbitrary byte string that must be equal to the value
	// used when encrypting.
	Label []byte
}

// PKCS1v15DecrypterOpts is for passing options to PKCS#1 v1.5 decryption using
// the crypto.Decrypter interface.
type PKCS1v15DecryptOptions struct {
	// SessionKeyLen is the length of the session key that is being
	// decrypted. If not zero, then a padding error during decryption will
	// cause a random plaintext of this length to be returned rather than
	// an error. These alternatives happen in constant time.
	SessionKeyLen int
}

// PSSOptions contains options for creating and verifying PSS signatures.

// PSSOptions包含用于创建和认证PSS签名的参数。
type PSSOptions struct {
	// SaltLength controls the length of the salt used in the PSS
	// signature. It can either be a number of bytes, or one of the special
	// PSSSaltLength constants.
	SaltLength int

	// Hash, if not zero, overrides the hash function passed to SignPSS.
	// This is the only way to specify the hash function when using the
	// crypto.Signer interface.
	Hash crypto.Hash
}

type PrecomputedValues struct {
	Dp, Dq *big.Int // D mod (P-1) (or mod Q-1)
	Qinv   *big.Int // Q^-1 mod P

	// CRTValues is used for the 3rd and subsequent primes. Due to a
	// historical accident, the CRT for the first two primes is handled
	// differently in PKCS#1 and interoperability is sufficiently
	// important that we mirror this.
	CRTValues []CRTValue
}

// A PrivateKey represents an RSA key

// 代表一个RSA私钥。
type PrivateKey struct {
	PublicKey // public part.
	D         *big.Int   // private exponent
	Primes    []*big.Int // prime factors of N, has >= 2 elements.

	// Precomputed contains precomputed values that speed up private
	// operations, if available.
	Precomputed PrecomputedValues
}

// A PublicKey represents the public part of an RSA key.

// 代表一个RSA公钥。
type PublicKey struct {
	N *big.Int // modulus
	E int      // public exponent
}

// OAEP is parameterised by a hash function that is used as a random oracle.
// Encryption and decryption of a given message must use the same hash function
// and sha256.New() is a reasonable choice.
//
// The random parameter, if not nil, is used to blind the private-key operation
// and avoid timing side-channel attacks. Blinding is purely internal to this
// function – the random data need not match that used when encrypting.
//
// The label parameter must match the value given when encrypting. See
// EncryptOAEP for details.

// DecryptOAEP解密RSA-OAEP算法加密的数据。如果random不是nil，函数会注意规避时间
// 侧信道攻击。
func DecryptOAEP(hash hash.Hash, random io.Reader, priv *PrivateKey, ciphertext []byte, label []byte) ([]byte, error)

// DecryptPKCS1v15 decrypts a plaintext using RSA and the padding scheme from
// PKCS#1 v1.5. If rand != nil, it uses RSA blinding to avoid timing
// side-channel attacks.
//
// Note that whether this function returns an error or not discloses secret
// information. If an attacker can cause this function to run repeatedly and
// learn whether each instance returned an error then they can decrypt and forge
// signatures as if they had the private key. See DecryptPKCS1v15SessionKey for
// a way of solving this problem.

// DecryptPKCS1v15使用PKCS#1 v1.5规定的填充方案和RSA算法解密密文。如果random不是
// nil，函数会注意规避时间侧信道攻击。
func DecryptPKCS1v15(rand io.Reader, priv *PrivateKey, ciphertext []byte) ([]byte, error)

// DecryptPKCS1v15SessionKey decrypts a session key using RSA and the padding
// scheme from PKCS#1 v1.5. If rand != nil, it uses RSA blinding to avoid timing
// side-channel attacks. It returns an error if the ciphertext is the wrong
// length or if the ciphertext is greater than the public modulus. Otherwise, no
// error is returned. If the padding is valid, the resulting plaintext message
// is copied into key. Otherwise, key is unchanged. These alternatives occur in
// constant time. It is intended that the user of this function generate a
// random session key beforehand and continue the protocol with the resulting
// value. This will remove any possibility that an attacker can learn any
// information about the plaintext. See ``Chosen Ciphertext Attacks Against
// Protocols Based on the RSA Encryption Standard PKCS #1'', Daniel
// Bleichenbacher, Advances in Cryptology (Crypto '98).
//
// Note that if the session key is too small then it may be possible for an
// attacker to brute-force it. If they can do that then they can learn whether a
// random value was used (because it'll be different for the same ciphertext)
// and thus whether the padding was correct. This defeats the point of this
// function. Using at least a 16-byte key will protect against this attack.

// DecryptPKCS1v15SessionKey使用PKCS#1 v1.5规定的填充方案和RSA算法解密会话密钥。
// 如果random不是nil，函数会注意规避时间侧信道攻击。
//
// 如果密文长度不对，或者如果密文比公共模数的长度还长，会返回错误；否则，不会返
// 回任何错误。如果填充是合法的，生成的明文信息会拷贝进key；否则，key不会被修改
// 。这些情况都会在固定时间内出现（规避时间侧信道攻击）。本函数的目的是让程序的
// 使用者事先生成一个随机的会话密钥，并用运行时的值继续协议。这样可以避免任何攻
// 击者从明文窃取信息的可能性。
//
// 参见”Chosen Ciphertext Attacks Against Protocols Based on the RSA Encryption
// Standard PKCS #1”。
func DecryptPKCS1v15SessionKey(rand io.Reader, priv *PrivateKey, ciphertext []byte, key []byte) error

// EncryptOAEP encrypts the given message with RSA-OAEP.
//
// OAEP is parameterised by a hash function that is used as a random oracle.
// Encryption and decryption of a given message must use the same hash function
// and sha256.New() is a reasonable choice.
//
// The random parameter is used as a source of entropy to ensure that encrypting
// the same message twice doesn't result in the same ciphertext.
//
// The label parameter may contain arbitrary data that will not be encrypted,
// but which gives important context to the message. For example, if a given
// public key is used to decrypt two types of messages then distinct label
// values could be used to ensure that a ciphertext for one purpose cannot be
// used for another by an attacker. If not required it can be empty.
//
// The message must be no longer than the length of the public modulus less
// twice the hash length plus 2.

// 采用RSA-OAEP算法加密给出的数据。数据不能超过((公共模数的长度)-2*( hash长度
// )+2)字节。
func EncryptOAEP(hash hash.Hash, random io.Reader, pub *PublicKey, msg []byte, label []byte) ([]byte, error)

// EncryptPKCS1v15 encrypts the given message with RSA and the padding
// scheme from PKCS#1 v1.5.  The message must be no longer than the
// length of the public modulus minus 11 bytes.
//
// The rand parameter is used as a source of entropy to ensure that
// encrypting the same message twice doesn't result in the same
// ciphertext.
//
// WARNING: use of this function to encrypt plaintexts other than
// session keys is dangerous. Use RSA OAEP in new protocols.

// EncryptPKCS1v15使用PKCS#1 v1.5规定的填充方案和RSA算法加密msg。信息不能超过((
// 公共模数的长度)-11)字节。注意：使用本函数加密明文（而不是会话密钥）是危险的，
// 请尽量在新协议中使用RSA OAEP。
func EncryptPKCS1v15(rand io.Reader, pub *PublicKey, msg []byte) ([]byte, error)

// GenerateKey generates an RSA keypair of the given bit size using the
// random source random (for example, crypto/rand.Reader).

// GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥。
func GenerateKey(random io.Reader, bits int) (*PrivateKey, error)

// GenerateMultiPrimeKey generates a multi-prime RSA keypair of the given bit
// size and the given random source, as suggested in [1]. Although the public
// keys are compatible (actually, indistinguishable) from the 2-prime case, the
// private keys are not. Thus it may not be possible to export multi-prime
// private keys in certain formats or to subsequently import them into other
// code.
//
// Table 1 in [2] suggests maximum numbers of primes for a given size.
//
// [1] US patent 4405829 (1972, expired) [2]
// http://www.cacr.math.uwaterloo.ca/techreports/2006/cacr2006-16.pdf

// GenerateMultiPrimeKey使用指定的字位数生成一对多质数的RSA密钥，参见US patent
// 4405829。虽然公钥可以和二质数情况下的公钥兼容（事实上，不能区分两种公钥），私
// 钥却不行。因此有可能无法生成特定格式的多质数的密钥对，或不能将生成的密钥用在
// 其他（语言的）代码里。
//
// http://www.cacr.math.uwaterloo.ca/techreports/2006/cacr2006-16.pdf中的Table 1
// 说明了给定字位数的密钥可以接受的质数最大数量。
func GenerateMultiPrimeKey(random io.Reader, nprimes int, bits int) (*PrivateKey, error)

// SignPKCS1v15 calculates the signature of hashed using
// RSASSA-PKCS1-V1_5-SIGN from RSA PKCS#1 v1.5.  Note that hashed must
// be the result of hashing the input message using the given hash
// function. If hash is zero, hashed is signed directly. This isn't
// advisable except for interoperability.
//
// If rand is not nil then RSA blinding will be used to avoid timing
// side-channel attacks.
//
// This function is deterministic. Thus, if the set of possible
// messages is small, an attacker may be able to build a map from
// messages to signatures and identify the signed messages. As ever,
// signatures provide authenticity, not confidentiality.

// SignPKCS1v15使用RSA PKCS#1 v1.5规定的RSASSA-PKCS1-V1_5-SIGN签名方案计算签名。
// 注意hashed必须是使用提供给本函数的hash参数对（要签名的）原始数据进行hash的结
// 果。
func SignPKCS1v15(rand io.Reader, priv *PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error)

// SignPSS calculates the signature of hashed using RSASSA-PSS [1]. Note that
// hashed must be the result of hashing the input message using the given hash
// function. The opts argument may be nil, in which case sensible defaults are
// used.

// SignPSS采用RSASSA-PSS方案计算签名。注意hashed必须是使用提供给本函数的hash参数
// 对（要签名的）原始数据进行hash的结果。opts参数可以为nil，此时会使用默认参数。
func SignPSS(rand io.Reader, priv *PrivateKey, hash crypto.Hash, hashed []byte, opts *PSSOptions) ([]byte, error)

// VerifyPKCS1v15 verifies an RSA PKCS#1 v1.5 signature.
// hashed is the result of hashing the input message using the given hash
// function and sig is the signature. A valid signature is indicated by
// returning a nil error. If hash is zero then hashed is used directly. This
// isn't advisable except for interoperability.

// VerifyPKCS1v15认证RSA PKCS#1 v1.5签名。hashed是使用提供的hash参数对（要签名的
// ）原始数据进行hash的结果。合法的签名会返回nil，否则表示签名不合法。
func VerifyPKCS1v15(pub *PublicKey, hash crypto.Hash, hashed []byte, sig []byte) error

// VerifyPSS verifies a PSS signature. hashed is the result of hashing the input
// message using the given hash function and sig is the signature. A valid
// signature is indicated by returning a nil error. The opts argument may be
// nil, in which case sensible defaults are used.

// VerifyPSS认证一个PSS签名。hashed是使用提供给本函数的hash参数对（要签名的）原
// 始数据进行hash的结果。合法的签名会返回nil，否则表示签名不合法。opts参数可以为
// nil，此时会使用默认参数。
func VerifyPSS(pub *PublicKey, hash crypto.Hash, hashed []byte, sig []byte, opts *PSSOptions) error

// HashFunc returns pssOpts.Hash so that PSSOptions implements
// crypto.SignerOpts.
func (pssOpts *PSSOptions) HashFunc() crypto.Hash

// Decrypt decrypts ciphertext with priv. If opts is nil or of type
// *PKCS1v15DecryptOptions then PKCS#1 v1.5 decryption is performed. Otherwise
// opts must have type *OAEPOptions and OAEP decryption is done.
func (priv *PrivateKey) Decrypt(rand io.Reader, ciphertext []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error)

// Precompute performs some calculations that speed up private key operations
// in the future.

// Precompute方法会预先进行一些计算，以加速未来的私钥的操作。
func (priv *PrivateKey) Precompute()

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() crypto.PublicKey

// Sign signs msg with priv, reading randomness from rand. If opts is a
// *PSSOptions then the PSS algorithm will be used, otherwise PKCS#1 v1.5 will
// be used. This method is intended to support keys where the private part is
// kept in, for example, a hardware module. Common uses should use the Sign*
// functions in this package.
func (priv *PrivateKey) Sign(rand io.Reader, msg []byte, opts crypto.SignerOpts) ([]byte, error)

// Validate performs basic sanity checks on the key. It returns nil if the key
// is valid, or else an error describing a problem.

// Validate方法进行密钥的完整性检查。如果密钥合法会返回nil，否则会返回说明问题的
// error值。
func (priv *PrivateKey) Validate() error

