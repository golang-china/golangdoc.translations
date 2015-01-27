// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package html provides functions for escaping and unescaping HTML text.
package html

// EscapeString escapes special characters like "<" to become "&lt;". It escapes
// only five such characters: <, >, &, ' and ". UnescapeString(EscapeString(s)) ==
// s always holds, but the converse isn't always true.
func EscapeString(s string) string

// UnescapeString unescapes entities like "&lt;" to become "<". It unescapes a
// larger range of entities than EscapeString escapes. For example, "&aacute;"
// unescapes to "á", as does "&#225;" and "&xE1;". UnescapeString(EscapeString(s))
// == s always holds, but the converse isn't always true.
func UnescapeString(s string) string
