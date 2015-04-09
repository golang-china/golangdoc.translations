// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package nycbcra provides non-alpha-premultiplied Y'CbCr-with-alpha image and
// color types.
package nycbcra

// ColorModel is the Model for non-alpha-premultiplied Y'CbCr-with-alpha colors.
var ColorModel color.Model = color.ModelFunc(nYCbCrAModel)

// Color represents a non-alpha-premultiplied Y'CbCr-with-alpha color, having 8
// bits each for one luma, two chroma and one alpha component.
type Color struct {
	color.YCbCr
	A uint8
}

func (c Color) RGBA() (r, g, b, a uint32)

// Image is an in-memory image of non-alpha-premultiplied Y'CbCr-with-alpha colors.
// A and AStride are analogous to the Y and YStride fields of the embedded YCbCr.
type Image struct {
	image.YCbCr
	A       []uint8
	AStride int
}

// New returns a new Image with the given bounds and subsample ratio.
func New(r image.Rectangle, subsampleRatio image.YCbCrSubsampleRatio) *Image

// AOffset returns the index of the first element of A that corresponds to the
// pixel at (x, y).
func (p *Image) AOffset(x, y int) int

func (p *Image) At(x, y int) color.Color

func (p *Image) ColorModel() color.Model

func (p *Image) NYCbCrAAt(x, y int) Color

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *Image) Opaque() bool

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Image) SubImage(r image.Rectangle) image.Image
