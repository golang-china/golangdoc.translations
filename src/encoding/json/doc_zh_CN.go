// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package json implements encoding and decoding of JSON objects as defined in RFC
// 4627. The mapping between JSON objects and Go values is described in the
// documentation for the Marshal and Unmarshal functions.
//
// See "JSON and Go" for an introduction to this package:
// http://golang.org/doc/articles/json_and_go.html
package json

// Compact appends to dst the JSON-encoded src with insignificant space characters
// elided.
func Compact(dst *bytes.Buffer, src []byte) error

// HTMLEscape appends to dst the JSON-encoded src with <, >, &, U+2028 and U+2029
// characters inside string literals changed to \u003c, \u003e, \u0026, \u2028,
// \u2029 so that the JSON will be safe to embed inside HTML <script> tags. For
// historical reasons, web browsers don't honor standard HTML escaping within
// <script> tags, so an alternative JSON encoding must be used.
func HTMLEscape(dst *bytes.Buffer, src []byte)

// Indent appends to dst an indented form of the JSON-encoded src. Each element in
// a JSON object or array begins on a new, indented line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
// The data appended to dst does not begin with the prefix nor any indentation, and
// has no trailing newline, to make it easier to embed inside other formatted JSON
// data.
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error

// Marshal returns the JSON encoding of v.
//
// Marshal traverses the value v recursively. If an encountered value implements
// the Marshaler interface and is not a nil pointer, Marshal calls its MarshalJSON
// method to produce JSON. The nil pointer exception is not strictly necessary but
// mimics a similar, necessary exception in the behavior of UnmarshalJSON.
//
// Otherwise, Marshal uses the following type-dependent default encodings:
//
// Boolean values encode as JSON booleans.
//
// Floating point, integer, and Number values encode as JSON numbers.
//
// String values encode as JSON strings coerced to valid UTF-8, replacing invalid
// bytes with the Unicode replacement rune. The angle brackets "<" and ">" are
// escaped to "\u003c" and "\u003e" to keep some browsers from misinterpreting JSON
// output as HTML. Ampersand "&" is also escaped to "\u0026" for the same reason.
//
// Array and slice values encode as JSON arrays, except that []byte encodes as a
// base64-encoded string, and a nil slice encodes as the null JSON object.
//
// Struct values encode as JSON objects. Each exported struct field becomes a
// member of the object unless
//
//	- the field's tag is "-", or
//	- the field is empty and its tag specifies the "omitempty" option.
//
// The empty values are false, 0, any nil pointer or interface value, and any
// array, slice, map, or string of length zero. The object's default key string is
// the struct field name but can be specified in the struct field's tag value. The
// "json" key in the struct field's tag value is the key name, followed by an
// optional comma and options. Examples:
//
//	// Field is ignored by this package.
//	Field int `json:"-"`
//
//	// Field appears in JSON as key "myName".
//	Field int `json:"myName"`
//
//	// Field appears in JSON as key "myName" and
//	// the field is omitted from the object if its value is empty,
//	// as defined above.
//	Field int `json:"myName,omitempty"`
//
//	// Field appears in JSON as key "Field" (the default), but
//	// the field is skipped if empty.
//	// Note the leading comma.
//	Field int `json:",omitempty"`
//
// The "string" option signals that a field is stored as JSON inside a JSON-encoded
// string. It applies only to fields of string, floating point, or integer types.
// This extra level of encoding is sometimes used when communicating with
// JavaScript programs:
//
//	Int64String int64 `json:",string"`
//
// The key name will be used if it's a non-empty string consisting of only Unicode
// letters, digits, dollar signs, percent signs, hyphens, underscores and slashes.
//
// Anonymous struct fields are usually marshaled as if their inner exported fields
// were fields in the outer struct, subject to the usual Go visibility rules
// amended as described in the next paragraph. An anonymous struct field with a
// name given in its JSON tag is treated as having that name, rather than being
// anonymous. An anonymous struct field of interface type is treated the same as
// having that type as its name, rather than being anonymous.
//
// The Go visibility rules for struct fields are amended for JSON when deciding
// which field to marshal or unmarshal. If there are multiple fields at the same
// level, and that level is the least nested (and would therefore be the nesting
// level selected by the usual Go rules), the following extra rules apply:
//
// 1) Of those fields, if any are JSON-tagged, only tagged fields are considered,
// even if there are multiple untagged fields that would otherwise conflict. 2) If
// there is exactly one field (tagged or not according to the first rule), that is
// selected. 3) Otherwise there are multiple fields, and all are ignored; no error
// occurs.
//
// Handling of anonymous struct fields is new in Go 1.1. Prior to Go 1.1, anonymous
// struct fields were ignored. To force ignoring of an anonymous struct field in
// both current and earlier versions, give the field a JSON tag of "-".
//
// Map values encode as JSON objects. The map's key type must be string; the object
// keys are used directly as map keys.
//
// Pointer values encode as the value pointed to. A nil pointer encodes as the null
// JSON object.
//
// Interface values encode as the value contained in the interface. A nil interface
// value encodes as the null JSON object.
//
// Channel, complex, and function values cannot be encoded in JSON. Attempting to
// encode such a value causes Marshal to return an UnsupportedTypeError.
//
// JSON cannot represent cyclic data structures and Marshal does not handle them.
// Passing cyclic structures to Marshal will result in an infinite recursion.

// Marshal returns the JSON encoding of v.
//
// Marshal traverses the value v recursively. If an encountered value implements
// the Marshaler interface and is not a nil pointer, Marshal calls its MarshalJSON
// method to produce JSON. The nil pointer exception is not strictly necessary but
// mimics a similar, necessary exception in the behavior of UnmarshalJSON.
//
// Otherwise, Marshal uses the following type-dependent default encodings:
//
// Boolean values encode as JSON booleans.
//
// Floating point, integer, and Number values encode as JSON numbers.
//
// String values encode as JSON strings coerced to valid UTF-8, replacing invalid
// bytes with the Unicode replacement rune. The angle brackets "<" and ">" are
// escaped to "\u003c" and "\u003e" to keep some browsers from misinterpreting JSON
// output as HTML. Ampersand "&" is also escaped to "\u0026" for the same reason.
//
// Array and slice values encode as JSON arrays, except that []byte encodes as a
// base64-encoded string, and a nil slice encodes as the null JSON object.
//
// Struct values encode as JSON objects. Each exported struct field becomes a
// member of the object unless
//
//	- the field's tag is "-", or
//	- the field is empty and its tag specifies the "omitempty" option.
//
// The empty values are false, 0, any nil pointer or interface value, and any
// array, slice, map, or string of length zero. The object's default key string is
// the struct field name but can be specified in the struct field's tag value. The
// "json" key in the struct field's tag value is the key name, followed by an
// optional comma and options. Examples:
//
//	// Field is ignored by this package.
//	Field int `json:"-"`
//
//	// Field appears in JSON as key "myName".
//	Field int `json:"myName"`
//
//	// Field appears in JSON as key "myName" and
//	// the field is omitted from the object if its value is empty,
//	// as defined above.
//	Field int `json:"myName,omitempty"`
//
//	// Field appears in JSON as key "Field" (the default), but
//	// the field is skipped if empty.
//	// Note the leading comma.
//	Field int `json:",omitempty"`
//
// The "string" option signals that a field is stored as JSON inside a JSON-encoded
// string. It applies only to fields of string, floating point, integer, or boolean
// types. This extra level of encoding is sometimes used when communicating with
// JavaScript programs:
//
//	Int64String int64 `json:",string"`
//
// The key name will be used if it's a non-empty string consisting of only Unicode
// letters, digits, dollar signs, percent signs, hyphens, underscores and slashes.
//
// Anonymous struct fields are usually marshaled as if their inner exported fields
// were fields in the outer struct, subject to the usual Go visibility rules
// amended as described in the next paragraph. An anonymous struct field with a
// name given in its JSON tag is treated as having that name, rather than being
// anonymous. An anonymous struct field of interface type is treated the same as
// having that type as its name, rather than being anonymous.
//
// The Go visibility rules for struct fields are amended for JSON when deciding
// which field to marshal or unmarshal. If there are multiple fields at the same
// level, and that level is the least nested (and would therefore be the nesting
// level selected by the usual Go rules), the following extra rules apply:
//
// 1) Of those fields, if any are JSON-tagged, only tagged fields are considered,
// even if there are multiple untagged fields that would otherwise conflict. 2) If
// there is exactly one field (tagged or not according to the first rule), that is
// selected. 3) Otherwise there are multiple fields, and all are ignored; no error
// occurs.
//
// Handling of anonymous struct fields is new in Go 1.1. Prior to Go 1.1, anonymous
// struct fields were ignored. To force ignoring of an anonymous struct field in
// both current and earlier versions, give the field a JSON tag of "-".
//
// Map values encode as JSON objects. The map's key type must be string; the object
// keys are used directly as map keys.
//
// Pointer values encode as the value pointed to. A nil pointer encodes as the null
// JSON object.
//
// Interface values encode as the value contained in the interface. A nil interface
// value encodes as the null JSON object.
//
// Channel, complex, and function values cannot be encoded in JSON. Attempting to
// encode such a value causes Marshal to return an UnsupportedTypeError.
//
// JSON cannot represent cyclic data structures and Marshal does not handle them.
// Passing cyclic structures to Marshal will result in an infinite recursion.
func Marshal(v interface{}) ([]byte, error)

// MarshalIndent is like Marshal but applies Indent to format the output.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)

// Unmarshal parses the JSON-encoded data and stores the result in the value
// pointed to by v.
//
// Unmarshal uses the inverse of the encodings that Marshal uses, allocating maps,
// slices, and pointers as necessary, with the following additional rules:
//
// To unmarshal JSON into a pointer, Unmarshal first handles the case of the JSON
// being the JSON literal null. In that case, Unmarshal sets the pointer to nil.
// Otherwise, Unmarshal unmarshals the JSON into the value pointed at by the
// pointer. If the pointer is nil, Unmarshal allocates a new value for it to point
// to.
//
// To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the
// keys used by Marshal (either the struct field name or its tag), preferring an
// exact match but also accepting a case-insensitive match.
//
// To unmarshal JSON into an interface value, Unmarshal stores one of these in the
// interface value:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	string, for JSON strings
//	[]interface{}, for JSON arrays
//	map[string]interface{}, for JSON objects
//	nil for JSON null
//
// If a JSON value is not appropriate for a given target type, or if a JSON number
// overflows the target type, Unmarshal skips that field and completes the
// unmarshalling as best it can. If no more serious errors are encountered,
// Unmarshal returns an UnmarshalTypeError describing the earliest such error.
//
// The JSON null value unmarshals into an interface, map, pointer, or slice by
// setting that Go value to nil. Because null is often used in JSON to mean ``not
// present,'' unmarshaling a JSON null into any other Go type has no effect on the
// value and produces no error.
//
// When unmarshaling quoted strings, invalid UTF-8 or invalid UTF-16 surrogate
// pairs are not treated as an error. Instead, they are replaced by the Unicode
// replacement character U+FFFD.
func Unmarshal(data []byte, v interface{}) error

// A Decoder reads and decodes JSON objects from an input stream.
type Decoder struct {
	// contains filtered or unexported fields
}

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may read data from r beyond the
// JSON values requested.
func NewDecoder(r io.Reader) *Decoder

// Buffered returns a reader of the data remaining in the Decoder's buffer. The
// reader is valid until the next call to Decode.
func (dec *Decoder) Buffered() io.Reader

// Decode reads the next JSON-encoded value from its input and stores it in the
// value pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion of JSON
// into a Go value.
func (dec *Decoder) Decode(v interface{}) error

// UseNumber causes the Decoder to unmarshal a number into an interface{} as a
// Number instead of as a float64.
func (dec *Decoder) UseNumber()

// An Encoder writes JSON objects to an output stream.
type Encoder struct {
	// contains filtered or unexported fields
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder

// Encode writes the JSON encoding of v to the stream, followed by a newline
// character.
//
// See the documentation for Marshal for details about the conversion of Go values
// to JSON.
func (enc *Encoder) Encode(v interface{}) error

// Before Go 1.2, an InvalidUTF8Error was returned by Marshal when attempting to
// encode a string value with invalid UTF-8 sequences. As of Go 1.2, Marshal
// instead coerces the string to valid UTF-8 by replacing invalid bytes with the
// Unicode replacement rune U+FFFD. This error is no longer generated but is kept
// for backwards compatibility with programs that might mention it.
type InvalidUTF8Error struct {
	S string // the whole string value that caused the error
}

func (e *InvalidUTF8Error) Error() string

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal. (The
// argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string

// Marshaler is the interface implemented by objects that can marshal themselves
// into valid JSON.
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type MarshalerError struct {
	Type reflect.Type
	Err  error
}

func (e *MarshalerError) Error() string

// A Number represents a JSON number literal.
type Number string

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error)

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error)

// String returns the literal text of the number.
func (n Number) String() string

// RawMessage is a raw encoded JSON object. It implements Marshaler and Unmarshaler
// and can be used to delay JSON decoding or precompute a JSON encoding.
type RawMessage []byte

// MarshalJSON returns *m as the JSON encoding of m.
func (m *RawMessage) MarshalJSON() ([]byte, error)

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error

// A SyntaxError is a description of a JSON syntax error.
type SyntaxError struct {
	Offset int64 // error occurred after reading Offset bytes
	// contains filtered or unexported fields
}

func (e *SyntaxError) Error() string

// An UnmarshalFieldError describes a JSON object key that led to an unexported
// (and therefore unwritable) struct field. (No longer used; kept for
// compatibility.)
type UnmarshalFieldError struct {
	Key   string
	Type  reflect.Type
	Field reflect.StructField
}

func (e *UnmarshalFieldError) Error() string

// An UnmarshalTypeError describes a JSON value that was not appropriate for a
// value of a specific Go type.
type UnmarshalTypeError struct {
	Value string       // description of JSON value - "bool", "array", "number -5"
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e *UnmarshalTypeError) Error() string

// Unmarshaler is the interface implemented by objects that can unmarshal a JSON
// description of themselves. The input can be assumed to be a valid encoding of a
// JSON value. UnmarshalJSON must copy the JSON data if it wishes to retain the
// data after returning.
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

// An UnsupportedTypeError is returned by Marshal when attempting to encode an
// unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string
