// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package flags implements top-level flags and the usage message for the
// assembler.

// Package flags implements top-level flags and the usage message for the
// assembler.
package flags // import "cmd/asm/internal/flags"

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

var (
    D   MultiFlag
    I   MultiFlag
)

var (
    Debug      = flag.Bool("debug", false, "dump instructions as they are parsed")
    OutputFile = flag.String("o", "", "output file; default foo.6 for /a/b/c/foo.s on amd64")
    PrintOut   = flag.Bool("S", false, "print assembly and machine code")
    TrimPath   = flag.String("trimpath", "", "remove prefix from recorded source file paths")
    Shared     = flag.Bool("shared", false, "generate code that can be linked into a shared library")
    Dynlink    = flag.Bool("dynlink", false, "support references to Go symbols defined in other shared libraries")
    AllErrors  = flag.Bool("e", false, "no limit on number of errors reported")
)

// MultiFlag allows setting a value multiple times to collect a list, as in
// -I=dir1 -I=dir2.
type MultiFlag []string

func Parse()

func Usage()

func (*MultiFlag) Set(val string) error

func (*MultiFlag) String() string

