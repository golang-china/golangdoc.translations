// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package base32 implements base32 encoding as specified by RFC 4648.

// Package base32 implements base32
// encoding as specified by RFC 4648.
package base32

// HexEncoding is the ``Extended Hex Alphabet'' defined in RFC 4648. It is
// typically used in DNS.

// HexEncoding is the ``Extended Hex
// Alphabet'' defined in RFC 4648. It is
// typically used in DNS.
var HexEncoding = NewEncoding(encodeHex)

// StdEncoding is the standard base32 encoding, as defined in RFC 4648.

// StdEncoding is the standard base32
// encoding, as defined in RFC 4648.
var StdEncoding = NewEncoding(encodeStd)

// NewDecoder constructs a new base32 stream decoder.

// NewDecoder constructs a new base32
// stream decoder.
func NewDecoder(enc *Encoding, r io.Reader) io.Reader

// NewEncoder returns a new base32 stream encoder. Data written to the returned
// writer will be encoded using enc and then written to w. Base32 encodings operate
// in 5-byte blocks; when finished writing, the caller must Close the returned
// encoder to flush any partially written blocks.

// NewEncoder returns a new base32 stream
// encoder. Data written to the returned
// writer will be encoded using enc and
// then written to w. Base32 encodings
// operate in 5-byte blocks; when finished
// writing, the caller must Close the
// returned encoder to flush any partially
// written blocks.
func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser

type CorruptInputError int64

func (e CorruptInputError) Error() string

// An Encoding is a radix 32 encoding/decoding scheme, defined by a 32-character
// alphabet. The most common is the "base32" encoding introduced for SASL GSSAPI
// and standardized in RFC 4648. The alternate "base32hex" encoding is used in
// DNSSEC.

// An Encoding is a radix 32
// encoding/decoding scheme, defined by a
// 32-character alphabet. The most common
// is the "base32" encoding introduced for
// SASL GSSAPI and standardized in RFC
// 4648. The alternate "base32hex" encoding
// is used in DNSSEC.
type Encoding struct {
	// contains filtered or unexported fields
}

// NewEncoding returns a new Encoding defined by the given alphabet, which must be
// a 32-byte string.

// NewEncoding returns a new Encoding
// defined by the given alphabet, which
// must be a 32-byte string.
func NewEncoding(encoder string) *Encoding

// Decode decodes src using the encoding enc. It writes at most
// DecodedLen(len(src)) bytes to dst and returns the number of bytes written. If
// src contains invalid base32 data, it will return the number of bytes
// successfully written and CorruptInputError. New line characters (\r and \n) are
// ignored.

// Decode decodes src using the encoding
// enc. It writes at most
// DecodedLen(len(src)) bytes to dst and
// returns the number of bytes written. If
// src contains invalid base32 data, it
// will return the number of bytes
// successfully written and
// CorruptInputError. New line characters
// (\r and \n) are ignored.
func (enc *Encoding) Decode(dst, src []byte) (n int, err error)

// DecodeString returns the bytes represented by the base32 string s.

// DecodeString returns the bytes
// represented by the base32 string s.
func (enc *Encoding) DecodeString(s string) ([]byte, error)

// DecodedLen returns the maximum length in bytes of the decoded data corresponding
// to n bytes of base32-encoded data.

// DecodedLen returns the maximum length in
// bytes of the decoded data corresponding
// to n bytes of base32-encoded data.
func (enc *Encoding) DecodedLen(n int) int

// Encode encodes src using the encoding enc, writing EncodedLen(len(src)) bytes to
// dst.
//
// The encoding pads the output to a multiple of 8 bytes, so Encode is not
// appropriate for use on individual blocks of a large data stream. Use
// NewEncoder() instead.

// Encode encodes src using the encoding
// enc, writing EncodedLen(len(src)) bytes
// to dst.
//
// The encoding pads the output to a
// multiple of 8 bytes, so Encode is not
// appropriate for use on individual blocks
// of a large data stream. Use NewEncoder()
// instead.
func (enc *Encoding) Encode(dst, src []byte)

// EncodeToString returns the base32 encoding of src.

// EncodeToString returns the base32
// encoding of src.
func (enc *Encoding) EncodeToString(src []byte) string

// EncodedLen returns the length in bytes of the base32 encoding of an input buffer
// of length n.

// EncodedLen returns the length in bytes
// of the base32 encoding of an input
// buffer of length n.
func (enc *Encoding) EncodedLen(n int) int
