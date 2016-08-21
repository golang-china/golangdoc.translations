// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package types declares the data types and implements
// the algorithms for type-checking of Go packages. Use
// Config.Check to invoke the type checker for a package.
// Alternatively, create a new type checked with NewChecker
// and invoke it incrementally by calling Checker.Files.
//
// Type-checking consists of several interdependent phases:
//
// Name resolution maps each identifier (ast.Ident) in the program to the
// language object (Object) it denotes.
// Use Info.{Defs,Uses,Implicits} for the results of name resolution.
//
// Constant folding computes the exact constant value (constant.Value)
// for every expression (ast.Expr) that is a compile-time constant.
// Use Info.Types[expr].Value for the results of constant folding.
//
// Type inference computes the type (Type) of every expression (ast.Expr)
// and checks for compliance with the language specification.
// Use Info.Types[expr].Type for the results of type inference.
//
// For a tutorial, see https://golang.org/s/types-tutorial.

// Package types declares the data types and implements
// the algorithms for type-checking of Go packages. Use
// Config.Check to invoke the type checker for a package.
// Alternatively, create a new type checked with NewChecker
// and invoke it incrementally by calling Checker.Files.
//
// Type-checking consists of several interdependent phases:
//
// Name resolution maps each identifier (ast.Ident) in the program to the
// language object (Object) it denotes.
// Use Info.{Defs,Uses,Implicits} for the results of name resolution.
//
// Constant folding computes the exact constant value (constant.Value)
// for every expression (ast.Expr) that is a compile-time constant.
// Use Info.Types[expr].Value for the results of constant folding.
//
// Type inference computes the type (Type) of every expression (ast.Expr)
// and checks for compliance with the language specification.
// Use Info.Types[expr].Type for the results of type inference.
//
// For a tutorial, see https://golang.org/s/types-tutorial.
package types // import "go/types"

import (
    "bytes"
    "container/heap"
    "fmt"
    "go/ast"
    "go/constant"
    "go/parser"
    "go/token"
    "io"
    "math"
    "sort"
    "strconv"
    "strings"
    "sync"
    "testing"
    "unicode"
)

const (
    FieldVal   SelectionKind = iota // x.f is a struct field selector
    MethodVal                       // x.f is a method selector
    MethodExpr                      // x.f is a method expression
)

const (
    Invalid BasicKind = iota // type is invalid

    // predeclared types
    Bool
    Int
    Int8
    Int16
    Int32
    Int64
    Uint
    Uint8
    Uint16
    Uint32
    Uint64
    Uintptr
    Float32
    Float64
    Complex64
    Complex128
    String
    UnsafePointer

    // types for untyped values
    UntypedBool
    UntypedInt
    UntypedRune
    UntypedFloat
    UntypedComplex
    UntypedString
    UntypedNil

    // aliases
    Byte = Uint8
    Rune = Int32
)

// Properties of basic types.
const (
    IsBoolean BasicInfo = 1 << iota
    IsInteger
    IsUnsigned
    IsFloat
    IsComplex
    IsString
    IsUntyped

    IsOrdered   = IsInteger | IsFloat | IsString
    IsNumeric   = IsInteger | IsFloat | IsComplex
    IsConstType = IsBoolean | IsNumeric | IsString
)

// The direction of a channel is indicated by one of these constants.
const (
    SendRecv ChanDir = iota
    SendOnly
    RecvOnly
)

var Typ = []*Basic{
    Invalid: {Invalid, 0, "invalid type"},

    Bool:          {Bool, IsBoolean, "bool"},
    Int:           {Int, IsInteger, "int"},
    Int8:          {Int8, IsInteger, "int8"},
    Int16:         {Int16, IsInteger, "int16"},
    Int32:         {Int32, IsInteger, "int32"},
    Int64:         {Int64, IsInteger, "int64"},
    Uint:          {Uint, IsInteger | IsUnsigned, "uint"},
    Uint8:         {Uint8, IsInteger | IsUnsigned, "uint8"},
    Uint16:        {Uint16, IsInteger | IsUnsigned, "uint16"},
    Uint32:        {Uint32, IsInteger | IsUnsigned, "uint32"},
    Uint64:        {Uint64, IsInteger | IsUnsigned, "uint64"},
    Uintptr:       {Uintptr, IsInteger | IsUnsigned, "uintptr"},
    Float32:       {Float32, IsFloat, "float32"},
    Float64:       {Float64, IsFloat, "float64"},
    Complex64:     {Complex64, IsComplex, "complex64"},
    Complex128:    {Complex128, IsComplex, "complex128"},
    String:        {String, IsString, "string"},
    UnsafePointer: {UnsafePointer, 0, "Pointer"},

    UntypedBool:    {UntypedBool, IsBoolean | IsUntyped, "untyped bool"},
    UntypedInt:     {UntypedInt, IsInteger | IsUntyped, "untyped int"},
    UntypedRune:    {UntypedRune, IsInteger | IsUntyped, "untyped rune"},
    UntypedFloat:   {UntypedFloat, IsFloat | IsUntyped, "untyped float"},
    UntypedComplex: {UntypedComplex, IsComplex | IsUntyped, "untyped complex"},
    UntypedString:  {UntypedString, IsString | IsUntyped, "untyped string"},
    UntypedNil:     {UntypedNil, IsUntyped, "untyped nil"},
}

var (
    Universe *Scope
    Unsafe   *Package
)

// An Array represents an array type.
type Array struct {
    len  int64
    elem Type
}

// A Basic represents a basic type.
type Basic struct {
    kind BasicKind
    info BasicInfo
    name string
}

// BasicInfo is a set of flags describing properties of a basic type.
type BasicInfo int

// BasicKind describes the kind of basic type.
type BasicKind int

// A Builtin represents a built-in function.
// Builtins don't have a valid type.
type Builtin struct {
    id builtinId
}

// A Chan represents a channel type.
type Chan struct {
    dir  ChanDir
    elem Type
}

// A ChanDir value indicates a channel direction.
type ChanDir int

// A Checker maintains the state of the type checker.
// It must be created with NewChecker.
type Checker struct {
    // package information
    // (initialized by NewChecker, valid for the life-time of checker)
    conf *Config
    fset *token.FileSet
    pkg  *Package

    objMap map[Object]*declInfo // maps package-level object to declaration info

    // information collected during type-checking of a set of package files
    // (initialized by Files, valid only for the duration of check.Files;
    // maps and lists are allocated on demand)
    files            []*ast.File                       // package files
    unusedDotImports map[*Scope]map[*Package]token.Pos // positions of unused dot-imported packages for each file scope

    firstErr error                 // first error encountered
    methods  map[string][]*Func    // maps type names to associated methods
    untyped  map[ast.Expr]exprInfo // map of expressions without final type
    funcs    []funcInfo            // list of functions to type-check
    delayed  []func()              // delayed checks requiring fully setup types

    pos token.Pos // if valid, identifiers are looked up as if at position pos (used by Eval)

    // debugging
    indent int // indentation for tracing
}

// A Config specifies the configuration for type checking.
// The zero value for Config is a ready-to-use default configuration.
type Config struct {
    // If IgnoreFuncBodies is set, function bodies are not
    // type-checked.
    IgnoreFuncBodies bool

    // If FakeImportC is set, `import "C"` (for packages requiring Cgo)
    // declares an empty "C" package and errors are omitted for qualified
    // identifiers referring to package C (which won't find an object).
    // This feature is intended for the standard library cmd/api tool.
    //
    // Caution: Effects may be unpredictable due to follow-up errors.
    //          Do not use casually!
    FakeImportC bool

    // If Error != nil, it is called with each error found
    // during type checking; err has dynamic type Error.
    // Secondary errors (for instance, to enumerate all types
    // involved in an invalid recursive type declaration) have
    // error strings that start with a '\t' character.
    // If Error == nil, type-checking stops with the first
    // error found.
    Error func(err error)

    // An importer is used to import packages referred to from
    // import declarations.
    // If the installed importer implements ImporterFrom, the type
    // checker calls ImportFrom instead of Import.
    // The type checker reports an error if an importer is needed
    // but none was installed.
    Importer Importer

    // If Sizes != nil, it provides the sizing functions for package unsafe.
    // Otherwise &StdSizes{WordSize: 8, MaxAlign: 8} is used instead.
    Sizes Sizes

    // If DisableUnusedImportCheck is set, packages are not checked
    // for unused imports.
    DisableUnusedImportCheck bool
}

// A Const represents a declared constant.
type Const struct {
    val     constant.Value
    visited bool // for initialization cycle detection
}

// An Error describes a type-checking error; it implements the error interface.
// A "soft" error is an error that still permits a valid interpretation of a
// package (such as "unused variable"); "hard" errors may lead to unpredictable
// behavior if ignored.
type Error struct {
    Fset *token.FileSet // file set for interpretation of Pos
    Pos  token.Pos      // error position
    Msg  string         // error message
    Soft bool           // if set, error is "soft"
}

// A Func represents a declared function, concrete method, or abstract
// (interface) method.  Its Type() is always a *Signature.
// An abstract method may belong to many interfaces due to embedding.
type Func struct {
}

// ImportMode is reserved for future use.
type ImportMode int

// An Importer resolves import paths to Packages.
//
// CAUTION: This interface does not support the import of locally
// vendored packages. See https://golang.org/s/go15vendor.
// If possible, external implementations should implement ImporterFrom.
type Importer interface {
    // Import returns the imported package for the given import
    // path, or an error if the package couldn't be imported.
    // Two calls to Import with the same path return the same
    // package.
    Import(path string) (*Package, error)
}

// An ImporterFrom resolves import paths to packages; it
// supports vendoring per https://golang.org/s/go15vendor.
// Use go/importer to obtain an ImporterFrom implementation.
type ImporterFrom interface {
    // Importer is present for backward-compatibility. Calling
    // Import(path) is the same as calling ImportFrom(path, "", 0);
    // i.e., locally vendored packages may not be found.
    // The types package does not call Import if an ImporterFrom
    // is present.
    Importer

    // ImportFrom returns the imported package for the given import
    // path when imported by the package in srcDir, or an error
    // if the package couldn't be imported. The mode value must
    // be 0; it is reserved for future use.
    // Two calls to ImportFrom with the same path and srcDir return
    // the same package.
    ImportFrom(path, srcDir string, mode ImportMode) (*Package, error)
}

// Info holds result type information for a type-checked package.
// Only the information for which a map is provided is collected.
// If the package has type errors, the collected information may
// be incomplete.
type Info struct {
    // Types maps expressions to their types, and for constant
    // expressions, their values. Invalid expressions are omitted.
    //
    // For (possibly parenthesized) identifiers denoting built-in
    // functions, the recorded signatures are call-site specific:
    // if the call result is not a constant, the recorded type is
    // an argument-specific signature. Otherwise, the recorded type
    // is invalid.
    //
    // Identifiers on the lhs of declarations (i.e., the identifiers
    // which are being declared) are collected in the Defs map.
    // Identifiers denoting packages are collected in the Uses maps.
    Types map[ast.Expr]TypeAndValue

    // Defs maps identifiers to the objects they define (including
    // package names, dots "." of dot-imports, and blank "_" identifiers).
    // For identifiers that do not denote objects (e.g., the package name
    // in package clauses, or symbolic variables t in t := x.(type) of
    // type switch headers), the corresponding objects are nil.
    //
    // For an anonymous field, Defs returns the field *Var it defines.
    //
    // Invariant: Defs[id] == nil || Defs[id].Pos() == id.Pos()
    Defs map[*ast.Ident]Object

    // Uses maps identifiers to the objects they denote.
    //
    // For an anonymous field, Uses returns the *TypeName it denotes.
    //
    // Invariant: Uses[id].Pos() != id.Pos()
    Uses map[*ast.Ident]Object

    // Implicits maps nodes to their implicitly declared objects, if any.
    // The following node and object types may appear:
    //
    //	node               declared object
    //
    //	*ast.ImportSpec    *PkgName for dot-imports and imports without renames
    //	*ast.CaseClause    type-specific *Var for each type switch case clause (incl. default)
    //      *ast.Field         anonymous parameter *Var
    //
    Implicits map[ast.Node]Object

    // Selections maps selector expressions (excluding qualified identifiers)
    // to their corresponding selections.
    Selections map[*ast.SelectorExpr]*Selection

    // Scopes maps ast.Nodes to the scopes they define. Package scopes are not
    // associated with a specific node but with all files belonging to a package.
    // Thus, the package scope can be found in the type-checked Package object.
    // Scopes nest, with the Universe scope being the outermost scope, enclosing
    // the package scope, which contains (one or more) files scopes, which enclose
    // function scopes which in turn enclose statement and function literal scopes.
    // Note that even though package-level functions are declared in the package
    // scope, the function scopes are embedded in the file scope of the file
    // containing the function declaration.
    //
    // The following node types may appear in Scopes:
    //
    //	*ast.File
    //	*ast.FuncType
    //	*ast.BlockStmt
    //	*ast.IfStmt
    //	*ast.SwitchStmt
    //	*ast.TypeSwitchStmt
    //	*ast.CaseClause
    //	*ast.CommClause
    //	*ast.ForStmt
    //	*ast.RangeStmt
    //
    Scopes map[ast.Node]*Scope

    // InitOrder is the list of package-level initializers in the order in which
    // they must be executed. Initializers referring to variables related by an
    // initialization dependency appear in topological order, the others appear
    // in source order. Variables without an initialization expression do not
    // appear in this list.
    InitOrder []*Initializer
}

// An Initializer describes a package-level variable, or a list of variables in
// case of a multi-valued initialization expression, and the corresponding
// initialization expression.
type Initializer struct {
    Lhs []*Var // var Lhs = Rhs
    Rhs ast.Expr
}

// An Interface represents an interface type.
type Interface struct {
    methods   []*Func  // ordered list of explicitly declared methods
    embeddeds []*Named // ordered list of explicitly embedded types

    allMethods []*Func // ordered list of methods declared with or embedded in this interface (TODO(gri): replace with mset)
}

// A Label represents a declared label.
type Label struct {
    used bool // set if the label was used
}

// A Map represents a map type.
type Map struct {
    key, elem Type
}

// A MethodSet is an ordered set of concrete or abstract (interface) methods; a
// method is a MethodVal selection, and they are ordered by ascending
// m.Obj().Id(). The zero value for a MethodSet is a ready-to-use empty method
// set.
type MethodSet struct {
    list []*Selection
}

// A Named represents a named type.
type Named struct {
    obj        *TypeName // corresponding declared object
    underlying Type      // possibly a *Named during setup; never a *Named once set up completely
    methods    []*Func   // methods declared for this type (not the method set of this type)
}

// Nil represents the predeclared value nil.
type Nil struct {
}

// An Object describes a named language entity such as a package,
// constant, type, variable, function (incl. methods), or label.
// All objects implement the Object interface.
type Object interface {
    Parent() *Scope // scope in which this object is declared
    Pos() token.Pos // position of object identifier in declaration
    Pkg() *Package  // nil for objects in the Universe scope and labels
    Name() string   // package local object name
    Type() Type     // object type
    Exported() bool // reports whether the name starts with a capital letter
    Id() string     // object id (see Id below)

    // String returns a human-readable string of the object.
    String() string

    // order reflects a package-level object's source order: if object
    // a is before object b in the source, then a.order() < b.order().
    // order returns a value > 0 for package-level objects; it returns
    // 0 for all other objects (including objects in file scopes).
    order() uint32

    // setOrder sets the order number of the object. It must be > 0.
    setOrder(uint32)

    // setParent sets the parent scope of the object.
    setParent(*Scope)

    // sameId reports whether obj.Id() and Id(pkg, name) are the same.
    sameId(pkg *Package, name string) bool

    // scopePos returns the start position of the scope of this Object
    scopePos() token.Pos

    // setScopePos sets the start position of the scope for this Object.
    setScopePos(pos token.Pos)
}

// A Package describes a Go package.
type Package struct {
    path     string
    name     string
    scope    *Scope
    complete bool
    imports  []*Package
    fake     bool // scope lookup errors are silently dropped if package is fake (internal use only)
}

// A PkgName represents an imported Go package.
type PkgName struct {
    imported *Package
    used     bool // set if the package was used
}

// A Pointer represents a pointer type.
type Pointer struct {
    base Type // element type
}

// A Qualifier controls how named package-level objects are printed in
// calls to TypeString, ObjectString, and SelectionString.
//
// These three formatting routines call the Qualifier for each
// package-level object O, and if the Qualifier returns a non-empty
// string p, the object is printed in the form p.O.
// If it returns an empty string, only the object name O is printed.
//
// Using a nil Qualifier is equivalent to using (*Package).Path: the
// object is qualified by the import path, e.g., "encoding/json.Marshal".
type Qualifier func(*Package) string

// A Scope maintains a set of objects and links to its containing
// (parent) and contained (children) scopes. Objects may be inserted
// and looked up by name. The zero value for Scope is a ready-to-use
// empty scope.
type Scope struct {
    parent   *Scope
    children []*Scope
    elems    map[string]Object // lazily allocated
    pos, end token.Pos         // scope extent; may be invalid
    comment  string            // for debugging only
}

// A Selection describes a selector expression x.f. For the declarations:
//
//     type T struct{ x int; E }
//     type E struct{}
//     func (e E) m() {}
//     var p *T
//
// the following relations exist:
//
//     Selector    Kind          Recv    Obj    Type               Index     Indirect
//
//     p.x         FieldVal      T       x      int                {0}       true
//     p.m         MethodVal     *T      m      func (e *T) m()    {1, 0}    true
//     T.m         MethodExpr    T       m      func m(_ T)        {1, 0}    false
type Selection struct {
    kind     SelectionKind
    recv     Type   // type of x
    obj      Object // object denoted by x.f
    index    []int  // path from x to x.f
    indirect bool   // set if there was any pointer indirection on the path
}

// SelectionKind describes the kind of a selector expression x.f
// (excluding qualified identifiers).
type SelectionKind int

// A Signature represents a (non-builtin) function or method type.
type Signature struct {
    // We need to keep the scope in Signature (rather than passing it around
    // and store it in the Func Object) because when type-checking a function
    // literal we call the general type checker which returns a general Type.
    // We then unpack the *Signature and use the scope for the literal body.
    scope    *Scope // function scope, present for package-local signatures
    recv     *Var   // nil if not a method
    params   *Tuple // (incoming) parameters from left to right; or nil
    results  *Tuple // (outgoing) results from left to right; or nil
    variadic bool   // true if the last parameter's type is of the form ...T (or string, for append built-in only)
}

// Sizes defines the sizing functions for package unsafe.
type Sizes interface {
    // Alignof returns the alignment of a variable of type T.
    // Alignof must implement the alignment guarantees required by the spec.
    Alignof(T Type) int64

    // Offsetsof returns the offsets of the given struct fields, in bytes.
    // Offsetsof must implement the offset guarantees required by the spec.
    Offsetsof(fields []*Var) []int64

    // Sizeof returns the size of a variable of type T.
    // Sizeof must implement the size guarantees required by the spec.
    Sizeof(T Type) int64
}

// A Slice represents a slice type.
type Slice struct {
    elem Type
}

// StdSizes is a convenience type for creating commonly used Sizes.
// It makes the following simplifying assumptions:
//
//     - The size of explicitly sized basic types (int16, etc.) is the
//       specified size.
//     - The size of strings and interfaces is 2*WordSize.
//     - The size of slices is 3*WordSize.
//     - The size of an array of n elements corresponds to the size of
//       a struct of n consecutive fields of the array's element type.
//      - The size of a struct is the offset of the last field plus that
//       field's size. As with all element types, if the struct is used
//       in an array its size must first be aligned to a multiple of the
//       struct's alignment.
//     - All other types have size WordSize.
//     - Arrays and structs are aligned per spec definition; all other
//       types are naturally aligned with a maximum alignment MaxAlign.
//
// *StdSizes implements Sizes.
type StdSizes struct {
    WordSize int64 // word size in bytes - must be >= 4 (32bits)
    MaxAlign int64 // maximum alignment in bytes - must be >= 1
}

// A Struct represents a struct type.
type Struct struct {
    fields      []*Var
    tags        []string  // field tags; nil if there are no tags
    offsets     []int64   // field offsets in bytes, lazily initialized
    offsetsOnce sync.Once // for threadsafe lazy initialization of offsets
}

// A Tuple represents an ordered list of variables; a nil *Tuple is a valid
// (empty) tuple. Tuples are used as components of signatures and to represent
// the type of multiple assignments; they are not first class types of Go.
type Tuple struct {
    vars []*Var
}

// A Type represents a type of Go.
// All types implement the Type interface.
type Type interface {
    // Underlying returns the underlying type of a type.
    Underlying() Type

    // String returns a string representation of a type.
    String() string
}

// TypeAndValue reports the type and value (for constants)
// of the corresponding expression.
type TypeAndValue struct {
    mode  operandMode
    Type  Type
    Value constant.Value
}

// A TypeName represents a declared type.
type TypeName struct {
}

// A Variable represents a declared variable (including function parameters and
// results, and struct fields).
type Var struct {
    anonymous bool // if set, the variable is an anonymous struct field, and name is the type name
    visited   bool // for initialization cycle detection
    isField   bool // var is struct field
    used      bool // set if the variable was used
}

// AssertableTo reports whether a value of type V can be asserted to have type
// T.
func AssertableTo(V *Interface, T Type) bool

// AssignableTo reports whether a value of type V is assignable to a variable of
// type T.
func AssignableTo(V, T Type) bool

// Comparable reports whether values of type T are comparable.
func Comparable(T Type) bool

// ConvertibleTo reports whether a value of type V is convertible to a value of
// type T.
func ConvertibleTo(V, T Type) bool

// DefPredeclaredTestFuncs defines the assert and trace built-ins.
// These built-ins are intended for debugging and testing of this
// package only.
func DefPredeclaredTestFuncs()

// Eval returns the type and, if constant, the value for the
// expression expr, evaluated at position pos of package pkg,
// which must have been derived from type-checking an AST with
// complete position information relative to the provided file
// set.
//
// If the expression contains function literals, their bodies
// are ignored (i.e., the bodies are not type-checked).
//
// If pkg == nil, the Universe scope is used and the provided
// position pos is ignored. If pkg != nil, and pos is invalid,
// the package scope is used. Otherwise, pos must belong to the
// package.
//
// An error is returned if pos is not within the package or
// if the node cannot be evaluated.
//
// Note: Eval should not be used instead of running Check to compute
// types and values, but in addition to Check. Eval will re-evaluate
// its argument each time, and it also does not know about the context
// in which an expression is used (e.g., an assignment). Thus, top-
// level untyped constants will return an untyped type rather then the
// respective context-specific type.
func Eval(fset *token.FileSet, pkg *Package, pos token.Pos, expr string) (tv TypeAndValue, err error)

// ExprString returns the (possibly simplified) string representation for x.
func ExprString(x ast.Expr) string

// Id returns name if it is exported, otherwise it
// returns the name qualified with the package path.
func Id(pkg *Package, name string) string

// Identical reports whether x and y are identical.
func Identical(x, y Type) bool

// Implements reports whether type V implements interface T.
func Implements(V Type, T *Interface) bool

// IsInterface reports whether typ is an interface type.
func IsInterface(typ Type) bool

// LookupFieldOrMethod looks up a field or method with given package and name
// in T and returns the corresponding *Var or *Func, an index sequence, and a
// bool indicating if there were any pointer indirections on the path to the
// field or method. If addressable is set, T is the type of an addressable
// variable (only matters for method lookups).
//
// The last index entry is the field or method index in the (possibly embedded)
// type where the entry was found, either:
//
//     1) the list of declared methods of a named type; or
//     2) the list of all methods (method set) of an interface type; or
//     3) the list of fields of a struct type.
//
// The earlier index entries are the indices of the anonymous struct fields
// traversed to get to the found entry, starting at depth 0.
//
// If no entry is found, a nil object is returned. In this case, the returned
// index and indirect values have the following meaning:
//
//     - If index != nil, the index sequence points to an ambiguous entry
//     (the same name appeared more than once at the same embedding level).
//
//     - If indirect is set, a method with a pointer receiver type was found
//      but there was no pointer on the path from the actual receiver type to
//     the method's formal receiver base type, nor was the receiver addressable.
func LookupFieldOrMethod(T Type, addressable bool, pkg *Package, name string) (obj Object, index []int, indirect bool)

// MissingMethod returns (nil, false) if V implements T, otherwise it
// returns a missing method required by T and whether it is missing or
// just has the wrong type.
//
// For non-interface types V, or if static is set, V implements T if all
// methods of T are present in V. Otherwise (V is an interface and static
// is not set), MissingMethod only checks that methods of T which are also
// present in V have matching types (e.g., for a type assertion x.(T) where
// x is of interface type V).
func MissingMethod(V Type, T *Interface, static bool) (method *Func, wrongType bool)

// NewArray returns a new array type for the given element type and length.
func NewArray(elem Type, len int64) *Array

// NewChan returns a new channel type for the given direction and element type.
func NewChan(dir ChanDir, elem Type) *Chan

// NewChecker returns a new Checker instance for a given package.
// Package files may be added incrementally via checker.Files.
func NewChecker(conf *Config, fset *token.FileSet, pkg *Package, info *Info) *Checker

func NewConst(pos token.Pos, pkg *Package, name string, typ Type, val constant.Value) *Const

func NewField(pos token.Pos, pkg *Package, name string, typ Type, anonymous bool) *Var

func NewFunc(pos token.Pos, pkg *Package, name string, sig *Signature) *Func

// NewInterface returns a new interface for the given methods and embedded
// types.
func NewInterface(methods []*Func, embeddeds []*Named) *Interface

func NewLabel(pos token.Pos, pkg *Package, name string) *Label

// NewMap returns a new map for the given key and element types.
func NewMap(key, elem Type) *Map

// NewMethodSet returns the method set for the given type T.
// It always returns a non-nil method set, even if it is empty.
func NewMethodSet(T Type) *MethodSet

// NewNamed returns a new named type for the given type name, underlying type,
// and associated methods. The underlying type must not be a *Named.
func NewNamed(obj *TypeName, underlying Type, methods []*Func) *Named

// NewPackage returns a new Package for the given package path and name;
// the name must not be the blank identifier.
// The package is not complete and contains no explicit imports.
func NewPackage(path, name string) *Package

func NewParam(pos token.Pos, pkg *Package, name string, typ Type) *Var

func NewPkgName(pos token.Pos, pkg *Package, name string, imported *Package) *PkgName

// NewPointer returns a new pointer type for the given element (base) type.
func NewPointer(elem Type) *Pointer

// NewScope returns a new, empty scope contained in the given parent
// scope, if any.  The comment is for debugging only.
func NewScope(parent *Scope, pos, end token.Pos, comment string) *Scope

// NewSignature returns a new function type for the given receiver, parameters,
// and results, either of which may be nil. If variadic is set, the function
// is variadic, it must have at least one parameter, and the last parameter
// must be of unnamed slice type.
func NewSignature(recv *Var, params, results *Tuple, variadic bool) *Signature

// NewSlice returns a new slice type for the given element type.
func NewSlice(elem Type) *Slice

// NewStruct returns a new struct with the given fields and corresponding field
// tags. If a field with index i has a tag, tags[i] must be that tag, but
// len(tags) may be only as long as required to hold the tag with the largest
// index i. Consequently, if no field has a tag, tags may be nil.
func NewStruct(fields []*Var, tags []string) *Struct

// NewTuple returns a new tuple for the given variables.
func NewTuple(x ...*Var) *Tuple

func NewTypeName(pos token.Pos, pkg *Package, name string, typ Type) *TypeName

func NewVar(pos token.Pos, pkg *Package, name string, typ Type) *Var

// ObjectString returns the string form of obj.
// The Qualifier controls the printing of
// package-level objects, and may be nil.
func ObjectString(obj Object, qf Qualifier) string

// RelativeTo(pkg) returns a Qualifier that fully qualifies members of
// all packages other than pkg.
func RelativeTo(pkg *Package) Qualifier

// SelectionString returns the string form of s.
// The Qualifier controls the printing of
// package-level objects, and may be nil.
//
// Examples:
//     "field (T) f int"
//     "method (T) f(X) Y"
//     "method expr (T) f(X) Y"
func SelectionString(s *Selection, qf Qualifier) string

func TestAssignOp(t *testing.T)

func TestZeroTok(t *testing.T)

// TypeString returns the string representation of typ.
// The Qualifier controls the printing of
// package-level objects, and may be nil.
func TypeString(typ Type, qf Qualifier) string

// WriteExpr writes the (possibly simplified) string representation for x to
// buf.
func WriteExpr(buf *bytes.Buffer, x ast.Expr)

// WriteSignature writes the representation of the signature sig to buf,
// without a leading "func" keyword.
// The Qualifier controls the printing of
// package-level objects, and may be nil.
func WriteSignature(buf *bytes.Buffer, sig *Signature, qf Qualifier)

// WriteType writes the string representation of typ to buf.
// The Qualifier controls the printing of
// package-level objects, and may be nil.
func WriteType(buf *bytes.Buffer, typ Type, qf Qualifier)

// Elem returns element type of array a.
func (*Array) Elem() Type

// Len returns the length of array a.
func (*Array) Len() int64

func (*Array) String() string

func (*Array) Underlying() Type

// Info returns information about properties of basic type b.
func (*Basic) Info() BasicInfo

// Kind returns the kind of basic type b.
func (*Basic) Kind() BasicKind

// Name returns the name of basic type b.
func (*Basic) Name() string

func (*Basic) String() string

func (*Basic) Underlying() Type

func (*Builtin) String() string

// Dir returns the direction of channel c.
func (*Chan) Dir() ChanDir

// Elem returns the element type of channel c.
func (*Chan) Elem() Type

func (*Chan) String() string

func (*Chan) Underlying() Type

// Files checks the provided files as part of the checker's package.
func (*Checker) Files(files []*ast.File) (err error)

// Check type-checks a package and returns the resulting package object and
// the first error if any. Additionally, if info != nil, Check populates each
// of the non-nil maps in the Info struct.
//
// The package is marked as complete if no errors occurred, otherwise it is
// incomplete. See Config.Error for controlling behavior in the presence of
// errors.
//
// The package is specified by a list of *ast.Files and corresponding
// file set, and the package path the package is identified with.
// The clean path must not be empty or dot (".").
func (*Config) Check(path string, fset *token.FileSet, files []*ast.File, info *Info) (*Package, error)

func (*Const) String() string

func (*Const) Val() constant.Value

// FullName returns the package- or receiver-type-qualified name of
// function or method obj.
func (*Func) FullName() string

func (*Func) Scope() *Scope

func (*Func) String() string

// ObjectOf returns the object denoted by the specified id,
// or nil if not found.
//
// If id is an anonymous struct field, ObjectOf returns the field (*Var)
// it uses, not the type (*TypeName) it defines.
//
// Precondition: the Uses and Defs maps are populated.
func (*Info) ObjectOf(id *ast.Ident) Object

// TypeOf returns the type of expression e, or nil if not found.
// Precondition: the Types, Uses and Defs maps are populated.
func (*Info) TypeOf(e ast.Expr) Type

func (*Initializer) String() string

// Complete computes the interface's method set. It must be called by users of
// NewInterface after the interface's embedded types are fully defined and
// before using the interface type in any way other than to form other types.
// Complete returns the receiver.
func (*Interface) Complete() *Interface

// Embedded returns the i'th embedded type of interface t for 0 <= i <
// t.NumEmbeddeds(). The types are ordered by the corresponding TypeName's
// unique Id.
func (*Interface) Embedded(i int) *Named

// Empty returns true if t is the empty interface.
func (*Interface) Empty() bool

// ExplicitMethod returns the i'th explicitly declared method of interface t for
// 0 <= i < t.NumExplicitMethods(). The methods are ordered by their unique Id.
func (*Interface) ExplicitMethod(i int) *Func

// Method returns the i'th method of interface t for 0 <= i < t.NumMethods().
// The methods are ordered by their unique Id.
func (*Interface) Method(i int) *Func

// NumEmbeddeds returns the number of embedded types in interface t.
func (*Interface) NumEmbeddeds() int

// NumExplicitMethods returns the number of explicitly declared methods of
// interface t.
func (*Interface) NumExplicitMethods() int

// NumMethods returns the total number of methods of interface t.
func (*Interface) NumMethods() int

func (*Interface) String() string

func (*Interface) Underlying() Type

func (*Label) String() string

// Elem returns the element type of map m.
func (*Map) Elem() Type

// Key returns the key type of map m.
func (*Map) Key() Type

func (*Map) String() string

func (*Map) Underlying() Type

// At returns the i'th method in s for 0 <= i < s.Len().
func (*MethodSet) At(i int) *Selection

// Len returns the number of methods in s.
func (*MethodSet) Len() int

// Lookup returns the method with matching package and name, or nil if not
// found.
func (*MethodSet) Lookup(pkg *Package, name string) *Selection

func (*MethodSet) String() string

// AddMethod adds method m unless it is already in the method list.
// TODO(gri) find a better solution instead of providing this function
func (*Named) AddMethod(m *Func)

// Method returns the i'th method of named type t for 0 <= i < t.NumMethods().
func (*Named) Method(i int) *Func

// NumMethods returns the number of explicit methods whose receiver is named
// type t.
func (*Named) NumMethods() int

// TypeName returns the type name for the named type t.
func (*Named) Obj() *TypeName

// SetUnderlying sets the underlying type and marks t as complete. TODO(gri)
// determine if there's a better solution rather than providing this function
func (*Named) SetUnderlying(underlying Type)

func (*Named) String() string

func (*Named) Underlying() Type

func (*Nil) String() string

// A package is complete if its scope contains (at least) all
// exported objects; otherwise it is incomplete.
func (*Package) Complete() bool

// Imports returns the list of packages directly imported by
// pkg; the list is in source order. Package unsafe is excluded.
//
// If pkg was loaded from export data, Imports includes packages that
// provide package-level objects referenced by pkg.  This may be more or
// less than the set of packages directly imported by pkg's source code.
func (*Package) Imports() []*Package

// MarkComplete marks a package as complete.
func (*Package) MarkComplete()

// Name returns the package name.
func (*Package) Name() string

// Path returns the package path.
func (*Package) Path() string

// Scope returns the (complete or incomplete) package scope
// holding the objects declared at package level (TypeNames,
// Consts, Vars, and Funcs).
func (*Package) Scope() *Scope

// SetImports sets the list of explicitly imported packages to list.
// It is the caller's responsibility to make sure list elements are unique.
func (*Package) SetImports(list []*Package)

// SetName sets the package name.
func (*Package) SetName(name string)

func (*Package) String() string

// Imported returns the package that was imported. It is distinct from Pkg(),
// which is the package containing the import statement.
func (*PkgName) Imported() *Package

func (*PkgName) String() string

// Elem returns the element type for the given pointer p.
func (*Pointer) Elem() Type

func (*Pointer) String() string

func (*Pointer) Underlying() Type

// Child returns the i'th child scope for 0 <= i < NumChildren().
func (*Scope) Child(i int) *Scope

// Contains returns true if pos is within the scope's extent.
// The result is guaranteed to be valid only if the type-checked
// AST has complete position information.
func (*Scope) Contains(pos token.Pos) bool

func (*Scope) End() token.Pos

// Innermost returns the innermost (child) scope containing
// pos. If pos is not within any scope, the result is nil.
// The result is also nil for the Universe scope.
// The result is guaranteed to be valid only if the type-checked
// AST has complete position information.
func (*Scope) Innermost(pos token.Pos) *Scope

// Insert attempts to insert an object obj into scope s.
// If s already contains an alternative object alt with
// the same name, Insert leaves s unchanged and returns alt.
// Otherwise it inserts obj, sets the object's parent scope
// if not already set, and returns nil.
func (*Scope) Insert(obj Object) Object

// Len() returns the number of scope elements.
func (*Scope) Len() int

// Lookup returns the object in scope s with the given name if such an
// object exists; otherwise the result is nil.
func (*Scope) Lookup(name string) Object

// LookupParent follows the parent chain of scopes starting with s until
// it finds a scope where Lookup(name) returns a non-nil object, and then
// returns that scope and object. If a valid position pos is provided,
// only objects that were declared at or before pos are considered.
// If no such scope and object exists, the result is (nil, nil).
//
// Note that obj.Parent() may be different from the returned scope if the
// object was inserted into the scope and already had a parent at that
// time (see Insert, below). This can only happen for dot-imported objects
// whose scope is the scope of the package that exported them.
func (*Scope) LookupParent(name string, pos token.Pos) (*Scope, Object)

// Names returns the scope's element names in sorted order.
func (*Scope) Names() []string

// NumChildren() returns the number of scopes nested in s.
func (*Scope) NumChildren() int

// Parent returns the scope's containing (parent) scope.
func (*Scope) Parent() *Scope

// Pos and End describe the scope's source code extent [pos, end).
// The results are guaranteed to be valid only if the type-checked
// AST has complete position information. The extent is undefined
// for Universe and package scopes.
func (*Scope) Pos() token.Pos

// String returns a string representation of the scope, for debugging.
func (*Scope) String() string

// WriteTo writes a string representation of the scope to w,
// with the scope elements sorted by name.
// The level of indentation is controlled by n >= 0, with
// n == 0 for no indentation.
// If recurse is set, it also writes nested (children) scopes.
func (*Scope) WriteTo(w io.Writer, n int, recurse bool)

// Index describes the path from x to f in x.f.
// The last index entry is the field or method index of the type declaring f;
// either:
//
//     1) the list of declared methods of a named type; or
//     2) the list of methods of an interface type; or
//     3) the list of fields of a struct type.
//
// The earlier index entries are the indices of the embedded fields implicitly
// traversed to get from (the type of) x to f, starting at embedding depth 0.
func (*Selection) Index() []int

// Indirect reports whether any pointer indirection was required to get from
// x to f in x.f.
func (*Selection) Indirect() bool

// Kind returns the selection kind.
func (*Selection) Kind() SelectionKind

// Obj returns the object denoted by x.f; a *Var for
// a field selection, and a *Func in all other cases.
func (*Selection) Obj() Object

// Recv returns the type of x in x.f.
func (*Selection) Recv() Type

func (*Selection) String() string

// Type returns the type of x.f, which may be different from the type of f.
// See Selection for more information.
func (*Selection) Type() Type

// Params returns the parameters of signature s, or nil.
func (*Signature) Params() *Tuple

// Recv returns the receiver of signature s (if a method), or nil if a
// function.
//
// For an abstract method, Recv returns the enclosing interface either
// as a *Named or an *Interface.  Due to embedding, an interface may
// contain methods whose receiver type is a different interface.
func (*Signature) Recv() *Var

// Results returns the results of signature s, or nil.
func (*Signature) Results() *Tuple

func (*Signature) String() string

func (*Signature) Underlying() Type

// Variadic reports whether the signature s is variadic.
func (*Signature) Variadic() bool

// Elem returns the element type of slice s.
func (*Slice) Elem() Type

func (*Slice) String() string

func (*Slice) Underlying() Type

func (*StdSizes) Alignof(T Type) int64

func (*StdSizes) Offsetsof(fields []*Var) []int64

func (*StdSizes) Sizeof(T Type) int64

// Field returns the i'th field for 0 <= i < NumFields().
func (*Struct) Field(i int) *Var

// NumFields returns the number of fields in the struct (including blank and
// anonymous fields).
func (*Struct) NumFields() int

func (*Struct) String() string

// Tag returns the i'th field tag for 0 <= i < NumFields().
func (*Struct) Tag(i int) string

func (*Struct) Underlying() Type

// At returns the i'th variable of tuple t.
func (*Tuple) At(i int) *Var

// Len returns the number variables of tuple t.
func (*Tuple) Len() int

func (*Tuple) String() string

func (*Tuple) Underlying() Type

func (*TypeName) String() string

func (*Var) Anonymous() bool

func (*Var) IsField() bool

func (*Var) String() string

// Error returns an error string formatted as follows:
// filename:line:column: message
func (Error) Error() string

// Addressable reports whether the corresponding expression
// is addressable (https://golang.org/ref/spec#Address_operators).
func (TypeAndValue) Addressable() bool

// Assignable reports whether the corresponding expression
// is assignable to (provided a value of the right type).
func (TypeAndValue) Assignable() bool

// HasOk reports whether the corresponding expression may be
// used on the lhs of a comma-ok assignment.
func (TypeAndValue) HasOk() bool

// IsBuiltin reports whether the corresponding expression denotes
// a (possibly parenthesized) built-in function.
func (TypeAndValue) IsBuiltin() bool

// IsNil reports whether the corresponding expression denotes the
// predeclared value nil.
func (TypeAndValue) IsNil() bool

// IsType reports whether the corresponding expression specifies a type.
func (TypeAndValue) IsType() bool

// IsValue reports whether the corresponding expression is a value.
// Builtins are not considered values. Constant values have a non-
// nil Value.
func (TypeAndValue) IsValue() bool

// IsVoid reports whether the corresponding expression
// is a function call without results.
func (TypeAndValue) IsVoid() bool

