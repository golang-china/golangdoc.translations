// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package webp implements a decoder for WEBP images.
//
// WEBP is defined at: https://developers.google.com/speed/webp/docs/riff_container
package webp

// Decode reads a WEBP image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig returns the color model and dimensions of a WEBP image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error)
