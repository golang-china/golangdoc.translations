// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package asm implements the parser and instruction generator for the assembler.
// TODO: Split apart?
package asm // import "cmd/asm/internal/asm"

import (
    "bytes"
    "cmd/asm/internal/arch"
    "cmd/asm/internal/flags"
    "cmd/asm/internal/lex"
    "cmd/internal/obj"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strconv"
    "strings"
    "testing"
    "text/scanner"
    "unicode/utf8"
)

// EOF represents the end of input.
var EOF = lex.Make(scanner.EOF, "EOF")

type Parser struct {
    lex           lex.TokenReader
    lineNum       int   // Line number in source file.
    histLineNum   int32 // Cumulative line number across source files.
    errorLine     int32 // (Cumulative) line number of last error.
    errorCount    int   // Number of errors.
    pc            int64 // virtual PC; count of Progs; doesn't advance for GLOBL or DATA.
    input         []lex.Token
    inputPos      int
    pendingLabels []string // Labels to attach to next instruction.
    labels        map[string]*obj.Prog
    toPatch       []Patch
    addr          []obj.Addr
    arch          *arch.Arch
    ctxt          *obj.Link
    firstProg     *obj.Prog
    lastProg      *obj.Prog
    dataAddr      map[string]int64 // Most recent address for DATA for this symbol.
    isJump        bool             // Instruction being assembled is a jump.
    errorWriter   io.Writer
}

type Patch struct {
    prog  *obj.Prog
    label string
}

func NewParser(ctxt *obj.Link, ar *arch.Arch, lexer lex.TokenReader) *Parser

func Test386EndToEnd(t *testing.T)

func Test386OperandParser(t *testing.T)

func TestAMD64Encoder(t *testing.T)

func TestAMD64EndToEnd(t *testing.T)

func TestAMD64Errors(t *testing.T)

func TestAMD64OperandParser(t *testing.T)

func TestARM64EndToEnd(t *testing.T)

func TestARM64OperandParser(t *testing.T)

func TestARMEndToEnd(t *testing.T)

func TestARMOperandParser(t *testing.T)

func TestBadExpr(t *testing.T)

func TestErroneous(t *testing.T)

func TestExpr(t *testing.T)

func TestMIPS64EndToEnd(t *testing.T)

func TestMIPS64OperandParser(t *testing.T)

func TestPPC64EndToEnd(t *testing.T)

func TestPPC64OperandParser(t *testing.T)

func (*Parser) Parse() (*obj.Prog, bool)

