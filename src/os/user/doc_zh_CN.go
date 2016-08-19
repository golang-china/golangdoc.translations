// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package user allows user account lookups by name or id.

// Package user allows user account lookups by name or id.
package user

import (
    "C"
    "errors"
    "fmt"
    "os"
    "runtime"
    "strconv"
    "strings"
    "syscall"
    "unsafe"
)

// Group represents a grouping of users.
//
// On POSIX systems Gid contains a decimal number
// representing the group ID.
type Group struct {
	Gid  string // group ID
	Name string // group name
}


// UnknownGroupError is returned by LookupGroup when
// a group cannot be found.
type UnknownGroupError string


// UnknownGroupIdError is returned by LookupGroupId when
// a group cannot be found.
type UnknownGroupIdError string


// UnknownUserError is returned by Lookup when
// a user cannot be found.
type UnknownUserError string


// UnknownUserIdError is returned by LookupId when
// a user cannot be found.
type UnknownUserIdError int


// User represents a user account.
//
// On posix systems Uid and Gid contain a decimal number
// representing uid and gid. On windows Uid and Gid
// contain security identifier (SID) in a string format.
// On Plan 9, Uid, Gid, Username, and Name will be the
// contents of /dev/user.

// User represents a user account.
//
// On POSIX systems Uid and Gid contain a decimal number
// representing uid and gid. On windows Uid and Gid
// contain security identifier (SID) in a string format.
// On Plan 9, Uid, Gid, Username, and Name will be the
// contents of /dev/user.
type User struct {
	Uid      string // user ID
	Gid      string // primary group ID
	Username string
	Name     string
	HomeDir  string
}


// Current returns the current user.
func Current() (*User, error)

// Lookup looks up a user by username. If the user cannot be found, the
// returned error is of type UnknownUserError.
func Lookup(username string) (*User, error)

// LookupGroup looks up a group by name. If the group cannot be found, the
// returned error is of type UnknownGroupError.
func LookupGroup(name string) (*Group, error)

// LookupGroupId looks up a group by groupid. If the group cannot be found, the
// returned error is of type UnknownGroupIdError.
func LookupGroupId(gid string) (*Group, error)

// LookupId looks up a user by userid. If the user cannot be found, the
// returned error is of type UnknownUserIdError.
func LookupId(uid string) (*User, error)

// GroupIds returns the list of group IDs that the user is a member of.
func (*User) GroupIds() ([]string, error)

func (UnknownGroupError) Error() string

func (UnknownGroupIdError) Error() string

func (UnknownUserError) Error() string

func (UnknownUserIdError) Error() string

