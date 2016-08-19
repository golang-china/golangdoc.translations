// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package arm

import (
    "cmd/compile/internal/gc"
    "cmd/compile/internal/ssa"
    "cmd/internal/obj"
    "cmd/internal/obj/arm"
    "fmt"
)


const (
	Falsecond = iota
	Truecond
	Delbranch
	Keepbranch
)



const (
	Join = iota
	Split
	End
	Branch
	Setcond
	Toolong
)



const (
	NREGVAR = 32
)



const (
	ODynam = 1 << 0
	OPtrto = 1 << 1
)



const (
	RightRdwr = gc.RightRead | gc.RightWrite
)



type Joininfo struct {
	start *gc.Flow
	last  *gc.Flow
	end   *gc.Flow
	len   int
}


func BtoF(b uint64) int

func BtoR(b uint64) int

func Main()

//  *	bit	reg
//  *	0	R0
//  *	1	R1
//  *	...	...
//  *	10	R10
//  *	12  R12
//  *
//  *	bit	reg
//  *	18	F2
//  *	19	F3
//  *	...	...
//  *	31	F15
func RtoB(r int) uint64

