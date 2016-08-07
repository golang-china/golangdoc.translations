// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package hmac implements the Keyed-Hash Message Authentication Code (HMAC) as
// defined in U.S. Federal Information Processing Standards Publication 198.
// An HMAC is a cryptographic hash that uses a key to sign a message.
// The receiver verifies the hash by recomputing it using the same key.
//
// Receivers should be careful to use Equal to compare MACs in order to avoid
// timing side-channels:
//
//     // CheckMAC reports whether messageMAC is a valid HMAC tag for message.
//     func CheckMAC(message, messageMAC, key []byte) bool {
//         mac := hmac.New(sha256.New, key)
//         mac.Write(message)
//         expectedMAC := mac.Sum(nil)
//         return hmac.Equal(messageMAC, expectedMAC)
//     }

// hmac包实现了U.S. Federal Information Processing Standards Publication
// 198规定的HMAC（加密哈希信息认证码）。
//
// HMAC是使用key标记信息的加密hash。接收者使用相同的key逆运算来认证hash。
//
// 出于安全目的，接收者应使用Equal函数比较认证码：
//
//     // 如果messageMAC是message的合法HMAC标签，函数返回真
//     func CheckMAC(message, messageMAC, key []byte) bool {
//         mac := hmac.New(sha256.New, key)
//         mac.Write(message)
//         expectedMAC := mac.Sum(nil)
//         return hmac.Equal(messageMAC, expectedMAC)
//     }
package hmac

import (
    "crypto/subtle"
    "hash"
)

// Equal compares two MACs for equality without leaking timing information.

// 比较两个MAC是否相同，而不会泄露对比时间信息。（以规避时间侧信道攻击：指通过计
// 算比较时花费的时间的长短来获取密码的信息，用于密码破解）
func Equal(mac1, mac2 []byte) bool

// New returns a new HMAC hash using the given hash.Hash type and key.

// New函数返回一个采用hash.Hash作为底层hash接口、key作为密钥的HMAC算法的hash接口
// 。
func New(h func() hash.Hash, key []byte) hash.Hash

