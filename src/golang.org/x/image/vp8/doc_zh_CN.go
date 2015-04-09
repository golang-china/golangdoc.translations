// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package vp8 implements a decoder for the VP8 lossy image format.
//
// The VP8 specification is RFC 6386.
package vp8

// Decoder decodes VP8 bitstreams into frames. Decoding one frame consists of
// calling Init, DecodeFrameHeader and then DecodeFrame in that order. A Decoder
// can be re-used to decode multiple frames.
type Decoder struct {
	// contains filtered or unexported fields
}

// NewDecoder returns a new Decoder.
func NewDecoder() *Decoder

// DecodeFrame decodes the frame and returns it as an YCbCr image. The image's
// contents are valid up until the next call to Decoder.Init.
func (d *Decoder) DecodeFrame() (*image.YCbCr, error)

// DecodeFrameHeader decodes the frame header.
func (d *Decoder) DecodeFrameHeader() (fh FrameHeader, err error)

// Init initializes the decoder to read at most n bytes from r.
func (d *Decoder) Init(r io.Reader, n int)

// FrameHeader is a frame header, as specified in section 9.1.
type FrameHeader struct {
	KeyFrame          bool
	VersionNumber     uint8
	ShowFrame         bool
	FirstPartitionLen uint32
	Width             int
	Height            int
	XScale            uint8
	YScale            uint8
}
