// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package hpack implements HPACK, a compression format for
// efficiently representing HTTP header fields in the context of HTTP/2.
//
// See http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-09
package hpack // import "internal/golang.org/x/net/http2/hpack"

import (
    "bufio"
    "bytes"
    "encoding/hex"
    "errors"
    "fmt"
    "io"
    "math/rand"
    "reflect"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "testing"
    "time"
)

// ErrInvalidHuffman is returned for errors found decoding
// Huffman-encoded strings.
var ErrInvalidHuffman = errors.New("hpack: invalid Huffman-encoded data")

// ErrStringLength is returned by Decoder.Write when the max string length
// (as configured by Decoder.SetMaxStringLength) would be violated.
var ErrStringLength = errors.New("hpack: string too long")

// A Decoder is the decoding context for incremental processing of
// header blocks.
type Decoder struct {
    dynTab dynamicTable
    emit   func(f HeaderField)

    emitEnabled bool // whether calls to emit are enabled
    maxStrLen   int  // 0 means unlimited

    // buf is the unparsed buffer. It's only written to
    // saveBuf if it was truncated in the middle of a header
    // block. Because it's usually not owned, we can only
    // process it under Write.
    buf []byte // not owned; only valid during Write

    // saveBuf is previous data passed to Write which we weren't able
    // to fully parse before. Unlike buf, we own this data.
    saveBuf bytes.Buffer
}

// A DecodingError is something the spec defines as a decoding error.
type DecodingError struct {
    Err error
}

type Encoder struct {
    dynTab dynamicTable
    // minSize is the minimum table size set by
    // SetMaxDynamicTableSize after the previous Header Table Size
    // Update.
    minSize uint32
    // maxSizeLimit is the maximum table size this encoder
    // supports. This will protect the encoder from too large
    // size.
    maxSizeLimit uint32
    // tableSizeUpdate indicates whether "Header Table Size
    // Update" is required.
    tableSizeUpdate bool
    w               io.Writer
    buf             []byte
}

// A HeaderField is a name-value pair. Both the name and value are
// treated as opaque sequences of octets.
type HeaderField struct {
    Name, Value string

    // Sensitive means that this header field should never be
    // indexed.
    Sensitive bool
}

// An InvalidIndexError is returned when an encoder references a table
// entry before the static table or after the end of the dynamic table.
type InvalidIndexError int

// AppendHuffmanString appends s, as encoded in Huffman codes, to dst
// and returns the extended buffer.
func AppendHuffmanString(dst []byte, s string) []byte

// HuffmanDecode decodes the string in v and writes the expanded
// result to w, returning the number of bytes written to w and the
// Write call's return value. At most one Write call is made.
func HuffmanDecode(w io.Writer, v []byte) (int, error)

// HuffmanDecodeToString decodes the string in v.
func HuffmanDecodeToString(v []byte) (string, error)

// HuffmanEncodeLength returns the number of bytes required to encode
// s in Huffman codes. The result is round up to byte boundary.
func HuffmanEncodeLength(s string) uint64

// NewDecoder returns a new decoder with the provided maximum dynamic
// table size. The emitFunc will be called for each valid field
// parsed, in the same goroutine as calls to Write, before Write returns.
func NewDecoder(maxDynamicTableSize uint32, emitFunc func(f HeaderField)) *Decoder

// NewEncoder returns a new Encoder which performs HPACK encoding. An
// encoded data is written to w.
func NewEncoder(w io.Writer) *Encoder

func TestAppendHpackString(t *testing.T)

func TestAppendHuffmanString(t *testing.T)

func TestAppendIndexed(t *testing.T)

func TestAppendIndexedName(t *testing.T)

func TestAppendNewName(t *testing.T)

func TestAppendTableSize(t *testing.T)

func TestAppendVarInt(t *testing.T)

// C.3 Request Examples without Huffman Coding
// http://http2.github.io/http2-spec/compression.html#rfc.section.C.3
func TestDecodeC3_NoHuffman(t *testing.T)

// C.4 Request Examples with Huffman Coding
// http://http2.github.io/http2-spec/compression.html#rfc.section.C.4
func TestDecodeC4_Huffman(t *testing.T)

// http://http2.github.io/http2-spec/compression.html#rfc.section.C.5
// "This section shows several consecutive header lists, corresponding
// to HTTP responses, on the same connection. The HTTP/2 setting
// parameter SETTINGS_HEADER_TABLE_SIZE is set to the value of 256
// octets, causing some evictions to occur."
func TestDecodeC5_ResponsesNoHuff(t *testing.T)

// http://http2.github.io/http2-spec/compression.html#rfc.section.C.6
// "This section shows the same examples as the previous section, but
// using Huffman encoding for the literal values. The HTTP/2 setting
// parameter SETTINGS_HEADER_TABLE_SIZE is set to the value of 256
// octets, causing some evictions to occur. The eviction mechanism
// uses the length of the decoded literal values, so the same
// evictions occurs as in the previous section."
func TestDecodeC6_ResponsesHuffman(t *testing.T)

func TestDecoderDecode(t *testing.T)

func TestDynamicTableAt(t *testing.T)

func TestDynamicTableSearch(t *testing.T)

func TestDynamicTableSizeEvict(t *testing.T)

func TestEmitEnabled(t *testing.T)

func TestEncoderSearchTable(t *testing.T)

func TestEncoderSetMaxDynamicTableSize(t *testing.T)

func TestEncoderSetMaxDynamicTableSizeLimit(t *testing.T)

func TestEncoderTableSizeUpdate(t *testing.T)

func TestEncoderWriteField(t *testing.T)

func TestHuffmanDecode(t *testing.T)

func TestHuffmanDecodeFuzz(t *testing.T)

// Fuzz crash, originally reported at https://github.com/bradfitz/http2/issues/56
func TestHuffmanFuzzCrash(t *testing.T)

func TestHuffmanMaxStrLen(t *testing.T)

func TestHuffmanRoundtripStress(t *testing.T)

func TestReadVarInt(t *testing.T)

func TestSaveBufLimit(t *testing.T)

func TestStaticTable(t *testing.T)

func (*Decoder) Close() error

// Decode decodes an entire block.
//
// TODO: remove this method and make it incremental later? This is
// easier for debugging now.
func (*Decoder) DecodeFull(p []byte) ([]HeaderField, error)

// EmitEnabled reports whether calls to the emitFunc provided to NewDecoder
// are currently enabled. The default is true.
func (*Decoder) EmitEnabled() bool

// SetAllowedMaxDynamicTableSize sets the upper bound that the encoded
// stream (via dynamic table size updates) may set the maximum size
// to.
func (*Decoder) SetAllowedMaxDynamicTableSize(v uint32)

// SetEmitEnabled controls whether the emitFunc provided to NewDecoder
// should be called. The default is true.
//
// This facility exists to let servers enforce MAX_HEADER_LIST_SIZE
// while still decoding and keeping in-sync with decoder state, but
// without doing unnecessary decompression or generating unnecessary
// garbage for header fields past the limit.
func (*Decoder) SetEmitEnabled(v bool)

// SetEmitFunc changes the callback used when new header fields
// are decoded.
// It must be non-nil. It does not affect EmitEnabled.
func (*Decoder) SetEmitFunc(emitFunc func(f HeaderField))

func (*Decoder) SetMaxDynamicTableSize(v uint32)

// SetMaxStringLength sets the maximum size of a HeaderField name or
// value string. If a string exceeds this length (even after any
// decompression), Write will return ErrStringLength.
// A value of 0 means unlimited and is the default from NewDecoder.
func (*Decoder) SetMaxStringLength(n int)

func (*Decoder) Write(p []byte) (n int, err error)

// SetMaxDynamicTableSize changes the dynamic header table size to v.
// The actual size is bounded by the value passed to
// SetMaxDynamicTableSizeLimit.
func (*Encoder) SetMaxDynamicTableSize(v uint32)

// SetMaxDynamicTableSizeLimit changes the maximum value that can be
// specified in SetMaxDynamicTableSize to v. By default, it is set to
// 4096, which is the same size of the default dynamic header table
// size described in HPACK specification. If the current maximum
// dynamic header table size is strictly greater than v, "Header Table
// Size Update" will be done in the next WriteField call and the
// maximum dynamic header table size is truncated to v.
func (*Encoder) SetMaxDynamicTableSizeLimit(v uint32)

// WriteField encodes f into a single Write to e's underlying Writer.
// This function may also produce bytes for "Header Table Size Update"
// if necessary.  If produced, it is done before encoding f.
func (*Encoder) WriteField(f HeaderField) error

func (DecodingError) Error() string

func (HeaderField) String() string

func (InvalidIndexError) Error() string

