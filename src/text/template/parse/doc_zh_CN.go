// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package parse builds parse trees for templates as defined by text/template and
// html/template. Clients should use those packages to construct templates rather
// than this one, which provides shared internal data structures not intended for
// general use.
package parse

// IsEmptyTree reports whether this tree (node) is empty of everything but space.
func IsEmptyTree(n Node) bool

// Parse returns a map from template name to parse.Tree, created by parsing the
// templates described in the argument string. The top-level template will be given
// the specified name. If an error is encountered, parsing stops and an empty map
// is returned with the error.
func Parse(name, text, leftDelim, rightDelim string, funcs ...map[string]interface{}) (treeSet map[string]*Tree, err error)

// ActionNode holds an action (something bounded by delimiters). Control actions
// have their own nodes; ActionNode represents simple ones such as field
// evaluations and parenthesized pipelines.
type ActionNode struct {
	NodeType
	Pos

	Line int       // The line number in the input (deprecated; kept for compatibility)
	Pipe *PipeNode // The pipeline in the action.
	// contains filtered or unexported fields
}

func (a *ActionNode) Copy() Node

func (a *ActionNode) String() string

// BoolNode holds a boolean constant.
type BoolNode struct {
	NodeType
	Pos

	True bool // The value of the boolean constant.
	// contains filtered or unexported fields
}

func (b *BoolNode) Copy() Node

func (b *BoolNode) String() string

// BranchNode is the common representation of if, range, and with.
type BranchNode struct {
	NodeType
	Pos

	Line     int       // The line number in the input (deprecated; kept for compatibility)
	Pipe     *PipeNode // The pipeline to be evaluated.
	List     *ListNode // What to execute if the value is non-empty.
	ElseList *ListNode // What to execute if the value is empty (nil if absent).
	// contains filtered or unexported fields
}

func (b *BranchNode) Copy() Node

func (b *BranchNode) String() string

// ChainNode holds a term followed by a chain of field accesses (identifier
// starting with '.'). The names may be chained ('.x.y'). The periods are dropped
// from each ident.
type ChainNode struct {
	NodeType
	Pos

	Node  Node
	Field []string // The identifiers in lexical order.
	// contains filtered or unexported fields
}

// Add adds the named field (which should start with a period) to the end of the
// chain.
func (c *ChainNode) Add(field string)

func (c *ChainNode) Copy() Node

func (c *ChainNode) String() string

// CommandNode holds a command (a pipeline inside an evaluating action).
type CommandNode struct {
	NodeType
	Pos

	Args []Node // Arguments in lexical order: Identifier, field, or constant.
	// contains filtered or unexported fields
}

func (c *CommandNode) Copy() Node

func (c *CommandNode) String() string

// DotNode holds the special identifier '.'.
type DotNode struct {
	NodeType
	Pos
	// contains filtered or unexported fields
}

func (d *DotNode) Copy() Node

func (d *DotNode) String() string

func (d *DotNode) Type() NodeType

// FieldNode holds a field (identifier starting with '.'). The names may be chained
// ('.x.y'). The period is dropped from each ident.
type FieldNode struct {
	NodeType
	Pos

	Ident []string // The identifiers in lexical order.
	// contains filtered or unexported fields
}

func (f *FieldNode) Copy() Node

func (f *FieldNode) String() string

// IdentifierNode holds an identifier.
type IdentifierNode struct {
	NodeType
	Pos

	Ident string // The identifier's name.
	// contains filtered or unexported fields
}

// NewIdentifier returns a new IdentifierNode with the given identifier name.
func NewIdentifier(ident string) *IdentifierNode

func (i *IdentifierNode) Copy() Node

// SetPos sets the position. NewIdentifier is a public method so we can't modify
// its signature. Chained for convenience. TODO: fix one day?
func (i *IdentifierNode) SetPos(pos Pos) *IdentifierNode

// SetTree sets the parent tree for the node. NewIdentifier is a public method so
// we can't modify its signature. Chained for convenience. TODO: fix one day?
func (i *IdentifierNode) SetTree(t *Tree) *IdentifierNode

func (i *IdentifierNode) String() string

// IfNode represents an {{if}} action and its commands.
type IfNode struct {
	BranchNode
}

func (i *IfNode) Copy() Node

// ListNode holds a sequence of nodes.
type ListNode struct {
	NodeType
	Pos

	Nodes []Node // The element nodes in lexical order.
	// contains filtered or unexported fields
}

func (l *ListNode) Copy() Node

func (l *ListNode) CopyList() *ListNode

func (l *ListNode) String() string

// NilNode holds the special identifier 'nil' representing an untyped nil constant.
type NilNode struct {
	NodeType
	Pos
	// contains filtered or unexported fields
}

func (n *NilNode) Copy() Node

func (n *NilNode) String() string

func (n *NilNode) Type() NodeType

// A Node is an element in the parse tree. The interface is trivial. The interface
// contains an unexported method so that only types local to this package can
// satisfy it.
type Node interface {
	Type() NodeType
	String() string
	// Copy does a deep copy of the Node and all its components.
	// To avoid type assertions, some XxxNodes also have specialized
	// CopyXxx methods that return *XxxNode.
	Copy() Node
	Position() Pos // byte position of start of node in full original input string
	// contains filtered or unexported methods
}

// NodeType identifies the type of a parse tree node.
type NodeType int

const (
	NodeText    NodeType = iota // Plain text.
	NodeAction                  // A non-control action such as a field evaluation.
	NodeBool                    // A boolean constant.
	NodeChain                   // A sequence of field accesses.
	NodeCommand                 // An element of a pipeline.
	NodeDot                     // The cursor, dot.

	NodeField      // A field or method name.
	NodeIdentifier // An identifier; always a function name.
	NodeIf         // An if action.
	NodeList       // A list of Nodes.
	NodeNil        // An untyped nil constant.
	NodeNumber     // A numerical constant.
	NodePipe       // A pipeline of commands.
	NodeRange      // A range action.
	NodeString     // A string constant.
	NodeTemplate   // A template invocation action.
	NodeVariable   // A $ variable.
	NodeWith       // A with action.
)

// Type returns itself and provides an easy default implementation for embedding in
// a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType

// NumberNode holds a number: signed or unsigned integer, float, or complex. The
// value is parsed and stored under all the types that can represent the value.
// This simulates in a small amount of code the behavior of Go's ideal constants.
type NumberNode struct {
	NodeType
	Pos

	IsInt      bool       // Number has an integral value.
	IsUint     bool       // Number has an unsigned integral value.
	IsFloat    bool       // Number has a floating-point value.
	IsComplex  bool       // Number is complex.
	Int64      int64      // The signed integer value.
	Uint64     uint64     // The unsigned integer value.
	Float64    float64    // The floating-point value.
	Complex128 complex128 // The complex value.
	Text       string     // The original textual representation from the input.
	// contains filtered or unexported fields
}

func (n *NumberNode) Copy() Node

func (n *NumberNode) String() string

// PipeNode holds a pipeline with optional declaration
type PipeNode struct {
	NodeType
	Pos

	Line int             // The line number in the input (deprecated; kept for compatibility)
	Decl []*VariableNode // Variable declarations in lexical order.
	Cmds []*CommandNode  // The commands in lexical order.
	// contains filtered or unexported fields
}

func (p *PipeNode) Copy() Node

func (p *PipeNode) CopyPipe() *PipeNode

func (p *PipeNode) String() string

// Pos represents a byte position in the original input text from which this
// template was parsed.
type Pos int

func (p Pos) Position() Pos

// RangeNode represents a {{range}} action and its commands.
type RangeNode struct {
	BranchNode
}

func (r *RangeNode) Copy() Node

// StringNode holds a string constant. The value has been "unquoted".
type StringNode struct {
	NodeType
	Pos

	Quoted string // The original text of the string, with quotes.
	Text   string // The string, after quote processing.
	// contains filtered or unexported fields
}

func (s *StringNode) Copy() Node

func (s *StringNode) String() string

// TemplateNode represents a {{template}} action.
type TemplateNode struct {
	NodeType
	Pos

	Line int       // The line number in the input (deprecated; kept for compatibility)
	Name string    // The name of the template (unquoted).
	Pipe *PipeNode // The command to evaluate as dot for the template.
	// contains filtered or unexported fields
}

func (t *TemplateNode) Copy() Node

func (t *TemplateNode) String() string

// TextNode holds plain text.
type TextNode struct {
	NodeType
	Pos

	Text []byte // The text; may span newlines.
	// contains filtered or unexported fields
}

func (t *TextNode) Copy() Node

func (t *TextNode) String() string

// Tree is the representation of a single parsed template.
type Tree struct {
	Name      string    // name of the template represented by the tree.
	ParseName string    // name of the top-level template during parsing, for error messages.
	Root      *ListNode // top-level root of the tree.
	// contains filtered or unexported fields
}

// New allocates a new parse tree with the given name.
func New(name string, funcs ...map[string]interface{}) *Tree

// Copy returns a copy of the Tree. Any parsing state is discarded.
func (t *Tree) Copy() *Tree

// ErrorContext returns a textual representation of the location of the node in the
// input text. The receiver is only used when the node does not have a pointer to
// the tree inside, which can occur in old code.
func (t *Tree) ErrorContext(n Node) (location, context string)

// Parse parses the template definition string to construct a representation of the
// template for execution. If either action delimiter string is empty, the default
// ("{{" or "}}") is used. Embedded template definitions are added to the treeSet
// map.
func (t *Tree) Parse(text, leftDelim, rightDelim string, treeSet map[string]*Tree, funcs ...map[string]interface{}) (tree *Tree, err error)

// VariableNode holds a list of variable names, possibly with chained field
// accesses. The dollar sign is part of the (first) name.
type VariableNode struct {
	NodeType
	Pos

	Ident []string // Variable name and fields in lexical order.
	// contains filtered or unexported fields
}

func (v *VariableNode) Copy() Node

func (v *VariableNode) String() string

// WithNode represents a {{with}} action and its commands.
type WithNode struct {
	BranchNode
}

func (w *WithNode) Copy() Node
