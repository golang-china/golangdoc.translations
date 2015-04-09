// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package lzw implements the Lempel-Ziv-Welch compressed data format, described in
// T. A. Welch, ``A Technique for High-Performance Data Compression'', Computer,
// 17(6) (June 1984), pp 8-19.
//
// In particular, it implements LZW as used by the TIFF file format, including an
// "off by one" algorithmic difference when compared to standard LZW.
package lzw

// NewReader creates a new io.ReadCloser. Reads from the returned io.ReadCloser
// read and decompress data from r. It is the caller's responsibility to call Close
// on the ReadCloser when finished reading. The number of bits to use for literal
// codes, litWidth, must be in the range [2,8] and is typically 8.
func NewReader(r io.Reader, order Order, litWidth int) io.ReadCloser

// Order specifies the bit ordering in an LZW data stream.
type Order int

const (
	// LSB means Least Significant Bits first, as used in the GIF file format.
	LSB Order = iota
	// MSB means Most Significant Bits first, as used in the TIFF and PDF
	// file formats.
	MSB
)
