// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package png implements a PNG image decoder and encoder.
//
// The PNG specification is at http://www.w3.org/TR/PNG/.

// png 包实现了PNG图像的编码和解码.
//
// PNG的具体说明在http://www.w3.org/TR/PNG/。
package png

// Decode reads a PNG image from r and returns it as an image.Image. The type of
// Image returned depends on the PNG contents.
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig returns the color model and dimensions of a PNG image without
// decoding the entire image.

// DecodeConfig返回颜色模型，没有解码整个图像，获得了PNG图片的尺寸。
func DecodeConfig(r io.Reader) (image.Config, error)

// Encode writes the Image m to w in PNG format. Any Image may be encoded, but
// images that are not image.NRGBA might be encoded lossily.

// Encode将图片m以PNG的格式写到w中。任何图片都可以被编码，但是哪些不是 image.NRGBA
// 的图片编码可能是有损的。
func Encode(w io.Writer, m image.Image) error

type CompressionLevel int

const (
	DefaultCompression CompressionLevel = 0
	NoCompression      CompressionLevel = -1
	BestSpeed          CompressionLevel = -2
	BestCompression    CompressionLevel = -3
)

// Encoder configures encoding PNG images.
type Encoder struct {
	CompressionLevel CompressionLevel
}

// Encode writes the Image m to w in PNG format.

// Encode 将图像 m 以 PNG 格式写入 w。
func (enc *Encoder) Encode(w io.Writer, m image.Image) error

// A FormatError reports that the input is not a valid PNG.

// FormatError会提示输入并不是一个合法的PNG。
type FormatError string

func (e FormatError) Error() string

// An UnsupportedError reports that the input uses a valid but unimplemented PNG
// feature.

// UnsupportedError会提示输入使用一个合法的，但是未实现的PNG特性。
type UnsupportedError string

func (e UnsupportedError) Error() string
