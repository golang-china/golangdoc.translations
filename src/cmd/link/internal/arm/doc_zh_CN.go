// +build ingore

package arm // import "cmd/link/internal/arm"

import (
    "cmd/internal/obj"
    "cmd/link/internal/ld"
    "fmt"
    "log"
)

// Used by ../internal/ld/dwarf.go
const (
    DWARFREGSP = 13
    DWARFREGLR = 14
)

const (
    MaxAlign  = 8 // max data alignment
    FuncAlign = 4 // single-instruction alignment
    MINLC     = 4
)

func Main()

