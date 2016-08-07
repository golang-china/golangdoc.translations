// +build ingore

package ppc64 // import "cmd/internal/obj/ppc64"

import (
    "cmd/internal/obj"
    "encoding/binary"
    "fmt"
    "log"
    "math"
    "sort"
)

const (
    AADD = obj.ABasePPC64 + obj.A_ARCHSPECIFIC + iota
    AADDCC
    AADDV
    AADDVCC
    AADDC
    AADDCCC
    AADDCV
    AADDCVCC
    AADDME
    AADDMECC
    AADDMEVCC
    AADDMEV
    AADDE
    AADDECC
    AADDEVCC
    AADDEV
    AADDZE
    AADDZECC
    AADDZEVCC
    AADDZEV
    AAND
    AANDCC
    AANDN
    AANDNCC
    ABC
    ABCL
    ABEQ
    ABGE
    ABGT
    ABLE
    ABLT
    ABNE
    ABVC
    ABVS
    ACMP
    ACMPU
    ACNTLZW
    ACNTLZWCC
    ACRAND
    ACRANDN
    ACREQV
    ACRNAND
    ACRNOR
    ACROR
    ACRORN
    ACRXOR
    ADIVW
    ADIVWCC
    ADIVWVCC
    ADIVWV
    ADIVWU
    ADIVWUCC
    ADIVWUVCC
    ADIVWUV
    AEQV
    AEQVCC
    AEXTSB
    AEXTSBCC
    AEXTSH
    AEXTSHCC
    AFABS
    AFABSCC
    AFADD
    AFADDCC
    AFADDS
    AFADDSCC
    AFCMPO
    AFCMPU
    AFCTIW
    AFCTIWCC
    AFCTIWZ
    AFCTIWZCC
    AFDIV
    AFDIVCC
    AFDIVS
    AFDIVSCC
    AFMADD
    AFMADDCC
    AFMADDS
    AFMADDSCC
    AFMOVD
    AFMOVDCC
    AFMOVDU
    AFMOVS
    AFMOVSU
    AFMSUB
    AFMSUBCC
    AFMSUBS
    AFMSUBSCC
    AFMUL
    AFMULCC
    AFMULS
    AFMULSCC
    AFNABS
    AFNABSCC
    AFNEG
    AFNEGCC
    AFNMADD
    AFNMADDCC
    AFNMADDS
    AFNMADDSCC
    AFNMSUB
    AFNMSUBCC
    AFNMSUBS
    AFNMSUBSCC
    AFRSP
    AFRSPCC
    AFSUB
    AFSUBCC
    AFSUBS
    AFSUBSCC
    AMOVMW
    ALSW
    ALWAR
    AMOVWBR
    AMOVB
    AMOVBU
    AMOVBZ
    AMOVBZU
    AMOVH
    AMOVHBR
    AMOVHU
    AMOVHZ
    AMOVHZU
    AMOVW
    AMOVWU
    AMOVFL
    AMOVCRFS
    AMTFSB0
    AMTFSB0CC
    AMTFSB1
    AMTFSB1CC
    AMULHW
    AMULHWCC
    AMULHWU
    AMULHWUCC
    AMULLW
    AMULLWCC
    AMULLWVCC
    AMULLWV
    ANAND
    ANANDCC
    ANEG
    ANEGCC
    ANEGVCC
    ANEGV
    ANOR
    ANORCC
    AOR
    AORCC
    AORN
    AORNCC
    AREM
    AREMCC
    AREMV
    AREMVCC
    AREMU
    AREMUCC
    AREMUV
    AREMUVCC
    ARFI
    ARLWMI
    ARLWMICC
    ARLWNM
    ARLWNMCC
    ASLW
    ASLWCC
    ASRW
    ASRAW
    ASRAWCC
    ASRWCC
    ASTSW
    ASTWCCC
    ASUB
    ASUBCC
    ASUBVCC
    ASUBC
    ASUBCCC
    ASUBCV
    ASUBCVCC
    ASUBME
    ASUBMECC
    ASUBMEVCC
    ASUBMEV
    ASUBV
    ASUBE
    ASUBECC
    ASUBEV
    ASUBEVCC
    ASUBZE
    ASUBZECC
    ASUBZEVCC
    ASUBZEV
    ASYNC
    AXOR
    AXORCC

    ADCBF
    ADCBI
    ADCBST
    ADCBT
    ADCBTST
    ADCBZ
    AECIWX
    AECOWX
    AEIEIO
    AICBI
    AISYNC
    APTESYNC
    ATLBIE
    ATLBIEL
    ATLBSYNC
    ATW

    ASYSCALL
    AWORD

    ARFCI

    /* optional on 32-bit */
    AFRES
    AFRESCC
    AFRSQRTE
    AFRSQRTECC
    AFSEL
    AFSELCC
    AFSQRT
    AFSQRTCC
    AFSQRTS
    AFSQRTSCC

    ACNTLZD
    ACNTLZDCC
    ACMPW /* CMP with L=0 */
    ACMPWU
    ADIVD
    ADIVDCC
    ADIVDVCC
    ADIVDV
    ADIVDU
    ADIVDUCC
    ADIVDUVCC
    ADIVDUV
    AEXTSW
    AEXTSWCC
    /* AFCFIW; AFCFIWCC */
    AFCFID
    AFCFIDCC
    AFCTID
    AFCTIDCC
    AFCTIDZ
    AFCTIDZCC
    ALDAR
    AMOVD
    AMOVDU
    AMOVWZ
    AMOVWZU
    AMULHD
    AMULHDCC
    AMULHDU
    AMULHDUCC
    AMULLD
    AMULLDCC
    AMULLDVCC
    AMULLDV
    ARFID
    ARLDMI
    ARLDMICC
    ARLDC
    ARLDCCC
    ARLDCR
    ARLDCRCC
    ARLDCL
    ARLDCLCC
    ASLBIA
    ASLBIE
    ASLBMFEE
    ASLBMFEV
    ASLBMTE
    ASLD
    ASLDCC
    ASRD
    ASRAD
    ASRADCC
    ASRDCC
    ASTDCCC
    ATD

    /* 64-bit pseudo operation */
    ADWORD
    AREMD
    AREMDCC
    AREMDV
    AREMDVCC
    AREMDU
    AREMDUCC
    AREMDUV
    AREMDUVCC

    /* more 64-bit operations */
    AHRFID

    ALAST

    // aliases
    ABR = obj.AJMP
    ABL = obj.ACALL
)

// * GENERAL:
//  *
//  * compiler allocates R3 up as temps
//  * compiler allocates register variables R7-R27
//  * compiler allocates external registers R30 down
//  *
//  * compiler allocates register variables F17-F26
//  * compiler allocates external registers F26 down
const (
    BIG = 32768 - 8
)

const (
    C_NONE = iota
    C_REG
    C_FREG
    C_CREG
    C_SPR /* special processor register */
    C_ZCON
    C_SCON   /* 16 bit signed */
    C_UCON   /* 32 bit signed, low 16 bits 0 */
    C_ADDCON /* -0x8000 <= v < 0 */
    C_ANDCON /* 0 < v <= 0xFFFF */
    C_LCON   /* other 32 */
    C_DCON   /* other 64 (could subdivide further) */
    C_SACON  /* $n(REG) where n <= int16 */
    C_SECON
    C_LACON /* $n(REG) where int16 < n <= int32 */
    C_LECON
    C_DACON /* $n(REG) where int32 < n */
    C_SBRA
    C_LBRA
    C_LBRAPIC
    C_SAUTO
    C_LAUTO
    C_SEXT
    C_LEXT
    C_ZOREG
    C_SOREG
    C_LOREG
    C_FPSCR
    C_MSR
    C_XER
    C_LR
    C_CTR
    C_ANY
    C_GOK
    C_ADDR
    C_GOTADDR
    C_TLS_LE
    C_TLS_IE
    C_TEXTSIZE

    C_NCLASS /* must be the last */
)

const (
    D_FORM = iota
    DS_FORM
)

const (
    FuncAlign = 8
)

const (
    /* mark flags */
    LABEL   = 1 << 0
    LEAF    = 1 << 1
    FLOAT   = 1 << 2
    BRANCH  = 1 << 3
    LOAD    = 1 << 4
    FCMP    = 1 << 5
    SYNC    = 1 << 6
    LIST    = 1 << 7
    FOLL    = 1 << 8
    NOSCHED = 1 << 9
)

// * powerpc 64
const (
    NSNAME = 8
    NSYM   = 50
    NREG   = 32 /* number of general registers */
    NFREG  = 32 /* number of floating point registers */
)

const (
    /* each rhs is OPVCC(_, _, _, _) */
    OP_ADD    = 31<<26 | 266<<1 | 0<<10 | 0
    OP_ADDI   = 14<<26 | 0<<1 | 0<<10 | 0
    OP_ADDIS  = 15<<26 | 0<<1 | 0<<10 | 0
    OP_ANDI   = 28<<26 | 0<<1 | 0<<10 | 0
    OP_EXTSB  = 31<<26 | 954<<1 | 0<<10 | 0
    OP_EXTSH  = 31<<26 | 922<<1 | 0<<10 | 0
    OP_EXTSW  = 31<<26 | 986<<1 | 0<<10 | 0
    OP_MCRF   = 19<<26 | 0<<1 | 0<<10 | 0
    OP_MCRFS  = 63<<26 | 64<<1 | 0<<10 | 0
    OP_MCRXR  = 31<<26 | 512<<1 | 0<<10 | 0
    OP_MFCR   = 31<<26 | 19<<1 | 0<<10 | 0
    OP_MFFS   = 63<<26 | 583<<1 | 0<<10 | 0
    OP_MFMSR  = 31<<26 | 83<<1 | 0<<10 | 0
    OP_MFSPR  = 31<<26 | 339<<1 | 0<<10 | 0
    OP_MFSR   = 31<<26 | 595<<1 | 0<<10 | 0
    OP_MFSRIN = 31<<26 | 659<<1 | 0<<10 | 0
    OP_MTCRF  = 31<<26 | 144<<1 | 0<<10 | 0
    OP_MTFSF  = 63<<26 | 711<<1 | 0<<10 | 0
    OP_MTFSFI = 63<<26 | 134<<1 | 0<<10 | 0
    OP_MTMSR  = 31<<26 | 146<<1 | 0<<10 | 0
    OP_MTMSRD = 31<<26 | 178<<1 | 0<<10 | 0
    OP_MTSPR  = 31<<26 | 467<<1 | 0<<10 | 0
    OP_MTSR   = 31<<26 | 210<<1 | 0<<10 | 0
    OP_MTSRIN = 31<<26 | 242<<1 | 0<<10 | 0
    OP_MULLW  = 31<<26 | 235<<1 | 0<<10 | 0
    OP_MULLD  = 31<<26 | 233<<1 | 0<<10 | 0
    OP_OR     = 31<<26 | 444<<1 | 0<<10 | 0
    OP_ORI    = 24<<26 | 0<<1 | 0<<10 | 0
    OP_ORIS   = 25<<26 | 0<<1 | 0<<10 | 0
    OP_RLWINM = 21<<26 | 0<<1 | 0<<10 | 0
    OP_SUBF   = 31<<26 | 40<<1 | 0<<10 | 0
    OP_RLDIC  = 30<<26 | 4<<1 | 0<<10 | 0
    OP_RLDICR = 30<<26 | 2<<1 | 0<<10 | 0
    OP_RLDICL = 30<<26 | 0<<1 | 0<<10 | 0
)

const (
    REG_R0 = obj.RBasePPC64 + iota
    REG_R1
    REG_R2
    REG_R3
    REG_R4
    REG_R5
    REG_R6
    REG_R7
    REG_R8
    REG_R9
    REG_R10
    REG_R11
    REG_R12
    REG_R13
    REG_R14
    REG_R15
    REG_R16
    REG_R17
    REG_R18
    REG_R19
    REG_R20
    REG_R21
    REG_R22
    REG_R23
    REG_R24
    REG_R25
    REG_R26
    REG_R27
    REG_R28
    REG_R29
    REG_R30
    REG_R31

    REG_F0
    REG_F1
    REG_F2
    REG_F3
    REG_F4
    REG_F5
    REG_F6
    REG_F7
    REG_F8
    REG_F9
    REG_F10
    REG_F11
    REG_F12
    REG_F13
    REG_F14
    REG_F15
    REG_F16
    REG_F17
    REG_F18
    REG_F19
    REG_F20
    REG_F21
    REG_F22
    REG_F23
    REG_F24
    REG_F25
    REG_F26
    REG_F27
    REG_F28
    REG_F29
    REG_F30
    REG_F31

    REG_CR0
    REG_CR1
    REG_CR2
    REG_CR3
    REG_CR4
    REG_CR5
    REG_CR6
    REG_CR7

    REG_MSR
    REG_FPSCR
    REG_CR

    REG_SPECIAL = REG_CR0

    REG_SPR0 = obj.RBasePPC64 + 1024 // first of 1024 registers
    REG_DCR0 = obj.RBasePPC64 + 2048 // first of 1024 registers

    REG_XER = REG_SPR0 + 1
    REG_LR  = REG_SPR0 + 8
    REG_CTR = REG_SPR0 + 9

    REGZERO  = REG_R0 /* set to zero */
    REGSP    = REG_R1
    REGSB    = REG_R2
    REGRET   = REG_R3
    REGARG   = -1      /* -1 disables passing the first argument in register */
    REGRT1   = REG_R3  /* reserved for runtime, duffzero and duffcopy */
    REGRT2   = REG_R4  /* reserved for runtime, duffcopy */
    REGMIN   = REG_R7  /* register variables allocated from here to REGMAX */
    REGCTXT  = REG_R11 /* context for closures */
    REGTLS   = REG_R13 /* C ABI TLS base pointer */
    REGMAX   = REG_R27
    REGEXT   = REG_R30 /* external registers allocated from here down */
    REGG     = REG_R30 /* G */
    REGTMP   = REG_R31 /* used by the linker */
    FREGRET  = REG_F0
    FREGMIN  = REG_F17 /* first register variable */
    FREGMAX  = REG_F26 /* last register variable for 9g only */
    FREGEXT  = REG_F26 /* first external register */
    FREGCVI  = REG_F27 /* floating conversion constant */
    FREGZERO = REG_F28 /* both float and double */
    FREGHALF = REG_F29 /* double */
    FREGONE  = REG_F30 /* double */
    FREGTWO  = REG_F31 /* double */
)

var Anames = []string{
    obj.A_ARCHSPECIFIC: "ADD",
    "ADDCC",
    "ADDV",
    "ADDVCC",
    "ADDC",
    "ADDCCC",
    "ADDCV",
    "ADDCVCC",
    "ADDME",
    "ADDMECC",
    "ADDMEVCC",
    "ADDMEV",
    "ADDE",
    "ADDECC",
    "ADDEVCC",
    "ADDEV",
    "ADDZE",
    "ADDZECC",
    "ADDZEVCC",
    "ADDZEV",
    "AND",
    "ANDCC",
    "ANDN",
    "ANDNCC",
    "BC",
    "BCL",
    "BEQ",
    "BGE",
    "BGT",
    "BLE",
    "BLT",
    "BNE",
    "BVC",
    "BVS",
    "CMP",
    "CMPU",
    "CNTLZW",
    "CNTLZWCC",
    "CRAND",
    "CRANDN",
    "CREQV",
    "CRNAND",
    "CRNOR",
    "CROR",
    "CRORN",
    "CRXOR",
    "DIVW",
    "DIVWCC",
    "DIVWVCC",
    "DIVWV",
    "DIVWU",
    "DIVWUCC",
    "DIVWUVCC",
    "DIVWUV",
    "EQV",
    "EQVCC",
    "EXTSB",
    "EXTSBCC",
    "EXTSH",
    "EXTSHCC",
    "FABS",
    "FABSCC",
    "FADD",
    "FADDCC",
    "FADDS",
    "FADDSCC",
    "FCMPO",
    "FCMPU",
    "FCTIW",
    "FCTIWCC",
    "FCTIWZ",
    "FCTIWZCC",
    "FDIV",
    "FDIVCC",
    "FDIVS",
    "FDIVSCC",
    "FMADD",
    "FMADDCC",
    "FMADDS",
    "FMADDSCC",
    "FMOVD",
    "FMOVDCC",
    "FMOVDU",
    "FMOVS",
    "FMOVSU",
    "FMSUB",
    "FMSUBCC",
    "FMSUBS",
    "FMSUBSCC",
    "FMUL",
    "FMULCC",
    "FMULS",
    "FMULSCC",
    "FNABS",
    "FNABSCC",
    "FNEG",
    "FNEGCC",
    "FNMADD",
    "FNMADDCC",
    "FNMADDS",
    "FNMADDSCC",
    "FNMSUB",
    "FNMSUBCC",
    "FNMSUBS",
    "FNMSUBSCC",
    "FRSP",
    "FRSPCC",
    "FSUB",
    "FSUBCC",
    "FSUBS",
    "FSUBSCC",
    "MOVMW",
    "LSW",
    "LWAR",
    "MOVWBR",
    "MOVB",
    "MOVBU",
    "MOVBZ",
    "MOVBZU",
    "MOVH",
    "MOVHBR",
    "MOVHU",
    "MOVHZ",
    "MOVHZU",
    "MOVW",
    "MOVWU",
    "MOVFL",
    "MOVCRFS",
    "MTFSB0",
    "MTFSB0CC",
    "MTFSB1",
    "MTFSB1CC",
    "MULHW",
    "MULHWCC",
    "MULHWU",
    "MULHWUCC",
    "MULLW",
    "MULLWCC",
    "MULLWVCC",
    "MULLWV",
    "NAND",
    "NANDCC",
    "NEG",
    "NEGCC",
    "NEGVCC",
    "NEGV",
    "NOR",
    "NORCC",
    "OR",
    "ORCC",
    "ORN",
    "ORNCC",
    "REM",
    "REMCC",
    "REMV",
    "REMVCC",
    "REMU",
    "REMUCC",
    "REMUV",
    "REMUVCC",
    "RFI",
    "RLWMI",
    "RLWMICC",
    "RLWNM",
    "RLWNMCC",
    "SLW",
    "SLWCC",
    "SRW",
    "SRAW",
    "SRAWCC",
    "SRWCC",
    "STSW",
    "STWCCC",
    "SUB",
    "SUBCC",
    "SUBVCC",
    "SUBC",
    "SUBCCC",
    "SUBCV",
    "SUBCVCC",
    "SUBME",
    "SUBMECC",
    "SUBMEVCC",
    "SUBMEV",
    "SUBV",
    "SUBE",
    "SUBECC",
    "SUBEV",
    "SUBEVCC",
    "SUBZE",
    "SUBZECC",
    "SUBZEVCC",
    "SUBZEV",
    "SYNC",
    "XOR",
    "XORCC",
    "DCBF",
    "DCBI",
    "DCBST",
    "DCBT",
    "DCBTST",
    "DCBZ",
    "ECIWX",
    "ECOWX",
    "EIEIO",
    "ICBI",
    "ISYNC",
    "PTESYNC",
    "TLBIE",
    "TLBIEL",
    "TLBSYNC",
    "TW",
    "SYSCALL",
    "WORD",
    "RFCI",
    "FRES",
    "FRESCC",
    "FRSQRTE",
    "FRSQRTECC",
    "FSEL",
    "FSELCC",
    "FSQRT",
    "FSQRTCC",
    "FSQRTS",
    "FSQRTSCC",
    "CNTLZD",
    "CNTLZDCC",
    "CMPW",
    "CMPWU",
    "DIVD",
    "DIVDCC",
    "DIVDVCC",
    "DIVDV",
    "DIVDU",
    "DIVDUCC",
    "DIVDUVCC",
    "DIVDUV",
    "EXTSW",
    "EXTSWCC",
    "FCFID",
    "FCFIDCC",
    "FCTID",
    "FCTIDCC",
    "FCTIDZ",
    "FCTIDZCC",
    "LDAR",
    "MOVD",
    "MOVDU",
    "MOVWZ",
    "MOVWZU",
    "MULHD",
    "MULHDCC",
    "MULHDU",
    "MULHDUCC",
    "MULLD",
    "MULLDCC",
    "MULLDVCC",
    "MULLDV",
    "RFID",
    "RLDMI",
    "RLDMICC",
    "RLDC",
    "RLDCCC",
    "RLDCR",
    "RLDCRCC",
    "RLDCL",
    "RLDCLCC",
    "SLBIA",
    "SLBIE",
    "SLBMFEE",
    "SLBMFEV",
    "SLBMTE",
    "SLD",
    "SLDCC",
    "SRD",
    "SRAD",
    "SRADCC",
    "SRDCC",
    "STDCCC",
    "TD",
    "DWORD",
    "REMD",
    "REMDCC",
    "REMDV",
    "REMDVCC",
    "REMDU",
    "REMDUCC",
    "REMDUV",
    "REMDUVCC",
    "HRFID",
    "LAST",
}

var Linkppc64 = obj.LinkArch{
    ByteOrder:  binary.BigEndian,
    Name:       "ppc64",
    Thechar:    '9',
    Preprocess: preprocess,
    Assemble:   span9,
    Follow:     follow,
    Progedit:   progedit,
    Minlc:      4,
    Ptrsize:    8,
    Regsize:    8,
}

var Linkppc64le = obj.LinkArch{
    ByteOrder:  binary.LittleEndian,
    Name:       "ppc64le",
    Thechar:    '9',
    Preprocess: preprocess,
    Assemble:   span9,
    Follow:     follow,
    Progedit:   progedit,
    Minlc:      4,
    Ptrsize:    8,
    Regsize:    8,
}

type Oprang struct {
    start []Optab
    stop  []Optab
}

type Optab struct {
    as    int16
    a1    uint8
    a2    uint8
    a3    uint8
    a4    uint8
    type_ int8
    size  int8
    param int16
}

func AOP_IRR(op uint32, d uint32, a uint32, simm uint32) uint32

// the order is dest, a/s, b/imm for both arithmetic and logical operations
func AOP_RRR(op uint32, d uint32, a uint32, b uint32) uint32

func DRconv(a int) string

func LOP_IRR(op uint32, a uint32, s uint32, uimm uint32) uint32

func LOP_RRR(op uint32, a uint32, s uint32, b uint32) uint32

func OP(o uint32, xo uint32) uint32

func OPCC(o uint32, xo uint32, rc uint32) uint32

func OPVCC(o uint32, xo uint32, oe uint32, rc uint32) uint32

func OP_BC(op uint32, bo uint32, bi uint32, bd uint32, aa uint32) uint32

func OP_BCR(op uint32, bo uint32, bi uint32) uint32

func OP_BR(op uint32, li uint32, aa uint32) uint32

func OP_RLW(op uint32, a uint32, s uint32, sh uint32, mb uint32, me uint32) uint32

func Rconv(r int) string

