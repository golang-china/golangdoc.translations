// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

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
	Stack
}

// A Macro represents the definition of a #defined macro.
type Macro struct {
}

// A ScanToken represents an input item. It is a simple wrapping of rune, as
// returned by text/scanner.Scanner, plus a couple of extra values.
type ScanToken rune

// A Slice reads from a slice of Tokens.
type Slice struct {
}

// A Stack is a stack of TokenReaders. As the top TokenReader hits EOF,
// it resumes reading the next one down.
type Stack struct {
}

// A Token is a scan token plus its string value.
// A macro is stored as a sequence of Tokens with spaces stripped.
type Token struct {
	ScanToken
}

// A TokenReader is like a reader, but returns lex tokens of type Token. It also
// can tell you what the text of the most recently returned token is, and where
// it was found. The underlying scanner elides all spaces except newline, so the
// input looks like a stream of Tokens; original spacing is lost but we don't
// need it.
type TokenReader interface {
	// Next returns the next token.
	Next()ScanToken

	// The following methods all refer to the most recent token returned by
	// Next. Text returns the original string representation of the token.
	Text()string

	// File reports the source file name of the token.
	File()string

	// Line reports the source line number of the token.
	Line()int

	// Col reports the source column number of the token.
	Col()int

	// SetPos sets the file and line number.
	SetPos(line int, file string)

	// Close does any teardown required.
	Close()
}

// A Tokenizer is a simple wrapping of text/scanner.Scanner, configured
// for our purposes and made a TokenReader. It forms the lowest level,
// turning text from readers into tokens.
type Tokenizer struct {
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

// NewInput returns an Input from the given path.

// NewInput returns a
func NewInput(name string) *Input

// NewLexer returns a lexer for the named file and the given link context.
func NewLexer(name string, ctxt *obj.Link) TokenReader

func NewSlice(fileName string, line int, tokens []Token) *Slice

func NewTokenizer(name string, r io.Reader, file *os.File) *Tokenizer

// Tokenize turns a string into a list of Tokens; used to parse the -D flag and
// in tests.
func Tokenize(str string) []Token

func (in *Input) Close()

func (in *Input) Error(args ...interface{})

func (in *Input) Next() ScanToken

func (in *Input) Push(r TokenReader)

func (in *Input) Text() string

func (s *Slice) Close()

func (s *Slice) Col() int

func (s *Slice) File() string

func (s *Slice) Line() int

func (s *Slice) Next() ScanToken

func (s *Slice) SetPos(line int, file string)

func (s *Slice) Text() string

func (s *Stack) Close()

func (s *Stack) Col() int

func (s *Stack) File() string

func (s *Stack) Line() int

func (s *Stack) Next() ScanToken

// Push adds tr to the top (end) of the input stack. (Popping happens
// automatically.)
func (s *Stack) Push(tr TokenReader)

func (s *Stack) SetPos(line int, file string)

func (s *Stack) Text() string

func (t *Tokenizer) Close()

func (t *Tokenizer) Col() int

func (t *Tokenizer) File() string

func (t *Tokenizer) Line() int

func (t *Tokenizer) Next() ScanToken

func (t *Tokenizer) SetPos(line int, file string)

func (t *Tokenizer) Text() string

func (t ScanToken) String() string

func (l Token) String() string

