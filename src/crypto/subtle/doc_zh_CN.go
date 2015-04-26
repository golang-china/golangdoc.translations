// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package subtle implements functions that are often useful in cryptographic code
// but require careful thought to use correctly.

// Package subtle implements functions that are often useful in cryptographic
//
//	code but require careful thought to use correctly.
package subtle

// ConstantTimeByteEq returns 1 if x == y and 0 otherwise.

// 如果x == y返回1，否则返回0。
func ConstantTimeByteEq(x, y uint8) int

// ConstantTimeCompare returns 1 iff the two slices, x and y, have equal contents.
// The time taken is a function of the length of the slices and is independent of
// the contents.

// 如果x、y的长度和内容都相同返回1；否则返回0。消耗的时间正比于切片长度而与内容无关。
func ConstantTimeCompare(x, y []byte) int

// ConstantTimeCopy copies the contents of y into x (a slice of equal length) if v
// == 1. If v == 0, x is left unchanged. Its behavior is undefined if v takes any
// other value.

// 如果v == 1,则将y的内容拷贝到x；如果v ==
// 0，x不作修改；其他情况的行为是未定义并应避免的。
func ConstantTimeCopy(v int, x, y []byte)

// ConstantTimeEq returns 1 if x == y and 0 otherwise.

// 如果x == y返回1，否则返回0。
func ConstantTimeEq(x, y int32) int

// ConstantTimeLessOrEq returns 1 if x <= y and 0 otherwise. Its behavior is
// undefined if x or y are negative or > 2**31 - 1.

// 如果x <=
// y返回1，否则返回0；如果x或y为负数，或者大于2**31-1，函数行为是未定义的。
func ConstantTimeLessOrEq(x, y int) int

// ConstantTimeSelect returns x if v is 1 and y if v is 0. Its behavior is
// undefined if v takes any other value.

// 如果v == 1，返回x；如果v ==
// 0，返回y；其他情况的行为是未定义并应避免的。
func ConstantTimeSelect(v, x, y int) int
