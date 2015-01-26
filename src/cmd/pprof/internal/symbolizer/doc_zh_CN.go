// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package symbolizer provides a routine to populate a profile with symbol, file
// and line number information. It relies on the addr2liner and demangler packages
// to do the actual work.

// Package symbolizer provides a routine to
// populate a profile with symbol, file and
// line number information. It relies on
// the addr2liner and demangler packages to
// do the actual work.
package symbolizer

// Symbolize adds symbol and line number information to all locations in a profile.
// mode enables some options to control symbolization. Currently only recognizes
// "force", which causes it to overwrite any existing data.

// Symbolize adds symbol and line number
// information to all locations in a
// profile. mode enables some options to
// control symbolization. Currently only
// recognizes "force", which causes it to
// overwrite any existing data.
func Symbolize(mode string, prof *profile.Profile, obj plugin.ObjTool, ui plugin.UI) error
