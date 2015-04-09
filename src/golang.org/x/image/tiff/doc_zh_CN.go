// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package tiff implements a TIFF image decoder and encoder.
//
// The TIFF specification is at
// http://partners.adobe.com/public/developer/en/tiff/TIFF6.pdf
package tiff

// Decode reads a TIFF image from r and returns it as an image.Image. The type of
// Image returned depends on the contents of the TIFF.
func Decode(r io.Reader) (img image.Image, err error)

// DecodeConfig returns the color model and dimensions of a TIFF image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error)

// Encode writes the image m to w. opt determines the options used for encoding,
// such as the compression type. If opt is nil, an uncompressed image is written.
func Encode(w io.Writer, m image.Image, opt *Options) error

// CompressionType describes the type of compression used in Options.
type CompressionType int

const (
	Uncompressed CompressionType = iota
	Deflate
)

// A FormatError reports that the input is not a valid TIFF image.
type FormatError string

func (e FormatError) Error() string

// An InternalError reports that an internal error was encountered.
type InternalError string

func (e InternalError) Error() string

// Options are the encoding parameters.
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
type UnsupportedError string

func (e UnsupportedError) Error() string
