// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package asm implements the parser and instruction generator for the
// assembler. TODO: Split apart?
package asm

import (
	"bytes"
	"cmd/asm/internal/arch"
	"cmd/asm/internal/flags"
	"cmd/asm/internal/lex"
	"cmd/internal/obj"
	"cmd/internal/sys"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"text/scanner"
	"unicode/utf8"
)

// EOF represents the end of input.
var EOF = lex.Make(scanner.EOF, "EOF")

type Parser struct {
}

type Patch struct {
}

func NewParser(ctxt *obj.Link, ar *arch.Arch, lexer lex.TokenReader) *Parser

func (p *Parser) Parse() (*obj.Prog, bool)

