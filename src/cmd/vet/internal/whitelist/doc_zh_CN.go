// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package whitelist defines exceptions for the vet tool.

// Package whitelist defines exceptions for the vet tool.
package whitelist

// UnkeyedLiteral are types that are actually slices, but
// syntactically, we cannot tell whether the Typ in pkg.Typ{1, 2, 3}
// is a slice or a struct, so we whitelist all the standard package
// library's exported slice types.

// UnkeyedLiteral is a white list of types in the standard packages
// that are used with unkeyed literals we deem to be acceptable.
var UnkeyedLiteral = map[string]bool{

	"image/color.Alpha16": true,
	"image/color.Alpha":   true,
	"image/color.CMYK":    true,
	"image/color.Gray16":  true,
	"image/color.Gray":    true,
	"image/color.NRGBA64": true,
	"image/color.NRGBA":   true,
	"image/color.NYCbCrA": true,
	"image/color.RGBA64":  true,
	"image/color.RGBA":    true,
	"image/color.YCbCr":   true,
	"image.Point":         true,
	"image.Rectangle":     true,
	"image.Uniform":       true,

	"unicode.Range16": true,
}


