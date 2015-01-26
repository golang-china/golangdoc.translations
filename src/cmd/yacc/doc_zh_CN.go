// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Yacc is a version of yacc for Go. It is written in Go and generates parsers
// written in Go.
//
// Usage:
//
//	go tool yacc args...
//
// It is largely transliterated from the Inferno version written in Limbo which in
// turn was largely transliterated from the Plan 9 version written in C and
// documented at
//
//	http://plan9.bell-labs.com/magic/man2html/1/yacc
//
// Adepts of the original yacc will have no trouble adapting to this form of the
// tool.
//
// The directory $GOROOT/cmd/yacc/testdata/expr is a yacc program for a very simple
// expression parser. See expr.y and main.go in that directory for examples of how
// to write and build yacc programs.
//
// The generated parser is reentrant. Parse expects to be given an argument that
// conforms to the following interface:
//
//	type yyLexer interface {
//		Lex(lval *yySymType) int
//		Error(e string)
//	}
//
// Lex should return the token identifier, and place other token information in
// lval (which replaces the usual yylval). Error is equivalent to yyerror in the
// original yacc.
//
// Code inside the parser may refer to the variable yylex, which holds the yyLexer
// passed to Parse.
//
// Multiple grammars compiled into a single program should be placed in distinct
// packages. If that is impossible, the "-p prefix" flag to yacc sets the prefix,
// by default yy, that begins the names of symbols, including types, the parser,
// and the lexer, generated and referenced by yacc's generated code. Setting it to
// distinct values allows multiple grammars to be placed in a single package.

// Yacc is a version of yacc for Go. It is
// written in Go and generates parsers
// written in Go.
//
// Usage:
//
//	go tool yacc args...
//
// It is largely transliterated from the
// Inferno version written in Limbo which
// in turn was largely transliterated from
// the Plan 9 version written in C and
// documented at
//
//	http://plan9.bell-labs.com/magic/man2html/1/yacc
//
// Adepts of the original yacc will have no
// trouble adapting to this form of the
// tool.
//
// The directory
// $GOROOT/cmd/yacc/testdata/expr is a yacc
// program for a very simple expression
// parser. See expr.y and main.go in that
// directory for examples of how to write
// and build yacc programs.
//
// The generated parser is reentrant. Parse
// expects to be given an argument that
// conforms to the following interface:
//
//	type yyLexer interface {
//		Lex(lval *yySymType) int
//		Error(e string)
//	}
//
// Lex should return the token identifier,
// and place other token information in
// lval (which replaces the usual yylval).
// Error is equivalent to yyerror in the
// original yacc.
//
// Code inside the parser may refer to the
// variable yylex, which holds the yyLexer
// passed to Parse.
//
// Multiple grammars compiled into a single
// program should be placed in distinct
// packages. If that is impossible, the "-p
// prefix" flag to yacc sets the prefix, by
// default yy, that begins the names of
// symbols, including types, the parser,
// and the lexer, generated and referenced
// by yacc's generated code. Setting it to
// distinct values allows multiple grammars
// to be placed in a single package.
package main

// the following are adjustable according to memory size

// the following are adjustable according
// to memory size
const (
	ACTSIZE  = 30000
	NSTATES  = 2000
	TEMPSIZE = 2000

	SYMINC   = 50  // increase for non-term or term
	RULEINC  = 50  // increase for max rule length prodptr[i]
	PRODINC  = 100 // increase for productions     prodptr
	WSETINC  = 50  // increase for working sets    wsets
	STATEINC = 200 // increase for states          statemem

	NAMESIZE = 50
	NTYPES   = 63
	ISIZE    = 400

	PRIVATE = 0xE000 // unicode private use

	NTBASE     = 010000
	ERRCODE    = 8190
	ACCEPTCODE = 8191
	YYLEXUNK   = 3
	TOKSTART   = 4 //index of first defined token
)

// no, left, right, binary assoc.

// no, left, right, binary assoc.
const (
	NOASC = iota
	LASC
	RASC
	BASC
)

// flags for state generation

// flags for state generation
const (
	DONE = iota
	MUSTDO
	MUSTLOOKAHEAD
)

// flags for a rule having an action, and being reduced

// flags for a rule having an action, and
// being reduced
const (
	ACTFLAG = 1 << (iota + 2)
	REDFLAG
)

// parse tokens

// parse tokens
const (
	IDENTIFIER = PRIVATE + iota
	MARK
	TERM
	LEFT
	RIGHT
	BINARY
	PREC
	LCURLY
	IDENTCOLON
	NUMBER
	START
	TYPEDEF
	TYPENAME
	UNION
)

const EMPTY = 1

const ENDFILE = 0

const EOF = -1

const NOMORE = -1000

const OK = 1

const WHOKNOWS = 0

// macros for getting associativity and precedence levels

// macros for getting associativity and
// precedence levels
func ASSOC(i int) int

func PLEVEL(i int) int

// macros for setting associativity and precedence levels

// macros for setting associativity and
// precedence levels
func SETASC(i, j int) int

func SETPLEV(i, j int) int

func SETTYPE(i, j int) int

func TYPE(i int) int

type Item struct {
	// contains filtered or unexported fields
}

// structure declarations

// structure declarations
type Lkset []int

type Pitem struct {
	// contains filtered or unexported fields
}

type Resrv struct {
	// contains filtered or unexported fields
}

type Symb struct {
	// contains filtered or unexported fields
}

type Wset struct {
	// contains filtered or unexported fields
}
