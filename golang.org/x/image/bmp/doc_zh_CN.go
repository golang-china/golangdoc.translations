// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package bmp implements a BMP image decoder and encoder.
//
// The BMP specification is at http://www.digicamsoft.com/bmp/bmp.html.

// bmp 包实现了 BMP 图像格式的编码器和解码器.
//
// BMP 图像格式规范 http://www.digicamsoft.com/bmp/bmp.html.
package bmp

// ErrUnsupported means that the input BMP image uses a valid but unsupported
// feature.

// ErrUnsupported 表示输入的 BMP 图像使用了当前包尚不支持的 BMP 特性.
var ErrUnsupported = errors.New("bmp: unsupported BMP image")

// Decode reads a BMP image from r and returns it as an image.Image. Limitation:
// The file must be 8, 24 or 32 bits per pixel.

// Decode 从 r 读取 BMP 图像, 并返回 image.Image.
// 限制: BMP图像必须是 8, 24 或 32 bit 深度.
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig returns the color model and dimensions of a BMP image without
// decoding the entire image. Limitation: The file must be 8, 24 or 32 bits per
// pixel.

// DecodeConfig 返回颜色模型和图像尺寸, 但是并不解码图像像素数据.
// 限制: BMP图像必须是 8, 24 或 32 bit 深度.
func DecodeConfig(r io.Reader) (image.Config, error)

// Encode writes the image m to w in BMP format.

// Encode 将 m 以 BMP 格式输出到 w.
func Encode(w io.Writer, m image.Image) error
