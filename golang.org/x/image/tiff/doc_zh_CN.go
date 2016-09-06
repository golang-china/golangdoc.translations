// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package tiff implements a TIFF image decoder and encoder.
//
// The TIFF specification is at
// http://partners.adobe.com/public/developer/en/tiff/TIFF6.pdf

// tiff 包实现了 TIFF 图像格式的编码器和解码器.
//
// TIFF 图像格式规范
// http://partners.adobe.com/public/developer/en/tiff/TIFF6.pdf
package tiff

// Decode reads a TIFF image from r and returns it as an image.Image. The type of
// Image returned depends on the contents of the TIFF.

// Decode 从 r 读取 TIFF 图像, 并返回 image.Image.
// 返回的图像类型依赖于输入的 TIFF 图像格式.
func Decode(r io.Reader) (img image.Image, err error)

// DecodeConfig returns the color model and dimensions of a TIFF image without
// decoding the entire image.

// DecodeConfig 返回颜色模型和图像尺寸, 但是并不解码图像像素数据.
func DecodeConfig(r io.Reader) (image.Config, error)

// Encode writes the image m to w. opt determines the options used for encoding,
// such as the compression type. If opt is nil, an uncompressed image is written.

// Encode 将 m 以 TIFF 格式输出到 w.
// opt 用于指定编码的参数, 比如 压缩类型.
// 如果 opt 为空, 则默认无压缩.
func Encode(w io.Writer, m image.Image, opt *Options) error

// CompressionType describes the type of compression used in Options.

// CompressionType 表示编码图像时的压缩类型, 在 Options 使用.
type CompressionType int

const (
	Uncompressed CompressionType = iota
	Deflate
)

// A FormatError reports that the input is not a valid TIFF image.

// FormatError 表示输入不是一个有效的 TIFF 图像.
type FormatError string

func (e FormatError) Error() string

// An InternalError reports that an internal error was encountered.

// InternalError 表示内部产生的错误.
type InternalError string

func (e InternalError) Error() string

// Options are the encoding parameters.

// Options 表示编码参数.
type Options struct {
	// Compression is the type of compression used.
	Compression CompressionType
	// Predictor determines whether a differencing predictor is used;
	// if true, instead of each pixel's color, the color difference to the
	// preceding one is saved.  This improves the compression for certain
	// types of images and compressors. For example, it works well for
	// photos with Deflate compression.
	Predictor bool
}

// An UnsupportedError reports that the input uses a valid but unimplemented
// feature.

// UnsupportedError 表示输入是有效的 TIFF 图像, 但是使用了当前包不支持的特性.
type UnsupportedError string

func (e UnsupportedError) Error() string
