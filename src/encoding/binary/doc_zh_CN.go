// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package binary implements simple translation between numbers and byte
// sequences and encoding and decoding of varints.
//
// Numbers are translated by reading and writing fixed-size values.
// A fixed-size value is either a fixed-size arithmetic
// type (int8, uint8, int16, float32, complex64, ...)
// or an array or struct containing only fixed-size values.
//
// The varint functions encode and decode single integer values using
// a variable-length encoding; smaller values require fewer bytes.
// For a specification, see
// https://developers.google.com/protocol-buffers/docs/encoding.
//
// This package favors simplicity over efficiency. Clients that require
// high-performance serialization, especially for large data structures,
// should look at more advanced solutions such as the encoding/gob
// package or protocol buffers.

// binary包实现了简单的数字与字节序列的转换以及变长值的编解码。
//
// 数字翻译为定长值来读写，一个定长值，要么是固定长度的数字类型（int8, uint8,
// int16, float32, complex64, ...）或者只包含定长值的结构体或者数组。
//
// 变长值是使用一到多个字节编码整数的方法，绝对值较小的数字会占用较少的字节数。
// 详情请参见：http://code.google.com/apis/protocolbuffers/docs/encoding.html。
//
// 本包相对于效率更注重简单。如果需要高效的序列化，特别是数据结构较复杂的，请参
// 见更高级的解决方法，例如encoding/gob包，或者采用协议缓存。
package binary

import (
	"errors"
	"io"
	"math"
	"reflect"
)

// MaxVarintLenN is the maximum length of a varint-encoded N-bit integer.

// 变长编码N位整数的最大字节数。
const (
	MaxVarintLen16 = 3
	MaxVarintLen32 = 5
	MaxVarintLen64 = 10
)

// BigEndian is the big-endian implementation of ByteOrder.
var BigEndian bigEndian

// LittleEndian is the little-endian implementation of ByteOrder.
var LittleEndian littleEndian

// A ByteOrder specifies how to convert byte sequences into
// 16-, 32-, or 64-bit unsigned integers.

// ByteOrder规定了如何将字节序列和
// 16、32或64比特的无符号整数互相转化。
type ByteOrder interface {
	Uint16([]byte)uint16
	Uint32([]byte)uint32
	Uint64([]byte)uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String()string
}

// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.

// 将一个uint64数字编码写入buf并返回写入的长度，如果buf太小，则会panic。
func PutUvarint(buf []byte, x uint64) int

// PutVarint encodes an int64 into buf and returns the number of bytes written.
// If the buffer is too small, PutVarint will panic.

// 将一个int64数字编码写入buf并返回写入的长度，如果buf太小，则会panic。
func PutVarint(buf []byte, x int64) int

// Read reads structured binary data from r into data.
// Data must be a pointer to a fixed-size value or a slice
// of fixed-size values.
// Bytes read from r are decoded using the specified byte order
// and written to successive fields of the data.
// When reading into structs, the field data for fields with
// blank (_) field names is skipped; i.e., blank field names
// may be used for padding.
// When reading into a struct, all non-blank fields must be exported.
//
// The error is EOF only if no bytes were read.
// If an EOF happens after reading some but not all the bytes,
// Read returns ErrUnexpectedEOF.

// 从r中读取binary编码的数据并赋给data，data必须是一个指向定长值的指针或者定长值
// 的切片。从r读取的字节使用order指定的字节序解码并写入data的字段里当写入结构体
// 是，名字中有'_'的字段会被跳过，这些字段可用于填充（内存空间）。
func Read(r io.Reader, order ByteOrder, data interface{}) error

// ReadUvarint reads an encoded unsigned integer from r and returns it as a
// uint64.

// 从r读取一个编码后的无符号整数，并返回该整数。
func ReadUvarint(r io.ByteReader) (uint64, error)

// ReadVarint reads an encoded signed integer from r and returns it as an int64.

// 从r读取一个编码后的有符号整数，并返回该整数。
func ReadVarint(r io.ByteReader) (int64, error)

// Size returns how many bytes Write would generate to encode the value v, which
// must be a fixed-size value or a slice of fixed-size values, or a pointer to
// such data. If v is neither of these, Size returns -1.

// 返回v编码后会占用多少字节，注意v必须是定长值、定长值的切片、定长值的指针。
func Size(v interface{}) int

// Uvarint decodes a uint64 from buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 meaning:
//
// 	n == 0: buf too small
// 	n  < 0: value larger than 64 bits (overflow)
//              and -n is the number of bytes read

// 从buf解码一个uint64，返回该数字和读取的字节长度，如果发生了错误，该数字为0而
// 读取长度n返回值的意思是：
//
//     n == 0: buf不完整，太短了
//     n  < 0: 值太大了，64比特装不下（溢出），-n为读取的字节数
func Uvarint(buf []byte) (uint64, int)

// Varint decodes an int64 from buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 with the following meaning:
//
// 	n == 0: buf too small
// 	n  < 0: value larger than 64 bits (overflow)
//              and -n is the number of bytes read

// 从buf解码一个int64，返回该数字和读取的字节长度，如果发生了错误，该数字为0而读
// 取长度n返回值的意思是：
//
//     n == 0: buf不完整，太短了
//     n  < 0: 值太大了，64比特装不下（溢出），-n为读取的字节数
func Varint(buf []byte) (int64, int)

// Write writes the binary representation of data into w.
// Data must be a fixed-size value or a slice of fixed-size
// values, or a pointer to such data.
// Bytes written to w are encoded using the specified byte order
// and read from successive fields of the data.
// When writing structs, zero values are written for fields
// with blank (_) field names.

// 将data的binary编码格式写入w，data必须是定长值、定长值的切片、定长值的指针。
// order指定写入数据的字节序，写入结构体时，名字中有'_'的字段会置为0。
func Write(w io.Writer, order ByteOrder, data interface{}) error

