// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package jpeg implements a JPEG image decoder and encoder.
//
// JPEG is defined in ITU-T T.81: http://www.w3.org/Graphics/JPEG/itu-t81.pdf.

// jpeg包实现了jpeg格式图像的编解码。JPEG格式参见http://www.w3.org/Graphics/JPEG/itu-t81.pdf
package jpeg

// DefaultQuality is the default quality encoding parameter.

// DefaultQuality是默认的编码质量参数。
const DefaultQuality = 75

// Decode reads a JPEG image from r and returns it as an image.Image.

// 从r读取一幅jpeg格式的图像并解码返回该图像。
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig returns the color model and dimensions of a JPEG image without
// decoding the entire image.

// 返回JPEG图像的色彩模型和尺寸；函数不会解码整个图像。
func DecodeConfig(r io.Reader) (image.Config, error)

// Encode writes the Image m to w in JPEG 4:2:0 baseline format with the given
// options. Default parameters are used if a nil *Options is passed.

// Encode函数将采用JPEG
// 4:2:0基线格式和指定的编码质量将图像写入w。如果o为nil将使用DefaultQuality。
func Encode(w io.Writer, m image.Image, o *Options) error

// A FormatError reports that the input is not a valid JPEG.

// 当输入流不是合法的jpeg格式图像时，就会返回FormatError类型的错误。
type FormatError string

func (e FormatError) Error() string

// Options are the encoding parameters. Quality ranges from 1 to 100 inclusive,
// higher is better.

// Options是编码质量参数。取值范围[1,100]，越大图像编码质量越高。
type Options struct {
	Quality int
}

// Reader is deprecated.

// 如果提供的io.Reader接口没有ReadByte方法，Decode函数会为该接口附加一个缓冲。
type Reader interface {
	io.ByteReader
	io.Reader
}

// An UnsupportedError reports that the input uses a valid but unimplemented JPEG
// feature.

// 当输入流使用了合法但尚不支持的jpeg特性的时候，就会返回UnsupportedError类型的错误。
type UnsupportedError string

func (e UnsupportedError) Error() string
