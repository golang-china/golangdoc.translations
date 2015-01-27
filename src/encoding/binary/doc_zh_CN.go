// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package binary implements simple translation between numbers and byte sequences
// and encoding and decoding of varints.
//
// Numbers are translated by reading and writing fixed-size values. A fixed-size
// value is either a fixed-size arithmetic type (int8, uint8, int16, float32,
// complex64, ...) or an array or struct containing only fixed-size values.
//
// The varint functions encode and decode single integer values using a
// variable-length encoding; smaller values require fewer bytes. For a
// specification, see
// http://code.google.com/apis/protocolbuffers/docs/encoding.html.
//
// This package favors simplicity over efficiency. Clients that require
// high-performance serialization, especially for large data structures, should
// look at more advanced solutions such as the encoding/gob package or protocol
// buffers.
package binary

// MaxVarintLenN is the maximum length of a varint-encoded N-bit integer.
const (
	MaxVarintLen16 = 3
	MaxVarintLen32 = 5
	MaxVarintLen64 = 10
)

// BigEndian is the big-endian implementation of ByteOrder.
var BigEndian bigEndian

// LittleEndian is the little-endian implementation of ByteOrder.
var LittleEndian littleEndian

// PutUvarint encodes a uint64 into buf and returns the number of bytes written. If
// the buffer is too small, PutUvarint will panic.
func PutUvarint(buf []byte, x uint64) int

// PutVarint encodes an int64 into buf and returns the number of bytes written. If
// the buffer is too small, PutVarint will panic.
func PutVarint(buf []byte, x int64) int

// Read reads structured binary data from r into data. Data must be a pointer to a
// fixed-size value or a slice of fixed-size values. Bytes read from r are decoded
// using the specified byte order and written to successive fields of the data.
// When reading into structs, the field data for fields with blank (_) field names
// is skipped; i.e., blank field names may be used for padding. When reading into a
// struct, all non-blank fields must be exported.
func Read(r io.Reader, order ByteOrder, data interface{}) error

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func ReadUvarint(r io.ByteReader) (uint64, error)

// ReadVarint reads an encoded signed integer from r and returns it as an int64.
func ReadVarint(r io.ByteReader) (int64, error)

// Size returns how many bytes Write would generate to encode the value v, which
// must be a fixed-size value or a slice of fixed-size values, or a pointer to such
// data. If v is neither of these, Size returns -1.
func Size(v interface{}) int

// Uvarint decodes a uint64 from buf and returns that value and the number of bytes
// read (> 0). If an error occurred, the value is 0 and the number of bytes n is <=
// 0 meaning:
//
//		n == 0: buf too small
//		n  < 0: value larger than 64 bits (overflow)
//	             and -n is the number of bytes read
func Uvarint(buf []byte) (uint64, int)

// Varint decodes an int64 from buf and returns that value and the number of bytes
// read (> 0). If an error occurred, the value is 0 and the number of bytes n is <=
// 0 with the following meaning:
//
//		n == 0: buf too small
//		n  < 0: value larger than 64 bits (overflow)
//	             and -n is the number of bytes read
func Varint(buf []byte) (int64, int)

// Write writes the binary representation of data into w. Data must be a fixed-size
// value or a slice of fixed-size values, or a pointer to such data. Bytes written
// to w are encoded using the specified byte order and read from successive fields
// of the data. When writing structs, zero values are written for fields with blank
// (_) field names.
func Write(w io.Writer, order ByteOrder, data interface{}) error

// A ByteOrder specifies how to convert byte sequences into 16-, 32-, or 64-bit
// unsigned integers.
type ByteOrder interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String() string
}
