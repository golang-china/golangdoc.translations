// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package ecdsa implements the Elliptic Curve Digital Signature Algorithm, as
// defined in FIPS 186-3.
package ecdsa

// Sign signs an arbitrary length hash (which should be the result of hashing a
// larger message) using the private key, priv. It returns the signature as a pair
// of integers. The security of the private key depends on the entropy of rand.
func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error)

// Verify verifies the signature in r, s of hash using the public key, pub. Its
// return value records whether the signature is valid.
func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool

// PrivateKey represents a ECDSA private key.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

// GenerateKey generates a public and private key pair.
func GenerateKey(c elliptic.Curve, rand io.Reader) (priv *PrivateKey, err error)

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() crypto.PublicKey

// Sign signs msg with priv, reading randomness from rand. This method is intended
// to support keys where the private part is kept in, for example, a hardware
// module. Common uses should use the Sign function in this package directly.
func (priv *PrivateKey) Sign(rand io.Reader, msg []byte, opts crypto.SignerOpts) ([]byte, error)

// PublicKey represents an ECDSA public key.
type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}
