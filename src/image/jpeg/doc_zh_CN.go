// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package jpeg implements a JPEG image decoder and encoder.
//
// JPEG is defined in ITU-T T.81: http://www.w3.org/Graphics/JPEG/itu-t81.pdf.

// Package jpeg implements a JPEG image
// decoder and encoder.
//
// JPEG is defined in ITU-T T.81:
// http://www.w3.org/Graphics/JPEG/itu-t81.pdf.
package jpeg

// DefaultQuality is the default quality encoding parameter.

// DefaultQuality is the default quality
// encoding parameter.
const DefaultQuality = 75

// Decode reads a JPEG image from r and returns it as an image.Image.

// Decode reads a JPEG image from r and
// returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig returns the color model and dimensions of a JPEG image without
// decoding the entire image.

// DecodeConfig returns the color model and
// dimensions of a JPEG image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error)

// Encode writes the Image m to w in JPEG 4:2:0 baseline format with the given
// options. Default parameters are used if a nil *Options is passed.

// Encode writes the Image m to w in JPEG
// 4:2:0 baseline format with the given
// options. Default parameters are used if
// a nil *Options is passed.
func Encode(w io.Writer, m image.Image, o *Options) error

// A FormatError reports that the input is not a valid JPEG.

// A FormatError reports that the input is
// not a valid JPEG.
type FormatError string

func (e FormatError) Error() string

// Options are the encoding parameters. Quality ranges from 1 to 100 inclusive,
// higher is better.

// Options are the encoding parameters.
// Quality ranges from 1 to 100 inclusive,
// higher is better.
type Options struct {
	Quality int
}

// Reader is deprecated.

// Reader is deprecated.
type Reader interface {
	io.ByteReader
	io.Reader
}

// An UnsupportedError reports that the input uses a valid but unimplemented JPEG
// feature.

// An UnsupportedError reports that the
// input uses a valid but unimplemented
// JPEG feature.
type UnsupportedError string

func (e UnsupportedError) Error() string
