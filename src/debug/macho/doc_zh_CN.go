// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package macho implements access to Mach-O object files.

// Package macho implements access to
// Mach-O object files.
package macho

const (
	Magic32  uint32 = 0xfeedface
	Magic64  uint32 = 0xfeedfacf
	MagicFat uint32 = 0xcafebabe
)

// ErrNotFat is returned from NewFatFile or OpenFat when the file is not a
// universal binary but may be a thin binary, based on its magic number.

// ErrNotFat is returned from NewFatFile or
// OpenFat when the file is not a universal
// binary but may be a thin binary, based
// on its magic number.
var ErrNotFat = &FormatError{0, "not a fat Mach-O file", nil}

// A Cpu is a Mach-O cpu type.

// A Cpu is a Mach-O cpu type.
type Cpu uint32

const (
	Cpu386   Cpu = 7
	CpuAmd64 Cpu = Cpu386 | cpuArch64
	CpuArm   Cpu = 12
	CpuPpc   Cpu = 18
	CpuPpc64 Cpu = CpuPpc | cpuArch64
)

func (i Cpu) GoString() string

func (i Cpu) String() string

// A Dylib represents a Mach-O load dynamic library command.

// A Dylib represents a Mach-O load dynamic
// library command.
type Dylib struct {
	LoadBytes
	Name           string
	Time           uint32
	CurrentVersion uint32
	CompatVersion  uint32
}

// A DylibCmd is a Mach-O load dynamic library command.

// A DylibCmd is a Mach-O load dynamic
// library command.
type DylibCmd struct {
	Cmd            LoadCmd
	Len            uint32
	Name           uint32
	Time           uint32
	CurrentVersion uint32
	CompatVersion  uint32
}

// A Dysymtab represents a Mach-O dynamic symbol table command.

// A Dysymtab represents a Mach-O dynamic
// symbol table command.
type Dysymtab struct {
	LoadBytes
	DysymtabCmd
	IndirectSyms []uint32 // indices into Symtab.Syms
}

// A DysymtabCmd is a Mach-O dynamic symbol table command.

// A DysymtabCmd is a Mach-O dynamic symbol
// table command.
type DysymtabCmd struct {
	Cmd            LoadCmd
	Len            uint32
	Ilocalsym      uint32
	Nlocalsym      uint32
	Iextdefsym     uint32
	Nextdefsym     uint32
	Iundefsym      uint32
	Nundefsym      uint32
	Tocoffset      uint32
	Ntoc           uint32
	Modtaboff      uint32
	Nmodtab        uint32
	Extrefsymoff   uint32
	Nextrefsyms    uint32
	Indirectsymoff uint32
	Nindirectsyms  uint32
	Extreloff      uint32
	Nextrel        uint32
	Locreloff      uint32
	Nlocrel        uint32
}

// A FatArch is a Mach-O File inside a FatFile.

// A FatArch is a Mach-O File inside a
// FatFile.
type FatArch struct {
	FatArchHeader
	*File
}

// A FatArchHeader represents a fat header for a specific image architecture.

// A FatArchHeader represents a fat header
// for a specific image architecture.
type FatArchHeader struct {
	Cpu    Cpu
	SubCpu uint32
	Offset uint32
	Size   uint32
	Align  uint32
}

// A FatFile is a Mach-O universal binary that contains at least one architecture.

// A FatFile is a Mach-O universal binary
// that contains at least one architecture.
type FatFile struct {
	Magic  uint32
	Arches []FatArch
	// contains filtered or unexported fields
}

// NewFatFile creates a new FatFile for accessing all the Mach-O images in a
// universal binary. The Mach-O binary is expected to start at position 0 in the
// ReaderAt.

// NewFatFile creates a new FatFile for
// accessing all the Mach-O images in a
// universal binary. The Mach-O binary is
// expected to start at position 0 in the
// ReaderAt.
func NewFatFile(r io.ReaderAt) (*FatFile, error)

// OpenFat opens the named file using os.Open and prepares it for use as a Mach-O
// universal binary.

// OpenFat opens the named file using
// os.Open and prepares it for use as a
// Mach-O universal binary.
func OpenFat(name string) (ff *FatFile, err error)

func (ff *FatFile) Close() error

// A File represents an open Mach-O file.

// A File represents an open Mach-O file.
type File struct {
	FileHeader
	ByteOrder binary.ByteOrder
	Loads     []Load
	Sections  []*Section

	Symtab   *Symtab
	Dysymtab *Dysymtab
	// contains filtered or unexported fields
}

// NewFile creates a new File for accessing a Mach-O binary in an underlying
// reader. The Mach-O binary is expected to start at position 0 in the ReaderAt.

// NewFile creates a new File for accessing
// a Mach-O binary in an underlying reader.
// The Mach-O binary is expected to start
// at position 0 in the ReaderAt.
func NewFile(r io.ReaderAt) (*File, error)

// Open opens the named file using os.Open and prepares it for use as a Mach-O
// binary.

// Open opens the named file using os.Open
// and prepares it for use as a Mach-O
// binary.
func Open(name string) (*File, error)

// Close closes the File. If the File was created using NewFile directly instead of
// Open, Close has no effect.

// Close closes the File. If the File was
// created using NewFile directly instead
// of Open, Close has no effect.
func (f *File) Close() error

// DWARF returns the DWARF debug information for the Mach-O file.

// DWARF returns the DWARF debug
// information for the Mach-O file.
func (f *File) DWARF() (*dwarf.Data, error)

// ImportedLibraries returns the paths of all libraries referred to by the binary f
// that are expected to be linked with the binary at dynamic link time.

// ImportedLibraries returns the paths of
// all libraries referred to by the binary
// f that are expected to be linked with
// the binary at dynamic link time.
func (f *File) ImportedLibraries() ([]string, error)

// ImportedSymbols returns the names of all symbols referred to by the binary f
// that are expected to be satisfied by other libraries at dynamic load time.

// ImportedSymbols returns the names of all
// symbols referred to by the binary f that
// are expected to be satisfied by other
// libraries at dynamic load time.
func (f *File) ImportedSymbols() ([]string, error)

// Section returns the first section with the given name, or nil if no such section
// exists.

// Section returns the first section with
// the given name, or nil if no such
// section exists.
func (f *File) Section(name string) *Section

// Segment returns the first Segment with the given name, or nil if no such segment
// exists.

// Segment returns the first Segment with
// the given name, or nil if no such
// segment exists.
func (f *File) Segment(name string) *Segment

// A FileHeader represents a Mach-O file header.

// A FileHeader represents a Mach-O file
// header.
type FileHeader struct {
	Magic  uint32
	Cpu    Cpu
	SubCpu uint32
	Type   Type
	Ncmd   uint32
	Cmdsz  uint32
	Flags  uint32
}

// FormatError is returned by some operations if the data does not have the correct
// format for an object file.

// FormatError is returned by some
// operations if the data does not have the
// correct format for an object file.
type FormatError struct {
	// contains filtered or unexported fields
}

func (e *FormatError) Error() string

// A Load represents any Mach-O load command.

// A Load represents any Mach-O load
// command.
type Load interface {
	Raw() []byte
}

// A LoadBytes is the uninterpreted bytes of a Mach-O load command.

// A LoadBytes is the uninterpreted bytes
// of a Mach-O load command.
type LoadBytes []byte

func (b LoadBytes) Raw() []byte

// A LoadCmd is a Mach-O load command.

// A LoadCmd is a Mach-O load command.
type LoadCmd uint32

const (
	LoadCmdSegment    LoadCmd = 1
	LoadCmdSymtab     LoadCmd = 2
	LoadCmdThread     LoadCmd = 4
	LoadCmdUnixThread LoadCmd = 5 // thread+stack
	LoadCmdDysymtab   LoadCmd = 11
	LoadCmdDylib      LoadCmd = 12
	LoadCmdDylinker   LoadCmd = 15
	LoadCmdSegment64  LoadCmd = 25
)

func (i LoadCmd) GoString() string

func (i LoadCmd) String() string

// An Nlist32 is a Mach-O 32-bit symbol table entry.

// An Nlist32 is a Mach-O 32-bit symbol
// table entry.
type Nlist32 struct {
	Name  uint32
	Type  uint8
	Sect  uint8
	Desc  uint16
	Value uint32
}

// An Nlist64 is a Mach-O 64-bit symbol table entry.

// An Nlist64 is a Mach-O 64-bit symbol
// table entry.
type Nlist64 struct {
	Name  uint32
	Type  uint8
	Sect  uint8
	Desc  uint16
	Value uint64
}

// Regs386 is the Mach-O 386 register structure.

// Regs386 is the Mach-O 386 register
// structure.
type Regs386 struct {
	AX    uint32
	BX    uint32
	CX    uint32
	DX    uint32
	DI    uint32
	SI    uint32
	BP    uint32
	SP    uint32
	SS    uint32
	FLAGS uint32
	IP    uint32
	CS    uint32
	DS    uint32
	ES    uint32
	FS    uint32
	GS    uint32
}

// RegsAMD64 is the Mach-O AMD64 register structure.

// RegsAMD64 is the Mach-O AMD64 register
// structure.
type RegsAMD64 struct {
	AX    uint64
	BX    uint64
	CX    uint64
	DX    uint64
	DI    uint64
	SI    uint64
	BP    uint64
	SP    uint64
	R8    uint64
	R9    uint64
	R10   uint64
	R11   uint64
	R12   uint64
	R13   uint64
	R14   uint64
	R15   uint64
	IP    uint64
	FLAGS uint64
	CS    uint64
	FS    uint64
	GS    uint64
}

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

// Data reads and returns the contents of the Mach-O section.

// Data reads and returns the contents of
// the Mach-O section.
func (s *Section) Data() ([]byte, error)

// Open returns a new ReadSeeker reading the Mach-O section.

// Open returns a new ReadSeeker reading
// the Mach-O section.
func (s *Section) Open() io.ReadSeeker

// A Section32 is a 32-bit Mach-O section header.

// A Section32 is a 32-bit Mach-O section
// header.
type Section32 struct {
	Name     [16]byte
	Seg      [16]byte
	Addr     uint32
	Size     uint32
	Offset   uint32
	Align    uint32
	Reloff   uint32
	Nreloc   uint32
	Flags    uint32
	Reserve1 uint32
	Reserve2 uint32
}

// A Section32 is a 64-bit Mach-O section header.

// A Section32 is a 64-bit Mach-O section
// header.
type Section64 struct {
	Name     [16]byte
	Seg      [16]byte
	Addr     uint64
	Size     uint64
	Offset   uint32
	Align    uint32
	Reloff   uint32
	Nreloc   uint32
	Flags    uint32
	Reserve1 uint32
	Reserve2 uint32
	Reserve3 uint32
}

type SectionHeader struct {
	Name   string
	Seg    string
	Addr   uint64
	Size   uint64
	Offset uint32
	Align  uint32
	Reloff uint32
	Nreloc uint32
	Flags  uint32
}

// A Segment represents a Mach-O 32-bit or 64-bit load segment command.

// A Segment represents a Mach-O 32-bit or
// 64-bit load segment command.
type Segment struct {
	LoadBytes
	SegmentHeader

	// Embed ReaderAt for ReadAt method.
	// Do not embed SectionReader directly
	// to avoid having Read and Seek.
	// If a client wants Read and Seek it must use
	// Open() to avoid fighting over the seek offset
	// with other clients.
	io.ReaderAt
	// contains filtered or unexported fields
}

// Data reads and returns the contents of the segment.

// Data reads and returns the contents of
// the segment.
func (s *Segment) Data() ([]byte, error)

// Open returns a new ReadSeeker reading the segment.

// Open returns a new ReadSeeker reading
// the segment.
func (s *Segment) Open() io.ReadSeeker

// A Segment32 is a 32-bit Mach-O segment load command.

// A Segment32 is a 32-bit Mach-O segment
// load command.
type Segment32 struct {
	Cmd     LoadCmd
	Len     uint32
	Name    [16]byte
	Addr    uint32
	Memsz   uint32
	Offset  uint32
	Filesz  uint32
	Maxprot uint32
	Prot    uint32
	Nsect   uint32
	Flag    uint32
}

// A Segment64 is a 64-bit Mach-O segment load command.

// A Segment64 is a 64-bit Mach-O segment
// load command.
type Segment64 struct {
	Cmd     LoadCmd
	Len     uint32
	Name    [16]byte
	Addr    uint64
	Memsz   uint64
	Offset  uint64
	Filesz  uint64
	Maxprot uint32
	Prot    uint32
	Nsect   uint32
	Flag    uint32
}

// A SegmentHeader is the header for a Mach-O 32-bit or 64-bit load segment
// command.

// A SegmentHeader is the header for a
// Mach-O 32-bit or 64-bit load segment
// command.
type SegmentHeader struct {
	Cmd     LoadCmd
	Len     uint32
	Name    string
	Addr    uint64
	Memsz   uint64
	Offset  uint64
	Filesz  uint64
	Maxprot uint32
	Prot    uint32
	Nsect   uint32
	Flag    uint32
}

// A Symbol is a Mach-O 32-bit or 64-bit symbol table entry.

// A Symbol is a Mach-O 32-bit or 64-bit
// symbol table entry.
type Symbol struct {
	Name  string
	Type  uint8
	Sect  uint8
	Desc  uint16
	Value uint64
}

// A Symtab represents a Mach-O symbol table command.

// A Symtab represents a Mach-O symbol
// table command.
type Symtab struct {
	LoadBytes
	SymtabCmd
	Syms []Symbol
}

// A SymtabCmd is a Mach-O symbol table command.

// A SymtabCmd is a Mach-O symbol table
// command.
type SymtabCmd struct {
	Cmd     LoadCmd
	Len     uint32
	Symoff  uint32
	Nsyms   uint32
	Stroff  uint32
	Strsize uint32
}

// A Thread is a Mach-O thread state command.

// A Thread is a Mach-O thread state
// command.
type Thread struct {
	Cmd  LoadCmd
	Len  uint32
	Type uint32
	Data []uint32
}

// A Type is the Mach-O file type, e.g. an object file, executable, or dynamic
// library.

// A Type is the Mach-O file type, e.g. an
// object file, executable, or dynamic
// library.
type Type uint32

const (
	TypeObj    Type = 1
	TypeExec   Type = 2
	TypeDylib  Type = 6
	TypeBundle Type = 8
)
