// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package singleflight provides a duplicate function call suppression
// mechanism.

// Package singleflight provides a duplicate function call suppression
// mechanism.
package singleflight // import "internal/singleflight"

import (
    "errors"
    "fmt"
    "sync"
    "sync/atomic"
    "testing"
    "time"
)

// Group represents a class of work and forms a namespace in
// which units of work can be executed with duplicate suppression.
type Group struct {
    mu  sync.Mutex       // protects m
    m   map[string]*call // lazily initialized
}

// Result holds the results of Do, so they can be passed
// on a channel.
type Result struct {
    Val    interface{}
    Err    error
    Shared bool
}

func TestDo(t *testing.T)

func TestDoDupSuppress(t *testing.T)

func TestDoErr(t *testing.T)

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
// The return value shared indicates whether v was given to multiple callers.
func (*Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)

// DoChan is like Do but returns a channel that will receive the
// results when they are ready.
func (*Group) DoChan(key string, fn func() (interface{}, error)) (<-chan Result)

// Forget tells the singleflight to forget about a key.  Future calls
// to Do for this key will call the function rather than waiting for
// an earlier call to complete.
func (*Group) Forget(key string)

