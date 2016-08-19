// Copyright 2009 The Go Authors. All rights reserved.
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


// LocalhostCert is a PEM-encoded TLS cert with SAN IPs "127.0.0.1" and "[::1]",
// expiring at Jan 29 16:00:00 2084 GMT. generated from src/crypto/tls: go run
// generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com --ca
// --start-date "Jan 1 00:00:00 1970" --duration=1000000h

// LocalhostCert is a PEM-encoded TLS cert with SAN IPs "127.0.0.1" and "[::1]",
// expiring at Jan 29 16:00:00 2084 GMT. generated from src/crypto/tls: go run
// generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com --ca
// --start-date "Jan 1 00:00:00 1970" --duration=1000000h
var LocalhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIICEzCCAXygAwIBAgIQMIMChMLGrR+QvmQvpwAU6zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9SjY1bIw4
iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZBl2+XsDul
rKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQABo2gwZjAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAuBgNVHREEJzAlggtleGFtcGxlLmNvbYcEfwAAAYcQAAAAAAAAAAAAAAAA
AAAAATANBgkqhkiG9w0BAQsFAAOBgQCEcetwO59EWk7WiJsG4x8SY+UIAA+flUI9
tyC4lNhbcF2Idq9greZwbYCqTTTr2XiRNSMLCOjKyI7ukPoPjo16ocHj+P3vZGfs
h1fIw3cSS2OolhloGw/XM6RWPWtPAlGykKLciQrBru5NAPvCMsb/I1DAceTiotQM
fblo6RBxUQ==
-----END CERTIFICATE-----`)


// LocalhostKey is the private key for localhostCert.
var LocalhostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9
SjY1bIw4iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZB
l2+XsDulrKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQAB
AoGAGRzwwir7XvBOAy5tM/uV6e+Zf6anZzus1s1Y1ClbjbE6HXbnWWF/wbZGOpet
3Zm4vD6MXc7jpTLryzTQIvVdfQbRc6+MUVeLKwZatTXtdZrhu+Jk7hx0nTPy8Jcb
uJqFk541aEw+mMogY/xEcfbWd6IOkp+4xqjlFLBEDytgbIECQQDvH/E6nk+hgN4H
qzzVtxxr397vWrjrIgPbJpQvBsafG7b0dA4AFjwVbFLmQcj2PprIMmPcQrooz8vp
jy4SHEg1AkEA/v13/5M47K9vCxmb8QeD/asydfsgS5TeuNi8DoUBEmiSJwma7FXY
fFUtxuvL7XvjwjN5B30pNEbc6Iuyt7y4MQJBAIt21su4b3sjXNueLKH85Q+phy2U
fQtuUE9txblTu14q3N7gHRZB4ZMhFYyDy8CKrN2cPg/Fvyt0Xlp/DoCzjA0CQQDU
y2ptGsuSmgUtWj3NM9xuwYPm+Z/F84K6+ARYiZ6PYj013sovGKUFfYAqVXVlxtIX
qyUBnu3X9ps8ZfjLZO7BAkEAlT4R5Yl6cGhaJQYZHOde3JEMhNRcVFMO8dJDaFeo
f9Oeos0UUothgiDktdQHxdNEwLjQf7lJJBzV+5OtwswCWA==
-----END RSA PRIVATE KEY-----`)


// FlushAfterChunkWriter signals from the caller of NewChunkedWriter
// that each chunk should be followed by a flush. It is used by the
// http.Transport code to keep the buffering behavior for headers and
// trailers, but flush out chunks aggressively in the middle for
// request bodies which may be generated slowly. See Issue 6574.
type FlushAfterChunkWriter struct {
}


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

