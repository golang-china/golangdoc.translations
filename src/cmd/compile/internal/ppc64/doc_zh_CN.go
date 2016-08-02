// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package ppc64 // import "cmd/compile/internal/ppc64"

import (
    "cmd/compile/internal/big"
    "cmd/compile/internal/gc"
    "cmd/internal/obj"
    "cmd/internal/obj/ppc64"
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

// Many Power ISA arithmetic and logical instructions come in four
// standard variants.  These bits let us map between variants.
const (
    V_CC = 1 << 0 // xCC (affect CR field 0 flags)
    V_V  = 1 << 1 // xV (affect SO and OV flags)
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

