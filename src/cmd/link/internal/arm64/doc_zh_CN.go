// +build ingore

package arm64

import (
	"cmd/internal/obj"
	"cmd/internal/sys"
	"cmd/link/internal/ld"
	"encoding/binary"
	"fmt"
	"log"
)

//  Used by ../internal/ld/dwarf.go

// Used by ../internal/ld/dwarf.go
const (
	DWARFREGSP = 31
	DWARFREGLR = 30
)

const (
	MaxAlign  = 32 // max data alignment
	MinAlign  = 1  // min data alignment
	FuncAlign = 8
)

func Main()

