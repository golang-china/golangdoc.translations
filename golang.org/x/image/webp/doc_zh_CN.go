// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package webp implements a decoder for WEBP images.
//
// WEBP is defined at: https://developers.google.com/speed/webp/docs/riff_container

// webp 包实现了 WEBP 图像格式的解码器.
//
// WEBP 图像格式规范
// https://developers.google.com/speed/webp/docs/riff_container
package webp

// Decode reads a WEBP image from r and returns it as an image.Image.

// Decode 从 r 读取 WEBP 图像, 并返回 image.Image.
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig returns the color model and dimensions of a WEBP image without
// decoding the entire image.

// DecodeConfig 返回颜色模型和图像尺寸, 但是并不解码图像像素数据.
func DecodeConfig(r io.Reader) (image.Config, error)
