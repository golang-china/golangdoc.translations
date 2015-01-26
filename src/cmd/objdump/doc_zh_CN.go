// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Objdump disassembles executable files.
//
// Usage:
//
//	go tool objdump [-s symregexp] binary
//
// Objdump prints a disassembly of all text symbols (code) in the binary. If the -s
// option is present, objdump only disassembles symbols with names matching the
// regular expression.
//
// Alternate usage:
//
//	go tool objdump binary start end
//
// In this mode, objdump disassembles the binary starting at the start address and
// stopping at the end address. The start and end addresses are program counters
// written in hexadecimal with optional leading 0x prefix. In this mode, objdump
// prints a sequence of stanzas of the form:
//
//	file:line
//	 address: assembly
//	 address: assembly
//	 ...
//
// Each stanza gives the disassembly for a contiguous range of addresses all mapped
// to the same original source file and line number. This mode is intended for use
// by pprof.

// Objdump disassembles executable files.
//
// Usage:
//
//	go tool objdump [-s symregexp] binary
//
// Objdump prints a disassembly of all text
// symbols (code) in the binary. If the -s
// option is present, objdump only
// disassembles symbols with names matching
// the regular expression.
//
// Alternate usage:
//
//	go tool objdump binary start end
//
// In this mode, objdump disassembles the
// binary starting at the start address and
// stopping at the end address. The start
// and end addresses are program counters
// written in hexadecimal with optional
// leading 0x prefix. In this mode, objdump
// prints a sequence of stanzas of the
// form:
//
//	file:line
//	 address: assembly
//	 address: assembly
//	 ...
//
// Each stanza gives the disassembly for a
// contiguous range of addresses all mapped
// to the same original source file and
// line number. This mode is intended for
// use by pprof.
package main
