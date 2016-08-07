// Copyright The Go Authors. All rights reserved.
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

// image实现了基本的2D图片库。
//
// 基本接口叫作Image。图片的色彩定义在image/color包。
//
// Image接口可以通过调用如NewRGBA和NewPaletted函数等获得；也可以通过调用Decode函
// 数解码包含GIF、JPEG或PNG格式图像数据的输入流获得。解码任何具体图像类型之前都
// 必须注册对应类型的解码函数。注册过程一般是作为包初始化的副作用，放在包的init
// 函数里。因此，要解码PNG图像，只需在程序的main包里嵌入如下代码：
//
//     import _ "image/png"
//
// _表示导入包但不使用包中的变量/函数/类型，只是为了包初始化函数的副作用。
//
// 参见http://golang.org/doc/articles/image_package.html
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

// ErrFormat说明解码时遇到了未知的格式。
var ErrFormat = errors.New("image: unknown format")

// ZP is the zero Point.
var ZP Point

// ZR is the zero Rectangle.
var ZR Rectangle

// Alpha is an in-memory image whose At method returns color.Alpha values.

// Alpha类型代表一幅内存中的图像，其At方法返回color.Alpha类型的值。
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

// Alpha16类型代表一幅内存中的图像，其At方法返回color.Alpha16类型的值。
type Alpha16 struct {
    // Pix holds the image's pixels, as alpha values in big-endian format. The pixel at
    // (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
    Pix []uint8
    // Stride is the Pix stride (in bytes) between vertically adjacent pixels.
    Stride int
    // Rect is the image's bounds.
    Rect Rectangle
}

// Config holds an image's color model and dimensions.

// Config保管图像的色彩模型和尺寸信息。
type Config struct {
    ColorModel    color.Model
    Width, Height int
}

// Gray is an in-memory image whose At method returns color.Gray values.

// Gray类型代表一幅内存中的图像，其At方法返回color.Gray类型的值。
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

// Gray16类型代表一幅内存中的图像，其At方法返回color.Gray16类型的值。
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

// Image接口表示一个采用某色彩模型的颜色构成的有限矩形网格（即一幅图像）。
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

// NRGBA类型代表一幅内存中的图像，其At方法返回color.NRGBA类型的值。
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

// NRGBA64类型代表一幅内存中的图像，其At方法返回color.NRGBA64类型的值。
type NRGBA64 struct {
    // Pix holds the image's pixels, in R, G, B, A order and big-endian format. The pixel at
    // (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*8].
    Pix []uint8
    // Stride is the Pix stride (in bytes) between vertically adjacent pixels.
    Stride int
    // Rect is the image's bounds.
    Rect Rectangle
}

// Paletted is an in-memory image of uint8 indices into a given palette.

// Paletted类型是一幅采用uint8类型索引调色板的内存中的图像。
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

// PalettedImage接口代表一幅图像，它的像素可能来自一个有限的调色板。
//
// 如果有对象m满足PalettedImage接口，且m.ColorModel()返回的color.Model接口底层为
// 一个Palette类型值（记为p），则m.At(x, y)返回值应等于p[m.ColorIndexAt(x, y)]。
// 如果m的色彩模型不是Palette，则ColorIndexAt的行为是不确定的。
type PalettedImage interface {
    // ColorIndexAt returns the palette index of the pixel at (x, y).
    ColorIndexAt(x, y int) uint8
    Image
}

// A Point is an X, Y coordinate pair. The axes increase right and down.

// Point是X,
// Y坐标对。坐标轴是向右（X）向下（Y）的。既可以表示点，也可以表示向量。
//
//     var ZP Point
//
// ZP是原点。
type Point struct {
    X, Y int
}

// RGBA is an in-memory image whose At method returns color.RGBA values.

// RGBA类型代表一幅内存中的图像，其At方法返回color.RGBA类型的值。
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

// RGBA64类型代表一幅内存中的图像，其At方法返回color.RGBA64类型的值
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

// Rectangle代表一个矩形。该矩形包含所有满足Min.X <= X < Max.X且Min.Y <= Y <
// Max.Y的点。如果两个字段满足Min.X <= Max.X且Min.Y <= Max.Y，就称该实例为规范格
// 式的。矩形的方法，当输入是规范格式时，总是返回规范格式的输出。
//
//     var ZR Rectangle
//
// ZR是矩形的零值。
type Rectangle struct {
    Min, Max Point
}

// Uniform is an infinite-sized Image of uniform color.
// It implements the color.Color, color.Model, and Image interfaces.

// Uniform类型代表一块面积无限大的具有同一色彩的图像。它实现了color.Color、
// color.Model和Image等接口。
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

// YcbCr代表采用Y'CbCr色彩模型的一幅内存中的图像。每个像素都对应一个Y采样，但每
// 个Cb/Cr采样对应多个像素。Ystride是两个垂直相邻的像素之间的Y组分的索引增量。
// CStride是两个映射到单独的色度采样的垂直相邻的像素之间的Cb/Cr组分的索引增量。
// 虽然不作绝对要求，但Ystride字段和len(Y)一般应为8的倍数，并且：
//
//     For 4:4:4, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/1.
//     For 4:2:2, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/2.
//     For 4:2:0, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/4.
//     For 4:4:0, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/2.
type YCbCr struct {
    Y, Cb, Cr      []uint8
    YStride        int
    CStride        int
    SubsampleRatio YCbCrSubsampleRatio
    Rect           Rectangle
}

// YCbCrSubsampleRatio is the chroma subsample ratio used in a YCbCr image.

// YcbCrSubsampleRatio是YCbCr图像的色度二次采样比率。
//
//     const (
//         YCbCrSubsampleRatio444 YCbCrSubsampleRatio = iota
//         YCbCrSubsampleRatio422
//         YCbCrSubsampleRatio420
//         YCbCrSubsampleRatio440
//     )
type YCbCrSubsampleRatio int

// Decode decodes an image that has been encoded in a registered format.
// The string returned is the format name used during format registration.
// Format registration is typically done by an init function in the codec-
// specific package.

// DecodeConfig函数解码并返回一个采用某种已注册格式编码的图像。字符串返回值是该
// 格式注册时的名字。格式一般是在该编码格式的包的init函数中注册的。
func Decode(r io.Reader) (Image, string, error)

// DecodeConfig decodes the color model and dimensions of an image that has
// been encoded in a registered format. The string returned is the format name
// used during format registration. Format registration is typically done by
// an init function in the codec-specific package.

// DecodeConfig函数解码并返回一个采用某种已注册格式编码的图像的色彩模型和尺寸。
// 字符串返回值是该格式注册时的名字。格式一般是在该编码格式的包的init函数中注册
// 的。
func DecodeConfig(r io.Reader) (Config, string, error)

// NewAlpha returns a new Alpha image with the given bounds.

// NewAlpha函数创建并返回一个具有指定宽度和高度的Alpha。
func NewAlpha(r Rectangle) *Alpha

// NewAlpha16 returns a new Alpha16 image with the given bounds.

// NewAlpha16函数创建并返回一个具有指定范围的Alpha16。
func NewAlpha16(r Rectangle) *Alpha16

// NewGray returns a new Gray image with the given bounds.

// NewGray函数创建并返回一个具有指定范围的Gray。
func NewGray(r Rectangle) *Gray

// NewGray16 returns a new Gray16 image with the given bounds.

// NewGray16函数创建并返回一个具有指定范围的Gray16。
func NewGray16(r Rectangle) *Gray16

// NewNRGBA returns a new NRGBA image with the given bounds.

// NewNRGBA函数创建并返回一个具有指定范围的NRGBA。
func NewNRGBA(r Rectangle) *NRGBA

// NewNRGBA64 returns a new NRGBA64 image with the given bounds.

// NewNRGBA64函数创建并返回一个具有指定范围的NRGBA64。
func NewNRGBA64(r Rectangle) *NRGBA64

// NewPaletted returns a new Paletted image with the given width, height and
// palette.

// NewPaletted函数创建并返回一个具有指定范围、调色板的Paletted。
func NewPaletted(r Rectangle, p color.Palette) *Paletted

// NewRGBA returns a new RGBA image with the given bounds.

// NewRGBA函数创建并返回一个具有指定范围的RGBA。
func NewRGBA(r Rectangle) *RGBA

// NewRGBA64 returns a new RGBA64 image with the given bounds.

// NewRGBA64函数创建并返回一个具有指定范围的RGBA64
func NewRGBA64(r Rectangle) *RGBA64

func NewUniform(c color.Color) *Uniform

// NewYCbCr returns a new YCbCr image with the given bounds and subsample
// ratio.

// NewYCbCr函数创建并返回一个具有指定宽度、高度和二次采样率的YcbCr。
func NewYCbCr(r Rectangle, subsampleRatio YCbCrSubsampleRatio) *YCbCr

// Pt is shorthand for Point{X, Y}.

// 返回Point{X , Y}
func Pt(X, Y int) Point

// Rect is shorthand for Rectangle{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.

// 返回一个矩形Rectangle{Pt(x0, y0), Pt(x1, y1)}。
func Rect(x0, y0, x1, y1 int) Rectangle

// RegisterFormat registers an image format for use by Decode.
// Name is the name of the format, like "jpeg" or "png".
// Magic is the magic prefix that identifies the format's encoding. The magic
// string can contain "?" wildcards that each match any one byte.
// Decode is the function that decodes the encoded image.
// DecodeConfig is the function that decodes just its configuration.

// RegisterFormat注册一个供Decode函数使用的图片格式。name是格式的名字，如"jpeg"
// 或"png"；magic是该格式编码的魔术前缀，该字符串可以包含"?"通配符，每个通配符匹
// 配一个字节；decode函数用于解码图片；decodeConfig函数只解码图片的配置。
func RegisterFormat(name, magic string, decode func(io.Reader) (Image, error), decodeConfig func(io.Reader) (Config, error))

func (*Alpha) AlphaAt(x, y int) color.Alpha

func (*Alpha) At(x, y int) color.Color

func (*Alpha) Bounds() Rectangle

func (*Alpha) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*Alpha) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*Alpha) PixOffset(x, y int) int

func (*Alpha) Set(x, y int, c color.Color)

func (*Alpha) SetAlpha(x, y int, c color.Alpha)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*Alpha) SubImage(r Rectangle) Image

func (*Alpha16) Alpha16At(x, y int) color.Alpha16

func (*Alpha16) At(x, y int) color.Color

func (*Alpha16) Bounds() Rectangle

func (*Alpha16) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*Alpha16) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*Alpha16) PixOffset(x, y int) int

func (*Alpha16) Set(x, y int, c color.Color)

func (*Alpha16) SetAlpha16(x, y int, c color.Alpha16)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*Alpha16) SubImage(r Rectangle) Image

func (*Gray) At(x, y int) color.Color

func (*Gray) Bounds() Rectangle

func (*Gray) ColorModel() color.Model

func (*Gray) GrayAt(x, y int) color.Gray

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*Gray) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*Gray) PixOffset(x, y int) int

func (*Gray) Set(x, y int, c color.Color)

func (*Gray) SetGray(x, y int, c color.Gray)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*Gray) SubImage(r Rectangle) Image

func (*Gray16) At(x, y int) color.Color

func (*Gray16) Bounds() Rectangle

func (*Gray16) ColorModel() color.Model

func (*Gray16) Gray16At(x, y int) color.Gray16

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*Gray16) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*Gray16) PixOffset(x, y int) int

func (*Gray16) Set(x, y int, c color.Color)

func (*Gray16) SetGray16(x, y int, c color.Gray16)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*Gray16) SubImage(r Rectangle) Image

func (*NRGBA) At(x, y int) color.Color

func (*NRGBA) Bounds() Rectangle

func (*NRGBA) ColorModel() color.Model

func (*NRGBA) NRGBAAt(x, y int) color.NRGBA

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*NRGBA) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*NRGBA) PixOffset(x, y int) int

func (*NRGBA) Set(x, y int, c color.Color)

func (*NRGBA) SetNRGBA(x, y int, c color.NRGBA)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*NRGBA) SubImage(r Rectangle) Image

func (*NRGBA64) At(x, y int) color.Color

func (*NRGBA64) Bounds() Rectangle

func (*NRGBA64) ColorModel() color.Model

func (*NRGBA64) NRGBA64At(x, y int) color.NRGBA64

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*NRGBA64) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*NRGBA64) PixOffset(x, y int) int

func (*NRGBA64) Set(x, y int, c color.Color)

func (*NRGBA64) SetNRGBA64(x, y int, c color.NRGBA64)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*NRGBA64) SubImage(r Rectangle) Image

func (*Paletted) At(x, y int) color.Color

func (*Paletted) Bounds() Rectangle

func (*Paletted) ColorIndexAt(x, y int) uint8

func (*Paletted) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*Paletted) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*Paletted) PixOffset(x, y int) int

func (*Paletted) Set(x, y int, c color.Color)

func (*Paletted) SetColorIndex(x, y int, index uint8)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*Paletted) SubImage(r Rectangle) Image

func (*RGBA) At(x, y int) color.Color

func (*RGBA) Bounds() Rectangle

func (*RGBA) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*RGBA) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*RGBA) PixOffset(x, y int) int

func (*RGBA) RGBAAt(x, y int) color.RGBA

func (*RGBA) Set(x, y int, c color.Color)

func (*RGBA) SetRGBA(x, y int, c color.RGBA)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*RGBA) SubImage(r Rectangle) Image

func (*RGBA64) At(x, y int) color.Color

func (*RGBA64) Bounds() Rectangle

func (*RGBA64) ColorModel() color.Model

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告图像是否是完全不透明的。
func (*RGBA64) Opaque() bool

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).

// PixOffset方法返回像素(x,
// y)的数据起始位置在Pix字段的偏移量/索引。
func (*RGBA64) PixOffset(x, y int) int

func (*RGBA64) RGBA64At(x, y int) color.RGBA64

func (*RGBA64) Set(x, y int, c color.Color)

func (*RGBA64) SetRGBA64(x, y int, c color.RGBA64)

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*RGBA64) SubImage(r Rectangle) Image

func (*Uniform) At(x, y int) color.Color

func (*Uniform) Bounds() Rectangle

func (*Uniform) ColorModel() color.Model

func (*Uniform) Convert(color.Color) color.Color

// Opaque scans the entire image and reports whether it is fully opaque.

// Opaque方法扫描整个图像并报告该图像是否是完全不透明的。
func (*Uniform) Opaque() bool

func (*Uniform) RGBA() (r, g, b, a uint32)

func (*YCbCr) At(x, y int) color.Color

func (*YCbCr) Bounds() Rectangle

// COffset returns the index of the first element of Cb or Cr that corresponds
// to the pixel at (x, y).

// 像素(X,
// Y)的Cb或Cr（色度）组分的数据起始位置在Cb/Cr字段的偏移量/索引。
func (*YCbCr) COffset(x, y int) int

func (*YCbCr) ColorModel() color.Model

func (*YCbCr) Opaque() bool

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.

// SubImage方法返回代表原图像一部分（r的范围）的新图像。返回值和原图像的像素数据
// 是共用的。
func (*YCbCr) SubImage(r Rectangle) Image

func (*YCbCr) YCbCrAt(x, y int) color.YCbCr

// YOffset returns the index of the first element of Y that corresponds to
// the pixel at (x, y).

// 像素(X,
// Y)的Y（亮度）组分的数据起始位置在Y字段的偏移量/索引。
func (*YCbCr) YOffset(x, y int) int

// Add returns the vector p+q.

// 返回点Point{p.X+q.X, p.Y+q.Y}
func (Point) Add(q Point) Point

// Div returns the vector p/k.

// 返回点Point{p.X/k, p.Y/k }
func (Point) Div(k int) Point

// Eq reports whether p and q are equal.

// 报告p和q是否相同。
func (Point) Eq(q Point) bool

// In reports whether p is in r.

// 报告p是否在r范围内。
func (Point) In(r Rectangle) bool

// Mod returns the point q in r such that p.X-q.X is a multiple of r's width
// and p.Y-q.Y is a multiple of r's height.

// 返回r范围内的某点q，满足p.X-q.X是r宽度的倍数，p.Y-q.Y是r高度的倍数。
func (Point) Mod(r Rectangle) Point

// Mul returns the vector p*k.

// 返回点Point{p.X*k, p.Y*k}
func (Point) Mul(k int) Point

// String returns a string representation of p like "(3,4)".

// 返回p的字符串表示。格式为"(3,4)"
func (Point) String() string

// Sub returns the vector p-q.

// 返回点Point{p.X-q.X, p.Y-q.Y}
func (Point) Sub(q Point) Point

// Add returns the rectangle r translated by p.

// 返回矩形按p（作为向量）平移后的新矩形。
func (Rectangle) Add(p Point) Rectangle

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.

// 返回矩形的规范版本（左上&右下），方法必要时会交换坐标的最大值和最小值。
func (Rectangle) Canon() Rectangle

// Dx returns r's width.

// 返回r的宽度。
func (Rectangle) Dx() int

// Dy returns r's height.

// 返回r的高度。
func (Rectangle) Dy() int

// Empty reports whether the rectangle contains no points.

// 报告矩形是否为空矩形。（即内部不包含点的矩形）
func (Rectangle) Empty() bool

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.

// 报告两个矩形是否相同。
func (Rectangle) Eq(s Rectangle) bool

// In reports whether every point in r is in s.

// 如果r包含的所有点都在s内，则返回真；否则返回假。
func (Rectangle) In(s Rectangle) bool

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.

// 返回去掉矩形四周宽度n的框的矩形，n可为负数。如果n过大将返回靠近r中心位置的空
// 矩形。
func (Rectangle) Inset(n int) Rectangle

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.

// 返回两个矩形的交集矩形（同时被r和s包含的最大矩形）；如果r和s没有重叠会返回
// Rectangle零值。
func (Rectangle) Intersect(s Rectangle) Rectangle

// Overlaps reports whether r and s have a non-empty intersection.

// 如果r和s有非空的交集，则返回真；否则返回假。
func (Rectangle) Overlaps(s Rectangle) bool

// Size returns r's width and height.

// 返回r的宽度w和高度h构成的点Point{w, h}。
func (Rectangle) Size() Point

// String returns a string representation of r like "(3,4)-(6,5)".

// 返回矩形的字符串表示，格式为"(3,4)-(6,5)"。
func (Rectangle) String() string

// Sub returns the rectangle r translated by -p.

// 返回矩形按p（作为向量）反向平移后的新矩形。
func (Rectangle) Sub(p Point) Rectangle

// Union returns the smallest rectangle that contains both r and s.

// 返回同时包含r和s的最小矩形。
func (Rectangle) Union(s Rectangle) Rectangle

func (YCbCrSubsampleRatio) String() string

