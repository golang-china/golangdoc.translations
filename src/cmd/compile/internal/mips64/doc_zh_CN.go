// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package mips64

import (
    "cmd/compile/internal/big"
    "cmd/compile/internal/gc"
    "cmd/internal/obj"
    "cmd/internal/obj/mips"
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


func BtoF(b uint64) int

func BtoR(b uint64) int

func Main()

//  * track register variables including external registers:
//  *	bit	reg
//  *	0	R0
//  *	1	R1
//  *	...	...
//  *	31	R31
//  *	32+0	F0
//  *	32+1	F1
//  *	...	...
//  *	32+31	F31
func RtoB(r int) uint64

