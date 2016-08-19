// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package lex implements lexical analysis for the assembler.

// Package lex implements lexical analysis for the assembler.
package lex

import (
    "cmd/asm/internal/flags"
    "cmd/internal/obj"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "text/scanner"
    "unicode"
)


const (
	// Asm defines some two-character lexemes. We make up
	// a rune/ScanToken value for them - ugly but simple.
	LSH ScanToken = -1000 - iota // << Left shift.
	RSH                          // >> Logical right shift.
	ARR                          // -> Used on ARM for shift type 3, arithmetic right shift.
	ROT                          // @> Used on ARM for shift type 4, rotate right.

)


// Input is the main input: a stack of readers and some macro definitions.
// It also handles #include processing (by pushing onto the input stack)
// and parses and instantiates macro definitions.
type Input struct {
	includes        []string
	beginningOfLine bool
	ifdefStack      []bool
	macros          map[string]*Macro
	text            string // Text of last token returned by Next.
	peek            bool
	peekToken       ScanToken
	peekText        string
}


// A Macro represents the definition of a #defined macro.
type Macro struct {
	name   string   // The #define name.
	args   []string // Formal arguments.
	tokens []Token  // Body of macro.
}


// A ScanToken represents an input item. It is a simple wrapping of rune, as
// returned by text/scanner.Scanner, plus a couple of extra values.
type ScanToken rune


// A Slice reads from a slice of Tokens.
type Slice struct {
	tokens   []Token
	fileName string
	line     int
	pos      int
}


// A Stack is a stack of TokenReaders. As the top TokenReader hits EOF,
// it resumes reading the next one down.
type Stack struct {
	tr []TokenReader
}


// A Token is a scan token plus its string value.
// A macro is stored as a sequence of Tokens with spaces stripped.
type Token struct {
	text string
}


// A TokenReader is like a reader, but returns lex tokens of type Token. It also
// can tell you what the text of the most recently returned token is, and where
// it was found. The underlying scanner elides all spaces except newline, so the
// input looks like a stream of Tokens; original spacing is lost but we don't
// need it.

// A TokenReader is like a reader, but returns lex tokens of type Token. It also
// can tell you what the text of the most recently returned token is, and where
// it was found. The underlying scanner elides all spaces except newline, so the
// input looks like a stream of Tokens; original spacing is lost but we don't
// need it.
type TokenReader interface {
	// Next returns the next token.
	Next() ScanToken
	// The following methods all refer to the most recent token returned by Next.
	// Text returns the original string representation of the token.
	Text() string
	// File reports the source file name of the token.
	File() string
	// Line reports the source line number of the token.
	Line() int
	// Col reports the source column number of the token.
	Col() int
	// SetPos sets the file and line number.
	SetPos(line int, file string)
	// Close does any teardown required.
	Close()
}


// A Tokenizer is a simple wrapping of text/scanner.Scanner, configured
// for our purposes and made a TokenReader. It forms the lowest level,
// turning text from readers into tokens.
type Tokenizer struct {
	tok      ScanToken
	s        *scanner.Scanner
	line     int
	fileName string
	file     *os.File // If non-nil, file descriptor to close.
}


// HistLine reports the cumulative source line number of the token,
// for use in the Prog structure for the linker. (It's always handling the
// instruction from the current lex line.)
// It returns int32 because that's what type ../asm prefers.
func HistLine() int32

// InitHist sets the line count to 1, for reproducible testing.
func InitHist()

// IsRegisterShift reports whether the token is one of the ARM register shift
// operators.
func IsRegisterShift(r ScanToken) bool

// Make returns a Token with the given rune (ScanToken) and text representation.
func Make(token ScanToken, text string) Token

// NewInput returns a

// NewInput returns an Input from the given path.
func NewInput(name string) *Input

// NewLexer returns a lexer for the named file and the given link context.
func NewLexer(name string, ctxt *obj.Link) TokenReader

func NewSlice(fileName string, line int, tokens []Token) *Slice

func NewTokenizer(name string, r io.Reader, file *os.File) *Tokenizer

// Tokenize turns a string into a list of Tokens; used to parse the -D flag and
// in tests.
func Tokenize(str string) []Token

func (*Input) Close()

func (*Input) Error(args ...interface{})

func (*Input) Next() ScanToken

func (*Input) Push(r TokenReader)

func (*Input) Text() string

func (*Slice) Close()

func (*Slice) Col() int

func (*Slice) File() string

func (*Slice) Line() int

func (*Slice) Next() ScanToken

func (*Slice) SetPos(line int, file string)

func (*Slice) Text() string

func (*Stack) Close()

func (*Stack) Col() int

func (*Stack) File() string

func (*Stack) Line() int

func (*Stack) Next() ScanToken

// Push adds tr to the top (end) of the input stack. (Popping happens
// automatically.)
func (*Stack) Push(tr TokenReader)

func (*Stack) SetPos(line int, file string)

func (*Stack) Text() string

func (*Tokenizer) Close()

func (*Tokenizer) Col() int

func (*Tokenizer) File() string

func (*Tokenizer) Line() int

func (*Tokenizer) Next() ScanToken

func (*Tokenizer) SetPos(line int, file string)

func (*Tokenizer) Text() string

func (ScanToken) String() string

func (Token) String() string

