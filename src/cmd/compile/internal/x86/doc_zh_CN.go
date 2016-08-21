// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package x86 // import "cmd/compile/internal/x86"

import (
    "cmd/compile/internal/big"
    "cmd/compile/internal/gc"
    "cmd/internal/obj"
    "cmd/internal/obj/x86"
    "fmt"
    "os"
)

// foptoas flags
const (
    Frev  = 1 << 0
    Fpop  = 1 << 1
    Fpop2 = 1 << 2
)

const (
    NREGVAR = 16 /* 8 integer + 8 floating */
)

const (
    REGEXT = 0
)

var (
    AX               = RtoB(x86.REG_AX)
    BX               = RtoB(x86.REG_BX)
    CX               = RtoB(x86.REG_CX)
    DX               = RtoB(x86.REG_DX)
    DI               = RtoB(x86.REG_DI)
    SI               = RtoB(x86.REG_SI)
    LeftRdwr  uint32 = gc.LeftRead | gc.LeftWrite
    RightRdwr uint32 = gc.RightRead | gc.RightWrite
)

var MAXWIDTH int64 = (1 << 32) - 1

func BtoF(b uint64) int

func BtoR(b uint64) int

func FtoB(f int) uint64

func Main()

func RtoB(r int) uint64

