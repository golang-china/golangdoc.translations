// +build ingore

package arm

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
	DWARFREGSP = 13
	DWARFREGLR = 14
)



const (
	MaxAlign  = 8 // max data alignment
	MinAlign  = 1 // min data alignment
	FuncAlign = 4 // single-instruction alignment

)


func Main()

