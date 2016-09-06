// +build ingore

package x86

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
	DWARFREGSP = 4
	DWARFREGLR = 8
)

const (
	MaxAlign  = 32 // max data alignment
	MinAlign  = 1  // min data alignment
	FuncAlign = 16
)

func Main()

