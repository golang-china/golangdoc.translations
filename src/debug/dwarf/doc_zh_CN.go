// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package dwarf provides access to DWARF debugging information loaded from
// executable files, as defined in the DWARF 2.0 Standard at
// http://dwarfstd.org/doc/dwarf-2.0.0.pdf

// Package dwarf provides access to DWARF
// debugging information loaded from
// executable files, as defined in the
// DWARF 2.0 Standard at
// http://dwarfstd.org/doc/dwarf-2.0.0.pdf
package dwarf

// An AddrType represents a machine address type.

// An AddrType represents a machine address
// type.
type AddrType struct {
	BasicType
}

// An ArrayType represents a fixed size array type.

// An ArrayType represents a fixed size
// array type.
type ArrayType struct {
	CommonType
	Type          Type
	StrideBitSize int64 // if > 0, number of bits to hold each element
	Count         int64 // if == -1, an incomplete array, like char x[].
}

func (t *ArrayType) Size() int64

func (t *ArrayType) String() string

// An Attr identifies the attribute type in a DWARF Entry's Field.

// An Attr identifies the attribute type in
// a DWARF Entry's Field.
type Attr uint32

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

func (a Attr) GoString() string

func (a Attr) String() string

// A BasicType holds fields common to all basic types.

// A BasicType holds fields common to all
// basic types.
type BasicType struct {
	CommonType
	BitSize   int64
	BitOffset int64
}

func (b *BasicType) Basic() *BasicType

func (t *BasicType) String() string

// A BoolType represents a boolean type.

// A BoolType represents a boolean type.
type BoolType struct {
	BasicType
}

// A CharType represents a signed character type.

// A CharType represents a signed character
// type.
type CharType struct {
	BasicType
}

// A CommonType holds fields common to multiple types. If a field is not known or
// not applicable for a given type, the zero value is used.

// A CommonType holds fields common to
// multiple types. If a field is not known
// or not applicable for a given type, the
// zero value is used.
type CommonType struct {
	ByteSize int64  // size of value of this type, in bytes
	Name     string // name that can be used to refer to type
}

func (c *CommonType) Common() *CommonType

func (c *CommonType) Size() int64

// A ComplexType represents a complex floating point type.

// A ComplexType represents a complex
// floating point type.
type ComplexType struct {
	BasicType
}

// Data represents the DWARF debugging information loaded from an executable file
// (for example, an ELF or Mach-O executable).

// Data represents the DWARF debugging
// information loaded from an executable
// file (for example, an ELF or Mach-O
// executable).
type Data struct {
	// contains filtered or unexported fields
}

// New returns a new Data object initialized from the given parameters. Rather than
// calling this function directly, clients should typically use the DWARF method of
// the File type of the appropriate package debug/elf, debug/macho, or debug/pe.
//
// The []byte arguments are the data from the corresponding debug section in the
// object file; for example, for an ELF object, abbrev is the contents of the
// ".debug_abbrev" section.

// New returns a new Data object
// initialized from the given parameters.
// Rather than calling this function
// directly, clients should typically use
// the DWARF method of the File type of the
// appropriate package debug/elf,
// debug/macho, or debug/pe.
//
// The []byte arguments are the data from
// the corresponding debug section in the
// object file; for example, for an ELF
// object, abbrev is the contents of the
// ".debug_abbrev" section.
func New(abbrev, aranges, frame, info, line, pubnames, ranges, str []byte) (*Data, error)

// AddTypes will add one .debug_types section to the DWARF data. A typical object
// with DWARF version 4 debug info will have multiple .debug_types sections. The
// name is used for error reporting only, and serves to distinguish one
// .debug_types section from another.

// AddTypes will add one .debug_types
// section to the DWARF data. A typical
// object with DWARF version 4 debug info
// will have multiple .debug_types
// sections. The name is used for error
// reporting only, and serves to
// distinguish one .debug_types section
// from another.
func (d *Data) AddTypes(name string, types []byte) error

// Reader returns a new Reader for Data. The reader is positioned at byte offset 0
// in the DWARF ``info'' section.

// Reader returns a new Reader for Data.
// The reader is positioned at byte offset
// 0 in the DWARF ``info'' section.
func (d *Data) Reader() *Reader

// Type reads the type at off in the DWARF ``info'' section.

// Type reads the type at off in the DWARF
// ``info'' section.
func (d *Data) Type(off Offset) (Type, error)

type DecodeError struct {
	Name   string
	Offset Offset
	Err    string
}

func (e DecodeError) Error() string

// A DotDotDotType represents the variadic ... function parameter.

// A DotDotDotType represents the variadic
// ... function parameter.
type DotDotDotType struct {
	CommonType
}

func (t *DotDotDotType) String() string

// An entry is a sequence of attribute/value pairs.

// An entry is a sequence of
// attribute/value pairs.
type Entry struct {
	Offset   Offset // offset of Entry in DWARF info
	Tag      Tag    // tag (kind of Entry)
	Children bool   // whether Entry is followed by children
	Field    []Field
}

// Val returns the value associated with attribute Attr in Entry, or nil if there
// is no such attribute.
//
// A common idiom is to merge the check for nil return with the check that the
// value has the expected dynamic type, as in:
//
//	v, ok := e.Val(AttrSibling).(int64);

// Val returns the value associated with
// attribute Attr in Entry, or nil if there
// is no such attribute.
//
// A common idiom is to merge the check for
// nil return with the check that the value
// has the expected dynamic type, as in:
//
//	v, ok := e.Val(AttrSibling).(int64);
func (e *Entry) Val(a Attr) interface{}

// An EnumType represents an enumerated type. The only indication of its native
// integer type is its ByteSize (inside CommonType).

// An EnumType represents an enumerated
// type. The only indication of its native
// integer type is its ByteSize (inside
// CommonType).
type EnumType struct {
	CommonType
	EnumName string
	Val      []*EnumValue
}

func (t *EnumType) String() string

// An EnumValue represents a single enumeration value.

// An EnumValue represents a single
// enumeration value.
type EnumValue struct {
	Name string
	Val  int64
}

// A Field is a single attribute/value pair in an Entry.

// A Field is a single attribute/value pair
// in an Entry.
type Field struct {
	Attr Attr
	Val  interface{}
}

// A FloatType represents a floating point type.

// A FloatType represents a floating point
// type.
type FloatType struct {
	BasicType
}

// A FuncType represents a function type.

// A FuncType represents a function type.
type FuncType struct {
	CommonType
	ReturnType Type
	ParamType  []Type
}

func (t *FuncType) String() string

// An IntType represents a signed integer type.

// An IntType represents a signed integer
// type.
type IntType struct {
	BasicType
}

// An Offset represents the location of an Entry within the DWARF info. (See
// Reader.Seek.)

// An Offset represents the location of an
// Entry within the DWARF info. (See
// Reader.Seek.)
type Offset uint32

// A PtrType represents a pointer type.

// A PtrType represents a pointer type.
type PtrType struct {
	CommonType
	Type Type
}

func (t *PtrType) String() string

// A QualType represents a type that has the C/C++ "const", "restrict", or
// "volatile" qualifier.

// A QualType represents a type that has
// the C/C++ "const", "restrict", or
// "volatile" qualifier.
type QualType struct {
	CommonType
	Qual string
	Type Type
}

func (t *QualType) Size() int64

func (t *QualType) String() string

// A Reader allows reading Entry structures from a DWARF ``info'' section. The
// Entry structures are arranged in a tree. The Reader's Next function return
// successive entries from a pre-order traversal of the tree. If an entry has
// children, its Children field will be true, and the children follow, terminated
// by an Entry with Tag 0.

// A Reader allows reading Entry structures
// from a DWARF ``info'' section. The Entry
// structures are arranged in a tree. The
// Reader's Next function return successive
// entries from a pre-order traversal of
// the tree. If an entry has children, its
// Children field will be true, and the
// children follow, terminated by an Entry
// with Tag 0.
type Reader struct {
	// contains filtered or unexported fields
}

// Next reads the next entry from the encoded entry stream. It returns nil, nil
// when it reaches the end of the section. It returns an error if the current
// offset is invalid or the data at the offset cannot be decoded as a valid Entry.

// Next reads the next entry from the
// encoded entry stream. It returns nil,
// nil when it reaches the end of the
// section. It returns an error if the
// current offset is invalid or the data at
// the offset cannot be decoded as a valid
// Entry.
func (r *Reader) Next() (*Entry, error)

// Seek positions the Reader at offset off in the encoded entry stream. Offset 0
// can be used to denote the first entry.

// Seek positions the Reader at offset off
// in the encoded entry stream. Offset 0
// can be used to denote the first entry.
func (r *Reader) Seek(off Offset)

// SkipChildren skips over the child entries associated with the last Entry
// returned by Next. If that Entry did not have children or Next has not been
// called, SkipChildren is a no-op.

// SkipChildren skips over the child
// entries associated with the last Entry
// returned by Next. If that Entry did not
// have children or Next has not been
// called, SkipChildren is a no-op.
func (r *Reader) SkipChildren()

// A StructField represents a field in a struct, union, or C++ class type.

// A StructField represents a field in a
// struct, union, or C++ class type.
type StructField struct {
	Name       string
	Type       Type
	ByteOffset int64
	ByteSize   int64
	BitOffset  int64 // within the ByteSize bytes at ByteOffset
	BitSize    int64 // zero if not a bit field
}

// A StructType represents a struct, union, or C++ class type.

// A StructType represents a struct, union,
// or C++ class type.
type StructType struct {
	CommonType
	StructName string
	Kind       string // "struct", "union", or "class".
	Field      []*StructField
	Incomplete bool // if true, struct, union, class is declared but not defined
}

func (t *StructType) Defn() string

func (t *StructType) String() string

// A Tag is the classification (the type) of an Entry.

// A Tag is the classification (the type)
// of an Entry.
type Tag uint32

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

func (t Tag) GoString() string

func (t Tag) String() string

// A Type conventionally represents a pointer to any of the specific Type
// structures (CharType, StructType, etc.).

// A Type conventionally represents a
// pointer to any of the specific Type
// structures (CharType, StructType, etc.).
type Type interface {
	Common() *CommonType
	String() string
	Size() int64
}

// A TypedefType represents a named type.

// A TypedefType represents a named type.
type TypedefType struct {
	CommonType
	Type Type
}

func (t *TypedefType) Size() int64

func (t *TypedefType) String() string

// A UcharType represents an unsigned character type.

// A UcharType represents an unsigned
// character type.
type UcharType struct {
	BasicType
}

// A UintType represents an unsigned integer type.

// A UintType represents an unsigned
// integer type.
type UintType struct {
	BasicType
}

// An UnspecifiedType represents an implicit, unknown, ambiguous or nonexistent
// type.

// An UnspecifiedType represents an
// implicit, unknown, ambiguous or
// nonexistent type.
type UnspecifiedType struct {
	BasicType
}

// A VoidType represents the C void type.

// A VoidType represents the C void type.
type VoidType struct {
	CommonType
}

func (t *VoidType) String() string
