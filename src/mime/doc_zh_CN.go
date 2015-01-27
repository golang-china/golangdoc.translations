// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package mime implements parts of the MIME spec.
package mime

// AddExtensionType sets the MIME type associated with the extension ext to typ.
// The extension should begin with a leading dot, as in ".html".
func AddExtensionType(ext, typ string) error

// FormatMediaType serializes mediatype t and the parameters param as a media type
// conforming to RFC 2045 and RFC 2616. The type and parameter names are written in
// lower-case. When any of the arguments result in a standard violation then
// FormatMediaType returns the empty string.
func FormatMediaType(t string, param map[string]string) string

// ParseMediaType parses a media type value and any optional parameters, per RFC
// 1521. Media types are the values in Content-Type and Content-Disposition headers
// (RFC 2183). On success, ParseMediaType returns the media type converted to
// lowercase and trimmed of white space and a non-nil map. The returned map,
// params, maps from the lowercase attribute to the attribute value with its case
// preserved.
func ParseMediaType(v string) (mediatype string, params map[string]string, err error)

// TypeByExtension returns the MIME type associated with the file extension ext.
// The extension ext should begin with a leading dot, as in ".html". When ext has
// no associated type, TypeByExtension returns "".
//
// Extensions are looked up first case-sensitively, then case-insensitively.
//
// The built-in table is small but on unix it is augmented by the local system's
// mime.types file(s) if available under one or more of these names:
//
//	/etc/mime.types
//	/etc/apache2/mime.types
//	/etc/apache/mime.types
//
// On Windows, MIME types are extracted from the registry.
//
// Text types have the charset parameter set to "utf-8" by default.
func TypeByExtension(ext string) string
