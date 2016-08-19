// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package ast declares the types used to represent syntax trees for Go
// packages.

// Package ast declares the types used to represent syntax trees for Go
// packages.
package ast

import (
    "bytes"
    "fmt"
    "go/scanner"
    "go/token"
    "io"
    "os"
    "reflect"
    "sort"
    "strconv"
    "strings"
    "unicode"
    "unicode/utf8"
)

// The list of possible Object kinds.
const (
	Bad ObjKind = iota // for error handling
	Pkg                // package
	Con                // constant
	Typ                // type
	Var                // variable
	Fun                // function or method
	Lbl                // label

)



const (
	// If set, duplicate function declarations are excluded.
	FilterFuncDuplicates MergeMode = 1 << iota
	// If set, comments that are not associated with a specific
	// AST node (as Doc or Comment) are excluded.
	FilterUnassociatedComments
	// If set, duplicate import declarations are excluded.
	FilterImportDuplicates
)



const (
	SEND ChanDir = 1 << iota
	RECV
)


// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.

// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.
type (
	// An ArrayType node represents an array or slice type.
	ArrayType struct {
		Lbrack token.Pos // position of "["
		Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
		Elt    Expr      // element type
	}

	// A StructType node represents a struct type.
	StructType struct {
		Struct     token.Pos  // position of "struct" keyword
		Fields     *FieldList // list of field declarations
		Incomplete bool       // true if (source) fields are missing in the Fields list
	}

	// A FuncType node represents a function type.
	FuncType struct {
		Func    token.Pos  // position of "func" keyword (token.NoPos if there is no "func")
		Params  *FieldList // (incoming) parameters; non-nil
		Results *FieldList // (outgoing) results; or nil
	}

	// An InterfaceType node represents an interface type.
	InterfaceType struct {
		Interface  token.Pos  // position of "interface" keyword
		Methods    *FieldList // list of methods
		Incomplete bool       // true if (source) methods are missing in the Methods list
	}

	// A MapType node represents a map type.
	MapType struct {
		Map   token.Pos // position of "map" keyword
		Key   Expr
		Value Expr
	}

	// A ChanType node represents a channel type.
	ChanType struct {
		Begin token.Pos // position of "chan" keyword or "<-" (whichever comes first)
		Arrow token.Pos // position of "<-" (token.NoPos if there is no "<-")
		Dir   ChanDir   // channel direction
		Value Expr      // value type
	}
)


// A declaration is represented by one of the following declaration nodes.

// A declaration is represented by one of the following declaration nodes.
type (
	// A BadDecl node is a placeholder for declarations containing
	// syntax errors for which no correct declaration nodes can be
	// created.

	// A BadDecl node is a placeholder for declarations containing
	// syntax errors for which no correct declaration nodes can be
	// created.
	BadDecl struct {
		From, To token.Pos // position range of bad declaration
	}

	// A GenDecl node (generic declaration node) represents an import,
	// constant, type or variable declaration. A valid Lparen position
	// (Lparen.Line > 0) indicates a parenthesized declaration.
	//
	// Relationship between Tok value and Specs element type:
	//
	//     token.IMPORT  *ImportSpec
	//     token.CONST   *ValueSpec
	//     token.TYPE    *TypeSpec
	//     token.VAR     *ValueSpec

	// A GenDecl node (generic declaration node) represents an import,
	// constant, type or variable declaration. A valid Lparen position
	// (Lparen.Line > 0) indicates a parenthesized declaration.
	//
	// Relationship between Tok value and Specs element type:
	//
	//     token.IMPORT  *ImportSpec
	//     token.CONST   *ValueSpec
	//     token.TYPE    *TypeSpec
	//     token.VAR     *ValueSpec
	GenDecl struct {
		Doc    *CommentGroup // associated documentation; or nil
		TokPos token.Pos     // position of Tok
		Tok    token.Token   // IMPORT, CONST, TYPE, VAR
		Lparen token.Pos     // position of '(', if any
		Specs  []Spec
		Rparen token.Pos // position of ')', if any
	}

	// A FuncDecl node represents a function declaration.
	FuncDecl struct {
		Doc  *CommentGroup // associated documentation; or nil
		Recv *FieldList    // receiver (methods); or nil (functions)
		Name *Ident        // function/method name
		Type *FuncType     // function signature: parameters, results, and position of "func" keyword
		Body *BlockStmt    // function body; or nil (forward declaration)
	}
)


// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.
type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.

	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	BadExpr struct {
		From, To token.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
		Obj     *Object   // denoted object; or nil
	}

	// An Ellipsis node stands for the "..." type in a
	// parameter list or the "..." length in an array type.

	// An Ellipsis node stands for the "..." type in a
	// parameter list or the "..." length in an array type.
	Ellipsis struct {
		Ellipsis token.Pos // position of "..."
		Elt      Expr      // ellipsis element type (parameter lists only); or nil
	}

	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		ValuePos token.Pos   // literal position
		Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
		Value    string      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
	}

	// A FuncLit node represents a function literal.
	FuncLit struct {
		Type *FuncType  // function type
		Body *BlockStmt // function body
	}

	// A CompositeLit node represents a composite literal.
	CompositeLit struct {
		Type   Expr      // literal type; or nil
		Lbrace token.Pos // position of "{"
		Elts   []Expr    // list of composite elements; or nil
		Rbrace token.Pos // position of "}"
	}

	// A ParenExpr node represents a parenthesized expression.
	ParenExpr struct {
		Lparen token.Pos // position of "("
		X      Expr      // parenthesized expression
		Rparen token.Pos // position of ")"
	}

	// A SelectorExpr node represents an expression followed by a selector.
	SelectorExpr struct {
		X   Expr   // expression
		Sel *Ident // field selector
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		X      Expr      // expression
		Lbrack token.Pos // position of "["
		Index  Expr      // index expression
		Rbrack token.Pos // position of "]"
	}

	// An SliceExpr node represents an expression followed by slice indices.
	SliceExpr struct {
		X      Expr      // expression
		Lbrack token.Pos // position of "["
		Low    Expr      // begin of slice range; or nil
		High   Expr      // end of slice range; or nil
		Max    Expr      // maximum capacity of slice; or nil
		Slice3 bool      // true if 3-index slice (2 colons present)
		Rbrack token.Pos // position of "]"
	}

	// A TypeAssertExpr node represents an expression followed by a
	// type assertion.

	// A TypeAssertExpr node represents an expression followed by a
	// type assertion.
	TypeAssertExpr struct {
		X      Expr      // expression
		Lparen token.Pos // position of "("
		Type   Expr      // asserted type; nil means type switch X.(type)
		Rparen token.Pos // position of ")"
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Fun      Expr      // function expression
		Lparen   token.Pos // position of "("
		Args     []Expr    // function arguments; or nil
		Ellipsis token.Pos // position of "...", if any
		Rparen   token.Pos // position of ")"
	}

	// A StarExpr node represents an expression of the form "*" Expression.
	// Semantically it could be a unary "*" expression, or a pointer type.

	// A StarExpr node represents an expression of the form "*" Expression.
	// Semantically it could be a unary "*" expression, or a pointer type.
	StarExpr struct {
		Star token.Pos // position of "*"
		X    Expr      // operand
	}

	// A UnaryExpr node represents a unary expression.
	// Unary "*" expressions are represented via StarExpr nodes.

	// A UnaryExpr node represents a unary expression.
	// Unary "*" expressions are represented via StarExpr nodes.
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr        // operand
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.
	KeyValueExpr struct {
		Key   Expr
		Colon token.Pos // position of ":"
		Value Expr
	}
)


// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.

	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	BadStmt struct {
		From, To token.Pos // position range of bad statement
	}

	// A DeclStmt node represents a declaration in a statement list.
	DeclStmt struct {
		Decl Decl // *GenDecl with CONST, TYPE, or VAR token
	}

	// An EmptyStmt node represents an empty statement.
	// The "position" of the empty statement is the position
	// of the immediately following (explicit or implicit) semicolon.

	// An EmptyStmt node represents an empty statement.
	// The "position" of the empty statement is the position
	// of the immediately following (explicit or implicit) semicolon.
	EmptyStmt struct {
		Semicolon token.Pos // position of following ";"
		Implicit  bool      // if set, ";" was omitted in the source
	}

	// A LabeledStmt node represents a labeled statement.
	LabeledStmt struct {
		Label *Ident
		Colon token.Pos // position of ":"
		Stmt  Stmt
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	ExprStmt struct {
		X Expr // expression
	}

	// A SendStmt node represents a send statement.
	SendStmt struct {
		Chan  Expr
		Arrow token.Pos // position of "<-"
		Value Expr
	}

	// An IncDecStmt node represents an increment or decrement statement.
	IncDecStmt struct {
		X      Expr
		TokPos token.Pos   // position of Tok
		Tok    token.Token // INC or DEC
	}

	// An AssignStmt node represents an assignment or
	// a short variable declaration.

	// An AssignStmt node represents an assignment or
	// a short variable declaration.
	AssignStmt struct {
		Lhs    []Expr
		TokPos token.Pos   // position of Tok
		Tok    token.Token // assignment token, DEFINE
		Rhs    []Expr
	}

	// A GoStmt node represents a go statement.
	GoStmt struct {
		Go   token.Pos // position of "go" keyword
		Call *CallExpr
	}

	// A DeferStmt node represents a defer statement.
	DeferStmt struct {
		Defer token.Pos // position of "defer" keyword
		Call  *CallExpr
	}

	// A ReturnStmt node represents a return statement.
	ReturnStmt struct {
		Return  token.Pos // position of "return" keyword
		Results []Expr    // result expressions; or nil
	}

	// A BranchStmt node represents a break, continue, goto,
	// or fallthrough statement.

	// A BranchStmt node represents a break, continue, goto,
	// or fallthrough statement.
	BranchStmt struct {
		TokPos token.Pos   // position of Tok
		Tok    token.Token // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
		Label  *Ident      // label name; or nil
	}

	// A BlockStmt node represents a braced statement list.
	BlockStmt struct {
		Lbrace token.Pos // position of "{"
		List   []Stmt
		Rbrace token.Pos // position of "}"
	}

	// An IfStmt node represents an if statement.
	IfStmt struct {
		If   token.Pos // position of "if" keyword
		Init Stmt      // initialization statement; or nil
		Cond Expr      // condition
		Body *BlockStmt
		Else Stmt // else branch; or nil
	}

	// A CaseClause represents a case of an expression or type switch statement.
	CaseClause struct {
		Case  token.Pos // position of "case" or "default" keyword
		List  []Expr    // list of expressions or types; nil means default case
		Colon token.Pos // position of ":"
		Body  []Stmt    // statement list; or nil
	}

	// A SwitchStmt node represents an expression switch statement.
	SwitchStmt struct {
		Switch token.Pos  // position of "switch" keyword
		Init   Stmt       // initialization statement; or nil
		Tag    Expr       // tag expression; or nil
		Body   *BlockStmt // CaseClauses only
	}

	// An TypeSwitchStmt node represents a type switch statement.
	TypeSwitchStmt struct {
		Switch token.Pos  // position of "switch" keyword
		Init   Stmt       // initialization statement; or nil
		Assign Stmt       // x := y.(type) or y.(type)
		Body   *BlockStmt // CaseClauses only
	}

	// A CommClause node represents a case of a select statement.
	CommClause struct {
		Case  token.Pos // position of "case" or "default" keyword
		Comm  Stmt      // send or receive statement; nil means default case
		Colon token.Pos // position of ":"
		Body  []Stmt    // statement list; or nil
	}

	// An SelectStmt node represents a select statement.
	SelectStmt struct {
		Select token.Pos  // position of "select" keyword
		Body   *BlockStmt // CommClauses only
	}

	// A ForStmt represents a for statement.
	ForStmt struct {
		For  token.Pos // position of "for" keyword
		Init Stmt      // initialization statement; or nil
		Cond Expr      // condition; or nil
		Post Stmt      // post iteration statement; or nil
		Body *BlockStmt
	}

	// A RangeStmt represents a for statement with a range clause.
	RangeStmt struct {
		For        token.Pos   // position of "for" keyword
		Key, Value Expr        // Key, Value may be nil
		TokPos     token.Pos   // position of Tok; invalid if Key == nil
		Tok        token.Token // ILLEGAL if Key == nil, ASSIGN, DEFINE
		X          Expr        // value to range over
		Body       *BlockStmt
	}
)


// The direction of a channel type is indicated by one
// of the following constants.

// The direction of a channel type is indicated by one
// of the following constants.
type ChanDir int


// A Comment node represents a single //-style or /*-style comment.
type Comment struct {
	Slash token.Pos // position of "/" starting the comment
	Text  string    // comment text (excluding '\n' for //-style comments)
}


// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
type CommentGroup struct {
	List []*Comment // len(List) > 0
}


// A CommentMap maps an AST node to a list of comment groups
// associated with it. See NewCommentMap for a description of
// the association.

// A CommentMap maps an AST node to a list of comment groups
// associated with it. See NewCommentMap for a description of
// the association.
type CommentMap map[Node][]*CommentGroup


// All declaration nodes implement the Decl interface.
type Decl interface {
	Node
	declNode()
}


// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}


// A Field represents a Field declaration list in a struct type,
// a method list in an interface type, or a parameter/result declaration
// in a signature.

// A Field represents a Field declaration list in a struct type,
// a method list in an interface type, or a parameter/result declaration
// in a signature.
type Field struct {
	Doc     *CommentGroup // associated documentation; or nil
	Names   []*Ident      // field/method/parameter names; or nil if anonymous field
	Type    Expr          // field/method/parameter type
	Tag     *BasicLit     // field tag; or nil
	Comment *CommentGroup // line comments; or nil
}


// A FieldFilter may be provided to Fprint to control the output.
type FieldFilter func(name string, value reflect.Value) bool


// A FieldList represents a list of Fields, enclosed by parentheses or braces.
type FieldList struct {
	Opening token.Pos // position of opening parenthesis/brace, if any
	List    []*Field  // field list; or nil
	Closing token.Pos // position of closing parenthesis/brace, if any
}


// A File node represents a Go source file.
//
// The Comments list contains all comments in the source file in order of
// appearance, including the comments that are pointed to from other nodes
// via Doc and Comment fields.

// A File node represents a Go source file.
//
// The Comments list contains all comments in the source file in order of
// appearance, including the comments that are pointed to from other nodes
// via Doc and Comment fields.
type File struct {
	Doc        *CommentGroup   // associated documentation; or nil
	Package    token.Pos       // position of "package" keyword
	Name       *Ident          // package name
	Decls      []Decl          // top-level declarations; or nil
	Scope      *Scope          // package scope (this file only)
	Imports    []*ImportSpec   // imports in this file
	Unresolved []*Ident        // unresolved identifiers in this file
	Comments   []*CommentGroup // list of all comments in the source file
}



type Filter func(string) bool


// An Importer resolves import paths to package Objects.
// The imports map records the packages already imported,
// indexed by package id (canonical import path).
// An Importer must determine the canonical import path and
// check the map to see if it is already present in the imports map.
// If so, the Importer can return the map entry.  Otherwise, the
// Importer should load the package data for the given path into
// a new *Object (pkg), record pkg in the imports map, and then
// return pkg.

// An Importer resolves import paths to package Objects.
// The imports map records the packages already imported,
// indexed by package id (canonical import path).
// An Importer must determine the canonical import path and
// check the map to see if it is already present in the imports map.
// If so, the Importer can return the map entry. Otherwise, the
// Importer should load the package data for the given path into
// a new *Object (pkg), record pkg in the imports map, and then
// return pkg.
type Importer func(imports map[string]*Object, path string) (pkg *Object, err error)


// The MergeMode flags control the behavior of MergePackageFiles.
type MergeMode uint


// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}


// ObjKind describes what an object represents.
type ObjKind int


// An Object describes a named language entity such as a package, constant,
// type, variable, function (incl. methods), or label.
//
// The Data fields contains object-specific data:
//
//     Kind    Data type         Data value
//     Pkg	*types.Package    package scope
//     Con     int               iota for the respective declaration
//     Con     != nil            constant value
//     Typ     *Scope            (used as method scope during type checking - transient)

// An Object describes a named language entity such as a package, constant,
// type, variable, function (incl. methods), or label.
//
// The Data fields contains object-specific data:
//
//     Kind    Data type         Data value
//     Pkg	*types.Package    package scope
//     Con     int               iota for the respective declaration
//     Con     != nil            constant value
//     Typ     *Scope            (used as method scope during type checking - transient)
type Object struct {
	Kind ObjKind
	Name string      // declared name
	Decl interface{} // corresponding Field, XxxSpec, FuncDecl, LabeledStmt, AssignStmt, Scope; or nil
	Data interface{} // object-specific data; or nil
	Type interface{} // placeholder for type information; may be nil
}


// A Package node represents a set of source files
// collectively building a Go package.

// A Package node represents a set of source files
// collectively building a Go package.
type Package struct {
	Name    string             // package name
	Scope   *Scope             // package scope across all files
	Imports map[string]*Object // map of package id -> package object
	Files   map[string]*File   // Go source files by filename
}


// A Scope maintains the set of named language entities declared
// in the scope and a link to the immediately surrounding (outer)
// scope.

// A Scope maintains the set of named language entities declared
// in the scope and a link to the immediately surrounding (outer)
// scope.
type Scope struct {
	Outer   *Scope
	Objects map[string]*Object
}


// A Spec node represents a single (non-parenthesized) import,
// constant, type, or variable declaration.

// A Spec node represents a single (non-parenthesized) import,
// constant, type, or variable declaration.
type (
	// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
	Spec interface {
		Node
		specNode()
	}

	// An ImportSpec node represents a single package import.
	ImportSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Name    *Ident        // local package name (including "."); or nil
		Path    *BasicLit     // import path
		Comment *CommentGroup // line comments; or nil
		EndPos  token.Pos     // end of spec (overrides Path.Pos if nonzero)
	}

	// A ValueSpec node represents a constant or variable declaration
	// (ConstSpec or VarSpec production).

	// A ValueSpec node represents a constant or variable declaration
	// (ConstSpec or VarSpec production).
	ValueSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Names   []*Ident      // value names (len(Names) > 0)
		Type    Expr          // value type; or nil
		Values  []Expr        // initial values; or nil
		Comment *CommentGroup // line comments; or nil
	}

	// A TypeSpec node represents a type declaration (TypeSpec production).
	TypeSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Name    *Ident        // type name
		Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
		Comment *CommentGroup // line comments; or nil
	}
)


// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}


// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}


// FileExports trims the AST for a Go source file in place such that
// only exported nodes remain: all top-level identifiers which are not exported
// and their associated information (such as type, initial value, or function
// body) are removed. Non-exported fields and methods of exported types are
// stripped. The File.Comments list is not changed.
//
// FileExports reports whether there are exported declarations.
func FileExports(src *File) bool

// FilterDecl trims the AST for a Go declaration in place by removing
// all names (including struct field and interface method names, but
// not from parameter lists) that don't pass through the filter f.
//
// FilterDecl reports whether there are any declared names left after
// filtering.
func FilterDecl(decl Decl, f Filter) bool

// FilterFile trims the AST for a Go file in place by removing all
// names from top-level declarations (including struct field and
// interface method names, but not from parameter lists) that don't
// pass through the filter f. If the declaration is empty afterwards,
// the declaration is removed from the AST. Import declarations are
// always removed. The File.Comments list is not changed.
//
// FilterFile reports whether there are any top-level declarations
// left after filtering.
func FilterFile(src *File, f Filter) bool

// FilterPackage trims the AST for a Go package in place by removing
// all names from top-level declarations (including struct field and
// interface method names, but not from parameter lists) that don't
// pass through the filter f. If the declaration is empty afterwards,
// the declaration is removed from the AST. The pkg.Files list is not
// changed, so that file names and top-level package comments don't get
// lost.
//
// FilterPackage reports whether there are any top-level declarations
// left after filtering.
func FilterPackage(pkg *Package, f Filter) bool

// Fprint prints the (sub-)tree starting at AST node x to w.
// If fset != nil, position information is interpreted relative
// to that file set. Otherwise positions are printed as integer
// values (file set specific offsets).
//
// A non-nil FieldFilter f may be provided to control the output:
// struct fields for which f(fieldname, fieldvalue) is true are
// printed; all others are filtered from the output. Unexported
// struct fields are never printed.
func Fprint(w io.Writer, fset *token.FileSet, x interface{}, f FieldFilter) error

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// recursively for each of the non-nil children of node, followed by a
// call of f(nil).
func Inspect(node Node, f func(Node) bool)

// IsExported reports whether name is an exported Go symbol
// (that is, whether it begins with an upper-case letter).
func IsExported(name string) bool

// MergePackageFiles creates a file AST by merging the ASTs of the
// files belonging to a package. The mode flags control merging behavior.
func MergePackageFiles(pkg *Package, mode MergeMode) *File

// NewCommentMap creates a new comment map by associating comment groups
// of the comments list with the nodes of the AST specified by node.
//
// A comment group g is associated with a node n if:
//
//   - g starts on the same line as n ends
//   - g starts on the line immediately following n, and there is
//     at least one empty line after g and before the next node
//   - g starts before n and is not associated to the node before n
//     via the previous rules
//
// NewCommentMap tries to associate a comment group to the "largest"
// node possible: For instance, if the comment is a line comment
// trailing an assignment, the comment is associated with the entire
// assignment rather than just the last operand in the assignment.
func NewCommentMap(fset *token.FileSet, node Node, comments []*CommentGroup) CommentMap

// NewIdent creates a new Ident without position.
// Useful for ASTs generated by code other than the Go parser.
func NewIdent(name string) *Ident

// NewObj creates a new object of a given kind and name.
func NewObj(kind ObjKind, name string) *Object

// NewPackage creates a new Package node from a set of File nodes. It resolves
// unresolved identifiers across files and updates each file's Unresolved list
// accordingly. If a non-nil importer and universe scope are provided, they are
// used to resolve identifiers not declared in any of the package files. Any
// remaining unresolved identifiers are reported as undeclared. If the files
// belong to different packages, one package name is selected and files with
// different package names are reported and then ignored.
// The result is a package node and a scanner.ErrorList if there were errors.
func NewPackage(fset *token.FileSet, files map[string]*File, importer Importer, universe *Scope) (*Package, error)

// NewScope creates a new scope nested in the outer scope.
func NewScope(outer *Scope) *Scope

// NotNilFilter returns true for field values that are not nil;
// it returns false otherwise.
func NotNilFilter(_ string, v reflect.Value) bool

// PackageExports trims the AST for a Go package in place such that
// only exported nodes remain. The pkg.Files list is not changed, so that
// file names and top-level package comments don't get lost.
//
// PackageExports reports whether there are exported declarations;
// it returns false otherwise.
func PackageExports(pkg *Package) bool

// Print prints x to standard output, skipping nil fields.
// Print(fset, x) is the same as Fprint(os.Stdout, fset, x, NotNilFilter).
func Print(fset *token.FileSet, x interface{}) error

// SortImports sorts runs of consecutive import lines in import blocks in f. It
// also removes duplicate imports when it is possible to do so without data
// loss.
func SortImports(fset *token.FileSet, f *File)

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
func Walk(v Visitor, node Node)

func (*ArrayType) End() token.Pos

func (*ArrayType) Pos() token.Pos

func (*AssignStmt) End() token.Pos

func (*AssignStmt) Pos() token.Pos

func (*BadDecl) End() token.Pos

// Pos and End implementations for declaration nodes.
func (*BadDecl) Pos() token.Pos

func (*BadExpr) End() token.Pos

// Pos and End implementations for expression/type nodes.
func (*BadExpr) Pos() token.Pos

func (*BadStmt) End() token.Pos

// Pos and End implementations for statement nodes.
func (*BadStmt) Pos() token.Pos

func (*BasicLit) End() token.Pos

func (*BasicLit) Pos() token.Pos

func (*BinaryExpr) End() token.Pos

func (*BinaryExpr) Pos() token.Pos

func (*BlockStmt) End() token.Pos

func (*BlockStmt) Pos() token.Pos

func (*BranchStmt) End() token.Pos

func (*BranchStmt) Pos() token.Pos

func (*CallExpr) End() token.Pos

func (*CallExpr) Pos() token.Pos

func (*CaseClause) End() token.Pos

func (*CaseClause) Pos() token.Pos

func (*ChanType) End() token.Pos

func (*ChanType) Pos() token.Pos

func (*CommClause) End() token.Pos

func (*CommClause) Pos() token.Pos

func (*Comment) End() token.Pos

func (*Comment) Pos() token.Pos

func (*CommentGroup) End() token.Pos

func (*CommentGroup) Pos() token.Pos

// Text returns the text of the comment.
// Comment markers (//, /*, and */), the first space of a line comment, and
// leading and trailing empty lines are removed. Multiple empty lines are
// reduced to one, and trailing space on lines is trimmed. Unless the result
// is empty, it is newline-terminated.
func (*CommentGroup) Text() string

func (*CompositeLit) End() token.Pos

func (*CompositeLit) Pos() token.Pos

func (*DeclStmt) End() token.Pos

func (*DeclStmt) Pos() token.Pos

func (*DeferStmt) End() token.Pos

func (*DeferStmt) Pos() token.Pos

func (*Ellipsis) End() token.Pos

func (*Ellipsis) Pos() token.Pos

func (*EmptyStmt) End() token.Pos

func (*EmptyStmt) Pos() token.Pos

func (*ExprStmt) End() token.Pos

func (*ExprStmt) Pos() token.Pos

func (*Field) End() token.Pos

func (*Field) Pos() token.Pos

func (*FieldList) End() token.Pos

// NumFields returns the number of (named and anonymous fields) in a FieldList.
func (*FieldList) NumFields() int

func (*FieldList) Pos() token.Pos

func (*File) End() token.Pos

func (*File) Pos() token.Pos

func (*ForStmt) End() token.Pos

func (*ForStmt) Pos() token.Pos

func (*FuncDecl) End() token.Pos

func (*FuncDecl) Pos() token.Pos

func (*FuncLit) End() token.Pos

func (*FuncLit) Pos() token.Pos

func (*FuncType) End() token.Pos

func (*FuncType) Pos() token.Pos

func (*GenDecl) End() token.Pos

func (*GenDecl) Pos() token.Pos

func (*GoStmt) End() token.Pos

func (*GoStmt) Pos() token.Pos

func (*Ident) End() token.Pos

// IsExported reports whether id is an exported Go symbol
// (that is, whether it begins with an uppercase letter).
func (*Ident) IsExported() bool

func (*Ident) Pos() token.Pos

func (*Ident) String() string

func (*IfStmt) End() token.Pos

func (*IfStmt) Pos() token.Pos

func (*ImportSpec) End() token.Pos

// Pos and End implementations for spec nodes.
func (*ImportSpec) Pos() token.Pos

func (*IncDecStmt) End() token.Pos

func (*IncDecStmt) Pos() token.Pos

func (*IndexExpr) End() token.Pos

func (*IndexExpr) Pos() token.Pos

func (*InterfaceType) End() token.Pos

func (*InterfaceType) Pos() token.Pos

func (*KeyValueExpr) End() token.Pos

func (*KeyValueExpr) Pos() token.Pos

func (*LabeledStmt) End() token.Pos

func (*LabeledStmt) Pos() token.Pos

func (*MapType) End() token.Pos

func (*MapType) Pos() token.Pos

// Pos computes the source position of the declaration of an object name.
// The result may be an invalid position if it cannot be computed
// (obj.Decl may be nil or not correct).
func (*Object) Pos() token.Pos

func (*Package) End() token.Pos

func (*Package) Pos() token.Pos

func (*ParenExpr) End() token.Pos

func (*ParenExpr) Pos() token.Pos

func (*RangeStmt) End() token.Pos

func (*RangeStmt) Pos() token.Pos

func (*ReturnStmt) End() token.Pos

func (*ReturnStmt) Pos() token.Pos

// Insert attempts to insert a named object obj into the scope s.
// If the scope already contains an object alt with the same name,
// Insert leaves the scope unchanged and returns alt. Otherwise
// it inserts obj and returns nil.
func (*Scope) Insert(obj *Object) (alt *Object)

// Lookup returns the object with the given name if it is
// found in scope s, otherwise it returns nil. Outer scopes
// are ignored.
func (*Scope) Lookup(name string) *Object

// Debugging support
func (*Scope) String() string

func (*SelectStmt) End() token.Pos

func (*SelectStmt) Pos() token.Pos

func (*SelectorExpr) End() token.Pos

func (*SelectorExpr) Pos() token.Pos

func (*SendStmt) End() token.Pos

func (*SendStmt) Pos() token.Pos

func (*SliceExpr) End() token.Pos

func (*SliceExpr) Pos() token.Pos

func (*StarExpr) End() token.Pos

func (*StarExpr) Pos() token.Pos

func (*StructType) End() token.Pos

func (*StructType) Pos() token.Pos

func (*SwitchStmt) End() token.Pos

func (*SwitchStmt) Pos() token.Pos

func (*TypeAssertExpr) End() token.Pos

func (*TypeAssertExpr) Pos() token.Pos

func (*TypeSpec) End() token.Pos

func (*TypeSpec) Pos() token.Pos

func (*TypeSwitchStmt) End() token.Pos

func (*TypeSwitchStmt) Pos() token.Pos

func (*UnaryExpr) End() token.Pos

func (*UnaryExpr) Pos() token.Pos

func (*ValueSpec) End() token.Pos

func (*ValueSpec) Pos() token.Pos

// Comments returns the list of comment groups in the comment map.
// The result is sorted is source order.

// Comments returns the list of comment groups in the comment map.
// The result is sorted in source order.
func (CommentMap) Comments() []*CommentGroup

// Filter returns a new comment map consisting of only those
// entries of cmap for which a corresponding node exists in
// the AST specified by node.
func (CommentMap) Filter(node Node) CommentMap

func (CommentMap) String() string

// Update replaces an old node in the comment map with the new node
// and returns the new node. Comments that were associated with the
// old node are associated with the new node.
func (CommentMap) Update(old, new Node) Node

func (ObjKind) String() string

