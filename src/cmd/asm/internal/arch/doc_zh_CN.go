// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package arch

import (
    "cmd/internal/obj"
    "cmd/internal/obj/arm"
    "cmd/internal/obj/arm64"
    "cmd/internal/obj/mips"
    "cmd/internal/obj/ppc64"
    "cmd/internal/obj/s390x"
    "cmd/internal/obj/x86"
    "fmt"
    "strings"
)

// Pseudo-registers whose names are the constant name without the leading R.
const (
	RFP = -(iota + 1)
	RSB
	RSP
	RPC
)


// Arch wraps the link architecture object with more architecture-specific
// information.

// Arch wraps the link architecture object with more architecture-specific
// information.
type Arch struct {

	// Map of instruction names to enumeration.
	Instructions map[string]obj.As
	// Map of register names to enumeration.
	Register map[string]int16
	// Table of register prefix names. These are things like R for R(0) and SPR for SPR(268).
	RegisterPrefix map[string]bool
	// RegisterNumber converts R(10) into arm.REG_R10.
	RegisterNumber func(string, int16) (int16, bool)
	// Instruction is a jump.
	IsJump func(word string) bool
}


// ARM64Suffix handles the special suffix for the ARM64.
// It returns a boolean to indicate success; failure means
// cond was unrecognized.
func ARM64Suffix(prog *obj.Prog, cond string) bool

// ARMConditionCodes handles the special condition code situation for the ARM.
// It returns a boolean to indicate success; failure means cond was
// unrecognized.
func ARMConditionCodes(prog *obj.Prog, cond string) bool

// ARMMRCOffset implements the peculiar encoding of the MRC and MCR
// instructions. The difference between MRC and MCR is represented by a bit high
// in the word, not in the usual way by the opcode itself. Asm must use AMRC for
// both instructions, so we return the opcode for MRC so that asm doesn't need
// to import obj/arm.
func ARMMRCOffset(op obj.As, cond string, x0, x1, x2, x3, x4, x5 int64) (offset int64, op0 obj.As, ok bool)

// IsAMD4OP reports whether the op (as defined by an ppc64.A* constant) is
// The FMADD-like instructions behave similarly.
func IsAMD4OP(op obj.As) bool

// IsARM64CMP reports whether the op (as defined by an arm.A* constant) is
// one of the comparison instructions that require special handling.
func IsARM64CMP(op obj.As) bool

// IsARM64STLXR reports whether the op (as defined by an arm64.A*
// constant) is one of the STLXR-like instructions that require special
// handling.
func IsARM64STLXR(op obj.As) bool

// IsARMCMP reports whether the op (as defined by an arm.A* constant) is
// one of the comparison instructions that require special handling.
func IsARMCMP(op obj.As) bool

// IsARMFloatCmp reports whether the op is a floating comparison instruction.
func IsARMFloatCmp(op obj.As) bool

// IsARMMRC reports whether the op (as defined by an arm.A* constant) is
// MRC or MCR
func IsARMMRC(op obj.As) bool

// IsARMMULA reports whether the op (as defined by an arm.A* constant) is
// MULA, MULAWT or MULAWB, the 4-operand instructions.
func IsARMMULA(op obj.As) bool

// IsARMSTREX reports whether the op (as defined by an arm.A* constant) is
// one of the STREX-like instructions that require special handling.
func IsARMSTREX(op obj.As) bool

// IsMIPS64CMP reports whether the op (as defined by an mips.A* constant) is
// one of the CMP instructions that require special handling.
func IsMIPS64CMP(op obj.As) bool

// IsMIPS64MUL reports whether the op (as defined by an mips.A* constant) is
// one of the MUL/DIV/REM instructions that require special handling.
func IsMIPS64MUL(op obj.As) bool

// IsPPC64CMP reports whether the op (as defined by an ppc64.A* constant) is
// one of the CMP instructions that require special handling.
func IsPPC64CMP(op obj.As) bool

// IsPPC64NEG reports whether the op (as defined by an ppc64.A* constant) is
// one of the NEG-like instructions that require special handling.
func IsPPC64NEG(op obj.As) bool

// IsPPC64RLD reports whether the op (as defined by an ppc64.A* constant) is
// one of the RLD-like instructions that require special handling.
// The FMADD-like instructions behave similarly.
func IsPPC64RLD(op obj.As) bool

// IsS390xCMP reports whether the op (as defined by an s390x.A* constant) is
// one of the CMP instructions that require special handling.
func IsS390xCMP(op obj.As) bool

// IsS390xNEG reports whether the op (as defined by an s390x.A* constant) is
// one of the NEG-like instructions that require special handling.
func IsS390xNEG(op obj.As) bool

// IsS390xRLD reports whether the op (as defined by an s390x.A* constant) is
// one of the RLD-like instructions that require special handling.
// The FMADD-like instructions behave similarly.
func IsS390xRLD(op obj.As) bool

// IsS390xWithIndex reports whether the op (as defined by an s390x.A* constant)
// refers to an instruction which takes an index as its first argument.
func IsS390xWithIndex(op obj.As) bool

// IsS390xWithLength reports whether the op (as defined by an s390x.A* constant)
// refers to an instruction which takes a length as its first argument.
func IsS390xWithLength(op obj.As) bool

// ParseARM64Suffix parses the suffix attached to an ARM64 instruction.
// The input is a single string consisting of period-separated condition
// codes, such as ".P.W". An initial period is ignored.
func ParseARM64Suffix(cond string) (uint8, bool)

// ParseARMCondition parses the conditions attached to an ARM instruction.
// The input is a single string consisting of period-separated condition
// codes, such as ".P.W". An initial period is ignored.
func ParseARMCondition(cond string) (uint8, bool)

// Set configures the architecture specified by GOARCH and returns its
// representation. It returns nil if GOARCH is not recognized.
func Set(GOARCH string) *Arch

