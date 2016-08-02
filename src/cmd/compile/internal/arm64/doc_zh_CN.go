// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package arm64 // import "cmd/compile/internal/arm64"

import (
    "cmd/compile/internal/gc"
    "cmd/internal/obj"
    "cmd/internal/obj/arm64"
    "fmt"
)

const (
    LeftRdwr  uint32 = gc.LeftRead | gc.LeftWrite
    RightRdwr uint32 = gc.RightRead | gc.RightWrite
)

const (
    NREGVAR = 64 /* 32 general + 32 floating */
)

const (
    ODynam   = 1 << 0
    OAddable = 1 << 1
)

var MAXWIDTH int64 = 1 << 50

func BtoF(b uint64) int

func BtoR(b uint64) int

func Main()

//  * track register variables including external registers:
//  		reg
//  		R0
//  		R1
//  	..	...
//  		R31
//  		F0
//  		F1
//  	..	...
//  		F31
func RtoB(r int) uint64

