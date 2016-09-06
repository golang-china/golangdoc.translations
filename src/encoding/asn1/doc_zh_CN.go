// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package asn1 implements parsing of DER-encoded ASN.1 data structures,
// as defined in ITU-T Rec X.690.
//
// See also ``A Layman's Guide to a Subset of ASN.1, BER, and DER,''
// http://luca.ntop.org/Teaching/Appunti/asn1.html.

// asn1包实现了DER编码的ASN.1数据结构的解析，参见ITU-T Rec X.690。
//
// 其他细节参见"A Layman's Guide to a Subset of ASN.1, BER, and DER"。
//
// 网址http://luca.ntop.org/Teaching/Appunti/asn1.html
package asn1

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// ASN.1 class types represent the namespace of the tag.
const (
	ClassUniversal       = 0
	ClassApplication     = 1
	ClassContextSpecific = 2
	ClassPrivate         = 3
)

// ASN.1 tags represent the type of the following object.
const (
	TagBoolean         = 1
	TagInteger         = 2
	TagBitString       = 3
	TagOctetString     = 4
	TagOID             = 6
	TagEnum            = 10
	TagUTF8String      = 12
	TagSequence        = 16
	TagSet             = 17
	TagPrintableString = 19
	TagT61String       = 20
	TagIA5String       = 22
	TagUTCTime         = 23
	TagGeneralizedTime = 24
	TagGeneralString   = 27
)

// BitString is the structure to use when you want an ASN.1 BIT STRING type. A
// bit string is padded up to the nearest byte in memory and the number of
// valid bits is recorded. Padding bits will be zero.

// BitString类型是用于表示ASN.1 BIT STRING类型的结构体。字位流补齐到最近的字节数
// 保存在内存里并记录合法字位数，补齐的位可以为0个。
type BitString struct {
	Bytes     []byte // bits packed into bytes.
	BitLength int    // length in bits.
}

// An Enumerated is represented as a plain int.

// Enumerated表示一个明文整数。
type Enumerated int

// A Flag accepts any data and is set to true if present.

// Flag接收任何数据，如果数据存在就设自身为真。
type Flag bool

// An ObjectIdentifier represents an ASN.1 OBJECT IDENTIFIER.

// ObjectIdentifier类型用于表示ASN.1 OBJECT IDENTIFIER类型。
type ObjectIdentifier []int

// RawContent is used to signal that the undecoded, DER data needs to be
// preserved for a struct. To use it, the first field of the struct must have
// this type. It's an error for any of the other fields to have this type.

// RawContent用于标记未解码的应被结构体保留的DER数据。如要使用它，结构体的第一个
// 字段必须是本类型，其它字段不能是本类型。
type RawContent []byte

// A RawValue represents an undecoded ASN.1 object.

// RawValue代表一个未解码的ASN.1对象。
type RawValue struct {
	Class, Tag int
	IsCompound bool
	Bytes      []byte
	FullBytes  []byte // includes the tag and length
}

// A StructuralError suggests that the ASN.1 data is valid, but the Go type
// which is receiving it doesn't match.

// StructuralError表示ASN.1数据合法但接收的Go类型不匹配。
type StructuralError struct {
	Msg string
}

// A SyntaxError suggests that the ASN.1 data is invalid.

// SyntaxErrorLeixing表示ASN.1数据不合法。
type SyntaxError struct {
	Msg string
}

// Marshal returns the ASN.1 encoding of val.
//
// In addition to the struct tags recognised by Unmarshal, the following can be
// used:
//
// 	ia5:		causes strings to be marshaled as ASN.1, IA5 strings
// 	omitempty:	causes empty slices to be skipped
// 	printable:	causes strings to be marshaled as ASN.1, PrintableString strings.
// 	utf8:		causes strings to be marshaled as ASN.1, UTF8 strings

// Marshal函数返回val的ASN.1编码。
//
// 此外还提供了供Unmarshal函数识别的结构体标签，可用如下标签：
//
//     ia5:           使字符串序列化为ASN.1 IA5String类型
//     omitempty:     使空切片被跳过
//     printable:     使字符串序列化为ASN.1 PrintableString类型
//     utf8:          使字符串序列化为ASN.1 UTF8字符串
func Marshal(val interface{}) ([]byte, error)

// Unmarshal parses the DER-encoded ASN.1 data structure b and uses the reflect
// package to fill in an arbitrary value pointed at by val. Because Unmarshal
// uses the reflect package, the structs being written to must use upper case
// field names.
//
// An ASN.1 INTEGER can be written to an int, int32, int64, or *big.Int (from
// the math/big package). If the encoded value does not fit in the Go type,
// Unmarshal returns a parse error.
//
// An ASN.1 BIT STRING can be written to a BitString.
//
// An ASN.1 OCTET STRING can be written to a []byte.
//
// An ASN.1 OBJECT IDENTIFIER can be written to an ObjectIdentifier.
//
// An ASN.1 ENUMERATED can be written to an Enumerated.
//
// An ASN.1 UTCTIME or GENERALIZEDTIME can be written to a time.Time.
//
// An ASN.1 PrintableString or IA5String can be written to a string.
//
// Any of the above ASN.1 values can be written to an interface{}. The value
// stored in the interface has the corresponding Go type. For integers, that
// type is int64.
//
// An ASN.1 SEQUENCE OF x or SET OF x can be written to a slice if an x can be
// written to the slice's element type.
//
// An ASN.1 SEQUENCE or SET can be written to a struct if each of the elements
// in the sequence can be written to the corresponding element in the struct.
//
// The following tags on struct fields have special meaning to Unmarshal:
//
// 	application	specifies that a APPLICATION tag is used
// 	default:x	sets the default value for optional integer fields
// 	explicit	specifies that an additional, explicit tag wraps the implicit one
// 	optional	marks the field as ASN.1 OPTIONAL
// 	set		causes a SET, rather than a SEQUENCE type to be expected
// 	tag:x		specifies the ASN.1 tag number; implies ASN.1 CONTEXT SPECIFIC
//
// If the type of the first field of a structure is RawContent then the raw ASN1
// contents of the struct will be stored in it.
//
// If the type name of a slice element ends with "SET" then it's treated as if
// the "set" tag was set on it. This can be used with nested slices where a
// struct tag cannot be given.
//
// Other ASN.1 types are not supported; if it encounters them, Unmarshal returns
// a parse error.

// Unmarshal函数解析DER编码的ASN.1结构体数据并使用reflect包填写val指向的任意类型
// 值。因为本函数使用了reflect包，结构体必须使用大写字母起始的字段名。
//
// ASN.1 INTEGER 类型值可以写入int、int32、int64或*big.Int（math/big包）类型。类
// 型不匹配会返回解析错误。
//
// ASN.1 BIT STRING类型值可以写入BitString类型。
//
// ASN.1 OCTET STRING类型值可以写入[]byte类型。
//
// ASN.1 OBJECT IDENTIFIER类型值可以写入ObjectIdentifier类型。
//
// ASN.1 ENUMERATED类型值可以写入Enumerated类型。
//
// ASN.1 UTCTIME类型值或GENERALIZEDTIME 类型值可以写入time.Time类型。
//
// ASN.1 PrintableString类型值或者IA5String类型值可以写入string类型。
//
// 以上任一ASN.1类型值都可写入interface{}类型。保存在接口里的类型为对应的Go类型
// ，ASN.1整型对应int64。
//
// 如果类型x可以写入切片的成员类型，则类型x的ASN.1 SEQUENCE或SET类型可以写入该切
// 片。
//
// ASN.1 SEQUENCE或SET类型如果其每一个成员都可以写入某结构体的对应字段，则可以写
// 入该结构体
//
// 对Unmarshal函数，下列字段标签有特殊含义：
//
// 	application    指明使用了APPLICATION标签
// 	default:x      设置一个可选整数字段的默认值
// 	explicit       给一个隐式的标签设置一个额外的显式标签
// 	optional       标记字段为ASN.1 OPTIONAL的
// 	set            表示期望一个SET而不是SEQUENCE类型
// 	tag:x          指定ASN.1标签码，隐含ASN.1 CONTEXT SPECIFIC
//
// 如果结构体的第一个字段的类型为RawContent，则会将原始ASN1结构体内容包存在该字
// 段。
//
// 如果切片成员的类型名以"SET"结尾，则视为该字段有"set"标签。这是给不能使用标签
// 的嵌套切片使用的。
//
// 其它ASN.1类型不支持，如果遭遇这些类型，Unmarshal返回解析错误。
func Unmarshal(b []byte, val interface{}) (rest []byte, err error)

// UnmarshalWithParams allows field parameters to be specified for the
// top-level element. The form of the params is the same as the field tags.

// UnmarshalWithParams允许指定val顶层成员的字段参数，格式和字段标签相同。
func UnmarshalWithParams(b []byte, val interface{}, params string) (rest []byte, err error)

// At returns the bit at the given index. If the index is out of range it
// returns false.

// At方法发挥index位置的字位，如果index出界则返回0。
func (b BitString) At(i int) int

// RightAlign returns a slice where the padding bits are at the beginning. The
// slice may share memory with the BitString.

// RightAlign方法返回b表示的字位流的右对齐版本（即补位在开始部分）切片，该切片可
// 能和b共享底层内存。
func (b BitString) RightAlign() []byte

// Equal reports whether oi and other represent the same identifier.

// 如果oi和other代表同一个标识符，Equal方法返回真。
func (oi ObjectIdentifier) Equal(other ObjectIdentifier) bool

func (oi ObjectIdentifier) String() string

func (e StructuralError) Error() string

func (e SyntaxError) Error() string

