// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package image implements a basic 2-D image library.
//
// The fundamental interface is called Image. An Image contains colors, which
// are described in the image/color package.
//
// Values of the Image interface are created either by calling functions such
// as NewRGBA and NewPaletted, or by calling Decode on an io.Reader containing
// image data in a format such as GIF, JPEG or PNG. Decoding any particular
// image format requires the prior registration of a decoder function.
// Registration is typically automatic as a side effect of initializing that
// format's package so that, to decode a PNG image, it suffices to have
//     import _ "image/png"
// in a program's main package. The _ means to import a package purely for its
// initialization side effects.
//
// See "The Go image package" for more details:
// https://golang.org/doc/articles/image_package.html

// Package image implements a basic 2-D image library.
//
// The fundamental interface is called Image. An Image contains colors, which
// are described in the image/color package.
//
// Values of the Image interface are created either by calling functions such
// as NewRGBA and NewPaletted, or by calling Decode on an io.Reader containing
// image data in a format such as GIF, JPEG or PNG. Decoding any particular
// image format requires the prior registration of a decoder function.
// Registration is typically automatic as a side effect of initializing that
// format's package so that, to decode a PNG image, it suffices to have
//     import _ "image/png"
// in a program's main package. The _ means to import a package purely for its
// initialization side effects.
//
// See "The Go image package" for more details:
// https://golang.org/doc/articles/image_package.html
package image

import (
    "bufio"
    "errors"
    "image/color"
    "io"
    "strconv"
)


const (
	YCbCrSubsampleRatio444 YCbCrSubsampleRatio = iota
	YCbCrSubsampleRatio422
	YCbCrSubsampleRatio420
	YCbCrSubsampleRatio440
	YCbCrSubsampleRatio411
	YCbCrSubsampleRatio410
)



var (
	// Black is an opaque black uniform image.
	Black = NewUniform(color.Black)
	// White is an opaque white uniform image.
	White = NewUniform(color.White)
	// Transparent is a fully transparent uniform image.
	Transparent = NewUniform(color.Transparent)
	// Opaque is a fully opaque uniform image.
	Opaque = NewUniform(color.Opaque)
)


// ErrFormat indicates that decoding encountered an unknown format.
var ErrFormat = errors.New("image: unknown format")


// ZP is the zero Point.
var ZP Point


// ZR is the zero Rectangle.
var ZR Rectangle


// Alpha is an in-memory image whose At method returns color.Alpha values.
type Alpha struct {
	// Pix holds the image's pixels, as alpha values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// Alpha16 is an in-memory image whose At method returns color.Alpha16 values.
type Alpha16 struct {
	// Pix holds the image's pixels, as alpha values in big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// CMYK is an in-memory image whose At method returns color.CMYK values.
type CMYK struct {
	// Pix holds the image's pixels, in C, M, Y, K order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// Config holds an image's color model and dimensions.
type Config struct {
	ColorModel    color.Model
	Width, Height int
}


// Gray is an in-memory image whose At method returns color.Gray values.
type Gray struct {
	// Pix holds the image's pixels, as gray values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// Gray16 is an in-memory image whose At method returns color.Gray16 values.
type Gray16 struct {
	// Pix holds the image's pixels, as gray values in big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// Image is a finite rectangular grid of color.Color values taken from a color
// model.
type Image interface {
	// ColorModel returns the Image's color model.
	ColorModel() color.Model
	// Bounds returns the domain for which At can return non-zero color.
	// The bounds do not necessarily contain the point (0, 0).
	Bounds() Rectangle
	// At returns the color of the pixel at (x, y).
	// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
	// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
	At(x, y int) color.Color
}


// NRGBA is an in-memory image whose At method returns color.NRGBA values.
type NRGBA struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// NRGBA64 is an in-memory image whose At method returns color.NRGBA64 values.
type NRGBA64 struct {
	// Pix holds the image's pixels, in R, G, B, A order and big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*8].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// NYCbCrA is an in-memory image of non-alpha-premultiplied Y'CbCr-with-alpha
// colors. A and AStride are analogous to the Y and YStride fields of the
// embedded YCbCr.
type NYCbCrA struct {
	A       []uint8
	AStride int
}


// Paletted is an in-memory image of uint8 indices into a given palette.
type Paletted struct {
	// Pix holds the image's pixels, as palette indices. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
	// Palette is the image's palette.
	Palette color.Palette
}


// PalettedImage is an image whose colors may come from a limited palette.
// If m is a PalettedImage and m.ColorModel() returns a color.Palette p,
// then m.At(x, y) should be equivalent to p[m.ColorIndexAt(x, y)]. If m's
// color model is not a color.Palette, then ColorIndexAt's behavior is
// undefined.
type PalettedImage interface {
	// ColorIndexAt returns the palette index of the pixel at (x, y).
	ColorIndexAt(x, y int) uint8
	Image
}


// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y int
}


// RGBA is an in-memory image whose At method returns color.RGBA values.
type RGBA struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// RGBA64 is an in-memory image whose At method returns color.RGBA64 values.
type RGBA64 struct {
	// Pix holds the image's pixels, in R, G, B, A order and big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*8].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}


// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
//
// A Rectangle is also an Image whose bounds are the rectangle itself. At
// returns color.Opaque for points in the rectangle and color.Transparent
// otherwise.
type Rectangle struct {
	Min, Max Point
}


// Uniform is an infinite-sized Image of uniform color.
// It implements the color.Color, color.Model, and Image interfaces.
type Uniform struct {
	C color.Color
}


// YCbCr is an in-memory image of Y'CbCr colors. There is one Y sample per
// pixel, but each Cb and Cr sample can span one or more pixels.
// YStride is the Y slice index delta between vertically adjacent pixels.
// CStride is the Cb and Cr slice index delta between vertically adjacent pixels
// that map to separate chroma samples.
// It is not an absolute requirement, but YStride and len(Y) are typically
// multiples of 8, and:
//     For 4:4:4, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/1.
//     For 4:2:2, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/2.
//     For 4:2:0, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/4.
//     For 4:4:0, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/2.
//     For 4:1:1, CStride == YStride/4 && len(Cb) == len(Cr) == len(Y)/4.
//     For 4:1:0, CStride == YStride/4 && len(Cb) == len(Cr) == len(Y)/8.

// YCbCr is an in-memory image of Y'CbCr colors. There is one Y sample per
// pixel, but each Cb and Cr sample can span one or more pixels.
// YStride is the Y slice index delta between vertically adjacent pixels.
// CStride is the Cb and Cr slice index delta between vertically adjacent pixels
// that map to separate chroma samples.
// It is not an absolute requirement, but YStride and len(Y) are typically
// multiples of 8, and:
//     For 4:4:4, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/1.
//     For 4:2:2, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/2.
//     For 4:2:0, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/4.
//     For 4:4:0, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/2.
//     For 4:1:1, CStride == YStride/4 && len(Cb) == len(Cr) == len(Y)/4.
//     For 4:1:0, CStride == YStride/4 && len(Cb) == len(Cr) == len(Y)/8.
type YCbCr struct {
	Y, Cb, Cr      []uint8
	YStride        int
	CStride        int
	SubsampleRatio YCbCrSubsampleRatio
	Rect           Rectangle
}


// YCbCrSubsampleRatio is the chroma subsample ratio used in a YCbCr image.
type YCbCrSubsampleRatio int


// Decode decodes an image that has been encoded in a registered format.
// The string returned is the format name used during format registration.
// Format registration is typically done by an init function in the codec-
// specific package.
func Decode(r io.Reader) (Image, string, error)

// DecodeConfig decodes the color model and dimensions of an image that has
// been encoded in a registered format. The string returned is the format name
// used during format registration. Format registration is typically done by
// an init function in the codec-specific package.
func DecodeConfig(r io.Reader) (Config, string, error)

// NewAlpha returns a new Alpha image with the given bounds.
func NewAlpha(r Rectangle) *Alpha

// NewAlpha16 returns a new Alpha16 image with the given bounds.
func NewAlpha16(r Rectangle) *Alpha16

// NewCMYK returns a new CMYK image with the given bounds.
func NewCMYK(r Rectangle) *CMYK

// NewGray returns a new Gray image with the given bounds.
func NewGray(r Rectangle) *Gray

// NewGray16 returns a new Gray16 image with the given bounds.
func NewGray16(r Rectangle) *Gray16

// NewNRGBA returns a new NRGBA image with the given bounds.
func NewNRGBA(r Rectangle) *NRGBA

// NewNRGBA64 returns a new NRGBA64 image with the given bounds.
func NewNRGBA64(r Rectangle) *NRGBA64

// NewNYCbCrA returns a new NYCbCrA image with the given bounds and subsample
// ratio.
func NewNYCbCrA(r Rectangle, subsampleRatio YCbCrSubsampleRatio) *NYCbCrA

// NewPaletted returns a new Paletted image with the given width, height and
// palette.
func NewPaletted(r Rectangle, p color.Palette) *Paletted

// NewRGBA returns a new RGBA image with the given bounds.
func NewRGBA(r Rectangle) *RGBA

// NewRGBA64 returns a new RGBA64 image with the given bounds.
func NewRGBA64(r Rectangle) *RGBA64

func NewUniform(c color.Color) *Uniform

// NewYCbCr returns a new YCbCr image with the given bounds and subsample
// ratio.
func NewYCbCr(r Rectangle, subsampleRatio YCbCrSubsampleRatio) *YCbCr

// Pt is shorthand for Point{X, Y}.
func Pt(X, Y int) Point

// Rect is shorthand for Rectangle{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect(x0, y0, x1, y1 int) Rectangle

// RegisterFormat registers an image format for use by Decode.
// Name is the name of the format, like "jpeg" or "png".
// Magic is the magic prefix that identifies the format's encoding. The magic
// string can contain "?" wildcards that each match any one byte.
// Decode is the function that decodes the encoded image.
// DecodeConfig is the function that decodes just its configuration.
func RegisterFormat(name, magic string, decode func(io.Reader) (Image, error), decodeConfig func(io.Reader) (Config, error))

func (*Alpha) AlphaAt(x, y int) color.Alpha

func (*Alpha) At(x, y int) color.Color

func (*Alpha) Bounds() Rectangle

func (*Alpha) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.
func (*Alpha) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*Alpha) PixOffset(x, y int) int

func (*Alpha) Set(x, y int, c color.Color)

func (*Alpha) SetAlpha(x, y int, c color.Alpha)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*Alpha) SubImage(r Rectangle) Image

func (*Alpha16) Alpha16At(x, y int) color.Alpha16

func (*Alpha16) At(x, y int) color.Color

func (*Alpha16) Bounds() Rectangle

func (*Alpha16) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.
func (*Alpha16) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*Alpha16) PixOffset(x, y int) int

func (*Alpha16) Set(x, y int, c color.Color)

func (*Alpha16) SetAlpha16(x, y int, c color.Alpha16)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*Alpha16) SubImage(r Rectangle) Image

func (*CMYK) At(x, y int) color.Color

func (*CMYK) Bounds() Rectangle

func (*CMYK) CMYKAt(x, y int) color.CMYK

func (*CMYK) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.
func (*CMYK) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*CMYK) PixOffset(x, y int) int

func (*CMYK) Set(x, y int, c color.Color)

func (*CMYK) SetCMYK(x, y int, c color.CMYK)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*CMYK) SubImage(r Rectangle) Image

func (*Gray) At(x, y int) color.Color

func (*Gray) Bounds() Rectangle

func (*Gray) ColorModel() color.Model

func (*Gray) GrayAt(x, y int) color.Gray

// Opaque scans the entire image and reports whether it is fully opaque.
func (*Gray) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*Gray) PixOffset(x, y int) int

func (*Gray) Set(x, y int, c color.Color)

func (*Gray) SetGray(x, y int, c color.Gray)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*Gray) SubImage(r Rectangle) Image

func (*Gray16) At(x, y int) color.Color

func (*Gray16) Bounds() Rectangle

func (*Gray16) ColorModel() color.Model

func (*Gray16) Gray16At(x, y int) color.Gray16

// Opaque scans the entire image and reports whether it is fully opaque.
func (*Gray16) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*Gray16) PixOffset(x, y int) int

func (*Gray16) Set(x, y int, c color.Color)

func (*Gray16) SetGray16(x, y int, c color.Gray16)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*Gray16) SubImage(r Rectangle) Image

func (*NRGBA) At(x, y int) color.Color

func (*NRGBA) Bounds() Rectangle

func (*NRGBA) ColorModel() color.Model

func (*NRGBA) NRGBAAt(x, y int) color.NRGBA

// Opaque scans the entire image and reports whether it is fully opaque.
func (*NRGBA) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*NRGBA) PixOffset(x, y int) int

func (*NRGBA) Set(x, y int, c color.Color)

func (*NRGBA) SetNRGBA(x, y int, c color.NRGBA)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*NRGBA) SubImage(r Rectangle) Image

func (*NRGBA64) At(x, y int) color.Color

func (*NRGBA64) Bounds() Rectangle

func (*NRGBA64) ColorModel() color.Model

func (*NRGBA64) NRGBA64At(x, y int) color.NRGBA64

// Opaque scans the entire image and reports whether it is fully opaque.
func (*NRGBA64) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*NRGBA64) PixOffset(x, y int) int

func (*NRGBA64) Set(x, y int, c color.Color)

func (*NRGBA64) SetNRGBA64(x, y int, c color.NRGBA64)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*NRGBA64) SubImage(r Rectangle) Image

// AOffset returns the index of the first element of A that corresponds to the
// pixel at (x, y).
func (*NYCbCrA) AOffset(x, y int) int

func (*NYCbCrA) At(x, y int) color.Color

func (*NYCbCrA) ColorModel() color.Model

func (*NYCbCrA) NYCbCrAAt(x, y int) color.NYCbCrA

// Opaque scans the entire image and reports whether it is fully opaque.
func (*NYCbCrA) Opaque() bool

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*NYCbCrA) SubImage(r Rectangle) Image

func (*Paletted) At(x, y int) color.Color

func (*Paletted) Bounds() Rectangle

func (*Paletted) ColorIndexAt(x, y int) uint8

func (*Paletted) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.
func (*Paletted) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*Paletted) PixOffset(x, y int) int

func (*Paletted) Set(x, y int, c color.Color)

func (*Paletted) SetColorIndex(x, y int, index uint8)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*Paletted) SubImage(r Rectangle) Image

func (*RGBA) At(x, y int) color.Color

func (*RGBA) Bounds() Rectangle

func (*RGBA) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.
func (*RGBA) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*RGBA) PixOffset(x, y int) int

func (*RGBA) RGBAAt(x, y int) color.RGBA

func (*RGBA) Set(x, y int, c color.Color)

func (*RGBA) SetRGBA(x, y int, c color.RGBA)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*RGBA) SubImage(r Rectangle) Image

func (*RGBA64) At(x, y int) color.Color

func (*RGBA64) Bounds() Rectangle

func (*RGBA64) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.
func (*RGBA64) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (*RGBA64) PixOffset(x, y int) int

func (*RGBA64) RGBA64At(x, y int) color.RGBA64

func (*RGBA64) Set(x, y int, c color.Color)

func (*RGBA64) SetRGBA64(x, y int, c color.RGBA64)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*RGBA64) SubImage(r Rectangle) Image

func (*Uniform) At(x, y int) color.Color

func (*Uniform) Bounds() Rectangle

func (*Uniform) ColorModel() color.Model

func (*Uniform) Convert(color.Color) color.Color

// Opaque scans the entire image and reports whether it is fully opaque.
func (*Uniform) Opaque() bool

func (*Uniform) RGBA() (r, g, b, a uint32)

func (*YCbCr) At(x, y int) color.Color

func (*YCbCr) Bounds() Rectangle

// COffset returns the index of the first element of Cb or Cr that corresponds
// to the pixel at (x, y).
func (*YCbCr) COffset(x, y int) int

func (*YCbCr) ColorModel() color.Model

func (*YCbCr) Opaque() bool

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (*YCbCr) SubImage(r Rectangle) Image

func (*YCbCr) YCbCrAt(x, y int) color.YCbCr

// YOffset returns the index of the first element of Y that corresponds to
// the pixel at (x, y).
func (*YCbCr) YOffset(x, y int) int

// Add returns the vector p+q.
func (Point) Add(q Point) Point

// Div returns the vector p/k.
func (Point) Div(k int) Point

// Eq reports whether p and q are equal.
func (Point) Eq(q Point) bool

// In reports whether p is in r.
func (Point) In(r Rectangle) bool

// Mod returns the point q in r such that p.X-q.X is a multiple of r's width
// and p.Y-q.Y is a multiple of r's height.
func (Point) Mod(r Rectangle) Point

// Mul returns the vector p*k.
func (Point) Mul(k int) Point

// String returns a string representation of p like "(3,4)".
func (Point) String() string

// Sub returns the vector p-q.
func (Point) Sub(q Point) Point

// Add returns the rectangle r translated by p.
func (Rectangle) Add(p Point) Rectangle

// At implements the Image interface.
func (Rectangle) At(x, y int) color.Color

// Bounds implements the Image interface.
func (Rectangle) Bounds() Rectangle

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (Rectangle) Canon() Rectangle

// ColorModel implements the Image interface.
func (Rectangle) ColorModel() color.Model

// Dx returns r's width.
func (Rectangle) Dx() int

// Dy returns r's height.
func (Rectangle) Dy() int

// Empty reports whether the rectangle contains no points.
func (Rectangle) Empty() bool

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (Rectangle) Eq(s Rectangle) bool

// In reports whether every point in r is in s.
func (Rectangle) In(s Rectangle) bool

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (Rectangle) Inset(n int) Rectangle

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (Rectangle) Intersect(s Rectangle) Rectangle

// Overlaps reports whether r and s have a non-empty intersection.
func (Rectangle) Overlaps(s Rectangle) bool

// Size returns r's width and height.
func (Rectangle) Size() Point

// String returns a string representation of r like "(3,4)-(6,5)".
func (Rectangle) String() string

// Sub returns the rectangle r translated by -p.
func (Rectangle) Sub(p Point) Rectangle

// Union returns the smallest rectangle that contains both r and s.
func (Rectangle) Union(s Rectangle) Rectangle

func (YCbCrSubsampleRatio) String() string

