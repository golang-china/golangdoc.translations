// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package expvar provides a standardized interface to public variables, such as
// operation counters in servers. It exposes these variables via HTTP at
// /debug/vars in JSON format.
//
// Operations to set or modify these public variables are atomic.
//
// In addition to adding the HTTP handler, this package registers the following
// variables:
//
//	cmdline   os.Args
//	memstats  runtime.Memstats
//
// The package is sometimes only imported for the side effect of registering its
// HTTP handler and the above variables. To use it this way, link this package into
// your program:
//
//	import _ "expvar"
package expvar

// Do calls f for each exported variable. The global variable map is locked during
// the iteration, but existing entries may be concurrently updated.
func Do(f func(KeyValue))

// Publish declares a named exported variable. This should be called from a
// package's init function when it creates its Vars. If the name is already
// registered then this will log.Panic.
func Publish(name string, v Var)

// Float is a 64-bit float variable that satisfies the Var interface.
type Float struct {
	// contains filtered or unexported fields
}

func NewFloat(name string) *Float

// Add adds delta to v.
func (v *Float) Add(delta float64)

// Set sets v to value.
func (v *Float) Set(value float64)

func (v *Float) String() string

// Func implements Var by calling the function and formatting the returned value
// using JSON.
type Func func() interface{}

func (f Func) String() string

// Int is a 64-bit integer variable that satisfies the Var interface.
type Int struct {
	// contains filtered or unexported fields
}

func NewInt(name string) *Int

func (v *Int) Add(delta int64)

func (v *Int) Set(value int64)

func (v *Int) String() string

// KeyValue represents a single entry in a Map.
type KeyValue struct {
	Key   string
	Value Var
}

// Map is a string-to-Var map variable that satisfies the Var interface.
type Map struct {
	// contains filtered or unexported fields
}

func NewMap(name string) *Map

func (v *Map) Add(key string, delta int64)

// AddFloat adds delta to the *Float value stored under the given map key.
func (v *Map) AddFloat(key string, delta float64)

// Do calls f for each entry in the map. The map is locked during the iteration,
// but existing entries may be concurrently updated.
func (v *Map) Do(f func(KeyValue))

func (v *Map) Get(key string) Var

func (v *Map) Init() *Map

func (v *Map) Set(key string, av Var)

func (v *Map) String() string

// String is a string variable, and satisfies the Var interface.
type String struct {
	// contains filtered or unexported fields
}

func NewString(name string) *String

func (v *String) Set(value string)

func (v *String) String() string

// Var is an abstract type for all exported variables.
type Var interface {
	String() string
}

// Get retrieves a named exported variable.
func Get(name string) Var
