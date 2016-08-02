// +build ingore

package x86 // import "cmd/link/internal/x86"

import (
    "cmd/internal/obj"
    "cmd/link/internal/ld"
    "fmt"
    "log"
)

//  Used by ../internal/ld/dwarf.go
const (
    DWARFREGSP = 4
    DWARFREGLR = 8
)

const (
    PtrSize   = 4
    MaxAlign  = 32 // max data alignment
    FuncAlign = 16
    MINLC     = 1
)

func Main()

