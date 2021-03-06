// +build ingore

package amd64

import (
	"cmd/internal/obj"
	"cmd/internal/sys"
	"cmd/link/internal/ld"
	"debug/elf"
	"fmt"
	"log"
)

//  Used by ../internal/ld/dwarf.go

// Used by ../internal/ld/dwarf.go
const (
	DWARFREGSP = 7
	DWARFREGLR = 16
)

const (
	MaxAlign  = 32 // max data alignment
	MinAlign  = 1  // min data alignment
	FuncAlign = 16
)

func Addcall(ctxt *ld.Link, s *ld.LSym, t *ld.LSym) int64

func Main()

func PADDR(x uint32) uint32

