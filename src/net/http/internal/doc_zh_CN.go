// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package internal contains HTTP internals shared by net/http and
// net/http/httputil.

// internal 包含 net/http 和 net/http/httputil 共享的 HTTP 内部函数.
package internal

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "io"
)

var ErrLineTooLong = errors.New("header line too long")

// NewChunkedReader returns a new chunkedReader that translates the data read
// from r out of HTTP "chunked" format before returning it. The chunkedReader
// returns io.EOF when the final 0-length chunk is read.
//
// NewChunkedReader is not needed by normal applications. The http package
// automatically decodes chunking when reading response bodies.

// NewChunkedReader返回一个新的chunkedReader。这个chunkedReader能翻译从r中的HTTP
// “chunked” 获取到的数据，并且返回数据。 chunkedReader当读取到最后的0长度的
// chunk的时候返回io.EOF。
//
// NewChunkedReader在通常的应用中并不需要。http包会在读取回复的消息体的时候自动
// 解码。
func NewChunkedReader(r io.Reader) io.Reader

// NewChunkedWriter returns a new chunkedWriter that translates writes into HTTP
// "chunked" format before writing them to w. Closing the returned chunkedWriter
// sends the final 0-length chunk that marks the end of the stream.
//
// NewChunkedWriter is not needed by normal applications. The http
// package adds chunking automatically if handlers don't set a
// Content-Length header. Using newChunkedWriter inside a handler
// would result in double chunking or chunking with a Content-Length
// length, both of which are wrong.

// NewChunkedWriter返回一个新的chunkWriter，这个chunkWriter会将HTTP的“chunked”
// 格式进行转化 然后写入w中。关闭返回的chunkedWriter，发送0长度的chunk就标志着流
// 的结束。
//
// 一般的应用并不需要使用NewChunkedWriter。如果不设置Content-Length头的话，http
// 包会自动增加chunk。 在handler中使用NewChunkedWriter会导致重复块，或者有
// Content-length长度的块，这两种都是错误的。
func NewChunkedWriter(w io.Writer) io.WriteCloser

