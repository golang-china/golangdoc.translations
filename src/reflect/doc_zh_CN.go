// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package reflect implements run-time reflection, allowing a program to manipulate
// objects with arbitrary types. The typical use is to take a value with static
// type interface{} and extract its dynamic type information by calling TypeOf,
// which returns a Type.
//
// A call to ValueOf returns a Value representing the run-time data. Zero takes a
// Type and returns a Value representing a zero value for that type.
//
// See "The Laws of Reflection" for an introduction to reflection in Go:
// http://golang.org/doc/articles/laws_of_reflection.html

// Package reflect implements run-time
// reflection, allowing a program to
// manipulate objects with arbitrary types.
// The typical use is to take a value with
// static type interface{} and extract its
// dynamic type information by calling
// TypeOf, which returns a Type.
//
// A call to ValueOf returns a Value
// representing the run-time data. Zero
// takes a Type and returns a Value
// representing a zero value for that type.
//
// See "The Laws of Reflection" for an
// introduction to reflection in Go:
// http://golang.org/doc/articles/laws_of_reflection.html
package reflect

// Copy copies the contents of src into dst until either dst has been filled or src
// has been exhausted. It returns the number of elements copied. Dst and src each
// must have kind Slice or Array, and dst and src must have the same element type.

// Copy copies the contents of src into dst
// until either dst has been filled or src
// has been exhausted. It returns the
// number of elements copied. Dst and src
// each must have kind Slice or Array, and
// dst and src must have the same element
// type.
func Copy(dst, src Value) int

// DeepEqual tests for deep equality. It uses normal == equality where possible but
// will scan elements of arrays, slices, maps, and fields of structs. In maps, keys
// are compared with == but elements use deep equality. DeepEqual correctly handles
// recursive types. Functions are equal only if they are both nil. An empty slice
// is not equal to a nil slice.

// DeepEqual tests for deep equality. It
// uses normal == equality where possible
// but will scan elements of arrays,
// slices, maps, and fields of structs. In
// maps, keys are compared with == but
// elements use deep equality. DeepEqual
// correctly handles recursive types.
// Functions are equal only if they are
// both nil. An empty slice is not equal to
// a nil slice.
func DeepEqual(a1, a2 interface{}) bool

// Select executes a select operation described by the list of cases. Like the Go
// select statement, it blocks until at least one of the cases can proceed, makes a
// uniform pseudo-random choice, and then executes that case. It returns the index
// of the chosen case and, if that case was a receive operation, the value received
// and a boolean indicating whether the value corresponds to a send on the channel
// (as opposed to a zero value received because the channel is closed).

// Select executes a select operation
// described by the list of cases. Like the
// Go select statement, it blocks until at
// least one of the cases can proceed,
// makes a uniform pseudo-random choice,
// and then executes that case. It returns
// the index of the chosen case and, if
// that case was a receive operation, the
// value received and a boolean indicating
// whether the value corresponds to a send
// on the channel (as opposed to a zero
// value received because the channel is
// closed).
func Select(cases []SelectCase) (chosen int, recv Value, recvOK bool)

// ChanDir represents a channel type's direction.

// ChanDir represents a channel type's
// direction.
type ChanDir int

const (
	RecvDir ChanDir             = 1 << iota // <-chan
	SendDir                                 // chan<-
	BothDir = RecvDir | SendDir             // chan
)

func (d ChanDir) String() string

// A Kind represents the specific kind of type that a Type represents. The zero
// Kind is not a valid kind.

// A Kind represents the specific kind of
// type that a Type represents. The zero
// Kind is not a valid kind.
type Kind uint

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)

func (k Kind) String() string

// Method represents a single method.

// Method represents a single method.
type Method struct {
	// Name is the method name.
	// PkgPath is the package path that qualifies a lower case (unexported)
	// method name.  It is empty for upper case (exported) method names.
	// The combination of PkgPath and Name uniquely identifies a method
	// in a method set.
	// See http://golang.org/ref/spec#Uniqueness_of_identifiers
	Name    string
	PkgPath string

	Type  Type  // method type
	Func  Value // func with receiver as first argument
	Index int   // index for Type.Method
}

// A SelectCase describes a single case in a select operation. The kind of case
// depends on Dir, the communication direction.
//
// If Dir is SelectDefault, the case represents a default case. Chan and Send must
// be zero Values.
//
// If Dir is SelectSend, the case represents a send operation. Normally Chan's
// underlying value must be a channel, and Send's underlying value must be
// assignable to the channel's element type. As a special case, if Chan is a zero
// Value, then the case is ignored, and the field Send will also be ignored and may
// be either zero or non-zero.
//
// If Dir is SelectRecv, the case represents a receive operation. Normally Chan's
// underlying value must be a channel and Send must be a zero Value. If Chan is a
// zero Value, then the case is ignored, but Send must still be a zero Value. When
// a receive operation is selected, the received Value is returned by Select.

// A SelectCase describes a single case in
// a select operation. The kind of case
// depends on Dir, the communication
// direction.
//
// If Dir is SelectDefault, the case
// represents a default case. Chan and Send
// must be zero Values.
//
// If Dir is SelectSend, the case
// represents a send operation. Normally
// Chan's underlying value must be a
// channel, and Send's underlying value
// must be assignable to the channel's
// element type. As a special case, if Chan
// is a zero Value, then the case is
// ignored, and the field Send will also be
// ignored and may be either zero or
// non-zero.
//
// If Dir is SelectRecv, the case
// represents a receive operation. Normally
// Chan's underlying value must be a
// channel and Send must be a zero Value.
// If Chan is a zero Value, then the case
// is ignored, but Send must still be a
// zero Value. When a receive operation is
// selected, the received Value is returned
// by Select.
type SelectCase struct {
	Dir  SelectDir // direction of case
	Chan Value     // channel to use (for send or receive)
	Send Value     // value to send (for send)
}

// A SelectDir describes the communication direction of a select case.

// A SelectDir describes the communication
// direction of a select case.
type SelectDir int

const (
	_             SelectDir = iota
	SelectSend              // case Chan <- Send
	SelectRecv              // case <-Chan:
	SelectDefault           // default
)

// SliceHeader is the runtime representation of a slice. It cannot be used safely
// or portably and its representation may change in a later release. Moreover, the
// Data field is not sufficient to guarantee the data it references will not be
// garbage collected, so programs must keep a separate, correctly typed pointer to
// the underlying data.

// SliceHeader is the runtime
// representation of a slice. It cannot be
// used safely or portably and its
// representation may change in a later
// release. Moreover, the Data field is not
// sufficient to guarantee the data it
// references will not be garbage
// collected, so programs must keep a
// separate, correctly typed pointer to the
// underlying data.
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// StringHeader is the runtime representation of a string. It cannot be used safely
// or portably and its representation may change in a later release. Moreover, the
// Data field is not sufficient to guarantee the data it references will not be
// garbage collected, so programs must keep a separate, correctly typed pointer to
// the underlying data.

// StringHeader is the runtime
// representation of a string. It cannot be
// used safely or portably and its
// representation may change in a later
// release. Moreover, the Data field is not
// sufficient to guarantee the data it
// references will not be garbage
// collected, so programs must keep a
// separate, correctly typed pointer to the
// underlying data.
type StringHeader struct {
	Data uintptr
	Len  int
}

// A StructField describes a single field in a struct.

// A StructField describes a single field
// in a struct.
type StructField struct {
	// Name is the field name.
	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name.  It is empty for upper case (exported) field names.
	// See http://golang.org/ref/spec#Uniqueness_of_identifiers
	Name    string
	PkgPath string

	Type      Type      // field type
	Tag       StructTag // field tag string
	Offset    uintptr   // offset within struct, in bytes
	Index     []int     // index sequence for Type.FieldByIndex
	Anonymous bool      // is an embedded field
}

// A StructTag is the tag string in a struct field.
//
// By convention, tag strings are a concatenation of optionally space-separated
// key:"value" pairs. Each key is a non-empty string consisting of non-control
// characters other than space (U+0020 ' '), quote (U+0022 '"'), and colon (U+003A
// ':'). Each value is quoted using U+0022 '"' characters and Go string literal
// syntax.

// A StructTag is the tag string in a
// struct field.
//
// By convention, tag strings are a
// concatenation of optionally
// space-separated key:"value" pairs. Each
// key is a non-empty string consisting of
// non-control characters other than space
// (U+0020 ' '), quote (U+0022 '"'), and
// colon (U+003A ':'). Each value is quoted
// using U+0022 '"' characters and Go
// string literal syntax.
type StructTag string

// Get returns the value associated with key in the tag string. If there is no such
// key in the tag, Get returns the empty string. If the tag does not have the
// conventional format, the value returned by Get is unspecified.

// Get returns the value associated with
// key in the tag string. If there is no
// such key in the tag, Get returns the
// empty string. If the tag does not have
// the conventional format, the value
// returned by Get is unspecified.
func (tag StructTag) Get(key string) string

// Type is the representation of a Go type.
//
// Not all methods apply to all kinds of types. Restrictions, if any, are noted in
// the documentation for each method. Use the Kind method to find out the kind of
// type before calling kind-specific methods. Calling a method inappropriate to the
// kind of type causes a run-time panic.

// Type is the representation of a Go type.
//
// Not all methods apply to all kinds of
// types. Restrictions, if any, are noted
// in the documentation for each method.
// Use the Kind method to find out the kind
// of type before calling kind-specific
// methods. Calling a method inappropriate
// to the kind of type causes a run-time
// panic.
type Type interface {

	// Align returns the alignment in bytes of a value of
	// this type when allocated in memory.
	Align() int

	// FieldAlign returns the alignment in bytes of a value of
	// this type when used as a field in a struct.
	FieldAlign() int

	// Method returns the i'th method in the type's method set.
	// It panics if i is not in the range [0, NumMethod()).
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	Method(int) Method

	// MethodByName returns the method with that name in the type's
	// method set and a boolean indicating if the method was found.
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	MethodByName(string) (Method, bool)

	// NumMethod returns the number of methods in the type's method set.
	NumMethod() int

	// Name returns the type's name within its package.
	// It returns an empty string for unnamed types.
	Name() string

	// PkgPath returns a named type's package path, that is, the import path
	// that uniquely identifies the package, such as "encoding/base64".
	// If the type was predeclared (string, error) or unnamed (*T, struct{}, []int),
	// the package path will be the empty string.
	PkgPath() string

	// Size returns the number of bytes needed to store
	// a value of the given type; it is analogous to unsafe.Sizeof.
	Size() uintptr

	// String returns a string representation of the type.
	// The string representation may use shortened package names
	// (e.g., base64 instead of "encoding/base64") and is not
	// guaranteed to be unique among types.  To test for equality,
	// compare the Types directly.
	String() string

	// Kind returns the specific kind of this type.
	Kind() Kind

	// Implements returns true if the type implements the interface type u.
	Implements(u Type) bool

	// AssignableTo returns true if a value of the type is assignable to type u.
	AssignableTo(u Type) bool

	// ConvertibleTo returns true if a value of the type is convertible to type u.
	ConvertibleTo(u Type) bool

	// Comparable returns true if values of this type are comparable.
	Comparable() bool

	// Bits returns the size of the type in bits.
	// It panics if the type's Kind is not one of the
	// sized or unsized Int, Uint, Float, or Complex kinds.
	Bits() int

	// ChanDir returns a channel type's direction.
	// It panics if the type's Kind is not Chan.
	ChanDir() ChanDir

	// IsVariadic returns true if a function type's final input parameter
	// is a "..." parameter.  If so, t.In(t.NumIn() - 1) returns the parameter's
	// implicit actual type []T.
	//
	// For concreteness, if t represents func(x int, y ... float64), then
	//
	//	t.NumIn() == 2
	//	t.In(0) is the reflect.Type for "int"
	//	t.In(1) is the reflect.Type for "[]float64"
	//	t.IsVariadic() == true
	//
	// IsVariadic panics if the type's Kind is not Func.
	IsVariadic() bool

	// Elem returns a type's element type.
	// It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
	Elem() Type

	// Field returns a struct type's i'th field.
	// It panics if the type's Kind is not Struct.
	// It panics if i is not in the range [0, NumField()).
	Field(i int) StructField

	// FieldByIndex returns the nested field corresponding
	// to the index sequence.  It is equivalent to calling Field
	// successively for each index i.
	// It panics if the type's Kind is not Struct.
	FieldByIndex(index []int) StructField

	// FieldByName returns the struct field with the given name
	// and a boolean indicating if the field was found.
	FieldByName(name string) (StructField, bool)

	// FieldByNameFunc returns the first struct field with a name
	// that satisfies the match function and a boolean indicating if
	// the field was found.
	FieldByNameFunc(match func(string) bool) (StructField, bool)

	// In returns the type of a function type's i'th input parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumIn()).
	In(i int) Type

	// Key returns a map type's key type.
	// It panics if the type's Kind is not Map.
	Key() Type

	// Len returns an array type's length.
	// It panics if the type's Kind is not Array.
	Len() int

	// NumField returns a struct type's field count.
	// It panics if the type's Kind is not Struct.
	NumField() int

	// NumIn returns a function type's input parameter count.
	// It panics if the type's Kind is not Func.
	NumIn() int

	// NumOut returns a function type's output parameter count.
	// It panics if the type's Kind is not Func.
	NumOut() int

	// Out returns the type of a function type's i'th output parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumOut()).
	Out(i int) Type
	// contains filtered or unexported methods
}

// ChanOf returns the channel type with the given direction and element type. For
// example, if t represents int, ChanOf(RecvDir, t) represents <-chan int.
//
// The gc runtime imposes a limit of 64 kB on channel element types. If t's size is
// equal to or exceeds this limit, ChanOf panics.

// ChanOf returns the channel type with the
// given direction and element type. For
// example, if t represents int,
// ChanOf(RecvDir, t) represents <-chan
// int.
//
// The gc runtime imposes a limit of 64 kB
// on channel element types. If t's size is
// equal to or exceeds this limit, ChanOf
// panics.
func ChanOf(dir ChanDir, t Type) Type

// MapOf returns the map type with the given key and element types. For example, if
// k represents int and e represents string, MapOf(k, e) represents map[int]string.
//
// If the key type is not a valid map key type (that is, if it does not implement
// Go's == operator), MapOf panics.

// MapOf returns the map type with the
// given key and element types. For
// example, if k represents int and e
// represents string, MapOf(k, e)
// represents map[int]string.
//
// If the key type is not a valid map key
// type (that is, if it does not implement
// Go's == operator), MapOf panics.
func MapOf(key, elem Type) Type

// PtrTo returns the pointer type with element t. For example, if t represents type
// Foo, PtrTo(t) represents *Foo.

// PtrTo returns the pointer type with
// element t. For example, if t represents
// type Foo, PtrTo(t) represents *Foo.
func PtrTo(t Type) Type

// SliceOf returns the slice type with element type t. For example, if t represents
// int, SliceOf(t) represents []int.

// SliceOf returns the slice type with
// element type t. For example, if t
// represents int, SliceOf(t) represents
// []int.
func SliceOf(t Type) Type

// TypeOf returns the reflection Type of the value in the interface{}. TypeOf(nil)
// returns nil.

// TypeOf returns the reflection Type of
// the value in the interface{}.
// TypeOf(nil) returns nil.
func TypeOf(i interface{}) Type

// Value is the reflection interface to a Go value.
//
// Not all methods apply to all kinds of values. Restrictions, if any, are noted in
// the documentation for each method. Use the Kind method to find out the kind of
// value before calling kind-specific methods. Calling a method inappropriate to
// the kind of type causes a run time panic.
//
// The zero Value represents no value. Its IsValid method returns false, its Kind
// method returns Invalid, its String method returns "<invalid Value>", and all
// other methods panic. Most functions and methods never return an invalid value.
// If one does, its documentation states the conditions explicitly.
//
// A Value can be used concurrently by multiple goroutines provided that the
// underlying Go value can be used concurrently for the equivalent direct
// operations.

// Value is the reflection interface to a
// Go value.
//
// Not all methods apply to all kinds of
// values. Restrictions, if any, are noted
// in the documentation for each method.
// Use the Kind method to find out the kind
// of value before calling kind-specific
// methods. Calling a method inappropriate
// to the kind of type causes a run time
// panic.
//
// The zero Value represents no value. Its
// IsValid method returns false, its Kind
// method returns Invalid, its String
// method returns "<invalid Value>", and
// all other methods panic. Most functions
// and methods never return an invalid
// value. If one does, its documentation
// states the conditions explicitly.
//
// A Value can be used concurrently by
// multiple goroutines provided that the
// underlying Go value can be used
// concurrently for the equivalent direct
// operations.
//
// Using == on two Values does not compare
// the underlying values they represent,
// but rather the contents of the Value
// structs. To compare two Values, compare
// the results of the Interface method.
type Value struct {
	// contains filtered or unexported fields
}

// Append appends the values x to a slice s and returns the resulting slice. As in
// Go, each x's value must be assignable to the slice's element type.

// Append appends the values x to a slice s
// and returns the resulting slice. As in
// Go, each x's value must be assignable to
// the slice's element type.
func Append(s Value, x ...Value) Value

// AppendSlice appends a slice t to a slice s and returns the resulting slice. The
// slices s and t must have the same element type.

// AppendSlice appends a slice t to a slice
// s and returns the resulting slice. The
// slices s and t must have the same
// element type.
func AppendSlice(s, t Value) Value

// Indirect returns the value that v points to. If v is a nil pointer, Indirect
// returns a zero Value. If v is not a pointer, Indirect returns v.

// Indirect returns the value that v points
// to. If v is a nil pointer, Indirect
// returns a zero Value. If v is not a
// pointer, Indirect returns v.
func Indirect(v Value) Value

// MakeChan creates a new channel with the specified type and buffer size.

// MakeChan creates a new channel with the
// specified type and buffer size.
func MakeChan(typ Type, buffer int) Value

// MakeFunc returns a new function of the given Type that wraps the function fn.
// When called, that new function does the following:
//
//	- converts its arguments to a slice of Values.
//	- runs results := fn(args).
//	- returns the results as a slice of Values, one per formal result.
//
// The implementation fn can assume that the argument Value slice has the number
// and type of arguments given by typ. If typ describes a variadic function, the
// final Value is itself a slice representing the variadic arguments, as in the
// body of a variadic function. The result Value slice returned by fn must have the
// number and type of results given by typ.
//
// The Value.Call method allows the caller to invoke a typed function in terms of
// Values; in contrast, MakeFunc allows the caller to implement a typed function in
// terms of Values.
//
// The Examples section of the documentation includes an illustration of how to use
// MakeFunc to build a swap function for different types.

// MakeFunc returns a new function of the
// given Type that wraps the function fn.
// When called, that new function does the
// following:
//
//	- converts its arguments to a slice of Values.
//	- runs results := fn(args).
//	- returns the results as a slice of Values, one per formal result.
//
// The implementation fn can assume that
// the argument Value slice has the number
// and type of arguments given by typ. If
// typ describes a variadic function, the
// final Value is itself a slice
// representing the variadic arguments, as
// in the body of a variadic function. The
// result Value slice returned by fn must
// have the number and type of results
// given by typ.
//
// The Value.Call method allows the caller
// to invoke a typed function in terms of
// Values; in contrast, MakeFunc allows the
// caller to implement a typed function in
// terms of Values.
//
// The Examples section of the
// documentation includes an illustration
// of how to use MakeFunc to build a swap
// function for different types.
func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value

// MakeMap creates a new map of the specified type.

// MakeMap creates a new map of the
// specified type.
func MakeMap(typ Type) Value

// MakeSlice creates a new zero-initialized slice value for the specified slice
// type, length, and capacity.

// MakeSlice creates a new zero-initialized
// slice value for the specified slice
// type, length, and capacity.
func MakeSlice(typ Type, len, cap int) Value

// New returns a Value representing a pointer to a new zero value for the specified
// type. That is, the returned Value's Type is PtrTo(typ).

// New returns a Value representing a
// pointer to a new zero value for the
// specified type. That is, the returned
// Value's Type is PtrTo(typ).
func New(typ Type) Value

// NewAt returns a Value representing a pointer to a value of the specified type,
// using p as that pointer.

// NewAt returns a Value representing a
// pointer to a value of the specified
// type, using p as that pointer.
func NewAt(typ Type, p unsafe.Pointer) Value

// ValueOf returns a new Value initialized to the concrete value stored in the
// interface i. ValueOf(nil) returns the zero Value.

// ValueOf returns a new Value initialized
// to the concrete value stored in the
// interface i. ValueOf(nil) returns the
// zero Value.
func ValueOf(i interface{}) Value

// Zero returns a Value representing the zero value for the specified type. The
// result is different from the zero value of the Value struct, which represents no
// value at all. For example, Zero(TypeOf(42)) returns a Value with Kind Int and
// value 0. The returned value is neither addressable nor settable.

// Zero returns a Value representing the
// zero value for the specified type. The
// result is different from the zero value
// of the Value struct, which represents no
// value at all. For example,
// Zero(TypeOf(42)) returns a Value with
// Kind Int and value 0. The returned value
// is neither addressable nor settable.
func Zero(typ Type) Value

// Addr returns a pointer value representing the address of v. It panics if
// CanAddr() returns false. Addr is typically used to obtain a pointer to a struct
// field or slice element in order to call a method that requires a pointer
// receiver.

// Addr returns a pointer value
// representing the address of v. It panics
// if CanAddr() returns false. Addr is
// typically used to obtain a pointer to a
// struct field or slice element in order
// to call a method that requires a pointer
// receiver.
func (v Value) Addr() Value

// Bool returns v's underlying value. It panics if v's kind is not Bool.

// Bool returns v's underlying value. It
// panics if v's kind is not Bool.
func (v Value) Bool() bool

// Bytes returns v's underlying value. It panics if v's underlying value is not a
// slice of bytes.

// Bytes returns v's underlying value. It
// panics if v's underlying value is not a
// slice of bytes.
func (v Value) Bytes() []byte

// Call calls the function v with the input arguments in. For example, if len(in)
// == 3, v.Call(in) represents the Go call v(in[0], in[1], in[2]). Call panics if
// v's Kind is not Func. It returns the output results as Values. As in Go, each
// input argument must be assignable to the type of the function's corresponding
// input parameter. If v is a variadic function, Call creates the variadic slice
// parameter itself, copying in the corresponding values.

// Call calls the function v with the input
// arguments in. For example, if len(in) ==
// 3, v.Call(in) represents the Go call
// v(in[0], in[1], in[2]). Call panics if
// v's Kind is not Func. It returns the
// output results as Values. As in Go, each
// input argument must be assignable to the
// type of the function's corresponding
// input parameter. If v is a variadic
// function, Call creates the variadic
// slice parameter itself, copying in the
// corresponding values.
func (v Value) Call(in []Value) []Value

// CallSlice calls the variadic function v with the input arguments in, assigning
// the slice in[len(in)-1] to v's final variadic argument. For example, if len(in)
// == 3, v.Call(in) represents the Go call v(in[0], in[1], in[2]...). Call panics
// if v's Kind is not Func or if v is not variadic. It returns the output results
// as Values. As in Go, each input argument must be assignable to the type of the
// function's corresponding input parameter.

// CallSlice calls the variadic function v
// with the input arguments in, assigning
// the slice in[len(in)-1] to v's final
// variadic argument. For example, if
// len(in) == 3, v.Call(in) represents the
// Go call v(in[0], in[1], in[2]...). Call
// panics if v's Kind is not Func or if v
// is not variadic. It returns the output
// results as Values. As in Go, each input
// argument must be assignable to the type
// of the function's corresponding input
// parameter.
func (v Value) CallSlice(in []Value) []Value

// CanAddr returns true if the value's address can be obtained with Addr. Such
// values are called addressable. A value is addressable if it is an element of a
// slice, an element of an addressable array, a field of an addressable struct, or
// the result of dereferencing a pointer. If CanAddr returns false, calling Addr
// will panic.

// CanAddr returns true if the value's
// address can be obtained with Addr. Such
// values are called addressable. A value
// is addressable if it is an element of a
// slice, an element of an addressable
// array, a field of an addressable struct,
// or the result of dereferencing a
// pointer. If CanAddr returns false,
// calling Addr will panic.
func (v Value) CanAddr() bool

// CanInterface returns true if Interface can be used without panicking.

// CanInterface returns true if Interface
// can be used without panicking.
func (v Value) CanInterface() bool

// CanSet returns true if the value of v can be changed. A Value can be changed
// only if it is addressable and was not obtained by the use of unexported struct
// fields. If CanSet returns false, calling Set or any type-specific setter (e.g.,
// SetBool, SetInt64) will panic.

// CanSet returns true if the value of v
// can be changed. A Value can be changed
// only if it is addressable and was not
// obtained by the use of unexported struct
// fields. If CanSet returns false, calling
// Set or any type-specific setter (e.g.,
// SetBool, SetInt64) will panic.
func (v Value) CanSet() bool

// Cap returns v's capacity. It panics if v's Kind is not Array, Chan, or Slice.

// Cap returns v's capacity. It panics if
// v's Kind is not Array, Chan, or Slice.
func (v Value) Cap() int

// Close closes the channel v. It panics if v's Kind is not Chan.

// Close closes the channel v. It panics if
// v's Kind is not Chan.
func (v Value) Close()

// Complex returns v's underlying value, as a complex128. It panics if v's Kind is
// not Complex64 or Complex128

// Complex returns v's underlying value, as
// a complex128. It panics if v's Kind is
// not Complex64 or Complex128
func (v Value) Complex() complex128

// Convert returns the value v converted to type t. If the usual Go conversion
// rules do not allow conversion of the value v to type t, Convert panics.

// Convert returns the value v converted to
// type t. If the usual Go conversion rules
// do not allow conversion of the value v
// to type t, Convert panics.
func (v Value) Convert(t Type) Value

// Elem returns the value that the interface v contains or that the pointer v
// points to. It panics if v's Kind is not Interface or Ptr. It returns the zero
// Value if v is nil.

// Elem returns the value that the
// interface v contains or that the pointer
// v points to. It panics if v's Kind is
// not Interface or Ptr. It returns the
// zero Value if v is nil.
func (v Value) Elem() Value

// Field returns the i'th field of the struct v. It panics if v's Kind is not
// Struct or i is out of range.

// Field returns the i'th field of the
// struct v. It panics if v's Kind is not
// Struct or i is out of range.
func (v Value) Field(i int) Value

// FieldByIndex returns the nested field corresponding to index. It panics if v's
// Kind is not struct.

// FieldByIndex returns the nested field
// corresponding to index. It panics if v's
// Kind is not struct.
func (v Value) FieldByIndex(index []int) Value

// FieldByName returns the struct field with the given name. It returns the zero
// Value if no field was found. It panics if v's Kind is not struct.

// FieldByName returns the struct field
// with the given name. It returns the zero
// Value if no field was found. It panics
// if v's Kind is not struct.
func (v Value) FieldByName(name string) Value

// FieldByNameFunc returns the struct field with a name that satisfies the match
// function. It panics if v's Kind is not struct. It returns the zero Value if no
// field was found.

// FieldByNameFunc returns the struct field
// with a name that satisfies the match
// function. It panics if v's Kind is not
// struct. It returns the zero Value if no
// field was found.
func (v Value) FieldByNameFunc(match func(string) bool) Value

// Float returns v's underlying value, as a float64. It panics if v's Kind is not
// Float32 or Float64

// Float returns v's underlying value, as a
// float64. It panics if v's Kind is not
// Float32 or Float64
func (v Value) Float() float64

// Index returns v's i'th element. It panics if v's Kind is not Array, Slice, or
// String or i is out of range.

// Index returns v's i'th element. It
// panics if v's Kind is not Array, Slice,
// or String or i is out of range.
func (v Value) Index(i int) Value

// Int returns v's underlying value, as an int64. It panics if v's Kind is not Int,
// Int8, Int16, Int32, or Int64.

// Int returns v's underlying value, as an
// int64. It panics if v's Kind is not Int,
// Int8, Int16, Int32, or Int64.
func (v Value) Int() int64

// Interface returns v's current value as an interface{}. It is equivalent to:
//
//	var i interface{} = (v's underlying value)
//
// It panics if the Value was obtained by accessing unexported struct fields.

// Interface returns v's current value as
// an interface{}. It is equivalent to:
//
//	var i interface{} = (v's underlying value)
//
// It panics if the Value was obtained by
// accessing unexported struct fields.
func (v Value) Interface() (i interface{})

// InterfaceData returns the interface v's value as a uintptr pair. It panics if
// v's Kind is not Interface.

// InterfaceData returns the interface v's
// value as a uintptr pair. It panics if
// v's Kind is not Interface.
func (v Value) InterfaceData() [2]uintptr

// IsNil reports whether its argument v is nil. The argument must be a chan, func,
// interface, map, pointer, or slice value; if it is not, IsNil panics. Note that
// IsNil is not always equivalent to a regular comparison with nil in Go. For
// example, if v was created by calling ValueOf with an uninitialized interface
// variable i, i==nil will be true but v.IsNil will panic as v will be the zero
// Value.

// IsNil reports whether its argument v is
// nil. The argument must be a chan, func,
// interface, map, pointer, or slice value;
// if it is not, IsNil panics. Note that
// IsNil is not always equivalent to a
// regular comparison with nil in Go. For
// example, if v was created by calling
// ValueOf with an uninitialized interface
// variable i, i==nil will be true but
// v.IsNil will panic as v will be the zero
// Value.
func (v Value) IsNil() bool

// IsValid returns true if v represents a value. It returns false if v is the zero
// Value. If IsValid returns false, all other methods except String panic. Most
// functions and methods never return an invalid value. If one does, its
// documentation states the conditions explicitly.

// IsValid returns true if v represents a
// value. It returns false if v is the zero
// Value. If IsValid returns false, all
// other methods except String panic. Most
// functions and methods never return an
// invalid value. If one does, its
// documentation states the conditions
// explicitly.
func (v Value) IsValid() bool

// Kind returns v's Kind. If v is the zero Value (IsValid returns false), Kind
// returns Invalid.

// Kind returns v's Kind. If v is the zero
// Value (IsValid returns false), Kind
// returns Invalid.
func (v Value) Kind() Kind

// Len returns v's length. It panics if v's Kind is not Array, Chan, Map, Slice, or
// String.

// Len returns v's length. It panics if v's
// Kind is not Array, Chan, Map, Slice, or
// String.
func (v Value) Len() int

// MapIndex returns the value associated with key in the map v. It panics if v's
// Kind is not Map. It returns the zero Value if key is not found in the map or if
// v represents a nil map. As in Go, the key's value must be assignable to the
// map's key type.

// MapIndex returns the value associated
// with key in the map v. It panics if v's
// Kind is not Map. It returns the zero
// Value if key is not found in the map or
// if v represents a nil map. As in Go, the
// key's value must be assignable to the
// map's key type.
func (v Value) MapIndex(key Value) Value

// MapKeys returns a slice containing all the keys present in the map, in
// unspecified order. It panics if v's Kind is not Map. It returns an empty slice
// if v represents a nil map.

// MapKeys returns a slice containing all
// the keys present in the map, in
// unspecified order. It panics if v's Kind
// is not Map. It returns an empty slice if
// v represents a nil map.
func (v Value) MapKeys() []Value

// Method returns a function value corresponding to v's i'th method. The arguments
// to a Call on the returned function should not include a receiver; the returned
// function will always use v as the receiver. Method panics if i is out of range
// or if v is a nil interface value.

// Method returns a function value
// corresponding to v's i'th method. The
// arguments to a Call on the returned
// function should not include a receiver;
// the returned function will always use v
// as the receiver. Method panics if i is
// out of range or if v is a nil interface
// value.
func (v Value) Method(i int) Value

// MethodByName returns a function value corresponding to the method of v with the
// given name. The arguments to a Call on the returned function should not include
// a receiver; the returned function will always use v as the receiver. It returns
// the zero Value if no method was found.

// MethodByName returns a function value
// corresponding to the method of v with
// the given name. The arguments to a Call
// on the returned function should not
// include a receiver; the returned
// function will always use v as the
// receiver. It returns the zero Value if
// no method was found.
func (v Value) MethodByName(name string) Value

// NumField returns the number of fields in the struct v. It panics if v's Kind is
// not Struct.

// NumField returns the number of fields in
// the struct v. It panics if v's Kind is
// not Struct.
func (v Value) NumField() int

// NumMethod returns the number of methods in the value's method set.

// NumMethod returns the number of methods
// in the value's method set.
func (v Value) NumMethod() int

// OverflowComplex returns true if the complex128 x cannot be represented by v's
// type. It panics if v's Kind is not Complex64 or Complex128.

// OverflowComplex returns true if the
// complex128 x cannot be represented by
// v's type. It panics if v's Kind is not
// Complex64 or Complex128.
func (v Value) OverflowComplex(x complex128) bool

// OverflowFloat returns true if the float64 x cannot be represented by v's type.
// It panics if v's Kind is not Float32 or Float64.

// OverflowFloat returns true if the
// float64 x cannot be represented by v's
// type. It panics if v's Kind is not
// Float32 or Float64.
func (v Value) OverflowFloat(x float64) bool

// OverflowInt returns true if the int64 x cannot be represented by v's type. It
// panics if v's Kind is not Int, Int8, int16, Int32, or Int64.

// OverflowInt returns true if the int64 x
// cannot be represented by v's type. It
// panics if v's Kind is not Int, Int8,
// int16, Int32, or Int64.
func (v Value) OverflowInt(x int64) bool

// OverflowUint returns true if the uint64 x cannot be represented by v's type. It
// panics if v's Kind is not Uint, Uintptr, Uint8, Uint16, Uint32, or Uint64.

// OverflowUint returns true if the uint64
// x cannot be represented by v's type. It
// panics if v's Kind is not Uint, Uintptr,
// Uint8, Uint16, Uint32, or Uint64.
func (v Value) OverflowUint(x uint64) bool

// Pointer returns v's value as a uintptr. It returns uintptr instead of
// unsafe.Pointer so that code using reflect cannot obtain unsafe.Pointers without
// importing the unsafe package explicitly. It panics if v's Kind is not Chan,
// Func, Map, Ptr, Slice, or UnsafePointer.
//
// If v's Kind is Func, the returned pointer is an underlying code pointer, but not
// necessarily enough to identify a single function uniquely. The only guarantee is
// that the result is zero if and only if v is a nil func Value.
//
// If v's Kind is Slice, the returned pointer is to the first element of the slice.
// If the slice is nil the returned value is 0. If the slice is empty but non-nil
// the return value is non-zero.

// Pointer returns v's value as a uintptr.
// It returns uintptr instead of
// unsafe.Pointer so that code using
// reflect cannot obtain unsafe.Pointers
// without importing the unsafe package
// explicitly. It panics if v's Kind is not
// Chan, Func, Map, Ptr, Slice, or
// UnsafePointer.
//
// If v's Kind is Func, the returned
// pointer is an underlying code pointer,
// but not necessarily enough to identify a
// single function uniquely. The only
// guarantee is that the result is zero if
// and only if v is a nil func Value.
//
// If v's Kind is Slice, the returned
// pointer is to the first element of the
// slice. If the slice is nil the returned
// value is 0. If the slice is empty but
// non-nil the return value is non-zero.
func (v Value) Pointer() uintptr

// Recv receives and returns a value from the channel v. It panics if v's Kind is
// not Chan. The receive blocks until a value is ready. The boolean value ok is
// true if the value x corresponds to a send on the channel, false if it is a zero
// value received because the channel is closed.

// Recv receives and returns a value from
// the channel v. It panics if v's Kind is
// not Chan. The receive blocks until a
// value is ready. The boolean value ok is
// true if the value x corresponds to a
// send on the channel, false if it is a
// zero value received because the channel
// is closed.
func (v Value) Recv() (x Value, ok bool)

// Send sends x on the channel v. It panics if v's kind is not Chan or if x's type
// is not the same type as v's element type. As in Go, x's value must be assignable
// to the channel's element type.

// Send sends x on the channel v. It panics
// if v's kind is not Chan or if x's type
// is not the same type as v's element
// type. As in Go, x's value must be
// assignable to the channel's element
// type.
func (v Value) Send(x Value)

// Set assigns x to the value v. It panics if CanSet returns false. As in Go, x's
// value must be assignable to v's type.

// Set assigns x to the value v. It panics
// if CanSet returns false. As in Go, x's
// value must be assignable to v's type.
func (v Value) Set(x Value)

// SetBool sets v's underlying value. It panics if v's Kind is not Bool or if
// CanSet() is false.

// SetBool sets v's underlying value. It
// panics if v's Kind is not Bool or if
// CanSet() is false.
func (v Value) SetBool(x bool)

// SetBytes sets v's underlying value. It panics if v's underlying value is not a
// slice of bytes.

// SetBytes sets v's underlying value. It
// panics if v's underlying value is not a
// slice of bytes.
func (v Value) SetBytes(x []byte)

// SetCap sets v's capacity to n. It panics if v's Kind is not Slice or if n is
// smaller than the length or greater than the capacity of the slice.

// SetCap sets v's capacity to n. It panics
// if v's Kind is not Slice or if n is
// smaller than the length or greater than
// the capacity of the slice.
func (v Value) SetCap(n int)

// SetComplex sets v's underlying value to x. It panics if v's Kind is not
// Complex64 or Complex128, or if CanSet() is false.

// SetComplex sets v's underlying value to
// x. It panics if v's Kind is not
// Complex64 or Complex128, or if CanSet()
// is false.
func (v Value) SetComplex(x complex128)

// SetFloat sets v's underlying value to x. It panics if v's Kind is not Float32 or
// Float64, or if CanSet() is false.

// SetFloat sets v's underlying value to x.
// It panics if v's Kind is not Float32 or
// Float64, or if CanSet() is false.
func (v Value) SetFloat(x float64)

// SetInt sets v's underlying value to x. It panics if v's Kind is not Int, Int8,
// Int16, Int32, or Int64, or if CanSet() is false.

// SetInt sets v's underlying value to x.
// It panics if v's Kind is not Int, Int8,
// Int16, Int32, or Int64, or if CanSet()
// is false.
func (v Value) SetInt(x int64)

// SetLen sets v's length to n. It panics if v's Kind is not Slice or if n is
// negative or greater than the capacity of the slice.

// SetLen sets v's length to n. It panics
// if v's Kind is not Slice or if n is
// negative or greater than the capacity of
// the slice.
func (v Value) SetLen(n int)

// SetMapIndex sets the value associated with key in the map v to val. It panics if
// v's Kind is not Map. If val is the zero Value, SetMapIndex deletes the key from
// the map. Otherwise if v holds a nil map, SetMapIndex will panic. As in Go, key's
// value must be assignable to the map's key type, and val's value must be
// assignable to the map's value type.

// SetMapIndex sets the value associated
// with key in the map v to val. It panics
// if v's Kind is not Map. If val is the
// zero Value, SetMapIndex deletes the key
// from the map. Otherwise if v holds a nil
// map, SetMapIndex will panic. As in Go,
// key's value must be assignable to the
// map's key type, and val's value must be
// assignable to the map's value type.
func (v Value) SetMapIndex(key, val Value)

// SetPointer sets the unsafe.Pointer value v to x. It panics if v's Kind is not
// UnsafePointer.

// SetPointer sets the unsafe.Pointer value
// v to x. It panics if v's Kind is not
// UnsafePointer.
func (v Value) SetPointer(x unsafe.Pointer)

// SetString sets v's underlying value to x. It panics if v's Kind is not String or
// if CanSet() is false.

// SetString sets v's underlying value to
// x. It panics if v's Kind is not String
// or if CanSet() is false.
func (v Value) SetString(x string)

// SetUint sets v's underlying value to x. It panics if v's Kind is not Uint,
// Uintptr, Uint8, Uint16, Uint32, or Uint64, or if CanSet() is false.

// SetUint sets v's underlying value to x.
// It panics if v's Kind is not Uint,
// Uintptr, Uint8, Uint16, Uint32, or
// Uint64, or if CanSet() is false.
func (v Value) SetUint(x uint64)

// Slice returns v[i:j]. It panics if v's Kind is not Array, Slice or String, or if
// v is an unaddressable array, or if the indexes are out of bounds.

// Slice returns v[i:j]. It panics if v's
// Kind is not Array, Slice or String, or
// if v is an unaddressable array, or if
// the indexes are out of bounds.
func (v Value) Slice(i, j int) Value

// Slice3 is the 3-index form of the slice operation: it returns v[i:j:k]. It
// panics if v's Kind is not Array or Slice, or if v is an unaddressable array, or
// if the indexes are out of bounds.

// Slice3 is the 3-index form of the slice
// operation: it returns v[i:j:k]. It
// panics if v's Kind is not Array or
// Slice, or if v is an unaddressable
// array, or if the indexes are out of
// bounds.
func (v Value) Slice3(i, j, k int) Value

// String returns the string v's underlying value, as a string. String is a special
// case because of Go's String method convention. Unlike the other getters, it does
// not panic if v's Kind is not String. Instead, it returns a string of the form
// "<T value>" where T is v's type.

// String returns the string v's underlying
// value, as a string. String is a special
// case because of Go's String method
// convention. Unlike the other getters, it
// does not panic if v's Kind is not
// String. Instead, it returns a string of
// the form "<T value>" where T is v's
// type.
func (v Value) String() string

// TryRecv attempts to receive a value from the channel v but will not block. It
// panics if v's Kind is not Chan. If the receive delivers a value, x is the
// transferred value and ok is true. If the receive cannot finish without blocking,
// x is the zero Value and ok is false. If the channel is closed, x is the zero
// value for the channel's element type and ok is false.

// TryRecv attempts to receive a value from
// the channel v but will not block. It
// panics if v's Kind is not Chan. If the
// receive delivers a value, x is the
// transferred value and ok is true. If the
// receive cannot finish without blocking,
// x is the zero Value and ok is false. If
// the channel is closed, x is the zero
// value for the channel's element type and
// ok is false.
func (v Value) TryRecv() (x Value, ok bool)

// TrySend attempts to send x on the channel v but will not block. It panics if v's
// Kind is not Chan. It returns true if the value was sent, false otherwise. As in
// Go, x's value must be assignable to the channel's element type.

// TrySend attempts to send x on the
// channel v but will not block. It panics
// if v's Kind is not Chan. It returns true
// if the value was sent, false otherwise.
// As in Go, x's value must be assignable
// to the channel's element type.
func (v Value) TrySend(x Value) bool

// Type returns v's type.

// Type returns v's type.
func (v Value) Type() Type

// Uint returns v's underlying value, as a uint64. It panics if v's Kind is not
// Uint, Uintptr, Uint8, Uint16, Uint32, or Uint64.

// Uint returns v's underlying value, as a
// uint64. It panics if v's Kind is not
// Uint, Uintptr, Uint8, Uint16, Uint32, or
// Uint64.
func (v Value) Uint() uint64

// UnsafeAddr returns a pointer to v's data. It is for advanced clients that also
// import the "unsafe" package. It panics if v is not addressable.

// UnsafeAddr returns a pointer to v's
// data. It is for advanced clients that
// also import the "unsafe" package. It
// panics if v is not addressable.
func (v Value) UnsafeAddr() uintptr

// A ValueError occurs when a Value method is invoked on a Value that does not
// support it. Such cases are documented in the description of each method.

// A ValueError occurs when a Value method
// is invoked on a Value that does not
// support it. Such cases are documented in
// the description of each method.
type ValueError struct {
	Method string
	Kind   Kind
}

func (e *ValueError) Error() string
