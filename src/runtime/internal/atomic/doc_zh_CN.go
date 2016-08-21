// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package atomic // import "runtime/internal/atomic"

import (
    "runtime/internal/sys"
    "unsafe"
)

// go:noescape
func And8(ptr *uint8, val uint8)

// go:noescape
func Cas(ptr *uint32, old, new uint32) bool

// go:noescape
func Cas64(ptr *uint64, old, new uint64) bool

// NO go:noescape annotation; see atomic_pointer.go.
func Casp1(ptr *unsafe.Pointer, old, new unsafe.Pointer) bool

// go:noescape
func Casuintptr(ptr *uintptr, old, new uintptr) bool

// go:nosplit
// go:noinline
func Load(ptr *uint32) uint32

// go:noescape
func Load64(ptr *uint64) uint64

// go:noescape
func Loadint64(ptr *int64) int64

// go:nosplit
// go:noinline
func Loadp(ptr unsafe.Pointer) unsafe.Pointer

// go:noescape
func Loaduint(ptr *uint) uint

// go:noescape
func Loaduintptr(ptr *uintptr) uintptr

// go:noescape
func Or8(ptr *uint8, val uint8)

// go:noescape
func Store(ptr *uint32, val uint32)

// go:noescape
func Store64(ptr *uint64, val uint64)

// NO go:noescape annotation; see atomic_pointer.go.
func Storep1(ptr unsafe.Pointer, val unsafe.Pointer)

// go:noescape
func Storeuintptr(ptr *uintptr, new uintptr)

// go:noescape
func Xadd(ptr *uint32, delta int32) uint32

// go:nosplit
func Xadd64(ptr *uint64, delta int64) uint64

// go:noescape
func Xaddint64(ptr *int64, delta int64) int64

// go:noescape
func Xadduintptr(ptr *uintptr, delta uintptr) uintptr

// go:noescape
func Xchg(ptr *uint32, new uint32) uint32

// go:nosplit
func Xchg64(ptr *uint64, new uint64) uint64

// go:noescape
func Xchguintptr(ptr *uintptr, new uintptr) uintptr

