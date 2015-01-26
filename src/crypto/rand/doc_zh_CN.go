// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package rand implements a cryptographically secure pseudorandom number
// generator.

// Package rand implements a
// cryptographically secure pseudorandom
// number generator.
package rand

// Reader is a global, shared instance of a cryptographically strong pseudo-random
// generator. On Unix-like systems, Reader reads from /dev/urandom. On Windows
// systems, Reader uses the CryptGenRandom API.

// Reader is a global, shared instance of a
// cryptographically strong pseudo-random
// generator. On Unix-like systems, Reader
// reads from /dev/urandom. On Windows
// systems, Reader uses the CryptGenRandom
// API.
var Reader io.Reader

// Int returns a uniform random value in [0, max). It panics if max <= 0.

// Int returns a uniform random value in
// [0, max). It panics if max <= 0.
func Int(rand io.Reader, max *big.Int) (n *big.Int, err error)

// Prime returns a number, p, of the given size, such that p is prime with high
// probability. Prime will return error for any error returned by rand.Read or if
// bits < 2.

// Prime returns a number, p, of the given
// size, such that p is prime with high
// probability. Prime will return error for
// any error returned by rand.Read or if
// bits < 2.
func Prime(rand io.Reader, bits int) (p *big.Int, err error)

// Read is a helper function that calls Reader.Read using io.ReadFull. On return, n
// == len(b) if and only if err == nil.

// Read is a helper function that calls
// Reader.Read using io.ReadFull. On
// return, n == len(b) if and only if err
// == nil.
func Read(b []byte) (n int, err error)
