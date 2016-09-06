// +build ingore

package s390x

import (
	"cmd/internal/obj"
	"cmd/internal/sys"
	"cmd/link/internal/ld"
	"debug/elf"
	"fmt"
)

//  Used by ../internal/ld/dwarf.go
const (
	DWARFREGSP = 15
	DWARFREGLR = 14
)

const (
	MaxAlign  = 32 // max data alignment
	MinAlign  = 2  // min data alignment
	FuncAlign = 16
)

func Main()

