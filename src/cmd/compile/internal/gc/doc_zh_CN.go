// +build ingore

package gc

import (
	"bufio"
	"bytes"
	"cmd/compile/internal/big"
	"cmd/compile/internal/ssa"
	"cmd/internal/bio"
	"cmd/internal/gcprog"
	"cmd/internal/obj"
	"cmd/internal/obj/ppc64"
	"cmd/internal/sys"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"html"
	"io"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	// These values are known by runtime.
	ANOEQ AlgKind = iota
	AMEM0
	AMEM8
	AMEM16
	AMEM32
	AMEM64
	AMEM128
	ASTRING
	AINTER
	ANILINTER
	AFLOAT32
	AFLOAT64
	ACPLX64
	ACPLX128

	// Type can be compared/hashed as regular memory.

	// These values are known by runtime.
	// The MEMx and NOEQx values must run in parallel.  See algtype.
	AMEM AlgKind = 100

	// Type needs special comparison/hashing functions.
	ASPECIAL AlgKind = -1
)

// architecture-independent object file output
const (
	ArhdrSize = 60
)

const (
	BITS = 3
	NVAR = BITS * 64
)

// Builds a type representing a Bucket structure for
// the given map type. This type is not visible to users -
// we include only enough information to generate a correct GC
// program for it.
// Make sure this stays in sync with ../../../../runtime/hashmap.go!

// Builds a type representing a Bucket structure for
// the given map type.  This type is not visible to users -
// we include only enough information to generate a correct GC
// program for it.
// Make sure this stays in sync with ../../runtime/hashmap.go!
const (
	BUCKETSIZE = 8
	MAXKEYSIZE = 128
	MAXVALSIZE = 128
)

// Cost parameters
const (
	CLOAD = 5 // cost of load
	CREF  = 5 // cost of reference if not registerized
	LOOP  = 3 // loop execution count (applied in popt.go)
)

const (
	CTxxx Ctype = iota
	CTINT
	CTRUNE
	CTFLT
	CTCPLX
	CTSTR
	CTBOOL
	CTNIL
)

const (
	// types of channel
	// must match ../../../../reflect/type.go:/ChanDir
	Crecv ChanDir = 1 << 0
	Csend ChanDir = 1 << 1
	Cboth ChanDir = Crecv | Csend
)

const (
	EOF = -1
	BOM = 0xFEFF
)

const (
	EscFuncUnknown = 0 + iota
	EscFuncPlanned
	EscFuncStarted
	EscFuncTagged
)

// Escape constants are numbered in order of increasing "escapiness"
// to help make inferences be monotonic. With the exception of
// EscNever which is sticky, eX < eY means that eY is more exposed
// than eX, and hence replaces it in a conservative analysis.

// Escape constants are numbered in order of increasing "escapiness"
// to help make inferences be monotonic.  With the exception of
// EscNever which is sticky, eX < eY means that eY is more exposed
// than eX, and hence replaces it in a conservative analysis.
const (
	EscUnknown = iota
	EscNone    // Does not escape to heap, result, or parameters.
	EscReturn  // Is returned or reachable from returned.
	EscScope   // Allocated in an inner loop scope, assigned to an outer loop scope,

	// which allows the construction of non-escaping but arbitrarily large
	// linked data structures (i.e., not eligible for allocation in a fixed-size
	// stack frame).
	EscHeap           // Reachable from the heap
	EscNever          // By construction will not escape.
	EscBits           = 3
	EscMask           = (1 << EscBits) - 1
	EscContentEscapes = 1 << EscBits // value obtained by indirect of parameter escapes to heap
	EscReturnBits     = EscBits + 1
)

const (
	Etop      = 1 << iota // evaluated at statement level
	Erv                   // evaluated in value context
	Etype                 // evaluated in type context
	Ecall                 // call-only expressions are ok
	Efnstruct             // multivalue function returns are ok
	Easgn                 // assigning to expression
	Ecomplit              // type in composite literal
)

// Format conversions
//
// 		%L int		Line numbers
//
// 		%E int		etype values (aka 'Kind')
//
// 		%O int		Node Opcodes
// 			Flags: "%#O": print go syntax. (automatic unless fmtmode == FDbg)
//
// 		%J Node*	Node details
// 			Flags: "%hJ" suppresses things not relevant until walk.
//
// 		%V Val*		Constant values
//
// 		%S Sym*		Symbols
// 			Flags: +,- #: mode (see below)
// 				"%hS"	unqualified identifier in any mode
// 				"%hhS"  in export mode: unqualified identifier if exported, qualified if not
//
// 		%T Type*	Types
// 			Flags: +,- #: mode (see below)
// 				'l' definition instead of name.
// 				'h' omit "func" and receiver in function types
// 				'u' (only in -/Sym mode) print type identifiers wit package name instead of prefix.
//
// 		%N Node*	Nodes
// 			Flags: +,- #: mode (see below)
// 				'h' (only in +/debug mode) suppress recursion
// 				'l' (only in Error mode) print "foo (type Bar)"
//
// 		%H Nodes	Nodes
// 			Flags: those of %N
// 				','  separate items with ',' instead of ';'
//
// 	  In mparith2.go and mparith3.go:
// 			%B Mpint*	Big integers
// 			%F Mpflt*	Big floats
//
// 	  %S, %T and %N obey use the following flags to set the format mode:

// Format conversions
//
// 	  %L int        Line numbers
//
// 	  %E int        etype values (aka 'Kind')
//
// 	  %O int        Node Opcodes
// 	      Flags: "%#O": print go syntax. (automatic unless fmtmode == FDbg)
//
// 	  %J Node*    Node details
// 	      Flags: "%hJ" suppresses things not relevant until walk.
//
// 	  %V Val*        Constant values
//
// 	  %S Sym*        Symbols
// 	      Flags: +,- #: mode (see below)
// 	          "%hS"    unqualified identifier in any mode
// 	          "%hhS"  in export mode: unqualified identifier if exported, qualified if not
//
// 	  %T Type*    Types
// 	      Flags: +,- #: mode (see below)
// 	          'l' definition instead of name.
// 	          'h' omit "func" and receiver in function types
// 	          'u' (only in -/Sym mode) print type identifiers wit package name instead of prefix.
//
// 	  %N Node*    Nodes
// 	      Flags: +,- #: mode (see below)
// 	          'h' (only in +/debug mode) suppress recursion
// 	          'l' (only in Error mode) print "foo (type Bar)"
//
// 	  %H NodeList*    NodeLists
// 	      Flags: those of %N
// 	          ','  separate items with ',' instead of ';'
//
// 	In mparith2.go and mparith3.go:
// 	      %B Mpint*    Big integers
// 	      %F Mpflt*    Big floats
//
// 	%S, %T and %N obey use the following flags to set the format mode:
const (
	FErr = iota
	FDbg
	FExp
	FTypeId
)

const (
	FmtWidth    FmtFlag = 1 << iota
	FmtLeft             // "-"
	FmtSharp            // "#"
	FmtSign             // "+"
	FmtUnsigned         // "u"
	FmtShort            // "h"
	FmtLong             // "l"
	FmtComma            // ","
	FmtByte             // "hh"
	FmtBody             // for printing export bodies
)

const (
	FunargNone    Funarg = iota
	FunargRcvr           // receiver
	FunargParams         // input parameters
	FunargResults        // output results
)

// FNV-1 hash function constants.
const (
	H0 = 2166136261
	Hp = 16777619
)

// static initialization
const (
	InitNotStarted = 0
	InitDone       = 1
	InitPending    = 2
)

const (
	// names and literals
	LNAME = utf8.RuneSelf + iota
	LLITERAL

	// operator-based operations
	LOPER
	LASOP
	LINCOP

	// miscellaneous
	LCOLAS
	LCOMM
	LDDD

	// keywords
	LBREAK
	LCASE
	LCHAN
	LCONST
	LCONTINUE
	LDEFAULT
	LDEFER
	LELSE
	LFALL
	LFOR
	LFUNC
	LGO
	LGOTO
	LIF
	LIMPORT
	LINTERFACE
	LMAP
	LPACKAGE
	LRANGE
	LRETURN
	LSELECT
	LSTRUCT
	LSWITCH
	LTYPE
	LVAR
	LIGNORE
)

// MaxFlowProg is the maximum size program (counted in instructions) for which
// the flow code will build a graph. Functions larger than this limit will not
// have flow graphs and consequently will not be optimized.
const MaxFlowProg = 50000

// The Plan 9 C compilers used a limit of 600 regions,
// but the yacc-generated parser in y.go has 3100 regions.
// We set MaxRgn large enough to handle that.
// There's not a huge cost to having too many regions:
// the main processing traces the live area for each variable,
// which is limited by the number of variables times the area,
// not the raw region count. If there are many regions, they
// are almost certainly small and easy to trace.
// The only operation that scales with region count is the
// sorting by cost, which uses sort.Sort and is therefore
// guaranteed n log n.
const MaxRgn = 6000

// There appear to be some loops in the escape graph, causing
// arbitrary recursion into deeper and deeper levels.
// Cut this off safely by making minLevel sticky: once you
// get that deep, you cannot go down any further but you also
// cannot go up any further. This is a conservative fix.
// Making minLevel smaller (more negative) would handle more
// complex chains of indirections followed by address-of operations,
// at the cost of repeating the traversal once for each additional
// allowed level when a loop is encountered. Using -2 suffices to
// pass all the tests we have written so far, which we assume matches
// the level of complexity we want the escape analysis code to handle.
const (
	MinLevel = -2
)

const (
	// Maximum size in bits for Mpints before signalling
	// overflow and also mantissa precision for Mpflts.
	Mpprec = 512

	// Turn on for constant arithmetic debugging output.
	Mpdebug = false
)

const NOTALOOPDEPTH = -1

const (
	Nointerface       Pragma = 1 << iota
	Noescape                 // func parameters don't escape
	Norace                   // func must not have race detector annotations
	Nosplit                  // func should not execute on separate stack
	Noinline                 // func should not be inlined
	Systemstack              // func must run on system stack
	Nowritebarrier           // emit compiler error instead of write barrier
	Nowritebarrierrec        // error on write barrier in this or recursive callees
	CgoUnsafeArgs            // treat a pointer to one arg as a pointer to them all
	UintptrEscapes           // pointers converted to uintptr escape
)

// Node ops.
const (
	OXXX = Op(iota)

	// names
	ONAME    // var, const or func name
	ONONAME  // unnamed arg or return value: f(int, string) (int, error) { etc }
	OTYPE    // type name
	OPACK    // import
	OLITERAL // literal

	// expressions
	OADD             // Left + Right
	OSUB             // Left - Right
	OOR              // Left | Right
	OXOR             // Left ^ Right
	OADDSTR          // +{List} (string addition, list elements are strings)
	OADDR            // &Left
	OANDAND          // Left && Right
	OAPPEND          // append(List)
	OARRAYBYTESTR    // Type(Left) (Type is string, Left is a []byte)
	OARRAYBYTESTRTMP // Type(Left) (Type is string, Left is a []byte, ephemeral)
	OARRAYRUNESTR    // Type(Left) (Type is string, Left is a []rune)
	OSTRARRAYBYTE    // Type(Left) (Type is []byte, Left is a string)
	OSTRARRAYBYTETMP // Type(Left) (Type is []byte, Left is a string, ephemeral)
	OSTRARRAYRUNE    // Type(Left) (Type is []rune, Left is a string)
	OAS              // Left = Right or (if Colas=true) Left := Right
	OAS2             // List = Rlist (x, y, z = a, b, c)
	OAS2FUNC         // List = Rlist (x, y = f())
	OAS2RECV         // List = Rlist (x, ok = <-c)
	OAS2MAPR         // List = Rlist (x, ok = m["foo"])
	OAS2DOTTYPE      // List = Rlist (x, ok = I.(int))
	OASOP            // Left Etype= Right (x += y)
	OASWB            // Left = Right (with write barrier)
	OCALL            // Left(List) (function call, method call or type conversion)
	OCALLFUNC        // Left(List) (function call f(args))
	OCALLMETH        // Left(List) (direct method call x.Method(args))
	OCALLINTER       // Left(List) (interface method call x.Method(args))
	OCALLPART        // Left.Right (method expression x.Method, not called)
	OCAP             // cap(Left)
	OCLOSE           // close(Left)
	OCLOSURE         // func Type { Body } (func literal)
	OCMPIFACE        // Left Etype Right (interface comparison, x == y or x != y)
	OCMPSTR          // Left Etype Right (string comparison, x == y, x < y, etc)
	OCOMPLIT         // Right{List} (composite literal, not yet lowered to specific form)
	OMAPLIT          // Type{List} (composite literal, Type is map)
	OSTRUCTLIT       // Type{List} (composite literal, Type is struct)
	OARRAYLIT        // Type{List} (composite literal, Type is array or slice)
	OPTRLIT          // &Left (left is composite literal)
	OCONV            // Type(Left) (type conversion)
	OCONVIFACE       // Type(Left) (type conversion, to interface)
	OCONVNOP         // Type(Left) (type conversion, no effect)
	OCOPY            // copy(Left, Right)
	ODCL             // var Left (declares Left of type Left.Type)

	// Used during parsing but don't last.
	ODCLFUNC   // func f() or func (r) f()
	ODCLFIELD  // struct field, interface field, or func/method argument/return value.
	ODCLCONST  // const pi = 3.14
	ODCLTYPE   // type Int int
	ODELETE    // delete(Left, Right)
	ODOT       // Left.Sym (Left is of struct type)
	ODOTPTR    // Left.Sym (Left is of pointer to struct type)
	ODOTMETH   // Left.Sym (Left is non-interface, Right is method name)
	ODOTINTER  // Left.Sym (Left is interface, Right is method name)
	OXDOT      // Left.Sym (before rewrite to one of the preceding)
	ODOTTYPE   // Left.Right or Left.Type (.Right during parsing, .Type once resolved)
	ODOTTYPE2  // Left.Right or Left.Type (.Right during parsing, .Type once resolved; on rhs of OAS2DOTTYPE)
	OEQ        // Left == Right
	ONE        // Left != Right
	OLT        // Left < Right
	OLE        // Left <= Right
	OGE        // Left >= Right
	OGT        // Left > Right
	OIND       // *Left
	OINDEX     // Left[Right] (index of array or slice)
	OINDEXMAP  // Left[Right] (index of map)
	OKEY       // Left:Right (key:value in struct/array/map literal, or slice index pair)
	OLEN       // len(Left)
	OMAKE      // make(List) (before type checking converts to one of the following)
	OMAKECHAN  // make(Type, Left) (type is chan)
	OMAKEMAP   // make(Type, Left) (type is map)
	OMAKESLICE // make(Type, Left, Right) (type is slice)
	OMUL       // Left * Right
	ODIV       // Left / Right
	OMOD       // Left % Right
	OLSH       // Left << Right
	ORSH       // Left >> Right
	OAND       // Left & Right
	OANDNOT    // Left &^ Right
	ONEW       // new(Left)
	ONOT       // !Left
	OCOM       // ^Left
	OPLUS      // +Left
	OMINUS     // -Left
	OOROR      // Left || Right
	OPANIC     // panic(Left)
	OPRINT     // print(List)
	OPRINTN    // println(List)
	OPAREN     // (Left)
	OSEND      // Left <- Right
	OSLICE     // Left[Right.Left : Right.Right] (Left is untypechecked or slice; Right.Op==OKEY)
	OSLICEARR  // Left[Right.Left : Right.Right] (Left is array)
	OSLICESTR  // Left[Right.Left : Right.Right] (Left is string)
	OSLICE3    // Left[R.Left : R.R.Left : R.R.R] (R=Right; Left is untypedchecked or slice; R.Op and R.R.Op==OKEY)
	OSLICE3ARR // Left[R.Left : R.R.Left : R.R.R] (R=Right; Left is array; R.Op and R.R.Op==OKEY)
	ORECOVER   // recover()
	ORECV      // <-Left
	ORUNESTR   // Type(Left) (Type is string, Left is rune)
	OSELRECV   // Left = <-Right.Left: (appears as .Left of OCASE; Right.Op == ORECV)
	OSELRECV2  // List = <-Right.Left: (apperas as .Left of OCASE; count(List) == 2, Right.Op == ORECV)
	OIOTA      // iota
	OREAL      // real(Left)
	OIMAG      // imag(Left)
	OCOMPLEX   // complex(Left, Right)

	// statements
	OBLOCK    // { List } (block of code)
	OBREAK    // break
	OCASE     // case List: Nbody (select case after processing; List==nil means default)
	OXCASE    // case List: Nbody (select case before processing; List==nil means default)
	OCONTINUE // continue
	ODEFER    // defer Left (Left must be call)
	OEMPTY    // no-op (empty statement)
	OFALL     // fallthrough (after processing)
	OXFALL    // fallthrough (before processing)
	OFOR      // for Ninit; Left; Right { Nbody }
	OGOTO     // goto Left
	OIF       // if Ninit; Left { Nbody } else { Rlist }
	OLABEL    // Left:
	OPROC     // go Left (Left must be call)
	ORANGE    // for List = range Right { Nbody }
	ORETURN   // return List
	OSELECT   // select { List } (List is list of OXCASE or OCASE)
	OSWITCH   // switch Ninit; Left { List } (List is a list of OXCASE or OCASE)
	OTYPESW   // List = Left.(type) (appears as .Left of OSWITCH)

	// types
	OTCHAN   // chan int
	OTMAP    // map[string]int
	OTSTRUCT // struct{}
	OTINTER  // interface{}
	OTFUNC   // func()
	OTARRAY  // []int, [8]int, [N]int or [...]int

	// misc
	ODDD        // func f(args ...int) or f(l...) or var a = [...]int{0, 1, 2}.
	ODDDARG     // func f(args ...int), introduced by escape analysis.
	OINLCALL    // intermediary representation of an inlined call.
	OEFACE      // itable and data words of an empty-interface value.
	OITAB       // itable word of an interface value.
	OSPTR       // base pointer of a slice or string.
	OCLOSUREVAR // variable reference at beginning of closure function
	OCFUNC      // reference to c function pointer (not go func value)
	OCHECKNIL   // emit code to ensure pointer/interface not nil
	OVARKILL    // variable is dead
	OVARLIVE    // variable is alive

	// thearch-specific registers
	OREGISTER // a register, such as AX.
	OINDREG   // offset plus indirect of a register, such as 8(SP).

	// arch-specific opcodes
	OCMP    // compare: ACMP.
	ODEC    // decrement: ADEC.
	OINC    // increment: AINC.
	OEXTEND // extend: ACWD/ACDQ/ACQO.
	OHMUL   // high mul: AMUL/AIMUL for unsigned/signed (OMUL uses AIMUL for both).
	OLROT   // left rotate: AROL.
	ORROTC  // right rotate-carry: ARCR.
	ORETJMP // return to other function
	OPS     // compare parity set (for x86 NaN check)
	OPC     // compare parity clear (for x86 NaN check)
	OSQRT   // sqrt(float64), on systems that have hw support
	OGETG   // runtime.getg() (read g pointer)
	OEND
)

const (
	// Precedences of binary operators (must be > 0).
	PCOMM OpPrec = 1 + iota
	POROR
	PANDAND
	PCMP
	PADD
	PMUL
)

const (
	// Pseudo-op, like TEXT, GLOBL, TYPE, PCDATA, FUNCDATA.
	Pseudo = 1 << 1

	// There's nothing to say about the instruction,
	// but it's still okay to see.
	OK = 1 << 2

	// Size of right-side write, or right-side read if no write.
	SizeB = 1 << 3
	SizeW = 1 << 4
	SizeL = 1 << 5
	SizeQ = 1 << 6
	SizeF = 1 << 7
	SizeD = 1 << 8

	// Left side (Prog.from): address taken, read, write.
	LeftAddr  = 1 << 9
	LeftRead  = 1 << 10
	LeftWrite = 1 << 11

	// Register in middle (Prog.reg); only ever read. (arm, ppc64)
	RegRead    = 1 << 12
	CanRegRead = 1 << 13

	// Right side (Prog.to): address taken, read, write.
	RightAddr  = 1 << 14
	RightRead  = 1 << 15
	RightWrite = 1 << 16

	// Instruction kinds
	Move  = 1 << 17 // straight move
	Conv  = 1 << 18 // size conversion
	Cjmp  = 1 << 19 // conditional jump
	Break = 1 << 20 // breaks control flow (no fallthrough)
	Call  = 1 << 21 // function call
	Jump  = 1 << 22 // jump
	Skip  = 1 << 23 // data instruction

	// Set, use, or kill of carry bit. Kill means we never look at the carry bit
	// after this kind of instruction. Originally for understanding ADC, RCR,
	// and so on, but now also tracks set, use, and kill of the zero and
	// overflow bits as well. TODO rename to {Set,Use,Kill}Flags

	// Set, use, or kill of carry bit. Kill means we never look at the carry bit
	// after this kind of instruction.
	SetCarry  = 1 << 24
	UseCarry  = 1 << 25
	KillCarry = 1 << 26

	// Special cases for register use. (amd64, 386)
	ShiftCX  = 1 << 27 // possible shift by CX
	ImulAXDX = 1 << 28 // possible multiply into DX:AX

	// Instruction updates whichever of from/to is type D_OREG. (ppc64)
	PostInc = 1 << 29
)

const (
	Pxxx      Class = iota
	PEXTERN         // global variable
	PAUTO           // local variables
	PAUTOHEAP       // local variable or parameter moved to heap
	PPARAM          // input arguments
	PPARAMOUT       // output results
	PFUNC           // global function
	PDISCARD        // discard during parse of duplicate import
)

const (
	SymExport SymFlags = 1 << iota // to be exported
	SymPackage
	SymExported // already written out by export
	SymUniq
	SymSiggen
	SymAsm
	SymAlgGen
)

const (
	Txxx = iota
	TINT8
	TUINT8
	TINT16
	TUINT16
	TINT32
	TUINT32
	TINT64
	TUINT64
	TINT
	TUINT
	TUINTPTR
	TCOMPLEX64
	TCOMPLEX128
	TFLOAT32
	TFLOAT64
	TBOOL
	TPTR32
	TPTR64
	TFUNC
	TSLICE
	TARRAY
	TSTRUCT
	TCHAN
	TMAP
	TINTER
	TFORW
	TANY
	TSTRING
	TUNSAFEPTR

	// pseudo-types for literals
	TIDEAL
	TNIL
	TBLANK

	// pseudo-types for frame layout

	// pseudo-type for frame layout
	TFUNCARGS
	TCHANARGS
	TINTERMETH

	// pseudo-types for import/export
	TDDDFIELD // wrapper: contained type is a ... field
	NTYPE
)

const (
	UINF            = 100
	BADWIDTH        = -1000000000
	MaxStackVarSize = 10 * 1024 * 1024
)

const (
	UNVISITED = 0
	VISITED   = 1
)

const (
	WORDBITS  = 32
	WORDMASK  = WORDBITS - 1
	WORDSHIFT = 5
)

// note this is the runtime representation
// of the compilers arrays.
//
// typedef	struct
// {					// must not move anything
// 	uchar	array[8];	// pointer to data
// 	uchar	nel[4];		// number of elements
// 	uchar	cap[4];		// allocated number of elements
// } Array;

// note this is the runtime representation
// of the compilers arrays.
//
// typedef    struct
// {                    // must not move anything
//     uchar    array[8];    // pointer to data
//     uchar    nel[4];        // number of elements
//     uchar    cap[4];        // allocated number of elements
// } Array;
var Array_array int // runtime offsetof(Array,array) - same for String

var Array_cap int // runtime offsetof(Array,cap)

var Array_nel int // runtime offsetof(Array,nel) - same for String

var Ctxt *obj.Link

var Curfn *Node

var Debug [256]int

var Debug_checknil int

var (
	Debug_export int // if set, print debugging information about export data
)

var Debug_gcprog int // set by -d gcprog

var Debug_typeassert int

var Deferproc *Node

var Deferreturn *Node

var Disable_checknil int

var Funcdepth int32 // len(funcstack) during parsing, but then forced to be the same later during compilation

var (
	Isint     [NTYPE]bool
	Isfloat   [NTYPE]bool
	Iscomplex [NTYPE]bool
)

var Maxarg int64

var Maxintval [NTYPE]*Mpint

var Minintval [NTYPE]*Mpint

var Nacl bool

var Newproc *Node

var Ostats OptStats

var Panicindex *Node

var Pc *obj.Prog

var Runtimepkg *Pkg // package runtime

var Simtype [NTYPE]EType

var Stksize int64 // stack size for current frame

var Thearch Arch

var Tptr EType // either TPTR32 or TPTR64

// Types stores pointers to predeclared named types.
//
// It also stores pointers to several special types:
//
// 	- Types[TANY] is the placeholder "any" type recognized by substArgTypes.
// 	- Types[TBLANK] represents the blank variable's type.
// 	- Types[TIDEAL] represents untyped numeric constants.
// 	- Types[TNIL] represents the predeclared "nil" value's type.
// 	- Types[TUNSAFEPTR] is package unsafe's Pointer type.
var Types [NTYPE]*Type

var Widthint int

var Widthptr int

var Widthreg int

// AlgKind describes the kind of algorithms used for comparing and
// hashing a Type.
type AlgKind int

type Arch struct {
	LinkArch            *obj.LinkArch
	REGSP               int
	REGCTXT             int
	REGCALLX            int // BX
	REGCALLX2           int // AX
	REGRETURN           int // AX
	REGMIN              int
	REGMAX              int
	REGZERO             int // architectural zero register, if available
	FREGMIN             int
	FREGMAX             int
	MAXWIDTH            int64
	ReservedRegs        []int
	AddIndex            func(*Node, int64, *Node) bool // optional
	Betypeinit          func()
	Bgen_float          func(*Node, bool, int, *obj.Prog) // optional
	Cgen64              func(*Node, *Node)                // only on 32-bit systems
	Cgenindex           func(*Node, *Node, bool) *obj.Prog
	Cgen_bmul           func(Op, *Node, *Node, *Node) bool
	Cgen_float          func(*Node, *Node) // optional
	Cgen_hmul           func(*Node, *Node, *Node)
	RightShiftWithCarry func(*Node, uint, *Node)  // only on systems without RROTC instruction
	AddSetCarry         func(*Node, *Node, *Node) // only on systems when ADD does not update carry flag
	Cgen_shift          func(Op, bool, *Node, *Node, *Node)
	Clearfat            func(*Node)
	Cmp64               func(*Node, *Node, Op, int, *obj.Prog) // only on 32-bit systems
	Defframe            func(*obj.Prog)
	Dodiv               func(Op, *Node, *Node, *Node)
	Excise              func(*Flow)
	Expandchecks        func(*obj.Prog)
	Getg                func(*Node)
	Gins                func(obj.As, *Node, *Node) *obj.Prog

	// Ginscmp generates code comparing n1 to n2 and jumping away if op is
	// satisfied. The returned prog should be Patch'ed with the jump target. If
	// op is not satisfied, code falls through to the next emitted instruction.
	// Likely is the branch prediction hint: +1 for likely, -1 for unlikely, 0
	// for no opinion.
	//
	// Ginscmp must be able to handle all kinds of arguments for n1 and n2, not
	// just simple registers, although it can assume that there are no function
	// calls needed during the evaluation, and on 32-bit systems the values are
	// guaranteed not to be 64-bit values, so no in-memory temporaries are
	// necessary.
	Ginscmp func(op Op, t *Type, n1, n2 *Node, likely int) *obj.Prog

	// Ginsboolval inserts instructions to convert the result
	// of a just-completed comparison to a boolean value.
	// The first argument is the conditional jump instruction
	// corresponding to the desired value.
	// The second argument is the destination.
	// If not present, Ginsboolval will be emulated with jumps.
	Ginsboolval  func(obj.As, *Node)
	Ginscon      func(obj.As, int64, *Node)
	Ginsnop      func()
	Gmove        func(*Node, *Node)
	Igenindex    func(*Node, *Node, bool) *obj.Prog
	Peep         func(*obj.Prog)
	Proginfo     func(*obj.Prog) // fills in Prog.Info
	Regtyp       func(*obj.Addr) bool
	Sameaddr     func(*obj.Addr, *obj.Addr) bool
	Smallindir   func(*obj.Addr, *obj.Addr) bool
	Stackaddr    func(*obj.Addr) bool
	Blockcopy    func(*Node, *Node, int64, int64, int64)
	Sudoaddable  func(obj.As, *Node, *obj.Addr) bool
	Sudoclean    func()
	Excludedregs func() uint64
	RtoB         func(int) uint64
	FtoB         func(int) uint64
	BtoR         func(uint64) int
	BtoF         func(uint64) int
	Optoas       func(Op, *Type) obj.As
	Doregbits    func(int) uint64
	Regnames     func(*int) []string
	Use387       bool // should 8g use 387 FP instructions instead of sse2.

	// SSARegToReg maps ssa register numbers to obj register numbers.
	SSARegToReg []int16

	// SSAMarkMoves marks any MOVXconst ops that need to avoid clobbering flags.
	SSAMarkMoves func(*SSAGenState, *ssa.Block)

	// SSAGenValue emits Prog(s) for the Value.
	SSAGenValue func(*SSAGenState, *ssa.Value)

	// SSAGenBlock emits end-of-block Progs. SSAGenValue should be called
	// for all values in the block before SSAGenBlock.
	SSAGenBlock func(s *SSAGenState, b, next *ssa.Block)
}

// ArrayType contains Type fields specific to array types.
type ArrayType struct {
	Elem        *Type // element type
	Bound       int64 // number of elements; <0 if unknown yet
	Haspointers uint8 // 0 unknown, 1 no, 2 yes
}

// An ordinary basic block.
//
// Instructions are threaded together in a doubly-linked list. To iterate in
// program order follow the link pointer from the first node and stop after the
// last node has been visited
//
// 	for p = bb.first; ; p = p.link {
// 	  ...
// 	  if p == bb.last {
// 	    break
// 	  }
// 	}
//
// To iterate in reverse program order by following the opt pointer from the
// last node
//
// 	for p = bb.last; p != nil; p = p.opt {
// 	  ...
// 	}

// An ordinary basic block.
//
// Instructions are threaded together in a doubly-linked list. To iterate in
// program order follow the link pointer from the first node and stop after the
// last node has been visited
//
// 	for(p = bb->first;; p = p->link) {
// 	  ...
// 	  if(p == bb->last)
// 	    break;
// 	}
//
// To iterate in reverse program order by following the opt pointer from the
// last node
//
// 	for(p = bb->last; p != nil; p = p->opt) {
// 	  ...
// 	}
type BasicBlock struct {
}

// Bits represents a set of Vars, stored as a bit set of var numbers
// (the index in vars, or equivalently v.id).
type Bits struct {
}

// Branch is an unresolved branch.
type Branch struct {
	P *obj.Prog  // branch instruction
	B *ssa.Block // target
}

// ChanArgsType contains Type fields specific to TCHANARGS types.
type ChanArgsType struct {
	T *Type // reference to a chan type whose elements need a width check
}

// ChanDir is whether a channel can send, receive, or both.
type ChanDir uint8

// ChanType contains Type fields specific to channel types.
type ChanType struct {
	Elem *Type   // element type
	Dir  ChanDir // channel direction
}

// The Class of a variable/function describes the "storage class"
// of a variable or function. During parsing, storage classes are
// called declaration contexts.
type Class uint8

// Ctype describes the constant kind of an "ideal" (untyped) constant.
type Ctype int8

// DDDFieldType contains Type fields specific to TDDDFIELD types.
type DDDFieldType struct {
	T *Type // reference to a slice type for ... args
}

// A Dlist stores a pointer to a TFIELD Type embedded within
// a TSTRUCT or TINTER Type.
type Dlist struct {
}

// EType describes a kind of type.
type EType uint8

type Error struct {
}

type EscState struct {
}

// An EscStep documents one step in the path from memory
// that is heap allocated to the (alleged) reason for the
// heap allocation.
type EscStep struct {
}

// A Field represents a field in a struct or a method in an interface or
// associated with a named type.
type Field struct {
	Nointerface bool
	Embedded    uint8 // embedded field
	Funarg      Funarg
	Broke       bool // broken field definition
	Isddd       bool // field is ... argument
	Sym         *Sym
	Nname       *Node
	Type        *Type // field type

	// Offset in bytes of this field or method within its enclosing struct
	// or interface Type.
	Offset int64
	Note   string // literal string annotation
}

// Fields is a pointer to a slice of *Field.
// This saves space in Types that do not have fields or methods
// compared to a simple slice of *Field.
type Fields struct {
}

type FloatingEQNEJump struct {
	Jump  obj.As
	Index int
}

type Flow struct {
	Prog   *obj.Prog // actual instruction
	P1     *Flow     // predecessors of this instruction: p1,
	P2     *Flow     // and then p2 linked though p2link.
	P2link *Flow
	S1     *Flow // successors of this instruction (at most two: s1 and s2).
	S2     *Flow
	Link   *Flow       // next instruction in function code
	Active int32       // usable by client
	Id     int32       // sequence number in flow graph
	Rpo    int32       // reverse post ordering
	Loop   uint16      // x5 for every loop
	Refset bool        // diagnostic generated
	Data   interface{} // for use by client
}

// A FmtFlag value is a set of flags (or 0).
// They control how the Xconv functions format their values.
// See the respective function's documentation for details.
type FmtFlag int

// ForwardType contains Type fields specific to forward types.
type ForwardType struct {
	Copyto      []*Node // where to copy the eventual value to
	Embedlineno int32   // first use of this type as an embedded type
}

// Fnstruct records the kind of function argument
type Funarg uint8

// Func holds Node fields used only with function-like nodes.
type Func struct {
	Shortname     *Node
	Enter         Nodes // for example, allocate and initialize memory for escaping parameters
	Exit          Nodes
	Cvars         Nodes   // closure params
	Dcl           []*Node // autodcl for this func/closure
	Inldcl        Nodes   // copy of dcl for use in inlining
	Closgen       int
	Outerfunc     *Node // outer function (for closure)
	FieldTrack    map[*Sym]struct{}
	Ntype         *Node // signature
	Top           int   // top context (Ecall, Eproc, etc)
	Closure       *Node // OCLOSURE <-> ODCLFUNC
	FCurfn        *Node
	Nname         *Node
	Inl           Nodes // copy of the body for use in inlining
	InlCost       int32
	Depth         int32
	Endlineno     int32
	WBLineno      int32  // line number of first write barrier
	Pragma        Pragma // go:xxx function annotations
	Dupok         bool   // duplicate definitions ok
	Wrapper       bool   // is method wrapper
	Needctxt      bool   // function uses context register (has closure variables)
	ReflectMethod bool   // function calls reflect.Type.Method or MethodByName
}

// // FuncArgsType contains Type fields specific to TFUNCARGS types.
type FuncArgsType struct {
	T *Type // reference to a func type whose elements need a width check
}

// FuncType contains Type fields specific to func types.
type FuncType struct {
	Receiver *Type // function receiver
	Results  *Type // function results
	Params   *Type // function params
	Nname    *Node

	// Argwid is the total width of the function receiver, params, and results.
	// It gets calculated via a temporary TFUNCARGS type.
	// Note that TFUNC's Width is Widthptr.
	Argwid   int64
	Outnamed bool
}

type GCProg struct {
}

type Graph struct {
	Start *Flow
	Num   int

	// After calling flowrpo, rpo lists the flow nodes in reverse postorder,
	// and each non-dead Flow node f has g->rpo[f->rpo] == f.
	Rpo []*Flow
}

type InitEntry struct {
	Xoffset int64 // struct, array only
	Expr    *Node // bytes of run-time computed expressions
}

type InitPlan struct {
	E []InitEntry
}

// InterMethType contains Type fields specific to interface method psuedo-types.
type InterMethType struct {
	Nname *Node
}

// InterType contains Type fields specific to interface types.
type InterType struct {
}

// Iter provides an abstraction for iterating across struct fields and
// interface methods.
type Iter struct {
}

type Label struct {
	Sym *Sym
	Def *Node
	Use []*Node

	// for use during gen
	Gotopc   *obj.Prog // pointer to unresolved gotos
	Labelpc  *obj.Prog // pointer to code
	Breakpc  *obj.Prog // pointer to code
	Continpc *obj.Prog // pointer to code
	Used     bool
}

// A Level encodes the reference state and context applied to (stack, heap)
// allocated memory.
//
// value is the overall sum of *(1) and &(-1) operations encountered along a
// path from a destination (sink, return value) to a source (allocation,
// parameter).
//
// suffixValue is the maximum-copy-started-suffix-level applied to a sink. For
// example: sink = x.left.left --> level=2, x is dereferenced twice and does not
// escape to sink. sink = &Node{x} --> level=-1, x is accessible from sink via
// one "address of" sink = &Node{&Node{x}} --> level=-2, x is accessible from
// sink via two "address of" sink = &Node{&Node{x.left}} --> level=-1, but x is
// NOT accessible from sink because it was indirected and then copied. (The copy
// operations are sometimes implicit in the source code; in this case, value of
// x.left was copied into a field of a newly allocated Node)
//
// There's one of these for each Node, and the integer values rarely exceed even
// what can be stored in 4 bits, never mind 8.
type Level struct {
}

// A collection of global state used by liveness analysis.
type Liveness struct {
}

// argument passing to/from
// smagic and umagic
type Magic struct {
	W   int // input for both - width
	S   int // output for both - shift
	Bad int // output for both - unexpected failure

	// magic multiplier for signed literal divisors
	Sd int64 // input - literal divisor
	Sm int64 // output - multiplier

	// magic multiplier for unsigned literal divisors
	Ud uint64 // input - literal divisor
	Um uint64 // output - multiplier
	Ua int    // output - adder
}

// MapType contains Type fields specific to maps.
type MapType struct {
	Key    *Type // Key type
	Val    *Type // Val (elem) type
	Bucket *Type // internal struct type representing a hash bucket
	Hmap   *Type // internal struct type representing the Hmap (map header object)
	Hiter  *Type // internal struct type representing hash iterator state
}

// Mpcplx represents a complex constant.
type Mpcplx struct {
	Real Mpflt
	Imag Mpflt
}

// Mpflt represents a floating-point constant.
type Mpflt struct {
	Val big.Float
}

// Mpint represents an integer constant.
type Mpint struct {
	Val  big.Int
	Ovf  bool // set if Val overflowed compiler limit (sticky)
	Rune bool // set if syntax indicates default type rune
}

// Name holds Node fields used only by named nodes (ONAME, OPACK, OLABEL,
// ODCLFIELD, some OLITERAL).

// Name holds Node fields used only by named nodes (ONAME, OPACK, some
// OLITERAL).
type Name struct {
	Pack      *Node  // real package for import . names
	Pkg       *Pkg   // pkg for OPACK nodes
	Heapaddr  *Node  // temp holding heap address of param (could move to Param?)
	Inlvar    *Node  // ONAME substitute while inlining (could move to Param?)
	Defn      *Node  // initializing assignment
	Curfn     *Node  // function for local variables
	Param     *Param // additional fields for ONAME, ODCLFIELD
	Decldepth int32  // declaration loop depth, increased for every loop or label
	Vargen    int32  // unique name for ONAME within a function.  Function outputs are numbered starting at one.
	Iota      int32  // value if this name is iota
	Funcdepth int32
	Method    bool // OCALLMETH name
	Readonly  bool
	Captured  bool // is the variable captured by a closure
	Byval     bool // is the variable captured by value or by reference
	Needzero  bool // if it contains pointers, needs to be zeroed on function entry
	Keepalive bool // mark value live across unknown assembly call
}

type NilVal struct {
}

// A Node is a single node in the syntax tree.
// Actually the syntax tree is a syntax DAG, because there is only one
// node with Op=ONAME for a given instance of a variable x.
// The same is true for Op=OTYPE and Op=OLITERAL.
type Node struct {
	// Tree structure.
	// Generic recursive walks should follow these fields.
	Left  *Node
	Right *Node
	Ninit Nodes
	Nbody Nodes
	List  Nodes
	Rlist Nodes

	// most nodes
	Type *Type
	Orig *Node // original form, for printing, and tracking copies of ONAMEs

	// func
	Func *Func

	// ONAME
	Name *Name
	Sym  *Sym        // various
	E    interface{} // Opt or Val, see methods below

	// Various. Usually an offset into a struct. For example, ONAME nodes
	// that refer to local variables use it to identify their stack frame
	// position. ODOT, ODOTPTR, and OINDREG use it to indicate offset
	// relative to their base address. ONAME nodes on the left side of an
	// OKEY within an OSTRUCTLIT use it to store the named field's offset.
	// OXCASE and OXFALL use it to validate the use of fallthrough.
	// Possibly still more uses. If you find any, document them.
	Xoffset int64
	Lineno  int32

	// OREGISTER, OINDREG
	Reg       int16
	Esc       uint16 // EscXXX
	Op        Op
	Ullman    uint8 // sethi/ullman number
	Addable   bool  // addressable
	Etype     EType // op for OASOP, etype for OTYPE, exclam for export, 6g saved reg, ChanDir for OTCHAN
	Bounded   bool  // bounds check unnecessary
	NonNil    bool  // guaranteed to be non-nil
	Class     Class // PPARAM, PAUTO, PEXTERN, etc
	Embedded  uint8 // ODCLFIELD embedded type
	Colas     bool  // OAS resulting from :=
	Diag      uint8 // already printed error about this
	Noescape  bool  // func arguments do not escape; TODO(rsc): move Noescape to Func struct (see CL 7360)
	Walkdef   uint8
	Typecheck uint8
	Local     bool
	Dodata    uint8
	Initorder uint8
	Used      bool
	Isddd     bool // is the argument variadic
	Implicit  bool
	Addrtaken bool // address taken, even if not moved to heap
	Assigned  bool // is the variable ever assigned to
	Likely    int8 // likeliness of if statement
}

type NodeEscState struct {
	Curfn             *Node
	Escflowsrc        []EscStep // flow(this, src)
	Escretval         Nodes     // on OCALLxxx, list of dummy return values
	Escloopdepth      int32     // -1: global, 0: return variables, 1:function top level, increased inside function for every loop or label to mark scopes
	Esclevel          Level
	Walkgen           uint32
	Maxextraloopdepth int32
}

// Nodes is a pointer to a slice of *Node.
// For fields that are not used in most nodes, this is used instead of
// a slice to save space.
type Nodes struct {
}

type Op uint8

type OpPrec int

type OptStats struct {
	Ncvtreg int32
	Nspill  int32
	Nreload int32
	Ndelmov int32
	Nvar    int32
	Naddr   int32
}

// Order holds state during the ordering process.
type Order struct {
}

type Param struct {
	Ntype *Node

	// ONAME PAUTOHEAP
	Stackcopy *Node // the PPARAM/PPARAMOUT on-stack slot (moved func params only)

	// ONAME PPARAM
	Field *Field // TFIELD in arg struct

	// ONAME closure linkage Consider:
	//
	// 	func f() {
	// 		x := 1 // x1
	// 		func() {
	// 			use(x) // x2
	// 			func() {
	// 				use(x) // x3
	// 				--- parser is here ---
	// 			}()
	// 		}()
	// 	}
	//
	// There is an original declaration of x and then a chain of mentions of x
	// leading into the current function. Each time x is mentioned in a new
	// closure, we create a variable representing x for use in that specific
	// closure, since the way you get to x is different in each closure.
	//
	// Let's number the specific variables as shown in the code: x1 is the
	// original x, x2 is when mentioned in the closure, and x3 is when mentioned
	// in the closure in the closure.
	//
	// We keep these linked (assume N > 1):
	//
	// 	- x1.Defn = original declaration statement for x (like most variables)
	// 	- x1.Innermost = current innermost closure x (in this case x3), or nil for none
	// 	- x1.isClosureVar() = false
	//
	// 	- xN.Defn = x1, N > 1
	// 	- xN.isClosureVar() = true, N > 1
	// 	- x2.Outer = nil
	// 	- xN.Outer = x(N-1), N > 2
	//
	// When we look up x in the symbol table, we always get x1. Then we can use
	// x1.Innermost (if not nil) to get the x for the innermost known closure
	// function, but the first reference in a closure will find either no
	// x1.Innermost or an x1.Innermost with .Funcdepth < Funcdepth. In that
	// case, a new xN must be created, linked in with:
	//
	// 	xN.Defn = x1
	// 	xN.Outer = x1.Innermost
	// 	x1.Innermost = xN
	//
	// When we finish the function, we'll process its closure variables and find
	// xN and pop it off the list using:
	//
	// 	x1 := xN.Defn
	// 	x1.Innermost = xN.Outer
	//
	// We leave xN.Innermost set so that we can still get to the original
	// variable quickly. Not shown here, but once we're done parsing a function
	// and no longer need xN.Outer for the lexical x reference links as
	// described above, closurebody recomputes xN.Outer as the semantic x
	// reference link tree, even filling in x in intermediate closures that
	// might not have mentioned it along the way to inner closures that did. See
	// closurebody for details.
	//
	// During the eventual compilation, then, for closure variables we have:
	//
	// 	xN.Defn = original variable
	// 	xN.Outer = variable captured in next outward scope
	// 	           to make closure where xN appears
	//
	// Because of the sharding of pieces of the node, x.Defn means x.Name.Defn
	// and x.Innermost/Outer means x.Name.Param.Innermost/Outer.
	Innermost *Node
	Outer     *Node
}

type Pkg struct {
	Name     string // package name, e.g. "sys"
	Path     string // string literal used in import statement, e.g. "runtime/internal/sys"
	Pathsym  *obj.LSym
	Prefix   string // escaped path for use in symbol table
	Imported bool   // export data of this package was parsed
	Exported bool   // import line written in export data
	Direct   bool   // imported directly
	Safe     bool   // whether the package is marked as safe
	Syms     map[string]*Sym
}

type Pragma uint16

// PtrType contains Type fields specific to pointer types.
type PtrType struct {
	Elem *Type // element type
}

// A Reg is a wrapper around a single Prog (one instruction) that holds
// register optimization information while the optimizer runs.
// r->prog is the instruction.
type Reg struct {
}

// A Rgn represents a single regopt variable over a region of code
// where a register could potentially be dedicated to that variable.
// The code encompassed by a Rgn is defined by the flow graph,
// starting at enter, flood-filling forward while varno is refahead
// and backward while varno is refbehind, and following branches.
// A single variable may be represented by multiple disjoint Rgns and
// each Rgn may choose a different register for that variable.
// Registers are allocated to regions greedily in order of descending
// cost.
type Rgn struct {
}

// SSAGenState contains state needed during Prog generation.
type SSAGenState struct {
	// Branches remembers all the branch instructions we've seen
	// and where they would like to go.
	Branches []Branch
}

type Sig struct {
}

// SliceType contains Type fields specific to slice types.
type SliceType struct {
	Elem *Type // element type
}

// StructType contains Type fields specific to struct types.
type StructType struct {
	// Maps have three associated internal structs (see struct MapType).
	// Map links such structs back to their map type.
	Map         *Type
	Funarg      Funarg // type of function arguments for arg struct
	Haspointers uint8  // 0 unknown, 1 no, 2 yes
}

// Sym represents an object name. Most commonly, this is a Go identifier naming
// an object declared within a package, but Syms are also used to name internal
// synthesized objects.
//
// As a special exception, field and method names that are exported use the Sym
// associated with localpkg instead of the package that declared them. This
// allows using Sym pointer equality to test for Go identifier uniqueness when
// handling selector expressions.
type Sym struct {
	Flags     SymFlags
	Link      *Sym
	Importdef *Pkg   // where imported definition was found
	Linkname  string // link name

	// saved and restored by dcopy
	Pkg        *Pkg
	Name       string // variable name
	Def        *Node  // definition: ONAME OTYPE OPACK or OLITERAL
	Block      int32  // blocknumber to catch redeclaration
	Lastlineno int32  // last declaration for diagnostic
	Label      *Label // corresponding label (ephemeral)
	Origpkg    *Pkg   // original package for . import
	Lsym       *obj.LSym
	Fsym       *Sym // funcsym
}

type SymFlags uint8

// code to help generate trampoline
// functions for methods on embedded
// subtypes.
// these are approx the same as
// the corresponding adddot routines
// except that they expect to be called
// with unique tasks and they return
// the actual methods.
type Symlink struct {
}

type TempVar struct {
}

// A Type represents a Go type.
type Type struct {
	// Extra contains extra etype-specific fields. As an optimization, those
	// etype-specific structs which contain exactly one pointer-shaped field are
	// stored as values rather than pointers when possible.
	//
	// TMAP: *MapType TFORW: *ForwardType TFUNC: *FuncType TINTERMETHOD:
	// InterMethType TSTRUCT: *StructType TINTER: *InterType TDDDFIELD:
	// DDDFieldType TFUNCARGS: FuncArgsType TCHANARGS: ChanArgsType TCHAN:
	// *ChanType TPTR32, TPTR64: PtrType TARRAY: *ArrayType TSLICE: SliceType
	Extra interface{}

	// Width is the width of this Type in bytes.
	Width      int64
	Nod        *Node // canonical OTYPE node
	Orig       *Type // original type (type literal or predefined type)
	Sym        *Sym  // symbol containing name, for named types
	Vargen     int32 // unique name for OTYPE/ONAME
	Lineno     int32 // line at which this type was declared, implicitly or explicitly
	Etype      EType // kind of type
	Noalg      bool  // suppress hash and eq algorithm generation
	Trecur     uint8 // to detect loops
	Printed    bool  // prevent duplicate export printing
	Local      bool  // created in this file
	Deferwidth bool
	Broke      bool  // broken type definition.
	Align      uint8 // the required alignment of this type, in bytes
}

type Val struct {
	// U contains one of: bool bool when n.ValCtype() == CTBOOL *Mpint int when
	// n.ValCtype() == CTINT, rune when n.ValCtype() == CTRUNE *Mpflt float when
	// n.ValCtype() == CTFLT *Mpcplx pair of floats when n.ValCtype() == CTCPLX
	// string string when n.ValCtype() == CTSTR *Nilval when n.ValCtype() ==
	// CTNIL
	U interface{}
}

// A Var represents a single variable that may be stored in a register.
// That variable may itself correspond to a hardware register,
// to represent the use of registers in the unoptimized instruction stream.
type Var struct {
}

// AddAux adds the offset in the aux fields (AuxInt and Aux) of v to a.
func AddAux(a *obj.Addr, v *ssa.Value)

func AddAux2(a *obj.Addr, v *ssa.Value, offset int64)

func Afunclit(a *obj.Addr, n *Node)

// generate:
// 	res = &n;
// The generated code checks that the result is not nil.

// generate:
//     res = &n;
// The generated code checks that the result is not nil.
func Agen(n *Node, res *Node)

// allocate a register (reusing res if possible) and generate
// 	a = &n
// The caller must call Regfree(a).
// The generated code checks that the result is not nil.

// allocate a register (reusing res if possible) and generate
//     a = &n
// The caller must call Regfree(a).
// The generated code checks that the result is not nil.
func Agenr(n *Node, a *Node, res *Node)

func Anyregalloc() bool

// compute total size of f's in/out arguments.
func Argsize(t *Type) int

func AtExit(f func())

// AutoVar returns a *Node and int64 representing the auto variable and offset
// within it where v should be spilled.
func AutoVar(v *ssa.Value) (*Node, int64)

// Bgen generates code for branches:
//
// 	if n == wantTrue {
// 		goto to
// 	}

// Bgen generates code for branches:
//
//     if n == wantTrue {
//         goto to
//     }
func Bgen(n *Node, wantTrue bool, likely int, to *obj.Prog)

// Bitno reports the lowest index of a 1 bit in b.
// It calls Fatalf if there is no 1 bit.
func Bitno(b uint64) int

// Brcom returns !(op).
// For example, Brcom(==) is !=.
func Brcom(op Op) Op

// Brrev returns reverse(op).
// For example, Brrev(<) is >.
func Brrev(op Op) Op

// Bvgen generates code for calculating boolean values:
// 	res = n == wantTrue

// Bvgen generates code for calculating boolean values:
//     res = n == wantTrue
func Bvgen(n, res *Node, wantTrue bool)

// generate:
// 	res = n;
// simplifies and calls Thearch.Gmove.
// if wb is true, need to emit write barriers.

// generate:
//     res = n;
// simplifies and calls Thearch.Gmove.
// if wb is true, need to emit write barriers.
func Cgen(n, res *Node)

// CgenTemp creates a temporary node, assigns n to it, and returns it.
func CgenTemp(n *Node) *Node

// generate:
// 	res, resok = x.(T)
// n.Left is x
// n.Type is T

// generate:
//     res, resok = x.(T)
// n.Left is x
// n.Type is T
func Cgen_As2dottype(n, res, resok *Node)

func Cgen_as(nl, nr *Node)

func Cgen_as_wb(nl, nr *Node, wb bool)

func Cgen_checknil(n *Node)

// generate:
// 	res = iface{typ, data}
// n->left is typ
// n->right is data

// generate:
//     res = iface{typ, data}
// n->left is typ
// n->right is data
func Cgen_eface(n *Node, res *Node)

// allocate a register (reusing res if possible) and generate
// 	a = n
// The caller must call Regfree(a).

// allocate a register (reusing res if possible) and generate
//     a = n
// The caller must call Regfree(a).
func Cgenr(n *Node, a *Node, res *Node)

func Clearp(p *obj.Prog)

// clearslim generates code to zero a slim node.
func Clearslim(n *Node)

func Complexgen(n *Node, res *Node)

func Complexmove(f *Node, t *Node)

func Complexop(n *Node, res *Node) bool

// Componentgen copies a composite value by moving its individual components.
// Slices, strings and interfaces are supported. Small structs or arrays with
// elements of basic type are also supported. nr is nil when assigning a zero
// value.
func Componentgen(nr, nl *Node) bool

func Datastring(s string, a *obj.Addr)

// gather series of offsets
// >=0 is direct addressed field
// <0 is pointer to next field (+1)
func Dotoffset(n *Node, oary []int64, nn **Node) int

func Dump(s string, n *Node)

func Dumpit(str string, r0 *Flow, isreg int)

// Eqtype reports whether t1 and t2 are identical, following the spec rules.
//
// Any cyclic type must go through a named type, and if one is
// named, it is only identical to the other if they are the same
// pointer (t1 == t2), so there's no chance of chasing cycles
// ad infinitum, so no need for a depth counter.

// Return 1 if t1 and t2 are identical, following the spec rules.
//
// Any cyclic type must go through a named type, and if one is
// named, it is only identical to the other if they are the same
// pointer (t1 == t2), so there's no chance of chasing cycles
// ad infinitum, so no need for a depth counter.
func Eqtype(t1, t2 *Type) bool

func Exit(code int)

func Fatalf(fmt_ string, args ...interface{})

func Fixlargeoffset(n *Node)

func Fldconv(f *Field, flag FmtFlag) string

func Flowend(graph *Graph)

func Flowstart(firstp *obj.Prog, newData func() interface{}) *Graph

func Flusherrors()

func Gbranch(as obj.As, t *Type, likely int) *obj.Prog

// compile statements
func Genlist(l Nodes)

func GetReg(r int) int

// generate:
// 	call f
// 	proc=-1	normal call but no return
// 	proc=0	normal call
// 	proc=1	goroutine run in new proc
// 	proc=2	defer call save away stack
// 	proc=3	normal call to C pointer (not Go func value)

// generate:
//     call f
//     proc=-1    normal call but no return
//     proc=0    normal call
//     proc=1    goroutine run in new proc
//     proc=2    defer call save away stack
//     proc=3    normal call to C pointer (not Go func value)
func Ginscall(f *Node, proc int)

func Gvardef(n *Node)

func Gvarkill(n *Node)

func Gvarlive(n *Node)

// Igen computes the address &n, stores it in a register r,
// and rewrites a to refer to *r. The chosen r may be the
// stack pointer, it may be borrowed from res, or it may
// be a newly allocated register. The caller must call Regfree(a)
// to free r when the address is no longer needed.
// The generated code ensures that &n is not nil.
func Igen(n *Node, a *Node, res *Node)

// Import populates importpkg from the serialized package data.
func Import(in *bufio.Reader)

// Is this a 64-bit type?
func Is64(t *Type) bool

func Isconst(n *Node, ct Ctype) bool

func Isfat(t *Type) bool

// Is this node a memory operand?
func Ismem(n *Node) bool

// IterFields returns the first field or method in struct or interface type t
// and an Iter value to continue iterating across the rest.
func IterFields(t *Type) (*Field, Iter)

func LOAD(r *Reg, z int) uint64

func Linksym(s *Sym) *obj.LSym

func Lookup(name string) *Sym

func LookupBytes(name []byte) *Sym

// LookupN looks up the symbol starting with prefix and ending with
// the decimal n. If prefix is too long, LookupN panics.
func LookupN(prefix string, n int) *Sym

func Lookupf(format string, a ...interface{}) *Sym

func Mfree(n *Node)

func Mgen(n *Node, n1 *Node, rg *Node)

// Naddr rewrites a to refer to n.
// It assumes that a is zeroed on entry.
func Naddr(a *obj.Addr, n *Node)

// Fmt '%N': Nodes.
// Flags: 'l' suffix with "(type %T)" where possible
// 	  '+h' in debug mode, don't recurse, no multiline output

// Fmt '%N': Nodes.
// Flags: 'l' suffix with "(type %T)" where possible
//       '+h' in debug mode, don't recurse, no multiline output
func Nconv(n *Node, flag FmtFlag) string

// NegOne returns a Node of type t with value -1.
func NegOne(t *Type) *Node

// Is a conversion between t1 and t2 a no-op?
func Noconv(t1 *Type, t2 *Type) bool

func Nod(op Op, nleft *Node, nright *Node) *Node

// NodSym makes a Node with Op op and with the Left field set to left
// and the Sym field set to sym. This is for ODOT and friends.
func NodSym(op Op, left *Node, sym *Sym) *Node

func Nodbool(b bool) *Node

func Nodconst(n *Node, t *Type, v int64)

func Nodindreg(n *Node, t *Type, r int)

func Nodintconst(v int64) *Node

func Nodreg(n *Node, t *Type, r int)

// p is a call instruction. Does the call fail to return?
func Noreturn(p *obj.Prog) bool

func Patch(p *obj.Prog, to *obj.Prog)

func Pkglookup(name string, pkg *Pkg) *Sym

func Prog(as obj.As) *obj.Prog

// Ptrto returns the Type *t.
// The returned struct must not be modified.
func Ptrto(t *Type) *Type

// allocate register of type t, leave in n.
// if o != N, o may be reusable register.
// caller must Regfree(n).
func Regalloc(n *Node, t *Type, o *Node)

func Regdump()

func Regfree(n *Node)

// Reginuse reports whether r is in use.
func Reginuse(r int) bool

// Regrealloc(n) undoes the effect of Regfree(n),
// so that a register can be given up but then reclaimed.
func Regrealloc(n *Node)

func Rnd(o int64, r int64) int64

func SSAGenFPJump(s *SSAGenState, b, next *ssa.Block, jumps *[2][2]FloatingEQNEJump)

// SSARegNum returns the register (in cmd/internal/obj numbering) to
// which v has been allocated. Panics if v is not assigned to a
// register.
// TODO: Make this panic again once it stops happening routinely.
func SSARegNum(v *ssa.Value) int16

func STORE(r *Reg, z int) uint64

func Samereg(a *Node, b *Node) bool

func SetReg(r, v int)

func Setmaxarg(t *Type, extra int32)

// even simpler simtype; get rid of ptr, bool.
// assuming that the front end has rejected
// all the invalid conversions (like ptr -> bool)
func Simsimtype(t *Type) EType

// magic number for signed division
// see hacker's delight chapter 10
func Smagic(m *Magic)

func Smallintconst(n *Node) bool

func Sysfunc(name string) *Node

// Fmt "%T": types.
// Flags: 'l' print definition, not name
// 	  'h' omit 'func' and receiver from function types, short type names
// 	  'u' package name, not prefix (FTypeId mode, sticky)

// Fmt "%T": types.
// Flags: 'l' print definition, not name
//       'h' omit 'func' and receiver from function types, short type names
//       'u' package name, not prefix (FTypeId mode, sticky)
func Tconv(t *Type, flag FmtFlag) string

// make a new off the books
func Tempname(nn *Node, t *Type)

// magic number for unsigned division
// see hacker's delight chapter 10
func Umagic(m *Magic)

func Uniqp(r *Flow) *Flow

func Uniqs(r *Flow) *Flow

func Warn(fmt_ string, args ...interface{})

func Warnl(line int32, fmt_ string, args ...interface{})

func Yyerror(format string, args ...interface{})

func (f *Field) Copy() *Field

// End returns the offset of the first byte immediately after this field.
func (f *Field) End() int64

// Append appends entries to f.
func (f *Fields) Append(s ...*Field)

// Index returns the i'th element of Fields.
// It panics if f does not have at least i+1 elements.
func (f *Fields) Index(i int) *Field

// Iter returns the first field in fs and an Iter value to continue iterating
// across its successor fields.
// Deprecated: New code should use Slice instead.
func (fs *Fields) Iter() (*Field, Iter)

// Len returns the number of entries in f.
func (f *Fields) Len() int

// Set sets f to a slice.
// This takes ownership of the slice.
func (f *Fields) Set(s []*Field)

// Slice returns the entries in f as a slice.
// Changes to the slice entries will be reflected in f.
func (f *Fields) Slice() []*Field

// Next returns the next field or method, if any.
func (i *Iter) Next() *Field

func (a *Mpflt) Add(b *Mpflt)

func (a *Mpflt) AddFloat64(c float64)

func (a *Mpflt) Cmp(b *Mpflt) int

func (a *Mpflt) CmpFloat64(c float64) int

func (a *Mpflt) Float32() float64

func (a *Mpflt) Float64() float64

func (a *Mpflt) Mul(b *Mpflt)

func (a *Mpflt) MulFloat64(c float64)

func (a *Mpflt) Neg()

func (a *Mpflt) Quo(b *Mpflt)

func (a *Mpflt) Set(b *Mpflt)

func (a *Mpflt) SetFloat64(c float64)

func (a *Mpflt) SetInt(b *Mpint)

// floating point input
// required syntax is [+-]d*[.]d*[e[+-]d*] or [+-]0xH*[e[+-]d*]
func (a *Mpflt) SetString(as string)

func (f *Mpflt) String() string

func (a *Mpflt) Sub(b *Mpflt)

func (a *Mpint) Add(b *Mpint)

func (a *Mpint) And(b *Mpint)

func (a *Mpint) AndNot(b *Mpint)

func (a *Mpint) Cmp(b *Mpint) int

func (a *Mpint) CmpInt64(c int64) int

func (a *Mpint) Int64() int64

func (a *Mpint) Lsh(b *Mpint)

func (a *Mpint) Mul(b *Mpint)

func (a *Mpint) Neg()

func (a *Mpint) Or(b *Mpint)

func (a *Mpint) Quo(b *Mpint)

func (a *Mpint) Rem(b *Mpint)

func (a *Mpint) Rsh(b *Mpint)

func (a *Mpint) Set(b *Mpint)

func (a *Mpint) SetFloat(b *Mpflt) int

func (a *Mpint) SetInt64(c int64)

func (a *Mpint) SetOverflow()

func (a *Mpint) SetString(as string)

func (x *Mpint) String() string

func (a *Mpint) Sub(b *Mpint)

func (a *Mpint) Xor(b *Mpint)

// Bool returns n as an bool.
// n must be an boolean constant.
func (n *Node) Bool() bool

// Convconst converts constant node n to type t and
// places the result in con.
func (n *Node) Convconst(con *Node, t *Type)

func (n *Node) HasBreak() bool

// Int64 returns n as an int64.
// n must be an integer or rune constant.
func (n *Node) Int64() int64

// IntLiteral returns the Node's literal value as an integer.

// IntLiteral returns the Node's literal value as an interger.
func (n *Node) IntLiteral() (x int64, ok bool)

func (n *Node) IsOutputParamHeapAddr() bool

func (n *Node) Line() string

func (n *Node) NotLiveAtEnd() bool

// Opt returns the optimizer data for the node.
func (n *Node) Opt() interface{}

// SetBigInt sets n's value to x.
// n must be an integer constant.
func (n *Node) SetBigInt(x *big.Int)

func (n *Node) SetHasBreak(b bool)

// SetInt sets n's value to i.
// n must be an integer constant.
func (n *Node) SetInt(i int64)

func (n *Node) SetNotLiveAtEnd(b bool)

// SetOpt sets the optimizer data for the node, which must not have been used
// with SetVal. SetOpt(nil) is ignored for Vals to simplify call sites that are
// clearing Opts.
func (n *Node) SetOpt(x interface{})

// SetSliceBounds sets n's slice bounds, where n is a slice expression. n must
// be a slice expression. If max is non-nil, n must be a full slice expression.
func (n *Node) SetSliceBounds(low, high, max *Node)

// SetVal sets the Val for the node, which must not have been used with SetOpt.
func (n *Node) SetVal(v Val)

// SliceBounds returns n's slice bounds: low, high, and max in
// expr[low:high:max]. n must be a slice expression. max is nil if n is a simple
// slice expression.
func (n *Node) SliceBounds() (low, high, max *Node)

func (n *Node) String() string

func (n *Node) Typ() ssa.Type

// Val returns the Val for the node.
func (n *Node) Val() Val

// Append appends entries to Nodes.
// If a slice is passed in, this will take ownership of it.
func (n *Nodes) Append(a ...*Node)

// AppendNodes appends the contents of *n2 to n, then clears n2.
func (n *Nodes) AppendNodes(n2 *Nodes)

// MoveNodes sets n to the contents of n2, then clears n2.
func (n *Nodes) MoveNodes(n2 *Nodes)

// Set sets n to a slice.
// This takes ownership of the slice.
func (n *Nodes) Set(s []*Node)

// Set1 sets n to a slice containing a single node.
func (n *Nodes) Set1(node *Node)

func (pkg *Pkg) Lookup(name string) *Sym

func (pkg *Pkg) LookupBytes(name []byte) *Sym

// Pc returns the current Prog.
func (s *SSAGenState) Pc() *obj.Prog

// SetLineno sets the current source line number.
func (s *SSAGenState) SetLineno(l int32)

func (s *Sym) String() string

func (t *Type) Alignment() int64

func (t *Type) AllMethods() *Fields

// ArgWidth returns the total aligned argument size for a function.
// It includes the receiver, parameters, and results.
func (t *Type) ArgWidth() int64

// ChanArgs returns the channel type for TCHANARGS type t.
func (t *Type) ChanArgs() *Type

// ChanDir returns the direction of a channel type t.
// The direction will be one of Crecv, Csend, or Cboth.
func (t *Type) ChanDir() ChanDir

// ChanType returns t's extra channel-specific fields.
func (t *Type) ChanType() *ChanType

// Compare compares types for purposes of the SSA back
// end, returning an ssa.Cmp (one of CMPlt, CMPeq, CMPgt).
// The answers are correct for an optimizer
// or code generator, but not necessarily typechecking.
// The order chosen is arbitrary, only consistency and division
// into equivalence classes (Types that compare CMPeq) matters.
func (t *Type) Compare(u ssa.Type) ssa.Cmp

// Copy returns a shallow copy of the Type.
func (t *Type) Copy() *Type

// DDDField returns the slice ... type for TDDDFIELD type t.
func (t *Type) DDDField() *Type

// Elem returns the type of elements of t.
// Usable with pointers, channels, arrays, and slices.
func (t *Type) Elem() *Type

func (t *Type) ElemType() ssa.Type

// Field returns the i'th field/method of struct/interface type t.
func (t *Type) Field(i int) *Field

func (t *Type) FieldName(i int) string

func (t *Type) FieldOff(i int) int64

// FieldSlice returns a slice of containing all fields/methods of
// struct/interface type t.
func (t *Type) FieldSlice() []*Field

func (t *Type) FieldType(i int) ssa.Type

func (t *Type) Fields() *Fields

// ForwardType returns t's extra forward-type-specific fields.
func (t *Type) ForwardType() *ForwardType

// FuncArgs returns the channel type for TFUNCARGS type t.
func (t *Type) FuncArgs() *Type

// FuncType returns t's extra func-specific fields.
func (t *Type) FuncType() *FuncType

// IncomparableField returns an incomparable Field of struct Type t, if any.
func (t *Type) IncomparableField() *Field

func (t *Type) IsArray() bool

func (t *Type) IsBoolean() bool

func (t *Type) IsChan() bool

// IsComparable reports whether t is a comparable type.
func (t *Type) IsComparable() bool

func (t *Type) IsComplex() bool

// IsEmptyInterface reports whether t is an empty interface type.
func (t *Type) IsEmptyInterface() bool

func (t *Type) IsFlags() bool

func (t *Type) IsFloat() bool

// IsFuncArgStruct reports whether t is a struct representing function
// parameters.
func (t *Type) IsFuncArgStruct() bool

func (t *Type) IsInteger() bool

func (t *Type) IsInterface() bool

// IsKind reports whether t is a Type of the specified kind.
func (t *Type) IsKind(et EType) bool

func (t *Type) IsMap() bool

func (t *Type) IsMemory() bool

// IsPtr reports whether t is a regular Go pointer type.
// This does not include unsafe.Pointer.
func (t *Type) IsPtr() bool

// IsPtrShaped reports whether t is represented by a single machine pointer. In
// addition to regular Go pointer types, this includes map, channel, and
// function types and unsafe.Pointer. It does not include array or struct types
// that consist of a single pointer shaped type. TODO(mdempsky): Should it? See
// golang.org/issue/15028.
func (t *Type) IsPtrShaped() bool

// IsRegularMemory reports whether t can be compared/hashed as regular memory.
func (t *Type) IsRegularMemory() bool

func (t *Type) IsSigned() bool

func (t *Type) IsSlice() bool

func (t *Type) IsString() bool

func (t *Type) IsStruct() bool

// IsUnsafePtr reports whether t is an unsafe pointer.
func (t *Type) IsUnsafePtr() bool

// IsUntyped reports whether t is an untyped type.
func (t *Type) IsUntyped() bool

func (t *Type) IsVoid() bool

// Key returns the key type of map type t.
func (t *Type) Key() *Type

// MapType returns t's extra map-specific fields.
func (t *Type) MapType() *MapType

func (t *Type) Methods() *Fields

// Nname returns the associated function's nname.
func (t *Type) Nname() *Node

func (t *Type) NumElem() int64

func (t *Type) NumFields() int

func (t *Type) Params() *Type

func (t *Type) ParamsP() **Type

func (t *Type) PtrTo() ssa.Type

// Recv returns the receiver of function type t, if any.
func (t *Type) Recv() *Field

func (t *Type) Recvs() *Type

func (t *Type) RecvsP() **Type

func (t *Type) Results() *Type

func (t *Type) ResultsP() **Type

// SetFields sets struct/interface type t's fields/methods to fields.
func (t *Type) SetFields(fields []*Field)

// Nname sets the associated function's nname.
func (t *Type) SetNname(n *Node)

// SetNumElem sets the number of elements in an array type.
// The only allowed use is on array types created with typDDDArray.
// For other uses, create a new array with typArray instead.
func (t *Type) SetNumElem(n int64)

func (t *Type) SimpleString() string

func (t *Type) Size() int64

func (t *Type) String() string

// StructType returns t's extra struct-specific fields.
func (t *Type) StructType() *StructType

// Val returns the value type of map type t.
func (t *Type) Val() *Type

// String returns a space-separated list of the variables represented by bits.
func (bits Bits) String() string

func (c ChanDir) CanRecv() bool

func (c ChanDir) CanSend() bool

func (et EType) String() string

// Addr returns the address of the i'th element of Nodes.
// It panics if n does not have at least i+1 elements.
func (n Nodes) Addr(i int) **Node

// First returns the first element of Nodes (same as n.Index(0)).
// It panics if n has no elements.
func (n Nodes) First() *Node

// Index returns the i'th element of Nodes.
// It panics if n does not have at least i+1 elements.
func (n Nodes) Index(i int) *Node

// Len returns the number of entries in Nodes.
func (n Nodes) Len() int

// Second returns the second element of Nodes (same as n.Index(1)).
// It panics if n has fewer than two elements.
func (n Nodes) Second() *Node

// SetIndex sets the i'th element of Nodes to node.
// It panics if n does not have at least i+1 elements.
func (n Nodes) SetIndex(i int, node *Node)

// Slice returns the entries in Nodes as a slice.
// Changes to the slice entries (as in s[i] = n) will be reflected in
// the Nodes.
func (n Nodes) Slice() []*Node

func (n Nodes) String() string

func (o Op) GoString() string

// IsSlice3 reports whether o is a slice3 op (OSLICE3, OSLICE3ARR).
// o must be a slicing op.
func (o Op) IsSlice3() bool

func (o Op) String() string

func (v Val) Ctype() Ctype

