// +build ingore

package s390x

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/s390x"
	"fmt"
)

const (
	NREGVAR = 32 /* 16 general + 16 floating*/
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
//  *	...	...
//  *	15	R15
//  *	16+0	F0
//  *	16+1	F1
//  *	...	...
//  *	16+15	F15
func RtoB(r int) uint64

