// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package race contains helper functions for manually instrumenting code for the
// race detector.
//
// The runtime package intentionally exports these functions only in the race build;
// this package exports them unconditionally but without the "race" build tag they
// are no-ops.
package race // import "internal/race"

import (
    "runtime"
    "unsafe"
)

const Enabled = false

const Enabled = true

func Acquire(addr unsafe.Pointer)

func Disable()

func Enable()

func Read(addr unsafe.Pointer)

func ReadRange(addr unsafe.Pointer, len int)

func Release(addr unsafe.Pointer)

func ReleaseMerge(addr unsafe.Pointer)

func Write(addr unsafe.Pointer)

func WriteRange(addr unsafe.Pointer, len int)

