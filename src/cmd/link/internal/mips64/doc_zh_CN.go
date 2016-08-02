// +build ingore

package mips64 // import "cmd/link/internal/mips64"

import (
    "cmd/internal/obj"
    "cmd/link/internal/ld"
    "encoding/binary"
    "fmt"
    "log"
)

//  Used by ../internal/ld/dwarf.go
const (
    DWARFREGSP = 29
    DWARFREGLR = 31
)

const (
    MaxAlign  = 32 // max data alignment
    FuncAlign = 8
    MINLC     = 4
)

func Main()

