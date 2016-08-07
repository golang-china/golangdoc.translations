// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package printer implements printing of AST nodes.

// Package printer implements printing of AST nodes.
package printer

import (
    "bytes"
    "fmt"
    "go/ast"
    "go/token"
    "io"
    "os"
    "strconv"
    "strings"
    "text/tabwriter"
    "unicode"
    "unicode/utf8"
)

const (
    RawFormat Mode = 1 << iota // do not use a tabwriter; if set, UseSpaces is ignored
    TabIndent                  // use tabs for indentation independent of UseSpaces
    UseSpaces                  // use spaces instead of tabs for alignment
    SourcePos                  // emit //line comments to preserve original source positions
)

// A CommentedNode bundles an AST node and corresponding comments. It may be
// provided as argument to any of the Fprint functions.
type CommentedNode struct {
    Node     interface{} // *ast.File, or ast.Expr, ast.Decl, ast.Spec, or ast.Stmt
    Comments []*ast.CommentGroup
}

// A Config node controls the output of Fprint.
type Config struct {
    Mode     Mode // default: 0
    Tabwidth int  // default: 8
    Indent   int  // default: 0 (all code is indented at least by this much)
}

// A Mode value is a set of flags (or 0). They control printing.
type Mode uint

// Fprint "pretty-prints" an AST node to output. It calls Config.Fprint with
// default settings.
func Fprint(output io.Writer, fset *token.FileSet, node interface{}) error

// Fprint "pretty-prints" an AST node to output for a given configuration cfg.
// Position information is interpreted relative to the file set fset. The node
// type must be *ast.File, *CommentedNode, []ast.Decl, []ast.Stmt, or
// assignment-compatible to ast.Expr, ast.Decl, ast.Spec, or ast.Stmt.
func (*Config) Fprint(output io.Writer, fset *token.FileSet, node interface{}) error

