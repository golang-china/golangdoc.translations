// +build ingore

// Package json implements encoding and decoding of JSON as defined in
// RFC 4627. The mapping between JSON and Go values is described
// in the documentation for the Marshal and Unmarshal functions.
//
// See "JSON and Go" for an introduction to this package:
// https://golang.org/doc/articles/json_and_go.html

// json包实现了json对象的编解码，参见RFC
// 4627。Json对象和go类型的映射关系请参见Marshal和Unmarshal函数的文档。
//
// 参见"JSON and
// Go"获取本包的一个介绍：http://golang.org/doc/articles/json_and_go.html
package json

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

// A Decoder reads and decodes JSON values from an input stream.

// Decoder从输入流解码json对象
type Decoder struct {
}

// A Delim is a JSON array or object delimiter, one of [ ] { or }.
type Delim rune

// An Encoder writes JSON values to an output stream.

// Encoder将json对象写入输出流。
type Encoder struct {
}

// Before Go 1.2, an InvalidUTF8Error was returned by Marshal when
// attempting to encode a string value with invalid UTF-8 sequences.
// As of Go 1.2, Marshal instead coerces the string to valid UTF-8 by
// replacing invalid bytes with the Unicode replacement rune U+FFFD.
// This error is no longer generated but is kept for backwards compatibility
// with programs that might mention it.

// Go 1.2之前版本，当试图编码一个包含非法utf-8序列的字符串时会返回本错误。Go 1.2
// 及之后版本，编码器会强行将非法字节替换为unicode字符U+FFFD来使字符串合法。本错
// 误已不会再出现，但出于向后兼容考虑而保留。
type InvalidUTF8Error struct {
	S string // the whole string value that caused the error
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)

// InvalidUnmarshalError用于描述一个传递给解码器的非法参数。（解码器的参数必须是
// 非nil指针）
type InvalidUnmarshalError struct {
	Type reflect.Type
}

// Marshaler is the interface implemented by types that
// can marshal themselves into valid JSON.

// 实现了Marshaler接口的类型可以将自身序列化为合法的json描述。
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type MarshalerError struct {
	Type reflect.Type
	Err  error
}

// A Number represents a JSON number literal.

// Number类型代表一个json数字字面量。
type Number string

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.

// RawMessage类型是一个保持原本编码的json对象。本类型实现了Marshaler和
// Unmarshaler接口，用于延迟json的解码或者预计算json的编码。
type RawMessage []byte

// A SyntaxError is a description of a JSON syntax error.

// SyntaxError表示一个json语法错误。
type SyntaxError struct {
	Offset int64 // error occurred after reading Offset bytes
}

// A Token holds a value of one of these types:
//
// 	Delim, for the four JSON delimiters [ ] { }
// 	bool, for JSON booleans
// 	float64, for JSON numbers
// 	Number, for JSON numbers
// 	string, for JSON string literals
// 	nil, for JSON null
type Token interface {
}

// An UnmarshalFieldError describes a JSON object key that
// led to an unexported (and therefore unwritable) struct field.
// (No longer used; kept for compatibility.)

// UnmarshalFieldError表示一个json对象的键指向一个非导出字段。（因此不能写入；已
// 不再使用，出于兼容保留）
type UnmarshalFieldError struct {
	Key   string
	Type  reflect.Type
	Field reflect.StructField
}

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.

// UnmarshalTypeError表示一个json值不能转化为特定的go类型的值。
type UnmarshalTypeError struct {
	Value  string       // description of JSON value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Offset int64        // error occurred after reading Offset bytes
}

// Unmarshaler is the interface implemented by types
// that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.

// 实现了Unmarshaler接口的对象可以将自身的json描述反序列化。该方法可以认为输入是
// 合法的json字符串。如果要在方法返回后保存自身的json数据，必须进行拷贝。
type Unmarshaler interface {
	UnmarshalJSON([]byte)error
}

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.

// UnsupportedTypeError表示试图编码一个不支持类型的值。
type UnsupportedTypeError struct {
	Type reflect.Type
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

// Compact appends to dst the JSON-encoded src with
// insignificant space characters elided.

// Compact函数将json编码的src中无用的空白字符剔除后写入dst。
func Compact(dst *bytes.Buffer, src []byte) error

// HTMLEscape appends to dst the JSON-encoded src with <, >, &, U+2028 and
// U+2029 characters inside string literals changed to \u003c, \u003e, \u0026,
// \u2028, \u2029 so that the JSON will be safe to embed inside HTML <script>
// tags. For historical reasons, web browsers don't honor standard HTML escaping
// within <script> tags, so an alternative JSON encoding must be used.

// HTMLEscape 函数将json编码的src中的<、>、&、U+2028 和U+2029字符替换为\u003c、
// \u003e、\u0026、\u2028、\u2029 转义字符串，以便json编码可以安全的嵌入HTML的
// <script>标签里。因为历史原因，网络浏览器不支持在<script>标签中使用标准HTML转
// 义， 因此必须使用另一种json编码方案。
func HTMLEscape(dst *bytes.Buffer, src []byte)

// Indent appends to dst an indented form of the JSON-encoded src. Each element
// in a JSON object or array begins on a new, indented line beginning with
// prefix followed by one or more copies of indent according to the indentation
// nesting. The data appended to dst does not begin with the prefix nor any
// indentation, to make it easier to embed inside other formatted JSON data.
// Although leading space characters (space, tab, carriage return, newline) at
// the beginning of src are dropped, trailing space characters at the end of src
// are preserved and copied to dst. For example, if src has no trailing spaces,
// neither will dst; if src ends in a trailing newline, so will dst.

// Indent函数将json编码的调整缩进之后写入dst。每一个json元素/数组都另起一行开始
// ，以prefix为起始，一或多个indent缩进（数目看嵌套层数）。写入dst的数据起始没有
// prefix字符，也没有indent字符，最后也不换行，因此可以更好的嵌入其他格式化后的
// json数据里。
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error

// Marshal returns the JSON encoding of v.
//
// Marshal traverses the value v recursively. If an encountered value implements
// the Marshaler interface and is not a nil pointer, Marshal calls its
// MarshalJSON method to produce JSON. If no MarshalJSON method is present but
// the value implements encoding.TextMarshaler instead, Marshal calls its
// MarshalText method. The nil pointer exception is not strictly necessary but
// mimics a similar, necessary exception in the behavior of UnmarshalJSON.
//
// Otherwise, Marshal uses the following type-dependent default encodings:
//
// Boolean values encode as JSON booleans.
//
// Floating point, integer, and Number values encode as JSON numbers.
//
// String values encode as JSON strings coerced to valid UTF-8, replacing
// invalid bytes with the Unicode replacement rune. The angle brackets "<" and
// ">" are escaped to "\u003c" and "\u003e" to keep some browsers from
// misinterpreting JSON output as HTML. Ampersand "&" is also escaped to
// "\u0026" for the same reason. This escaping can be disabled using an Encoder
// with DisableHTMLEscaping.
//
// Array and slice values encode as JSON arrays, except that []byte encodes as a
// base64-encoded string, and a nil slice encodes as the null JSON value.
//
// Struct values encode as JSON objects. Each exported struct field becomes a
// member of the object unless
//
// 	- the field's tag is "-", or
// 	- the field is empty and its tag specifies the "omitempty" option.
//
// The empty values are false, 0, any nil pointer or interface value, and any
// array, slice, map, or string of length zero. The object's default key string
// is the struct field name but can be specified in the struct field's tag
// value. The "json" key in the struct field's tag value is the key name,
// followed by an optional comma and options. Examples:
//
// 	// Field is ignored by this package.
// 	Field int `json:"-"`
//
// 	// Field appears in JSON as key "myName".
// 	Field int `json:"myName"`
//
// 	// Field appears in JSON as key "myName" and
// 	// the field is omitted from the object if its value is empty,
// 	// as defined above.
// 	Field int `json:"myName,omitempty"`
//
// 	// Field appears in JSON as key "Field" (the default), but
// 	// the field is skipped if empty.
// 	// Note the leading comma.
// 	Field int `json:",omitempty"`
//
// The "string" option signals that a field is stored as JSON inside a
// JSON-encoded string. It applies only to fields of string, floating point,
// integer, or boolean types. This extra level of encoding is sometimes used
// when communicating with JavaScript programs:
//
// 	Int64String int64 `json:",string"`
//
// The key name will be used if it's a non-empty string consisting of only
// Unicode letters, digits, and ASCII punctuation except quotation marks,
// backslash, and comma.
//
// Anonymous struct fields are usually marshaled as if their inner exported
// fields were fields in the outer struct, subject to the usual Go visibility
// rules amended as described in the next paragraph. An anonymous struct field
// with a name given in its JSON tag is treated as having that name, rather than
// being anonymous. An anonymous struct field of interface type is treated the
// same as having that type as its name, rather than being anonymous.
//
// The Go visibility rules for struct fields are amended for JSON when deciding
// which field to marshal or unmarshal. If there are multiple fields at the same
// level, and that level is the least nested (and would therefore be the nesting
// level selected by the usual Go rules), the following extra rules apply:
//
// 1) Of those fields, if any are JSON-tagged, only tagged fields are
// considered, even if there are multiple untagged fields that would otherwise
// conflict. 2) If there is exactly one field (tagged or not according to the
// first rule), that is selected. 3) Otherwise there are multiple fields, and
// all are ignored; no error occurs.
//
// Handling of anonymous struct fields is new in Go 1.1. Prior to Go 1.1,
// anonymous struct fields were ignored. To force ignoring of an anonymous
// struct field in both current and earlier versions, give the field a JSON tag
// of "-".
//
// Map values encode as JSON objects. The map's key type must either be a
// string, an integer type, or implement encoding.TextMarshaler. The map keys
// are sorted and used as JSON object keys by applying the following rules,
// subject to the UTF-8 coercion described for string values above:
//
// 	- string keys are used directly
// 	- encoding.TextMarshalers are marshaled
// 	- integer keys are converted to strings
//
// Pointer values encode as the value pointed to. A nil pointer encodes as the
// null JSON value.
//
// Interface values encode as the value contained in the interface. A nil
// interface value encodes as the null JSON value.
//
// Channel, complex, and function values cannot be encoded in JSON. Attempting
// to encode such a value causes Marshal to return an UnsupportedTypeError.
//
// JSON cannot represent cyclic data structures and Marshal does not handle
// them. Passing cyclic structures to Marshal will result in an infinite
// recursion.

// Marshal returns the JSON encoding of v.
//
// Marshal traverses the value v recursively. If an encountered value implements
// the Marshaler interface and is not a nil pointer, Marshal calls its
// MarshalJSON method to produce JSON. If no MarshalJSON method is present but
// the value implements encoding.TextMarshaler instead, Marshal calls its
// MarshalText method. The nil pointer exception is not strictly necessary but
// mimics a similar, necessary exception in the behavior of UnmarshalJSON.
//
// Otherwise, Marshal uses the following type-dependent default encodings:
//
// Boolean values encode as JSON booleans.
//
// Floating point, integer, and Number values encode as JSON numbers.
//
// String values encode as JSON strings coerced to valid UTF-8, replacing
// invalid bytes with the Unicode replacement rune. The angle brackets "<" and
// ">" are escaped to "\u003c" and "\u003e" to keep some browsers from
// misinterpreting JSON output as HTML. Ampersand "&" is also escaped to
// "\u0026" for the same reason. This escaping can be disabled using an Encoder
// with DisableHTMLEscaping.
//
// Array and slice values encode as JSON arrays, except that []byte encodes as a
// base64-encoded string, and a nil slice encodes as the null JSON value.
//
// Struct values encode as JSON objects. Each exported struct field becomes a
// member of the object unless
//
// 	- the field's tag is "-", or
// 	- the field is empty and its tag specifies the "omitempty" option.
//
// The empty values are false, 0, any nil pointer or interface value, and any
// array, slice, map, or string of length zero. The object's default key string
// is the struct field name but can be specified in the struct field's tag
// value. The "json" key in the struct field's tag value is the key name,
// followed by an optional comma and options. Examples:
//
// 	// Field is ignored by this package.
// 	Field int `json:"-"`
//
// 	// Field appears in JSON as key "myName".
// 	Field int `json:"myName"`
//
// 	// Field appears in JSON as key "myName" and
// 	// the field is omitted from the object if its value is empty,
// 	// as defined above.
// 	Field int `json:"myName,omitempty"`
//
// 	// Field appears in JSON as key "Field" (the default), but
// 	// the field is skipped if empty.
// 	// Note the leading comma.
// 	Field int `json:",omitempty"`
//
// The "string" option signals that a field is stored as JSON inside a
// JSON-encoded string. It applies only to fields of string, floating point,
// integer, or boolean types. This extra level of encoding is sometimes used
// when communicating with JavaScript programs:
//
// 	Int64String int64 `json:",string"`
//
// The key name will be used if it's a non-empty string consisting of only
// Unicode letters, digits, dollar signs, percent signs, hyphens, underscores
// and slashes.
//
// Anonymous struct fields are usually marshaled as if their inner exported
// fields were fields in the outer struct, subject to the usual Go visibility
// rules amended as described in the next paragraph. An anonymous struct field
// with a name given in its JSON tag is treated as having that name, rather than
// being anonymous. An anonymous struct field of interface type is treated the
// same as having that type as its name, rather than being anonymous.
//
// The Go visibility rules for struct fields are amended for JSON when deciding
// which field to marshal or unmarshal. If there are multiple fields at the same
// level, and that level is the least nested (and would therefore be the nesting
// level selected by the usual Go rules), the following extra rules apply:
//
// 1) Of those fields, if any are JSON-tagged, only tagged fields are
// considered, even if there are multiple untagged fields that would otherwise
// conflict. 2) If there is exactly one field (tagged or not according to the
// first rule), that is selected. 3) Otherwise there are multiple fields, and
// all are ignored; no error occurs.
//
// Handling of anonymous struct fields is new in Go 1.1. Prior to Go 1.1,
// anonymous struct fields were ignored. To force ignoring of an anonymous
// struct field in both current and earlier versions, give the field a JSON tag
// of "-".
//
// Map values encode as JSON objects. The map's key type must either be a
// string, an integer type, or implement encoding.TextMarshaler. The map keys
// are sorted and used as JSON object keys by applying the following rules,
// subject to the UTF-8 coercion described for string values above:
//
// 	- string keys are used directly
// 	- encoding.TextMarshalers are marshaled
// 	- integer keys are converted to strings
//
// Pointer values encode as the value pointed to. A nil pointer encodes as the
// null JSON value.
//
// Interface values encode as the value contained in the interface. A nil
// interface value encodes as the null JSON value.
//
// Channel, complex, and function values cannot be encoded in JSON. Attempting
// to encode such a value causes Marshal to return an UnsupportedTypeError.
//
// JSON cannot represent cyclic data structures and Marshal does not handle
// them. Passing cyclic structures to Marshal will result in an infinite
// recursion.
func Marshal(v interface{}) ([]byte, error)

// MarshalIndent is like Marshal but applies Indent to format the output.

// MarshalIndent类似Marshal但会使用缩进将输出格式化。
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may
// read data from r beyond the JSON values requested.

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may read data from r beyond the
// JSON values requested.
func NewDecoder(r io.Reader) *Decoder

// NewEncoder returns a new encoder that writes to w.

// NewEncoder创建一个将数据写入w的*Encoder。
func NewEncoder(w io.Writer) *Encoder

// Unmarshal parses the JSON-encoded data and stores the result in the value
// pointed to by v.
//
// Unmarshal uses the inverse of the encodings that Marshal uses, allocating
// maps, slices, and pointers as necessary, with the following additional rules:
//
// To unmarshal JSON into a pointer, Unmarshal first handles the case of the
// JSON being the JSON literal null. In that case, Unmarshal sets the pointer to
// nil. Otherwise, Unmarshal unmarshals the JSON into the value pointed at by
// the pointer. If the pointer is nil, Unmarshal allocates a new value for it to
// point to.
//
// To unmarshal JSON into a struct, Unmarshal matches incoming object keys to
// the keys used by Marshal (either the struct field name or its tag),
// preferring an exact match but also accepting a case-insensitive match.
// Unmarshal will only set exported fields of the struct.
//
// To unmarshal JSON into an interface value, Unmarshal stores one of these in
// the interface value:
//
// 	bool, for JSON booleans
// 	float64, for JSON numbers
// 	string, for JSON strings
// 	[]interface{}, for JSON arrays
// 	map[string]interface{}, for JSON objects
// 	nil for JSON null
//
// To unmarshal a JSON array into a slice, Unmarshal resets the slice length to
// zero and then appends each element to the slice. As a special case, to
// unmarshal an empty JSON array into a slice, Unmarshal replaces the slice with
// a new empty slice.
//
// To unmarshal a JSON array into a Go array, Unmarshal decodes JSON array
// elements into corresponding Go array elements. If the Go array is smaller
// than the JSON array, the additional JSON array elements are discarded. If the
// JSON array is smaller than the Go array, the additional Go array elements are
// set to zero values.
//
// To unmarshal a JSON object into a map, Unmarshal first establishes a map to
// use. If the map is nil, Unmarshal allocates a new map. Otherwise Unmarshal
// reuses the existing map, keeping existing entries. Unmarshal then stores key-
// value pairs from the JSON object into the map. The map's key type must either
// be a string, an integer, or implement encoding.TextUnmarshaler.
//
// If a JSON value is not appropriate for a given target type, or if a JSON
// number overflows the target type, Unmarshal skips that field and completes
// the unmarshaling as best it can. If no more serious errors are encountered,
// Unmarshal returns an UnmarshalTypeError describing the earliest such error.
//
// The JSON null value unmarshals into an interface, map, pointer, or slice by
// setting that Go value to nil. Because null is often used in JSON to mean
// ``not present,'' unmarshaling a JSON null into any other Go type has no
// effect on the value and produces no error.
//
// When unmarshaling quoted strings, invalid UTF-8 or invalid UTF-16 surrogate
// pairs are not treated as an error. Instead, they are replaced by the Unicode
// replacement character U+FFFD.

// Unmarshal函数解析json编码的数据并将结果存入v指向的值。
//
// Unmarshal和Marshal做相反的操作，必要时申请映射、切片或指针，有如下的附加规则
// ：
//
// 要将json数据解码写入一个指针，Unmarshal函数首先处理json数据是json字面值null的
// 情况。此时，函数将指针设为nil；否则，函数将json数据解码写入指针指向的值；如果
// 指针本身是nil，函数会先申请一个值并使指针指向它。
//
// 要将json数据解码写入一个结构体，函数会匹配输入对象的键和Marshal使用的键（结构
// 体字段名或者它的标签指定的键名），优先选择精确的匹配，但也接受大小写不敏感的
// 匹配。
//
// 要将json数据解码写入一个接口类型值，函数会将数据解码为如下类型写入接口：
//
// 	Bool                   对应JSON布尔类型
// 	float64                对应JSON数字类型
// 	string                 对应JSON字符串类型
// 	[]interface{}          对应JSON数组
// 	map[string]interface{} 对应JSON对象
// 	nil                    对应JSON的null
//
// 如果一个JSON值不匹配给出的目标类型，或者如果一个json数字写入目标类型时溢出，
// Unmarshal函数会跳过该字段并尽量完成其余的解码操作。如果没有出现更加严重的错误
// ，本函数会返回一个描述第一个此类错误的详细信息的UnmarshalTypeError。
//
// JSON的null值解码为go的接口、指针、切片时会将它们设为nil，因为null在json里一般
// 表示“不存在”。 解码json的null值到其他go类型时，不会造成任何改变，也不会产生
// 错误。
//
// 当解码字符串时，不合法的utf-8或utf-16代理（字符）对不视为错误，而是将非法字符
// 替换为unicode字符U+FFFD。
func Unmarshal(data []byte, v interface{}) error

// Buffered returns a reader of the data remaining in the Decoder's
// buffer. The reader is valid until the next call to Decode.

// Buffered returns a reader of the data remaining in the Decoder's buffer. The
// reader is valid until the next call to Decode.
func (dec *Decoder) Buffered() io.Reader

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
//
// See the documentation for Unmarshal for details about
// the conversion of JSON into a Go value.

// Decode reads the next JSON-encoded value from its input and stores it in the
// value pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion of JSON
// into a Go value.
func (dec *Decoder) Decode(v interface{}) error

// More reports whether there is another element in the
// current array or object being parsed.
func (dec *Decoder) More() bool

// Token returns the next JSON token in the input stream.
// At the end of the input stream, Token returns nil, io.EOF.
//
// Token guarantees that the delimiters [ ] { } it returns are
// properly nested and matched: if Token encounters an unexpected
// delimiter in the input, it will return an error.
//
// The input stream consists of basic JSON values—bool, string,
// number, and null—along with delimiters [ ] { } of type Delim
// to mark the start and end of arrays and objects.
// Commas and colons are elided.
func (dec *Decoder) Token() (Token, error)

// UseNumber causes the Decoder to unmarshal a number into an interface{} as a
// Number instead of as a float64.
func (dec *Decoder) UseNumber()

// Encode writes the JSON encoding of v to the stream,
// followed by a newline character.
//
// See the documentation for Marshal for details about the
// conversion of Go values to JSON.

// Encode将v的json编码写入输出流，并会写入一个换行符，参见Marshal函数的文档获取
// 细节信息。
func (enc *Encoder) Encode(v interface{}) error

// SetEscapeHTML specifies whether problematic HTML characters should be escaped
// inside JSON quoted strings. The default behavior is to escape &, <, and > to
// \u0026, \u003c, and \u003e to avoid certain safety problems that can arise
// when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability of
// the output, SetEscapeHTML(false) disables this behavior.
func (enc *Encoder) SetEscapeHTML(on bool)

// SetIndent instructs the encoder to format each subsequent encoded value as if
// indented by the package-level function Indent(dst, src, prefix, indent).
// Calling SetIndent("", "") disables indentation.
func (enc *Encoder) SetIndent(prefix, indent string)

func (e *InvalidUTF8Error) Error() string

func (e *InvalidUnmarshalError) Error() string

func (e *MarshalerError) Error() string

// MarshalJSON returns *m as the JSON encoding of m.
func (m *RawMessage) MarshalJSON() ([]byte, error)

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error

func (e *SyntaxError) Error() string

func (e *UnmarshalFieldError) Error() string

func (e *UnmarshalTypeError) Error() string

func (e *UnsupportedTypeError) Error() string

func (e *UnsupportedValueError) Error() string

func (d Delim) String() string

// Float64 returns the number as a float64.

// 将该数字作为float64类型返回。
func (n Number) Float64() (float64, error)

// Int64 returns the number as an int64.

// 将该数字作为int64类型返回。
func (n Number) Int64() (int64, error)

// String returns the literal text of the number.

// 返回该数字的字面值文本表示。
func (n Number) String() string

