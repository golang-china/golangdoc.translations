// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// TODO(rsc): Handle go.typelink, go.track symbols. TODO(rsc): Do not handle $f64.
// and $f32. symbols. Instead, generate those from the compiler and assemblers as
// dupok data, and then remove autoData below.

// TODO(rsc): Handle go.typelink, go.track symbols. TODO(rsc): Do not handle $f64.
// and $f32. symbols. Instead, generate those from the compiler and assemblers as
// dupok data, and then remove autoData below.
package main

// TODO(rsc): Define full enumeration for relocation types.

// TODO(rsc): Define full enumeration for relocation types.
const (
	R_ADDR    = 1
	R_SIZE    = 2
	R_CALL    = 3
	R_CALLARM = 4
	R_CALLIND = 5
	R_CONST   = 6
	R_PCREL   = 7
)

// An Addr represents a virtual memory address, a file address, or a size. It must
// be a uint64, not a uintptr, so that a 32-bit linker can still generate a 64-bit
// binary. It must be unsigned in order to link programs placed at very large start
// addresses. Math involving Addrs must be checked carefully not to require
// negative numbers.

// An Addr represents a virtual memory address, a file address, or a size. It must
// be a uint64, not a uintptr, so that a 32-bit linker can still generate a 64-bit
// binary. It must be unsigned in order to link programs placed at very large start
// addresses. Math involving Addrs must be checked carefully not to require
// negative numbers.
type Addr uint64

// A PCIter implements iteration over PC-data tables.
//
//	var it PCIter
//	for it.Init(p, data); !it.Done; it.Next() {
//		it.Value holds from it.PC up to (but not including) it.NextPC
//	}
//	if it.Corrupt {
//		data was malformed
//	}

// A PCIter implements iteration over PC-data tables.
//
//	var it PCIter
//	for it.Init(p, data); !it.Done; it.Next() {
//		it.Value holds from it.PC up to (but not including) it.NextPC
//	}
//	if it.Corrupt {
//		data was malformed
//	}
type PCIter struct {
	PC      uint32
	NextPC  uint32
	Value   int32
	Done    bool
	Corrupt bool
	// contains filtered or unexported fields
}

// Init initializes the iteration. On return, if it.Done is true, the iteration is
// over. Otherwise it.Value applies in the pc range [it.PC, it.NextPC).

// Init initializes the iteration. On return, if it.Done is true, the iteration is
// over. Otherwise it.Value applies in the pc range [it.PC, it.NextPC).
func (it *PCIter) Init(p *Prog, buf []byte)

// Next steps forward one entry in the table. On return, if it.Done is true, the
// iteration is over. Otherwise it.Value applies in the pc range [it.PC,
// it.NextPC).

// Next steps forward one entry in the table. On return, if it.Done is true, the
// iteration is over. Otherwise it.Value applies in the pc range [it.PC,
// it.NextPC).
func (it *PCIter) Next()

// A Package is a Go package loaded from a file.

// A Package is a Go package loaded from a file.
type Package struct {
	*goobj.Package        // table of contents
	File           string // file name for reopening
	Syms           []*Sym // symbols defined by this package
}

// A Prog holds state for constructing an executable (program) image.
//
// The usual sequence of operations on a Prog is:
//
//	p.init()
//	p.scan(file)
//	p.dead()
//	p.runtime()
//	p.layout()
//	p.load()
//	p.debug()
//	p.write(w)
//
// p.init is in this file. The rest of the methods are in files named for the
// method. The convenience method p.link runs this sequence.

// A Prog holds state for constructing an executable (program) image.
//
// The usual sequence of operations on a Prog is:
//
//	p.init()
//	p.scan(file)
//	p.dead()
//	p.runtime()
//	p.layout()
//	p.load()
//	p.debug()
//	p.write(w)
//
// p.init is in this file. The rest of the methods are in files named for the
// method. The convenience method p.link runs this sequence.
type Prog struct {
	// Context
	GOOS     string       // target operating system
	GOARCH   string       // target architecture
	Format   string       // desired file format ("elf", "macho", ...)
	Error    func(string) // called to report an error (if set)
	NumError int          // number of errors printed
	StartSym string

	// Input
	Packages   map[string]*Package  // loaded packages, by import path
	Syms       map[goobj.SymID]*Sym // defined symbols, by symbol ID
	Missing    map[goobj.SymID]bool // missing symbols
	Dead       map[goobj.SymID]bool // symbols removed as dead
	SymOrder   []*Sym               // order syms were scanned
	MaxVersion int                  // max SymID.Version, for generating fresh symbol IDs

	// Output
	UnmappedSize Addr       // size of unmapped region at address 0
	HeaderSize   Addr       // size of object file header
	Entry        Addr       // virtual address where execution begins
	Segments     []*Segment // loaded memory segments
	// contains filtered or unexported fields
}

// A Section is part of a loaded memory segment.

// A Section is part of a loaded memory segment.
type Section struct {
	Name     string   // name of section: "text", "rodata", "noptrbss", and so on
	VirtAddr Addr     // virtual memory address of section base
	Size     Addr     // size of section in memory
	Align    Addr     // required alignment
	InFile   bool     // section has image data in file (like data, unlike bss)
	Syms     []*Sym   // symbols stored in section
	Segment  *Segment // segment containing section
}

// A Segment is a loaded memory segment. A Prog is expected to have segments named
// "text" and optionally "data", in that order, before any other segments.

// A Segment is a loaded memory segment. A Prog is expected to have segments named
// "text" and optionally "data", in that order, before any other segments.
type Segment struct {
	Name       string     // name of segment: "text", "data", ...
	VirtAddr   Addr       // virtual memory address of segment base
	VirtSize   Addr       // size of segment in memory
	FileOffset Addr       // file offset of segment base
	FileSize   Addr       // size of segment in file; can be less than VirtSize
	Sections   []*Section // sections inside segment
	Data       []byte     // raw data of segment image
}

// A Sym is a symbol defined in a loaded package.

// A Sym is a symbol defined in a loaded package.
type Sym struct {
	*goobj.Sym          // symbol metadata from package file
	Package    *Package // package defining symbol
	Section    *Section // section where symbol is placed in output program
	Addr       Addr     // virtual address of symbol in output program
	Bytes      []byte   // symbol data, for internally defined symbols
}

// A SymBuffer is a buffer for preparing the data image of a linker-generated
// symbol.

// A SymBuffer is a buffer for preparing the data image of a linker-generated
// symbol.
type SymBuffer struct {
	// contains filtered or unexported fields
}

// Addr sets the pointer-sized address at offset off to refer to symoff bytes past
// the start of sym. It returns the offset just beyond the address.

// Addr sets the pointer-sized address at offset off to refer to symoff bytes past
// the start of sym. It returns the offset just beyond the address.
func (b *SymBuffer) Addr(off int, sym goobj.SymID, symoff int64) int

// Bytes returns the buffer data.

// Bytes returns the buffer data.
func (b *SymBuffer) Bytes() []byte

// Init initializes the buffer for writing.

// Init initializes the buffer for writing.
func (b *SymBuffer) Init(p *Prog)

// Reloc returns the buffered relocations.

// Reloc returns the buffered relocations.
func (b *SymBuffer) Reloc() []goobj.Reloc

// SetSize sets the buffer's data size to n bytes.

// SetSize sets the buffer's data size to n bytes.
func (b *SymBuffer) SetSize(n int)

// Size returns the buffer's data size.

// Size returns the buffer's data size.
func (b *SymBuffer) Size() int

// Uint sets the size-byte unsigned integer at offset off to v. It returns the
// offset just beyond v.

// Uint sets the size-byte unsigned integer at offset off to v. It returns the
// offset just beyond v.
func (b *SymBuffer) Uint(off int, v uint64, size int) int

// Uint16 sets the uint16 at offset off to v. It returns the offset just beyond v.

// Uint16 sets the uint16 at offset off to v. It returns the offset just beyond v.
func (b *SymBuffer) Uint16(off int, v uint16) int

// Uint32 sets the uint32 at offset off to v. It returns the offset just beyond v.

// Uint32 sets the uint32 at offset off to v. It returns the offset just beyond v.
func (b *SymBuffer) Uint32(off int, v uint32) int

// Uint64 sets the uint64 at offset off to v. It returns the offset just beyond v.

// Uint64 sets the uint64 at offset off to v. It returns the offset just beyond v.
func (b *SymBuffer) Uint64(off int, v uint64) int

// Uint8 sets the uint8 at offset off to v. It returns the offset just beyond v.

// Uint8 sets the uint8 at offset off to v. It returns the offset just beyond v.
func (b *SymBuffer) Uint8(off int, v uint8) int
