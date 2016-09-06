// +build ingore

package mips64

import (
	"cmd/internal/obj"
	"cmd/internal/sys"
	"cmd/link/internal/ld"
	"fmt"
	"log"
)

//  Used by ../internal/ld/dwarf.go

// Used by ../internal/ld/dwarf.go
const (
	DWARFREGSP = 29
	DWARFREGLR = 31
)

const (
	MaxAlign  = 32 // max data alignment
	MinAlign  = 1  // min data alignment
	FuncAlign = 8
)

func Main()

