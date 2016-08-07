// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package objfile implements portable access to OS-specific executable files.

// Package objfile implements portable access to OS-specific executable files.
package objfile

import (
    "bufio"
    "cmd/internal/goobj"
    "cmd/internal/unvendor/golang.org/x/arch/arm/armasm"
    "cmd/internal/unvendor/golang.org/x/arch/x86/x86asm"
    "debug/elf"
    "debug/gosym"
    "debug/macho"
    "debug/pe"
    "debug/plan9obj"
    "encoding/binary"
    "fmt"
    "io"
    "os"
    "regexp"
    "sort"
    "strings"
    "text/tabwriter"
)

// Disasm is a disassembler for a given File.
type Disasm struct {
}

// A File is an opened executable file.
type File struct {
}

// A Sym is a symbol defined in an executable file.
type Sym struct {
    Name string // symbol name
    Addr uint64 // virtual address of symbol
    Size int64  // size in bytes
    Code rune   // nm code (T for text, D for data, and so on)
    Type string // XXX?
}

// Open opens the named file. The caller must call f.Close when the file is no
// longer needed.
func Open(name string) (*File, error)

// Decode disassembles the text segment range [start, end), calling f for each
// instruction.
func (*Disasm) Decode(start, end uint64, f func(pc, size uint64, file string, line int, text string))

// Print prints a disassembly of the file to w. If filter is non-nil, the
// disassembly only includes functions with names matching filter. The
// disassembly only includes functions that overlap the range [start, end).
func (*Disasm) Print(w io.Writer, filter *regexp.Regexp, start, end uint64)

func (*File) Close() error

// Disasm returns a disassembler for the file f.
func (*File) Disasm() (*Disasm, error)

func (*File) GOARCH() string

func (*File) PCLineTable() (*gosym.Table, error)

func (*File) Symbols() ([]Sym, error)

func (*File) Text() (uint64, []byte, error)

