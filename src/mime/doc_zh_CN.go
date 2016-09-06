// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package mime implements parts of the MIME spec.

// mime实现了MIME的部分规定。
package mime

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

const (
	// BEncoding represents Base64 encoding scheme as defined by RFC 2045.
	BEncoding = WordEncoder('b')

	// QEncoding represents the Q-encoding scheme as defined by RFC 2047.
	QEncoding = WordEncoder('q')
)

// A WordDecoder decodes MIME headers containing RFC 2047 encoded-words.
type WordDecoder struct {
	// CharsetReader, if non-nil, defines a function to generate
	// charset-conversion readers, converting from the provided charset into
	// UTF-8. Charsets are always lower-case. utf-8, iso-8859-1 and us-ascii
	// charsets are handled by default. One of the the CharsetReader's result
	// values must be non-nil.
	CharsetReader func(charset string, input io.Reader) (io.Reader, error)
}

// A WordEncoder is an RFC 2047 encoded-word encoder.
type WordEncoder byte

// AddExtensionType sets the MIME type associated with
// the extension ext to typ. The extension should begin with
// a leading dot, as in ".html".

// 函数将扩展名和mimetype建立偶联；扩展名应以点号开始，例如".html"。
func AddExtensionType(ext, typ string) error

// ExtensionsByType returns the extensions known to be associated with the MIME
// type typ. The returned extensions will each begin with a leading dot, as in
// ".html". When typ has no associated extensions, ExtensionsByType returns an
// nil slice.
func ExtensionsByType(typ string) ([]string, error)

// FormatMediaType serializes mediatype t and the parameters
// param as a media type conforming to RFC 2045 and RFC 2616.
// The type and parameter names are written in lower-case.
// When any of the arguments result in a standard violation then
// FormatMediaType returns the empty string.

// 函数根据RFC 2045和 RFC 2616的规定将媒体类型t和参数param连接为一个mime媒体类型
// ，类型和参数都采用小写字母。任一个参数不合法都会返回空字符串。
func FormatMediaType(t string, param map[string]string) string

// ParseMediaType parses a media type value and any optional
// parameters, per RFC 1521.  Media types are the values in
// Content-Type and Content-Disposition headers (RFC 2183).
// On success, ParseMediaType returns the media type converted
// to lowercase and trimmed of white space and a non-nil map.
// The returned map, params, maps from the lowercase
// attribute to the attribute value with its case preserved.

// 函数根据RFC 1521解析一个媒体类型值以及可能的参数。媒体类型值一般应为
// Content-Type和Conten-Disposition头域的值（参见RFC 2183）。成功的调用会返回小
// 写字母、去空格的媒体类型和一个非空的map。返回的map映射小写字母的属性和对应的
// 属性值。
func ParseMediaType(v string) (mediatype string, params map[string]string, err error)

// TypeByExtension returns the MIME type associated with the file extension ext.
// The extension ext should begin with a leading dot, as in ".html".
// When ext has no associated type, TypeByExtension returns "".
//
// Extensions are looked up first case-sensitively, then case-insensitively.
//
// The built-in table is small but on unix it is augmented by the local
// system's mime.types file(s) if available under one or more of these
// names:
//
//   /etc/mime.types
//   /etc/apache2/mime.types
//   /etc/apache/mime.types
//
// On Windows, MIME types are extracted from the registry.
//
// Text types have the charset parameter set to "utf-8" by default.

// 函数返回与扩展名偶联的MIME类型。扩展名应以点号开始，如".html"。如果扩展名未偶
// 联类型，函数会返回""。
//
// 内建的偶联表很小，但在unix系统会从本地系统的一或多个mime.types文件（参加下表
// ）进行增补。
//
// 	/etc/mime.types
// 	/etc/apache2/mime.types
// 	/etc/apache/mime.types
//
// Windows系统的mime类型从注册表获取。文本类型的字符集参数默认设置为"utf-8"。
func TypeByExtension(ext string) string

// Decode decodes an RFC 2047 encoded-word.
func (d *WordDecoder) Decode(word string) (string, error)

// DecodeHeader decodes all encoded-words of the given string. It returns an
// error if and only if CharsetReader of d returns an error.
func (d *WordDecoder) DecodeHeader(header string) (string, error)

// Encode returns the encoded-word form of s. If s is ASCII without special
// characters, it is returned unchanged. The provided charset is the IANA
// charset name of s. It is case insensitive.
func (e WordEncoder) Encode(charset, s string) string

