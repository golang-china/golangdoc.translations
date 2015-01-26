// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package color implements a basic color library.

// color 包实现了基本的颜色库。
package color

// Standard colors.

// 标准的颜色。
var (
	Black       = Gray16{0}
	White       = Gray16{0xffff}
	Transparent = Alpha16{0}
	Opaque      = Alpha16{0xffff}
)

// RGBToYCbCr converts an RGB triple to a Y'CbCr triple.

// RGBToYCbCr将RGB的三重色转换为Y'CbCr模型的三重色。
func RGBToYCbCr(r, g, b uint8) (uint8, uint8, uint8)

// YCbCrToRGB converts a Y'CbCr triple to an RGB triple.

// YCbCrToRGB将Y'CbCr上的三重色转变成RGB的三重色。
func YCbCrToRGB(y, cb, cr uint8) (uint8, uint8, uint8)

// Alpha represents an 8-bit alpha color.

// Alpha代表一个8-bit的透明度。
type Alpha struct {
	A uint8
}

func (c Alpha) RGBA() (r, g, b, a uint32)

// Alpha16 represents a 16-bit alpha color.

// Alpha16代表一个16位的透明度。
type Alpha16 struct {
	A uint16
}

func (c Alpha16) RGBA() (r, g, b, a uint32)

// Color can convert itself to alpha-premultiplied 16-bits per channel RGBA. The
// conversion may be lossy.

// Color可以将它自己转化成每个RGBA通道都预乘透明度。
// 这种转化可能是有损的。
type Color interface {
	// RGBA returns the alpha-premultiplied red, green, blue and alpha values
	// for the color. Each value ranges within [0, 0xFFFF], but is represented
	// by a uint32 so that multiplying by a blend factor up to 0xFFFF will not
	// overflow.
	RGBA() (r, g, b, a uint32)
}

// Gray represents an 8-bit grayscale color.

// Gray代表一个8-bit的灰度。
type Gray struct {
	Y uint8
}

func (c Gray) RGBA() (r, g, b, a uint32)

// Gray16 represents a 16-bit grayscale color.

// Gray16代表了一个16-bit的灰度。
type Gray16 struct {
	Y uint16
}

func (c Gray16) RGBA() (r, g, b, a uint32)

// Model can convert any Color to one from its own color model. The conversion may
// be lossy.

// Model可以在它自己的颜色模型中将一种颜色转化到另一种。
// 这种转换可能是有损的。
type Model interface {
	Convert(c Color) Color
}

// Models for the standard color types.

// 基本的颜色模型。
var (
	RGBAModel    Model = ModelFunc(rgbaModel)
	RGBA64Model  Model = ModelFunc(rgba64Model)
	NRGBAModel   Model = ModelFunc(nrgbaModel)
	NRGBA64Model Model = ModelFunc(nrgba64Model)
	AlphaModel   Model = ModelFunc(alphaModel)
	Alpha16Model Model = ModelFunc(alpha16Model)
	GrayModel    Model = ModelFunc(grayModel)
	Gray16Model  Model = ModelFunc(gray16Model)
)

// YCbCrModel is the Model for Y'CbCr colors.

// YCbCrModel是Y'CbCr颜色的模型。
var YCbCrModel Model = ModelFunc(yCbCrModel)

// ModelFunc returns a Model that invokes f to implement the conversion.

// ModelFunc返回一个Model，它可以调用f来实现转换。
func ModelFunc(f func(Color) Color) Model

// NRGBA represents a non-alpha-premultiplied 32-bit color.

// NRGBA代表一个没有32位透明度加乘的颜色。
type NRGBA struct {
	R, G, B, A uint8
}

func (c NRGBA) RGBA() (r, g, b, a uint32)

// NRGBA64 represents a non-alpha-premultiplied 64-bit color, having 16 bits for
// each of red, green, blue and alpha.

// NRGBA64代表无透明度加乘的64-bit的颜色，
// 它的每个红，绿，蓝，和透明度都是个16bit的数值。
type NRGBA64 struct {
	R, G, B, A uint16
}

func (c NRGBA64) RGBA() (r, g, b, a uint32)

// Palette is a palette of colors.

// Palette是颜色的调色板。
type Palette []Color

// Convert returns the palette color closest to c in Euclidean R,G,B space.

// Convert在Euclidean R,G,B空间中找到最接近c的调色板。
func (p Palette) Convert(c Color) Color

// Index returns the index of the palette color closest to c in Euclidean R,G,B
// space.

// Index在Euclidean R,G,B空间中找到最接近c的调色板对应的索引。
func (p Palette) Index(c Color) int

// RGBA represents a traditional 32-bit alpha-premultiplied color, having 8 bits
// for each of red, green, blue and alpha.

// RGBA代表一个传统的32位的预乘透明度的颜色，
// 它的每个红，绿，蓝，和透明度都是个8bit的数值。
type RGBA struct {
	R, G, B, A uint8
}

func (c RGBA) RGBA() (r, g, b, a uint32)

// RGBA64 represents a 64-bit alpha-premultiplied color, having 16 bits for each of
// red, green, blue and alpha.

// RGBA64代表一个64位的预乘透明度的颜色，
// 它的每个红，绿，蓝，和透明度都是个8bit的数值。
type RGBA64 struct {
	R, G, B, A uint16
}

func (c RGBA64) RGBA() (r, g, b, a uint32)

// YCbCr represents a fully opaque 24-bit Y'CbCr color, having 8 bits each for one
// luma and two chroma components.
//
// JPEG, VP8, the MPEG family and other codecs use this color model. Such codecs
// often use the terms YUV and Y'CbCr interchangeably, but strictly speaking, the
// term YUV applies only to analog video signals, and Y' (luma) is Y (luminance)
// after applying gamma correction.
//
// Conversion between RGB and Y'CbCr is lossy and there are multiple, slightly
// different formulae for converting between the two. This package follows the JFIF
// specification at http://www.w3.org/Graphics/JPEG/jfif3.pdf.

// YCbCr代表了完全不透明的24-bit的Y'CbCr的颜色，它的每个亮度和每两个色度分量是8位的。
//
// JPEG，VP8，MPEG家族和其他一些解码器使用这个颜色模式。每个解码器经常将YUV和Y'CbCr同等使用，
// 但是严格来说，YUV只是用于分析视频信号，Y' (luma)是Y
// (luminance)伽玛校正之后的结果。
//
// RGB和Y'CbCr之间的转换是有损的，并且转换的时候有许多细微的不同。这个包是遵循JFIF的说明：
// http://www.w3.org/Graphics/JPEG/jfif3.pdf。
type YCbCr struct {
	Y, Cb, Cr uint8
}

func (c YCbCr) RGBA() (uint32, uint32, uint32, uint32)
