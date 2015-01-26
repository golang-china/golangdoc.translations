// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package image implements a basic 2-D image library.
//
// The fundamental interface is called Image. An Image contains colors, which are
// described in the image/color package.
//
// Values of the Image interface are created either by calling functions such as
// NewRGBA and NewPaletted, or by calling Decode on an io.Reader containing image
// data in a format such as GIF, JPEG or PNG. Decoding any particular image format
// requires the prior registration of a decoder function. Registration is typically
// automatic as a side effect of initializing that format's package so that, to
// decode a PNG image, it suffices to have
//
//	import _ "image/png"
//
// in a program's main package. The _ means to import a package purely for its
// initialization side effects.
//
// See "The Go image package" for more details:
// http://golang.org/doc/articles/image_package.html

// Package image implements a basic 2-D
// image library.
//
// The fundamental interface is called
// Image. An Image contains colors, which
// are described in the image/color
// package.
//
// Values of the Image interface are
// created either by calling functions such
// as NewRGBA and NewPaletted, or by
// calling Decode on an io.Reader
// containing image data in a format such
// as GIF, JPEG or PNG. Decoding any
// particular image format requires the
// prior registration of a decoder
// function. Registration is typically
// automatic as a side effect of
// initializing that format's package so
// that, to decode a PNG image, it suffices
// to have
//
//	import _ "image/png"
//
// in a program's main package. The _ means
// to import a package purely for its
// initialization side effects.
//
// See "The Go image package" for more
// details:
// http://golang.org/doc/articles/image_package.html
package image

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

// ErrFormat indicates that decoding
// encountered an unknown format.
var ErrFormat = errors.New("image: unknown format")

// RegisterFormat registers an image format for use by Decode. Name is the name of
// the format, like "jpeg" or "png". Magic is the magic prefix that identifies the
// format's encoding. The magic string can contain "?" wildcards that each match
// any one byte. Decode is the function that decodes the encoded image.
// DecodeConfig is the function that decodes just its configuration.

// RegisterFormat registers an image format
// for use by Decode. Name is the name of
// the format, like "jpeg" or "png". Magic
// is the magic prefix that identifies the
// format's encoding. The magic string can
// contain "?" wildcards that each match
// any one byte. Decode is the function
// that decodes the encoded image.
// DecodeConfig is the function that
// decodes just its configuration.
func RegisterFormat(name, magic string, decode func(io.Reader) (Image, error), decodeConfig func(io.Reader) (Config, error))

// Alpha is an in-memory image whose At method returns color.Alpha values.

// Alpha is an in-memory image whose At
// method returns color.Alpha values.
type Alpha struct {
	// Pix holds the image's pixels, as alpha values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}

// NewAlpha returns a new Alpha with the given bounds.

// NewAlpha returns a new Alpha with the
// given bounds.
func NewAlpha(r Rectangle) *Alpha

func (p *Alpha) AlphaAt(x, y int) color.Alpha

func (p *Alpha) At(x, y int) color.Color

func (p *Alpha) Bounds() Rectangle

func (p *Alpha) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *Alpha) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *Alpha) PixOffset(x, y int) int

func (p *Alpha) Set(x, y int, c color.Color)

func (p *Alpha) SetAlpha(x, y int, c color.Alpha)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *Alpha) SubImage(r Rectangle) Image

// Alpha16 is an in-memory image whose At method returns color.Alpha64 values.

// Alpha16 is an in-memory image whose At
// method returns color.Alpha64 values.
type Alpha16 struct {
	// Pix holds the image's pixels, as alpha values in big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}

// NewAlpha16 returns a new Alpha16 with the given bounds.

// NewAlpha16 returns a new Alpha16 with
// the given bounds.
func NewAlpha16(r Rectangle) *Alpha16

func (p *Alpha16) Alpha16At(x, y int) color.Alpha16

func (p *Alpha16) At(x, y int) color.Color

func (p *Alpha16) Bounds() Rectangle

func (p *Alpha16) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *Alpha16) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *Alpha16) PixOffset(x, y int) int

func (p *Alpha16) Set(x, y int, c color.Color)

func (p *Alpha16) SetAlpha16(x, y int, c color.Alpha16)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *Alpha16) SubImage(r Rectangle) Image

// Config holds an image's color model and dimensions.

// Config holds an image's color model and
// dimensions.
type Config struct {
	ColorModel    color.Model
	Width, Height int
}

// DecodeConfig decodes the color model and dimensions of an image that has been
// encoded in a registered format. The string returned is the format name used
// during format registration. Format registration is typically done by an init
// function in the codec-specific package.

// DecodeConfig decodes the color model and
// dimensions of an image that has been
// encoded in a registered format. The
// string returned is the format name used
// during format registration. Format
// registration is typically done by an
// init function in the codec-specific
// package.
func DecodeConfig(r io.Reader) (Config, string, error)

// Gray is an in-memory image whose At method returns color.Gray values.

// Gray is an in-memory image whose At
// method returns color.Gray values.
type Gray struct {
	// Pix holds the image's pixels, as gray values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}

// NewGray returns a new Gray with the given bounds.

// NewGray returns a new Gray with the
// given bounds.
func NewGray(r Rectangle) *Gray

func (p *Gray) At(x, y int) color.Color

func (p *Gray) Bounds() Rectangle

func (p *Gray) ColorModel() color.Model

func (p *Gray) GrayAt(x, y int) color.Gray

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *Gray) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *Gray) PixOffset(x, y int) int

func (p *Gray) Set(x, y int, c color.Color)

func (p *Gray) SetGray(x, y int, c color.Gray)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *Gray) SubImage(r Rectangle) Image

// Gray16 is an in-memory image whose At method returns color.Gray16 values.

// Gray16 is an in-memory image whose At
// method returns color.Gray16 values.
type Gray16 struct {
	// Pix holds the image's pixels, as gray values in big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}

// NewGray16 returns a new Gray16 with the given bounds.

// NewGray16 returns a new Gray16 with the
// given bounds.
func NewGray16(r Rectangle) *Gray16

func (p *Gray16) At(x, y int) color.Color

func (p *Gray16) Bounds() Rectangle

func (p *Gray16) ColorModel() color.Model

func (p *Gray16) Gray16At(x, y int) color.Gray16

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *Gray16) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *Gray16) PixOffset(x, y int) int

func (p *Gray16) Set(x, y int, c color.Color)

func (p *Gray16) SetGray16(x, y int, c color.Gray16)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *Gray16) SubImage(r Rectangle) Image

// Image is a finite rectangular grid of color.Color values taken from a color
// model.

// Image is a finite rectangular grid of
// color.Color values taken from a color
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

// Decode decodes an image that has been encoded in a registered format. The string
// returned is the format name used during format registration. Format registration
// is typically done by an init function in the codec- specific package.

// Decode decodes an image that has been
// encoded in a registered format. The
// string returned is the format name used
// during format registration. Format
// registration is typically done by an
// init function in the codec- specific
// package.
func Decode(r io.Reader) (Image, string, error)

// NRGBA is an in-memory image whose At method returns color.NRGBA values.

// NRGBA is an in-memory image whose At
// method returns color.NRGBA values.
type NRGBA struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}

// NewNRGBA returns a new NRGBA with the given bounds.

// NewNRGBA returns a new NRGBA with the
// given bounds.
func NewNRGBA(r Rectangle) *NRGBA

func (p *NRGBA) At(x, y int) color.Color

func (p *NRGBA) Bounds() Rectangle

func (p *NRGBA) ColorModel() color.Model

func (p *NRGBA) NRGBAAt(x, y int) color.NRGBA

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *NRGBA) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *NRGBA) PixOffset(x, y int) int

func (p *NRGBA) Set(x, y int, c color.Color)

func (p *NRGBA) SetNRGBA(x, y int, c color.NRGBA)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *NRGBA) SubImage(r Rectangle) Image

// NRGBA64 is an in-memory image whose At method returns color.NRGBA64 values.

// NRGBA64 is an in-memory image whose At
// method returns color.NRGBA64 values.
type NRGBA64 struct {
	// Pix holds the image's pixels, in R, G, B, A order and big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*8].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}

// NewNRGBA64 returns a new NRGBA64 with the given bounds.

// NewNRGBA64 returns a new NRGBA64 with
// the given bounds.
func NewNRGBA64(r Rectangle) *NRGBA64

func (p *NRGBA64) At(x, y int) color.Color

func (p *NRGBA64) Bounds() Rectangle

func (p *NRGBA64) ColorModel() color.Model

func (p *NRGBA64) NRGBA64At(x, y int) color.NRGBA64

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *NRGBA64) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *NRGBA64) PixOffset(x, y int) int

func (p *NRGBA64) Set(x, y int, c color.Color)

func (p *NRGBA64) SetNRGBA64(x, y int, c color.NRGBA64)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *NRGBA64) SubImage(r Rectangle) Image

// Paletted is an in-memory image of uint8 indices into a given palette.

// Paletted is an in-memory image of uint8
// indices into a given palette.
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

// NewPaletted returns a new Paletted with the given width, height and palette.

// NewPaletted returns a new Paletted with
// the given width, height and palette.
func NewPaletted(r Rectangle, p color.Palette) *Paletted

func (p *Paletted) At(x, y int) color.Color

func (p *Paletted) Bounds() Rectangle

func (p *Paletted) ColorIndexAt(x, y int) uint8

func (p *Paletted) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *Paletted) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *Paletted) PixOffset(x, y int) int

func (p *Paletted) Set(x, y int, c color.Color)

func (p *Paletted) SetColorIndex(x, y int, index uint8)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *Paletted) SubImage(r Rectangle) Image

// PalettedImage is an image whose colors may come from a limited palette. If m is
// a PalettedImage and m.ColorModel() returns a PalettedColorModel p, then m.At(x,
// y) should be equivalent to p[m.ColorIndexAt(x, y)]. If m's color model is not a
// PalettedColorModel, then ColorIndexAt's behavior is undefined.

// PalettedImage is an image whose colors
// may come from a limited palette. If m is
// a PalettedImage and m.ColorModel()
// returns a PalettedColorModel p, then
// m.At(x, y) should be equivalent to
// p[m.ColorIndexAt(x, y)]. If m's color
// model is not a PalettedColorModel, then
// ColorIndexAt's behavior is undefined.
type PalettedImage interface {
	// ColorIndexAt returns the palette index of the pixel at (x, y).
	ColorIndexAt(x, y int) uint8
	Image
}

// A Point is an X, Y coordinate pair. The axes increase right and down.

// A Point is an X, Y coordinate pair. The
// axes increase right and down.
type Point struct {
	X, Y int
}

// ZP is the zero Point.

// ZP is the zero Point.
var ZP Point

// Pt is shorthand for Point{X, Y}.

// Pt is shorthand for Point{X, Y}.
func Pt(X, Y int) Point

// Add returns the vector p+q.

// Add returns the vector p+q.
func (p Point) Add(q Point) Point

// Div returns the vector p/k.

// Div returns the vector p/k.
func (p Point) Div(k int) Point

// Eq reports whether p and q are equal.

// Eq reports whether p and q are equal.
func (p Point) Eq(q Point) bool

// In reports whether p is in r.

// In reports whether p is in r.
func (p Point) In(r Rectangle) bool

// Mod returns the point q in r such that p.X-q.X is a multiple of r's width and
// p.Y-q.Y is a multiple of r's height.

// Mod returns the point q in r such that
// p.X-q.X is a multiple of r's width and
// p.Y-q.Y is a multiple of r's height.
func (p Point) Mod(r Rectangle) Point

// Mul returns the vector p*k.

// Mul returns the vector p*k.
func (p Point) Mul(k int) Point

// String returns a string representation of p like "(3,4)".

// String returns a string representation
// of p like "(3,4)".
func (p Point) String() string

// Sub returns the vector p-q.

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point

// RGBA is an in-memory image whose At method returns color.RGBA values.

// RGBA is an in-memory image whose At
// method returns color.RGBA values.
type RGBA struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}

// NewRGBA returns a new RGBA with the given bounds.

// NewRGBA returns a new RGBA with the
// given bounds.
func NewRGBA(r Rectangle) *RGBA

func (p *RGBA) At(x, y int) color.Color

func (p *RGBA) Bounds() Rectangle

func (p *RGBA) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *RGBA) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *RGBA) PixOffset(x, y int) int

func (p *RGBA) RGBAAt(x, y int) color.RGBA

func (p *RGBA) Set(x, y int, c color.Color)

func (p *RGBA) SetRGBA(x, y int, c color.RGBA)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *RGBA) SubImage(r Rectangle) Image

// RGBA64 is an in-memory image whose At method returns color.RGBA64 values.

// RGBA64 is an in-memory image whose At
// method returns color.RGBA64 values.
type RGBA64 struct {
	// Pix holds the image's pixels, in R, G, B, A order and big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*8].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect Rectangle
}

// NewRGBA64 returns a new RGBA64 with the given bounds.

// NewRGBA64 returns a new RGBA64 with the
// given bounds.
func NewRGBA64(r Rectangle) *RGBA64

func (p *RGBA64) At(x, y int) color.Color

func (p *RGBA64) Bounds() Rectangle

func (p *RGBA64) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (p *RGBA64) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to the
// pixel at (x, y).

// PixOffset returns the index of the first
// element of Pix that corresponds to the
// pixel at (x, y).
func (p *RGBA64) PixOffset(x, y int) int

func (p *RGBA64) RGBA64At(x, y int) color.RGBA64

func (p *RGBA64) Set(x, y int, c color.Color)

func (p *RGBA64) SetRGBA64(x, y int, c color.RGBA64)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *RGBA64) SubImage(r Rectangle) Image

// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y. It
// is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.

// A Rectangle contains the points with
// Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and
// likewise for Y. Points are always
// well-formed. A rectangle's methods
// always return well-formed outputs for
// well-formed inputs.
type Rectangle struct {
	Min, Max Point
}

// ZR is the zero Rectangle.

// ZR is the zero Rectangle.
var ZR Rectangle

// Rect is shorthand for Rectangle{Pt(x0, y0), Pt(x1, y1)}.

// Rect is shorthand for Rectangle{Pt(x0,
// y0), Pt(x1, y1)}.
func Rect(x0, y0, x1, y1 int) Rectangle

// Add returns the rectangle r translated by p.

// Add returns the rectangle r translated
// by p.
func (r Rectangle) Add(p Point) Rectangle

// Canon returns the canonical version of r. The returned rectangle has minimum and
// maximum coordinates swapped if necessary so that it is well-formed.

// Canon returns the canonical version of
// r. The returned rectangle has minimum
// and maximum coordinates swapped if
// necessary so that it is well-formed.
func (r Rectangle) Canon() Rectangle

// Dx returns r's width.

// Dx returns r's width.
func (r Rectangle) Dx() int

// Dy returns r's height.

// Dy returns r's height.
func (r Rectangle) Dy() int

// Empty reports whether the rectangle contains no points.

// Empty reports whether the rectangle
// contains no points.
func (r Rectangle) Empty() bool

// Eq reports whether r and s are equal.

// Eq reports whether r and s are equal.
func (r Rectangle) Eq(s Rectangle) bool

// In reports whether every point in r is in s.

// In reports whether every point in r is
// in s.
func (r Rectangle) In(s Rectangle) bool

// Inset returns the rectangle r inset by n, which may be negative. If either of
// r's dimensions is less than 2*n then an empty rectangle near the center of r
// will be returned.

// Inset returns the rectangle r inset by
// n, which may be negative. If either of
// r's dimensions is less than 2*n then an
// empty rectangle near the center of r
// will be returned.
func (r Rectangle) Inset(n int) Rectangle

// Intersect returns the largest rectangle contained by both r and s. If the two
// rectangles do not overlap then the zero rectangle will be returned.

// Intersect returns the largest rectangle
// contained by both r and s. If the two
// rectangles do not overlap then the zero
// rectangle will be returned.
func (r Rectangle) Intersect(s Rectangle) Rectangle

// Overlaps reports whether r and s have a non-empty intersection.

// Overlaps reports whether r and s have a
// non-empty intersection.
func (r Rectangle) Overlaps(s Rectangle) bool

// Size returns r's width and height.

// Size returns r's width and height.
func (r Rectangle) Size() Point

// String returns a string representation of r like "(3,4)-(6,5)".

// String returns a string representation
// of r like "(3,4)-(6,5)".
func (r Rectangle) String() string

// Sub returns the rectangle r translated by -p.

// Sub returns the rectangle r translated
// by -p.
func (r Rectangle) Sub(p Point) Rectangle

// Union returns the smallest rectangle that contains both r and s.

// Union returns the smallest rectangle
// that contains both r and s.
func (r Rectangle) Union(s Rectangle) Rectangle

// Uniform is an infinite-sized Image of uniform color. It implements the
// color.Color, color.Model, and Image interfaces.

// Uniform is an infinite-sized Image of
// uniform color. It implements the
// color.Color, color.Model, and Image
// interfaces.
type Uniform struct {
	C color.Color
}

func NewUniform(c color.Color) *Uniform

func (c *Uniform) At(x, y int) color.Color

func (c *Uniform) Bounds() Rectangle

func (c *Uniform) ColorModel() color.Model

func (c *Uniform) Convert(color.Color) color.Color

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque scans the entire image and
// reports whether it is fully opaque.
func (c *Uniform) Opaque() bool

func (c *Uniform) RGBA() (r, g, b, a uint32)

// YCbCr is an in-memory image of Y'CbCr colors. There is one Y sample per pixel,
// but each Cb and Cr sample can span one or more pixels. YStride is the Y slice
// index delta between vertically adjacent pixels. CStride is the Cb and Cr slice
// index delta between vertically adjacent pixels that map to separate chroma
// samples. It is not an absolute requirement, but YStride and len(Y) are typically
// multiples of 8, and:
//
//	For 4:4:4, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/1.
//	For 4:2:2, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/2.
//	For 4:2:0, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/4.
//	For 4:4:0, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/2.

// YCbCr is an in-memory image of Y'CbCr
// colors. There is one Y sample per pixel,
// but each Cb and Cr sample can span one
// or more pixels. YStride is the Y slice
// index delta between vertically adjacent
// pixels. CStride is the Cb and Cr slice
// index delta between vertically adjacent
// pixels that map to separate chroma
// samples. It is not an absolute
// requirement, but YStride and len(Y) are
// typically multiples of 8, and:
//
//	For 4:4:4, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/1.
//	For 4:2:2, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/2.
//	For 4:2:0, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/4.
//	For 4:4:0, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/2.
type YCbCr struct {
	Y, Cb, Cr      []uint8
	YStride        int
	CStride        int
	SubsampleRatio YCbCrSubsampleRatio
	Rect           Rectangle
}

// NewYCbCr returns a new YCbCr with the given bounds and subsample ratio.

// NewYCbCr returns a new YCbCr with the
// given bounds and subsample ratio.
func NewYCbCr(r Rectangle, subsampleRatio YCbCrSubsampleRatio) *YCbCr

func (p *YCbCr) At(x, y int) color.Color

func (p *YCbCr) Bounds() Rectangle

// COffset returns the index of the first element of Cb or Cr that corresponds to
// the pixel at (x, y).

// COffset returns the index of the first
// element of Cb or Cr that corresponds to
// the pixel at (x, y).
func (p *YCbCr) COffset(x, y int) int

func (p *YCbCr) ColorModel() color.Model

func (p *YCbCr) Opaque() bool

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage returns an image representing
// the portion of the image p visible
// through r. The returned value shares
// pixels with the original image.
func (p *YCbCr) SubImage(r Rectangle) Image

func (p *YCbCr) YCbCrAt(x, y int) color.YCbCr

// YOffset returns the index of the first element of Y that corresponds to the
// pixel at (x, y).

// YOffset returns the index of the first
// element of Y that corresponds to the
// pixel at (x, y).
func (p *YCbCr) YOffset(x, y int) int

// YCbCrSubsampleRatio is the chroma subsample ratio used in a YCbCr image.

// YCbCrSubsampleRatio is the chroma
// subsample ratio used in a YCbCr image.
type YCbCrSubsampleRatio int

const (
	YCbCrSubsampleRatio444 YCbCrSubsampleRatio = iota
	YCbCrSubsampleRatio422
	YCbCrSubsampleRatio420
	YCbCrSubsampleRatio440
)

func (s YCbCrSubsampleRatio) String() string
