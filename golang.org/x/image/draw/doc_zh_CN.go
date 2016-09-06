// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package draw provides image composition functions.
//
// See "The Go image/draw package" for an introduction to this package:
// http://golang.org/doc/articles/image_draw.html
//
// This package is a superset of and a drop-in replacement for the image/draw
// package in the standard library.
package draw

var (
	// NearestNeighbor is the nearest neighbor interpolator. It is very fast,
	// but usually gives very low quality results. When scaling up, the result
	// will look 'blocky'.
	NearestNeighbor = Interpolator(nnInterpolator{})

	// ApproxBiLinear is a mixture of the nearest neighbor and bi-linear
	// interpolators. It is fast, but usually gives medium quality results.
	//
	// It implements bi-linear interpolation when upscaling and a bi-linear
	// blend of the 4 nearest neighbor pixels when downscaling. This yields
	// nicer quality than nearest neighbor interpolation when upscaling, but
	// the time taken is independent of the number of source pixels, unlike the
	// bi-linear interpolator. When downscaling a large image, the performance
	// difference can be significant.
	ApproxBiLinear = Interpolator(ablInterpolator{})

	// BiLinear is the tent kernel. It is slow, but usually gives high quality
	// results.
	BiLinear = &Kernel{1, func(t float64) float64 {
		return 1 - t
	}}

	// CatmullRom is the Catmull-Rom kernel. It is very slow, but usually gives
	// very high quality results.
	//
	// It is an instance of the more general cubic BC-spline kernel with parameters
	// B=0 and C=0.5. See Mitchell and Netravali, "Reconstruction Filters in
	// Computer Graphics", Computer Graphics, Vol. 22, No. 4, pp. 221-228.
	CatmullRom = &Kernel{2, func(t float64) float64 {
		if t < 1 {
			return (1.5*t-2.5)*t*t + 1
		}
		return ((-0.5*t+2.5)*t-4)*t + 2
	}}
)

// Copy copies the part of the source image defined by src and sr and writes to the
// part of the destination image defined by dst and the translation of sr so that
// sr.Min translates to dp.
func Copy(dst Image, dp image.Point, src image.Image, sr image.Rectangle, opts *Options)

// Draw calls DrawMask with a nil mask.
func Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)

// DrawMask aligns r.Min in dst with sp in src and mp in mask and then replaces the
// rectangle r in dst with the result of a Porter-Duff composition. A nil mask is
// treated as opaque.
func DrawMask(dst Image, r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op Op)

// Drawer contains the Draw method.
type Drawer interface {
	// Draw aligns r.Min in dst with sp in src and then replaces the
	// rectangle r in dst with the result of drawing src on dst.
	Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point)
}

// FloydSteinberg is a Drawer that is the Src Op with Floyd-Steinberg error
// diffusion.
var FloydSteinberg Drawer = floydSteinberg{}

// Image is an image.Image with a Set method to change a single pixel.
type Image interface {
	image.Image
	Set(x, y int, c color.Color)
}

// Interpolator is an interpolation algorithm, when dst and src pixels don't have a
// 1:1 correspondence.
//
// Of the interpolators provided by this package:
//
//	- NearestNeighbor is fast but usually looks worst.
//	- CatmullRom is slow but usually looks best.
//	- ApproxBiLinear has reasonable speed and quality.
//
// The time taken depends on the size of dr. For kernel interpolators, the speed
// also depends on the size of sr, and so are often slower than non-kernel
// interpolators, especially when scaling down.
type Interpolator interface {
	Scaler
	Transformer
}

// Kernel is an interpolator that blends source pixels weighted by a symmetric
// kernel function.
type Kernel struct {
	// Support is the kernel support and must be >= 0. At(t) is assumed to be
	// zero when t >= Support.
	Support float64
	// At is the kernel function. It will only be called with t in the
	// range [0, Support).
	At func(t float64) float64
}

// NewScaler returns a Scaler that is optimized for scaling multiple times with the
// same fixed destination and source width and height.
func (q *Kernel) NewScaler(dw, dh, sw, sh int) Scaler

// Scale implements the Scaler interface.
func (q *Kernel) Scale(dst Image, dr image.Rectangle, src image.Image, sr image.Rectangle, opts *Options)

func (q *Kernel) Transform(dst Image, s2d *f64.Aff3, src image.Image, sr image.Rectangle, opts *Options)

// Op is a Porter-Duff compositing operator.
type Op int

const (
	// Over specifies ``(src in mask) over dst''.
	Over Op = Op(draw.Over)
	// Src specifies ``src in mask''.
	Src Op = Op(draw.Src)
)

// Draw implements the Drawer interface by calling the Draw function with this Op.
func (op Op) Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point)

// Options are optional parameters to Copy, Scale and Transform.
//
// A nil *Options means to use the default (zero) values of each field.
type Options struct {
	// Op is the compositing operator. The default value is Over.
	Op Op
}

// Quantizer produces a palette for an image.
type Quantizer interface {
	// Quantize appends up to cap(p) - len(p) colors to p and returns the
	// updated palette suitable for converting m to a paletted image.
	Quantize(p color.Palette, m image.Image) color.Palette
}

// Scaler scales the part of the source image defined by src and sr and writes to
// the part of the destination image defined by dst and dr.
//
// A Scaler is safe to use concurrently.
type Scaler interface {
	Scale(dst Image, dr image.Rectangle, src image.Image, sr image.Rectangle, opts *Options)
}

// Transformer transforms the part of the source image defined by src and sr and
// writes to the part of the destination image defined by dst and the affine
// transform m applied to sr.
//
// For example, if m is the matrix
//
// m00 m01 m02 m10 m11 m12
//
// then the src-space point (sx, sy) maps to the dst-space point (m00*sx + m01*sy +
// m02, m10*sx + m11*sy + m12).
//
// A Transformer is safe to use concurrently.
type Transformer interface {
	Transform(dst Image, m *f64.Aff3, src image.Image, sr image.Rectangle, opts *Options)
}
