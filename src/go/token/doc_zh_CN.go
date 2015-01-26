// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package token defines constants representing the lexical tokens of the Go
// programming language and basic operations on tokens (printing, predicates).

// Package token defines constants
// representing the lexical tokens of the
// Go programming language and basic
// operations on tokens (printing,
// predicates).
package token

// A set of constants for precedence-based expression parsing. Non-operators have
// lowest precedence, followed by operators starting with precedence 1 up to unary
// operators. The highest precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.

// A set of constants for precedence-based
// expression parsing. Non-operators have
// lowest precedence, followed by operators
// starting with precedence 1 up to unary
// operators. The highest precedence serves
// as "catch-all" precedence for selector,
// indexing, and other operator and
// delimiter tokens.
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// A File is a handle for a file belonging to a FileSet. A File has a name, size,
// and line offset table.

// A File is a handle for a file belonging
// to a FileSet. A File has a name, size,
// and line offset table.
type File struct {
	// contains filtered or unexported fields
}

// AddLine adds the line offset for a new line. The line offset must be larger than
// the offset for the previous line and smaller than the file size; otherwise the
// line offset is ignored.

// AddLine adds the line offset for a new
// line. The line offset must be larger
// than the offset for the previous line
// and smaller than the file size;
// otherwise the line offset is ignored.
func (f *File) AddLine(offset int)

// AddLineInfo adds alternative file and line number information for a given file
// offset. The offset must be larger than the offset for the previously added
// alternative line info and smaller than the file size; otherwise the information
// is ignored.
//
// AddLineInfo is typically used to register alternative position information for
// //line filename:line comments in source files.

// AddLineInfo adds alternative file and
// line number information for a given file
// offset. The offset must be larger than
// the offset for the previously added
// alternative line info and smaller than
// the file size; otherwise the information
// is ignored.
//
// AddLineInfo is typically used to
// register alternative position
// information for //line filename:line
// comments in source files.
func (f *File) AddLineInfo(offset int, filename string, line int)

// Base returns the base offset of file f as registered with AddFile.

// Base returns the base offset of file f
// as registered with AddFile.
func (f *File) Base() int

// Line returns the line number for the given file position p; p must be a Pos
// value in that file or NoPos.

// Line returns the line number for the
// given file position p; p must be a Pos
// value in that file or NoPos.
func (f *File) Line(p Pos) int

// LineCount returns the number of lines in file f.

// LineCount returns the number of lines in
// file f.
func (f *File) LineCount() int

// MergeLine merges a line with the following line. It is akin to replacing the
// newline character at the end of the line with a space (to not change the
// remaining offsets). To obtain the line number, consult e.g. Position.Line.
// MergeLine will panic if given an invalid line number.

// MergeLine merges a line with the
// following line. It is akin to replacing
// the newline character at the end of the
// line with a space (to not change the
// remaining offsets). To obtain the line
// number, consult e.g. Position.Line.
// MergeLine will panic if given an invalid
// line number.
func (f *File) MergeLine(line int)

// Name returns the file name of file f as registered with AddFile.

// Name returns the file name of file f as
// registered with AddFile.
func (f *File) Name() string

// Offset returns the offset for the given file position p; p must be a valid Pos
// value in that file. f.Offset(f.Pos(offset)) == offset.

// Offset returns the offset for the given
// file position p; p must be a valid Pos
// value in that file.
// f.Offset(f.Pos(offset)) == offset.
func (f *File) Offset(p Pos) int

// Pos returns the Pos value for the given file offset; the offset must be <=
// f.Size(). f.Pos(f.Offset(p)) == p.

// Pos returns the Pos value for the given
// file offset; the offset must be <=
// f.Size(). f.Pos(f.Offset(p)) == p.
func (f *File) Pos(offset int) Pos

// Position returns the Position value for the given file position p. Calling
// f.Position(p) is equivalent to calling f.PositionFor(p, true).

// Position returns the Position value for
// the given file position p. Calling
// f.Position(p) is equivalent to calling
// f.PositionFor(p, true).
func (f *File) Position(p Pos) (pos Position)

// PositionFor returns the Position value for the given file position p. If
// adjusted is set, the position may be adjusted by position-altering //line
// comments; otherwise those comments are ignored. p must be a Pos value in f or
// NoPos.

// PositionFor returns the Position value
// for the given file position p. If
// adjusted is set, the position may be
// adjusted by position-altering //line
// comments; otherwise those comments are
// ignored. p must be a Pos value in f or
// NoPos.
func (f *File) PositionFor(p Pos, adjusted bool) (pos Position)

// SetLines sets the line offsets for a file and returns true if successful. The
// line offsets are the offsets of the first character of each line; for instance
// for the content "ab\nc\n" the line offsets are {0, 3}. An empty file has an
// empty line offset table. Each line offset must be larger than the offset for the
// previous line and smaller than the file size; otherwise SetLines fails and
// returns false.

// SetLines sets the line offsets for a
// file and returns true if successful. The
// line offsets are the offsets of the
// first character of each line; for
// instance for the content "ab\nc\n" the
// line offsets are {0, 3}. An empty file
// has an empty line offset table. Each
// line offset must be larger than the
// offset for the previous line and smaller
// than the file size; otherwise SetLines
// fails and returns false.
func (f *File) SetLines(lines []int) bool

// SetLinesForContent sets the line offsets for the given file content. It ignores
// position-altering //line comments.

// SetLinesForContent sets the line offsets
// for the given file content. It ignores
// position-altering //line comments.
func (f *File) SetLinesForContent(content []byte)

// Size returns the size of file f as registered with AddFile.

// Size returns the size of file f as
// registered with AddFile.
func (f *File) Size() int

// A FileSet represents a set of source files. Methods of file sets are
// synchronized; multiple goroutines may invoke them concurrently.

// A FileSet represents a set of source
// files. Methods of file sets are
// synchronized; multiple goroutines may
// invoke them concurrently.
type FileSet struct {
	// contains filtered or unexported fields
}

// NewFileSet creates a new file set.

// NewFileSet creates a new file set.
func NewFileSet() *FileSet

// AddFile adds a new file with a given filename, base offset, and file size to the
// file set s and returns the file. Multiple files may have the same name. The base
// offset must not be smaller than the FileSet's Base(), and size must not be
// negative. As a special case, if a negative base is provided, the current value
// of the FileSet's Base() is used instead.
//
// Adding the file will set the file set's Base() value to base + size + 1 as the
// minimum base value for the next file. The following relationship exists between
// a Pos value p for a given file offset offs:
//
//	int(p) = base + offs
//
// with offs in the range [0, size] and thus p in the range [base, base+size]. For
// convenience, File.Pos may be used to create file-specific position values from a
// file offset.

// AddFile adds a new file with a given
// filename, base offset, and file size to
// the file set s and returns the file.
// Multiple files may have the same name.
// The base offset must not be smaller than
// the FileSet's Base(), and size must not
// be negative. As a special case, if a
// negative base is provided, the current
// value of the FileSet's Base() is used
// instead.
//
// Adding the file will set the file set's
// Base() value to base + size + 1 as the
// minimum base value for the next file.
// The following relationship exists
// between a Pos value p for a given file
// offset offs:
//
//	int(p) = base + offs
//
// with offs in the range [0, size] and
// thus p in the range [base, base+size].
// For convenience, File.Pos may be used to
// create file-specific position values
// from a file offset.
func (s *FileSet) AddFile(filename string, base, size int) *File

// Base returns the minimum base offset that must be provided to AddFile when
// adding the next file.

// Base returns the minimum base offset
// that must be provided to AddFile when
// adding the next file.
func (s *FileSet) Base() int

// File returns the file that contains the position p. If no such file is found
// (for instance for p == NoPos), the result is nil.

// File returns the file that contains the
// position p. If no such file is found
// (for instance for p == NoPos), the
// result is nil.
func (s *FileSet) File(p Pos) (f *File)

// Iterate calls f for the files in the file set in the order they were added until
// f returns false.

// Iterate calls f for the files in the
// file set in the order they were added
// until f returns false.
func (s *FileSet) Iterate(f func(*File) bool)

// Position converts a Pos p in the fileset into a Position value. Calling
// s.Position(p) is equivalent to calling s.PositionFor(p, true).

// Position converts a Pos p in the fileset
// into a Position value. Calling
// s.Position(p) is equivalent to calling
// s.PositionFor(p, true).
func (s *FileSet) Position(p Pos) (pos Position)

// PositionFor converts a Pos p in the fileset into a Position value. If adjusted
// is set, the position may be adjusted by position-altering //line comments;
// otherwise those comments are ignored. p must be a Pos value in s or NoPos.

// PositionFor converts a Pos p in the
// fileset into a Position value. If
// adjusted is set, the position may be
// adjusted by position-altering //line
// comments; otherwise those comments are
// ignored. p must be a Pos value in s or
// NoPos.
func (s *FileSet) PositionFor(p Pos, adjusted bool) (pos Position)

// Read calls decode to deserialize a file set into s; s must not be nil.

// Read calls decode to deserialize a file
// set into s; s must not be nil.
func (s *FileSet) Read(decode func(interface{}) error) error

// Write calls encode to serialize the file set s.

// Write calls encode to serialize the file
// set s.
func (s *FileSet) Write(encode func(interface{}) error) error

// Pos is a compact encoding of a source position within a file set. It can be
// converted into a Position for a more convenient, but much larger,
// representation.
//
// The Pos value for a given file is a number in the range [base, base+size], where
// base and size are specified when adding the file to the file set via AddFile.
//
// To create the Pos value for a specific source offset, first add the respective
// file to the current file set (via FileSet.AddFile) and then call
// File.Pos(offset) for that file. Given a Pos value p for a specific file set
// fset, the corresponding Position value is obtained by calling fset.Position(p).
//
// Pos values can be compared directly with the usual comparison operators: If two
// Pos values p and q are in the same file, comparing p and q is equivalent to
// comparing the respective source file offsets. If p and q are in different files,
// p < q is true if the file implied by p was added to the respective file set
// before the file implied by q.

// Pos is a compact encoding of a source
// position within a file set. It can be
// converted into a Position for a more
// convenient, but much larger,
// representation.
//
// The Pos value for a given file is a
// number in the range [base, base+size],
// where base and size are specified when
// adding the file to the file set via
// AddFile.
//
// To create the Pos value for a specific
// source offset, first add the respective
// file to the current file set (via
// FileSet.AddFile) and then call
// File.Pos(offset) for that file. Given a
// Pos value p for a specific file set
// fset, the corresponding Position value
// is obtained by calling fset.Position(p).
//
// Pos values can be compared directly with
// the usual comparison operators: If two
// Pos values p and q are in the same file,
// comparing p and q is equivalent to
// comparing the respective source file
// offsets. If p and q are in different
// files, p < q is true if the file implied
// by p was added to the respective file
// set before the file implied by q.
type Pos int

// The zero value for Pos is NoPos; there is no file and line information
// associated with it, and NoPos().IsValid() is false. NoPos is always smaller than
// any other Pos value. The corresponding Position value for NoPos is the zero
// value for Position.

// The zero value for Pos is NoPos; there
// is no file and line information
// associated with it, and
// NoPos().IsValid() is false. NoPos is
// always smaller than any other Pos value.
// The corresponding Position value for
// NoPos is the zero value for Position.
const NoPos Pos = 0

// IsValid returns true if the position is valid.

// IsValid returns true if the position is
// valid.
func (p Pos) IsValid() bool

// Position describes an arbitrary source position including the file, line, and
// column location. A Position is valid if the line number is > 0.

// Position describes an arbitrary source
// position including the file, line, and
// column location. A Position is valid if
// the line number is > 0.
type Position struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (character count)
}

// IsValid returns true if the position is valid.

// IsValid returns true if the position is
// valid.
func (pos *Position) IsValid() bool

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	line:column         valid position without file name
//	file                invalid position with file name
//	-                   invalid position without file name

// String returns a string in one of
// several forms:
//
//	file:line:column    valid position with file name
//	line:column         valid position without file name
//	file                invalid position with file name
//	-                   invalid position without file name
func (pos Position) String() string

// Token is the set of lexical tokens of the Go programming language.

// Token is the set of lexical tokens of
// the Go programming language.
type Token int

// The list of tokens.

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"

	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	// Keywords
	BREAK
	CASE
	CHAN
	CONST
	CONTINUE

	DEFAULT
	DEFER
	ELSE
	FALLTHROUGH
	FOR

	FUNC
	GO
	GOTO
	IF
	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN

	SELECT
	STRUCT
	SWITCH
	TYPE
	VAR
)

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).

// Lookup maps an identifier to its keyword
// token or IDENT (if not a keyword).
func Lookup(ident string) Token

// IsKeyword returns true for tokens corresponding to keywords; it returns false
// otherwise.

// IsKeyword returns true for tokens
// corresponding to keywords; it returns
// false otherwise.
func (tok Token) IsKeyword() bool

// IsLiteral returns true for tokens corresponding to identifiers and basic type
// literals; it returns false otherwise.

// IsLiteral returns true for tokens
// corresponding to identifiers and basic
// type literals; it returns false
// otherwise.
func (tok Token) IsLiteral() bool

// IsOperator returns true for tokens corresponding to operators and delimiters; it
// returns false otherwise.

// IsOperator returns true for tokens
// corresponding to operators and
// delimiters; it returns false otherwise.
func (tok Token) IsOperator() bool

// Precedence returns the operator precedence of the binary operator op. If op is
// not a binary operator, the result is LowestPrecedence.

// Precedence returns the operator
// precedence of the binary operator op. If
// op is not a binary operator, the result
// is LowestPrecedence.
func (op Token) Precedence() int

// String returns the string corresponding to the token tok. For operators,
// delimiters, and keywords the string is the actual token character sequence
// (e.g., for the token ADD, the string is "+"). For all other tokens the string
// corresponds to the token constant name (e.g. for the token IDENT, the string is
// "IDENT").

// String returns the string corresponding
// to the token tok. For operators,
// delimiters, and keywords the string is
// the actual token character sequence
// (e.g., for the token ADD, the string is
// "+"). For all other tokens the string
// corresponds to the token constant name
// (e.g. for the token IDENT, the string is
// "IDENT").
func (tok Token) String() string
