// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package token defines constants representing the lexical tokens of the Go
// programming language and basic operations on tokens (printing, predicates).

// token 包定义了表示 Go 编程语言词法的和基础运算符的常量标记.
package token

// A set of constants for precedence-based expression parsing. Non-operators have
// lowest precedence, followed by operators starting with precedence 1 up to unary
// operators. The highest precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.

// 一套表达式分析优先级常量. 非运算符的优先级最低,
// 其后是优先级为 1 开始的运算符, 一直到一元运算符.
// 优先级最高的用于 "catch-all", 总是优先于选择器, 索引以及其它运算符和定界符.
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// A File is a handle for a file belonging to a FileSet. A File has a name, size,
// and line offset table.

// File 是属于 FileSet 的文件句柄. 一个文件具有名字, 大小, 行偏移量表.
type File struct {
	// contains filtered or unexported fields
}

// AddLine adds the line offset for a new line. The line offset must be larger than
// the offset for the previous line and smaller than the file size; otherwise the
// line offset is ignored.

// AddLine 为新行增加行偏移量. 该行偏移量必须大于前一行的偏移量,
// 并且不得超过文件大小; 否则偏移量被忽略.
func (f *File) AddLine(offset int)

// AddLineInfo adds alternative file and line number information for a given file
// offset. The offset must be larger than the offset for the previously added
// alternative line info and smaller than the file size; otherwise the information
// is ignored.
//
// AddLineInfo is typically used to register alternative position information for
// //line filename:line comments in source files.

// AddLineInfo 给指定的文件添加备用文件名和行号信息.
// offset 必须大于那之前加入替换行信息的偏移, 并且必须不超过文件大小.
// 否则信息将被忽略.
//
// AddLineInfo 一般用于在源文件注册一个备用位置信息.
// //line filename:line 源文件注释
func (f *File) AddLineInfo(offset int, filename string, line int)

// Base returns the base offset of file f as registered with AddFile.

// Base 返回经 AddFile 注册的文件 f 的基础偏移量.
func (f *File) Base() int

// Line returns the line number for the given file position p; p must be a Pos
// value in that file or NoPos.

// Line 返回该文件位置 p 的行号.
// p 类型必须是 Pos 且值在这个文件中或 NoPos.
func (f *File) Line(p Pos) int

// LineCount returns the number of lines in file f.

// LineCOunt 返回文件 f 的行数.
func (f *File) LineCount() int

// MergeLine merges a line with the following line. It is akin to replacing the
// newline character at the end of the line with a space (to not change the
// remaining offsets). To obtain the line number, consult e.g. Position.Line.
// MergeLine will panic if given an invalid line number.

// MergeLine 合并 line 后的一行. 它类似于在该行末的地方替换换行符
// (不改变其余偏移量). 要获得的行号, 参考 Position.Line 的例子.
// 如果给出一个无效的行号 MergeLine 将产生 panic.
func (f *File) MergeLine(line int)

// Name returns the file name of file f as registered with AddFile.

// Name 返回经 AddFile 注册的文件 f 的文件名.
func (f *File) Name() string

// Offset returns the offset for the given file position p; p must be a valid Pos
// value in that file. f.Offset(f.Pos(offset)) == offset.

// Offset 返回给定文件位置 p 的偏移量; p 必须是文件中的有效 Pos 值.
// f.Offset(f.Pos(offset)) == offset.
func (f *File) Offset(p Pos) int

// Pos returns the Pos value for the given file offset; the offset must be <=
// f.Size(). f.Pos(f.Offset(p)) == p.

// Pos 返回给定文件偏移量的 Pos 值; offset 必须 <= f.Size().
// f.Pos(f.Offset(p)) == p.
func (f *File) Pos(offset int) Pos

// Position returns the Position value for the given file position p. Calling
// f.Position(p) is equivalent to calling f.PositionFor(p, true).

// Position 返回给定的文件位置 P 的 Position 值.
// 调用 f.Position(p) 等同调用 f.PositionFor(p, true).
func (f *File) Position(p Pos) (pos Position)

// PositionFor returns the Position value for the given file position p. If
// adjusted is set, the position may be adjusted by position-altering //line
// comments; otherwise those comments are ignored. p must be a Pos value in f or
// NoPos.

// PositionFor 返回给定的文件位置 P 的 Position 值. 如果设定了 adjusted,
// 可调整该位置被 //line comments; 的位移; 否则那些注释被忽略.
// p 必须是文件中的 Pos 值或 NoPos.
func (f *File) PositionFor(p Pos, adjusted bool) (pos Position)

// SetLines sets the line offsets for a file and returns true if successful. The
// line offsets are the offsets of the first character of each line; for instance
// for the content "ab\nc\n" the line offsets are {0, 3}. An empty file has an
// empty line offset table. Each line offset must be larger than the offset for the
// previous line and smaller than the file size; otherwise SetLines fails and
// returns false.

// SetLines 设置的所有行偏移量, 如果成功返回true.
// 行偏移量是每一行第一个字符的偏移量; 例如内容 "ab\nc\n" 的行偏移量为
// {0, 3}. 空文件也有空的行偏移量表. 每一行的偏移量必须大于上一行的,
// 并且小于文件尺寸; 否则 SetLines 失败并返回 false.
func (f *File) SetLines(lines []int) bool

// SetLinesForContent sets the line offsets for the given file content. It ignores
// position-altering //line comments.

// SetLinesForContent 设置给定文件内容的行偏移量.
// 忽略 //line comments 的位移.
func (f *File) SetLinesForContent(content []byte)

// Size returns the size of file f as registered with AddFile.

// Size 返回经 AddFile 注册的文件 f 的尺寸.
func (f *File) Size() int

// A FileSet represents a set of source files. Methods of file sets are
// synchronized; multiple goroutines may invoke them concurrently.

// FileSet 表示一个源文件集. 文件集方法是同步的; 可以多协程并发调用它们.
type FileSet struct {
	// contains filtered or unexported fields
}

// NewFileSet creates a new file set.

// NewFileSet 新建一个文件集.
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

// AddFile 新增一个文件到文件集, 给定该文件的文件名, 基本偏移量, 和文件尺寸,
// 并返回该文件. 多个文件可以同名. 基本偏移量不得小于文件集的 Base(), 也不能为负.
// 典型的, 如果提供了负数基本偏移量, 使用当前文件集的 Base() 替代.
//
// 增加文件会设置文件集的 Base() 为 base + size + 1, 作为下一个文件的最小 base 值.
// 在 Pos 值 p 和给定文件偏移量直接存在下述关系:
//
//	int(p) = base + offs
//
// offs 在 [0, size] 范围内, 因而 p 在 [base, base+size] 范围内.
// 为了方便, 可用 File.Pos 从一个文件偏移量创建文件专用的位置值.
func (s *FileSet) AddFile(filename string, base, size int) *File

// Base returns the minimum base offset that must be provided to AddFile when
// adding the next file.

// Base 返回最小基础偏移量, 添加下一个文件时, 可提供给 AddFile.
func (s *FileSet) Base() int

// File returns the file that contains the position p. If no such file is found
// (for instance for p == NoPos), the result is nil.

// File 返回包含位置 p 的文件. 如果未找到这样的文件 (例如 p == NoPos),
// 返回 nil.
func (s *FileSet) File(p Pos) (f *File)

// Iterate calls f for the files in the file set in the order they were added until
// f returns false.

// Iterate 按文件在文件集中的添加顺序调用 f, 直到 f 返回 false.
func (s *FileSet) Iterate(f func(*File) bool)

// Position converts a Pos p in the fileset into a Position value. Calling
// s.Position(p) is equivalent to calling s.PositionFor(p, true).
func (s *FileSet) Position(p Pos) (pos Position)

// PositionFor converts a Pos p in the fileset into a Position value. If adjusted
// is set, the position may be adjusted by position-altering //line comments;
// otherwise those comments are ignored. p must be a Pos value in s or NoPos.
func (s *FileSet) PositionFor(p Pos, adjusted bool) (pos Position)

// Read calls decode to deserialize a file set into s; s must not be nil.
func (s *FileSet) Read(decode func(interface{}) error) error

// Write calls encode to serialize the file set s.
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
type Pos int

// The zero value for Pos is NoPos; there is no file and line information
// associated with it, and NoPos().IsValid() is false. NoPos is always smaller than
// any other Pos value. The corresponding Position value for NoPos is the zero
// value for Position.
const NoPos Pos = 0

// IsValid returns true if the position is valid.
func (p Pos) IsValid() bool

// Position describes an arbitrary source position including the file, line, and
// column location. A Position is valid if the line number is > 0.
type Position struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (character count)
}

// IsValid returns true if the position is valid.
func (pos *Position) IsValid() bool

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	line:column         valid position without file name
//	file                invalid position with file name
//	-                   invalid position without file name
func (pos Position) String() string

// Token is the set of lexical tokens of the Go programming language.
type Token int

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
func Lookup(ident string) Token

// IsKeyword returns true for tokens corresponding to keywords; it returns false
// otherwise.

// IsKeyword 对关键字相应的标记返回 true; 其它返回 false.
func (tok Token) IsKeyword() bool

// IsLiteral returns true for tokens corresponding to identifiers and basic type
// literals; it returns false otherwise.

// IsLiteral 对标识符和基本字面类型相应的标记返回 true; 其它返回 false.
func (tok Token) IsLiteral() bool

// IsOperator returns true for tokens corresponding to operators and delimiters; it
// returns false otherwise.

// IsOperator 对运算符和定界符相应的标记返回 true; 其它返回 false.
func (tok Token) IsOperator() bool

// Precedence returns the operator precedence of the binary operator op. If op is
// not a binary operator, the result is LowestPrecedence.

// Precedence 返回二元运算符的优先级. 如果 op 不是二元操作,
// 结果是 LowestPrecedence.
func (op Token) Precedence() int

// String returns the string corresponding to the token tok. For operators,
// delimiters, and keywords the string is the actual token character sequence
// (e.g., for the token ADD, the string is "+"). For all other tokens the string
// corresponds to the token constant name (e.g. for the token IDENT, the string is
// "IDENT").

// String 返回标记 tok 相应的字符串. 对于运算符, 定界符, 以及关键字
// 该字符串是实际标记字符序列 (例如: 对于标记 ADD, 该字符串是 "+").
// 对于所有其它标记, 该字符串对应为标记常量名子
// (例如: 对于标记 IDENT, 该字符串为 "IDENT").
func (tok Token) String() string
