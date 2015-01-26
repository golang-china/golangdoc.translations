// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Addr2line is a minimal simulation of the GNU addr2line tool, just enough to
// support pprof.
//
// Usage:
//
//	go tool addr2line binary
//
// Addr2line reads hexadecimal addresses, one per line and with optional 0x prefix,
// from standard input. For each input address, addr2line prints two output lines,
// first the name of the function containing the address and second the file:line
// of the source code corresponding to that address.
//
// This tool is intended for use only by pprof; its interface may change or it may
// be deleted entirely in future releases.

// Addr2line is a minimal simulation of the
// GNU addr2line tool, just enough to
// support pprof.
//
// Usage:
//
//	go tool addr2line binary
//
// Addr2line reads hexadecimal addresses,
// one per line and with optional 0x
// prefix, from standard input. For each
// input address, addr2line prints two
// output lines, first the name of the
// function containing the address and
// second the file:line of the source code
// corresponding to that address.
//
// This tool is intended for use only by
// pprof; its interface may change or it
// may be deleted entirely in future
// releases.
package main
