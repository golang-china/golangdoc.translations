// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package errors implements functions to manipulate errors.

// error 包实现了用于错误处理的函数.
package errors

// New returns an error that formats as the given text.

// New 返回一个按给定文本格式化的错误。
func New(text string) error

