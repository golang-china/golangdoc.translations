// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

const (
	PROCESSOR_ARCHITECTURE_AMD64 = 9
	PROCESSOR_ARCHITECTURE_INTEL = 0
)

const (
	CheckExit = 1 << iota
	ShowOutput
	Background
)
