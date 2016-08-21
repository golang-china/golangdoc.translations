// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package pem implements the PEM data encoding, which originated in Privacy
// Enhanced Mail. The most common use of PEM encoding today is in TLS keys and
// certificates. See RFC 1421.

// pem包实现了PEM数据编码（源自保密增强邮件协议）。目前PEM编码主要用于TLS密钥和
// 证书。参见RFC 1421
package pem

import (
    "bytes"
    "encoding/base64"
    "errors"
    "io"
    "sort"
    "strings"
)

// A Block represents a PEM encoded structure.
//
// The encoded form is:
//    -----BEGIN Type-----
//    Headers
//    base64-encoded Bytes
//    -----END Type-----
// where Headers is a possibly empty sequence of Key: Value lines.

// Block代表PEM编码的结构。编码格式如下：
//
//     -----BEGIN Type-----
//     Headers
//     base64-encoded Bytes
//     -----END Type-----
//
// 其中Headers是可为空的多行键值对。
type Block struct {
    Type    string            // The type, taken from the preamble (i.e. "RSA PRIVATE KEY").
    Headers map[string]string // Optional headers.
    Bytes   []byte            // The decoded bytes of the contents. Typically a DER encoded ASN.1 structure.
}

// Decode will find the next PEM formatted block (certificate, private key
// etc) in the input. It returns that block and the remainder of the input. If
// no PEM data is found, p is nil and the whole of the input is returned in
// rest.

// Decode函数会从输入里查找到下一个PEM格式的块（证书、私钥等）。它返回解码得到的
// Block和剩余未解码的数据。如果未发现PEM数据，返回(nil, data)。
func Decode(data []byte) (p *Block, rest []byte)

func Encode(out io.Writer, b *Block) error

func EncodeToMemory(b *Block) []byte

