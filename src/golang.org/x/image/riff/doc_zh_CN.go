// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package riff implements the Resource Interchange File Format, used by media
// formats such as AVI, WAVE and WEBP.
//
// A RIFF stream contains a sequence of chunks. Each chunk consists of an 8-byte
// header (containing a 4-byte chunk type and a 4-byte chunk length), the chunk
// data (presented as an io.Reader), and some padding bytes.
//
// A detailed description of the format is at
// http://www.tactilemedia.com/info/MCI_Control_Info.html
package riff

// LIST is the "LIST" FourCC.
var LIST = FourCC{'L', 'I', 'S', 'T'}

// FourCC is a four character code.
type FourCC [4]byte

// NewListReader returns a LIST chunk's list type, such as "movi" or "wavl", and
// its chunks as a *Reader.
func NewListReader(chunkLen uint32, chunkData io.Reader) (listType FourCC, data *Reader, err error)

// NewReader returns the RIFF stream's form type, such as "AVI " or "WAVE", and its
// chunks as a *Reader.
func NewReader(r io.Reader) (formType FourCC, data *Reader, err error)

// Reader reads chunks from an underlying io.Reader.
type Reader struct {
	// contains filtered or unexported fields
}

// Next returns the next chunk's ID, length and data. It returns io.EOF if there
// are no more chunks. The io.Reader returned becomes stale after the next Next
// call, and should no longer be used.
//
// It is valid to call Next even if all of the previous chunk's data has not been
// read.
func (z *Reader) Next() (chunkID FourCC, chunkLen uint32, chunkData io.Reader, err error)
