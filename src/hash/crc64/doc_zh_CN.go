// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package crc64 implements the 64-bit cyclic redundancy check, or CRC-64,
// checksum. See http://en.wikipedia.org/wiki/Cyclic_redundancy_check for
// information.

// Package crc64 implements the 64-bit cyclic redundancy check, or CRC-64,
//
//	checksum. See http://en.wikipedia.org/wiki/Cyclic_redundancy_check for
//	information.
package crc64

// Predefined polynomials.

// 预定义的多项式。
//
//	const Size = 8
//
// CRC-64校验和的字节数。
const (
	// The ISO polynomial, defined in ISO 3309 and used in HDLC.
	ISO = 0xD800000000000000

	// The ECMA polynomial, defined in ECMA 182.
	ECMA = 0xC96C5795D7870F42
)

// The size of a CRC-64 checksum in bytes.
const Size = 8

// Checksum returns the CRC-64 checksum of data using the polynomial represented by
// the Table.

// 返回数据data使用tab代表的多项式计算出的CRC-64校验和。
func Checksum(data []byte, tab *Table) uint64

// New creates a new hash.Hash64 computing the CRC-64 checksum using the polynomial
// represented by the Table.

// 创建一个使用tab代表的多项式计算CRC-64校验和的hash.Hash64接口。
func New(tab *Table) hash.Hash64

// Update returns the result of adding the bytes in p to the crc.

// 返回将切片p的数据采用tab表示的多项式添加到crc之后计算出的新校验和。
func Update(crc uint64, tab *Table, p []byte) uint64

// Table is a 256-word table representing the polynomial for efficient processing.

// 长度256的uint64切片，代表一个用于高效运作的多项式。
type Table [256]uint64

// MakeTable returns the Table constructed from the specified polynomial.

// 返回一个代表poly指定的多项式的*Table。
func MakeTable(poly uint64) *Table
