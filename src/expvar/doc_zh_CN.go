// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package expvar provides a standardized interface to public variables, such
// as operation counters in servers. It exposes these variables via HTTP at
// /debug/vars in JSON format.
//
// Operations to set or modify these public variables are atomic.
//
// In addition to adding the HTTP handler, this package registers the
// following variables:
//
//     cmdline   os.Args
//     memstats  runtime.Memstats
//
// The package is sometimes only imported for the side effect of
// registering its HTTP handler and the above variables.  To use it
// this way, link this package into your program:
//     import _ "expvar"

// Package expvar provides a standardized interface to public variables, such
// as operation counters in servers. It exposes these variables via HTTP at
// /debug/vars in JSON format.
//
// Operations to set or modify these public variables are atomic.
//
// In addition to adding the HTTP handler, this package registers the
// following variables:
//
//     cmdline   os.Args
//     memstats  runtime.Memstats
//
// The package is sometimes only imported for the side effect of
// registering its HTTP handler and the above variables. To use it
// this way, link this package into your program:
//     import _ "expvar"
package expvar

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "math"
    "net/http"
    "os"
    "runtime"
    "sort"
    "strconv"
    "sync"
    "sync/atomic"
)

// Float is a 64-bit float variable that satisfies the Var interface.
type Float struct {
	f uint64
}


// Func implements Var by calling the function
// and formatting the returned value using JSON.
type Func func() interface{}


// Int is a 64-bit integer variable that satisfies the Var interface.
type Int struct {
	i int64
}


// KeyValue represents a single entry in a Map.
type KeyValue struct {
	Key   string
	Value Var
}


// Map is a string-to-Var map variable that satisfies the Var interface.
type Map struct {
	mu   sync.RWMutex
	m    map[string]Var
	keys []string // sorted
}


// String is a string variable, and satisfies the Var interface.
type String struct {
	mu sync.RWMutex
	s  string
}


// Var is an abstract type for all exported variables.
type Var interface {
	// String returns a valid JSON value for the variable.
	// Types with String methods that do not return valid JSON
	// (such as time.Time) must not be used as a Var.
	String() string
}


// Do calls f for each exported variable.
// The global variable map is locked during the iteration,
// but existing entries may be concurrently updated.
func Do(f func(KeyValue))

// Get retrieves a named exported variable.

// Get retrieves a named exported variable. It returns nil if the name has
// not been registered.
func Get(name string) Var

func NewFloat(name string) *Float

func NewInt(name string) *Int

func NewMap(name string) *Map

func NewString(name string) *String

// Publish declares a named exported variable. This should be called from a
// package's init function when it creates its Vars. If the name is already
// registered then this will log.Panic.
func Publish(name string, v Var)

// Add adds delta to v.
func (*Float) Add(delta float64)

// Set sets v to value.
func (*Float) Set(value float64)

func (*Float) String() string

func (*Int) Add(delta int64)

func (*Int) Set(value int64)

func (*Int) String() string

func (*Map) Add(key string, delta int64)

// AddFloat adds delta to the *Float value stored under the given map key.
func (*Map) AddFloat(key string, delta float64)

// Do calls f for each entry in the map.
// The map is locked during the iteration,
// but existing entries may be concurrently updated.
func (*Map) Do(f func(KeyValue))

func (*Map) Get(key string) Var

func (*Map) Init() *Map

func (*Map) Set(key string, av Var)

func (*Map) String() string

func (*String) Set(value string)

func (*String) String() string

func (Func) String() string

