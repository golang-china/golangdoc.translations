// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package goobj implements reading of Go object files and archives.
//
// TODO(rsc): Decide where this package should live. (golang.org/issue/6932)
// TODO(rsc): Decide the appropriate integer types for various fields.
// TODO(rsc): Write tests. (File format still up in the air a little.)

// Package goobj implements reading of Go object files and archives.
//
// TODO(rsc): Decide where this package should live. (golang.org/issue/6932)
// TODO(rsc): Decide the appropriate integer types for various fields.
// TODO(rsc): Write tests. (File format still up in the air a little.)
package goobj

import (
    "bufio"
    "bytes"
    "cmd/internal/obj"
    "errors"
    "fmt"
    "io"
    "strconv"
    "strings"
)

// Defined SymKind values.
// TODO(rsc): Give idiomatic Go names.
// TODO(rsc): Reduce the number of symbol types in the object files.
const (
	_ SymKind = iota
	// readonly, executable
	STEXT      SymKind = obj.STEXT
	SELFRXSECT SymKind = obj.SELFRXSECT
	// readonly, non-executable
	STYPE      SymKind = obj.STYPE
	SSTRING    SymKind = obj.SSTRING
	SGOSTRING  SymKind = obj.SGOSTRING
	SGOFUNC    SymKind = obj.SGOFUNC
	SRODATA    SymKind = obj.SRODATA
	SFUNCTAB   SymKind = obj.SFUNCTAB
	STYPELINK  SymKind = obj.STYPELINK
	SITABLINK  SymKind = obj.SITABLINK
	SSYMTAB    SymKind = obj.SSYMTAB // TODO: move to unmapped section
	SPCLNTAB   SymKind = obj.SPCLNTAB
	SELFROSECT SymKind = obj.SELFROSECT
	// writable, non-executable
	SMACHOPLT  SymKind = obj.SMACHOPLT
	SELFSECT   SymKind = obj.SELFSECT
	SMACHO     SymKind = obj.SMACHO // Mach-O __nl_symbol_ptr
	SMACHOGOT  SymKind = obj.SMACHOGOT
	SWINDOWS   SymKind = obj.SWINDOWS
	SELFGOT    SymKind = obj.SELFGOT
	SNOPTRDATA SymKind = obj.SNOPTRDATA
	SINITARR   SymKind = obj.SINITARR
	SDATA      SymKind = obj.SDATA
	SBSS       SymKind = obj.SBSS
	SNOPTRBSS  SymKind = obj.SNOPTRBSS
	STLSBSS    SymKind = obj.STLSBSS
	// not mapped
	SXREF             SymKind = obj.SXREF
	SMACHOSYMSTR      SymKind = obj.SMACHOSYMSTR
	SMACHOSYMTAB      SymKind = obj.SMACHOSYMTAB
	SMACHOINDIRECTPLT SymKind = obj.SMACHOINDIRECTPLT
	SMACHOINDIRECTGOT SymKind = obj.SMACHOINDIRECTGOT
	SFILE             SymKind = obj.SFILE
	SFILEPATH         SymKind = obj.SFILEPATH
	SCONST            SymKind = obj.SCONST
	SDYNIMPORT        SymKind = obj.SDYNIMPORT
	SHOSTOBJ          SymKind = obj.SHOSTOBJ
)


// A Data is a reference to data stored in an object file.
// It records the offset and size of the data, so that a client can
// read the data only if necessary.
type Data struct {
	Offset int64
	Size   int64
}


// Func contains additional per-symbol information specific to functions.
type Func struct {
	Args     int        // size in bytes of argument frame: inputs and outputs
	Frame    int        // size in bytes of local variable frame
	Leaf     bool       // function omits save of link register (ARM)
	NoSplit  bool       // function omits stack split prologue
	Var      []Var      // detail about local variables
	PCSP     Data       // PC → SP offset map
	PCFile   Data       // PC → file number map (index into File)
	PCLine   Data       // PC → line number map
	PCData   []Data     // PC → runtime support data map
	FuncData []FuncData // non-PC-specific runtime support data
	File     []string   // paths indexed by PCFile
}


// A FuncData is a single function-specific data value.
type FuncData struct {
	Sym    SymID // symbol holding data
	Offset int64 // offset into symbol for funcdata pointer
}


// A Package is a parsed Go object file or archive defining a Go package.
type Package struct {
	ImportPath string   // import path denoting this package
	Imports    []string // packages imported by this package
	SymRefs    []SymID  // list of symbol names and versions referred to by this pack
	Syms       []*Sym   // symbols defined by this package
	MaxVersion int      // maximum Version in any SymID in Syms
}


// A Reloc describes a relocation applied to a memory image to refer
// to an address within a particular symbol.
type Reloc struct {
	// The bytes at [Offset, Offset+Size) within the memory image
	// should be updated to refer to the address Add bytes after the start
	// of the symbol Sym.
	Offset int
	Size   int
	Sym    SymID
	Add    int

	// The Type records the form of address expected in the bytes
	// described by the previous fields: absolute, PC-relative, and so on.
	// TODO(rsc): The interpretation of Type is not exposed by this package.
	Type int
}


// A Sym is a named symbol in an object file.
type Sym struct {
	Kind  SymKind // kind of symbol
	DupOK bool    // are duplicate definitions okay?
	Size  int     // size of corresponding data
	Type  SymID   // symbol for Go type information
	Data  Data    // memory image of symbol
	Reloc []Reloc // relocations to apply to Data
	Func  *Func   // additional data for functions
}


// A SymID - the combination of Name and Version - uniquely identifies
// a symbol within a package.
type SymID struct {
	// Name is the name of a symbol.
	Name string

	// Version is zero for symbols with global visibility.
	// Symbols with only file visibility (such as file-level static
	// declarations in C) have a non-zero version distinguishing
	// a symbol in one file from a symbol of the same name
	// in another file
	Version int
}


// A SymKind describes the kind of memory represented by a symbol.
type SymKind int


// A Var describes a variable in a function stack frame: a declared
// local variable, an input argument, or an output result.
type Var struct {
	// The combination of Name, Kind, and Offset uniquely
	// identifies a variable in a function stack frame.
	// Using fewer of these - in particular, using only Name - does not.
	Name   string // Name of variable.
	Kind   int    // TODO(rsc): Define meaning.
	Offset int    // Frame offset. TODO(rsc): Define meaning.

	Type SymID // Go type for variable.
}


// Parse parses an object file or archive from r,
// assuming that its import path is pkgpath.
func Parse(r io.ReadSeeker, pkgpath string) (*Package, error)

func (SymID) String() string

func (SymKind) String() string

