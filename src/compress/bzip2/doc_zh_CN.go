// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package bzip2 implements bzip2 decompression.

// bzip2 包实现 bzip2 的解压缩.
package bzip2

import (
    "bufio"
    "io"
    "sort"
)

// A StructuralError is returned when the bzip2 data is found to be
// syntactically invalid.

// StructuralError 表示一个错误的 bzip2 数据结构.
type StructuralError string

// NewReader returns an io.Reader which decompresses bzip2 data from r.
// If r does not also implement io.ByteReader,
// the decompressor may read more data than necessary from r.

// NewReader 返回一个从 r 读取 bzip2 压缩数据并解压缩后返回给调用者的 io.Reader.
func NewReader(r io.Reader) io.Reader

func (StructuralError) Error() string

