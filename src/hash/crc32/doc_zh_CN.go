// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package crc32 implements the 32-bit cyclic redundancy check, or CRC-32,
// checksum. See http://en.wikipedia.org/wiki/Cyclic_redundancy_check for
// information.

// Package crc32 implements the 32-bit
// cyclic redundancy check, or CRC-32,
// checksum. See
// http://en.wikipedia.org/wiki/Cyclic_redundancy_check
// for information.
package crc32

// Predefined polynomials.

// Predefined polynomials.
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

// The size of a CRC-32 checksum in bytes.
const Size = 4

// IEEETable is the table for the IEEE polynomial.

// IEEETable is the table for the IEEE
// polynomial.
var IEEETable = makeTable(IEEE)

// Checksum returns the CRC-32 checksum of data using the polynomial represented by
// the Table.

// Checksum returns the CRC-32 checksum of
// data using the polynomial represented by
// the Table.
func Checksum(data []byte, tab *Table) uint32

// ChecksumIEEE returns the CRC-32 checksum of data using the IEEE polynomial.

// ChecksumIEEE returns the CRC-32 checksum
// of data using the IEEE polynomial.
func ChecksumIEEE(data []byte) uint32

// New creates a new hash.Hash32 computing the CRC-32 checksum using the polynomial
// represented by the Table.

// New creates a new hash.Hash32 computing
// the CRC-32 checksum using the polynomial
// represented by the Table.
func New(tab *Table) hash.Hash32

// NewIEEE creates a new hash.Hash32 computing the CRC-32 checksum using the IEEE
// polynomial.

// NewIEEE creates a new hash.Hash32
// computing the CRC-32 checksum using the
// IEEE polynomial.
func NewIEEE() hash.Hash32

// Update returns the result of adding the bytes in p to the crc.

// Update returns the result of adding the
// bytes in p to the crc.
func Update(crc uint32, tab *Table, p []byte) uint32

// Table is a 256-word table representing the polynomial for efficient processing.

// Table is a 256-word table representing
// the polynomial for efficient processing.
type Table [256]uint32

// MakeTable returns the Table constructed from the specified polynomial.

// MakeTable returns the Table constructed
// from the specified polynomial.
func MakeTable(poly uint32) *Table
