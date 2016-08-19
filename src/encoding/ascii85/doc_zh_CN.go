// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package ascii85 implements the ascii85 data encoding
// as used in the btoa tool and Adobe's PostScript and PDF document formats.

// ascii85 包是对 ascii85 的数据编码的实现.
// 被用于 btoa ( binary to ascii )工具， Adobe 的 PostScript 和PDF文档格式。
package ascii85

import (
    "io"
    "strconv"
)


type CorruptInputError int64


// Decode decodes src into dst, returning both the number
// of bytes written to dst and the number consumed from src.
// If src contains invalid ascii85 data, Decode will return the
// number of bytes successfully written and a CorruptInputError.
// Decode ignores space and control characters in src.
// Often, ascii85-encoded data is wrapped in <~ and ~> symbols.
// Decode expects these to have been stripped by the caller.
//
// If flush is true, Decode assumes that src represents the
// end of the input stream and processes it completely rather
// than wait for the completion of another 32-bit block.
//
// NewDecoder wraps an io.Reader interface around Decode.

// Decode 从源解码到目标，返回写入目标和源消耗的字节数. 若源包含无效 ascii85 数
// 据， Decode 将返回成功写入的字节数和 CorruptInputError 函数。 Decode 忽略源中
// 的空格和控制字符。 通常，ascii85 编码数据用 <~ 和 ~> 符号括起来。 Decode 期望
// 这些被调用器去除。
//
// 若 flush 为真， Decode 会假定源表现为输入流结束并立即处理，而不是等待另一个32
// 位块的结束。
//
// NewDecoder 包含一个 io.Reader 接口，区别于 Decode 。
func Decode(dst, src []byte, flush bool) (ndst, nsrc int, err error)

// Encode encodes src into at most MaxEncodedLen(len(src))
// bytes of dst, returning the actual number of bytes written.
//
// The encoding handles 4-byte chunks, using a special encoding
// for the last fragment, so Encode is not appropriate for use on
// individual blocks of a large data stream.  Use NewEncoder() instead.
//
// Often, ascii85-encoded data is wrapped in <~ and ~> symbols.
// Encode does not add these.

// Encode 编码源的最多 MaxEncodedLen(len(src)) 字节的到目标，
// 返回实际的写入字节数。
//
// Encode 通过对最后分段使用特殊的编码来操作4字节的数据块，
// 所以将它用在大型数据流的私有块上是不合适的。请用 NewEncoder() 替代。
//
// 通常， ascii85 编码的数据用符号 <~ 和 ~> 括起来。
// Encode 不加这些。
func Encode(dst, src []byte) int

// MaxEncodedLen returns the maximum length of an encoding of n source bytes.

// MaxEncodedLen 返回 n 源字节编码的最大长度.
func MaxEncodedLen(n int) int

// NewDecoder constructs a new ascii85 stream decoder.

// NewDecoder 构造一个新的 ascii85 流解码器.
func NewDecoder(r io.Reader) io.Reader

// NewEncoder returns a new ascii85 stream encoder.  Data written to
// the returned writer will be encoded and then written to w.
// Ascii85 encodings operate in 32-bit blocks; when finished
// writing, the caller must Close the returned encoder to flush any
// trailing partial block.

// NewEncoder 返回一个新的 ascii85 流编码器. 写入到返回的写入器中的数据将被编
// 码，然后写入到 w 中。 Ascii85 编码在32位块中操作；当完成写入时，调用者必须关
// 闭返回的编码器，去除所有尾部块。
func NewEncoder(w io.Writer) io.WriteCloser

func (CorruptInputError) Error() string

