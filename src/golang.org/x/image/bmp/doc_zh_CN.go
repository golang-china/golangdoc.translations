// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package bmp implements a BMP image decoder and encoder.
//
// The BMP specification is at http://www.digicamsoft.com/bmp/bmp.html.
package bmp

// ErrUnsupported means that the input BMP image uses a valid but unsupported
// feature.
var ErrUnsupported = errors.New("bmp: unsupported BMP image")

// Decode reads a BMP image from r and returns it as an image.Image. Limitation:
// The file must be 8, 24 or 32 bits per pixel.
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig returns the color model and dimensions of a BMP image without
// decoding the entire image. Limitation: The file must be 8, 24 or 32 bits per
// pixel.
func DecodeConfig(r io.Reader) (image.Config, error)

// Encode writes the image m to w in BMP format.
func Encode(w io.Writer, m image.Image) error
