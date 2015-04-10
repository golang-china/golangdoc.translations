// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package lzw implements the Lempel-Ziv-Welch compressed data format, described in
// T. A. Welch, ``A Technique for High-Performance Data Compression'', Computer,
// 17(6) (June 1984), pp 8-19.
//
// In particular, it implements LZW as used by the GIF and PDF file formats, which
// means variable-width codes up to 12 bits and the first two non-literal codes are
// a clear code and an EOF code.
//
// The TIFF file format uses a similar but incompatible version of the LZW
// algorithm. See the golang.org/x/image/tiff/lzw package for an implementation.

// lzw 包实现了 Lempel-Ziv-Welch 数据压缩格式, 这是一种 T. A. Welch 在 ``A Technique for High-Performance Data Compression''
// 一文(Computer, 17(6) (June 1984), pp 8-19) 提出的一种压缩格式.
//
// 本包实现了用于 GIF/TIFF/PDF 文件的 lzw 压缩格式, 这是一种最长达到 12 位的变长码,
// 头两个非字面码为 clear 和 EOF 码.
//
// TIFF 文件格式采用了类似但不兼容的 LZW 算法,
// 具体实现请参考 golang.org/x/image/tiff/lzw 包.
package lzw

// NewReader creates a new io.ReadCloser. Reads from the returned io.ReadCloser
// read and decompress data from r. If r does not also implement io.ByteReader, the
// decompressor may read more data than necessary from r. It is the caller's
// responsibility to call Close on the ReadCloser when finished reading. The number
// of bits to use for literal codes, litWidth, must be in the range [2,8] and is
// typically 8.

// NewReader 创建一个 io.ReadCloser, 它从 r 读取并解压数据.
// 调用者有责任在结束读取后调用返回值的 Close 方法;
// litWidth 指定字面码的位数, 必须在 [2,8] 范围内, 一般为8.
func NewReader(r io.Reader, order Order, litWidth int) io.ReadCloser

// NewWriter creates a new io.WriteCloser. Writes to the returned io.WriteCloser
// are compressed and written to w. It is the caller's responsibility to call Close
// on the WriteCloser when finished writing. The number of bits to use for literal
// codes, litWidth, must be in the range [2,8] and is typically 8.

// NewWriter 创建一个 io.WriteCloser, 它将数据压缩后写入 w.
// 调用者有责任在结束写入后调用返回值的 Close 方法;
// litWidth 指定字面码的位数, 必须在 [2,8] 范围内, 一般为 8.
func NewWriter(w io.Writer, order Order, litWidth int) io.WriteCloser

// Order specifies the bit ordering in an LZW data stream.

// Order 指定一个 lzw 数据流的位顺序.
type Order int

const (
	// LSB means Least Significant Bits first, as used in the GIF file format.
	LSB Order = iota
	// MSB means Most Significant Bits first, as used in the TIFF and PDF
	// file formats.
	MSB
)
