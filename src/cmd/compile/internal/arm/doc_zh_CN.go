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

// reg
//  		R0
//  		R1
//  		...
//  		R10
//  	12  R12
//  *
//  		reg
//  		F2
//  		F3
//  		...
//  		F15
func RtoB(r int) uint64

