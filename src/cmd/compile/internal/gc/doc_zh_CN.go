// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package gc // import "cmd/compile/internal/gc"

import (
    "bytes"
    "cmd/compile/internal/big"
    "cmd/internal/gcprog"
    "cmd/internal/obj"
    "cmd/internal/obj/ppc64"
    "crypto/md5"
    "encoding/binary"
    "flag"
    "fmt"
    "internal/testenv"
    "io"
    "io/ioutil"
    "log"
    "math"
    "os"
    "os/exec"
    "path"
    "reflect"
    "runtime"
    "runtime/pprof"
    "sort"
    "strconv"
    "strings"
    "testing"
    "unicode"
    "unicode/utf8"
)

const (
    // These values are known by runtime.
    // The MEMx and NOEQx values must run in parallel.  See algtype.
    AMEM = iota
    AMEM0
    AMEM8
    AMEM16
    AMEM32
    AMEM64
    AMEM128
    ANOEQ
    ANOEQ0
    ANOEQ8
    ANOEQ16
    ANOEQ32
    ANOEQ64
    ANOEQ128
    ASTRING
    AINTER
    ANILINTER
    ASLICE
    AFLOAT32
    AFLOAT64
    ACPLX64
    ACPLX128
    AUNK = 100
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
    // must match ../../pkg/nreflect/type.go:/Chandir
    Cxxx  = 0
    Crecv = 1 << 0
    Csend = 1 << 1
    Cboth = Crecv | Csend
)

const (
    EOF = -1
)

const (
    EscFuncUnknown = 0 + iota
    EscFuncPlanned
    EscFuncStarted
    EscFuncTagged
)

// Escape constants are numbered in order of increasing "escapiness"
// to help make inferences be monotonic.  With the exception of
// EscNever which is sticky, eX < eY means that eY is more exposed
// than eX, and hence replaces it in a conservative analysis.
const (
    EscUnknown = iota
    EscNone    // Does not escape to heap, result, or parameters.
    EscReturn  // Is returned or reachable from returned.
    EscScope   // Allocated in an inner loop scope, assigned to an outer loop scope,
    // which allows the construction of non-escaping but arbitrarily large linked
    // data structures (i.e., not eligible for allocation in a fixed-size stack frame).
    EscHeap           // Reachable from the heap
    EscNever          // By construction will not escape.
    EscBits           = 3
    EscMask           = (1 << EscBits) - 1
    EscContentEscapes = 1 << EscBits // value obtained by indirect of parameter escapes to heap
    EscReturnBits     = EscBits + 1
)

const (
    Etop      = 1 << 1 // evaluated at statement level
    Erv       = 1 << 2 // evaluated in value context
    Etype     = 1 << 3
    Ecall     = 1 << 4  // call-only expressions are ok
    Efnstruct = 1 << 5  // multivalue function returns are ok
    Eiota     = 1 << 6  // iota is ok
    Easgn     = 1 << 7  // assigning to expression
    Eindir    = 1 << 8  // indirecting through expression
    Eaddr     = 1 << 9  // taking address of expression
    Eproc     = 1 << 10 // inside a go statement
    Ecomplit  = 1 << 11 // type in composite literal
)

// Format conversions
//
//         %L int        Line numbers
//
//         %E int        etype values (aka 'Kind')
//
//         %O int        Node Opcodes
//             Flags: "%#O": print go syntax. (automatic unless fmtmode == FDbg)
//
//         %J Node*    Node details
//             Flags: "%hJ" suppresses things not relevant until walk.
//
//         %V Val*        Constant values
//
//         %S Sym*        Symbols
//             Flags: +,- #: mode (see below)
//                 "%hS"    unqualified identifier in any mode
//                 "%hhS"  in export mode: unqualified identifier if exported, qualified if not
//
//         %T Type*    Types
//             Flags: +,- #: mode (see below)
//                 'l' definition instead of name.
//                 'h' omit "func" and receiver in function types
//                 'u' (only in -/Sym mode) print type identifiers wit package name instead of prefix.
//
//         %N Node*    Nodes
//             Flags: +,- #: mode (see below)
//                 'h' (only in +/debug mode) suppress recursion
//                 'l' (only in Error mode) print "foo (type Bar)"
//
//         %H NodeList*    NodeLists
//             Flags: those of %N
//                 ','  separate items with ',' instead of ';'
//
//       In mparith2.go and mparith3.go:
//             %B Mpint*    Big integers
//             %F Mpflt*    Big floats
//
//       %S, %T and %N obey use the following flags to set the format mode:

// Format conversions
//
//       %L int        Line numbers
//
//       %E int        etype values (aka 'Kind')
//
//       %O int        Node Opcodes
//           Flags: "%#O": print go syntax. (automatic unless fmtmode == FDbg)
//
//       %J Node*    Node details
//           Flags: "%hJ" suppresses things not relevant until walk.
//
//       %V Val*        Constant values
//
//       %S Sym*        Symbols
//           Flags: +,- #: mode (see below)
//               "%hS"    unqualified identifier in any mode
//               "%hhS"  in export mode: unqualified identifier if exported, qualified if not
//
//       %T Type*    Types
//           Flags: +,- #: mode (see below)
//               'l' definition instead of name.
//               'h' omit "func" and receiver in function types
//               'u' (only in -/Sym mode) print type identifiers wit package name instead of prefix.
//
//       %N Node*    Nodes
//           Flags: +,- #: mode (see below)
//               'h' (only in +/debug mode) suppress recursion
//               'l' (only in Error mode) print "foo (type Bar)"
//
//       %H NodeList*    NodeLists
//           Flags: those of %N
//               ','  separate items with ',' instead of ';'
//
//     In mparith2.go and mparith3.go:
//           %B Mpint*    Big integers
//           %F Mpflt*    Big floats
//
//     %S, %T and %N obey use the following flags to set the format mode:
const (
    FErr = iota
    FDbg
    FExp
    FTypeId
)

// FNV-1 hash function constants.
const (
    H0  = 2166136261
    Hp  = 16777619
)

// static initialization
const (
    InitNotStarted = 0
    InitDone       = 1
    InitPending    = 2
)

const (
    LLITERAL = 57346 + iota
    LASOP
    LCOLAS
    LBREAK
    LCASE
    LCHAN
    LCONST
    LCONTINUE
    LDDD
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
    LNAME
    LPACKAGE
    LRANGE
    LRETURN
    LSELECT
    LSTRUCT
    LSWITCH
    LTYPE
    LVAR
    LANDAND
    LANDNOT
    LCOMM
    LDEC
    LEQ
    LGE
    LGT
    LIGNORE
    LINC
    LLE
    LLSH
    LLT
    LNE
    LOROR
    LRSH
)

const (
    MODEDYNAM = 1
    MODECONST = 2
)

// MaxFlowProg is the maximum size program (counted in instructions)
// for which the flow code will build a graph. Functions larger than this limit
// will not have flow graphs and consequently will not be optimized.
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

// The parser's maximum stack size.
// We have to use a #define macro here since yacc
// or bison will check for its definition and use
// a potentially smaller value if it is undefined.
const (
    NHUNK           = 50000
    BUFSIZ          = 8192
    NSYMB           = 500
    NHASH           = 1024
    MAXALIGN        = 7
    UINF            = 100
    PRIME1          = 3
    BADWIDTH        = -1000000000
    MaxStackVarSize = 10 * 1024 * 1024
)

const NOTALOOPDEPTH = -1

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
    OADDSTR          // Left + Right (string addition)
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
    ODCLFUNC  // func f() or func (r) f()
    ODCLFIELD // struct field, interface field, or func/method argument/return value.
    ODCLCONST // const pi = 3.14
    ODCLTYPE  // type Int int

    ODELETE    // delete(Left, Right)
    ODOT       // Left.Right (Left is of struct type)
    ODOTPTR    // Left.Right (Left is of pointer to struct type)
    ODOTMETH   // Left.Right (Left is non-interface, Right is method name)
    ODOTINTER  // Left.Right (Left is interface, Right is method name)
    OXDOT      // Left.Right (before rewrite to one of the preceding)
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
    OPARAM     // variant of ONAME for on-stack copy of a parameter or return value that escapes.
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
    // Pseudo-op, like TEXT, GLOBL, TYPE, PCDATA, FUNCDATA.
    Pseudo = 1 << 1

    // There's nothing to say about the instruction,
    // but it's still okay to see.
    OK  = 1 << 2

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

    // Set, use, or kill of carry bit.
    // Kill means we never look at the carry bit after this kind of instruction.
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
    PPARAM          // input arguments
    PPARAMOUT       // output results
    PPARAMREF       // closure variable reference
    PFUNC           // global function

    PDISCARD // discard during parse of duplicate import

    PHEAP = 1 << 7 // an extra bit to identify an escaped variable
)

const (
    SymExport   = 1 << 0 // to be exported
    SymPackage  = 1 << 1
    SymExported = 1 << 2 // already written out by export
    SymUniq     = 1 << 3
    SymSiggen   = 1 << 4
    SymAsm      = 1 << 5
    SymAlgGen   = 1 << 6
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
    TARRAY
    T_old_DARRAY // Doesn't seem to be used in existing code. Used now for Isddd export (see bexport.go). TODO(gri) rename.
    TSTRUCT
    TCHAN
    TMAP
    TINTER
    TFORW
    TFIELD
    TANY
    TSTRING
    TUNSAFEPTR

    // pseudo-types for literals
    TIDEAL
    TNIL
    TBLANK

    // pseudo-type for frame layout
    TFUNCARGS
    TCHANARGS
    TINTERMETH

    NTYPE
)

const (
    UNVISITED = 0
    VISITED   = 1
)

const (
    WORDSIZE  = 4
    WORDBITS  = 32
    WORDMASK  = WORDBITS - 1
    WORDSHIFT = 5
)

// note this is the runtime representation
// of the compilers arrays.
//
// typedef    struct
// {                    // must not move anything
//     uchar    array[8];    // pointer to data
//     uchar    nel[4];        // number of elements
//     uchar    cap[4];        // allocated number of elements
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

var (
    Debug_append int
    Debug_panic  int
    Debug_slice  int
    Debug_wb     int
)

var Debug_checknil int

var (
    Debug_export int // if set, print debugging information about export data

)

var Debug_gcprog int // set by -d gcprog


var Debug_typeassert int

var Deferproc *Node

var Deferreturn *Node

var Disable_checknil int

var Funcdepth int32

var (
    Isptr [NTYPE]bool

    Isint     [NTYPE]bool
    Isfloat   [NTYPE]bool
    Iscomplex [NTYPE]bool
    Issigned  [NTYPE]bool
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


var Types [NTYPE]*Type

var Widthint int

var Widthptr int

var Widthreg int

type Arch struct {
    Thechar      int
    Thestring    string
    Thelinkarch  *obj.LinkArch
    Typedefs     []Typedef
    REGSP        int
    REGCTXT      int
    REGCALLX     int // BX
    REGCALLX2    int // AX
    REGRETURN    int // AX
    REGMIN       int
    REGMAX       int
    REGZERO      int // architectural zero register, if available
    FREGMIN      int
    FREGMAX      int
    MAXWIDTH     int64
    ReservedRegs []int

    AddIndex     func(*Node, int64, *Node) bool // optional
    Betypeinit   func()
    Bgen_float   func(*Node, bool, int, *obj.Prog) // optional
    Cgen64       func(*Node, *Node)                // only on 32-bit systems
    Cgenindex    func(*Node, *Node, bool) *obj.Prog
    Cgen_bmul    func(Op, *Node, *Node, *Node) bool
    Cgen_float   func(*Node, *Node) // optional
    Cgen_hmul    func(*Node, *Node, *Node)
    Cgen_shift   func(Op, bool, *Node, *Node, *Node)
    Clearfat     func(*Node)
    Cmp64        func(*Node, *Node, Op, int, *obj.Prog) // only on 32-bit systems
    Defframe     func(*obj.Prog)
    Dodiv        func(Op, *Node, *Node, *Node)
    Excise       func(*Flow)
    Expandchecks func(*obj.Prog)
    Getg         func(*Node)
    Gins         func(int, *Node, *Node) *obj.Prog

    // Ginscmp generates code comparing n1 to n2 and jumping away if op is satisfied.
    // The returned prog should be Patch'ed with the jump target.
    // If op is not satisfied, code falls through to the next emitted instruction.
    // Likely is the branch prediction hint: +1 for likely, -1 for unlikely, 0 for no opinion.
    //
    // Ginscmp must be able to handle all kinds of arguments for n1 and n2,
    // not just simple registers, although it can assume that there are no
    // function calls needed during the evaluation, and on 32-bit systems
    // the values are guaranteed not to be 64-bit values, so no in-memory
    // temporaries are necessary.
    Ginscmp func(op Op, t *Type, n1, n2 *Node, likely int) *obj.Prog

    // Ginsboolval inserts instructions to convert the result
    // of a just-completed comparison to a boolean value.
    // The first argument is the conditional jump instruction
    // corresponding to the desired value.
    // The second argument is the destination.
    // If not present, Ginsboolval will be emulated with jumps.
    Ginsboolval func(int, *Node)

    Ginscon      func(int, int64, *Node)
    Ginsnop      func()
    Gmove        func(*Node, *Node)
    Igenindex    func(*Node, *Node, bool) *obj.Prog
    Linkarchinit func()
    Peep         func(*obj.Prog)
    Proginfo     func(*obj.Prog) // fills in Prog.Info
    Regtyp       func(*obj.Addr) bool
    Sameaddr     func(*obj.Addr, *obj.Addr) bool
    Smallindir   func(*obj.Addr, *obj.Addr) bool
    Stackaddr    func(*obj.Addr) bool
    Blockcopy    func(*Node, *Node, int64, int64, int64)
    Sudoaddable  func(int, *Node, *obj.Addr) bool
    Sudoclean    func()
    Excludedregs func() uint64
    RtoB         func(int) uint64
    FtoB         func(int) uint64
    BtoR         func(uint64) int
    BtoF         func(uint64) int
    Optoas       func(Op, *Type) int
    Doregbits    func(int) uint64
    Regnames     func(*int) []string
    Use387       bool // should 8g use 387 FP instructions instead of sse2.
}

// An ordinary basic block.
//
// Instructions are threaded together in a doubly-linked list.  To iterate in
// program order follow the link pointer from the first node and stop after the
// last node has been visited
//
//   for(p = bb->first;; p = p->link) {
//     ...
//     if(p == bb->last)
//       break;
//   }
//
// To iterate in reverse program order by following the opt pointer from the
// last node
//
//   for(p = bb->last; p != nil; p = p->opt) {
//     ...
//   }
type BasicBlock struct {
    pred            []*BasicBlock // predecessors; if none, probably start of CFG
    succ            []*BasicBlock // successors; if none, probably ends in return statement
    first           *obj.Prog     // first instruction in block
    last            *obj.Prog     // last instruction in block
    rpo             int           // reverse post-order number (also index in cfg)
    mark            int           // mark bit for traversals
    lastbitmapindex int           // for livenessepilogue

    // Computed during livenessprologue using only the content of
    // individual blocks:
    //
    //	uevar: upward exposed variables (used before set in block)
    //	varkill: killed variables (set in block)
    //	avarinit: addrtaken variables set or used (proof of initialization)
    uevar    Bvec
    varkill  Bvec
    avarinit Bvec

    // Computed during livenesssolve using control flow information:
    //
    //	livein: variables live at block entry
    //	liveout: variables live at block exit
    //	avarinitany: addrtaken variables possibly initialized at block exit
    //		(initialized in block or at exit from any predecessor block)
    //	avarinitall: addrtaken variables certainly initialized at block exit
    //		(initialized in block or at exit from all predecessor blocks)
    livein      Bvec
    liveout     Bvec
    avarinitany Bvec
    avarinitall Bvec
}

// Bits represents a set of Vars, stored as a bit set of var numbers
// (the index in vars, or equivalently v.id).
type Bits struct {
    b [BITS]uint64
}

// A Bvec is a bit vector.
type Bvec struct {
    n   int32    // number of bits in vector
    b   []uint32 // words holding bits
}

// The Class of a variable/function describes the "storage class"
// of a variable or function. During parsing, storage classes are
// called declaration contexts.
type Class uint8

// Ctype describes the constant kind of an "ideal" (untyped) constant.
type Ctype int8

type Dlist struct {
    field *Type
}

type EType uint8

type Error struct {
    lineno int
    seq    int
    msg    string
}

type EscState struct {
    // Fake node that all
    //   - return values and output variables
    //   - parameters on imported functions not marked 'safe'
    //   - assignments to global variables
    // flow to.
    theSink Node

    dsts      *NodeList // all dst nodes
    loopdepth int32     // for detecting nested loop scopes
    pdepth    int       // for debug printing in recursions.
    dstcount  int       // diagnostic
    edgecount int       // diagnostic
    noesc     *NodeList // list of possible non-escaping nodes, for printing
    recursive bool      // recursive function or group of mutually recursive functions.
    opts      []*Node   // nodes with .Opt initialized
    walkgen   uint32
}

type Flow struct {
    Prog   *obj.Prog // actual instruction
    P1     *Flow     // predecessors of this instruction: p1,
    P2     *Flow     // and then p2 linked though p2link.
    P2link *Flow
    S1     *Flow // successors of this instruction (at most two: s1 and s2).
    S2     *Flow
    Link   *Flow // next instruction in function code

    Active int32 // usable by client

    Id     int32  // sequence number in flow graph
    Rpo    int32  // reverse post ordering
    Loop   uint16 // x5 for every loop
    Refset bool   // diagnostic generated

    Data interface{} // for use by client
}

// Func holds Node fields used only with function-like nodes.
type Func struct {
    Shortname  *Node
    Enter      *NodeList
    Exit       *NodeList
    Cvars      *NodeList // closure params
    Dcl        *NodeList // autodcl for this func/closure
    Inldcl     *NodeList // copy of dcl for use in inlining
    Closgen    int
    Outerfunc  *Node
    Fieldtrack []*Type
    Outer      *Node // outer func for closure
    Ntype      *Node // signature
    Top        int   // top context (Ecall, Eproc, etc)
    Closure    *Node // OCLOSURE <-> ODCLFUNC
    FCurfn     *Node
    Nname      *Node

    Inl     *NodeList // copy of the body for use in inlining
    InlCost int32
    Depth   int32

    Endlineno int32

    Norace            bool // func must not have race detector annotations
    Nosplit           bool // func should not execute on separate stack
    Noinline          bool // func should not be inlined
    Nowritebarrier    bool // emit compiler error instead of write barrier
    Nowritebarrierrec bool // error on write barrier in this or recursive callees
    Dupok             bool // duplicate definitions ok
    Wrapper           bool // is method wrapper
    Needctxt          bool // function uses context register (has closure variables)
    Systemstack       bool // must run on system stack

    WBLineno int32 // line number of first write barrier
}

type GCProg struct {
    sym    *Sym
    symoff int
    w      gcprog.Writer
}

type Graph struct {
    Start *Flow
    Num   int

    // After calling flowrpo, rpo lists the flow nodes in reverse postorder,
    // and each non-dead Flow node f has g->rpo[f->rpo] == f.
    Rpo []*Flow
}

type Idir struct {
    link *Idir
    dir  string
}

type InitEntry struct {
    Xoffset int64 // struct, array only
    Expr    *Node // bytes of run-time computed expressions
}

type InitPlan struct {
    Lit  int64
    Zero int64
    Expr int64
    E    []InitEntry
}

type Io struct {
    infile     string
    bin        *obj.Biobuf
    cp         string // used for content when bin==nil
    last       int
    peekc      int
    peekc1     int // second peekc for ...
    nlsemi     bool
    eofnl      bool
    importsafe bool
}

type Iter struct {
    Done  int
    Tfunc *Type
    T     *Type
}

type Label struct {
    Sym  *Sym
    Def  *Node
    Use  []*Node
    Link *Label

    // for use during gen
    Gotopc   *obj.Prog // pointer to unresolved gotos
    Labelpc  *obj.Prog // pointer to code
    Breakpc  *obj.Prog // pointer to code
    Continpc *obj.Prog // pointer to code

    Used bool
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
    value, suffixValue int8
}

// A collection of global state used by liveness analysis.
type Liveness struct {
    fn   *Node
    ptxt *obj.Prog
    vars []*Node
    cfg  []*BasicBlock

    // An array with a bit vector for each safe point tracking live pointers
    // in the arguments and locals area, indexed by bb.rpo.
    argslivepointers []Bvec
    livepointers     []Bvec
}

// argument passing to/from
// smagic and umagic
type Magic struct {
    W   int // input for both - width
    S   int // output for both - shift
    Bad int // output for both - unexpected failure

    // magic multiplier for signed literal divisors
    Sd  int64 // input - literal divisor
    Sm  int64 // output - multiplier

    // magic multiplier for unsigned literal divisors
    Ud  uint64 // input - literal divisor
    Um  uint64 // output - multiplier
    Ua  int    // output - adder
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

// Name holds Node fields used only by named nodes (ONAME, OPACK, some
// OLITERAL).
type Name struct {
    Pack      *Node // real package for import . names
    Pkg       *Pkg  // pkg for OPACK nodes
    Heapaddr  *Node // temp holding heap address of param
    Inlvar    *Node // ONAME substitute while inlining
    Defn      *Node // initializing assignment
    Curfn     *Node // function for local variables
    Param     *Param
    Decldepth int32 // declaration loop depth, increased for every loop or label
    Vargen    int32 // unique name for ONAME within a function.  Function outputs are numbered starting at one.
    Iota      int32 // value if this name is iota
    Funcdepth int32
    Method    bool // OCALLMETH name
    Readonly  bool
    Captured  bool // is the variable captured by a closure
    Byval     bool // is the variable captured by value or by reference
    Needzero  bool // if it contains pointers, needs to be zeroed on function entry
    Keepalive bool // mark value live across unknown assembly call
}

type NilVal struct{}

// A Node is a single node in the syntax tree.
// Actually the syntax tree is a syntax DAG, because there is only one
// node with Op=ONAME for a given instance of a variable x.
// The same is true for Op=OTYPE and Op=OLITERAL.
type Node struct {
    // Tree structure.
    // Generic recursive walks should follow these fields.
    Left  *Node
    Right *Node
    Ninit *NodeList
    Nbody *NodeList
    List  *NodeList
    Rlist *NodeList

    // most nodes
    Type *Type
    Orig *Node // original form, for printing, and tracking copies of ONAMEs

    // func
    Func *Func

    // ONAME
    Name *Name

    Sym *Sym        // various
    E   interface{} // Opt or Val, see methods below

    Xoffset int64

    Lineno int32

    // OREGISTER, OINDREG
    Reg int16

    Esc uint16 // EscXXX

    Op          Op
    Nointerface bool
    Ullman      uint8 // sethi/ullman number
    Addable     bool  // addressable
    Etype       EType // op for OASOP, etype for OTYPE, exclam for export, 6g saved reg
    Bounded     bool  // bounds check unnecessary
    Class       Class // PPARAM, PAUTO, PEXTERN, etc
    Embedded    uint8 // ODCLFIELD embedded type
    Colas       bool  // OAS resulting from :=
    Diag        uint8 // already printed error about this
    Noescape    bool  // func arguments do not escape; TODO(rsc): move Noescape to Func struct (see CL 7360)
    Walkdef     uint8
    Typecheck   uint8
    Local       bool
    Dodata      uint8
    Initorder   uint8
    Used        bool
    Isddd       bool // is the argument variadic
    Implicit    bool
    Addrtaken   bool // address taken, even if not moved to heap
    Assigned    bool // is the variable ever assigned to
    Likely      int8 // likeliness of if statement
    Hasbreak    bool // has break statement
    hasVal      int8 // +1 for Val, -1 for Opt, 0 for not yet set
}

type NodeEscState struct {
    Curfn             *Node
    Escflowsrc        *NodeList // flow(this, src)
    Escretval         *NodeList // on OCALLxxx, list of dummy return values
    Escloopdepth      int32     // -1: global, 0: return variables, 1:function top level, increased inside function for every loop or label to mark scopes
    Esclevel          Level
    Walkgen           uint32
    Maxextraloopdepth int32
}

// A NodeList is a linked list of nodes.
// TODO(rsc): Some uses of NodeList should be made into slices.
// The remaining ones probably just need a simple linked list,
// not one with concatenation support.
type NodeList struct {
    N    *Node
    Next *NodeList
    End  *NodeList
}

type Op uint8

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
    out  *NodeList // list of generated statements
    temp *NodeList // head of stack of temporary variables
    free *NodeList // free list of NodeList* structs (for use in temp)
}

type Param struct {
    Ntype *Node

    // ONAME func param with PHEAP
    Outerexpr  *Node // expression copied into closure for variable
    Stackparam *Node // OPARAM node referring to stack copy of param

    // ONAME PPARAM
    Field *Type // TFIELD in arg struct

    // ONAME closure param with PPARAMREF
    Outer   *Node // outer PPARAMREF in nested closure
    Closure *Node // ONAME/PHEAP <-> ONAME/PPARAMREF
}

type Pkg struct {
    Name     string // package name
    Path     string // string literal used in import statement
    Pathsym  *Sym
    Prefix   string // escaped path for use in symbol table
    Imported bool   // export data of this package was parsed
    Exported bool   // import line written in export data
    Direct   bool   // imported directly
    Safe     bool   // whether the package is marked as safe
    Syms     map[string]*Sym
}

// A Reg is a wrapper around a single Prog (one instruction) that holds
// register optimization information while the optimizer runs.
// r->prog is the instruction.
type Reg struct {
    set  Bits // regopt variables written by this instruction.
    use1 Bits // regopt variables read by prog->from.
    use2 Bits // regopt variables read by prog->to.

    // refahead/refbehind are the regopt variables whose current
    // value may be used in the following/preceding instructions
    // up to a CALL (or the value is clobbered).
    refbehind Bits
    refahead  Bits

    // calahead/calbehind are similar, but for variables in
    // instructions that are reachable after hitting at least one
    // CALL.
    calbehind Bits
    calahead  Bits

    regdiff Bits
    act     Bits
    regu    uint64 // register used bitmap
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
    enter *Flow
    cost  int16
    varno int16
    regno int16
}

type Sig struct {
    name   string
    pkg    *Pkg
    isym   *Sym
    tsym   *Sym
    type_  *Type
    mtype  *Type
    offset int32
}

type Sym struct {
    Lexical   uint16
    Flags     uint8
    Link      *Sym
    Uniqgen   uint32
    Importdef *Pkg   // where imported definition was found
    Linkname  string // link name

    // saved and restored by dcopy
    Pkg        *Pkg
    Name       string // variable name
    Def        *Node  // definition: ONAME OTYPE OPACK or OLITERAL
    Label      *Label // corresponding label (ephemeral)
    Block      int32  // blocknumber to catch redeclaration
    Lastlineno int32  // last declaration for diagnostic
    Origpkg    *Pkg   // original package for . import
    Lsym       *obj.LSym
    Fsym       *Sym // funcsym
}

// code to help generate trampoline
// functions for methods on embedded
// subtypes.
// these are approx the same as
// the corresponding adddot routines
// except that they expect to be called
// with unique tasks and they return
// the actual methods.
type Symlink struct {
    field     *Type
    link      *Symlink
    good      bool
    followptr bool
}

type TempVar struct {
    node    *Node
    def     *Flow    // definition of temp var
    use     *Flow    // use list, chained through Flow.data
    merge   *TempVar // merge var with this one
    start   int64    // smallest Prog.pc in live range
    end     int64    // largest Prog.pc in live range
    addr    bool     // address taken - no accurate end
    removed bool     // removed from program
}

type Type struct {
    Etype       EType
    Nointerface bool
    Noalg       bool
    Chan        uint8
    Trecur      uint8 // to detect loops
    Printed     bool
    Embedded    uint8 // TFIELD embedded type
    Funarg      bool  // on TSTRUCT and TFIELD
    Copyany     bool
    Local       bool // created in this file
    Deferwidth  bool
    Broke       bool // broken type definition.
    Isddd       bool // TFIELD is ... argument
    Align       uint8
    Haspointers uint8 // 0 unknown, 1 no, 2 yes

    Nod    *Node // canonical OTYPE node
    Orig   *Type // original type (type literal or predefined type)
    Lineno int

    // TFUNC
    Thistuple int
    Outtuple  int
    Intuple   int
    Outnamed  bool

    Method  *Type
    Xmethod *Type

    Sym    *Sym
    Vargen int32 // unique name for OTYPE/ONAME

    Nname  *Node
    Argwid int64

    // most nodes
    Type  *Type // actual type for TFIELD, element type for TARRAY, TCHAN, TMAP, TPTRxx
    Width int64 // offset in TFIELD, width in all others

    // TFIELD
    Down  *Type   // next struct field, also key type in TMAP
    Outer *Type   // outer struct
    Note  *string // literal string annotation

    // TARRAY
    Bound int64 // negative is dynamic array

    // TMAP
    Bucket *Type // internal type representing a hash bucket
    Hmap   *Type // internal type representing a Hmap (map header object)
    Hiter  *Type // internal type representing hash iterator state
    Map    *Type // link from the above 3 internal types back to the map type.

    Maplineno   int32 // first use of TFORW as map key
    Embedlineno int32 // first use of TFORW as embedded type

    // for TFORW, where to copy the eventual value to
    Copyto []*Node

    Lastfn *Node // for usefield
}

// when a type's width should be known, we call checkwidth
// to compute it.  during a declaration like
//
//     type T *struct { next T }
//
// it is necessary to defer the calculation of the struct width
// until after T has been initialized to be a pointer to that struct.
// similarly, during import processing structs may be used
// before their definition.  in those situations, calling
// defercheckwidth() stops width calculations until
// resumecheckwidth() is called, at which point all the
// checkwidths that were deferred are executed.
// dowidth should only be called when the type's size
// is needed immediately.  checkwidth makes sure the
// size is evaluated eventually.
type TypeList struct {
    t    *Type
    next *TypeList
}

type TypePairList struct {
    t1   *Type
    t2   *Type
    next *TypePairList
}

type Typedef struct {
    Name   string
    Etype  EType
    Sameas EType
}

type Val struct {
    // U contains one of:
    // bool     bool when n.ValCtype() == CTBOOL
    // *Mpint   int when n.ValCtype() == CTINT, rune when n.ValCtype() == CTRUNE
    // *Mpflt   float when n.ValCtype() == CTFLT
    // *Mpcplx  pair of floats when n.ValCtype() == CTCPLX
    // string   string when n.ValCtype() == CTSTR
    // *Nilval  when n.ValCtype() == CTNIL
    U interface{}
}

// A Var represents a single variable that may be stored in a register.
// That variable may itself correspond to a hardware register,
// to represent the use of registers in the unoptimized instruction stream.
type Var struct {
    offset     int64
    node       *Node
    nextinnode *Var
    width      int
    id         int // index in vars
    name       int8
    etype      EType
    addr       int8
}

func Afunclit(a *obj.Addr, n *Node)

// generate:
//     res = &n;
// The generated code checks that the result is not nil.
func Agen(n *Node, res *Node)

// allocate a register (reusing res if possible) and generate
//     a = &n
// The caller must call Regfree(a).
// The generated code checks that the result is not nil.
func Agenr(n *Node, a *Node, res *Node)

func Anyregalloc() bool

// compute total size of f's in/out arguments.
func Argsize(t *Type) int

func AtExit(f func())

func Bconv(xval *Mpint, flag int) string

// Bgen generates code for branches:
//
//     if n == wantTrue {
//         goto to
//     }
func Bgen(n *Node, wantTrue bool, likely int, to *obj.Prog)

// Bitno reports the lowest index of a 1 bit in b.
// It calls Fatalf if there is no 1 bit.
func Bitno(b uint64) int

func Bputname(b *obj.Biobuf, s *obj.LSym)

// Brcom returns !(op).
// For example, Brcom(==) is !=.
func Brcom(op Op) Op

// Brrev returns reverse(op).
// For example, Brrev(<) is >.
func Brrev(op Op) Op

// Bvgen generates code for calculating boolean values:
//     res = n == wantTrue
func Bvgen(n, res *Node, wantTrue bool)

// generate:
//     res = n;
// simplifies and calls Thearch.Gmove.
// if wb is true, need to emit write barriers.
func Cgen(n, res *Node)

// CgenTemp creates a temporary node, assigns n to it, and returns it.
func CgenTemp(n *Node) *Node

// generate:
//     res, resok = x.(T)
// n.Left is x
// n.Type is T
func Cgen_As2dottype(n, res, resok *Node)

func Cgen_as(nl, nr *Node)

func Cgen_as_wb(nl, nr *Node, wb bool)

func Cgen_checknil(n *Node)

// generate:
//     res = iface{typ, data}
// n->left is typ
// n->right is data
func Cgen_eface(n *Node, res *Node)

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
// elements of basic type are also supported.
// nr is nil when assigning a zero value.
func Componentgen(nr, nl *Node) bool

// convert n, if literal, to type t.
// implicit conversion.
func Convlit(np **Node, t *Type)

func Datastring(s string, a *obj.Addr)

// gather series of offsets
// >=0 is direct addressed field
// <0 is pointer to next field (+1)
func Dotoffset(n *Node, oary []int64, nn **Node) int

func Dump(s string, n *Node)

func Dumpit(str string, r0 *Flow, isreg int)

// Fmt "%E": etype
func Econv(et EType) string

// Return 1 if t1 and t2 are identical, following the spec rules.
//
// Any cyclic type must go through a named type, and if one is
// named, it is only identical to the other if they are the same
// pointer (t1 == t2), so there's no chance of chasing cycles
// ad infinitum, so no need for a depth counter.
func Eqtype(t1 *Type, t2 *Type) bool

func Exit(code int)

// Export writes the export data for localpkg to out and returns the number of
// bytes written.
func Export(out *obj.Biobuf, trace bool) int

func Fatalf(fmt_ string, args ...interface{})

func Fconv(fvp *Mpflt, flag int) string

func Fixlargeoffset(n *Node)

func Flowend(graph *Graph)

func Flowstart(firstp *obj.Prog, newData func() interface{}) *Graph

func Flusherrors()

func Gbranch(as int, t *Type, likely int) *obj.Prog

// compile statements
func Genlist(l *NodeList)

func GetReg(r int) int

func Getoutarg(t *Type) **Type

// generate:
//     call f
//     proc=-1    normal call but no return
//     proc=0    normal call
//     proc=1    goroutine run in new proc
//     proc=2    defer call save away stack
//     proc=3    normal call to C pointer (not Go func value)

// generate:
//     call f
//     proc=-1    normal call but no return
//     proc=0    normal call
//     proc=1    goroutine run in new proc
//     proc=2    defer call save away stack
//     proc=3    normal call to C pointer (not Go func value)
func Ginscall(f *Node, proc int)

func Gvardef(n *Node)

// Fmt '%H': NodeList. Flags: all those of %N plus ',': separate with comma's
// instead of semicolons.
func Hconv(l *NodeList, flag int) string

// Igen computes the address &n, stores it in a register r,
// and rewrites a to refer to *r. The chosen r may be the
// stack pointer, it may be borrowed from res, or it may
// be a newly allocated register. The caller must call Regfree(a)
// to free r when the address is no longer needed.
// The generated code ensures that &n is not nil.
func Igen(n *Node, a *Node, res *Node)

// Import populates importpkg from the serialized package data.
func Import(in *obj.Biobuf)

// Is this a 64-bit type?
func Is64(t *Type) bool

func Isconst(n *Node, ct Ctype) bool

func Isfat(t *Type) bool

func Isfixedarray(t *Type) bool

func Isinter(t *Type) bool

// Is this node a memory operand?
func Ismem(n *Node) bool

func Isslice(t *Type) bool

func Istype(t *Type, et EType) bool

// Fmt "%J": Node details.
func Jconv(n *Node, flag int) string

func LOAD(r *Reg, z int) uint64

func Linksym(s *Sym) *obj.LSym

func Lookup(name string) *Sym

func LookupBytes(name []byte) *Sym

func Lookupf(format string, a ...interface{}) *Sym

func Main()

func Mfree(n *Node)

func Mgen(n *Node, n1 *Node, rg *Node)

func Mpcmpfixfix(a, b *Mpint) int

func Mpgetfix(a *Mpint) int64

func Mpmovecfix(a *Mpint, c int64)

func Mpmovecflt(a *Mpflt, c float64)

func Mpmovefixflt(a *Mpflt, b *Mpint)

// shift left by s (or right by -s)
func Mpshiftfix(a *Mpint, s int)

// Naddr rewrites a to refer to n.
// It assumes that a is zeroed on entry.
func Naddr(a *obj.Addr, n *Node)

// Fmt '%N': Nodes.
// Flags: 'l' suffix with "(type %T)" where possible
//       '+h' in debug mode, don't recurse, no multiline output
func Nconv(n *Node, flag int) string

// Is a conversion between t1 and t2 a no-op?
func Noconv(t1 *Type, t2 *Type) bool

func Nod(op Op, nleft *Node, nright *Node) *Node

func Nodbool(b bool) *Node

func Nodconst(n *Node, t *Type, v int64)

func Nodindreg(n *Node, t *Type, r int)

func Nodintconst(v int64) *Node

func Nodreg(n *Node, t *Type, r int)

// p is a call instruction. Does the call fail to return?
func Noreturn(p *obj.Prog) bool

// Fmt "%O":  Node opcodes
func Oconv(o int, flag int) string

func Patch(p *obj.Prog, to *obj.Prog)

func Pkglookup(name string, pkg *Pkg) *Sym

func Prog(as int) *obj.Prog

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

func STORE(r *Reg, z int) uint64

func Samereg(a *Node, b *Node) bool

// Fmt "%S": syms
// Flags:  "%hS" suppresses qualifying with package
func Sconv(s *Sym, flag int) string

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

// iterator to walk a structure declaration
func Structfirst(s *Iter, nn **Type) *Type

func Sysfunc(name string) *Node

// Fmt "%T": types.
// Flags: 'l' print definition, not name
//       'h' omit 'func' and receiver from function types, short type names
//       'u' package name, not prefix (FTypeId mode, sticky)
func Tconv(t *Type, flag int) string

// make a new off the books
func Tempname(nn *Node, t *Type)

// Test all code paths for cmpstackvarlt.
func TestCmpstackvar(t *testing.T)

func TestExprcmp(t *testing.T)

func TestFloatCompare(t *testing.T)

func TestFloatConvert(t *testing.T)

func TestListsort(t *testing.T)

// Make sure "hello world" does not link in all the
// fmt.scanf routines.  See issue 6853.
func TestScanfRemoval(t *testing.T)

func TestSortingByMethodNameAndPackagePath(t *testing.T)

// magic number for unsigned division
// see hacker's delight chapter 10
func Umagic(m *Magic)

func Uniqp(r *Flow) *Flow

func Uniqs(r *Flow) *Flow

// Fmt "%V": Values
func Vconv(v Val, flag int) string

func Warn(fmt_ string, args ...interface{})

func Warnl(line int, fmt_ string, args ...interface{})

func Yyerror(format string, args ...interface{})

func (*Mpflt) String() string

func (*Mpint) String() string

// Bool returns n as an bool.
// n must be an boolean constant.
func (*Node) Bool() bool

// Convconst converts constant node n to type t and
// places the result in con.
func (*Node) Convconst(con *Node, t *Type)

// Int returns n as an int.
// n must be an integer constant.
func (*Node) Int() int64

// IntLiteral returns the Node's literal value as an interger.
func (*Node) IntLiteral() (x int64, ok bool)

func (*Node) Line() string

// Opt returns the optimizer data for the node.
func (*Node) Opt() interface{}

// SetBigInt sets n's value to x.
// n must be an integer constant.
func (*Node) SetBigInt(x *big.Int)

// SetInt sets n's value to i.
// n must be an integer constant.
func (*Node) SetInt(i int64)

// SetOpt sets the optimizer data for the node, which must not have been used
// with SetVal. SetOpt(nil) is ignored for Vals to simplify call sites that are
// clearing Opts.
func (*Node) SetOpt(x interface{})

// SetVal sets the Val for the node, which must not have been used with SetOpt.
func (*Node) SetVal(v Val)

func (*Node) String() string

// Val returns the Val for the node.
func (*Node) Val() Val

func (*NodeList) String() string

func (*Pkg) Lookup(name string) *Sym

func (*Pkg) LookupBytes(name []byte) *Sym

func (*Sym) String() string

func (*Type) String() string

// String returns a space-separated list of the variables represented by bits.
func (Bits) String() string

func (Val) Ctype() Ctype

