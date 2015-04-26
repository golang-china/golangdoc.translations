// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package user allows user account lookups by name or id.

// user包允许通过名称或ID查询用户帐户。
package user

// UnknownUserError is returned by Lookup when a user cannot be found.

// 当找不到用户时，Lookup会返回UnknownUserError
type UnknownUserError string

func (e UnknownUserError) Error() string

// UnknownUserIdError is returned by LookupId when a user cannot be found.

// 当找不到用户时，LookupId会返回UnknownUserIdError
type UnknownUserIdError int

func (e UnknownUserIdError) Error() string

// User represents a user account.
//
// On posix systems Uid and Gid contain a decimal number representing uid and gid.
// On windows Uid and Gid contain security identifier (SID) in a string format. On
// Plan 9, Uid, Gid, Username, and Name will be the contents of /dev/user.

// User代表一个用户帐户。
//
// 在posix系统中Uid和Gid字段分别包含代表uid和gid的十进制数字。在windows系统中Uid和Gid包含字符串格式的安全标识符（SID）。在Plan
// 9系统中，Uid、Gid、Username和Name字段是/dev/user的内容。
type User struct {
	Uid      string // user id
	Gid      string // primary group id
	Username string
	Name     string
	HomeDir  string
}

// Current returns the current user.

// 返回当前的用户帐户。
func Current() (*User, error)

// Lookup looks up a user by username. If the user cannot be found, the returned
// error is of type UnknownUserError.

// 根据用户名查询用户。
func Lookup(username string) (*User, error)

// LookupId looks up a user by userid. If the user cannot be found, the returned
// error is of type UnknownUserIdError.

// 根据用户ID查询用户。
func LookupId(uid string) (*User, error)
