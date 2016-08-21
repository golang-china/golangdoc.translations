// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package crc32 implements the 32-bit cyclic redundancy check, or CRC-32,
// checksum. See http://en.wikipedia.org/wiki/Cyclic_redundancy_check for
// information.
//
// Polynomials are represented in LSB-first form also known as reversed
// representation.
//
// See
//
//
// http://en.wikipedia.org/wiki/Mathematics_of_cyclic_redundancy_checks#Reversed_representations_and_reciprocal_polynomials for information.

// crc32包实现了32位循环冗余校验（CRC-32）的校验和算法，参见：
//
// http://en.wikipedia.org/wiki/Cyclic_redundancy_check
package crc32

import (
    "hash"
    "sync"
)

// Predefined polynomials.

// 预定义的多项式。
//
//     const Size = 4
//
// CRC-32校验和的字节长度。
const (
    // IEEE is by far and away the most common CRC-32 polynomial.
    // Used by ethernet (IEEE 802.3), v.42, fddi, gzip, zip, png, ...
    IEEE = 0xedb88320

    // Castagnoli's polynomial, used in iSCSI.
    // Has better error detection characteristics than IEEE.
    // http://dx.doi.org/10.1109/26.231911
    Castagnoli = 0x82f63b78

    // Koopman's polynomial.
    // Also has better error detection characteristics than IEEE.
    // http://dx.doi.org/10.1109/DSN.2002.1028931
    Koopman = 0xeb31d82e
)

// The size of a CRC-32 checksum in bytes.
const Size = 4

// IEEETable is the table for the IEEE polynomial.

// IEEETable是IEEE多项式对应的Table。
var IEEETable = makeTable(IEEE)

// Table is a 256-word table representing the polynomial for efficient
// processing.

// 长度256的uint32切片，代表一个用于高效运作的多项式。
type Table [256]uint32

// Checksum returns the CRC-32 checksum of data
// using the polynomial represented by the Table.

// 返回数据data使用tab代表的多项式计算出的CRC-32校验和。
func Checksum(data []byte, tab *Table) uint32

// ChecksumIEEE returns the CRC-32 checksum of data
// using the IEEE polynomial.

// 返回数据data使用IEEE多项式计算出的CRC-32校验和。
func ChecksumIEEE(data []byte) uint32

// MakeTable returns a Table constructed from the specified polynomial.
// The contents of this Table must not be modified.

// 返回一个代表poly指定的多项式的Table。
func MakeTable(poly uint32) *Table

// New creates a new hash.Hash32 computing the CRC-32 checksum
// using the polynomial represented by the Table.
// Its Sum method will lay the value out in big-endian byte order.

// 创建一个使用tab代表的多项式计算CRC-32校验和的hash.Hash32接口。
func New(tab *Table) hash.Hash32

// NewIEEE creates a new hash.Hash32 computing the CRC-32 checksum
// using the IEEE polynomial.
// Its Sum method will lay the value out in big-endian byte order.

// 创建一个使用IEEE多项式计算CRC-32校验和的hash.Hash32接口。
func NewIEEE() hash.Hash32

// Update returns the result of adding the bytes in p to the crc.

// 返回将切片p的数据采用tab表示的多项式添加到crc之后计算出的新校验和。
func Update(crc uint32, tab *Table, p []byte) uint32

