// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package runtime


type Pointer uintptr // not really; filled in by compiler


func Alignof(any) uintptr

// return types here are ignored; see unsafe.go
func Offsetof(any) uintptr

func Sizeof(any) uintptr

