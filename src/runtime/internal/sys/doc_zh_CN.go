// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// package sys contains system- and configuration- and architecture-specific
// constants used by the runtime.
package sys

const (
	AMD64 ArchFamilyType = iota
	ARM
	ARM64
	I386
	MIPS64
	PPC64
	S390X
)

const (
	ArchFamily    = AMD64
	BigEndian     = 0
	CacheLineSize = 64
	PhysPageSize  = 4096
	PCQuantum     = 1
	Int64Align    = 8
	HugePageSize  = 1 << 21
	MinFrameSize  = 0
)

const DefaultGoroot = `/usr/local/Cellar/go/1.7/libexec`

const GOARCH = `amd64`

const Goarch386 = 0

const GoarchAmd64 = 1

const GoarchAmd64p32 = 0

const GoarchArm = 0

const GoarchArm64 = 0

const GoarchArm64be = 0

const GoarchArmbe = 0

const GoarchMips = 0

const GoarchMips64 = 0

const GoarchMips64le = 0

const GoarchMips64p32 = 0

const GoarchMips64p32le = 0

const GoarchMipsle = 0

const GoarchPpc = 0

const GoarchPpc64 = 0

const GoarchPpc64le = 0

const GoarchS390 = 0

const GoarchS390x = 0

const GoarchSparc = 0

const GoarchSparc64 = 0

const Goexperiment = ``

const PtrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

const RegSize = 4 << (^Uintreg(0) >> 63) // unsafe.Sizeof(uintreg(0)) but an ideal const

const SpAlign = 1*(1-GoarchArm64) + 16*GoarchArm64

const StackGuardMultiplier = 1

const TheVersion = `go1.7`

type ArchFamilyType int

type Uintreg uint64

