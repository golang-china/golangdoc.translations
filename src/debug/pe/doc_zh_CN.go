// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package pe implements access to PE (Microsoft Windows Portable Executable)
// files.

// Package pe implements access to PE
// (Microsoft Windows Portable Executable)
// files.
package pe

const (
	IMAGE_FILE_MACHINE_UNKNOWN   = 0x0
	IMAGE_FILE_MACHINE_AM33      = 0x1d3
	IMAGE_FILE_MACHINE_AMD64     = 0x8664
	IMAGE_FILE_MACHINE_ARM       = 0x1c0
	IMAGE_FILE_MACHINE_EBC       = 0xebc
	IMAGE_FILE_MACHINE_I386      = 0x14c
	IMAGE_FILE_MACHINE_IA64      = 0x200
	IMAGE_FILE_MACHINE_M32R      = 0x9041
	IMAGE_FILE_MACHINE_MIPS16    = 0x266
	IMAGE_FILE_MACHINE_MIPSFPU   = 0x366
	IMAGE_FILE_MACHINE_MIPSFPU16 = 0x466
	IMAGE_FILE_MACHINE_POWERPC   = 0x1f0
	IMAGE_FILE_MACHINE_POWERPCFP = 0x1f1
	IMAGE_FILE_MACHINE_R4000     = 0x166
	IMAGE_FILE_MACHINE_SH3       = 0x1a2
	IMAGE_FILE_MACHINE_SH3DSP    = 0x1a3
	IMAGE_FILE_MACHINE_SH4       = 0x1a6
	IMAGE_FILE_MACHINE_SH5       = 0x1a8
	IMAGE_FILE_MACHINE_THUMB     = 0x1c2
	IMAGE_FILE_MACHINE_WCEMIPSV2 = 0x169
)

const COFFSymbolSize = 18

type COFFSymbol struct {
	Name               [8]uint8
	Value              uint32
	SectionNumber      int16
	Type               uint16
	StorageClass       uint8
	NumberOfAuxSymbols uint8
}

type DataDirectory struct {
	VirtualAddress uint32
	Size           uint32
}

// A File represents an open PE file.

// A File represents an open PE file.
type File struct {
	FileHeader
	OptionalHeader interface{} // of type *OptionalHeader32 or *OptionalHeader64
	Sections       []*Section
	Symbols        []*Symbol
	// contains filtered or unexported fields
}

// NewFile creates a new File for accessing a PE binary in an underlying reader.

// NewFile creates a new File for accessing
// a PE binary in an underlying reader.
func NewFile(r io.ReaderAt) (*File, error)

// Open opens the named file using os.Open and prepares it for use as a PE binary.

// Open opens the named file using os.Open
// and prepares it for use as a PE binary.
func Open(name string) (*File, error)

// Close closes the File. If the File was created using NewFile directly instead of
// Open, Close has no effect.

// Close closes the File. If the File was
// created using NewFile directly instead
// of Open, Close has no effect.
func (f *File) Close() error

func (f *File) DWARF() (*dwarf.Data, error)

// ImportedLibraries returns the names of all libraries referred to by the binary f
// that are expected to be linked with the binary at dynamic link time.

// ImportedLibraries returns the names of
// all libraries referred to by the binary
// f that are expected to be linked with
// the binary at dynamic link time.
func (f *File) ImportedLibraries() ([]string, error)

// ImportedSymbols returns the names of all symbols referred to by the binary f
// that are expected to be satisfied by other libraries at dynamic load time. It
// does not return weak symbols.

// ImportedSymbols returns the names of all
// symbols referred to by the binary f that
// are expected to be satisfied by other
// libraries at dynamic load time. It does
// not return weak symbols.
func (f *File) ImportedSymbols() ([]string, error)

// Section returns the first section with the given name, or nil if no such section
// exists.

// Section returns the first section with
// the given name, or nil if no such
// section exists.
func (f *File) Section(name string) *Section

type FileHeader struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type FormatError struct {
	// contains filtered or unexported fields
}

func (e *FormatError) Error() string

type ImportDirectory struct {
	OriginalFirstThunk uint32
	TimeDateStamp      uint32
	ForwarderChain     uint32
	Name               uint32
	FirstThunk         uint32
	// contains filtered or unexported fields
}

type OptionalHeader32 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	BaseOfData                  uint32
	ImageBase                   uint32
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint32
	SizeOfStackCommit           uint32
	SizeOfHeapReserve           uint32
	SizeOfHeapCommit            uint32
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               [16]DataDirectory
}

type OptionalHeader64 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               [16]DataDirectory
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

// Data reads and returns the contents of the PE section.

// Data reads and returns the contents of
// the PE section.
func (s *Section) Data() ([]byte, error)

// Open returns a new ReadSeeker reading the PE section.

// Open returns a new ReadSeeker reading
// the PE section.
func (s *Section) Open() io.ReadSeeker

type SectionHeader struct {
	Name                 string
	VirtualSize          uint32
	VirtualAddress       uint32
	Size                 uint32
	Offset               uint32
	PointerToRelocations uint32
	PointerToLineNumbers uint32
	NumberOfRelocations  uint16
	NumberOfLineNumbers  uint16
	Characteristics      uint32
}

type SectionHeader32 struct {
	Name                 [8]uint8
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLineNumbers uint32
	NumberOfRelocations  uint16
	NumberOfLineNumbers  uint16
	Characteristics      uint32
}

type Symbol struct {
	Name          string
	Value         uint32
	SectionNumber int16
	Type          uint16
	StorageClass  uint8
}
