// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package vp8l implements a decoder for the VP8L lossless image format.
//
// The VP8L specification is at:
// https://developers.google.com/speed/webp/docs/riff_container
package vp8l

// Decode decodes a VP8L image from r.
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig decodes the color model and dimensions of a VP8L image from r.
func DecodeConfig(r io.Reader) (image.Config, error)
