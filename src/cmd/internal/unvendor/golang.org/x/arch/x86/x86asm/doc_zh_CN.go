// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package x86asm implements decoding of x86 machine code.
package x86asm // import "cmd/internal/unvendor/golang.org/x/arch/x86/x86asm"

import (
    "bufio"
    "bytes"
    "debug/elf"
    "encoding/binary"
    "encoding/hex"
    "errors"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "math/rand"
    "os"
    "os/exec"
    "regexp"
    "runtime"
    "strconv"
    "strings"
    "testing"
    "time"
)

const (
    // Metadata about the role of a prefix in an instruction.
    PrefixImplicit Prefix = 0x8000 // prefix is implied by instruction text
    PrefixIgnored  Prefix = 0x4000 // prefix is ignored: either irrelevant or overridden by a later prefix
    PrefixInvalid  Prefix = 0x2000 // prefix makes entire instruction invalid (bad LOCK)

    // Memory segment overrides.
    PrefixES Prefix = 0x26 // ES segment override
    PrefixCS Prefix = 0x2E // CS segment override
    PrefixSS Prefix = 0x36 // SS segment override
    PrefixDS Prefix = 0x3E // DS segment override
    PrefixFS Prefix = 0x64 // FS segment override
    PrefixGS Prefix = 0x65 // GS segment override

    // Branch prediction.
    PrefixPN Prefix = 0x12E // predict not taken (conditional branch only)
    PrefixPT Prefix = 0x13E // predict taken (conditional branch only)

    // Size attributes.
    PrefixDataSize Prefix = 0x66 // operand size override
    PrefixData16   Prefix = 0x166
    PrefixData32   Prefix = 0x266
    PrefixAddrSize Prefix = 0x67 // address size override
    PrefixAddr16   Prefix = 0x167
    PrefixAddr32   Prefix = 0x267

    // One of a kind.
    PrefixLOCK     Prefix = 0xF0 // lock
    PrefixREPN     Prefix = 0xF2 // repeat not zero
    PrefixXACQUIRE Prefix = 0x1F2
    PrefixBND      Prefix = 0x2F2
    PrefixREP      Prefix = 0xF3 // repeat
    PrefixXRELEASE Prefix = 0x1F3

    // The REX prefixes must be in the range [PrefixREX, PrefixREX+0x10).
    // the other bits are set or not according to the intended use.
    PrefixREX  Prefix = 0x40 // REX 64-bit extension prefix
    PrefixREXW Prefix = 0x08 // extension bit W (64-bit instruction width)
    PrefixREXR Prefix = 0x04 // extension bit R (r field in modrm)
    PrefixREXX Prefix = 0x02 // extension bit X (index field in sib)
    PrefixREXB Prefix = 0x01 // extension bit B (r/m field in modrm or base field in sib)
)

const (
    AAA
    AAD
    AAM
    AAS
    ADC
    ADD
    ADDPD
    ADDPS
    ADDSD
    ADDSS
    ADDSUBPD
    ADDSUBPS
    AESDEC
    AESDECLAST
    AESENC
    AESENCLAST
    AESIMC
    AESKEYGENASSIST
    AND
    ANDNPD
    ANDNPS
    ANDPD
    ANDPS
    ARPL
    BLENDPD
    BLENDPS
    BLENDVPD
    BLENDVPS
    BOUND
    BSF
    BSR
    BSWAP
    BT
    BTC
    BTR
    BTS
    CALL
    CBW
    CDQ
    CDQE
    CLC
    CLD
    CLFLUSH
    CLI
    CLTS
    CMC
    CMOVA
    CMOVAE
    CMOVB
    CMOVBE
    CMOVE
    CMOVG
    CMOVGE
    CMOVL
    CMOVLE
    CMOVNE
    CMOVNO
    CMOVNP
    CMOVNS
    CMOVO
    CMOVP
    CMOVS
    CMP
    CMPPD
    CMPPS
    CMPSB
    CMPSD
    CMPSD_XMM
    CMPSQ
    CMPSS
    CMPSW
    CMPXCHG
    CMPXCHG16B
    CMPXCHG8B
    COMISD
    COMISS
    CPUID
    CQO
    CRC32
    CVTDQ2PD
    CVTDQ2PS
    CVTPD2DQ
    CVTPD2PI
    CVTPD2PS
    CVTPI2PD
    CVTPI2PS
    CVTPS2DQ
    CVTPS2PD
    CVTPS2PI
    CVTSD2SI
    CVTSD2SS
    CVTSI2SD
    CVTSI2SS
    CVTSS2SD
    CVTSS2SI
    CVTTPD2DQ
    CVTTPD2PI
    CVTTPS2DQ
    CVTTPS2PI
    CVTTSD2SI
    CVTTSS2SI
    CWD
    CWDE
    DAA
    DAS
    DEC
    DIV
    DIVPD
    DIVPS
    DIVSD
    DIVSS
    DPPD
    DPPS
    EMMS
    ENTER
    EXTRACTPS
    F2XM1
    FABS
    FADD
    FADDP
    FBLD
    FBSTP
    FCHS
    FCMOVB
    FCMOVBE
    FCMOVE
    FCMOVNB
    FCMOVNBE
    FCMOVNE
    FCMOVNU
    FCMOVU
    FCOM
    FCOMI
    FCOMIP
    FCOMP
    FCOMPP
    FCOS
    FDECSTP
    FDIV
    FDIVP
    FDIVR
    FDIVRP
    FFREE
    FFREEP
    FIADD
    FICOM
    FICOMP
    FIDIV
    FIDIVR
    FILD
    FIMUL
    FINCSTP
    FIST
    FISTP
    FISTTP
    FISUB
    FISUBR
    FLD
    FLD1
    FLDCW
    FLDENV
    FLDL2E
    FLDL2T
    FLDLG2
    FLDPI
    FMUL
    FMULP
    FNCLEX
    FNINIT
    FNOP
    FNSAVE
    FNSTCW
    FNSTENV
    FNSTSW
    FPATAN
    FPREM
    FPREM1
    FPTAN
    FRNDINT
    FRSTOR
    FSCALE
    FSIN
    FSINCOS
    FSQRT
    FST
    FSTP
    FSUB
    FSUBP
    FSUBR
    FSUBRP
    FTST
    FUCOM
    FUCOMI
    FUCOMIP
    FUCOMP
    FUCOMPP
    FWAIT
    FXAM
    FXCH
    FXRSTOR
    FXRSTOR64
    FXSAVE
    FXSAVE64
    FXTRACT
    FYL2X
    FYL2XP1
    HADDPD
    HADDPS
    HLT
    HSUBPD
    HSUBPS
    ICEBP
    IDIV
    IMUL
    IN
    INC
    INSB
    INSD
    INSERTPS
    INSW
    INT
    INTO
    INVD
    INVLPG
    INVPCID
    IRET
    IRETD
    IRETQ
    JA
    JAE
    JB
    JBE
    JCXZ
    JE
    JECXZ
    JG
    JGE
    JL
    JLE
    JMP
    JNE
    JNO
    JNP
    JNS
    JO
    JP
    JRCXZ
    JS
    LAHF
    LAR
    LCALL
    LDDQU
    LDMXCSR
    LDS
    LEA
    LEAVE
    LES
    LFENCE
    LFS
    LGDT
    LGS
    LIDT
    LJMP
    LLDT
    LMSW
    LODSB
    LODSD
    LODSQ
    LODSW
    LOOP
    LOOPE
    LOOPNE
    LRET
    LSL
    LSS
    LTR
    LZCNT
    MASKMOVDQU
    MASKMOVQ
    MAXPD
    MAXPS
    MAXSD
    MAXSS
    MFENCE
    MINPD
    MINPS
    MINSD
    MINSS
    MONITOR
    MOV
    MOVAPD
    MOVAPS
    MOVBE
    MOVD
    MOVDDUP
    MOVDQ2Q
    MOVDQA
    MOVDQU
    MOVHLPS
    MOVHPD
    MOVHPS
    MOVLHPS
    MOVLPD
    MOVLPS
    MOVMSKPD
    MOVMSKPS
    MOVNTDQ
    MOVNTDQA
    MOVNTI
    MOVNTPD
    MOVNTPS
    MOVNTQ
    MOVNTSD
    MOVNTSS
    MOVQ
    MOVQ2DQ
    MOVSB
    MOVSD
    MOVSD_XMM
    MOVSHDUP
    MOVSLDUP
    MOVSQ
    MOVSS
    MOVSW
    MOVSX
    MOVSXD
    MOVUPD
    MOVUPS
    MOVZX
    MPSADBW
    MUL
    MULPD
    MULPS
    MULSD
    MULSS
    MWAIT
    NEG
    NOP
    NOT
    OR
    ORPD
    ORPS
    OUT
    OUTSB
    OUTSD
    OUTSW
    PABSB
    PABSD
    PABSW
    PACKSSDW
    PACKSSWB
    PACKUSDW
    PACKUSWB
    PADDB
    PADDD
    PADDQ
    PADDSB
    PADDSW
    PADDUSB
    PADDUSW
    PADDW
    PALIGNR
    PAND
    PANDN
    PAUSE
    PAVGB
    PAVGW
    PBLENDVB
    PBLENDW
    PCLMULQDQ
    PCMPEQB
    PCMPEQD
    PCMPEQQ
    PCMPEQW
    PCMPESTRI
    PCMPESTRM
    PCMPGTB
    PCMPGTD
    PCMPGTQ
    PCMPGTW
    PCMPISTRI
    PCMPISTRM
    PEXTRB
    PEXTRD
    PEXTRQ
    PEXTRW
    PHADDD
    PHADDSW
    PHADDW
    PHMINPOSUW
    PHSUBD
    PHSUBSW
    PHSUBW
    PINSRB
    PINSRD
    PINSRQ
    PINSRW
    PMADDUBSW
    PMADDWD
    PMAXSB
    PMAXSD
    PMAXSW
    PMAXUB
    PMAXUD
    PMAXUW
    PMINSB
    PMINSD
    PMINSW
    PMINUB
    PMINUD
    PMINUW
    PMOVMSKB
    PMOVSXBD
    PMOVSXBQ
    PMOVSXBW
    PMOVSXDQ
    PMOVSXWD
    PMOVSXWQ
    PMOVZXBD
    PMOVZXBQ
    PMOVZXBW
    PMOVZXDQ
    PMOVZXWD
    PMOVZXWQ
    PMULDQ
    PMULHRSW
    PMULHUW
    PMULHW
    PMULLD
    PMULLW
    PMULUDQ
    POP
    POPA
    POPAD
    POPCNT
    POPF
    POPFD
    POPFQ
    POR
    PREFETCHNTA
    PREFETCHT0
    PREFETCHT1
    PREFETCHT2
    PREFETCHW
    PSADBW
    PSHUFB
    PSHUFD
    PSHUFHW
    PSHUFLW
    PSHUFW
    PSIGNB
    PSIGND
    PSIGNW
    PSLLD
    PSLLDQ
    PSLLQ
    PSLLW
    PSRAD
    PSRAW
    PSRLD
    PSRLDQ
    PSRLQ
    PSRLW
    PSUBB
    PSUBD
    PSUBQ
    PSUBSB
    PSUBSW
    PSUBUSB
    PSUBUSW
    PSUBW
    PTEST
    PUNPCKHBW
    PUNPCKHDQ
    PUNPCKHQDQ
    PUNPCKHWD
    PUNPCKLBW
    PUNPCKLDQ
    PUNPCKLQDQ
    PUNPCKLWD
    PUSH
    PUSHA
    PUSHAD
    PUSHF
    PUSHFD
    PUSHFQ
    PXOR
    RCL
    RCPPS
    RCPSS
    RCR
    RDFSBASE
    RDGSBASE
    RDMSR
    RDPMC
    RDRAND
    RDTSC
    RDTSCP
    RET
    ROL
    ROR
    ROUNDPD
    ROUNDPS
    ROUNDSD
    ROUNDSS
    RSM
    RSQRTPS
    RSQRTSS
    SAHF
    SAR
    SBB
    SCASB
    SCASD
    SCASQ
    SCASW
    SETA
    SETAE
    SETB
    SETBE
    SETE
    SETG
    SETGE
    SETL
    SETLE
    SETNE
    SETNO
    SETNP
    SETNS
    SETO
    SETP
    SETS
    SFENCE
    SGDT
    SHL
    SHLD
    SHR
    SHRD
    SHUFPD
    SHUFPS
    SIDT
    SLDT
    SMSW
    SQRTPD
    SQRTPS
    SQRTSD
    SQRTSS
    STC
    STD
    STI
    STMXCSR
    STOSB
    STOSD
    STOSQ
    STOSW
    STR
    SUB
    SUBPD
    SUBPS
    SUBSD
    SUBSS
    SWAPGS
    SYSCALL
    SYSENTER
    SYSEXIT
    SYSRET
    TEST
    TZCNT
    UCOMISD
    UCOMISS
    UD1
    UD2
    UNPCKHPD
    UNPCKHPS
    UNPCKLPD
    UNPCKLPS
    VERR
    VERW
    WBINVD
    WRFSBASE
    WRGSBASE
    WRMSR
    XABORT
    XADD
    XBEGIN
    XCHG
    XEND
    XGETBV
    XLATB
    XOR
    XORPD
    XORPS
    XRSTOR
    XRSTOR64
    XRSTORS
    XRSTORS64
    XSAVE
    XSAVE64
    XSAVEC
    XSAVEC64
    XSAVEOPT
    XSAVEOPT64
    XSAVES
    XSAVES64
    XSETBV
    XTEST
)

const (

    // 8-bit
    AL
    CL
    DL
    BL
    AH
    CH
    DH
    BH
    SPB
    BPB
    SIB
    DIB
    R8B
    R9B
    R10B
    R11B
    R12B
    R13B
    R14B
    R15B

    // 16-bit
    AX
    CX
    DX
    BX
    SP
    BP
    SI
    DI
    R8W
    R9W
    R10W
    R11W
    R12W
    R13W
    R14W
    R15W

    // 32-bit
    EAX
    ECX
    EDX
    EBX
    ESP
    EBP
    ESI
    EDI
    R8L
    R9L
    R10L
    R11L
    R12L
    R13L
    R14L
    R15L

    // 64-bit
    RAX
    RCX
    RDX
    RBX
    RSP
    RBP
    RSI
    RDI
    R8
    R9
    R10
    R11
    R12
    R13
    R14
    R15

    // Instruction pointer.
    IP  // 16-bit
    EIP // 32-bit
    RIP // 64-bit

    // 387 floating point registers.
    F0
    F1
    F2
    F3
    F4
    F5
    F6
    F7

    // MMX registers.
    M0
    M1
    M2
    M3
    M4
    M5
    M6
    M7

    // XMM registers.
    X0
    X1
    X2
    X3
    X4
    X5
    X6
    X7
    X8
    X9
    X10
    X11
    X12
    X13
    X14
    X15

    // Segment registers.
    ES
    CS
    SS
    DS
    FS
    GS

    // System registers.
    GDTR
    IDTR
    LDTR
    MSW
    TASK

    // Control registers.
    CR0
    CR1
    CR2
    CR3
    CR4
    CR5
    CR6
    CR7
    CR8
    CR9
    CR10
    CR11
    CR12
    CR13
    CR14
    CR15

    // Debug registers.
    DR0
    DR1
    DR2
    DR3
    DR4
    DR5
    DR6
    DR7
    DR8
    DR9
    DR10
    DR11
    DR12
    DR13
    DR14
    DR15

    // Task registers.
    TR0
    TR1
    TR2
    TR3
    TR4
    TR5
    TR6
    TR7
)

// These are the errors returned by Decode.
var (
    ErrInvalidMode  = errors.New("invalid x86 mode in Decode")
    ErrTruncated    = errors.New("truncated instruction")
    ErrUnrecognized = errors.New("unrecognized instruction")
)

// An Arg is a single instruction argument,
// one of these types: Reg, Mem, Imm, Rel.
type Arg interface {
    String() string
    isArg()
}

// An Args holds the instruction arguments.
// If an instruction has fewer than 4 arguments,
// the final elements in the array are nil.
type Args [4]Arg

// An ExtDis is a connection between an external disassembler and a test.
type ExtDis struct {
    Arch     int
    Dec      chan ExtInst
    File     *os.File
    Size     int
    KeepFile bool
    Cmd      *exec.Cmd
}

// A ExtInst represents a single decoded instruction parsed
// from an external disassembler's output.
type ExtInst struct {
    addr uint32
    enc  [32]byte
    nenc int
    text string
}

// An Imm is an integer constant.
type Imm int64

// An Inst is a single instruction.
type Inst struct {
    Prefix   Prefixes // Prefixes applied to the instruction.
    Op       Op       // Opcode mnemonic
    Opcode   uint32   // Encoded opcode bits, left aligned (first byte is Opcode>>24, etc)
    Args     Args     // Instruction arguments, in Intel order
    Mode     int      // processor mode in bits: 16, 32, or 64
    AddrSize int      // address size in bits: 16, 32, or 64
    DataSize int      // operand size in bits: 16, 32, or 64
    MemBytes int      // size of memory argument in bytes: 1, 2, 4, 8, 16, and so on.
    Len      int      // length of encoded instruction in bytes
    PCRel    int      // length of PC-relative address in instruction encoding
    PCRelOff int      // index of start of PC-relative address in instruction encoding
}

// A Mem is a memory reference.
// The general form is Segment:[Base+Scale*Index+Disp].
type Mem struct {
    Segment Reg
    Base    Reg
    Scale   uint8
    Index   Reg
    Disp    int64
}

// An Op is an x86 opcode.
type Op uint32

// A Prefix represents an Intel instruction prefix.
// The low 8 bits are the actual prefix byte encoding,
// and the top 8 bits contain distinguishing bits and metadata.
type Prefix uint16

// Prefixes is an array of prefixes associated with a single instruction.
// The prefixes are listed in the same order as found in the instruction:
// each prefix byte corresponds to one slot in the array. The first zero
// in the array marks the end of the prefixes.
type Prefixes [14]Prefix

// A Reg is a single register.
// The zero Reg value has no name but indicates ``no register.''
type Reg uint8

// A Rel is an offset relative to the current instruction pointer.
type Rel int32

// Decode decodes the leading bytes in src as a single instruction.
// The mode arguments specifies the assumed processor mode:
// 16, 32, or 64 for 16-, 32-, and 64-bit execution modes.
func Decode(src []byte, mode int) (inst Inst, err error)

// GNUSyntax returns the GNU assembler syntax for the instruction, as defined by
// GNU binutils. This general form is often called ``AT&T syntax'' as a reference
// to AT&T System V Unix.
func GNUSyntax(inst Inst) string

// GoSyntax returns the Go assembler syntax for the instruction. The syntax was
// originally defined by Plan 9. The pc is the program counter of the instruction,
// used for expanding PC-relative addresses into absolute ones. The symname function
// queries the symbol table for the program being disassembled. Given a target address
// it returns the name and base address of the symbol containing the target, if
// any; otherwise it returns "", 0.
func GoSyntax(inst Inst, pc uint64, symname func(uint64) (string, uint64)) string

// IntelSyntax returns the Intel assembler syntax for the instruction, as defined
// by Intel's XED tool.
func IntelSyntax(inst Inst) string

func TestDecode(t *testing.T)

func TestObjdump320F(t *testing.T)

func TestObjdump320F38(t *testing.T)

func TestObjdump320F3A(t *testing.T)

func TestObjdump32Manual(t *testing.T)

func TestObjdump32ModRM(t *testing.T)

func TestObjdump32OneByte(t *testing.T)

func TestObjdump32Prefix(t *testing.T)

func TestObjdump32Testdata(t *testing.T)

func TestObjdump640F(t *testing.T)

func TestObjdump640F38(t *testing.T)

func TestObjdump640F3A(t *testing.T)

func TestObjdump64Manual(t *testing.T)

func TestObjdump64ModRM(t *testing.T)

func TestObjdump64OneByte(t *testing.T)

func TestObjdump64Prefix(t *testing.T)

func TestObjdump64REX0F(t *testing.T)

func TestObjdump64REX0F38(t *testing.T)

func TestObjdump64REX0F3A(t *testing.T)

func TestObjdump64REXModRM(t *testing.T)

func TestObjdump64REXOneByte(t *testing.T)

func TestObjdump64REXPrefix(t *testing.T)

func TestObjdump64REXTestdata(t *testing.T)

func TestObjdump64Testdata(t *testing.T)

func TestPlan9320F(t *testing.T)

func TestPlan9320F38(t *testing.T)

func TestPlan9320F3A(t *testing.T)

func TestPlan932Manual(t *testing.T)

func TestPlan932ModRM(t *testing.T)

func TestPlan932OneByte(t *testing.T)

func TestPlan932Prefix(t *testing.T)

func TestPlan932Testdata(t *testing.T)

func TestPlan9640F(t *testing.T)

func TestPlan9640F38(t *testing.T)

func TestPlan9640F3A(t *testing.T)

func TestPlan964Manual(t *testing.T)

func TestPlan964ModRM(t *testing.T)

func TestPlan964OneByte(t *testing.T)

func TestPlan964Prefix(t *testing.T)

func TestPlan964REX0F(t *testing.T)

func TestPlan964REX0F38(t *testing.T)

func TestPlan964REX0F3A(t *testing.T)

func TestPlan964REXModRM(t *testing.T)

func TestPlan964REXOneByte(t *testing.T)

func TestPlan964REXPrefix(t *testing.T)

func TestPlan964REXTestdata(t *testing.T)

func TestPlan964Testdata(t *testing.T)

func TestRegString(t *testing.T)

func TestXed320F(t *testing.T)

func TestXed320F38(t *testing.T)

func TestXed320F3A(t *testing.T)

func TestXed32Manual(t *testing.T)

func TestXed32ModRM(t *testing.T)

func TestXed32OneByte(t *testing.T)

func TestXed32Prefix(t *testing.T)

func TestXed32Testdata(t *testing.T)

func TestXed640F(t *testing.T)

func TestXed640F38(t *testing.T)

func TestXed640F3A(t *testing.T)

func TestXed64Manual(t *testing.T)

func TestXed64ModRM(t *testing.T)

func TestXed64OneByte(t *testing.T)

func TestXed64Prefix(t *testing.T)

func TestXed64REX0F(t *testing.T)

func TestXed64REX0F38(t *testing.T)

func TestXed64REX0F3A(t *testing.T)

func TestXed64REXModRM(t *testing.T)

func TestXed64REXOneByte(t *testing.T)

func TestXed64REXPrefix(t *testing.T)

func TestXed64REXTestdata(t *testing.T)

func TestXed64Testdata(t *testing.T)

// Run runs the given command - the external disassembler - and returns
// a buffered reader of its standard output.
func (*ExtDis) Run(cmd ...string) (*bufio.Reader, error)

// Wait waits for the command started with Run to exit.
func (*ExtDis) Wait() error

func (ExtInst) String() string

func (Imm) String() string

func (Inst) String() string

func (Mem) String() string

func (Op) String() string

// IsREX reports whether p is a REX prefix byte.
func (Prefix) IsREX() bool

func (Prefix) String() string

func (Reg) String() string

func (Rel) String() string

