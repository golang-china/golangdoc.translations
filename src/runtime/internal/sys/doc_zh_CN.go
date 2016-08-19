// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// package sys contains system- and configuration- and architecture-specific
// constants used by the runtime.

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



const GOARCH = `amd64`



const GOOS = `linux`



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



const GoosAndroid = 0



const GoosDarwin = 0



const GoosDragonfly = 0



const GoosFreebsd = 0



const GoosLinux = 1



const GoosNacl = 0



const GoosNetbsd = 0



const GoosOpenbsd = 0



const GoosPlan9 = 0



const GoosSolaris = 0



const GoosWindows = 0



const PtrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const



const RegSize = 4 << (^Uintreg(0) >> 63) // unsafe.Sizeof(uintreg(0)) but an ideal const



const SpAlign = 1*(1-GoarchArm64) + 16*GoarchArm64



type ArchFamilyType int



type Uintreg uint64


// Bswap32 returns its input with byte order reversed
// 0x01020304 -> 0x04030201
func Bswap32(x uint32) uint32

// Bswap64 returns its input with byte order reversed
// 0x0102030405060708 -> 0x0807060504030201
func Bswap64(x uint64) uint64

// Ctz16 counts trailing (low-order) zeroes,
// and if all are zero, then 16.
func Ctz16(x uint16) uint16

// Ctz32 counts trailing (low-order) zeroes,
// and if all are zero, then 32.
func Ctz32(x uint32) uint32

// Ctz64 counts trailing (low-order) zeroes,
// and if all are zero, then 64.
func Ctz64(x uint64) uint64

// Ctz8 counts trailing (low-order) zeroes,
// and if all are zero, then 8.
func Ctz8(x uint8) uint8

