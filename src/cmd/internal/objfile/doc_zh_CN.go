// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package objfile implements portable access to OS-specific executable files.
package objfile

// Disasm is a disassembler for a given File.
type Disasm struct {
	// contains filtered or unexported fields
}

// Decode disassembles the text segment range [start, end), calling f for each
// instruction.
func (d *Disasm) Decode(start, end uint64, f func(pc, size uint64, file string, line int, text string))

// Print prints a disassembly of the file to w. If filter is non-nil, the
// disassembly only includes functions with names matching filter. The disassembly
// only includes functions that overlap the range [start, end).
func (d *Disasm) Print(w io.Writer, filter *regexp.Regexp, start, end uint64)

// A File is an opened executable file.
type File struct {
	// contains filtered or unexported fields
}

// Open opens the named file. The caller must call f.Close when the file is no
// longer needed.
func Open(name string) (*File, error)

func (f *File) Close() error

// Disasm returns a disassembler for the file f.
func (f *File) Disasm() (*Disasm, error)

func (f *File) GOARCH() string

func (f *File) PCLineTable() (*gosym.Table, error)

func (f *File) Symbols() ([]Sym, error)

func (f *File) Text() (uint64, []byte, error)

// A Sym is a symbol defined in an executable file.
type Sym struct {
	Name string // symbol name
	Addr uint64 // virtual address of symbol
	Size int64  // size in bytes
	Code rune   // nm code (T for text, D for data, and so on)
	Type string // XXX?
}
