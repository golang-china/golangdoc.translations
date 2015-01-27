// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package gif implements a GIF image decoder and encoder.
//
// The GIF specification is at http://www.w3.org/Graphics/GIF/spec-gif89a.txt.

// gif 包实现了GIF图片的解码.
//
// GIF的说明文档在 http://www.w3.org/Graphics/GIF/spec-gif89a.txt。
package gif

// Decode reads a GIF image from r and returns the first embedded image as an
// image.Image.

// Decode从r中读取一个GIF图像，然后返回的image.Image是第一个嵌入的图。
func Decode(r io.Reader) (image.Image, error)

// DecodeConfig returns the global color model and dimensions of a GIF image
// without decoding the entire image.

// DecodeConfig不需要解码整个图像就可以返回全局的颜色模型和GIF图片的尺寸。
func DecodeConfig(r io.Reader) (image.Config, error)

// Encode writes the Image m to w in GIF format.
func Encode(w io.Writer, m image.Image, o *Options) error

// EncodeAll writes the images in g to w in GIF format with the given loop count
// and delay between frames.
func EncodeAll(w io.Writer, g *GIF) error

// GIF represents the possibly multiple images stored in a GIF file.

// GIF代表一个GIF文件上的多个图像。
type GIF struct {
	Image     []*image.Paletted // The successive images.
	Delay     []int             // The successive delay times, one per frame, in 100ths of a second.
	LoopCount int               // The loop count.
}

// DecodeAll reads a GIF image from r and returns the sequential frames and timing
// information.

// DecodeAll
// 从r上读取一个GIF图片，并且返回顺序的帧和时间信息。
func DecodeAll(r io.Reader) (*GIF, error)

// Options are the encoding parameters.
type Options struct {
	// NumColors is the maximum number of colors used in the image.
	// It ranges from 1 to 256.
	NumColors int

	// Quantizer is used to produce a palette with size NumColors.
	// palette.Plan9 is used in place of a nil Quantizer.
	Quantizer draw.Quantizer

	// Drawer is used to convert the source image to the desired palette.
	// draw.FloydSteinberg is used in place of a nil Drawer.
	Drawer draw.Drawer
}
