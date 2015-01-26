// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package plan9obj implements access to Plan 9 a.out object files.

// Package plan9obj implements access to
// Plan 9 a.out object files.
package plan9obj

const (
	Magic64 = 0x8000 // 64-bit expanded header

	Magic386   = (4*11+0)*11 + 7
	MagicAMD64 = (4*26+0)*26 + 7 + Magic64
	MagicARM   = (4*20+0)*20 + 7
)

// A File represents an open Plan 9 a.out file.

// A File represents an open Plan 9 a.out
// file.
type File struct {
	FileHeader
	Sections []*Section
	// contains filtered or unexported fields
}

// NewFile creates a new File for accessing a Plan 9 binary in an underlying
// reader. The Plan 9 binary is expected to start at position 0 in the ReaderAt.

// NewFile creates a new File for accessing
// a Plan 9 binary in an underlying reader.
// The Plan 9 binary is expected to start
// at position 0 in the ReaderAt.
func NewFile(r io.ReaderAt) (*File, error)

// Open opens the named file using os.Open and prepares it for use as a Plan 9
// a.out binary.

// Open opens the named file using os.Open
// and prepares it for use as a Plan 9
// a.out binary.
func Open(name string) (*File, error)

// Close closes the File. If the File was created using NewFile directly instead of
// Open, Close has no effect.

// Close closes the File. If the File was
// created using NewFile directly instead
// of Open, Close has no effect.
func (f *File) Close() error

// Section returns a section with the given name, or nil if no such section exists.

// Section returns a section with the given
// name, or nil if no such section exists.
func (f *File) Section(name string) *Section

// Symbols returns the symbol table for f.

// Symbols returns the symbol table for f.
func (f *File) Symbols() ([]Sym, error)

// A FileHeader represents a Plan 9 a.out file header.

// A FileHeader represents a Plan 9 a.out
// file header.
type FileHeader struct {
	Magic       uint32
	Bss         uint32
	Entry       uint64
	PtrSize     int
	LoadAddress uint64
	HdrSize     uint64
}

// A Section represents a single section in a Plan 9 a.out file.

// A Section represents a single section in
// a Plan 9 a.out file.
type Section struct {
	SectionHeader

	// Embed ReaderAt for ReadAt method.
	// Do not embed SectionReader directly
	// to avoid having Read and Seek.
	// If a client wants Read and Seek it must use
	// Open() to avoid fighting over the seek offset
	// with other clients.
	io.ReaderAt
	// contains filtered or unexported fields
}

// Data reads and returns the contents of the Plan 9 a.out section.

// Data reads and returns the contents of
// the Plan 9 a.out section.
func (s *Section) Data() ([]byte, error)

// Open returns a new ReadSeeker reading the Plan 9 a.out section.

// Open returns a new ReadSeeker reading
// the Plan 9 a.out section.
func (s *Section) Open() io.ReadSeeker

// A SectionHeader represents a single Plan 9 a.out section header. This structure
// doesn't exist on-disk, but eases navigation through the object file.

// A SectionHeader represents a single Plan
// 9 a.out section header. This structure
// doesn't exist on-disk, but eases
// navigation through the object file.
type SectionHeader struct {
	Name   string
	Size   uint32
	Offset uint32
}

// A Symbol represents an entry in a Plan 9 a.out symbol table section.

// A Symbol represents an entry in a Plan 9
// a.out symbol table section.
type Sym struct {
	Value uint64
	Type  rune
	Name  string
}
