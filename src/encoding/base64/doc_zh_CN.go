// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package base64 implements base64 encoding as specified by RFC 4648.

// base64实现了RFC 4648规定的base64编码。
package base64

import (
	"io"
	"strconv"
)

// StdEncoding is the standard base64 encoding, as defined in
// RFC 4648.

// RFC 4648定义的标准base64编码字符集。
var StdEncoding = NewEncoding(encodeStd)

// URLEncoding is the alternate base64 encoding defined in RFC 4648. It is
// typically used in URLs and file names.

// RFC 4648 定义的另一base64编码字符集，用于URL和文件名。
var URLEncoding = NewEncoding(encodeURL)

type CorruptInputError int64

// An Encoding is a radix 64 encoding/decoding scheme, defined by a
// 64-character alphabet.  The most common encoding is the "base64"
// encoding defined in RFC 4648 and used in MIME (RFC 2045) and PEM
// (RFC 1421).  RFC 4648 also defines an alternate encoding, which is
// the standard encoding with - and _ substituted for + and /.

// 双向的编码/解码协议，根据一个64字符的字符集定义，RFC
// 4648标准化了两种字符集。默认字符集用于MIME（RFC 2045）和PEM（RFC
// 1421）编码；另一种用于URL和文件名，用'-'和'_'替换了'+'和'/'。
type Encoding struct {
}

// NewDecoder constructs a new base64 stream decoder.

// 创建一个新的base64流解码器。
func NewDecoder(enc *Encoding, r io.Reader) io.Reader

// NewEncoder returns a new base64 stream encoder.  Data written to
// the returned writer will be encoded using enc and then written to w.
// Base64 encodings operate in 4-byte blocks; when finished
// writing, the caller must Close the returned encoder to flush any
// partially written blocks.

// 创建一个新的base64流编码器。写入的数据会在编码后再写入w，base32编码每3字节执
// 行一次编码操作；写入完毕后，使用者必须调用Close方法以便将未写入的数据从缓存中
// 刷新到w中。
func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser

// NewEncoding returns a new padded Encoding defined by the given alphabet,
// which must be a 64-byte string.
// The resulting Encoding uses the default padding character ('='),
// which may be changed or disabled via WithPadding.

// 使用给出的字符集生成一个*Encoding，字符集必须是64字节的字符串。
func NewEncoding(encoder string) *Encoding

// Decode decodes src using the encoding enc.  It writes at most
// DecodedLen(len(src)) bytes to dst and returns the number of bytes
// written.  If src contains invalid base64 data, it will return the
// number of bytes successfully written and CorruptInputError.
// New line characters (\r and \n) are ignored.

// 将src的数据解码后存入dst，最多写DecodedLen(len(src))字节数据到dst，并返回写入
// 的字节数。 如果src包含非法字符，将返回成功写入的字符数和CorruptInputError。换
// 行符（\r、\n）会被忽略。
func (*Encoding) Decode(dst, src []byte) (n int, err error)

// DecodeString returns the bytes represented by the base64 string s.

// 返回base64编码的字符串s代表的数据。
func (*Encoding) DecodeString(s string) ([]byte, error)

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base64-encoded data.

// 返回n字节base64编码的数据解码后的最大长度。
func (*Encoding) DecodedLen(n int) int

// Encode encodes src using the encoding enc, writing
// EncodedLen(len(src)) bytes to dst.
//
// The encoding pads the output to a multiple of 4 bytes,
// so Encode is not appropriate for use on individual blocks
// of a large data stream.  Use NewEncoder() instead.

// Encode encodes src using the encoding enc, writing EncodedLen(len(src)) bytes
// to dst.
//
// The encoding pads the output to a multiple of 4 bytes, so Encode is not
// appropriate for use on individual blocks of a large data stream. Use
// NewEncoder() instead.
func (*Encoding) Encode(dst, src []byte)

// EncodeToString returns the base64 encoding of src.
func (*Encoding) EncodeToString(src []byte) string

// EncodedLen returns the length in bytes of the base64 encoding of an input
// buffer of length n.
func (*Encoding) EncodedLen(n int) int

func (CorruptInputError) Error() string
