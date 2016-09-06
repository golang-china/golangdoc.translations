// +build ingore

// Package dwarf provides access to DWARF debugging information loaded from
// executable files, as defined in the DWARF 2.0 Standard at
// http://dwarfstd.org/doc/dwarf-2.0.0.pdf
package dwarf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
)

const (
	AttrSibling        Attr = 0x01
	AttrLocation       Attr = 0x02
	AttrName           Attr = 0x03
	AttrOrdering       Attr = 0x09
	AttrByteSize       Attr = 0x0B
	AttrBitOffset      Attr = 0x0C
	AttrBitSize        Attr = 0x0D
	AttrStmtList       Attr = 0x10
	AttrLowpc          Attr = 0x11
	AttrHighpc         Attr = 0x12
	AttrLanguage       Attr = 0x13
	AttrDiscr          Attr = 0x15
	AttrDiscrValue     Attr = 0x16
	AttrVisibility     Attr = 0x17
	AttrImport         Attr = 0x18
	AttrStringLength   Attr = 0x19
	AttrCommonRef      Attr = 0x1A
	AttrCompDir        Attr = 0x1B
	AttrConstValue     Attr = 0x1C
	AttrContainingType Attr = 0x1D
	AttrDefaultValue   Attr = 0x1E
	AttrInline         Attr = 0x20
	AttrIsOptional     Attr = 0x21
	AttrLowerBound     Attr = 0x22
	AttrProducer       Attr = 0x25
	AttrPrototyped     Attr = 0x27
	AttrReturnAddr     Attr = 0x2A
	AttrStartScope     Attr = 0x2C
	AttrStrideSize     Attr = 0x2E
	AttrUpperBound     Attr = 0x2F
	AttrAbstractOrigin Attr = 0x31
	AttrAccessibility  Attr = 0x32
	AttrAddrClass      Attr = 0x33
	AttrArtificial     Attr = 0x34
	AttrBaseTypes      Attr = 0x35
	AttrCalling        Attr = 0x36
	AttrCount          Attr = 0x37
	AttrDataMemberLoc  Attr = 0x38
	AttrDeclColumn     Attr = 0x39
	AttrDeclFile       Attr = 0x3A
	AttrDeclLine       Attr = 0x3B
	AttrDeclaration    Attr = 0x3C
	AttrDiscrList      Attr = 0x3D
	AttrEncoding       Attr = 0x3E
	AttrExternal       Attr = 0x3F
	AttrFrameBase      Attr = 0x40
	AttrFriend         Attr = 0x41
	AttrIdentifierCase Attr = 0x42
	AttrMacroInfo      Attr = 0x43
	AttrNamelistItem   Attr = 0x44
	AttrPriority       Attr = 0x45
	AttrSegment        Attr = 0x46
	AttrSpecification  Attr = 0x47
	AttrStaticLink     Attr = 0x48
	AttrType           Attr = 0x49
	AttrUseLocation    Attr = 0x4A
	AttrVarParam       Attr = 0x4B
	AttrVirtuality     Attr = 0x4C
	AttrVtableElemLoc  Attr = 0x4D
	AttrAllocated      Attr = 0x4E
	AttrAssociated     Attr = 0x4F
	AttrDataLocation   Attr = 0x50
	AttrStride         Attr = 0x51
	AttrEntrypc        Attr = 0x52
	AttrUseUTF8        Attr = 0x53
	AttrExtension      Attr = 0x54
	AttrRanges         Attr = 0x55
	AttrTrampoline     Attr = 0x56
	AttrCallColumn     Attr = 0x57
	AttrCallFile       Attr = 0x58
	AttrCallLine       Attr = 0x59
	AttrDescription    Attr = 0x5A
)

const (
	// ClassUnknown represents values of unknown DWARF class.
	ClassUnknown Class = iota

	// ClassAddress represents values of type uint64 that are
	// addresses on the target machine.
	ClassAddress

	// ClassBlock represents values of type []byte whose
	// interpretation depends on the attribute.
	ClassBlock

	// ClassConstant represents values of type int64 that are
	// constants. The interpretation of this constant depends on
	// the attribute.
	ClassConstant

	// ClassExprLoc represents values of type []byte that contain
	// an encoded DWARF expression or location description.
	ClassExprLoc

	// ClassFlag represents values of type bool.
	ClassFlag

	// ClassLinePtr represents values that are an int64 offset
	// into the "line" section.
	ClassLinePtr

	// ClassLocListPtr represents values that are an int64 offset
	// into the "loclist" section.
	ClassLocListPtr

	// ClassMacPtr represents values that are an int64 offset into
	// the "mac" section.
	ClassMacPtr

	// ClassMacPtr represents values that are an int64 offset into
	// the "rangelist" section.
	ClassRangeListPtr

	// ClassReference represents values that are an Offset offset
	// of an Entry in the info section (for use with Reader.Seek).
	// The DWARF specification combines ClassReference and
	// ClassReferenceSig into class "reference".
	ClassReference

	// ClassReferenceSig represents values that are a uint64 type
	// signature referencing a type Entry.
	ClassReferenceSig

	// ClassString represents values that are strings. If the
	// compilation unit specifies the AttrUseUTF8 flag (strongly
	// recommended), the string value will be encoded in UTF-8.
	// Otherwise, the encoding is unspecified.
	ClassString

	// ClassReferenceAlt represents values of type int64 that are
	// an offset into the DWARF "info" section of an alternate
	// object file.
	ClassReferenceAlt

	// ClassStringAlt represents values of type int64 that are an
	// offset into the DWARF string section of an alternate object
	// file.
	ClassStringAlt
)

const (
	TagArrayType              Tag = 0x01
	TagClassType              Tag = 0x02
	TagEntryPoint             Tag = 0x03
	TagEnumerationType        Tag = 0x04
	TagFormalParameter        Tag = 0x05
	TagImportedDeclaration    Tag = 0x08
	TagLabel                  Tag = 0x0A
	TagLexDwarfBlock          Tag = 0x0B
	TagMember                 Tag = 0x0D
	TagPointerType            Tag = 0x0F
	TagReferenceType          Tag = 0x10
	TagCompileUnit            Tag = 0x11
	TagStringType             Tag = 0x12
	TagStructType             Tag = 0x13
	TagSubroutineType         Tag = 0x15
	TagTypedef                Tag = 0x16
	TagUnionType              Tag = 0x17
	TagUnspecifiedParameters  Tag = 0x18
	TagVariant                Tag = 0x19
	TagCommonDwarfBlock       Tag = 0x1A
	TagCommonInclusion        Tag = 0x1B
	TagInheritance            Tag = 0x1C
	TagInlinedSubroutine      Tag = 0x1D
	TagModule                 Tag = 0x1E
	TagPtrToMemberType        Tag = 0x1F
	TagSetType                Tag = 0x20
	TagSubrangeType           Tag = 0x21
	TagWithStmt               Tag = 0x22
	TagAccessDeclaration      Tag = 0x23
	TagBaseType               Tag = 0x24
	TagCatchDwarfBlock        Tag = 0x25
	TagConstType              Tag = 0x26
	TagConstant               Tag = 0x27
	TagEnumerator             Tag = 0x28
	TagFileType               Tag = 0x29
	TagFriend                 Tag = 0x2A
	TagNamelist               Tag = 0x2B
	TagNamelistItem           Tag = 0x2C
	TagPackedType             Tag = 0x2D
	TagSubprogram             Tag = 0x2E
	TagTemplateTypeParameter  Tag = 0x2F
	TagTemplateValueParameter Tag = 0x30
	TagThrownType             Tag = 0x31
	TagTryDwarfBlock          Tag = 0x32
	TagVariantPart            Tag = 0x33
	TagVariable               Tag = 0x34
	TagVolatileType           Tag = 0x35

	// The following are new in DWARF 3.
	TagDwarfProcedure  Tag = 0x36
	TagRestrictType    Tag = 0x37
	TagInterfaceType   Tag = 0x38
	TagNamespace       Tag = 0x39
	TagImportedModule  Tag = 0x3A
	TagUnspecifiedType Tag = 0x3B
	TagPartialUnit     Tag = 0x3C
	TagImportedUnit    Tag = 0x3D
	TagMutableType     Tag = 0x3E // Later removed from DWARF.
	TagCondition       Tag = 0x3F
	TagSharedType      Tag = 0x40

	// The following are new in DWARF 4.
	TagTypeUnit            Tag = 0x41
	TagRvalueReferenceType Tag = 0x42
	TagTemplateAlias       Tag = 0x43
)

// ErrUnknownPC is the error returned by LineReader.ScanPC when the
// seek PC is not covered by any entry in the line table.
var ErrUnknownPC = errors.New("ErrUnknownPC")

// An AddrType represents a machine address type.
type AddrType struct {
	BasicType
}

// An ArrayType represents a fixed size array type.
type ArrayType struct {
	CommonType
	Type          Type
	StrideBitSize int64 // if > 0, number of bits to hold each element
	Count         int64 // if == -1, an incomplete array, like char x[].
}

// An Attr identifies the attribute type in a DWARF Entry's Field.
type Attr uint32

// A BasicType holds fields common to all basic types.
type BasicType struct {
	CommonType
	BitSize   int64
	BitOffset int64
}

// A BoolType represents a boolean type.
type BoolType struct {
	BasicType
}

// A CharType represents a signed character type.
type CharType struct {
	BasicType
}

// A Class is the DWARF 4 class of an attribute value.
//
// In general, a given attribute's value may take on one of several
// possible classes defined by DWARF, each of which leads to a
// slightly different interpretation of the attribute.
//
// DWARF version 4 distinguishes attribute value classes more finely
// than previous versions of DWARF. The reader will disambiguate
// coarser classes from earlier versions of DWARF into the appropriate
// DWARF 4 class. For example, DWARF 2 uses "constant" for constants
// as well as all types of section offsets, but the reader will
// canonicalize attributes in DWARF 2 files that refer to section
// offsets to one of the Class*Ptr classes, even though these classes
// were only defined in DWARF 3.
type Class int

// A CommonType holds fields common to multiple types.
// If a field is not known or not applicable for a given type,
// the zero value is used.

// A CommonType holds fields common to multiple types. If a field is not known
// or not applicable for a given type, the zero value is used.
type CommonType struct {
	ByteSize int64  // size of value of this type, in bytes
	Name     string // name that can be used to refer to type
}

// A ComplexType represents a complex floating point type.
type ComplexType struct {
	BasicType
}

// Data represents the DWARF debugging information loaded from an executable
// file (for example, an ELF or Mach-O executable).
type Data struct {
}

type DecodeError struct {
	Name   string
	Offset Offset
	Err    string
}

// A DotDotDotType represents the variadic ... function parameter.
type DotDotDotType struct {
	CommonType
}

// An entry is a sequence of attribute/value pairs.
type Entry struct {
	Offset   Offset // offset of Entry in DWARF info
	Tag      Tag    // tag (kind of Entry)
	Children bool   // whether Entry is followed by children
	Field    []Field
}

// An EnumType represents an enumerated type.
// The only indication of its native integer type is its ByteSize
// (inside CommonType).

// An EnumType represents an enumerated type. The only indication of its native
// integer type is its ByteSize (inside CommonType).
type EnumType struct {
	CommonType
	EnumName string
	Val      []*EnumValue
}

// An EnumValue represents a single enumeration value.
type EnumValue struct {
	Name string
	Val  int64
}

// A Field is a single attribute/value pair in an Entry.
//
// A value can be one of several "attribute classes" defined by DWARF.
// The Go types corresponding to each class are:
//
//    DWARF class       Go type        Class
//    -----------       -------        -----
//    address           uint64         ClassAddress
//    block             []byte         ClassBlock
//    constant          int64          ClassConstant
//    flag              bool           ClassFlag
//    reference
//      to info         dwarf.Offset   ClassReference
//      to type unit    uint64         ClassReferenceSig
//    string            string         ClassString
//    exprloc           []byte         ClassExprLoc
//    lineptr           int64          ClassLinePtr
//    loclistptr        int64          ClassLocListPtr
//    macptr            int64          ClassMacPtr
//    rangelistptr      int64          ClassRangeListPtr
//
// For unrecognized or vendor-defined attributes, Class may be
// ClassUnknown.

// A Field is a single attribute/value pair in an Entry.
type Field struct {
	Attr  Attr
	Val   interface{}
	Class Class
}

// A FloatType represents a floating point type.
type FloatType struct {
	BasicType
}

// A FuncType represents a function type.
type FuncType struct {
	CommonType
	ReturnType Type
	ParamType  []Type
}

// An IntType represents a signed integer type.
type IntType struct {
	BasicType
}

// A LineEntry is a row in a DWARF line table.
type LineEntry struct {
	// Address is the program-counter value of a machine
	// instruction generated by the compiler. This LineEntry
	// applies to each instruction from Address to just before the
	// Address of the next LineEntry.
	Address uint64

	// OpIndex is the index of an operation within a VLIW
	// instruction. The index of the first operation is 0. For
	// non-VLIW architectures, it will always be 0. Address and
	// OpIndex together form an operation pointer that can
	// reference any individual operation within the instruction
	// stream.
	OpIndex int

	// File is the source file corresponding to these
	// instructions.
	File *LineFile

	// Line is the source code line number corresponding to these
	// instructions. Lines are numbered beginning at 1. It may be
	// 0 if these instructions cannot be attributed to any source
	// line.
	Line int

	// Column is the column number within the source line of these
	// instructions. Columns are numbered beginning at 1. It may
	// be 0 to indicate the "left edge" of the line.
	Column int

	// IsStmt indicates that Address is a recommended breakpoint
	// location, such as the beginning of a line, statement, or a
	// distinct subpart of a statement.
	IsStmt bool

	// BasicBlock indicates that Address is the beginning of a
	// basic block.
	BasicBlock bool

	// PrologueEnd indicates that Address is one (of possibly
	// many) PCs where execution should be suspended for a
	// breakpoint on entry to the containing function.
	//
	// Added in DWARF 3.
	PrologueEnd bool

	// EpilogueBegin indicates that Address is one (of possibly
	// many) PCs where execution should be suspended for a
	// breakpoint on exit from this function.
	//
	// Added in DWARF 3.
	EpilogueBegin bool

	// ISA is the instruction set architecture for these
	// instructions. Possible ISA values should be defined by the
	// applicable ABI specification.
	//
	// Added in DWARF 3.
	ISA int

	// Discriminator is an arbitrary integer indicating the block
	// to which these instructions belong. It serves to
	// distinguish among multiple blocks that may all have with
	// the same source file, line, and column. Where only one
	// block exists for a given source position, it should be 0.
	//
	// Added in DWARF 3.
	Discriminator int

	// EndSequence indicates that Address is the first byte after
	// the end of a sequence of target machine instructions. If it
	// is set, only this and the Address field are meaningful. A
	// line number table may contain information for multiple
	// potentially disjoint instruction sequences. The last entry
	// in a line table should always have EndSequence set.
	EndSequence bool
}

// A LineFile is a source file referenced by a DWARF line table entry.
type LineFile struct {
	Name   string
	Mtime  uint64 // Implementation defined modification time, or 0 if unknown
	Length int    // File length, or 0 if unknown
}

// A LineReader reads a sequence of LineEntry structures from a DWARF
// "line" section for a single compilation unit. LineEntries occur in
// order of increasing PC and each LineEntry gives metadata for the
// instructions from that LineEntry's PC to just before the next
// LineEntry's PC. The last entry will have its EndSequence field set.
type LineReader struct {
}

// A LineReaderPos represents a position in a line table.
type LineReaderPos struct {
}

// An Offset represents the location of an Entry within the DWARF info.
// (See Reader.Seek.)

// An Offset represents the location of an Entry within the DWARF info. (See
// Reader.Seek.)
type Offset uint32

// A PtrType represents a pointer type.
type PtrType struct {
	CommonType
	Type Type
}

// A QualType represents a type that has the C/C++ "const", "restrict", or
// "volatile" qualifier.
type QualType struct {
	CommonType
	Qual string
	Type Type
}

// A Reader allows reading Entry structures from a DWARF ``info'' section. The
// Entry structures are arranged in a tree. The Reader's Next function return
// successive entries from a pre-order traversal of the tree. If an entry has
// children, its Children field will be true, and the children follow,
// terminated by an Entry with Tag 0.
type Reader struct {
}

// A StructField represents a field in a struct, union, or C++ class type.
type StructField struct {
	Name       string
	Type       Type
	ByteOffset int64
	ByteSize   int64
	BitOffset  int64 // within the ByteSize bytes at ByteOffset
	BitSize    int64 // zero if not a bit field
}

// A StructType represents a struct, union, or C++ class type.
type StructType struct {
	CommonType
	StructName string
	Kind       string // "struct", "union", or "class".
	Field      []*StructField
	Incomplete bool // if true, struct, union, class is declared but not defined
}

// A Tag is the classification (the type) of an Entry.
type Tag uint32

// A Type conventionally represents a pointer to any of the
// specific Type structures (CharType, StructType, etc.).

// A Type conventionally represents a pointer to any of the specific Type
// structures (CharType, StructType, etc.).
type Type interface {
	Common()*CommonType
	String()string
	Size()int64
}

// A TypedefType represents a named type.
type TypedefType struct {
	CommonType
	Type Type
}

// A UcharType represents an unsigned character type.
type UcharType struct {
	BasicType
}

// A UintType represents an unsigned integer type.
type UintType struct {
	BasicType
}

// An UnspecifiedType represents an implicit, unknown, ambiguous or nonexistent
// type.
type UnspecifiedType struct {
	BasicType
}

// A VoidType represents the C void type.
type VoidType struct {
	CommonType
}

// New returns a new Data object initialized from the given parameters. Rather
// than calling this function directly, clients should typically use the DWARF
// method of the File type of the appropriate package debug/elf, debug/macho, or
// debug/pe.
//
// The []byte arguments are the data from the corresponding debug section in the
// object file; for example, for an ELF object, abbrev is the contents of the
// ".debug_abbrev" section.
func New(abbrev, aranges, frame, info, line, pubnames, ranges, str []byte) (*Data, error)

func (t *ArrayType) Size() int64

func (t *ArrayType) String() string

func (b *BasicType) Basic() *BasicType

func (t *BasicType) String() string

func (c *CommonType) Common() *CommonType

func (c *CommonType) Size() int64

// AddTypes will add one .debug_types section to the DWARF data. A
// typical object with DWARF version 4 debug info will have multiple
// .debug_types sections. The name is used for error reporting only,
// and serves to distinguish one .debug_types section from another.

// AddTypes will add one .debug_types section to the DWARF data. A typical
// object with DWARF version 4 debug info will have multiple .debug_types
// sections. The name is used for error reporting only, and serves to
// distinguish one .debug_types section from another.
func (d *Data) AddTypes(name string, types []byte) error

// LineReader returns a new reader for the line table of compilation
// unit cu, which must be an Entry with tag TagCompileUnit.
//
// If this compilation unit has no line table, it returns nil, nil.
func (d *Data) LineReader(cu *Entry) (*LineReader, error)

// Ranges returns the PC ranges covered by e, a slice of [low,high) pairs.
// Only some entry types, such as TagCompileUnit or TagSubprogram, have PC
// ranges; for others, this will return nil with no error.
func (d *Data) Ranges(e *Entry) ([][2]uint64, error)

// Reader returns a new Reader for Data.
// The reader is positioned at byte offset 0 in the DWARF ``info'' section.

// Reader returns a new Reader for Data. The reader is positioned at byte offset
// 0 in the DWARF ``info'' section.
func (d *Data) Reader() *Reader

// Type reads the type at off in the DWARF ``info'' section.
func (d *Data) Type(off Offset) (Type, error)

func (t *DotDotDotType) String() string

// AttrField returns the Field associated with attribute Attr in
// Entry, or nil if there is no such attribute.
func (e *Entry) AttrField(a Attr) *Field

// Val returns the value associated with attribute Attr in Entry,
// or nil if there is no such attribute.
//
// A common idiom is to merge the check for nil return with
// the check that the value has the expected dynamic type, as in:
// 	v, ok := e.Val(AttrSibling).(int64)

// Val returns the value associated with attribute Attr in Entry, or nil if
// there is no such attribute.
//
// A common idiom is to merge the check for nil return with the check that the
// value has the expected dynamic type, as in:
//
// 	v, ok := e.Val(AttrSibling).(int64);
func (e *Entry) Val(a Attr) interface{}

func (t *EnumType) String() string

func (t *FuncType) String() string

// Next sets *entry to the next row in this line table and moves to
// the next row. If there are no more entries and the line table is
// properly terminated, it returns io.EOF.
//
// Rows are always in order of increasing entry.Address, but
// entry.Line may go forward or backward.
func (r *LineReader) Next(entry *LineEntry) error

// Reset repositions the line table reader at the beginning of the
// line table.
func (r *LineReader) Reset()

// Seek restores the line table reader to a position returned by Tell.
//
// The argument pos must have been returned by a call to Tell on this
// line table.
func (r *LineReader) Seek(pos LineReaderPos)

// SeekPC sets *entry to the LineEntry that includes pc and positions
// the reader on the next entry in the line table. If necessary, this
// will seek backwards to find pc.
//
// If pc is not covered by any entry in this line table, SeekPC
// returns ErrUnknownPC. In this case, *entry and the final seek
// position are unspecified.
//
// Note that DWARF line tables only permit sequential, forward scans.
// Hence, in the worst case, this takes time linear in the size of the
// line table. If the caller wishes to do repeated fast PC lookups, it
// should build an appropriate index of the line table.
func (r *LineReader) SeekPC(pc uint64, entry *LineEntry) error

// Tell returns the current position in the line table.
func (r *LineReader) Tell() LineReaderPos

func (t *PtrType) String() string

func (t *QualType) Size() int64

func (t *QualType) String() string

// AddressSize returns the size in bytes of addresses in the current compilation
// unit.
func (r *Reader) AddressSize() int

// Next reads the next entry from the encoded entry stream.
// It returns nil, nil when it reaches the end of the section.
// It returns an error if the current offset is invalid or the data at the
// offset cannot be decoded as a valid Entry.

// Next reads the next entry from the encoded entry stream. It returns nil, nil
// when it reaches the end of the section. It returns an error if the current
// offset is invalid or the data at the offset cannot be decoded as a valid
// Entry.
func (r *Reader) Next() (*Entry, error)

// Seek positions the Reader at offset off in the encoded entry stream.
// Offset 0 can be used to denote the first entry.

// Seek positions the Reader at offset off in the encoded entry stream. Offset 0
// can be used to denote the first entry.
func (r *Reader) Seek(off Offset)

// SeekPC returns the Entry for the compilation unit that includes pc,
// and positions the reader to read the children of that unit.  If pc
// is not covered by any unit, SeekPC returns ErrUnknownPC and the
// position of the reader is undefined.
//
// Because compilation units can describe multiple regions of the
// executable, in the worst case SeekPC must search through all the
// ranges in all the compilation units. Each call to SeekPC starts the
// search at the compilation unit of the last call, so in general
// looking up a series of PCs will be faster if they are sorted. If
// the caller wishes to do repeated fast PC lookups, it should build
// an appropriate index using the Ranges method.
func (r *Reader) SeekPC(pc uint64) (*Entry, error)

// SkipChildren skips over the child entries associated with
// the last Entry returned by Next. If that Entry did not have
// children or Next has not been called, SkipChildren is a no-op.

// SkipChildren skips over the child entries associated with the last Entry
// returned by Next. If that Entry did not have children or Next has not been
// called, SkipChildren is a no-op.
func (r *Reader) SkipChildren()

func (t *StructType) Defn() string

func (t *StructType) String() string

func (t *TypedefType) Size() int64

func (t *TypedefType) String() string

func (t *VoidType) String() string

func (a Attr) GoString() string

func (a Attr) String() string

func (i Class) GoString() string

func (i Class) String() string

func (e DecodeError) Error() string

func (t Tag) GoString() string

func (t Tag) String() string

