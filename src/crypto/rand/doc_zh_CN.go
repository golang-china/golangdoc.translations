// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package rand implements a cryptographically secure pseudorandom number
// generator.

// rand包实现了用于加解密的更安全的随机数生成器。
package rand

// Reader is a global, shared instance of a cryptographically strong pseudo-random
// generator. On Unix-like systems, Reader reads from /dev/urandom. On Windows
// systems, Reader uses the CryptGenRandom API.

// Reader是一个全局、共享的密码用强随机数生成器。在Unix类型系统中，会从/dev/urandom读取；而Windows中会调用CryptGenRandom
// API。
var Reader io.Reader

// Int returns a uniform random value in [0, max). It panics if max <= 0.

// 返回一个在[0,
// max)区间服从均匀分布的随机值，如果max<=0则会panic。
func Int(rand io.Reader, max *big.Int) (n *big.Int, err error)

// Prime returns a number, p, of the given size, such that p is prime with high
// probability. Prime will return error for any error returned by rand.Read or if
// bits < 2.

// 返回一个具有指定字位数的数字，该数字具有很高可能性是质数。如果从rand读取时出错，或者bits<2会返回错误。
func Prime(rand io.Reader, bits int) (p *big.Int, err error)

// Read is a helper function that calls Reader.Read using io.ReadFull. On return, n
// == len(b) if and only if err == nil.

// 本函数是一个使用io.ReadFull调用Reader.Read的辅助性函数。当且仅当err == nil时，返回值n == len(b)。
func Read(b []byte) (n int, err error)
