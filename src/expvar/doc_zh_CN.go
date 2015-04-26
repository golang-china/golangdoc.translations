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

// expvar包提供了公共变量的标准接口，如服务的操作计数器。本包通过HTTP在/debug/vars位置以JSON格式导出了这些变量。
//
// 对这些公共变量的读写操作都是原子级的。
//
// 为了增加HTTP处理器，本包注册了如下变量：
//
//	cmdline   os.Args
//	memstats  runtime.Memstats
//
// 有时候本包被导入只是为了获得本包注册HTTP处理器和上述变量的副作用。此时可以如下方式导入本包：
//
//	import _ "expvar"
package expvar

// Do calls f for each exported variable. The global variable map is locked during
// the iteration, but existing entries may be concurrently updated.

// Do对导出变量的每一条记录都调用f。迭代执行时会锁定全局变量映射，但已存在的记录可以同时更新。
func Do(f func(KeyValue))

// Publish declares a named exported variable. This should be called from a
// package's init function when it creates its Vars. If the name is already
// registered then this will log.Panic.

// Publish声明一个导出变量。必须在init函数里调用。如果name已经被注册，会调用log.Panic。
func Publish(name string, v Var)

// Float is a 64-bit float variable that satisfies the Var interface.

// Float代表一个64位浮点数变量，满足Var接口。
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

// Func通过调用函数并将结果编码为json，实现了Var接口。
type Func func() interface{}

func (f Func) String() string

// Int is a 64-bit integer variable that satisfies the Var interface.

// Int代表一个64位整数变量，满足Var接口。
type Int struct {
	// contains filtered or unexported fields
}

func NewInt(name string) *Int

func (v *Int) Add(delta int64)

func (v *Int) Set(value int64)

func (v *Int) String() string

// KeyValue represents a single entry in a Map.

// KeyValue代表Map中的一条记录。（键值对）
type KeyValue struct {
	Key   string
	Value Var
}

// Map is a string-to-Var map variable that satisfies the Var interface.

// Map代表一个string到Var的映射变量，满足Var接口。
type Map struct {
	// contains filtered or unexported fields
}

func NewMap(name string) *Map

func (v *Map) Add(key string, delta int64)

// AddFloat adds delta to the *Float value stored under the given map key.

// AddFloat向索引key对应的值（底层为*Float）修改为加上delta后的值。
func (v *Map) AddFloat(key string, delta float64)

// Do calls f for each entry in the map. The map is locked during the iteration,
// but existing entries may be concurrently updated.

// Do对映射的每一条记录都调用f。迭代执行时会锁定该映射，但已存在的记录可以同时更新。
func (v *Map) Do(f func(KeyValue))

func (v *Map) Get(key string) Var

func (v *Map) Init() *Map

func (v *Map) Set(key string, av Var)

func (v *Map) String() string

// String is a string variable, and satisfies the Var interface.

// String代表一个字符串变量，满足Var接口。
type String struct {
	// contains filtered or unexported fields
}

func NewString(name string) *String

func (v *String) Set(value string)

func (v *String) String() string

// Var is an abstract type for all exported variables.

// Var接口是所有导出变量的抽象类型。
type Var interface {
	String() string
}

// Get retrieves a named exported variable.

// Get获取名为name的导出变量。
func Get(name string) Var
