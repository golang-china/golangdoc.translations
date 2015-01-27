// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package svg provides tools related to handling of SVG files
package svg

// Massage enhances the SVG output from DOT to provide bettern panning inside a web
// browser. It uses the SVGPan library, which is accessed through the svgPan URL.
func Massage(in bytes.Buffer, svgPan string) string
