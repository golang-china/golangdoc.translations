// +build ingore

package ssa

import (
	"bytes"
	"cmd/internal/obj"
	"cmd/internal/obj/arm"
	"cmd/internal/obj/x86"
	"container/heap"
	"crypto/sha1"
	"fmt"
	"html"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	// When used to lookup up definitions in a sparse tree,
	// these adjustments to a block's entry (+adjust) and
	// exit (-adjust) numbers allow a distinction to be made
	// between assignments (typically branch-dependent
	// conditionals) occurring "before" the block (e.g., as inputs
	// to the block and its phi functions), "within" the block,
	// and "after" the block.
	AdjustBefore = -1 // defined before phi
	AdjustWithin = 0  // defined by phi
	AdjustAfter  = 1  // defined within block
)

const (
	BlockInvalid BlockKind = iota
	BlockAMD64EQ
	BlockAMD64NE
	BlockAMD64LT
	BlockAMD64LE
	BlockAMD64GT
	BlockAMD64GE
	BlockAMD64ULT
	BlockAMD64ULE
	BlockAMD64UGT
	BlockAMD64UGE
	BlockAMD64EQF
	BlockAMD64NEF
	BlockAMD64ORD
	BlockAMD64NAN
	BlockARMEQ
	BlockARMNE
	BlockARMLT
	BlockARMLE
	BlockARMGT
	BlockARMGE
	BlockARMULT
	BlockARMULE
	BlockARMUGT
	BlockARMUGE
	BlockPlain
	BlockIf
	BlockCall
	BlockDefer
	BlockCheck
	BlockRet
	BlockRetJmp
	BlockExit
	BlockFirst
)

const (
	BranchUnlikely = BranchPrediction(-1)
	BranchUnknown  = BranchPrediction(0)
	BranchLikely   = BranchPrediction(+1)
)

const (
	CMPlt = Cmp(-1)
	CMPeq = Cmp(0)
	CMPgt = Cmp(1)
)

// MaxStruct is the maximum number of fields a struct
// can have and still be SSAable.
const MaxStruct = 4

const (
	OpInvalid Op = iota
	OpAMD64ADDSS
	OpAMD64ADDSD
	OpAMD64SUBSS
	OpAMD64SUBSD
	OpAMD64MULSS
	OpAMD64MULSD
	OpAMD64DIVSS
	OpAMD64DIVSD
	OpAMD64MOVSSload
	OpAMD64MOVSDload
	OpAMD64MOVSSconst
	OpAMD64MOVSDconst
	OpAMD64MOVSSloadidx1
	OpAMD64MOVSSloadidx4
	OpAMD64MOVSDloadidx1
	OpAMD64MOVSDloadidx8
	OpAMD64MOVSSstore
	OpAMD64MOVSDstore
	OpAMD64MOVSSstoreidx1
	OpAMD64MOVSSstoreidx4
	OpAMD64MOVSDstoreidx1
	OpAMD64MOVSDstoreidx8
	OpAMD64ADDQ
	OpAMD64ADDL
	OpAMD64ADDQconst
	OpAMD64ADDLconst
	OpAMD64SUBQ
	OpAMD64SUBL
	OpAMD64SUBQconst
	OpAMD64SUBLconst
	OpAMD64MULQ
	OpAMD64MULL
	OpAMD64MULQconst
	OpAMD64MULLconst
	OpAMD64HMULQ
	OpAMD64HMULL
	OpAMD64HMULW
	OpAMD64HMULB
	OpAMD64HMULQU
	OpAMD64HMULLU
	OpAMD64HMULWU
	OpAMD64HMULBU
	OpAMD64AVGQU
	OpAMD64DIVQ
	OpAMD64DIVL
	OpAMD64DIVW
	OpAMD64DIVQU
	OpAMD64DIVLU
	OpAMD64DIVWU
	OpAMD64MODQ
	OpAMD64MODL
	OpAMD64MODW
	OpAMD64MODQU
	OpAMD64MODLU
	OpAMD64MODWU
	OpAMD64ANDQ
	OpAMD64ANDL
	OpAMD64ANDQconst
	OpAMD64ANDLconst
	OpAMD64ORQ
	OpAMD64ORL
	OpAMD64ORQconst
	OpAMD64ORLconst
	OpAMD64XORQ
	OpAMD64XORL
	OpAMD64XORQconst
	OpAMD64XORLconst
	OpAMD64CMPQ
	OpAMD64CMPL
	OpAMD64CMPW
	OpAMD64CMPB
	OpAMD64CMPQconst
	OpAMD64CMPLconst
	OpAMD64CMPWconst
	OpAMD64CMPBconst
	OpAMD64UCOMISS
	OpAMD64UCOMISD
	OpAMD64TESTQ
	OpAMD64TESTL
	OpAMD64TESTW
	OpAMD64TESTB
	OpAMD64TESTQconst
	OpAMD64TESTLconst
	OpAMD64TESTWconst
	OpAMD64TESTBconst
	OpAMD64SHLQ
	OpAMD64SHLL
	OpAMD64SHLQconst
	OpAMD64SHLLconst
	OpAMD64SHRQ
	OpAMD64SHRL
	OpAMD64SHRW
	OpAMD64SHRB
	OpAMD64SHRQconst
	OpAMD64SHRLconst
	OpAMD64SHRWconst
	OpAMD64SHRBconst
	OpAMD64SARQ
	OpAMD64SARL
	OpAMD64SARW
	OpAMD64SARB
	OpAMD64SARQconst
	OpAMD64SARLconst
	OpAMD64SARWconst
	OpAMD64SARBconst
	OpAMD64ROLQconst
	OpAMD64ROLLconst
	OpAMD64ROLWconst
	OpAMD64ROLBconst
	OpAMD64NEGQ
	OpAMD64NEGL
	OpAMD64NOTQ
	OpAMD64NOTL
	OpAMD64BSFQ
	OpAMD64BSFL
	OpAMD64BSFW
	OpAMD64BSRQ
	OpAMD64BSRL
	OpAMD64BSRW
	OpAMD64CMOVQEQconst
	OpAMD64CMOVLEQconst
	OpAMD64CMOVWEQconst
	OpAMD64CMOVQNEconst
	OpAMD64CMOVLNEconst
	OpAMD64CMOVWNEconst
	OpAMD64BSWAPQ
	OpAMD64BSWAPL
	OpAMD64SQRTSD
	OpAMD64SBBQcarrymask
	OpAMD64SBBLcarrymask
	OpAMD64SETEQ
	OpAMD64SETNE
	OpAMD64SETL
	OpAMD64SETLE
	OpAMD64SETG
	OpAMD64SETGE
	OpAMD64SETB
	OpAMD64SETBE
	OpAMD64SETA
	OpAMD64SETAE
	OpAMD64SETEQF
	OpAMD64SETNEF
	OpAMD64SETORD
	OpAMD64SETNAN
	OpAMD64SETGF
	OpAMD64SETGEF
	OpAMD64MOVBQSX
	OpAMD64MOVBQZX
	OpAMD64MOVWQSX
	OpAMD64MOVWQZX
	OpAMD64MOVLQSX
	OpAMD64MOVLQZX
	OpAMD64MOVLconst
	OpAMD64MOVQconst
	OpAMD64CVTTSD2SL
	OpAMD64CVTTSD2SQ
	OpAMD64CVTTSS2SL
	OpAMD64CVTTSS2SQ
	OpAMD64CVTSL2SS
	OpAMD64CVTSL2SD
	OpAMD64CVTSQ2SS
	OpAMD64CVTSQ2SD
	OpAMD64CVTSD2SS
	OpAMD64CVTSS2SD
	OpAMD64PXOR
	OpAMD64LEAQ
	OpAMD64LEAQ1
	OpAMD64LEAQ2
	OpAMD64LEAQ4
	OpAMD64LEAQ8
	OpAMD64MOVBload
	OpAMD64MOVBQSXload
	OpAMD64MOVWload
	OpAMD64MOVWQSXload
	OpAMD64MOVLload
	OpAMD64MOVLQSXload
	OpAMD64MOVQload
	OpAMD64MOVBstore
	OpAMD64MOVWstore
	OpAMD64MOVLstore
	OpAMD64MOVQstore
	OpAMD64MOVOload
	OpAMD64MOVOstore
	OpAMD64MOVBloadidx1
	OpAMD64MOVWloadidx1
	OpAMD64MOVWloadidx2
	OpAMD64MOVLloadidx1
	OpAMD64MOVLloadidx4
	OpAMD64MOVQloadidx1
	OpAMD64MOVQloadidx8
	OpAMD64MOVBstoreidx1
	OpAMD64MOVWstoreidx1
	OpAMD64MOVWstoreidx2
	OpAMD64MOVLstoreidx1
	OpAMD64MOVLstoreidx4
	OpAMD64MOVQstoreidx1
	OpAMD64MOVQstoreidx8
	OpAMD64MOVBstoreconst
	OpAMD64MOVWstoreconst
	OpAMD64MOVLstoreconst
	OpAMD64MOVQstoreconst
	OpAMD64MOVBstoreconstidx1
	OpAMD64MOVWstoreconstidx1
	OpAMD64MOVWstoreconstidx2
	OpAMD64MOVLstoreconstidx1
	OpAMD64MOVLstoreconstidx4
	OpAMD64MOVQstoreconstidx1
	OpAMD64MOVQstoreconstidx8
	OpAMD64DUFFZERO
	OpAMD64MOVOconst
	OpAMD64REPSTOSQ
	OpAMD64CALLstatic
	OpAMD64CALLclosure
	OpAMD64CALLdefer
	OpAMD64CALLgo
	OpAMD64CALLinter
	OpAMD64DUFFCOPY
	OpAMD64REPMOVSQ
	OpAMD64InvertFlags
	OpAMD64LoweredGetG
	OpAMD64LoweredGetClosurePtr
	OpAMD64LoweredNilCheck
	OpAMD64MOVQconvert
	OpAMD64FlagEQ
	OpAMD64FlagLT_ULT
	OpAMD64FlagLT_UGT
	OpAMD64FlagGT_UGT
	OpAMD64FlagGT_ULT
	OpARMADD
	OpARMADDconst
	OpARMMOVWconst
	OpARMCMP
	OpARMMOVWload
	OpARMMOVWstore
	OpARMCALLstatic
	OpARMLessThan
	OpAdd8
	OpAdd16
	OpAdd32
	OpAdd64
	OpAddPtr
	OpAdd32F
	OpAdd64F
	OpSub8
	OpSub16
	OpSub32
	OpSub64
	OpSubPtr
	OpSub32F
	OpSub64F
	OpMul8
	OpMul16
	OpMul32
	OpMul64
	OpMul32F
	OpMul64F
	OpDiv32F
	OpDiv64F
	OpHmul8
	OpHmul8u
	OpHmul16
	OpHmul16u
	OpHmul32
	OpHmul32u
	OpHmul64
	OpHmul64u
	OpAvg64u
	OpDiv8
	OpDiv8u
	OpDiv16
	OpDiv16u
	OpDiv32
	OpDiv32u
	OpDiv64
	OpDiv64u
	OpMod8
	OpMod8u
	OpMod16
	OpMod16u
	OpMod32
	OpMod32u
	OpMod64
	OpMod64u
	OpAnd8
	OpAnd16
	OpAnd32
	OpAnd64
	OpOr8
	OpOr16
	OpOr32
	OpOr64
	OpXor8
	OpXor16
	OpXor32
	OpXor64
	OpLsh8x8
	OpLsh8x16
	OpLsh8x32
	OpLsh8x64
	OpLsh16x8
	OpLsh16x16
	OpLsh16x32
	OpLsh16x64
	OpLsh32x8
	OpLsh32x16
	OpLsh32x32
	OpLsh32x64
	OpLsh64x8
	OpLsh64x16
	OpLsh64x32
	OpLsh64x64
	OpRsh8x8
	OpRsh8x16
	OpRsh8x32
	OpRsh8x64
	OpRsh16x8
	OpRsh16x16
	OpRsh16x32
	OpRsh16x64
	OpRsh32x8
	OpRsh32x16
	OpRsh32x32
	OpRsh32x64
	OpRsh64x8
	OpRsh64x16
	OpRsh64x32
	OpRsh64x64
	OpRsh8Ux8
	OpRsh8Ux16
	OpRsh8Ux32
	OpRsh8Ux64
	OpRsh16Ux8
	OpRsh16Ux16
	OpRsh16Ux32
	OpRsh16Ux64
	OpRsh32Ux8
	OpRsh32Ux16
	OpRsh32Ux32
	OpRsh32Ux64
	OpRsh64Ux8
	OpRsh64Ux16
	OpRsh64Ux32
	OpRsh64Ux64
	OpLrot8
	OpLrot16
	OpLrot32
	OpLrot64
	OpEq8
	OpEq16
	OpEq32
	OpEq64
	OpEqPtr
	OpEqInter
	OpEqSlice
	OpEq32F
	OpEq64F
	OpNeq8
	OpNeq16
	OpNeq32
	OpNeq64
	OpNeqPtr
	OpNeqInter
	OpNeqSlice
	OpNeq32F
	OpNeq64F
	OpLess8
	OpLess8U
	OpLess16
	OpLess16U
	OpLess32
	OpLess32U
	OpLess64
	OpLess64U
	OpLess32F
	OpLess64F
	OpLeq8
	OpLeq8U
	OpLeq16
	OpLeq16U
	OpLeq32
	OpLeq32U
	OpLeq64
	OpLeq64U
	OpLeq32F
	OpLeq64F
	OpGreater8
	OpGreater8U
	OpGreater16
	OpGreater16U
	OpGreater32
	OpGreater32U
	OpGreater64
	OpGreater64U
	OpGreater32F
	OpGreater64F
	OpGeq8
	OpGeq8U
	OpGeq16
	OpGeq16U
	OpGeq32
	OpGeq32U
	OpGeq64
	OpGeq64U
	OpGeq32F
	OpGeq64F
	OpAndB
	OpOrB
	OpEqB
	OpNeqB
	OpNot
	OpNeg8
	OpNeg16
	OpNeg32
	OpNeg64
	OpNeg32F
	OpNeg64F
	OpCom8
	OpCom16
	OpCom32
	OpCom64
	OpCtz16
	OpCtz32
	OpCtz64
	OpClz16
	OpClz32
	OpClz64
	OpBswap32
	OpBswap64
	OpSqrt
	OpPhi
	OpCopy
	OpConvert
	OpConstBool
	OpConstString
	OpConstNil
	OpConst8
	OpConst16
	OpConst32
	OpConst64
	OpConst32F
	OpConst64F
	OpConstInterface
	OpConstSlice
	OpInitMem
	OpArg
	OpAddr
	OpSP
	OpSB
	OpFunc
	OpLoad
	OpStore
	OpMove
	OpZero
	OpClosureCall
	OpStaticCall
	OpDeferCall
	OpGoCall
	OpInterCall
	OpSignExt8to16
	OpSignExt8to32
	OpSignExt8to64
	OpSignExt16to32
	OpSignExt16to64
	OpSignExt32to64
	OpZeroExt8to16
	OpZeroExt8to32
	OpZeroExt8to64
	OpZeroExt16to32
	OpZeroExt16to64
	OpZeroExt32to64
	OpTrunc16to8
	OpTrunc32to8
	OpTrunc32to16
	OpTrunc64to8
	OpTrunc64to16
	OpTrunc64to32
	OpCvt32to32F
	OpCvt32to64F
	OpCvt64to32F
	OpCvt64to64F
	OpCvt32Fto32
	OpCvt32Fto64
	OpCvt64Fto32
	OpCvt64Fto64
	OpCvt32Fto64F
	OpCvt64Fto32F
	OpIsNonNil
	OpIsInBounds
	OpIsSliceInBounds
	OpNilCheck
	OpGetG
	OpGetClosurePtr
	OpArrayIndex
	OpPtrIndex
	OpOffPtr
	OpSliceMake
	OpSlicePtr
	OpSliceLen
	OpSliceCap
	OpComplexMake
	OpComplexReal
	OpComplexImag
	OpStringMake
	OpStringPtr
	OpStringLen
	OpIMake
	OpITab
	OpIData
	OpStructMake0
	OpStructMake1
	OpStructMake2
	OpStructMake3
	OpStructMake4
	OpStructSelect
	OpStoreReg
	OpLoadReg
	OpFwdRef
	OpUnknown
	OpVarDef
	OpVarKill
	OpVarLive
	OpKeepAlive
)

const (
	ScorePhi = iota // towards top of block
	ScoreVarDef
	ScoreMemory
	ScoreDefault
	ScoreFlags
	ScoreControl // towards bottom of block
)

var BuildDebug int

var BuildStats int

var BuildTest int

// Debug output
var IntrinsicsDebug int

var IntrinsicsDisable bool

var (
	TypeInvalid = &CompilerType{Name: "invalid"}
	TypeMem     = &CompilerType{Name: "mem", Memory: true}
	TypeFlags   = &CompilerType{Name: "flags", Flags: true}
	TypeVoid    = &CompilerType{Name: "void", Void: true}
	TypeInt128  = &CompilerType{Name: "int128", size: 16, Int128: true}
)

// ArgSymbol is an aux value that encodes an argument or result
// variable's constant offset from FP (FP = SP + framesize).
type ArgSymbol struct {
	Typ  Type   // Go type
	Node GCNode // A *gc.Node referring to the argument/result variable.
}

// AutoSymbol is an aux value that encodes a local variable's
// constant offset from SP.
type AutoSymbol struct {
	Typ  Type   // Go type
	Node GCNode // A *gc.Node referring to a local (auto) variable.
}

// Block represents a basic block in the control flow graph of a function.
type Block struct {
	// A unique identifier for the block. The system will attempt to allocate
	// these IDs densely, but no guarantees.
	ID ID

	// Line number for block's control operation
	Line int32

	// The kind of block this is.
	Kind BlockKind

	// Likely direction for branches.
	// If BranchLikely, Succs[0] is the most likely branch taken.
	// If BranchUnlikely, Succs[1] is the most likely branch taken.
	// Ignored if len(Succs) < 2.
	// Fatal if not BranchUnknown and len(Succs) > 2.
	Likely BranchPrediction

	// After flagalloc, records whether flags are live at the end of the block.
	FlagsLiveAtEnd bool

	// Subsequent blocks, if any. The number and order depend on the block kind.
	Succs []Edge

	// Inverse of successors. The order is significant to Phi nodes in the
	// block. TODO: predecessors is a pain to maintain. Can we somehow order phi
	// arguments by block id and have this field computed explicitly when
	// needed?
	Preds []Edge

	// A value that determines how the block is exited. Its value depends on the
	// kind of the block. For instance, a BlockIf has a boolean control value
	// and BlockExit has a memory control value.
	Control *Value

	// Auxiliary info for the block. Its value depends on the Kind.
	Aux interface{}

	// The unordered set of Values that define the operation of this block. The
	// list must include the control value, if any. (TODO: need this last
	// condition?) After the scheduling pass, this list is ordered.
	Values []*Value

	// The containing function
	Func *Func
}

// 	kind           control    successors
// ------------------------------------------
//
// 	 Exit        return mem                []
// 	Plain               nil            [next]
// 	   If   a boolean Value      [then, else]
// 	 Call               mem  [nopanic, panic]  (control opcode should be OpCall or OpStaticCall)
type BlockKind int8

type BranchPrediction int8

// Cmp is a comparison between values a and b.
// -1 if a < b
//  0 if a == b
//  1 if a > b
type Cmp int8

// Special compiler-only types.
type CompilerType struct {
	Name   string
	Memory bool
	Flags  bool
	Void   bool
	Int128 bool
}

type Config struct {
	IntSize int64       // 4 or 8
	PtrSize int64       // 4 or 8
	HTML    *HTMLWriter // html writer, for debugging
}

// Edge represents a CFG edge.
// Example edges for b branching to either c or d.
// (c and d have other predecessors.)
//   b.Succs = [{c,3}, {d,1}]
//   c.Preds = [?, ?, ?, {b,0}]
//   d.Preds = [?, {b,1}, ?]
// These indexes allow us to edit the CFG in constant time.
// In addition, it informs phi ops in degenerate cases like:
// b:
//    if k then c else c
// c:
//    v = Phi(x, y)
// Then the indexes tell you whether x is chosen from
// the if or else branch from b.
//   b.Succs = [{c,0},{c,1}]
//   c.Preds = [{b,0},{b,1}]
// means x is chosen if k is true.
type Edge struct {
}

// ExternSymbol is an aux value that encodes a variable's
// constant offset from the static base pointer.
type ExternSymbol struct {
	Typ Type         // Go type
	Sym fmt.Stringer // A *gc.Sym referring to a global variable
}

type Frontend interface {
	TypeSource
	Logger

	// StringData returns a symbol pointing to the given string's contents.
	StringData(string)interface{} // returns *gc.Sym

	// Auto returns a Node for an auto variable of the given type.
	// The SSA compiler uses this function to allocate space for spills.
	Auto(Type)GCNode

	// Given the name for a compound type, returns the name we should use
	// for the parts of that compound type.
	SplitString(LocalSlot) (LocalSlot, LocalSlot)
	SplitInterface(LocalSlot) (LocalSlot, LocalSlot)
	SplitSlice(LocalSlot) (LocalSlot, LocalSlot, LocalSlot)
	SplitComplex(LocalSlot) (LocalSlot, LocalSlot)
	SplitStruct(LocalSlot, int)LocalSlot

	// Line returns a string describing the given line number.
	Line(int32)string
}

// A Func represents a Go func declaration (or function literal) and
// its body. This package compiles each Func independently.
type Func struct {
	Config     *Config     // architecture information
	Name       string      // e.g. bytes·Compare
	Type       Type        // type signature of the function.
	StaticData interface{} // associated static data, untouched by the ssa package
	Blocks     []*Block    // unordered set of all basic blocks (note: not indexable by ID)
	Entry      *Block      // the entry basic block

	// when register allocation is done, maps value ids to locations
	RegAlloc []Location

	// map from LocalSlot to set of Values that we want to store in that slot.
	NamedValues map[LocalSlot][]*Value

	// Names is a copy of NamedValues.Keys. We keep a separate list
	// of keys to make iteration order deterministic.
	Names []LocalSlot
}

// interface used to hold *gc.Node. We'd use *gc.Node directly but
// that would lead to an import cycle.
type GCNode interface {
	Typ()Type
	String()string
}

type HTMLWriter struct {
	Logger
	*os.File
}

type ID int32

// A LocalSlot is a location in the stack frame.
// It is (possibly a subpiece of) a PPARAM, PPARAMOUT, or PAUTO ONAME node.
type LocalSlot struct {
	N    GCNode // an ONAME *gc.Node representing a variable on the stack
	Type Type   // type of slot
	Off  int64  // offset of slot in N
}

// A place that an ssa variable can reside.
type Location interface {
	Name()string // name to use in assembly templates: %rax, 16(%rsp), ...
}

type Logger interface {
	// Logf logs a message from the compiler.
	Logf(string, ...interface{})

	// Log returns true if logging is not a no-op
	// some logging calls account for more than a few heap allocations.
	Log()bool

	// Fatal reports a compiler error and exits.
	Fatalf(line int32, msg string, args ...interface{})

	// Unimplemented reports that the function cannot be compiled.
	// It will be removed once SSA work is complete.
	Unimplementedf(line int32, msg string, args ...interface{})

	// Warnl writes compiler messages in the form expected by "errorcheck" tests
	Warnl(line int32, fmt_ string, args ...interface{})

	// Fowards the Debug_checknil flag from gc
	Debug_checknil()bool
}

// An Op encodes the specific operation that a Value performs. Opcodes'
// semantics can be modified by the type and aux fields of the Value. For
// instance, OpAdd can be 32 or 64 bit, signed or unsigned, float or complex,
// depending on Value.Type. Semantics of each op are described in the opcode
// files in gen/*Ops.go. There is one file for generic
// (architecture-independent) ops and one file for each architecture.
type Op int32

// RBTint32 is a red-black tree with data stored at internal nodes,
// following Tarjan, Data Structures and Network Algorithms,
// pp 48-52, using explicit rank instead of red and black.
// Deletion is not yet implemented because it is not yet needed.
// Extra operations glb, lub, glbEq, lubEq are provided for
// use in sparse lookup algorithms.
type RBTint32 struct {
}

// A Register is a machine register, like %rax.
// They are numbered densely from 0 (for each architecture).
type Register struct {
	Num int32
}

// A SparseTree is a tree of Blocks.
// It allows rapid ancestor queries,
// such as whether one block dominates another.
type SparseTree []SparseTreeNode

// A SparseTreeHelper contains indexing and allocation data
// structures common to a collection of SparseTreeMaps, as well
// as exposing some useful control-flow-related data to other
// packages, such as gc.
type SparseTreeHelper struct {
	Sdom   []SparseTreeNode // indexed by block.ID
	Po     []*Block         // exported data; the blocks, in a post-order
	Dom    []*Block         // exported data; the dominator of this block.
	Ponums []int32          // exported data; Po[Ponums[b.ID]] == b; the index of b in Po
}

// A SparseTreeMap encodes a subset of nodes within a tree
// used for sparse-ancestor queries.
//
// Combined with a SparseTreeHelper, this supports an Insert
// to add a tree node to the set and a Find operation to locate
// the nearest tree ancestor of a given node such that the
// ancestor is also in the set.
//
// Given a set of blocks {B1, B2, B3} within the dominator tree, established
// by stm.Insert()ing B1, B2, B3, etc, a query at block B
// (performed with stm.Find(stm, B, adjust, helper))
// will return the member of the set that is the nearest strict
// ancestor of B within the dominator tree, or nil if none exists.
// The expected complexity of this operation is the log of the size
// the set, given certain assumptions about sparsity (the log complexity
// could be guaranteed with additional data structures whose constant-
// factor overhead has not yet been justified.)
//
// The adjust parameter allows positioning of the insertion
// and lookup points within a block -- one of
// AdjustBefore, AdjustWithin, AdjustAfter,
// where lookups at AdjustWithin can find insertions at
// AdjustBefore in the same block, and lookups at AdjustAfter
// can find insertions at either AdjustBefore or AdjustWithin
// in the same block.  (Note that this assumes a gappy numbering
// such that exit number or exit number is separated from its
// nearest neighbor by at least 3).
//
// The Sparse Tree lookup algorithm is described by
// Paul F. Dietz. Maintaining order in a linked list. In
// Proceedings of the Fourteenth Annual ACM Symposium on
// Theory of Computing, pages 122–127, May 1982.
// and by
// Ben Wegbreit. Faster retrieval from context trees.
// Communications of the ACM, 19(9):526–529, September 1976.

// A SparseTreeMap encodes a subset of nodes within a tree used for
// sparse-ancestor queries.
//
// Combined with a SparseTreeHelper, this supports an Insert to add a tree node
// to the set and a Find operation to locate the nearest tree ancestor of a
// given node such that the ancestor is also in the set.
//
// Given a set of blocks {B1, B2, B3} within the dominator tree, established by
// stm.Insert()ing B1, B2, B3, etc, a query at block B (performed with
// stm.Find(stm, B, adjust, helper)) will return the member of the set that is
// the nearest strict ancestor of B within the dominator tree, or nil if none
// exists. The expected complexity of this operation is the log of the size the
// set, given certain assumptions about sparsity (the log complexity could be
// guaranteed with additional data structures whose constant- factor overhead
// has not yet been justified.)
//
// The adjust parameter allows positioning of the insertion and lookup points
// within a block -- one of AdjustBefore, AdjustWithin, AdjustAfter, where
// lookups at AdjustWithin can find insertions at AdjustBefore in the same
// block, and lookups at AdjustAfter can find insertions at either AdjustBefore
// or AdjustWithin in the same block. (Note that this assumes a gappy numbering
// such that exit number or exit number is separated from its nearest neighbor
// by at least 3).
//
// The Sparse Tree lookup algorithm is described by Paul F. Dietz. Maintaining
// order in a linked list. In Proceedings of the Fourteenth Annual ACM Symposium
// on Theory of Computing, pages 122–127, May 1982. and by Ben Wegbreit. Faster
// retrieval from context trees. Communications of the ACM, 19(9):526–529,
// September 1976.
type SparseTreeMap RBTint32

type SparseTreeNode struct {
}

// A type interface used to import cmd/internal/gc:Type
// Type instances are not guaranteed to be canonical.
type Type interface {
	Size()int64 // return the size in bytes
	Alignment()int64
	IsBoolean()bool // is a named or unnamed boolean type
	IsInteger()bool //  ... ditto for the others
	IsSigned()bool
	IsFloat()bool
	IsComplex()bool
	IsPtrShaped()bool
	IsString()bool
	IsSlice()bool
	IsArray()bool
	IsStruct()bool
	IsInterface()bool
	IsMemory()bool // special ssa-package-only types
	IsFlags()bool
	IsVoid()bool
	ElemType()Type         // given []T or *T or [n]T, return T
	PtrTo()Type            // given T, return *T
	NumFields()int         // # of fields of a struct
	FieldType(i int)Type   // type of ith field of the struct
	FieldOff(i int)int64   // offset of ith field of the struct
	FieldName(i int)string // name of ith field of the struct
	NumElem()int64         // # of elements of an array
	String()string
	SimpleString()string // a coarser generic description of T, e.g. T's underlying type
	Compare(Type)Cmp     // compare types, returning one of CMPlt, CMPeq, CMPgt.
}

type TypeSource interface {
	TypeBool()Type
	TypeInt8()Type
	TypeInt16()Type
	TypeInt32()Type
	TypeInt64()Type
	TypeUInt8()Type
	TypeUInt16()Type
	TypeUInt32()Type
	TypeUInt64()Type
	TypeInt()Type
	TypeFloat32()Type
	TypeFloat64()Type
	TypeUintptr()Type
	TypeString()Type
	TypeBytePtr()Type // TODO: use unsafe.Pointer instead?
	CanSSA(t Type)bool
}

// A ValAndOff is used by the several opcodes. It holds
// both a value and a pointer offset.
// A ValAndOff is intended to be encoded into an AuxInt field.
// The zero ValAndOff encodes a value of 0 and an offset of 0.
// The high 32 bits hold a value.
// The low 32 bits hold a pointer offset.
type ValAndOff int64

type ValHeap struct {
}

// A Value represents a value in the SSA representation of the program. The ID
// and Type fields must not be modified. The remainder may be modified if they
// preserve the value of the Value (e.g. changing a (mul 2 x) to an (add x x)).
type Value struct {
	// A unique identifier for the value. For performance we allocate these IDs
	// densely starting at 1. There is no guarantee that there won't be
	// occasional holes, though.
	ID ID

	// The operation that computes this value. See op.go.
	Op Op

	// The type of this value. Normally this will be a Go type, but there
	// are a few other pseudo-types, see type.go.
	Type Type

	// Auxiliary info for this value. The type of this information depends on
	// the opcode and type. AuxInt is used for integer values, Aux is used for
	// other values.
	AuxInt int64
	Aux    interface{}

	// Arguments of this value
	Args []*Value

	// Containing basic block
	Block *Block

	// Source line number
	Line int32

	// Use count. Each appearance in Value.Args and Block.Control counts once.
	Uses int32
}

// Compile is the main entry point for this package. Compile modifies f so that
// on return:
//
// 	· all Values in f map to 0 or 1 assembly instructions of the target architecture
// 	· the order of f.Blocks is the order to emit the Blocks
// 	· the order of b.Values is the order to emit the Values in each Block
// 	· f has a non-nil regAlloc field
func Compile(f *Func)

// NewConfig returns a new configuration object for the given architecture.
func NewConfig(arch string, fe Frontend, ctxt *obj.Link, optimize bool) *Config

func NewHTMLWriter(path string, logger Logger, funcname string) *HTMLWriter

// NewSparseTreeHelper returns a SparseTreeHelper for use
// in the gc package, for example in phi-function placement.
func NewSparseTreeHelper(f *Func) *SparseTreeHelper

// PhaseOption sets the specified flag in the specified ssa phase, returning
// empty string if this was successful or a string explaining the error if it
// was not. A version of the phase name with "_" replaced by " " is also checked
// for a match. If the phase name begins a '~' then the rest of the
// underscores-replaced-with-blanks version is used as a regular expression to
// match the phase name(s).
//
// Special cases that have turned out to be useful:
//
// 	ssa/check/on enables checking after each phase
// 	ssa/all/time enables time reporting for all phases
//
// See gc/lex.go for dissection of the option string. Example uses:
//
// GO_GCFLAGS=-d=ssa/generic_cse/time,ssa/generic_cse/stats,ssa/generic_cse/debug=3
// ./make.bash
//
// BOOT_GO_GCFLAGS=-d='ssa/~^.*scc$/off' GO_GCFLAGS='-d=ssa/~^.*scc$/off'
// ./make.bash
func PhaseOption(phase, flag string, val int) string

// StructMakeOp returns the opcode to construct a struct with the
// given number of fields.
func StructMakeOp(nf int) Op

func (s *ArgSymbol) String() string

func (s *AutoSymbol) String() string

// AddEdgeTo adds an edge from block b to block c. Used during building of the
// SSA graph; do not use on an already-completed SSA graph.
func (b *Block) AddEdgeTo(c *Block)

func (b *Block) Fatalf(msg string, args ...interface{})

func (b *Block) HTML() string

func (b *Block) Log() bool

func (b *Block) Logf(msg string, args ...interface{})

func (b *Block) LongHTML() string

// long form print
func (b *Block) LongString() string

// NewValue0 returns a new value in the block with no arguments and zero aux
// values.
func (b *Block) NewValue0(line int32, op Op, t Type) *Value

// NewValue returns a new value in the block with no arguments and an aux value.
func (b *Block) NewValue0A(line int32, op Op, t Type, aux interface{}) *Value

// NewValue returns a new value in the block with no arguments and an auxint
// value.
func (b *Block) NewValue0I(line int32, op Op, t Type, auxint int64) *Value

// NewValue returns a new value in the block with no arguments and both an
// auxint and aux values.
func (b *Block) NewValue0IA(line int32, op Op, t Type, auxint int64, aux interface{}) *Value

// NewValue1 returns a new value in the block with one argument and zero aux
// values.
func (b *Block) NewValue1(line int32, op Op, t Type, arg *Value) *Value

// NewValue1A returns a new value in the block with one argument and an aux
// value.
func (b *Block) NewValue1A(line int32, op Op, t Type, aux interface{}, arg *Value) *Value

// NewValue1I returns a new value in the block with one argument and an auxint
// value.
func (b *Block) NewValue1I(line int32, op Op, t Type, auxint int64, arg *Value) *Value

// NewValue1IA returns a new value in the block with one argument and both an
// auxint and aux values.
func (b *Block) NewValue1IA(line int32, op Op, t Type, auxint int64, aux interface{}, arg *Value) *Value

// NewValue2 returns a new value in the block with two arguments and zero aux
// values.
func (b *Block) NewValue2(line int32, op Op, t Type, arg0, arg1 *Value) *Value

// NewValue2I returns a new value in the block with two arguments and an auxint
// value.
func (b *Block) NewValue2I(line int32, op Op, t Type, auxint int64, arg0, arg1 *Value) *Value

// NewValue3 returns a new value in the block with three arguments and zero aux
// values.
func (b *Block) NewValue3(line int32, op Op, t Type, arg0, arg1, arg2 *Value) *Value

// NewValue3I returns a new value in the block with three arguments and an
// auxint value.
func (b *Block) NewValue3I(line int32, op Op, t Type, auxint int64, arg0, arg1, arg2 *Value) *Value

func (b *Block) SetControl(v *Value)

// short form print
func (b *Block) String() string

func (b *Block) Unimplementedf(msg string, args ...interface{})

func (t *CompilerType) Alignment() int64

func (t *CompilerType) Compare(u Type) Cmp

func (t *CompilerType) ElemType() Type

func (t *CompilerType) FieldName(i int) string

func (t *CompilerType) FieldOff(i int) int64

func (t *CompilerType) FieldType(i int) Type

func (t *CompilerType) IsArray() bool

func (t *CompilerType) IsBoolean() bool

func (t *CompilerType) IsComplex() bool

func (t *CompilerType) IsFlags() bool

func (t *CompilerType) IsFloat() bool

func (t *CompilerType) IsInteger() bool

func (t *CompilerType) IsInterface() bool

func (t *CompilerType) IsMemory() bool

func (t *CompilerType) IsPtrShaped() bool

func (t *CompilerType) IsSigned() bool

func (t *CompilerType) IsSlice() bool

func (t *CompilerType) IsString() bool

func (t *CompilerType) IsStruct() bool

func (t *CompilerType) IsVoid() bool

func (t *CompilerType) NumElem() int64

func (t *CompilerType) NumFields() int

func (t *CompilerType) PtrTo() Type

func (t *CompilerType) SimpleString() string

func (t *CompilerType) Size() int64

func (t *CompilerType) String() string

func (c *Config) DebugHashMatch(evname, name string) bool

func (c *Config) DebugNameMatch(evname, name string) bool

func (c *Config) Debug_checknil() bool

func (c *Config) Fatalf(line int32, msg string, args ...interface{})

func (c *Config) Frontend() Frontend

func (c *Config) Log() bool

func (c *Config) Logf(msg string, args ...interface{})

// NewFunc returns a new, empty function object.
// Caller must call f.Free() before calling NewFunc again.
func (c *Config) NewFunc() *Func

func (c *Config) SparsePhiCutoff() uint64

func (c *Config) Unimplementedf(line int32, msg string, args ...interface{})

func (c *Config) Warnl(line int32, msg string, args ...interface{})

func (s *ExternSymbol) String() string

// ConstInt returns an int constant representing its argument.
func (f *Func) ConstBool(line int32, t Type, c bool) *Value

func (f *Func) ConstEmptyString(line int32, t Type) *Value

func (f *Func) ConstFloat32(line int32, t Type, c float64) *Value

func (f *Func) ConstFloat64(line int32, t Type, c float64) *Value

func (f *Func) ConstInt16(line int32, t Type, c int16) *Value

func (f *Func) ConstInt32(line int32, t Type, c int32) *Value

func (f *Func) ConstInt64(line int32, t Type, c int64) *Value

func (f *Func) ConstInt8(line int32, t Type, c int8) *Value

func (f *Func) ConstInterface(line int32, t Type) *Value

func (f *Func) ConstNil(line int32, t Type) *Value

func (f *Func) ConstSlice(line int32, t Type) *Value

func (f *Func) Fatalf(msg string, args ...interface{})

func (f *Func) Free()

func (f *Func) HTML() string

func (f *Func) Log() bool

// logPassStat writes a string key and int value as a warning in a
// tab-separated format easily handled by spreadsheets or awk.
// file names, lines, and function names are included to provide enough (?)
// context to allow item-by-item comparisons across runs.
// For example:
// awk 'BEGIN {FS="\t"} $3~/TIME/{sum+=$4} END{print "t(ns)=",sum}' t.log
func (f *Func) LogStat(key string, args ...interface{})

func (f *Func) Logf(msg string, args ...interface{})

// newBlock allocates a new Block of the given kind and places it at the end of
// f.Blocks.
func (f *Func) NewBlock(kind BlockKind) *Block

// NumBlocks returns an integer larger than the id of any Block in the Func.
func (f *Func) NumBlocks() int

// NumValues returns an integer larger than the id of any Value in the Func.
func (f *Func) NumValues() int

func (f *Func) String() string

func (f *Func) Unimplementedf(msg string, args ...interface{})

func (w *HTMLWriter) Close()

func (w *HTMLWriter) Printf(msg string, v ...interface{})

// WriteColumn writes raw HTML in a column headed by title.
// It is intended for pre- and post-compilation log output.
func (w *HTMLWriter) WriteColumn(title string, html string)

// WriteFunc writes f in a column headed by title.
func (w *HTMLWriter) WriteFunc(title string, f *Func)

func (w *HTMLWriter) WriteString(s string)

// Find returns the data associated with key in the tree, or
// nil if key is not in the tree.
func (t *RBTint32) Find(key int32) interface{}

// Glb returns the greatest-lower-bound-exclusive of x and its associated
// data.  If x has no glb in the tree, then (0, nil) is returned.
func (t *RBTint32) Glb(x int32) (k int32, d interface{})

// GlbEq returns the greatest-lower-bound-inclusive of x and its associated
// data.  If x has no glbEQ in the tree, then (0, nil) is returned.
func (t *RBTint32) GlbEq(x int32) (k int32, d interface{})

// Insert adds key to the tree and associates key with data.
// If key was already in the tree, it updates the associated data.
// Insert returns the previous data associated with key,
// or nil if key was not present.
// Insert panics if data is nil.
func (t *RBTint32) Insert(key int32, data interface{}) interface{}

// IsEmpty reports whether t is empty.
func (t *RBTint32) IsEmpty() bool

// IsSingle reports whether t is a singleton (leaf).
func (t *RBTint32) IsSingle() bool

// Lub returns the least-upper-bound-exclusive of x and its associated
// data.  If x has no lub in the tree, then (0, nil) is returned.
func (t *RBTint32) Lub(x int32) (k int32, d interface{})

// LubEq returns the least-upper-bound-inclusive of x and its associated
// data.  If x has no lubEq in the tree, then (0, nil) is returned.
func (t *RBTint32) LubEq(x int32) (k int32, d interface{})

// Max returns the maximum element of t and its associated data.
// If t is empty, then (0, nil) is returned.
func (t *RBTint32) Max() (k int32, d interface{})

// Min returns the minimum element of t and its associated data.
// If t is empty, then (0, nil) is returned.
func (t *RBTint32) Min() (k int32, d interface{})

func (t *RBTint32) String() string

// VisitInOrder applies f to the key and data pairs in t,
// with keys ordered from smallest to largest.
func (t *RBTint32) VisitInOrder(f func(int32, interface{}))

func (r *Register) Name() string

func (h *SparseTreeHelper) NewTree() *SparseTreeMap

// Find returns the definition visible from block b, or nil if none can be
// found. Adjust indicates where the block should be searched. AdjustBefore
// searches before the phi functions of b. AdjustWithin searches starting at the
// phi functions of b. AdjustAfter searches starting at the exit from the block,
// including normal within-block definitions.
//
// Note that Finds are properly nested with Inserts: m.Insert(b, a) followed by
// m.Find(b, a) will not return the result of the insert, but m.Insert(b,
// AdjustBefore) followed by m.Find(b, AdjustWithin) will.
//
// Another way to think of this is that Find searches for inputs, Insert defines
// outputs.
func (m *SparseTreeMap) Find(b *Block, adjust int32, helper *SparseTreeHelper) interface{}

// Insert creates a definition within b with data x. adjust indicates where in
// the block should be inserted: AdjustBefore means defined at a phi function
// (visible Within or After in the same block) AdjustWithin means defined within
// the block (visible After in the same block) AdjustAfter means after the block
// (visible within child blocks)
func (m *SparseTreeMap) Insert(b *Block, adjust int32, x interface{}, helper *SparseTreeHelper)

func (m *SparseTreeMap) String() string

func (s *SparseTreeNode) Entry() int32

func (s *SparseTreeNode) Exit() int32

func (s *SparseTreeNode) String() string

func (h *ValHeap) Pop() interface{}

func (h *ValHeap) Push(x interface{})

func (v *Value) AddArg(w *Value)

func (v *Value) AddArgs(a ...*Value)

func (v *Value) AuxFloat() float64

func (v *Value) AuxInt16() int16

func (v *Value) AuxInt32() int32

func (v *Value) AuxInt8() int8

func (v *Value) AuxValAndOff() ValAndOff

func (v *Value) Fatalf(msg string, args ...interface{})

func (v *Value) HTML() string

func (v *Value) Log() bool

func (v *Value) Logf(msg string, args ...interface{})

func (v *Value) LongHTML() string

// long form print.  v# = opcode <type> [aux] args [: reg]
func (v *Value) LongString() string

func (v *Value) RemoveArg(i int)

func (v *Value) SetArg(i int, w *Value)

func (v *Value) SetArgs1(a *Value)

func (v *Value) SetArgs2(a *Value, b *Value)

// short form print. Just v#.
func (v *Value) String() string

func (v *Value) Unimplementedf(msg string, args ...interface{})

func (k BlockKind) String() string

func (e Edge) Block() *Block

func (s LocalSlot) Name() string

func (o Op) Asm() obj.As

func (o Op) String() string

// Child returns a child of x in the dominator tree, or
// nil if there are none. The choice of first child is
// arbitrary but repeatable.
func (t SparseTree) Child(x *Block) *Block

// Sibling returns a sibling of x in the dominator tree (i.e.,
// a node with the same immediate dominator) or nil if there
// are no remaining siblings in the arbitrary but repeatable
// order chosen. Because the Child-Sibling order is used
// to assign entry and exit numbers in the treewalk, those
// numbers are also consistent with this order (i.e.,
// Sibling(x) has entry number larger than x's exit number).
func (t SparseTree) Sibling(x *Block) *Block

func (x ValAndOff) Int64() int64

func (x ValAndOff) Off() int64

func (x ValAndOff) String() string

func (x ValAndOff) Val() int64

func (h ValHeap) Len() int

func (h ValHeap) Less(i, j int) bool

func (h ValHeap) Swap(i, j int)

