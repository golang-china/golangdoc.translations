// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

//     Package builtin provides documentation for Go's predeclared identifiers.
//     The items documented here are not actually in package builtin
//     but their descriptions here allow godoc to present documentation
//     for the language's special identifiers.

// builtin 包为Go的预声明标识符提供了文档.
//     此处列出的条目其实并不在 buildin 包中，对它们的描述只是为了让 godoc
//     给该语言的特殊标识符提供文档。
package builtin

// ComplexType is here for the purposes of documentation only. It is a
// stand-in for either complex type: complex64 or complex128.

// ComplexType 在此只用作文档目的。
// 它代表所有的复数类型：即 complex64 或 complex128。
type ComplexType complex64


// FloatType is here for the purposes of documentation only. It is a stand-in
// for either float type: float32 or float64.

// FloatType 在此只用作文档目的。
// 它代表所有的浮点数类型：即 float32 或 float64。
type FloatType float32


// IntegerType is here for the purposes of documentation only. It is a stand-in
// for any integer type: int, uint, int8 etc.

// IntegerType 在此只用作文档目的。
// 它代表所有的整数类型：如 int、uint、int8 等。
type IntegerType int


// Type is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.

// Type 在此只用作文档目的。
// 它代表所有Go的类型，但对于任何给定的函数请求来说，它都代表与其相同的类型。
type Type int


// Type1 is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.

// Type1 在此只用作文档目的。
// 它代表所有Go的类型，但对于任何给定的函数请求来说，它都代表与其相同的类型。
type Type1 int


