// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package asn1 implements parsing of DER-encoded ASN.1 data structures, as defined
// in ITU-T Rec X.690.
//
// See also ``A Layman's Guide to a Subset of ASN.1, BER, and DER,''
// http://luca.ntop.org/Teaching/Appunti/asn1.html.

// Package asn1 implements parsing of
// DER-encoded ASN.1 data structures, as
// defined in ITU-T Rec X.690.
//
// See also ``A Layman's Guide to a Subset
// of ASN.1, BER, and DER,''
// http://luca.ntop.org/Teaching/Appunti/asn1.html.
package asn1

// Marshal returns the ASN.1 encoding of val.
//
// In addition to the struct tags recognised by Unmarshal, the following can be
// used:
//
//	ia5:		causes strings to be marshaled as ASN.1, IA5 strings
//	omitempty:	causes empty slices to be skipped
//	printable:	causes strings to be marshaled as ASN.1, PrintableString strings.
//	utf8:		causes strings to be marshaled as ASN.1, UTF8 strings

// Marshal returns the ASN.1 encoding of
// val.
//
// In addition to the struct tags
// recognised by Unmarshal, the following
// can be used:
//
//	ia5:		causes strings to be marshaled as ASN.1, IA5 strings
//	omitempty:	causes empty slices to be skipped
//	printable:	causes strings to be marshaled as ASN.1, PrintableString strings.
//	utf8:		causes strings to be marshaled as ASN.1, UTF8 strings
func Marshal(val interface{}) ([]byte, error)

// Unmarshal parses the DER-encoded ASN.1 data structure b and uses the reflect
// package to fill in an arbitrary value pointed at by val. Because Unmarshal uses
// the reflect package, the structs being written to must use upper case field
// names.
//
// An ASN.1 INTEGER can be written to an int, int32, int64, or *big.Int (from the
// math/big package). If the encoded value does not fit in the Go type, Unmarshal
// returns a parse error.
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
// Any of the above ASN.1 values can be written to an interface{}. The value stored
// in the interface has the corresponding Go type. For integers, that type is
// int64.
//
// An ASN.1 SEQUENCE OF x or SET OF x can be written to a slice if an x can be
// written to the slice's element type.
//
// An ASN.1 SEQUENCE or SET can be written to a struct if each of the elements in
// the sequence can be written to the corresponding element in the struct.
//
// The following tags on struct fields have special meaning to Unmarshal:
//
//	application	specifies that a APPLICATION tag is used
//	default:x	sets the default value for optional integer fields
//	explicit	specifies that an additional, explicit tag wraps the implicit one
//	optional	marks the field as ASN.1 OPTIONAL
//	set		causes a SET, rather than a SEQUENCE type to be expected
//	tag:x		specifies the ASN.1 tag number; implies ASN.1 CONTEXT SPECIFIC
//
// If the type of the first field of a structure is RawContent then the raw ASN1
// contents of the struct will be stored in it.
//
// If the type name of a slice element ends with "SET" then it's treated as if the
// "set" tag was set on it. This can be used with nested slices where a struct tag
// cannot be given.
//
// Other ASN.1 types are not supported; if it encounters them, Unmarshal returns a
// parse error.

// Unmarshal parses the DER-encoded ASN.1
// data structure b and uses the reflect
// package to fill in an arbitrary value
// pointed at by val. Because Unmarshal
// uses the reflect package, the structs
// being written to must use upper case
// field names.
//
// An ASN.1 INTEGER can be written to an
// int, int32, int64, or *big.Int (from the
// math/big package). If the encoded value
// does not fit in the Go type, Unmarshal
// returns a parse error.
//
// An ASN.1 BIT STRING can be written to a
// BitString.
//
// An ASN.1 OCTET STRING can be written to
// a []byte.
//
// An ASN.1 OBJECT IDENTIFIER can be
// written to an ObjectIdentifier.
//
// An ASN.1 ENUMERATED can be written to an
// Enumerated.
//
// An ASN.1 UTCTIME or GENERALIZEDTIME can
// be written to a time.Time.
//
// An ASN.1 PrintableString or IA5String
// can be written to a string.
//
// Any of the above ASN.1 values can be
// written to an interface{}. The value
// stored in the interface has the
// corresponding Go type. For integers,
// that type is int64.
//
// An ASN.1 SEQUENCE OF x or SET OF x can
// be written to a slice if an x can be
// written to the slice's element type.
//
// An ASN.1 SEQUENCE or SET can be written
// to a struct if each of the elements in
// the sequence can be written to the
// corresponding element in the struct.
//
// The following tags on struct fields have
// special meaning to Unmarshal:
//
//	application	specifies that a APPLICATION tag is used
//	default:x	sets the default value for optional integer fields
//	explicit	specifies that an additional, explicit tag wraps the implicit one
//	optional	marks the field as ASN.1 OPTIONAL
//	set		causes a SET, rather than a SEQUENCE type to be expected
//	tag:x		specifies the ASN.1 tag number; implies ASN.1 CONTEXT SPECIFIC
//
// If the type of the first field of a
// structure is RawContent then the raw
// ASN1 contents of the struct will be
// stored in it.
//
// If the type name of a slice element ends
// with "SET" then it's treated as if the
// "set" tag was set on it. This can be
// used with nested slices where a struct
// tag cannot be given.
//
// Other ASN.1 types are not supported; if
// it encounters them, Unmarshal returns a
// parse error.
func Unmarshal(b []byte, val interface{}) (rest []byte, err error)

// UnmarshalWithParams allows field parameters to be specified for the top-level
// element. The form of the params is the same as the field tags.

// UnmarshalWithParams allows field
// parameters to be specified for the
// top-level element. The form of the
// params is the same as the field tags.
func UnmarshalWithParams(b []byte, val interface{}, params string) (rest []byte, err error)

// BitString is the structure to use when you want an ASN.1 BIT STRING type. A bit
// string is padded up to the nearest byte in memory and the number of valid bits
// is recorded. Padding bits will be zero.

// BitString is the structure to use when
// you want an ASN.1 BIT STRING type. A bit
// string is padded up to the nearest byte
// in memory and the number of valid bits
// is recorded. Padding bits will be zero.
type BitString struct {
	Bytes     []byte // bits packed into bytes.
	BitLength int    // length in bits.
}

// At returns the bit at the given index. If the index is out of range it returns
// false.

// At returns the bit at the given index.
// If the index is out of range it returns
// false.
func (b BitString) At(i int) int

// RightAlign returns a slice where the padding bits are at the beginning. The
// slice may share memory with the BitString.

// RightAlign returns a slice where the
// padding bits are at the beginning. The
// slice may share memory with the
// BitString.
func (b BitString) RightAlign() []byte

// An Enumerated is represented as a plain int.

// An Enumerated is represented as a plain
// int.
type Enumerated int

// A Flag accepts any data and is set to true if present.

// A Flag accepts any data and is set to
// true if present.
type Flag bool

// An ObjectIdentifier represents an ASN.1 OBJECT IDENTIFIER.

// An ObjectIdentifier represents an ASN.1
// OBJECT IDENTIFIER.
type ObjectIdentifier []int

// Equal reports whether oi and other represent the same identifier.

// Equal reports whether oi and other
// represent the same identifier.
func (oi ObjectIdentifier) Equal(other ObjectIdentifier) bool

func (oi ObjectIdentifier) String() string

// RawContent is used to signal that the undecoded, DER data needs to be preserved
// for a struct. To use it, the first field of the struct must have this type. It's
// an error for any of the other fields to have this type.

// RawContent is used to signal that the
// undecoded, DER data needs to be
// preserved for a struct. To use it, the
// first field of the struct must have this
// type. It's an error for any of the other
// fields to have this type.
type RawContent []byte

// A RawValue represents an undecoded ASN.1 object.

// A RawValue represents an undecoded ASN.1
// object.
type RawValue struct {
	Class, Tag int
	IsCompound bool
	Bytes      []byte
	FullBytes  []byte // includes the tag and length
}

// A StructuralError suggests that the ASN.1 data is valid, but the Go type which
// is receiving it doesn't match.

// A StructuralError suggests that the
// ASN.1 data is valid, but the Go type
// which is receiving it doesn't match.
type StructuralError struct {
	Msg string
}

func (e StructuralError) Error() string

// A SyntaxError suggests that the ASN.1 data is invalid.

// A SyntaxError suggests that the ASN.1
// data is invalid.
type SyntaxError struct {
	Msg string
}

func (e SyntaxError) Error() string
