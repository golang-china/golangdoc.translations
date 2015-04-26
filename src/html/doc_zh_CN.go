// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package html provides functions for escaping and unescaping HTML text.

// html包提供了用于转义和解转义HTML文本的函数。
package html

// EscapeString escapes special characters like "<" to become "&lt;". It escapes
// only five such characters: <, >, &, ' and ". UnescapeString(EscapeString(s)) ==
// s always holds, but the converse isn't always true.

// EscapeString函数将特定的一些字符转为逸码后的字符实体，如"<"变成"&lt;"。
//
// 它只会修改五个字符：<、>、&、'、"。
//
// UnescapeString(EscapeString(s)) ==
// s总是成立，但是两个函数顺序反过来则不一定成立。
func EscapeString(s string) string

// UnescapeString unescapes entities like "&lt;" to become "<". It unescapes a
// larger range of entities than EscapeString escapes. For example, "&aacute;"
// unescapes to "á", as does "&#225;" and "&xE1;". UnescapeString(EscapeString(s))
// == s always holds, but the converse isn't always true.

// UnescapeString函数将逸码的字符实体如"&lt;"修改为原字符"<"。它会解码一个很大范围内的字符实体，远比函数EscapeString转码范围大得多。例如"&aacute;"解码为"á"，"&#225;"和"&xE1;"也会解码为该字符。
func UnescapeString(s string) string
