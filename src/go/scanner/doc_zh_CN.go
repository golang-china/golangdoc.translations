// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package scanner implements a scanner for Go source text. It takes a []byte as
// source which can then be tokenized through repeated calls to the Scan method.
package scanner

// PrintError is a utility function that prints a list of errors to w, one error
// per line, if the err parameter is an ErrorList. Otherwise it prints the err
// string.
func PrintError(w io.Writer, err error)

// In an ErrorList, an error is represented by an *Error. The position Pos, if
// valid, points to the beginning of the offending token, and the error condition
// is described by Msg.
type Error struct {
	Pos token.Position
	Msg string
}

// Error implements the error interface.
func (e Error) Error() string

// An ErrorHandler may be provided to Scanner.Init. If a syntax error is
// encountered and a handler was installed, the handler is called with a position
// and an error message. The position points to the beginning of the offending
// token.
type ErrorHandler func(pos token.Position, msg string)

// ErrorList is a list of *Errors. The zero value for an ErrorList is an empty
// ErrorList ready to use.
type ErrorList []*Error

// Add adds an Error with given position and error message to an ErrorList.
func (p *ErrorList) Add(pos token.Position, msg string)

// Err returns an error equivalent to this error list. If the list is empty, Err
// returns nil.
func (p ErrorList) Err() error

// An ErrorList implements the error interface.
func (p ErrorList) Error() string

// ErrorList implements the sort Interface.
func (p ErrorList) Len() int

func (p ErrorList) Less(i, j int) bool

// RemoveMultiples sorts an ErrorList and removes all but the first error per line.
func (p *ErrorList) RemoveMultiples()

// Reset resets an ErrorList to no errors.
func (p *ErrorList) Reset()

// Sort sorts an ErrorList. *Error entries are sorted by position, other errors are
// sorted by error message, and before any *Error entry.
func (p ErrorList) Sort()

func (p ErrorList) Swap(i, j int)

// A mode value is a set of flags (or 0). They control scanner behavior.
type Mode uint

const (
	ScanComments Mode = 1 << iota // return comments as COMMENT tokens

)

// A Scanner holds the scanner's internal state while processing a given text. It
// can be allocated as part of another data structure but must be initialized via
// Init before use.
type Scanner struct {

	// public state - ok to modify
	ErrorCount int // number of errors encountered
	// contains filtered or unexported fields
}

// Init prepares the scanner s to tokenize the text src by setting the scanner at
// the beginning of src. The scanner uses the file set file for position
// information and it adds line information for each line. It is ok to re-use the
// same file when re-scanning the same file as line information which is already
// present is ignored. Init causes a panic if the file size does not match the src
// size.
//
// Calls to Scan will invoke the error handler err if they encounter a syntax error
// and err is not nil. Also, for each error encountered, the Scanner field
// ErrorCount is incremented by one. The mode parameter determines how comments are
// handled.
//
// Note that Init may call err if there is an error in the first character of the
// file.
func (s *Scanner) Init(file *token.File, src []byte, err ErrorHandler, mode Mode)

// Scan scans the next token and returns the token position, the token, and its
// literal string if applicable. The source end is indicated by token.EOF.
//
// If the returned token is a literal (token.IDENT, token.INT, token.FLOAT,
// token.IMAG, token.CHAR, token.STRING) or token.COMMENT, the literal string has
// the corresponding value.
//
// If the returned token is a keyword, the literal string is the keyword.
//
// If the returned token is token.SEMICOLON, the corresponding literal string is
// ";" if the semicolon was present in the source, and "\n" if the semicolon was
// inserted because of a newline or at EOF.
//
// If the returned token is token.ILLEGAL, the literal string is the offending
// character.
//
// In all other cases, Scan returns an empty literal string.
//
// For more tolerant parsing, Scan will return a valid token if possible even if a
// syntax error was encountered. Thus, even if the resulting token sequence
// contains no illegal tokens, a client may not assume that no error occurred.
// Instead it must check the scanner's ErrorCount or the number of calls of the
// error handler, if there was one installed.
//
// Scan adds line information to the file added to the file set with Init. Token
// positions are relative to that file and thus relative to the file set.
func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string)
